package main

import (
	"log"

	"github.com/RFC1918-hub/Hassio-Add-ons/chord-scraper/config"
	"github.com/RFC1918-hub/Hassio-Add-ons/chord-scraper/server"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Create and start server
	srv := server.New(cfg)
	if err := srv.Start(); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
