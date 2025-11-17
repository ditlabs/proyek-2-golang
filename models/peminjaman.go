package models

import "time"

type Peminjaman struct {
	ID                int       `json:"id"`
	PeminjamID        int       `json:"peminjam_id"`
	Peminjam          *User     `json:"peminjam,omitempty"`
	RuanganID         *int      `json:"ruangan_id"`
	Ruangan           *Ruangan  `json:"ruangan,omitempty"`
	TanggalMulai      time.Time `json:"tanggal_mulai"`
	TanggalSelesai    time.Time `json:"tanggal_selesai"`
	Keperluan         string    `json:"keperluan"`
	Status            string    `json:"status"`
	SuratDigitalURL   string    `json:"surat_digital_url"`
	VerifiedBy        *int      `json:"verified_by"`
	Verifier          *User     `json:"verifier,omitempty"`
	VerifiedAt        *time.Time `json:"verified_at"`
	CatatanVerifikasi string    `json:"catatan_verifikasi"`
	CreatedAt         time.Time `json:"created_at"`
	Barang            []PeminjamanBarangDetail `json:"barang,omitempty"`
}

type PeminjamanBarangDetail struct {
	ID         int    `json:"id"`
	BarangID   int    `json:"barang_id"`
	Barang     Barang `json:"barang"`
	Jumlah     int    `json:"jumlah"`
}

type CreatePeminjamanRequest struct {
	RuanganID       *int                    `json:"ruangan_id"`
	TanggalMulai    string                   `json:"tanggal_mulai"` // ISO 8601 format
	TanggalSelesai  string                   `json:"tanggal_selesai"`
	Keperluan       string                   `json:"keperluan"`
	SuratDigitalURL string                   `json:"surat_digital_url"`
	Barang          []CreatePeminjamanBarang `json:"barang"`
}

type CreatePeminjamanBarang struct {
	BarangID int `json:"barang_id"`
	Jumlah   int `json:"jumlah"`
}

type VerifikasiPeminjamanRequest struct {
	Status            string `json:"status"` // APPROVED atau REJECTED
	CatatanVerifikasi string `json:"catatan_verifikasi"`
}

type JadwalRuanganResponse struct {
	RuanganID   int       `json:"ruangan_id"`
	KodeRuangan string    `json:"kode_ruangan"`
	NamaRuangan string    `json:"nama_ruangan"`
	TanggalMulai time.Time `json:"tanggal_mulai"`
	TanggalSelesai time.Time `json:"tanggal_selesai"`
	Status       string    `json:"status"`
	Peminjam     string    `json:"peminjam"`
}

