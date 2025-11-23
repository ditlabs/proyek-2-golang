package handlers

import (
	"encoding/json"
	"net/http"
)

type InfoUmum struct {
	NamaInstansi string `json:"nama_instansi"`
	Alamat       string `json:"alamat"`
	Kontak       string `json:"kontak"`
	Deskripsi    string `json:"deskripsi"`
}

func InfoUmumHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	info := InfoUmum{
		NamaInstansi: "Universitas Contoh",
		Alamat:       "Jl. Pendidikan No. 1, Kota Edukasi",
		Kontak:       "(021) 12345678",
		Deskripsi:    "Sistem Informasi Sarana dan Prasarana untuk pengelolaan ruangan dan peminjaman di lingkungan kampus.",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(info)
}
