package handlers

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
)

// FormatManualHandler formats manual submission content to OnSong format with chord analysis
func FormatManualHandler(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	var req FormatManualRequest
	if err := json.Unmarshal(body, &req); err != nil {
		http.Error(w, "Invalid request format", http.StatusBadRequest)
		return
	}

	// Validate required fields
	if req.Song == "" {
		http.Error(w, "Song title is required", http.StatusBadRequest)
		return
	}

	if req.Content == "" {
		http.Error(w, "Content is required", http.StatusBadRequest)
		return
	}

	// Default artist if not provided
	artist := req.Artist
	if artist == "" {
		artist = "Unknown Artist"
	}

	log.Printf("Formatting manual submission: %s - %s", req.Song, artist)

	// Format using the same function as other sources
	formattedContent := FormatWithNashville(req.Song, artist, "", req.Content)

	// Return formatted content as plain text
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(formattedContent))
}
