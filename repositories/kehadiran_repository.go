package repositories

import (
	"database/sql"
	"backend-sarpras/models"
)

type KehadiranRepository struct {
	DB *sql.DB
}

func NewKehadiranRepository(db *sql.DB) *KehadiranRepository {
	return &KehadiranRepository{DB: db}
}

func (r *KehadiranRepository) Create(kehadiran *models.KehadiranPeminjam) error {
	query := `
		INSERT INTO kehadiran_peminjam (peminjaman_id, security_id, status_kehadiran, catatan)
		VALUES ($1, $2, $3, $4)
		RETURNING id, waktu_verifikasi
	`
	err := r.DB.QueryRow(
		query,
		kehadiran.PeminjamanID,
		kehadiran.SecurityID,
		kehadiran.StatusKehadiran,
		kehadiran.Catatan,
	).Scan(&kehadiran.ID, &kehadiran.WaktuVerifikasi)
	return err
}

func (r *KehadiranRepository) GetByPeminjamanID(peminjamanID int) (*models.KehadiranPeminjam, error) {
	kehadiran := &models.KehadiranPeminjam{}
	query := `
		SELECT id, peminjaman_id, security_id, status_kehadiran, waktu_verifikasi, catatan
		FROM kehadiran_peminjam
		WHERE peminjaman_id = $1
	`
	err := r.DB.QueryRow(query, peminjamanID).Scan(
		&kehadiran.ID,
		&kehadiran.PeminjamanID,
		&kehadiran.SecurityID,
		&kehadiran.StatusKehadiran,
		&kehadiran.WaktuVerifikasi,
		&kehadiran.Catatan,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return kehadiran, err
}

func (r *KehadiranRepository) GetAll(start, end string) ([]models.KehadiranPeminjam, error) {
	query := `
		SELECT id, peminjaman_id, security_id, status_kehadiran, waktu_verifikasi, catatan
		FROM kehadiran_peminjam
		WHERE waktu_verifikasi >= $1 AND waktu_verifikasi <= $2
		ORDER BY waktu_verifikasi DESC
	`
	rows, err := r.DB.Query(query, start, end)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var kehadiran []models.KehadiranPeminjam
	for rows.Next() {
		var k models.KehadiranPeminjam
		err := rows.Scan(
			&k.ID,
			&k.PeminjamanID,
			&k.SecurityID,
			&k.StatusKehadiran,
			&k.WaktuVerifikasi,
			&k.Catatan,
		)
		if err != nil {
			return nil, err
		}
		kehadiran = append(kehadiran, k)
	}
	return kehadiran, nil
}

