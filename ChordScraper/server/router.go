package server

import (
	"github.com/RFC1918-hub/Hassio-Add-ons/chord-scraper/config"
	"github.com/RFC1918-hub/Hassio-Add-ons/chord-scraper/handlers"
	"github.com/gorilla/mux"
)

// SetupRoutes configures all application routes
func SetupRoutes(cfg *config.Config) *mux.Router {
	router := mux.NewRouter()

	// Health check endpoint (no rate limiting)
	router.HandleFunc("/health", handlers.HealthHandler).Methods("GET")

	// API endpoints
	router.HandleFunc("/search", handlers.SearchHandler).Methods("GET")
	router.HandleFunc("/onsong", handlers.OnSongHandler).Methods("POST")
	router.HandleFunc("/worshipchords", handlers.WorshipChordsHandler).Methods("POST")
	router.HandleFunc("/format-manual", handlers.FormatManualHandler).Methods("POST")

	// Google Drive webhook endpoint with stricter rate limiting
	strictLimiter := NewRateLimiter(cfg.RateLimit.MaxUploads, cfg.RateLimit.WindowMinutes)
	router.Handle("/send-to-drive",
		strictLimiter.Middleware(
			handlers.NewGoogleDriveHandler(cfg),
		),
	).Methods("POST")

	// Static file serving (React app)
	// This should be last to act as a catch-all
	router.PathPrefix("/").HandlerFunc(handlers.StaticHandler)

	return router
}
