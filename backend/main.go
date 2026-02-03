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

// syncNotifier 同步完成时向所有订阅的 SSE 客户端广播，前端可静默刷新列表
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
		log.Fatalf("初始化数据库失败: %v", err)
	}
	defer db.Close()

	log.Printf("数据库已初始化: %s", cfg.DBPath)

	syncNotifier := &syncNotifier{}
	syncService := services.NewSyncService(db.Conn(), cfg.PostsPath, cfg.IsDev, syncNotifier, cfg.PostsRemoteURL)

	// 配置了远程仓库时先同步完成再监听，确保首屏能加载出文章；否则后台同步
	if cfg.IsDev {
		if cfg.PostsRemoteURL != "" {
			log.Println("开发模式：执行初始同步（含远程 clone），完成后启动...")
			if err := syncService.Sync(); err != nil {
				log.Printf("初始同步失败: %v", err)
			}
		} else {
			go func() {
				log.Println("开发模式：执行初始同步...")
				if err := syncService.Sync(); err != nil {
					log.Printf("初始同步失败: %v", err)
				}
			}()
		}
	}

	if cfg.SyncIntervalMinutes > 0 {
		interval := time.Duration(cfg.SyncIntervalMinutes) * time.Minute
		log.Printf("启用定期同步：每 %d 分钟从远程仓库拉取并更新", cfg.SyncIntervalMinutes)
		go func() {
			ticker := time.NewTicker(interval)
			defer ticker.Stop()
			for range ticker.C {
				if err := syncService.Sync(); err != nil {
					log.Printf("定期同步失败: %v", err)
				}
			}
		}()
	}

	postsHandler := handlers.NewPostsHandler(db.Conn(), cfg.PostsPath)
	webhookHandler := handlers.NewWebhookHandler(syncService, cfg.WebhookSecret)

	r := gin.Default()

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
		api.POST("/webhook", webhookHandler.HandleWebhook)
		api.GET("/posts", postsHandler.GetPosts)
		api.GET("/posts/:slug", postsHandler.GetPost)
		// 文章仓库内图片等静态资源，便于 Markdown 中相对路径正确展示
		api.GET("/posts-assets/*path", postsHandler.ServePostAsset)
		// SSE：同步完成后推送，前端可静默刷新列表而不整页重载
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
	log.Printf("服务器启动在 %s (模式: %s)", addr, cfg.Mode)
	log.Printf("文章目录: %s", cfg.PostsPath)

	if err := r.Run(addr); err != nil {
		log.Fatalf("启动服务器失败: %v", err)
	}
}
