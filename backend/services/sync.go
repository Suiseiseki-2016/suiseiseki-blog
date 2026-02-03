package services

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"blog-suiseiseki/utils"
)

// SyncEventNotifier 同步完成后通知前端（如 SSE），便于静默刷新列表而不整页重载
type SyncEventNotifier interface {
	Broadcast()
}

type SyncService struct {
	db        *sql.DB
	postsPath string
	isDev     bool
	notifier  SyncEventNotifier
}

func NewSyncService(db *sql.DB, postsPath string, isDev bool, notifier SyncEventNotifier) *SyncService {
	return &SyncService{
		db:        db,
		postsPath: postsPath,
		isDev:     isDev,
		notifier:  notifier,
	}
}

// Sync 同步文章：扫描目录，解析 Markdown，更新数据库；生产环境先 git pull
func (s *SyncService) Sync() error {
	log.Println("开始同步文章...")

	if !s.isDev {
		if err := s.gitPull(); err != nil {
			log.Printf("Git pull 失败: %v", err)
		}
	}

	files, err := s.scanMarkdownFiles()
	if err != nil {
		return fmt.Errorf("扫描文件失败: %w", err)
	}

	log.Printf("找到 %d 个Markdown文件", len(files))

	existingPaths, err := s.getExistingPaths()
	if err != nil {
		return fmt.Errorf("获取现有路径失败: %w", err)
	}

	processedPaths := make(map[string]bool)
	for _, filePath := range files {
		processedPaths[filePath] = true
		if err := s.processFile(filePath); err != nil {
			log.Printf("处理文件 %s 失败: %v", filePath, err)
		}
	}

	for path := range existingPaths {
		if !processedPaths[path] {
			if err := s.deletePost(path); err != nil {
				log.Printf("删除文章 %s 失败: %v", path, err)
			}
		}
	}

	log.Println("同步完成")
	if s.notifier != nil {
		s.notifier.Broadcast()
	}
	return nil
}

func (s *SyncService) gitPull() error {
	if _, err := os.Stat(filepath.Join(s.postsPath, ".git")); os.IsNotExist(err) {
		return fmt.Errorf("不是Git仓库: %s", s.postsPath)
	}

	cmd := exec.Command("git", "pull")
	cmd.Dir = s.postsPath
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("git pull 失败: %v, 输出: %s", err, string(output))
	}

	log.Printf("Git pull 成功: %s", string(output))
	return nil
}

func (s *SyncService) scanMarkdownFiles() ([]string, error) {
	var files []string

	err := filepath.Walk(s.postsPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && strings.HasSuffix(path, ".md") {
			base := filepath.Base(path)
			if strings.EqualFold(base, "README.md") {
				return nil
			}
			files = append(files, path)
		}

		return nil
	})

	return files, err
}

func (s *SyncService) getExistingPaths() (map[string]bool, error) {
	rows, err := s.db.Query("SELECT content_path FROM posts")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	paths := make(map[string]bool)
	for rows.Next() {
		var path string
		if err := rows.Scan(&path); err != nil {
			return nil, err
		}
		paths[path] = true
	}

	return paths, nil
}

func (s *SyncService) processFile(filePath string) error {
	fm, _, err := utils.ParseMarkdownFile(filePath)
	if err != nil {
		return fmt.Errorf("解析Markdown失败: %w", err)
	}

	slug := fm.Slug
	if slug == "" {
		slug = utils.GenerateSlug(filePath)
	}

	var publishedAt time.Time
	if fm.PublishedAt != "" {
		formats := []string{
			"2006-01-02",
			"2006-01-02 15:04:05",
			time.RFC3339,
		}
		for _, format := range formats {
			if t, err := time.Parse(format, fm.PublishedAt); err == nil {
				publishedAt = t
				break
			}
		}
	}
	if publishedAt.IsZero() {
		if info, err := os.Stat(filePath); err == nil {
			publishedAt = info.ModTime()
		} else {
			publishedAt = time.Now()
		}
	}

	query := `
		INSERT INTO posts (slug, title, summary, category, published_at, content_path, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, CURRENT_TIMESTAMP)
		ON CONFLICT(slug) DO UPDATE SET
			title = excluded.title,
			summary = excluded.summary,
			category = excluded.category,
			published_at = excluded.published_at,
			content_path = excluded.content_path,
			updated_at = CURRENT_TIMESTAMP
	`

	_, err = s.db.Exec(query, slug, fm.Title, fm.Summary, fm.Category, publishedAt, filePath)
	if err != nil {
		return fmt.Errorf("数据库操作失败: %w", err)
	}

	log.Printf("处理文章: %s (%s)", fm.Title, slug)
	return nil
}

func (s *SyncService) deletePost(contentPath string) error {
	_, err := s.db.Exec("DELETE FROM posts WHERE content_path = ?", contentPath)
	return err
}
