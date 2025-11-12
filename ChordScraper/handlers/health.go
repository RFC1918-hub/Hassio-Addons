package handlers

import (
	"net/http"
)

// HealthHandler returns a simple OK response for health checks
func HealthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}
