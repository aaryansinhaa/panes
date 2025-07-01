package main

import (
	"fmt"

	"github.com/aaryansinhaa/panes/utils/config"
	"github.com/aaryansinhaa/panes/utils/mcp"
	"github.com/aaryansinhaa/panes/utils/server"
)

func main() {
	fmt.Println("Hello, from Panes!")

	// loading config
	cfg := config.MustLoadConfig()

	mcp.Server()
	//start the server
	server.LoadServer(cfg)

}
