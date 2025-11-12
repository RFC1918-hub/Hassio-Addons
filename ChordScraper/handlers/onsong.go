package handlers

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/RFC1918-hub/Hassio-Add-ons/chord-scraper/pkg/ultimateguitar"
)

// OnSongHandler handles requests to convert tabs to OnSong format
func OnSongHandler(w http.ResponseWriter, r *http.Request) {
	// Read request body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// Parse request
	var req OnSongRequest
	if err := json.Unmarshal(body, &req); err != nil {
		http.Error(w, "Invalid JSON request", http.StatusBadRequest)
		return
	}

	// Validate tab ID
	if req.ID <= 0 {
		http.Error(w, "Invalid tab ID", http.StatusBadRequest)
		return
	}

	// Create scraper and fetch tab
	scraper := ultimateguitar.New()
	onsongContent, err := scraper.GetTabByIDAsOnSong(int64(req.ID))
	if err != nil {
		log.Printf("Error fetching tab %d: %v", req.ID, err)
		http.Error(w, "Failed to retrieve OnSong format", http.StatusInternalServerError)
		return
	}

	// Return plain text response
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(onsongContent))
}
