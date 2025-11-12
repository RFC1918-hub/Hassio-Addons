package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/RFC1918-hub/Hassio-Add-ons/chord-scraper/config"
)

// Server represents the HTTP server
type Server struct {
	httpServer *http.Server
	config     *config.Config
}

// New creates a new Server instance
func New(cfg *config.Config) *Server {
	return &Server{
		config: cfg,
	}
}

// Start initializes and starts the HTTP server
func (s *Server) Start() error {
	// Set up routes
	router := SetupRoutes(s.config)

	// Apply global middleware
	handler := LoggingMiddleware(router)
	handler = RecoveryMiddleware(handler)

	// Apply rate limiting
	rateLimiter := NewRateLimiter(s.config.RateLimit.MaxRequests, s.config.RateLimit.WindowMinutes)
	handler = rateLimiter.Middleware(handler)

	// Apply CORS
	corsHandler := CORSMiddleware(s.config)
	handler = corsHandler.Handler(handler)

	// Configure HTTP server
	s.httpServer = &http.Server{
		Addr:           fmt.Sprintf("0.0.0.0:%d", s.config.Port),
		Handler:        handler,
		ReadTimeout:    15 * time.Second,
		WriteTimeout:   15 * time.Second,
		IdleTimeout:    60 * time.Second,
		MaxHeaderBytes: 1 << 20, // 1 MB
	}

	// Start server in a goroutine
	go func() {
		log.Printf("Ultimate Guitar Scraper server starting on port %d", s.config.Port)
		if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal for graceful shutdown
	s.waitForShutdown()

	return nil
}

// waitForShutdown waits for interrupt signal and gracefully shuts down the server
func (s *Server) waitForShutdown() {
	// Create channel to listen for interrupt signals
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// Block until signal is received
	sig := <-quit
	log.Printf("Received signal: %v. Shutting down gracefully...", sig)

	// Create context with timeout for shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Attempt graceful shutdown
	if err := s.httpServer.Shutdown(ctx); err != nil {
		log.Printf("Server forced to shutdown: %v", err)
	}

	log.Println("Server stopped")
}
