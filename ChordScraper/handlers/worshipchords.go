package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/RFC1918-hub/Hassio-Add-ons/chord-scraper/pkg/worshipchords"
)

// WorshipChordsHandler handles requests to fetch worship chords
func WorshipChordsHandler(w http.ResponseWriter, r *http.Request) {
	// Read request body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// Parse request
	var req WorshipChordsRequest
	if err := json.Unmarshal(body, &req); err != nil {
		http.Error(w, "Invalid JSON request", http.StatusBadRequest)
		return
	}

	// Validate URL
	if req.URL == "" {
		http.Error(w, "Missing required parameter: url", http.StatusBadRequest)
		return
	}

	// Validate URL is from worshipchords.com
	if !strings.Contains(req.URL, "worshipchords.com") {
		http.Error(w, "Invalid URL: must be from worshipchords.com", http.StatusBadRequest)
		return
	}

	// Create client and fetch song
	client := worshipchords.New()
	song, err := client.GetSongFromURL(req.URL)
	if err != nil {
		log.Printf("Error fetching worshipchords URL %s: %v", req.URL, err)
		http.Error(w, "Failed to retrieve worshipchords format", http.StatusInternalServerError)
		return
	}

	// Format to OnSong format with chord analysis (same as Ultimate Guitar)
	formatted := FormatWithNashville(song.Title, song.Artist, song.Key, song.Content)

	// Return plain text response
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, formatted)
}
