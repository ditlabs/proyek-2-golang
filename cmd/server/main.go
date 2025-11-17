package main

import (
	"log"
	"net/http"
	"backend-sarpras/internal/config"
	"backend-sarpras/internal/db"
	"backend-sarpras/internal/router"
)

func main() {
	// Load konfigurasi dari environment variable
	cfg := config.Load()

	// Buka koneksi database Supabase
	conn := db.Open(cfg.DatabaseURL)
	defer conn.Close()

	// Setup router
	handler := router.New(conn)

	// Jalankan server HTTP
	log.Printf("Server running on http://localhost:%s", cfg.Port)
	if err := http.ListenAndServe(":"+cfg.Port, handler); err != nil {
		log.Fatal(err)
	}
}
