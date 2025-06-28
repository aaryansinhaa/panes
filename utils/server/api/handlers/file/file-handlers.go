package file

import (
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/aaryansinhaa/panes/utils/services/storage"
	"github.com/aaryansinhaa/panes/utils/types"
)

func sanitizeFilename(name string) string {
	name = filepath.Base(name)
	name = strings.ReplaceAll(name, " ", "_")
	return name
}

// FileUploadHandler handles file uploads
func FileUploadHandler(w http.ResponseWriter, r *http.Request, store *storage.SQLite) {
	slog.Info("uploading File")
	var log types.LogEntry
	r.ParseMultipartForm(10 << 20) // 10 MB limit

	file, handler, err := r.FormFile("file")
	if err != nil {
		slog.Error("Error retrieving the file from form", "error", err)
		http.Error(w, "Error retrieving the file from form", http.StatusBadRequest)
		log.Message = "Error retrieving the file from form" + err.Error()
		log.Type = "error"
		log.Action = "upload"
		log.ClientName = "admin"
		err := store.CreateLogEntry(log)
		if err != nil {
			slog.Error("Failed to log error", "error", err)
		}
		return
	}
	defer file.Close()

	if handler.Filename == "" {
		slog.Error("No file was uploaded")
		http.Error(w, "No file was uploaded", http.StatusBadRequest)
		log.Message = "No file was uploaded"
		log.Type = "error"
		log.Action = "upload"
		log.ClientName = "admin"
		err := store.CreateLogEntry(log)
		if err != nil {
			slog.Error("Failed to log error", "error", err)
		}
		return
	}
	//err = interfaces.FileMetadata.UploadFileMetadata()
	var fileMetaData types.FileMetadata
	fileMetaData.Filename = handler.Filename
	fileMetaData.OriginalName = handler.Filename
	fileMetaData.FilePath = filepath.Join("./uploads", sanitizeFilename(handler.Filename))
	fileMetaData.MimeType = handler.Header.Get("Content-Type")
	fileMetaData.FileSize = handler.Size
	fileMetaData.Owner = "system" // default owner, can be changed later
	// Save file metadata to the database
	err = store.UploadFileMetadata(fileMetaData)

	if err != nil {
		slog.Error("Failed to upload file metadata", "error", err)
		http.Error(w, "Could not save file metadata", http.StatusInternalServerError)
		log.Message = "Failed to upload file metadata: " + err.Error()
		log.Type = "error"
		log.Action = "upload"
		log.ClientName = "admin"
		err := store.CreateLogEntry(log)
		if err != nil {
			slog.Error("Failed to log error", "error", err)
		}
		return
	}
	safeFilename := sanitizeFilename(handler.Filename)
	slog.Info("uploaded file", "filename", safeFilename)
	slog.Info("file size", "size", handler.Size)
	slog.Info("file type", "type", handler.Header.Get("Content-Type"))
	slog.Info("Mime Header", "header", handler.Header)

	uploadDir := "./uploads"
	os.MkdirAll(uploadDir, os.ModePerm)

	dstPath := filepath.Join(uploadDir, safeFilename)
	dst, err := os.Create(dstPath)
	if err != nil {
		slog.Error("Failed to create destination file", "error", err)
		http.Error(w, "Could not save the file", http.StatusInternalServerError)
		log.Message = "Failed to create destination file: " + err.Error()
		log.Type = "error"
		log.Action = "upload"
		log.ClientName = "admin"
		err := store.CreateLogEntry(log)
		if err != nil {
			slog.Error("Failed to log error", "error", err)
		}
		return
	}
	defer dst.Close()

	_, err = io.Copy(dst, file)
	if err != nil {
		slog.Error("Failed to save uploaded file", "error", err)
		http.Error(w, "Could not save the file", http.StatusInternalServerError)
		log.Message = "Failed to save uploaded file: " + err.Error()
		log.Type = "error"
		log.Action = "upload"
		log.ClientName = "admin"
		err := store.CreateLogEntry(log)
		if err != nil {
			slog.Error("Failed to log error", "error", err)
		}
		return
	}

	slog.Info("File uploaded successfully", "filename", safeFilename)
	log.Message = "File uploaded successfully: " + safeFilename
	log.Type = "success"
	log.Action = "upload"
	log.ClientName = "admin"
	err = store.CreateLogEntry(log)
	if err != nil {
		slog.Error("Failed to log success", "error", err)
		http.Error(w, "Could not log success", http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "File uploaded and saved successfully", "filename": safeFilename})
}

