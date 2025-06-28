package api

import (
	"net/http"

	"github.com/aaryansinhaa/panes/utils/server/api/handlers"
	"github.com/aaryansinhaa/panes/utils/server/api/handlers/file"
	"github.com/aaryansinhaa/panes/utils/server/api/handlers/mcp"
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
	router.HandleFunc("POST /api/mcp/register", mcp.RegisterClientHandler)
	router.HandleFunc("GET /api/mcp/clients", mcp.ListClientsHandler)
	router.HandleFunc("DELETE /api/mcp/clients/delete/{clientId}", mcp.DeleteClientHandler)
	router.HandleFunc("PUT /api/mcp/clients/update/{clientId}", mcp.UpdateClientHandler)

	//activity service

	//llm related services

	return router
}
