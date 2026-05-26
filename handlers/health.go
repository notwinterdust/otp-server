package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/notwinterdust/otp-server/models"
)

const ServerVersion = "2.0.0"

func Health(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(models.HealthResponse{
		Status:  "ok",
		Version: ServerVersion,
	})
}
