package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"backend-sarpras/models"
	"backend-sarpras/repositories"
)

type RuanganHandler struct {
	RuanganRepo *repositories.RuanganRepository
}

func NewRuanganHandler(ruanganRepo *repositories.RuanganRepository) *RuanganHandler {
	return &RuanganHandler{RuanganRepo: ruanganRepo}
}

func (h *RuanganHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	ruangans, err := h.RuanganRepo.GetAll()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ruangans)
}

func (h *RuanganHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	id, err := extractIDFromPath(r.URL.Path, "/api/ruangan/")
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	ruangan, err := h.RuanganRepo.GetByID(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if ruangan == nil {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ruangan)
}

func (h *RuanganHandler) Create(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req models.CreateRuanganRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	ruangan := &models.Ruangan{
		KodeRuangan: req.KodeRuangan,
		NamaRuangan: req.NamaRuangan,
		Lokasi:      req.Lokasi,
		Kapasitas:   req.Kapasitas,
		Deskripsi:   req.Deskripsi,
	}

	if err := h.RuanganRepo.Create(ruangan); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(ruangan)
}

func (h *RuanganHandler) Update(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	id, err := extractIDFromPath(r.URL.Path, "/api/ruangan/")
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	var req models.UpdateRuanganRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	ruangan := &models.Ruangan{
		ID:          id,
		NamaRuangan: req.NamaRuangan,
		Lokasi:      req.Lokasi,
		Kapasitas:   req.Kapasitas,
		Deskripsi:   req.Deskripsi,
	}

	if err := h.RuanganRepo.Update(ruangan); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ruangan)
}

func (h *RuanganHandler) Delete(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	id, err := extractIDFromPath(r.URL.Path, "/api/ruangan/")
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	if err := h.RuanganRepo.Delete(id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func extractIDFromPath(path, prefix string) (int, error) {
	if len(path) <= len(prefix) {
		return 0, strconv.ErrSyntax
	}
	raw := strings.Trim(path[len(prefix):], "/")
	return strconv.Atoi(raw)
}

