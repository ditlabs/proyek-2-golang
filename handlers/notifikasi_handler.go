package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"backend-sarpras/middleware"
	"backend-sarpras/repositories"
)

type NotifikasiHandler struct {
	NotifikasiRepo *repositories.NotifikasiRepository
}

func NewNotifikasiHandler(notifikasiRepo *repositories.NotifikasiRepository) *NotifikasiHandler {
	return &NotifikasiHandler{NotifikasiRepo: notifikasiRepo}
}

func (h *NotifikasiHandler) GetMyNotifikasi(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	user := middleware.GetUserFromContext(r)
	if user == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	notifikasi, err := h.NotifikasiRepo.GetByPenerimaID(user.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(notifikasi)
}

func (h *NotifikasiHandler) GetUnreadCount(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	user := middleware.GetUserFromContext(r)
	if user == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	count, err := h.NotifikasiRepo.GetUnreadCount(user.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]int{"count": count})
}

func (h *NotifikasiHandler) MarkAsRead(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPatch {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	user := middleware.GetUserFromContext(r)
	if user == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Extract ID from path like /api/notifikasi/{id}/dibaca
	path := r.URL.Path
	startIdx := len("/api/notifikasi/")
	endIdx := len(path) - len("/dibaca")
	if endIdx <= startIdx {
		http.Error(w, "Invalid path", http.StatusBadRequest)
		return
	}
	idStr := path[startIdx:endIdx]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	err = h.NotifikasiRepo.MarkAsRead(id, user.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Notifikasi ditandai sebagai dibaca"})
}

