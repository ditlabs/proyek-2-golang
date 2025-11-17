package models

import "time"

type Notifikasi struct {
	ID              int        `json:"id"`
	PeminjamanID    *int       `json:"peminjaman_id"`
	PenerimaID      int        `json:"penerima_id"`
	JenisNotifikasi string     `json:"jenis_notifikasi"`
	Pesan           string     `json:"pesan"`
	WaktuKirim      time.Time  `json:"waktu_kirim"`
	Status          string     `json:"status"`
	Peminjaman      *Peminjaman `json:"peminjaman,omitempty"`
}

