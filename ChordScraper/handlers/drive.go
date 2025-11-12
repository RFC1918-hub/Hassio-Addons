package handlers

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/RFC1918-hub/Hassio-Add-ons/chord-scraper/config"
)

// GoogleDriveHandler returns an HTTP handler that proxies requests to the n8n webhook
func NewGoogleDriveHandler(cfg *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("=== Google Drive Upload Request Started ===")
		log.Printf("Remote Address: %s", r.RemoteAddr)
		log.Printf("Content-Type: %s", r.Header.Get("Content-Type"))

		// Read request body
		body, err := io.ReadAll(r.Body)
		if err != nil {
			log.Printf("ERROR: Failed to read request body: %v", err)
			http.Error(w, "Failed to read request body", http.StatusBadRequest)
			return
		}
		defer r.Body.Close()

		log.Printf("Request body size: %d bytes", len(body))

		// Parse request
		var req GoogleDriveRequest
		if err := json.Unmarshal(body, &req); err != nil {
			log.Printf("ERROR: Invalid JSON: %v", err)
			http.Error(w, "Invalid JSON request", http.StatusBadRequest)
			return
		}

		log.Printf("Parsed request - Song: %s, Artist: %s, ID: %s, Manual: %v",
			req.Song, req.Artist, req.ID, req.IsManualSubmission)

		// Validate required fields
		if req.Content == "" || req.Song == "" || req.Artist == "" || req.ID == "" {
			log.Printf("ERROR: Missing required fields - Content: %v, Song: %v, Artist: %v, ID: %v",
				req.Content != "", req.Song != "", req.Artist != "", req.ID != "")
			http.Error(w, "Missing required fields: content, song, artist, id", http.StatusBadRequest)
			return
		}

		// Log submission type
		// Note: We no longer format here since formatting is done in the preview step
		// Content coming from /send-to-drive is already formatted and potentially edited by the user
		if req.IsManualSubmission {
			log.Printf("Manual submission detected (content already formatted from preview)")
		}

		log.Printf("Forwarding request to n8n webhook: %s", cfg.WebhookURL)

		// Forward request to n8n webhook
		if err := forwardToWebhook(cfg.WebhookURL, req, w); err != nil {
			log.Printf("ERROR: Failed to forward to webhook: %v", err)
			// Error response already written by forwardToWebhook
		} else {
			log.Printf("=== Google Drive Upload Request Completed Successfully ===")
		}
	}
}

// forwardToWebhook sends the request to the n8n webhook and forwards the response
func forwardToWebhook(webhookURL string, req GoogleDriveRequest, w http.ResponseWriter) error {
	log.Printf("Marshaling request body for webhook...")
	// Marshal request to JSON
	reqBody, err := json.Marshal(req)
	if err != nil {
		log.Printf("ERROR: Failed to marshal request: %v", err)
		http.Error(w, "Failed to marshal request", http.StatusInternalServerError)
		return err
	}

	log.Printf("Request payload size: %d bytes", len(reqBody))

	// Create HTTP request
	log.Printf("Creating HTTP POST request to: %s", webhookURL)
	httpReq, err := http.NewRequest("POST", webhookURL, bytes.NewReader(reqBody))
	if err != nil {
		log.Printf("ERROR: Failed to create webhook request: %v", err)
		http.Error(w, "Failed to create webhook request", http.StatusInternalServerError)
		return err
	}

	// Set headers
	httpReq.Header.Set("Content-Type", "application/json")

	// Create HTTP client with timeout
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	// Execute request
	log.Printf("Sending request to n8n webhook...")
	resp, err := client.Do(httpReq)
	if err != nil {
		log.Printf("ERROR: Failed to connect to webhook: %v", err)
		log.Printf("ERROR: This could be a network connectivity issue or the webhook URL is incorrect")
		http.Error(w, "Failed to connect to Google Drive service", http.StatusInternalServerError)
		return err
	}
	defer resp.Body.Close()

	// Read response body
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("ERROR: Failed to read webhook response: %v", err)
		http.Error(w, "Failed to read webhook response", http.StatusInternalServerError)
		return err
	}

	log.Printf("n8n webhook response: Status=%d, Body size=%d bytes", resp.StatusCode, len(respBody))

	if len(respBody) > 0 && len(respBody) < 500 {
		log.Printf("n8n webhook response body: %s", string(respBody))
	}

	// Handle error responses
	if resp.StatusCode >= 400 {
		var errorMsg string
		if resp.StatusCode == http.StatusNotFound {
			errorMsg = "n8n webhook not found. Please check if the workflow is active and the webhook URL is correct."
			log.Printf("ERROR: Webhook returned 404 - workflow may not be active")
		} else {
			errorMsg = "Failed to send to Google Drive"
			log.Printf("ERROR: Webhook returned error status: %d", resp.StatusCode)
		}

		// Try to parse response as JSON
		var errorResponse map[string]interface{}
		if json.Unmarshal(respBody, &errorResponse) == nil {
			errorResponse["error"] = errorMsg
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(resp.StatusCode)
			json.NewEncoder(w).Encode(errorResponse)
		} else {
			// Return plain text error
			http.Error(w, errorMsg, resp.StatusCode)
		}
		return nil
	}

	// Forward successful response
	log.Printf("Webhook request successful, forwarding response to client")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(resp.StatusCode)
	w.Write(respBody)

	return nil
}
