package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"backend-sarpras/internal/config"
	"backend-sarpras/internal/db"
)

// Struct model untuk data ruangan
type Ruangan struct {
	ID          int    `json:"id"`
	KodeRuangan string `json:"kode_ruangan"`
	NamaRuangan string `json:"nama_ruangan"`
	Lokasi      string `json:"lokasi"`
	Kapasitas   int    `json:"kapasitas"`
	Deskripsi   string `json:"deskripsi"`
}

func main() {
	// Load konfigurasi dari environment variable
	cfg := config.Load()

	// Buka koneksi database Supabase
	conn := db.Open(cfg.DatabaseURL)
	defer conn.Close()

	// Router HTTP
	mux := http.NewServeMux()

	// Endpoint health check
	mux.HandleFunc("/api/health", func(w http.ResponseWriter, r *http.Request) {
		resp := map[string]string{"status": "ok"}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	})

	// Endpoint untuk mengambil data ruangan
	mux.HandleFunc("/api/ruangan", func(w http.ResponseWriter, r *http.Request) {
		handleListRuangan(w, r, conn)
	})

	// Jalankan server HTTP
	log.Printf("server running on http://localhost:%s", cfg.Port)
	if err := http.ListenAndServe(":"+cfg.Port, mux); err != nil {
		log.Fatal(err)
	}
}

// Handler untuk GET /api/ruangan
func handleListRuangan(w http.ResponseWriter, r *http.Request, dbConn *sql.DB) {
	// Hanya menerima method GET
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	// Query seluruh ruangan
	rows, err := dbConn.Query(`
		SELECT id, kode_ruangan, nama_ruangan, lokasi, kapasitas, deskripsi
		FROM ruangan ORDER BY id
	`)
	if err != nil {
		log.Println("query error:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	// Parsing hasil query ke slice struct
	var result []Ruangan
	for rows.Next() {
		var r Ruangan
		if err := rows.Scan(&r.ID, &r.KodeRuangan, &r.NamaRuangan, &r.Lokasi, &r.Kapasitas, &r.Deskripsi); err != nil {
			log.Println("scan error:", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		result = append(result, r)
	}

	// Kirim respons JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}