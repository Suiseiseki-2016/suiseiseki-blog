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
		t.Fatalf("create test db: %v", err)
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
	`, slug, title, "Test summary", "Test category", contentPath)
	if err != nil {
		t.Fatalf("insert test post: %v", err)
	}
}

func TestGetPosts(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	insertTestPost(t, db, "test-post-1", "Test Post 1", "/test/path/test-post-1.md")
	insertTestPost(t, db, "test-post-2", "Test Post 2", "/test/path/test-post-2.md")

	handler := NewPostsHandler(db, "/test/posts")

	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.GET("/api/posts", handler.GetPosts)

	req, _ := http.NewRequest("GET", "/api/posts", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("want status 200, got %d", w.Code)
	}

	var response map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("parse response: %v", err)
	}

	posts, ok := response["posts"].([]interface{})
	if !ok {
		t.Fatal("response missing posts field")
	}

	if len(posts) != 2 {
		t.Fatalf("want 2 posts, got %d", len(posts))
	}
}

func TestGetPost(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	slug := "test-post"

	tmpDir := t.TempDir()
	mdPath := filepath.Join(tmpDir, slug+".md")
	os.WriteFile(mdPath, []byte("# Test content\n\nThis is test content."), 0644)

	insertTestPost(t, db, slug, "Test Post", mdPath)

	handler := NewPostsHandler(db, tmpDir)

	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.GET("/api/posts/:slug", handler.GetPost)

	req, _ := http.NewRequest("GET", "/api/posts/"+slug, nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("want status 200, got %d: %s", w.Code, w.Body.String())
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
		t.Fatalf("want status 404, got %d", w.Code)
	}
}
