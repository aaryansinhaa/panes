package server

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/aaryansinhaa/panes/utils/config"
	"github.com/aaryansinhaa/panes/utils/server/api"
	"github.com/aaryansinhaa/panes/utils/services/storage"
	"github.com/aaryansinhaa/panes/utils/types"
)

func LoadServer(cfg *config.Config) {
	fmt.Println("Please enter the port number to run the server on (default is 8080):")
	var port string
	fmt.Scanln(&port)
	if port == "" {
		port = "8080"
	}

	fmt.Printf("Starting server on port %s...\n", port)

	// Connect to SQLite database
	store, err := storage.LoadSQLiteStorage(cfg)
	if err != nil {
		fmt.Printf("Failed to connect to storage: %v\n", err)
		return
	}
	defer store.Close()
	router := api.Router(store)

	server := &http.Server{
		Addr:    fmt.Sprintf("%s:%s", cfg.HTTPServer.Address, port),
		Handler: router,
	}
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

	go func() {
		fmt.Printf("MCP Server is running at http://%s:%s\n", cfg.HTTPServer.Address, port)
		slog.Info("MCP Server started")
		store.CreateLogEntry(types.LogEntry{
			Message:    "MCP Server started",
			Type:       "info",
			Action:     "start",
			ClientName: "admin",
		},
		)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("Server error", "error", err)
			store.CreateLogEntry(types.LogEntry{
				Message:    fmt.Sprintf("Server error: %v", err),
				Type:       "error",
				Action:     "start",
				ClientName: "admin",
			})
		}
	}()

	// Block until shutdown signal
	<-done
	slog.Info("Shutdown signal received, shutting down...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		slog.Error("Error shutting down server", "error", err)
	} else {
		slog.Info("HTTP server shut down cleanly")
	}
}
