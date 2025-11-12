package server

import (
	"log"
	"net/http"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/RFC1918-hub/Hassio-Add-ons/chord-scraper/config"
	"github.com/rs/cors"
	"golang.org/x/time/rate"
)

// LoggingMiddleware logs HTTP requests
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Call the next handler
		next.ServeHTTP(w, r)

		// Log the request
		log.Printf(
			"%s %s %s %s",
			r.Method,
			r.RequestURI,
			r.RemoteAddr,
			time.Since(start),
		)
	})
}

// RecoveryMiddleware recovers from panics and returns 500 error
func RecoveryMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("Panic recovered: %v", err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(w, r)
	})
}

// CORSMiddleware creates a CORS handler with allowed origins
func CORSMiddleware(cfg *config.Config) *cors.Cors {
	return cors.New(cors.Options{
		AllowedOrigins: cfg.AllowedOrigins,
		AllowedMethods: []string{
			http.MethodGet,
			http.MethodPost,
			http.MethodPut,
			http.MethodDelete,
			http.MethodOptions,
		},
		AllowedHeaders: []string{
			"Accept",
			"Authorization",
			"Content-Type",
			"X-CSRF-Token",
		},
		ExposedHeaders:     []string{"Link"},
		AllowCredentials:   true,
		MaxAge:             300, // 5 minutes
		OptionsPassthrough: false,
		Debug:              false,
		// Custom origin validator to support wildcard patterns
		AllowOriginFunc: func(origin string) bool {
			// Allow requests with no origin (same-origin, mobile apps, curl)
			if origin == "" {
				return true
			}

			// Check against allowed origins
			for _, allowed := range cfg.AllowedOrigins {
				// Support wildcard subdomain matching
				if strings.Contains(allowed, "*") {
					pattern := strings.ReplaceAll(allowed, ".", "\\.")
					pattern = strings.ReplaceAll(pattern, "*", ".*")
					pattern = "^" + pattern + "$"

					matched, err := regexp.MatchString(pattern, origin)
					if err == nil && matched {
						log.Printf("CORS allowed origin (wildcard): %s", origin)
						return true
					}
				} else if allowed == origin {
					log.Printf("CORS allowed origin: %s", origin)
					return true
				}
			}

			log.Printf("CORS blocked origin: %s", origin)
			return false
		},
	})
}

// RateLimiter implements token bucket rate limiting per IP
type RateLimiter struct {
	visitors map[string]*visitor
	mu       sync.RWMutex
	rate     rate.Limit
	burst    int
}

type visitor struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

// NewRateLimiter creates a new rate limiter
func NewRateLimiter(requestsPerWindow, windowMinutes int) *RateLimiter {
	// Convert requests per window to requests per second
	rps := float64(requestsPerWindow) / (float64(windowMinutes) * 60.0)

	rl := &RateLimiter{
		visitors: make(map[string]*visitor),
		rate:     rate.Limit(rps),
		burst:    requestsPerWindow,
	}

	// Start cleanup goroutine to remove old visitors
	go rl.cleanupVisitors()

	return rl
}

// Middleware returns a middleware function for rate limiting
func (rl *RateLimiter) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get client IP
		ip := getClientIP(r)

		// Get or create limiter for this IP
		limiter := rl.getVisitor(ip)

		// Check if request is allowed
		if !limiter.Allow() {
			http.Error(w, "Too many requests from this IP, please try again later.", http.StatusTooManyRequests)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// getVisitor retrieves or creates a visitor for the given IP
func (rl *RateLimiter) getVisitor(ip string) *rate.Limiter {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	v, exists := rl.visitors[ip]
	if !exists {
		limiter := rate.NewLimiter(rl.rate, rl.burst)
		rl.visitors[ip] = &visitor{
			limiter:  limiter,
			lastSeen: time.Now(),
		}
		return limiter
	}

	// Update last seen time
	v.lastSeen = time.Now()
	return v.limiter
}

// cleanupVisitors removes old visitors every 5 minutes
func (rl *RateLimiter) cleanupVisitors() {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		rl.mu.Lock()
		for ip, v := range rl.visitors {
			// Remove visitors not seen in the last 30 minutes
			if time.Since(v.lastSeen) > 30*time.Minute {
				delete(rl.visitors, ip)
			}
		}
		rl.mu.Unlock()
	}
}

// getClientIP extracts the client IP from the request
func getClientIP(r *http.Request) string {
	// Try X-Forwarded-For header first (for proxies)
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		ips := strings.Split(xff, ",")
		if len(ips) > 0 {
			return strings.TrimSpace(ips[0])
		}
	}

	// Try X-Real-IP header
	if xri := r.Header.Get("X-Real-IP"); xri != "" {
		return xri
	}

	// Fall back to RemoteAddr
	ip := r.RemoteAddr
	// Remove port if present
	if colon := strings.LastIndex(ip, ":"); colon != -1 {
		ip = ip[:colon]
	}

	return ip
}
