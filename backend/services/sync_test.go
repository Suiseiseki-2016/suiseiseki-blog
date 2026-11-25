package services

import (
	"os"
	"path/filepath"
	"testing"

	_ "github.com/mattn/go-sqlite3"

	"blog-suiseiseki/database"
)

func TestSyncService_Sync(t *testing.T) {
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test.db")
	postsDir := filepath.Join(tmpDir, "posts")

	os.MkdirAll(postsDir, 0755)

	db, err := database.New(dbPath)
	if err != nil {
		t.Fatalf("create db: %v", err)
	}
	defer db.Close()

	mdFile := filepath.Join(postsDir, "test-post.md")
	content := `---
title: Test Post
slug: test-post
summary: Test summary
category: test
published_at: 2024-01-01
---

# Test Post

Test content.`

	os.WriteFile(mdFile, []byte(content), 0644)

	syncService := NewSyncService(db.Conn(), postsDir, true, nil, "")

	if err := syncService.Sync(); err != nil {
		t.Fatalf("sync: %v", err)
	}

	var count int
	err = db.Conn().QueryRow("SELECT COUNT(*) FROM posts WHERE slug = ?", "test-post").Scan(&count)
	if err != nil {
		t.Fatalf("query: %v", err)
	}

	if count != 1 {
		t.Fatalf("want 1 post, got %d", count)
	}

	var title, slug string
	err = db.Conn().QueryRow("SELECT title, slug FROM posts WHERE slug = ?", "test-post").Scan(&title, &slug)
	if err != nil {
		t.Fatalf("query post: %v", err)
	}

	if title != "Test Post" {
		t.Errorf("want title %q, got %q", "Test Post", title)
	}
}

func TestSyncService_DeleteRemovedPosts(t *testing.T) {
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test.db")
	postsDir := filepath.Join(tmpDir, "posts")

	os.MkdirAll(postsDir, 0755)

	db, err := database.New(dbPath)
	if err != nil {
		t.Fatalf("create db: %v", err)
	}
	defer db.Close()

	mdFile := filepath.Join(postsDir, "old-post.md")
	os.WriteFile(mdFile, []byte(`---
title: Old Post
slug: old-post
---

# Old Post`), 0644)

	syncService := NewSyncService(db.Conn(), postsDir, true, nil, "")
	if err := syncService.Sync(); err != nil {
		t.Fatalf("first sync: %v", err)
	}

	os.Remove(mdFile)

	if err := syncService.Sync(); err != nil {
		t.Fatalf("second sync: %v", err)
	}

	var count int
	err = db.Conn().QueryRow("SELECT COUNT(*) FROM posts WHERE slug = ?", "old-post").Scan(&count)
	if err != nil {
		t.Fatalf("query: %v", err)
	}

	if count != 0 {
		t.Fatalf("want post deleted, still have %d rows", count)
	}
}

func TestSyncService_SlugRename(t *testing.T) {
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test.db")
	postsDir := filepath.Join(tmpDir, "posts")

	os.MkdirAll(postsDir, 0755)

	db, err := database.New(dbPath)
	if err != nil {
		t.Fatalf("create db: %v", err)
	}
	defer db.Close()

	mdFile := filepath.Join(postsDir, "post.md")
	if err := os.WriteFile(mdFile, []byte(`---
title: Rename Post
slug: old-slug
---

# Rename Post`), 0644); err != nil {
		t.Fatalf("write initial markdown: %v", err)
	}

	syncService := NewSyncService(db.Conn(), postsDir, true, nil, "")
	if err := syncService.Sync(); err != nil {
		t.Fatalf("first sync: %v", err)
	}

	if err := os.WriteFile(mdFile, []byte(`---
title: Rename Post
slug: new-slug
---

# Rename Post`), 0644); err != nil {
		t.Fatalf("overwrite markdown: %v", err)
	}

	if err := syncService.Sync(); err != nil {
		t.Fatalf("second sync: %v", err)
	}

	var totalCount, oldCount, newCount int
	if err := db.Conn().QueryRow("SELECT COUNT(*) FROM posts").Scan(&totalCount); err != nil {
		t.Fatalf("query total: %v", err)
	}
	if err := db.Conn().QueryRow("SELECT COUNT(*) FROM posts WHERE slug = ?", "old-slug").Scan(&oldCount); err != nil {
		t.Fatalf("query old slug: %v", err)
	}
	if err := db.Conn().QueryRow("SELECT COUNT(*) FROM posts WHERE slug = ?", "new-slug").Scan(&newCount); err != nil {
		t.Fatalf("query new slug: %v", err)
	}

	if totalCount != 1 || oldCount != 0 || newCount != 1 {
		t.Fatalf("unexpected row counts after slug rename: total=%d, old=%d, new=%d", totalCount, oldCount, newCount)
	}
}
