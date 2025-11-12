package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"strings"
)

const (
	ultimateGuitarSearchURL = "https://www.ultimate-guitar.com/search.php"
	userAgent               = "UGT_ANDROID/4.11.1 (Pixel; 8.1.0)"
)

// SearchHandler handles search requests to Ultimate Guitar
func SearchHandler(w http.ResponseWriter, r *http.Request) {
	// Get title query parameter
	title := r.URL.Query().Get("title")
	if title == "" {
		http.Error(w, "Missing required parameter: title", http.StatusBadRequest)
		return
	}

	// Fetch search results from Ultimate Guitar
	results, err := searchUltimateGuitar(title)
	if err != nil {
		log.Printf("Error searching Ultimate Guitar: %v", err)
		http.Error(w, "Failed to search Ultimate Guitar", http.StatusInternalServerError)
		return
	}

	// Return JSON response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(results)
}

// searchUltimateGuitar fetches and parses search results from Ultimate Guitar
func searchUltimateGuitar(title string) ([]SearchResult, error) {
	// Build search URL
	searchURL := fmt.Sprintf("%s?search_type=title&value=%s",
		ultimateGuitarSearchURL,
		url.QueryEscape(title))

	// Create HTTP request
	req, err := http.NewRequest("GET", searchURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	req.Header.Set("User-Agent", userAgent)
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")

	// Execute request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch search results: %w", err)
	}
	defer resp.Body.Close()

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// Parse HTML to extract JSON data
	results, err := parseSearchResults(string(body))
	if err != nil {
		return nil, fmt.Errorf("failed to parse search results: %w", err)
	}

	// Filter results to get top-rated Chords per artist
	filtered := filterTopResults(results)

	return filtered, nil
}

// parseSearchResults extracts JSON data from the HTML response
func parseSearchResults(html string) ([]SearchResult, error) {
	// Extract <div class="js-store" data-content="...">
	re := regexp.MustCompile(`<div class="js-store"[^>]*data-content="([^"]+)"`)
	matches := re.FindStringSubmatch(html)

	if len(matches) < 2 {
		return []SearchResult{}, nil // No results found
	}

	// Decode HTML entities
	dataContent := decodeHTMLEntities(matches[1])

	// Parse JSON
	var store struct {
		Store struct {
			Page struct {
				Data struct {
					Results []struct {
						ID         int     `json:"id"`
						SongName   string  `json:"song_name"`
						ArtistName string  `json:"artist_name"`
						Type       string  `json:"type"`
						TabURL     string  `json:"tab_url"`
						Rating     float64 `json:"rating"`
					} `json:"results"`
				} `json:"data"`
			} `json:"page"`
		} `json:"store"`
	}

	if err := json.Unmarshal([]byte(dataContent), &store); err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
	}

	// Convert to SearchResult slice
	results := make([]SearchResult, 0, len(store.Store.Page.Data.Results))
	for _, r := range store.Store.Page.Data.Results {
		results = append(results, SearchResult{
			ID:     r.ID,
			Song:   r.SongName,
			Artist: r.ArtistName,
			Type:   r.Type,
			URL:    r.TabURL,
			Rating: r.Rating,
		})
	}

	return results, nil
}

// decodeHTMLEntities decodes common HTML entities
func decodeHTMLEntities(s string) string {
	replacements := map[string]string{
		"&quot;": "\"",
		"&amp;":  "&",
		"&#39;":  "'",
		"&lt;":   "<",
		"&gt;":   ">",
	}

	result := s
	for entity, replacement := range replacements {
		result = strings.ReplaceAll(result, entity, replacement)
	}

	return result
}

// filterTopResults picks the top-rated Chords version per artist
func filterTopResults(results []SearchResult) []SearchResult {
	// Map to store top result per artist
	topResults := make(map[string]SearchResult)

	for _, r := range results {
		artist := r.Artist
		if artist == "" {
			artist = "Unknown"
		}

		current, exists := topResults[artist]
		isChords := strings.ToLower(r.Type) == "chords"
		currentIsChords := strings.ToLower(current.Type) == "chords"

		if !exists {
			// No result for this artist yet
			topResults[artist] = r
		} else if isChords && !currentIsChords {
			// Replace non-Chords with Chords version
			topResults[artist] = r
		} else if isChords && currentIsChords && r.Rating > current.Rating {
			// Both are Chords, pick higher rated
			topResults[artist] = r
		} else if !isChords && !currentIsChords && r.Rating > current.Rating {
			// Neither are Chords, pick higher rated (fallback)
			topResults[artist] = r
		}
	}

	// Convert map to slice
	filtered := make([]SearchResult, 0, len(topResults))
	for _, result := range topResults {
		filtered = append(filtered, result)
	}

	return filtered
}
