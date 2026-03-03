package main

import (
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"

	"blog-suiseiseki/config"
	"blog-suiseiseki/database"
	"blog-suiseiseki/handlers"
	"blog-suiseiseki/services"
)

// syncNotifier broadcasts to SSE subscribers when sync completes so the frontend can refresh the list.
type syncNotifier struct {
	mu   sync.Mutex
	subs []chan struct{}
}

func (n *syncNotifier) Subscribe() (ch <-chan struct{}, unsubscribe func()) {
	c := make(chan struct{}, 1)
	n.mu.Lock()
	n.subs = append(n.subs, c)
	n.mu.Unlock()
	return c, func() {
		n.mu.Lock()
		defer n.mu.Unlock()
		for i, s := range n.subs {
			if s == c {
				n.subs = append(n.subs[:i], n.subs[i+1:]...)
				close(c)
				break
			}
		}
	}
}

func (n *syncNotifier) Broadcast() {
	n.mu.Lock()
	subs := make([]chan struct{}, len(n.subs))
	copy(subs, n.subs)
	n.mu.Unlock()
	for _, ch := range subs {
		select {
		case ch <- struct{}{}:
		default:
		}
	}
}

func main() {
	cfg := config.Load()

	if cfg.IsDev {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	db, err := database.New(cfg.DBPath)
	if err != nil {
		log.Fatalf("db init failed: %v", err)
	}
	defer db.Close()

	log.Printf("db initialized: %s", cfg.DBPath)

	syncNotifier := &syncNotifier{}
	syncService := services.NewSyncService(db.Conn(), cfg.PostsPath, cfg.IsDev, syncNotifier, cfg.PostsRemoteURL)

	// run an initial synchronization on startup regardless of the mode. 之前
	// 仅在开发环境调用，这会导致生产部署后的第一批文章延迟直到
	// webhook 或定期同步触发。现在启动时就会执行一次，确保服务器
	// 启动后立刻拥有最新内容。
	log.Println("running initial sync...")
	if err := syncService.Sync(); err != nil {
		log.Printf("initial sync failed: %v", err)
	}

	if cfg.SyncIntervalMinutes > 0 {
		interval := time.Duration(cfg.SyncIntervalMinutes) * time.Minute
		log.Printf("sync: interval %d min", cfg.SyncIntervalMinutes)
		go func() {
			ticker := time.NewTicker(interval)
			defer ticker.Stop()
			for range ticker.C {
				if err := syncService.Sync(); err != nil {
					log.Printf("periodic sync failed: %v", err)
				}
			}
		}()
	}

	postsHandler := handlers.NewPostsHandler(db.Conn(), cfg.PostsPath)

	r := gin.Default()
	// Trust only local reverse proxy (Caddy/nginx); avoids "trusted all proxies" warning
	r.SetTrustedProxies([]string{"127.0.0.1", "::1"})

	if cfg.IsDev {
		r.Use(func(c *gin.Context) {
			c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
			c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
			c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Hub-Signature-256")
			if c.Request.Method == "OPTIONS" {
				c.AbortWithStatus(204)
				return
			}
			c.Next()
		})
	}

	api := r.Group("/api")
	{
		api.GET("/posts", postsHandler.GetPosts)
		api.GET("/posts/:slug", postsHandler.GetPost)
		// Static assets from posts repo for relative paths in Markdown
		api.GET("/posts-assets/*path", postsHandler.ServePostAsset)
		// SSE: push on sync completion so frontend can refresh list without full reload
		api.GET("/events", func(c *gin.Context) {
			ch, unsub := syncNotifier.Subscribe()
			defer unsub()
			c.Header("Content-Type", "text/event-stream")
			c.Header("Cache-Control", "no-cache")
			c.Header("Connection", "keep-alive")
			c.Header("X-Accel-Buffering", "no")
			c.Writer.Flush()
			for {
				select {
				case <-ch:
					c.SSEvent("sync_completed", nil)
					c.Writer.Flush()
				case <-c.Request.Context().Done():
					return
				}
			}
		})
	}

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	addr := ":" + cfg.Port
	log.Printf("server listening on %s (mode: %s)", addr, cfg.Mode)
	log.Printf("posts path: %s", cfg.PostsPath)

	if err := r.Run(addr); err != nil {
		log.Fatalf("server run failed: %v", err)
	}
}
