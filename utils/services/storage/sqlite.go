package storage

import (
	"database/sql"
	"fmt"
	"log/slog"

	"github.com/aaryansinhaa/panes/utils/config"
	"github.com/aaryansinhaa/panes/utils/types"
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

//File related services

// UploadFileMetadata uploads file metadata to the SQLite database
func (s *SQLite) UploadFileMetadata(fileMetaData types.FileMetadata) error {
	result, err := s.DB.Prepare(`INSERT INTO files (filename, original_name, file_path, mime_type, owner, file_size) 
	VALUES (?, ?, ?, ?, ?, ?)`)

	if err != nil {
		slog.Error("Failed to upload file metadata", "error", err)
		return err
	}
	_, err = result.Exec(fileMetaData.Filename, fileMetaData.OriginalName, fileMetaData.FilePath, fileMetaData.MimeType, fileMetaData.Owner, fileMetaData.FileSize)
	if err != nil {
		slog.Error("Failed to execute file metadata upload", "error", err)
		return err
	}
	slog.Info("File metadata uploaded successfully", "filename", fileMetaData.Filename)
	return nil
}

// ListFiles lists all uploaded files from the SQLite database
func (s *SQLite) ListFileMetadata() ([]types.FileMetadata, error) {
	rows, err := s.DB.Query("SELECT id, filename, original_name, file_path, mime_type, file_size, uploaded_at, owner FROM files")
	if err != nil {
		slog.Error("Failed to list files", "error", err)
		return nil, err
	}
	defer rows.Close()

	var files []types.FileMetadata
	for rows.Next() {
		var file types.FileMetadata
		if err := rows.Scan(&file.ID, &file.Filename, &file.OriginalName, &file.FilePath, &file.MimeType, &file.FileSize, &file.UploadedAt, &file.Owner); err != nil {
			slog.Error("Failed to scan file row", "error", err)
			return nil, err
		}
		files = append(files, file)
	}
	return files, nil
}

// SearchFilesByName searches for files by name pattern in the SQLite database
func (s *SQLite) SearchFilesByName(pattern string, limit int) ([]types.FileMetadata, error) {
	query := "SELECT id, filename, original_name, file_path, mime_type, file_size, uploaded_at, owner FROM files WHERE filename LIKE ? LIMIT ?"
	rows, err := s.DB.Query(query, "%"+pattern+"%", limit)
	if err != nil {
		slog.Error("Failed to search files by name", "error", err)
		return nil, err
	}
	defer rows.Close()

	var files []types.FileMetadata
	for rows.Next() {
		var file types.FileMetadata
		if err := rows.Scan(&file.ID, &file.Filename, &file.OriginalName, &file.FilePath, &file.MimeType, &file.FileSize, &file.UploadedAt, &file.Owner); err != nil {
			slog.Error("Failed to scan file row", "error", err)
			return nil, err
		}
		files = append(files, file)
	}
	return files, nil
}

// DeleteFileMetadata deletes file metadata by file name(which are bound to be unique), from the SQLite database
func (s *SQLite) DeleteFileMetadata(filename string) error {
	result, err := s.DB.Prepare("DELETE FROM files WHERE filename = ?")
	if err != nil {
		slog.Error("Failed to prepare delete statement", "error", err)
		return err
	}
	_, err = result.Exec(filename)
	if err != nil {
		slog.Error("Failed to delete file metadata", "error", err)
		return err
	}
	slog.Info("File metadata deleted successfully", "id", filename)
	return nil
}
