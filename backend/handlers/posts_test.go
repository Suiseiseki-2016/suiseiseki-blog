package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"

	"blog-suiseiseki/database"
)

func setupTestDB(t *testing.T) (*sql.DB, func()) {
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test.db")

	db, err := database.New(dbPath)
	if err != nil {
		t.Fatalf("创建测试数据库失败: %v", err)
	}

	return db.Conn(), func() {
		db.Close()
		os.Remove(dbPath)
	}
}

func insertTestPost(t *testing.T, db *sql.DB, slug, title string, contentPath string) {
	_, err := db.Exec(`
		INSERT INTO posts (slug, title, summary, category, published_at, content_path)
		VALUES (?, ?, ?, ?, datetime('now'), ?)
	`, slug, title, "测试摘要", "测试分类", contentPath)
	if err != nil {
		t.Fatalf("插入测试数据失败: %v", err)
	}
}

func TestGetPosts(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	insertTestPost(t, db, "test-post-1", "测试文章1", "/test/path/test-post-1.md")
	insertTestPost(t, db, "test-post-2", "测试文章2", "/test/path/test-post-2.md")

	handler := NewPostsHandler(db, "/test/posts")

	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.GET("/api/posts", handler.GetPosts)

	req, _ := http.NewRequest("GET", "/api/posts", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("期望状态码 200，得到 %d", w.Code)
	}

	var response map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("解析响应失败: %v", err)
	}

	posts, ok := response["posts"].([]interface{})
	if !ok {
		t.Fatal("响应中没有posts字段")
	}

	if len(posts) != 2 {
		t.Fatalf("期望2篇文章，得到 %d", len(posts))
	}
}

func TestGetPost(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	slug := "test-post"

	tmpDir := t.TempDir()
	mdPath := filepath.Join(tmpDir, slug+".md")
	os.WriteFile(mdPath, []byte("# 测试内容\n\n这是测试内容。"), 0644)

	insertTestPost(t, db, slug, "测试文章", mdPath)

	handler := NewPostsHandler(db, tmpDir)

	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.GET("/api/posts/:slug", handler.GetPost)

	req, _ := http.NewRequest("GET", "/api/posts/"+slug, nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("期望状态码 200，得到 %d: %s", w.Code, w.Body.String())
	}
}

func TestGetPostNotFound(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	handler := NewPostsHandler(db, "/test/posts")

	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.GET("/api/posts/:slug", handler.GetPost)

	req, _ := http.NewRequest("GET", "/api/posts/non-existent", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("期望状态码 404，得到 %d", w.Code)
	}
}
