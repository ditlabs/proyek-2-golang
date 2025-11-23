package models

import "time"

type KehadiranPeminjam struct {
	ID              int         `json:"id"`
	PeminjamanID    int         `json:"peminjaman_id"`
	Peminjaman      *Peminjaman `json:"peminjaman,omitempty"`
	SecurityID      *int        `json:"security_id"`
	Security        *User       `json:"security,omitempty"`
	StatusKehadiran string      `json:"status_kehadiran"` // HADIR, TIDAK_HADIR, TERLAMBAT
	WaktuVerifikasi *time.Time  `json:"waktu_verifikasi"`
	Catatan         string      `json:"catatan"`
}

type CreateKehadiranRequest struct {
	PeminjamanID    int    `json:"peminjaman_id"`
	StatusKehadiran string `json:"status_kehadiran"`
	Catatan         string `json:"catatan"`
}
