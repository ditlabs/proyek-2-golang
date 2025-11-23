package handlers

import (
	"backend-sarpras/middleware"
	"backend-sarpras/models"
	"backend-sarpras/repositories"
	"backend-sarpras/services"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type PeminjamanHandler struct {
	PeminjamanService *services.PeminjamanService
	PeminjamanRepo    *repositories.PeminjamanRepository
	RuanganRepo       *repositories.RuanganRepository
	UserRepo          *repositories.UserRepository
}

func NewPeminjamanHandler(
	peminjamanService *services.PeminjamanService,
	peminjamanRepo *repositories.PeminjamanRepository,
	ruanganRepo *repositories.RuanganRepository,
	userRepo *repositories.UserRepository,
) *PeminjamanHandler {
	return &PeminjamanHandler{
		PeminjamanService: peminjamanService,
		PeminjamanRepo:    peminjamanRepo,
		RuanganRepo:       ruanganRepo,
		UserRepo:          userRepo,
	}
}

func (h *PeminjamanHandler) Create(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	user := middleware.GetUserFromContext(r)
	if user == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var req models.CreatePeminjamanRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	peminjaman, err := h.PeminjamanService.CreatePeminjaman(&req, user.ID)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(peminjaman)
}

func (h *PeminjamanHandler) GetMyPeminjaman(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	user := middleware.GetUserFromContext(r)
	if user == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	peminjaman, err := h.PeminjamanRepo.GetByPeminjamID(user.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Enrich dengan data relasi
	for i := range peminjaman {
		if peminjaman[i].RuanganID != nil {
			ruangan, _ := h.RuanganRepo.GetByID(*peminjaman[i].RuanganID)
			peminjaman[i].Ruangan = ruangan
		}
		user, _ := h.UserRepo.GetByID(peminjaman[i].PeminjamID)
		if user != nil {
			user.PasswordHash = ""
			peminjaman[i].Peminjam = user
		}
		items, _ := h.PeminjamanRepo.GetPeminjamanBarang(peminjaman[i].ID)
		peminjaman[i].Barang = items
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(peminjaman)
}

func (h *PeminjamanHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	id, err := extractPeminjamanID(r.URL.Path)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	peminjaman, err := h.PeminjamanRepo.GetByID(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if peminjaman == nil {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}

	// Enrich dengan data relasi
	if peminjaman.RuanganID != nil {
		ruangan, _ := h.RuanganRepo.GetByID(*peminjaman.RuanganID)
		peminjaman.Ruangan = ruangan
	}
	user, _ := h.UserRepo.GetByID(peminjaman.PeminjamID)
	if user != nil {
		user.PasswordHash = ""
		peminjaman.Peminjam = user
	}
	items, _ := h.PeminjamanRepo.GetPeminjamanBarang(peminjaman.ID)
	peminjaman.Barang = items

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(peminjaman)
}

func (h *PeminjamanHandler) GetPending(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	peminjaman, err := h.PeminjamanRepo.GetPending()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Enrich dengan data relasi
	for i := range peminjaman {
		if peminjaman[i].RuanganID != nil {
			ruangan, _ := h.RuanganRepo.GetByID(*peminjaman[i].RuanganID)
			peminjaman[i].Ruangan = ruangan
		}
		user, _ := h.UserRepo.GetByID(peminjaman[i].PeminjamID)
		if user != nil {
			user.PasswordHash = ""
			peminjaman[i].Peminjam = user
		}
		items, _ := h.PeminjamanRepo.GetPeminjamanBarang(peminjaman[i].ID)
		peminjaman[i].Barang = items
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(peminjaman)
}

func (h *PeminjamanHandler) Verifikasi(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	user := middleware.GetUserFromContext(r)
	if user == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	id, err := extractPeminjamanID(strings.TrimSuffix(r.URL.Path, "/verifikasi"))
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	var req models.VerifikasiPeminjamanRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	err = h.PeminjamanService.VerifikasiPeminjaman(id, user.ID, &req)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Verifikasi berhasil"})
}

func (h *PeminjamanHandler) GetJadwalRuangan(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	startStr := r.URL.Query().Get("start")
	endStr := r.URL.Query().Get("end")

	var start, end time.Time
	var err error

	if startStr == "" {
		start = time.Now()
	} else {
		start, err = time.Parse(time.RFC3339, startStr)
		if err != nil {
			http.Error(w, "Invalid start date", http.StatusBadRequest)
			return
		}
	}

	if endStr == "" {
		end = start.AddDate(0, 1, 0) // default 1 bulan ke depan
	} else {
		end, err = time.Parse(time.RFC3339, endStr)
		if err != nil {
			http.Error(w, "Invalid end date", http.StatusBadRequest)
			return
		}
	}

	jadwal, err := h.PeminjamanRepo.GetJadwalRuangan(start, end)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(jadwal)
}

func (h *PeminjamanHandler) GetJadwalAktif(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	startStr := r.URL.Query().Get("start")
	endStr := r.URL.Query().Get("end")

	var start, end time.Time
	var err error

	if startStr == "" {
		start = time.Now()
	} else {
		start, err = time.Parse(time.RFC3339, startStr)
		if err != nil {
			http.Error(w, "Invalid start date", http.StatusBadRequest)
			return
		}
	}

	if endStr == "" {
		end = start.AddDate(0, 0, 7) // default 7 hari ke depan
	} else {
		end, err = time.Parse(time.RFC3339, endStr)
		if err != nil {
			http.Error(w, "Invalid end date", http.StatusBadRequest)
			return
		}
	}

	peminjaman, err := h.PeminjamanRepo.GetJadwalAktif(start, end)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Enrich dengan data relasi
	for i := range peminjaman {
		if peminjaman[i].RuanganID != nil {
			ruangan, _ := h.RuanganRepo.GetByID(*peminjaman[i].RuanganID)
			peminjaman[i].Ruangan = ruangan
		}
		user, _ := h.UserRepo.GetByID(peminjaman[i].PeminjamID)
		if user != nil {
			user.PasswordHash = ""
			peminjaman[i].Peminjam = user
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(peminjaman)
}

func (h *PeminjamanHandler) GetJadwalAktifBelumVerifikasi(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	startStr := r.URL.Query().Get("start")
	endStr := r.URL.Query().Get("end")

	var start, end time.Time
	var err error

	if startStr == "" {
		start = time.Now()
	} else {
		start, err = time.Parse(time.RFC3339, startStr)
		if err != nil {
			http.Error(w, "Invalid start date", http.StatusBadRequest)
			return
		}
	}

	if endStr == "" {
		end = start.AddDate(0, 0, 7)
	} else {
		end, err = time.Parse(time.RFC3339, endStr)
		if err != nil {
			http.Error(w, "Invalid end date", http.StatusBadRequest)
			return
		}
	}

	peminjaman, err := h.PeminjamanRepo.GetJadwalAktifBelumVerifikasi(start, end)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Enrich dengan data relasi
	for i := range peminjaman {
		if peminjaman[i].RuanganID != nil {
			ruangan, _ := h.RuanganRepo.GetByID(*peminjaman[i].RuanganID)
			peminjaman[i].Ruangan = ruangan
		}
		user, _ := h.UserRepo.GetByID(peminjaman[i].PeminjamID)
		if user != nil {
			user.PasswordHash = ""
			peminjaman[i].Peminjam = user
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(peminjaman)
}

func (h *PeminjamanHandler) GetLaporan(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	startStr := r.URL.Query().Get("start")
	endStr := r.URL.Query().Get("end")
	status := r.URL.Query().Get("status")

	var start, end time.Time
	var err error

	if startStr == "" {
		start = time.Now().AddDate(0, -1, 0) // default 1 bulan lalu
	} else {
		start, err = time.Parse(time.RFC3339, startStr)
		if err != nil {
			http.Error(w, "Invalid start date", http.StatusBadRequest)
			return
		}
	}

	if endStr == "" {
		end = time.Now()
	} else {
		end, err = time.Parse(time.RFC3339, endStr)
		if err != nil {
			http.Error(w, "Invalid end date", http.StatusBadRequest)
			return
		}
	}

	peminjaman, err := h.PeminjamanRepo.GetLaporan(start, end, status)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Enrich dengan data relasi
	for i := range peminjaman {
		if peminjaman[i].RuanganID != nil {
			ruangan, _ := h.RuanganRepo.GetByID(*peminjaman[i].RuanganID)
			peminjaman[i].Ruangan = ruangan
		}
		user, _ := h.UserRepo.GetByID(peminjaman[i].PeminjamID)
		if user != nil {
			user.PasswordHash = ""
			peminjaman[i].Peminjam = user
		}
		items, _ := h.PeminjamanRepo.GetPeminjamanBarang(peminjaman[i].ID)
		peminjaman[i].Barang = items
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(peminjaman)
}

func extractPeminjamanID(path string) (int, error) {
	prefix := "/api/peminjaman/"
	if !strings.HasPrefix(path, prefix) || len(path) <= len(prefix) {
		return 0, strconv.ErrSyntax
	}
	raw := strings.Trim(path[len(prefix):], "/")
	return strconv.Atoi(raw)
}
