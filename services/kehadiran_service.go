package services

import (
	"errors"
	"backend-sarpras/models"
	"backend-sarpras/repositories"
)

type KehadiranService struct {
	KehadiranRepo  *repositories.KehadiranRepository
	PeminjamanRepo *repositories.PeminjamanRepository
	LogRepo        *repositories.LogAktivitasRepository
}

func NewKehadiranService(
	kehadiranRepo *repositories.KehadiranRepository,
	peminjamanRepo *repositories.PeminjamanRepository,
	logRepo *repositories.LogAktivitasRepository,
) *KehadiranService {
	return &KehadiranService{
		KehadiranRepo:  kehadiranRepo,
		PeminjamanRepo: peminjamanRepo,
		LogRepo:        logRepo,
	}
}

func (s *KehadiranService) CreateKehadiran(req *models.CreateKehadiranRequest, securityID int) error {
	// Validasi peminjaman
	peminjaman, err := s.PeminjamanRepo.GetByID(req.PeminjamanID)
	if err != nil {
		return err
	}
	if peminjaman == nil {
		return errors.New("peminjaman tidak ditemukan")
	}

	if peminjaman.Status != "APPROVED" {
		return errors.New("peminjaman belum disetujui")
	}

	// Cek apakah sudah ada kehadiran
	existing, _ := s.KehadiranRepo.GetByPeminjamanID(req.PeminjamanID)
	if existing != nil {
		return errors.New("kehadiran sudah pernah diisi")
	}

	// Validasi status kehadiran
	if req.StatusKehadiran != "HADIR" && req.StatusKehadiran != "TIDAK_HADIR" && req.StatusKehadiran != "TERLAMBAT" {
		return errors.New("status kehadiran tidak valid")
	}

	// Buat kehadiran
	kehadiran := &models.KehadiranPeminjam{
		PeminjamanID:    req.PeminjamanID,
		SecurityID:      &securityID,
		StatusKehadiran: req.StatusKehadiran,
		Catatan:         req.Catatan,
	}

	err = s.KehadiranRepo.Create(kehadiran)
	if err != nil {
		return err
	}

	// Buat log aktivitas
	s.LogRepo.Create(&models.LogAktivitas{
		UserID:       &securityID,
		PeminjamanID: &req.PeminjamanID,
		Aksi:         "UPDATE_KEHADIRAN",
		Keterangan:   "Status kehadiran: " + req.StatusKehadiran,
	})

	return nil
}

