package handlers

import (
	"encoding/json"
	"net/http"
	"backend-sarpras/middleware"
	"backend-sarpras/models"
	"backend-sarpras/repositories"
	"backend-sarpras/services"
)

type KehadiranHandler struct {
	KehadiranService *services.KehadiranService
	KehadiranRepo    *repositories.KehadiranRepository
	PeminjamanRepo   *repositories.PeminjamanRepository
}

func NewKehadiranHandler(
	kehadiranService *services.KehadiranService,
	kehadiranRepo *repositories.KehadiranRepository,
	peminjamanRepo *repositories.PeminjamanRepository,
) *KehadiranHandler {
	return &KehadiranHandler{
		KehadiranService: kehadiranService,
		KehadiranRepo:    kehadiranRepo,
		PeminjamanRepo:   peminjamanRepo,
	}
}

func (h *KehadiranHandler) Create(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	user := middleware.GetUserFromContext(r)
	if user == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var req models.CreateKehadiranRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	err := h.KehadiranService.CreateKehadiran(&req, user.ID)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "Kehadiran berhasil dicatat"})
}

func (h *KehadiranHandler) GetByPeminjamanID(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	peminjamanIDStr := r.URL.Query().Get("peminjaman_id")
	if peminjamanIDStr == "" {
		http.Error(w, "peminjaman_id required", http.StatusBadRequest)
		return
	}

	// TODO: parse peminjamanID dari query string
	// Untuk sekarang, return empty atau implementasi sesuai kebutuhan
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode([]models.KehadiranPeminjam{})
}

