package types

type LogEntry struct {
	ID         int64
	Message    string `json:"message"`
	Type       string `json:"type"`
	Timestamp  string `json:"timestamp"`
	Action     string `json:"action"`
	ClientName string `json:"client_name"`
}

type Client struct {
	ID            int64
	ClientID      string `json:"client_id"`
	ClientName    string `json:"client_name"`
	ClientAPIHash string `json:"client_api_hash"` // encrypted/hashed API key
	Active        bool   `json:"active"`
	CreatedAt     string `json:"created_at"`
}

type FileMetadata struct {
	ID           int64
	OriginalName string `json:"original_name"`
	Filename     string `json:"filename"`
	FilePath     string `json:"file_path"`
	FileSize     int64  `json:"file_size"`
	MimeType     string `json:"mime_type"`
	UploadedAt   string `json:"uploaded_at"`
	Owner        string `json:"owner"` // default 'system'
}

type Permission struct {
	ID         int64
	ClientID   string `json:"client_id"`
	Resource   string `json:"resource"`    // e.g., "file", "logs"
	Action     string `json:"action"`      // e.g., "read", "write",
	ResourceID string `json:"resource_id"` // e.g., file ID or log ID
	CreatedAt  string `json:"created_at"`
}

