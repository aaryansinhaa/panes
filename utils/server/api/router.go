package api

import (
	"net/http"

	"github.com/aaryansinhaa/panes/utils/server/api/handlers"
	"github.com/aaryansinhaa/panes/utils/server/api/handlers/file"
)

func Router() *http.ServeMux {
	router := http.NewServeMux()

	//general routing
	router.HandleFunc("GET /api", handlers.IndexHandler)

	//file based services
	router.HandleFunc("POST /api/files/upload", file.FileUploadHandler)
	router.HandleFunc("GET /api/files/list", file.ListFilesHandler)
	router.HandleFunc("GET /api/files/list/{filename}", file.SearchFileHandler)
	router.HandleFunc("DELETE /api/files/delete/{filename}", file.DeleteFileHandler)

	//permission based services

	//mcp based services

	//activity service

	//llm related services

	return router
}
