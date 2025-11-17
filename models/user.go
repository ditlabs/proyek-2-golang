package models

import "time"

type User struct {
	ID           int         `json:"id"`
	Nama         string      `json:"nama"`
	Email        string      `json:"email"`
	PasswordHash string      `json:"-"` // tidak di-expose ke JSON
	Role         string      `json:"role"`
	OrganisasiID *int        `json:"organisasi_id"`
	Organisasi   *Organisasi `json:"organisasi,omitempty"`
	CreatedAt    time.Time   `json:"created_at"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token string `json:"token"`
	User  User   `json:"user"`
}

type RegisterRequest struct {
	Nama         string `json:"nama"`
	Email        string `json:"email"`
	Password     string `json:"password"`
	Role         string `json:"role"`
	OrganisasiID *int   `json:"organisasi_id"`
}
