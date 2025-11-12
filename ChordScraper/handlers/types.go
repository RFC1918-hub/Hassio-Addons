package handlers

// SearchResult represents a single search result from Ultimate Guitar
type SearchResult struct {
	ID     int     `json:"id"`
	Song   string  `json:"song"`
	Artist string  `json:"artist"`
	Type   string  `json:"type"`
	URL    string  `json:"url"`
	Rating float64 `json:"rating"`
}

// OnSongRequest represents the request body for /onsong endpoint
type OnSongRequest struct {
	ID int `json:"id"`
}

// WorshipChordsRequest represents the request body for /worshipchords endpoint
type WorshipChordsRequest struct {
	URL string `json:"url"`
}

// GoogleDriveRequest represents the request body for /send-to-drive endpoint
type GoogleDriveRequest struct {
	Content             string `json:"content"`
	Song                string `json:"song"`
	Artist              string `json:"artist"`
	ID                  string `json:"id"` // Changed to string to support "manual-123" and "wc-456" IDs
	IsManualSubmission  bool   `json:"isManualSubmission"`
	RequiresAutomation  bool   `json:"requiresAutomation"`
}

// FormatManualRequest represents the request body for /format-manual endpoint
type FormatManualRequest struct {
	Song    string `json:"song"`
	Artist  string `json:"artist"`
	Content string `json:"content"`
}