// ListFilesHandler lists uploaded files
func ListFilesHandler(w http.ResponseWriter, r *http.Request, s *storage.SQLite) {
	slog.Info("listing files")

	var log types.LogEntry

	var fileList []types.FileMetadata
	fileList, err := s.ListFileMetadata()
	if err != nil {
		slog.Error("Failed to list file metadata", "error", err)
		http.Error(w, "Could not retrieve file metadata", http.StatusInternalServerError)
		log.Message = "Failed to list file metadata: " + err.Error()
		log.Type = "error"
		log.Action = "list"
		log.ClientName = "admin"
		err := s.CreateLogEntry(log)
		if err != nil {
			slog.Error("Failed to log error", "error", err)
			http.Error(w, "Could not log error", http.StatusInternalServerError)
		}
		return
	}
	var fileNames []string
	for _, file := range fileList {
		fileNames = append(fileNames, file.Filename)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(map[string][]string{"files": fileNames})
	if err != nil {
		slog.Error("Failed to write JSON response", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
	log.Message = "Listed files successfully"
	log.Type = "success"
	log.Action = "list"
	log.ClientName = "admin"
	err = s.CreateLogEntry(log)
	if err != nil {
		slog.Error("Failed to log success", "error", err)
		http.Error(w, "Could not log success", http.StatusInternalServerError)
		return
	}
}

func DeleteFileHandler(w http.ResponseWriter, r *http.Request, s *storage.SQLite) {
	slog.Info("deleting file")

	var log types.LogEntry

	fileName := r.PathValue("filename")
	if fileName == "" {
		slog.Error("No filename provided")
		http.Error(w, "No filename provided", http.StatusBadRequest)
		log.Message = "No filename provided"
		log.Type = "error"
		log.Action = "delete"
		log.ClientName = "admin"
		err := s.CreateLogEntry(log)
		if err != nil {
			slog.Error("Failed to log error", "error", err)
			http.Error(w, "Could not log error", http.StatusInternalServerError)
		}
		return
	}
	safeFilename := sanitizeFilename(fileName)
	filePath := filepath.Join("./uploads", safeFilename)

	// Check if file exists before deleting
	if info, err := os.Stat(filePath); os.IsNotExist(err) || info.IsDir() {
		slog.Error("File not found or is a directory", "filename", safeFilename)
		http.Error(w, "File not found", http.StatusNotFound)
		log.Message = "File not found: " + safeFilename
		log.Type = "error"
		log.Action = "delete"
		log.ClientName = "admin"
		err = s.CreateLogEntry(log)
		if err != nil {
			slog.Error("Failed to log error", "error", err)
			http.Error(w, "Could not log error", http.StatusInternalServerError)
		}
		return
	}

	// Try to delete the file
	err := os.Remove(filePath)
	if err != nil {
		slog.Error("Failed to delete file", "error", err)
		http.Error(w, "Could not delete the file", http.StatusInternalServerError)
		log.Message = "Failed to delete file: " + safeFilename + " - " + err.Error()
		log.Type = "error"
		log.Action = "delete"
		log.ClientName = "admin"
		err = s.CreateLogEntry(log)
		if err != nil {
			slog.Error("Failed to log error", "error", err)
			http.Error(w, "Could not log error", http.StatusInternalServerError)
		}
		return
	}

	// Remove file metadata from the database
	err = s.DeleteFileMetadata(safeFilename)
	if err != nil {
		slog.Error("Failed to delete file metadata", "error", err)
		http.Error(w, "Could not delete file metadata", http.StatusInternalServerError)
		log.Message = "Failed to delete file metadata: " + safeFilename + " - " + err.Error()
		log.Type = "error"
		log.Action = "delete"
		log.ClientName = "admin"
		err = s.CreateLogEntry(log)
		if err != nil {
			slog.Error("Failed to log error", "error", err)
			http.Error(w, "Could not log error", http.StatusInternalServerError)
		}
		return
	}

	slog.Info("File deleted successfully", "filename", safeFilename)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message":  "File deleted successfully",
		"filename": safeFilename,
	})
	log.Message = "File deleted successfully: " + safeFilename
	log.Type = "success"
	log.Action = "delete"
	log.ClientName = "admin"
	err = s.CreateLogEntry(log)
	if err != nil {
		slog.Error("Failed to log success", "error", err)
		http.Error(w, "Could not log success", http.StatusInternalServerError)
		return
	}
}

// search for FileMetadata by filename
func SearchFileHandler(w http.ResponseWriter, r *http.Request, s *storage.SQLite) {
	fileName := r.PathValue("filename")
	var log types.LogEntry
	if fileName == "" {
		slog.Error("No filename provided")
		http.Error(w, "No filename provided", http.StatusBadRequest)
		log.Message = "No filename provided"
		log.Type = "error"
		log.Action = "search"
		log.ClientName = "admin"
		err := s.CreateLogEntry(log)
		if err != nil {
			slog.Error("Failed to log error", "error", err)
			http.Error(w, "Could not log error", http.StatusInternalServerError)
		}
		return
	}
	safeFilename := sanitizeFilename(fileName)
	fileMetadata, err := s.SearchFilesByName(safeFilename, 1)
	if err != nil {
		slog.Error("Failed to search file metadata", "error", err)
		http.Error(w, "Could not retrieve file metadata", http.StatusInternalServerError)
		log.Message = "Failed to search file metadata: " + err.Error()
		log.Type = "error"
		log.Action = "search"
		log.ClientName = "admin"
		err = s.CreateLogEntry(log)
		if err != nil {
			slog.Error("Failed to log error", "error", err)
			http.Error(w, "Could not log error", http.StatusInternalServerError)
		}
		return
	}
	if len(fileMetadata) == 0 {
		slog.Error("File not found", "filename", safeFilename)
		http.Error(w, "File not found", http.StatusNotFound)
		log.Message = "File not found: " + safeFilename
		log.Type = "error"
		log.Action = "search"
		log.ClientName = "admin"
		err = s.CreateLogEntry(log)
		if err != nil {
			slog.Error("Failed to log error", "error", err)
			http.Error(w, "Could not log error", http.StatusInternalServerError)
		}
		return
	}
	slog.Info("File found", "filename", safeFilename)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(fileMetadata[0])
	log.Message = "File found: " + safeFilename
	log.Type = "success"
	log.Action = "search"
	log.ClientName = "admin"
	err = s.CreateLogEntry(log)
	if err != nil {
		slog.Error("Failed to log success", "error", err)
		http.Error(w, "Could not log success", http.StatusInternalServerError)
		return
	}

}
