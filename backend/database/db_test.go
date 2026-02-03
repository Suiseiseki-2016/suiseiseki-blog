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
		t.Fatalf("create db: %v", err)
	}
	defer db.Close()

	var count int
	err = db.Conn().QueryRow("SELECT COUNT(*) FROM posts").Scan(&count)
	if err != nil {
		t.Fatalf("query posts: %v", err)
	}
	if count != 0 {
		t.Fatalf("new db posts table should be empty, got %d", count)
	}
}
