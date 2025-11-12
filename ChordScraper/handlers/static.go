package handlers

import (
	"log"
	"net/http"
	"os"
	"path/filepath"
)

const (
	// Frontend build directory
	frontendDir = "./build/frontend"
)

// StaticHandler serves the React frontend static files
func StaticHandler(w http.ResponseWriter, r *http.Request) {
	// Get the requested path
	path := r.URL.Path

	// Construct full file path
	fullPath := filepath.Join(frontendDir, path)

	// Check if file exists
	info, err := os.Stat(fullPath)
	if err != nil {
		// File doesn't exist or error accessing it
		// Serve index.html for React routing (SPA fallback)
		serveIndexHTML(w, r)
		return
	}

	// If it's a directory, serve index.html
	if info.IsDir() {
		indexPath := filepath.Join(fullPath, "index.html")
		if _, err := os.Stat(indexPath); err == nil {
			http.ServeFile(w, r, indexPath)
			return
		}
		// No index.html in directory, serve root index.html
		serveIndexHTML(w, r)
		return
	}

	// File exists, serve it
	http.ServeFile(w, r, fullPath)
}

// serveIndexHTML serves the root index.html file for React routing
func serveIndexHTML(w http.ResponseWriter, r *http.Request) {
	indexPath := filepath.Join(frontendDir, "index.html")

	// Check if index.html exists
	if _, err := os.Stat(indexPath); err != nil {
		log.Printf("Frontend build not found at %s. Did you run the frontend build?", frontendDir)
		http.Error(w, "Frontend not found. Please build the React application first.", http.StatusNotFound)
		return
	}

	// Serve index.html
	http.ServeFile(w, r, indexPath)
}
