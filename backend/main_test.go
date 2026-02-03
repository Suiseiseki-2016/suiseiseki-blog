package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	_ "github.com/mattn/go-sqlite3"

	"blog-suiseiseki/database"
)

func TestHealthEndpoint(t *testing.T) {
	os.Setenv("MODE", "dev")
	os.Setenv("PORT", "8080")
	os.Setenv("DB_PATH", filepath.Join(t.TempDir(), "test.db"))
	os.Setenv("POSTS_PATH", t.TempDir())

	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test.db")

	db, err := database.New(dbPath)
	if err != nil {
		t.Fatalf("创建数据库失败: %v", err)
	}
	defer db.Close()

	req, _ := http.NewRequest("GET", "/health", nil)
	w := httptest.NewRecorder()

	handler := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
	}

	handler(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("期望状态码 200，得到 %d", w.Code)
	}

	var response map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("解析响应失败: %v", err)
	}

	if response["status"] != "ok" {
		t.Errorf("期望状态 'ok'，得到 %v", response["status"])
	}
}
