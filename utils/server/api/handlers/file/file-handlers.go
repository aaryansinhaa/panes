package file

import (
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func sanitizeFilename(name string) string {
	name = filepath.Base(name)
	name = strings.ReplaceAll(name, " ", "_")
	return name
}

// FileUploadHandler handles file uploads
func FileUploadHandler(w http.ResponseWriter, r *http.Request) {
	slog.Info("uploading File")

	r.ParseMultipartForm(10 << 20) // 10 MB limit

	file, handler, err := r.FormFile("file")
	if err != nil {
		slog.Error("Error retrieving the file from form", "error", err)
		http.Error(w, "Error retrieving the file from form", http.StatusBadRequest)
		return
	}
	defer file.Close()

	if handler.Filename == "" {
		slog.Error("No file was uploaded")
		http.Error(w, "No file was uploaded", http.StatusBadRequest)
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
		return
	}
	defer dst.Close()

	_, err = io.Copy(dst, file)
	if err != nil {
		slog.Error("Failed to save uploaded file", "error", err)
		http.Error(w, "Could not save the file", http.StatusInternalServerError)
		return
	}

	slog.Info("File uploaded successfully", "filename", safeFilename)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "File uploaded and saved successfully", "filename": safeFilename})
}

// ListFilesHandler lists uploaded files
func ListFilesHandler(w http.ResponseWriter, r *http.Request) {
	slog.Info("listing files")

	uploadDir := "./uploads"
	files, err := os.ReadDir(uploadDir)
	if err != nil {
		slog.Error("Failed to read upload directory", "error", err)
		http.Error(w, "Could not read the upload directory", http.StatusInternalServerError)
		return
	}

	var fileList []string
	for _, file := range files {
		if !file.IsDir() {
			fileList = append(fileList, file.Name())
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(map[string][]string{"files": fileList})
	if err != nil {
		slog.Error("Failed to write JSON response", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

func DeleteFileHandler(w http.ResponseWriter, r *http.Request) {
	slog.Info("deleting file")

	fileName := r.PathValue("filename")
	if fileName == "" {
		slog.Error("No filename provided")
		http.Error(w, "No filename provided", http.StatusBadRequest)
		return
	}

	safeFilename := sanitizeFilename(fileName)
	filePath := filepath.Join("./uploads", safeFilename)

	// Check if file exists before deleting
	if info, err := os.Stat(filePath); os.IsNotExist(err) || info.IsDir() {
		slog.Error("File not found or is a directory", "filename", safeFilename)
		http.Error(w, "File not found", http.StatusNotFound)
		return
	}

	// Try to delete the file
	err := os.Remove(filePath)
	if err != nil {
		slog.Error("Failed to delete file", "error", err)
		http.Error(w, "Could not delete the file", http.StatusInternalServerError)
		return
	}

	slog.Info("File deleted successfully", "filename", safeFilename)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message":  "File deleted successfully",
		"filename": safeFilename,
	})
}

// search for a file
func SearchFileHandler(w http.ResponseWriter, r *http.Request) {
	slog.Info("searching for file")

	// Get the path variable
	fileName := r.PathValue("filename")
	if fileName == "" {
		slog.Error("No filename provided")
		http.Error(w, "No filename provided", http.StatusBadRequest)
		return
	}

	// Sanitize the filename to avoid path traversal
	safeFilename := sanitizeFilename(fileName)
	filePath := filepath.Join("./uploads", safeFilename)

	// Check if the file exists
	if info, err := os.Stat(filePath); os.IsNotExist(err) || info.IsDir() {
		slog.Error("File not found", "filename", safeFilename)
		http.Error(w, "File not found", http.StatusNotFound)
		return
	}

	// File exists
	slog.Info("File found", "filename", safeFilename)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message":  "File found",
		"filename": safeFilename,
	})
}
