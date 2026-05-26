package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/notwinterdust/otp-server/middleware"
	"github.com/notwinterdust/otp-server/models"
	"github.com/notwinterdust/otp-server/storage"
)

type AccountsHandler struct {
	DB *storage.DB
}

func (h *AccountsHandler) Sync(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.UserIDKey).(int64)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	var req models.SyncRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	if err := h.DB.SetAccounts(userID, req.Accounts); err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	accounts, err := h.DB.GetAccounts(userID)
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(models.SyncResponse{Accounts: accounts})
}

func (h *AccountsHandler) Pull(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.UserIDKey).(int64)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	accounts, err := h.DB.GetAccounts(userID)
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(models.SyncResponse{Accounts: accounts})
}
