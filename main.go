package main

import (
	"github.com/korvised/go-ecommerce/config"
	"github.com/korvised/go-ecommerce/modules/servers"
	"github.com/korvised/go-ecommerce/pkg/databases"
	"os"
)

func envPath() string {
	if len(os.Args) == 1 {
		return ".env"
	} else {
		return os.Args[1]
	}
}

func main() {
	// Load environment config
	cfg := config.LoadConfig(envPath())

	// Initial database connection
	db := databases.DbConnect(cfg.Db())
	defer db.Close()

	// Start server
	servers.NewServer(cfg, db).Start()
}
