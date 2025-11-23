package handlers

import (
	"backend-sarpras/middleware"
	"backend-sarpras/models"
	"backend-sarpras/repositories"
	"backend-sarpras/services"
	"encoding/json"
	"net/http"
)

type KehadiranHandler struct {
	KehadiranService *services.KehadiranService
	KehadiranRepo    *repositories.KehadiranRepository
	PeminjamanRepo   *repositories.PeminjamanRepository
	RuanganRepo      *repositories.RuanganRepository
	UserRepo         *repositories.UserRepository
}

func NewKehadiranHandler(
	kehadiranService *services.KehadiranService,
	kehadiranRepo *repositories.KehadiranRepository,
	peminjamanRepo *repositories.PeminjamanRepository,
	ruanganRepo *repositories.RuanganRepository,
	userRepo *repositories.UserRepository,
) *KehadiranHandler {
	return &KehadiranHandler{
		KehadiranService: kehadiranService,
		KehadiranRepo:    kehadiranRepo,
		PeminjamanRepo:   peminjamanRepo,
		RuanganRepo:      ruanganRepo,
		UserRepo:         userRepo,
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

func (h *KehadiranHandler) GetRiwayatBySecurity(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	user := middleware.GetUserFromContext(r)
	if user == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	kehadiranList, err := h.KehadiranRepo.GetBySecurityID(user.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	for i := range kehadiranList {
		if kehadiranList[i].PeminjamanID != 0 && h.PeminjamanRepo != nil {
			peminjaman, _ := h.PeminjamanRepo.GetByID(kehadiranList[i].PeminjamanID)
			if peminjaman != nil {
				// Enrich ruangan
				if peminjaman.RuanganID != nil && h.RuanganRepo != nil {
					ruangan, _ := h.RuanganRepo.GetByID(*peminjaman.RuanganID)
					peminjaman.Ruangan = ruangan
				}
				// Enrich peminjam
				if h.UserRepo != nil {
					peminjam, _ := h.UserRepo.GetByID(peminjaman.PeminjamID)
					if peminjam != nil {
						peminjam.PasswordHash = ""
						peminjaman.Peminjam = peminjam
					}
				}
			}
			kehadiranList[i].Peminjaman = peminjaman
		}
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(kehadiranList)
}
