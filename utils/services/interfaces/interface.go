package interfaces

type LogEntry interface {
	CreateLogEntry(message, logType, action, clientName string) error
	GetLogEntries(limit int) ([]LogEntry, error)
	GetLogEntryByID(id int64) (LogEntry, error)
}

type Client interface {
	CreateClient(clientID, clientName, clientAPIHash string) error
	GetClients(limit int) ([]Client, error)
	GetClientByID(clientID string) (Client, error)
	UpdateClient(clientID string, clientName string, clientAPIHash string) error
	DeleteClient(clientID string) error
}

type FileMetadata interface {
	UploadFileMetadata(originalName, filename, filePath, mimeType, owner string, fileSize int64) error
	ListFileMetadata(limit int) ([]FileMetadata, error)
	DeleteFileMetadata(id int64) error
	SearchFilesByName(pattern string, limit int) ([]FileMetadata, error)
}

type Permission interface {
	CreatePermission(clientID, resource, action, resourceID string) error
	GetPermissionsByClientID(clientID string) ([]Permission, error)
	GetPermissionByID(id int64) (Permission, error)
	DeletePermission(id int64) error
	UpdatePermissionByClientID(clientID, resource, action, resourceID string) error
	CheckPermission(clientID, resource, action string) (bool, error)
}
