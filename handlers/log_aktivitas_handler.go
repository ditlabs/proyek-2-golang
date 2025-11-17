package handlers

import (
	"encoding/json"
	"net/http"
	"backend-sarpras/repositories"
)

type LogAktivitasHandler struct {
	LogRepo *repositories.LogAktivitasRepository
}

func NewLogAktivitasHandler(logRepo *repositories.LogAktivitasRepository) *LogAktivitasHandler {
	return &LogAktivitasHandler{LogRepo: logRepo}
}

func (h *LogAktivitasHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	filter := r.URL.Query().Get("filter")
	logs, err := h.LogRepo.GetAll(filter)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(logs)
}

