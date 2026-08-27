package handlers

import (
	"encoding/json"
	"net/http"
)

// HealthHandler reports that the service is up. It requires no
// authentication and is used by uptime checks.
//
//	@Summary		Health check
//	@Description	Returns a fixed status payload confirming the API is running. Public endpoint, no authentication required.
//	@Tags			health
//	@Produce		json
//	@Success		200	{object}	map[string]string	"status ok payload"
//	@Router			/health [get]
func HealthHandler(w http.ResponseWriter, r *http.Request) {
	response := map[string]string{
		"status": "ok",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
