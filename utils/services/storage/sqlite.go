package storage

import (
	"database/sql"
	"fmt"
	"log/slog"

	"github.com/aaryansinhaa/panes/utils/config"
	_ "github.com/mattn/go-sqlite3"
)

type SQLite struct {
	DB *sql.DB
}

// Loads the essential SQLite connection and makes essential table for the application: Clients, Files, Permissions, and Logs
func LoadSQLiteStorage(cfg *config.Config) (*SQLite, error) {
	storage, err := sql.Open("sqlite3", cfg.StoragePath)
	if err != nil {
		return nil, err
	}
	fmt.Printf("Connecting to SQLite database at %s\n", cfg.StoragePath)
	_, err = storage.Exec(`CREATE TABLE IF NOT EXISTS logs (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		message TEXT NOT NULL,
		type TEXT NOT NULL,
		timestamp DATETIME DEFAULT CURRENT_TIMESTAMP,
		action TEXT NOT NULL,
		client_name TEXT NOT NULL
	)`)
	if err != nil {
		return nil, err
	}
	_, err = storage.Exec(`CREATE TABLE IF NOT EXISTS clients (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    client_id TEXT NOT NULL UNIQUE,
    client_name TEXT NOT NULL,
    client_api_hash TEXT NOT NULL,  -- encrypted/hashed API key
    active BOOLEAN DEFAULT TRUE,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	)`)
	if err != nil {
		return nil, err
	}

	_, err = storage.Exec(`CREATE TABLE IF NOT EXISTS files (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    filename TEXT NOT NULL UNIQUE,
    original_name TEXT NOT NULL,
    file_path TEXT NOT NULL,
    file_size INTEGER NOT NULL,
    mime_type TEXT,
    uploaded_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    owner TEXT DEFAULT 'system'
	)`)
	if err != nil {
		return nil, err
	}

	_, err = storage.Exec(`CREATE TABLE IF NOT EXISTS permissions (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    client_id TEXT NOT NULL,
    file_id INTEGER,  -- NULL means global permission
    permission_type TEXT NOT NULL, -- READ, WRITE, DELETE
    allowed BOOLEAN DEFAULT FALSE,
    granted_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (client_id) REFERENCES clients(client_id),
    FOREIGN KEY (file_id) REFERENCES files(id)
	)`)
	if err != nil {
		return nil, err
	}

	return &SQLite{DB: storage}, nil
}

// close the sqlite connection
func (s *SQLite) Close() error {
	if s.DB != nil {
		slog.Info("Closing SQLite database connection")
		return s.DB.Close()
	}
	return nil
}
