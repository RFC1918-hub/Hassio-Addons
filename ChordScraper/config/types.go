package config

// Config holds all application configuration
type Config struct {
	Port           int
	WebhookURL     string
	AllowedOrigins []string
	RateLimit      RateLimitConfig
}

// RateLimitConfig holds rate limiting settings
type RateLimitConfig struct {
	WindowMinutes int
	MaxRequests   int
	MaxUploads    int
}

// HomeAssistantOptions represents the structure of /data/options.json
type HomeAssistantOptions struct {
	WebhookURL     string `json:"webhook_url"`
	AllowedOrigins string `json:"allowed_origins"`
}
