package config

import (
	"encoding/json"
	"log"
	"os"
	"strconv"
	"strings"
)

const (
	// Default configuration values
	defaultPort           = 3000
	defaultWindowMinutes  = 15
	defaultMaxRequests    = 100
	defaultMaxUploads     = 50
	defaultWebhookURL     = "https://n8n-058ea47.peakhq.co.za/webhook/7da5ea58-40bd-4470-8347-e4fc3d7dd1c4"
	defaultAllowedOrigins = "http://localhost:3000,http://127.0.0.1:3000,https://chords.peakhq.co.za"

	// File paths
	homeAssistantOptionsPath = "/data/options.json"
)

// Load reads configuration from multiple sources:
// 1. Home Assistant options.json (if running as HA add-on)
// 2. Environment variables
// 3. Default values
func Load() (*Config, error) {
	config := &Config{
		Port: getEnvInt("PORT", defaultPort),
		RateLimit: RateLimitConfig{
			WindowMinutes: defaultWindowMinutes,
			MaxRequests:   defaultMaxRequests,
			MaxUploads:    defaultMaxUploads,
		},
	}

	// Try to load Home Assistant options
	haOptions, err := loadHomeAssistantOptions()
	if err != nil {
		log.Printf("No Home Assistant options found (this is normal for local dev): %v", err)
	}

	// Load webhook URL (priority: HA options > env > default)
	if haOptions != nil && haOptions.WebhookURL != "" {
		config.WebhookURL = haOptions.WebhookURL
		log.Printf("Loaded webhook URL from Home Assistant options")
	} else if envWebhook := os.Getenv("N8N_WEBHOOK_URL"); envWebhook != "" {
		config.WebhookURL = envWebhook
		log.Printf("Loaded webhook URL from environment variable")
	} else {
		config.WebhookURL = defaultWebhookURL
		log.Printf("Using default webhook URL")
	}

	// Load allowed origins (priority: HA options > env > default)
	var originsStr string
	if haOptions != nil && haOptions.AllowedOrigins != "" {
		originsStr = haOptions.AllowedOrigins
		log.Printf("Loaded allowed origins from Home Assistant options")
	} else if envOrigins := os.Getenv("ALLOWED_ORIGINS"); envOrigins != "" {
		originsStr = envOrigins
		log.Printf("Loaded allowed origins from environment variable")
	} else {
		originsStr = defaultAllowedOrigins
		log.Printf("Using default allowed origins")
	}

	// Parse allowed origins (comma-separated list)
	config.AllowedOrigins = parseAllowedOrigins(originsStr)

	log.Printf("Configuration loaded successfully:")
	log.Printf("  Port: %d", config.Port)
	log.Printf("  Webhook URL: %s", config.WebhookURL)
	log.Printf("  Allowed Origins: %v", config.AllowedOrigins)
	log.Printf("  Rate Limit: %d requests per %d minutes", config.RateLimit.MaxRequests, config.RateLimit.WindowMinutes)

	return config, nil
}

// loadHomeAssistantOptions reads the /data/options.json file
func loadHomeAssistantOptions() (*HomeAssistantOptions, error) {
	// Check if file exists
	if _, err := os.Stat(homeAssistantOptionsPath); os.IsNotExist(err) {
		return nil, err
	}

	// Read file
	data, err := os.ReadFile(homeAssistantOptionsPath)
	if err != nil {
		return nil, err
	}

	// Parse JSON
	var options HomeAssistantOptions
	if err := json.Unmarshal(data, &options); err != nil {
		return nil, err
	}

	return &options, nil
}

// parseAllowedOrigins splits comma-separated origins and trims whitespace
func parseAllowedOrigins(originsStr string) []string {
	parts := strings.Split(originsStr, ",")
	origins := make([]string, 0, len(parts))

	for _, origin := range parts {
		trimmed := strings.TrimSpace(origin)
		if trimmed != "" {
			origins = append(origins, trimmed)
		}
	}

	return origins
}

// getEnvInt gets an integer from environment variable or returns default
func getEnvInt(key string, defaultValue int) int {
	if val := os.Getenv(key); val != "" {
		if intVal, err := strconv.Atoi(val); err == nil {
			return intVal
		}
	}
	return defaultValue
}
