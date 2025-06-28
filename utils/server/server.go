package server

import (
	"fmt"
	"net/http"

	"github.com/aaryansinhaa/panes/utils/config"
	"github.com/aaryansinhaa/panes/utils/server/api"
)

func LoadServer(cfg *config.Config) {

	fmt.Println("Please enter the port number to run the server on (default is 8080):")
	var port string
	fmt.Scanln(&port)
	if port == "" {
		port = "8080"
	}

	fmt.Printf("Starting server on port %s...\n", port)

	router := api.Router()

	server := http.Server{
		Addr:    fmt.Sprintf("%s:%s", cfg.HTTPServer.Address, port),
		Handler: router,
	}

	fmt.Printf("MCP Server is running at http://%s:%s\n", cfg.HTTPServer.Address, port)

	if err := server.ListenAndServe(); err != nil {
		fmt.Printf("Failed to start server: %v\n", err)
		return
	}

	fmt.Println("Server started successfully!")

}
