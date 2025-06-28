package api

import (
	"net/http"

	"github.com/aaryansinhaa/panes/utils/server/api/handlers"
	"github.com/aaryansinhaa/panes/utils/server/api/handlers/file"
	"github.com/aaryansinhaa/panes/utils/server/api/handlers/mcp"
	"github.com/aaryansinhaa/panes/utils/services/storage"
)

func Router(s *storage.SQLite) *http.ServeMux {
	router := http.NewServeMux()

	//general routing
	router.HandleFunc("GET /api", handlers.IndexHandler)

	//file based services
	router.HandleFunc("POST /api/files/upload", func(w http.ResponseWriter, r *http.Request) {
		file.FileUploadHandler(w, r, s)
	})
	router.HandleFunc("GET /api/files/list", func(w http.ResponseWriter, r *http.Request) {
		file.ListFilesHandler(w, r, s)
	})
	router.HandleFunc("GET /api/files/list/{filename}", func(w http.ResponseWriter, r *http.Request) {
		file.SearchFileHandler(w, r, s)
	})
	router.HandleFunc("DELETE /api/files/delete/{filename}", func(w http.ResponseWriter, r *http.Request) {
		file.DeleteFileHandler(w, r, s)
	})

	//permission based services

	//mcp based services
	router.HandleFunc("POST /api/mcp/register", mcp.RegisterClientHandler)
	router.HandleFunc("GET /api/mcp/clients", mcp.ListClientsHandler)
	router.HandleFunc("DELETE /api/mcp/clients/delete/{clientId}", mcp.DeleteClientHandler)
	router.HandleFunc("PUT /api/mcp/clients/update/{clientId}", mcp.UpdateClientHandler)

	//activity service

	//llm related services

	return router
}
