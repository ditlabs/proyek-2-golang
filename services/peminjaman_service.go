package services

import (
	"errors"
	"time"
	"backend-sarpras/models"
	"backend-sarpras/repositories"
)

type PeminjamanService struct {
	PeminjamanRepo *repositories.PeminjamanRepository
	BarangRepo     *repositories.BarangRepository
	NotifikasiRepo *repositories.NotifikasiRepository
	LogRepo        *repositories.LogAktivitasRepository
	UserRepo       *repositories.UserRepository
}

func NewPeminjamanService(
	peminjamanRepo *repositories.PeminjamanRepository,
	barangRepo *repositories.BarangRepository,
	notifikasiRepo *repositories.NotifikasiRepository,
	logRepo *repositories.LogAktivitasRepository,
	userRepo *repositories.UserRepository,
) *PeminjamanService {
	return &PeminjamanService{
		PeminjamanRepo: peminjamanRepo,
		BarangRepo:     barangRepo,
		NotifikasiRepo: notifikasiRepo,
		LogRepo:        logRepo,
		UserRepo:       userRepo,
	}
}

func (s *PeminjamanService) CreatePeminjaman(req *models.CreatePeminjamanRequest, peminjamID int) (*models.Peminjaman, error) {
	// Validasi input
	if req.SuratDigitalURL == "" {
		return nil, errors.New("surat digital wajib diupload")
	}

	tanggalMulai, err := time.Parse(time.RFC3339, req.TanggalMulai)
	if err != nil {
		return nil, errors.New("format tanggal_mulai tidak valid")
	}

	tanggalSelesai, err := time.Parse(time.RFC3339, req.TanggalSelesai)
	if err != nil {
		return nil, errors.New("format tanggal_selesai tidak valid")
	}

	if tanggalSelesai.Before(tanggalMulai) {
		return nil, errors.New("tanggal_selesai harus setelah tanggal_mulai")
	}

	// Validasi stok barang jika ada
	for _, item := range req.Barang {
		barang, err := s.BarangRepo.GetByID(item.BarangID)
		if err != nil {
			return nil, err
		}
		if barang == nil {
			return nil, errors.New("barang tidak ditemukan")
		}
		if barang.JumlahTersedia < item.Jumlah {
			return nil, errors.New("stok barang tidak mencukupi")
		}
	}

	// Buat peminjaman
	peminjaman := &models.Peminjaman{
		PeminjamID:      peminjamID,
		RuanganID:       req.RuanganID,
		TanggalMulai:    tanggalMulai,
		TanggalSelesai:  tanggalSelesai,
		Keperluan:       req.Keperluan,
		Status:          "PENDING",
		SuratDigitalURL: req.SuratDigitalURL,
	}

	err = s.PeminjamanRepo.Create(peminjaman)
	if err != nil {
		return nil, err
	}

	// Simpan barang yang dipinjam
	for _, item := range req.Barang {
		err = s.PeminjamanRepo.CreatePeminjamanBarang(peminjaman.ID, item.BarangID, item.Jumlah)
		if err != nil {
			return nil, err
		}

		// Kurangi stok tersedia
		barang, _ := s.BarangRepo.GetByID(item.BarangID)
		newJumlah := barang.JumlahTersedia - item.Jumlah
		err = s.BarangRepo.UpdateJumlahTersedia(item.BarangID, newJumlah)
		if err != nil {
			return nil, err
		}
	}

	// Buat log aktivitas
	s.LogRepo.Create(&models.LogAktivitas{
		UserID:       &peminjamID,
		PeminjamanID: &peminjaman.ID,
		Aksi:         "CREATE_PEMINJAMAN",
		Keterangan:   "Pengajuan peminjaman baru dibuat",
	})

	// Buat notifikasi untuk Sarpras
	s.NotifikasiRepo.Create(&models.Notifikasi{
		PeminjamanID:    &peminjaman.ID,
		PenerimaID:      0, // TODO: ambil ID petugas sarpras
		JenisNotifikasi: "PENGAJUAN_DIBUAT",
		Pesan:           "Pengajuan peminjaman baru menunggu verifikasi",
		Status:          "TERKIRIM",
	})

	return peminjaman, nil
}

func (s *PeminjamanService) VerifikasiPeminjaman(peminjamanID int, verifierID int, req *models.VerifikasiPeminjamanRequest) error {
	peminjaman, err := s.PeminjamanRepo.GetByID(peminjamanID)
	if err != nil {
		return err
	}
	if peminjaman == nil {
		return errors.New("peminjaman tidak ditemukan")
	}

	if peminjaman.Status != "PENDING" {
		return errors.New("peminjaman sudah diverifikasi")
	}

	status := req.Status
	if status != "APPROVED" && status != "REJECTED" {
		return errors.New("status tidak valid")
	}

	// Update status
	err = s.PeminjamanRepo.UpdateStatus(peminjamanID, status, &verifierID, req.CatatanVerifikasi)
	if err != nil {
		return err
	}

	// Jika ditolak, kembalikan stok barang
	if status == "REJECTED" {
		items, _ := s.PeminjamanRepo.GetPeminjamanBarang(peminjamanID)
		for _, item := range items {
			barang, _ := s.BarangRepo.GetByID(item.BarangID)
			newJumlah := barang.JumlahTersedia + item.Jumlah
			s.BarangRepo.UpdateJumlahTersedia(item.BarangID, newJumlah)
		}
	}

	// Buat log aktivitas
	s.LogRepo.Create(&models.LogAktivitas{
		UserID:       &verifierID,
		PeminjamanID: &peminjamanID,
		Aksi:         "UPDATE_STATUS",
		Keterangan:   "Status peminjaman diubah menjadi " + status,
	})

	// Buat notifikasi untuk peminjam
	pesan := "Pengajuan peminjaman Anda telah " + status
	if status == "REJECTED" {
		pesan += ". Catatan: " + req.CatatanVerifikasi
	}
	s.NotifikasiRepo.Create(&models.Notifikasi{
		PeminjamanID:    &peminjamanID,
		PenerimaID:      peminjaman.PeminjamID,
		JenisNotifikasi: "STATUS_" + status,
		Pesan:           pesan,
		Status:          "TERKIRIM",
	})

	return nil
}

