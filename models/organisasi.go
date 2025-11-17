package models

import "time"

type Organisasi struct {
	ID        int       `json:"id"`
	Nama      string    `json:"nama"`
	Jenis     string    `json:"jenis"`
	Kontak    string    `json:"kontak"`
	CreatedAt time.Time `json:"created_at"`
}

