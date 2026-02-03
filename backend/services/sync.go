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

// SyncEventNotifier notifies the frontend (e.g. via SSE) when sync completes so the list can refresh.
type SyncEventNotifier interface {
	Broadcast()
}

type SyncService struct {
	db          *sql.DB
	postsPath   string
	remoteURL   string
	isDev       bool
	notifier    SyncEventNotifier
}

func NewSyncService(db *sql.DB, postsPath string, isDev bool, notifier SyncEventNotifier, remoteURL string) *SyncService {
	return &SyncService{
		db:        db,
		postsPath: postsPath,
		remoteURL: remoteURL,
		isDev:     isDev,
		notifier:  notifier,
	}
}

// ensurePostsFromRemote clones the remote repo when posts dir is empty or missing.
func (s *SyncService) ensurePostsFromRemote() error {
	if s.remoteURL == "" {
		return nil
	}

	// Dir missing: clone directly to postsPath
	if _, err := os.Stat(s.postsPath); os.IsNotExist(err) {
		parent := filepath.Dir(s.postsPath)
		if err := os.MkdirAll(parent, 0755); err != nil {
			return fmt.Errorf("failed to create dir: %w", err)
		}
		log.Printf("posts dir missing, cloning from remote: %s", s.remoteURL)
		if err := s.gitClone(s.remoteURL, s.postsPath); err != nil {
			return err
		}
		return nil
	}

	// Dir exists: check for any .md files
	files, err := s.scanMarkdownFiles()
	if err != nil {
		return nil // e.g. permission issues; don't block startup
	}
	if len(files) > 0 {
		return nil // already has posts, no clone needed
	}

	// Empty but already a git repo: git pull will run later
	if _, err := os.Stat(filepath.Join(s.postsPath, ".git")); err == nil {
		return nil
	}

	// posts empty and not a git repo: clone to temp dir then replace
	log.Printf("posts empty and not a git repo, cloning from remote: %s", s.remoteURL)
	tmpDir, err := os.MkdirTemp(filepath.Dir(s.postsPath), "posts-clone-*")
	if err != nil {
		return fmt.Errorf("failed to create temp dir: %w", err)
	}
	if err := s.gitClone(s.remoteURL, tmpDir); err != nil {
		os.RemoveAll(tmpDir)
		return err
	}
	if err := os.RemoveAll(s.postsPath); err != nil {
		os.RemoveAll(tmpDir)
		return fmt.Errorf("failed to clean posts dir: %w", err)
	}
	if err := os.Rename(tmpDir, s.postsPath); err != nil {
		os.RemoveAll(tmpDir)
		return fmt.Errorf("failed to replace posts dir: %w", err)
	}
	return nil
}

func (s *SyncService) gitClone(url, dest string) error {
	cmd := exec.Command("git", "clone", "--depth", "1", url, dest)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("git clone failed: %v, output: %s", err, string(output))
	}
	log.Printf("git clone ok: %s", string(output))
	return nil
}

// Sync: ensure posts from remote if needed, scan .md, update DB; prod runs git pull first
func (s *SyncService) Sync() error {
	log.Println("sync: starting...")

	if err := s.ensurePostsFromRemote(); err != nil {
		log.Printf("ensurePostsFromRemote failed: %v", err)
	}

	if !s.isDev {
		if err := s.gitPull(); err != nil {
			log.Printf("git pull failed: %v", err)
		}
	}

	files, err := s.scanMarkdownFiles()
	if err != nil {
		return fmt.Errorf("scan files failed: %w", err)
	}

	log.Printf("sync: found %d markdown file(s)", len(files))

	existingPaths, err := s.getExistingPaths()
	if err != nil {
		return fmt.Errorf("get existing paths failed: %w", err)
	}

	processedPaths := make(map[string]bool)
	for _, filePath := range files {
		processedPaths[filePath] = true
		if err := s.processFile(filePath); err != nil {
			log.Printf("process file %s failed: %v", filePath, err)
		}
	}

	for path := range existingPaths {
		if !processedPaths[path] {
			if err := s.deletePost(path); err != nil {
				log.Printf("delete post %s failed: %v", path, err)
			}
		}
	}

	log.Println("sync: done")
	if s.notifier != nil {
		s.notifier.Broadcast()
	}
	return nil
}

func (s *SyncService) gitPull() error {
	if _, err := os.Stat(filepath.Join(s.postsPath, ".git")); os.IsNotExist(err) {
		return fmt.Errorf("not a git repo: %s", s.postsPath)
	}

	cmd := exec.Command("git", "pull")
	cmd.Dir = s.postsPath
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("git pull failed: %v, output: %s", err, string(output))
	}

	log.Printf("git pull ok: %s", string(output))
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
		return fmt.Errorf("parse markdown failed: %w", err)
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
		return fmt.Errorf("db exec failed: %w", err)
	}

	log.Printf("sync post: %s (%s)", fm.Title, slug)
	return nil
}

func (s *SyncService) deletePost(contentPath string) error {
	_, err := s.db.Exec("DELETE FROM posts WHERE content_path = ?", contentPath)
	return err
}
