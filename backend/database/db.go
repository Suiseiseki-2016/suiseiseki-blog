package database

import (
	"database/sql"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
)

type DB struct {
	conn *sql.DB
}

func New(dbPath string) (*DB, error) {
	dir := filepath.Dir(dbPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, err
	}

	conn, err := sql.Open("sqlite3", dbPath+"?_journal_mode=WAL&_synchronous=NORMAL")
	if err != nil {
		return nil, err
	}

	db := &DB{conn: conn}

	if err := db.initSchema(); err != nil {
		return nil, err
	}

	return db, nil
}

func (db *DB) initSchema() error {
	schema := `
	CREATE TABLE IF NOT EXISTS posts (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		slug TEXT UNIQUE,
		title TEXT,
		summary TEXT,
		category TEXT,
		published_at DATETIME,
		content_path TEXT,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);

	CREATE INDEX IF NOT EXISTS idx_posts_published_at ON posts(published_at DESC);
	CREATE INDEX IF NOT EXISTS idx_posts_slug ON posts(slug);
	`

	_, err := db.conn.Exec(schema)
	return err
}

func (db *DB) Close() error {
	return db.conn.Close()
}

func (db *DB) Conn() *sql.DB {
	return db.conn
}
