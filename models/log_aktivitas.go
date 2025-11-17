package models

import "time"

type LogAktivitas struct {
	ID          int        `json:"id"`
	UserID      *int       `json:"user_id"`
	User        *User      `json:"user,omitempty"`
	PeminjamanID *int      `json:"peminjaman_id"`
	Peminjaman  *Peminjaman `json:"peminjaman,omitempty"`
	Aksi        string     `json:"aksi"`
	Keterangan  string     `json:"keterangan"`
	Waktu       time.Time  `json:"waktu"`
}

