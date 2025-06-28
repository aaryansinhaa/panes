package interfaces

import "github.com/aaryansinhaa/panes/utils/types"

type LogEntry interface {
	CreateLogEntry(message, logType, action, clientName string) error
	GetLogEntries(limit int) ([]LogEntry, error)
	DeleteLogEntry(id int64) error
	DeleteAllLogEntries() error
}

type FileMetadata interface {
	UploadFileMetadata(fileMetaData types.FileMetadata) error
	ListFileMetadata() ([]types.FileMetadata, error)
	DeleteFileMetadata(filename string) error
	SearchFilesByName(pattern string, limit int) ([]types.FileMetadata, error)
}

type Client interface {
	CreateClient(clientID, clientName, clientAPIHash string) error
	GetClients(limit int) ([]Client, error)
	GetClientByID(clientID string) (Client, error)
	UpdateClient(clientID string, clientName string, clientAPIHash string) error
	DeleteClient(clientID string) error
}

type Permission interface {
	CreatePermission(clientID, resource, action, resourceID string) error
	GetPermissionsByClientID(clientID string) ([]Permission, error)
	GetPermissionByID(id int64) (Permission, error)
	DeletePermission(id int64) error
	UpdatePermissionByClientID(clientID, resource, action, resourceID string) error
	CheckPermission(clientID, resource, action string) (bool, error)
}
