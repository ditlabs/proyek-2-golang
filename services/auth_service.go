package services

import (
	"backend-sarpras/models"
	"backend-sarpras/repositories"
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"time"
)

type AuthService struct {
	UserRepo *repositories.UserRepository
}

func NewAuthService(userRepo *repositories.UserRepository) *AuthService {
	return &AuthService{UserRepo: userRepo}
}

var jwtSecret = []byte("your-secret-key-change-in-production") // TODO: ambil dari env

func (s *AuthService) Login(email, password string) (*models.LoginResponse, error) {
	user, err := s.UserRepo.GetByEmail(email)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("invalid credentials")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	// Generate JWT token
	claims := jwt.MapClaims{
		"user_id": user.ID,
		"email":   user.Email,
		"role":    user.Role,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		return nil, err
	}

	// Hapus password hash dari response
	user.PasswordHash = ""

	return &models.LoginResponse{
		Token: tokenString,
		User:  *user,
	}, nil
}

func (s *AuthService) HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

func (s *AuthService) Register(req *models.RegisterRequest) (*models.User, error) {
	if req.Email == "" || req.Password == "" || req.Nama == "" {
		return nil, errors.New("nama, email, dan password wajib diisi")
	}
	if req.Role == "" {
		req.Role = "MAHASISWA"
	}

	existing, err := s.UserRepo.GetByEmail(req.Email)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, errors.New("email sudah terdaftar")
	}

	passwordHash, err := s.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	user := &models.User{
		Nama:         req.Nama,
		Email:        req.Email,
		PasswordHash: passwordHash,
		Role:         req.Role,
		OrganisasiID: req.OrganisasiID,
	}

	if err := s.UserRepo.Create(user); err != nil {
		return nil, err
	}

	user.PasswordHash = ""
	return user, nil
}
