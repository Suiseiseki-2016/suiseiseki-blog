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
		t.Fatalf("创建数据库失败: %v", err)
	}
	defer db.Close()

	mdFile := filepath.Join(postsDir, "test-post.md")
	content := `---
title: 测试文章
slug: test-post
summary: 测试摘要
category: 测试
published_at: 2024-01-01
---

# 测试文章

这是测试内容。`

	os.WriteFile(mdFile, []byte(content), 0644)

	syncService := NewSyncService(db.Conn(), postsDir, true, nil, "")

	if err := syncService.Sync(); err != nil {
		t.Fatalf("同步失败: %v", err)
	}

	var count int
	err = db.Conn().QueryRow("SELECT COUNT(*) FROM posts WHERE slug = ?", "test-post").Scan(&count)
	if err != nil {
		t.Fatalf("查询失败: %v", err)
	}

	if count != 1 {
		t.Fatalf("期望1篇文章，得到 %d", count)
	}

	var title, slug string
	err = db.Conn().QueryRow("SELECT title, slug FROM posts WHERE slug = ?", "test-post").Scan(&title, &slug)
	if err != nil {
		t.Fatalf("查询文章失败: %v", err)
	}

	if title != "测试文章" {
		t.Errorf("期望标题 '测试文章'，得到 %q", title)
	}
}

func TestSyncService_DeleteRemovedPosts(t *testing.T) {
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test.db")
	postsDir := filepath.Join(tmpDir, "posts")

	os.MkdirAll(postsDir, 0755)

	db, err := database.New(dbPath)
	if err != nil {
		t.Fatalf("创建数据库失败: %v", err)
	}
	defer db.Close()

	mdFile := filepath.Join(postsDir, "old-post.md")
	os.WriteFile(mdFile, []byte(`---
title: 旧文章
slug: old-post
---

# 旧文章`), 0644)

	syncService := NewSyncService(db.Conn(), postsDir, true, nil, "")
	if err := syncService.Sync(); err != nil {
		t.Fatalf("第一次同步失败: %v", err)
	}

	os.Remove(mdFile)

	if err := syncService.Sync(); err != nil {
		t.Fatalf("第二次同步失败: %v", err)
	}

	var count int
	err = db.Conn().QueryRow("SELECT COUNT(*) FROM posts WHERE slug = ?", "old-post").Scan(&count)
	if err != nil {
		t.Fatalf("查询失败: %v", err)
	}

	if count != 0 {
		t.Fatalf("期望文章被删除，但数据库中仍有 %d 条记录", count)
	}
}
