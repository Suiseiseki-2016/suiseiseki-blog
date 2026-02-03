package database

import (
	"path/filepath"
	"testing"
)

func TestNew(t *testing.T) {
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test.db")

	db, err := New(dbPath)
	if err != nil {
		t.Fatalf("创建数据库失败: %v", err)
	}
	defer db.Close()

	var count int
	err = db.Conn().QueryRow("SELECT COUNT(*) FROM posts").Scan(&count)
	if err != nil {
		t.Fatalf("查询posts表失败: %v", err)
	}
	if count != 0 {
		t.Fatalf("新数据库posts表应为空，得到 %d", count)
	}
}
