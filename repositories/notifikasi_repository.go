package repositories

import (
	"database/sql"
	"backend-sarpras/models"
)

type NotifikasiRepository struct {
	DB *sql.DB
}

func NewNotifikasiRepository(db *sql.DB) *NotifikasiRepository {
	return &NotifikasiRepository{DB: db}
}

func (r *NotifikasiRepository) Create(notifikasi *models.Notifikasi) error {
	query := `
		INSERT INTO notifikasi (peminjaman_id, penerima_id, jenis_notifikasi, pesan, status)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, waktu_kirim
	`
	err := r.DB.QueryRow(
		query,
		notifikasi.PeminjamanID,
		notifikasi.PenerimaID,
		notifikasi.JenisNotifikasi,
		notifikasi.Pesan,
		notifikasi.Status,
	).Scan(&notifikasi.ID, &notifikasi.WaktuKirim)
	return err
}

func (r *NotifikasiRepository) GetByPenerimaID(penerimaID int) ([]models.Notifikasi, error) {
	query := `
		SELECT id, peminjaman_id, penerima_id, jenis_notifikasi, pesan, waktu_kirim, status
		FROM notifikasi
		WHERE penerima_id = $1
		ORDER BY waktu_kirim DESC
		LIMIT 50
	`
	rows, err := r.DB.Query(query, penerimaID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var notifikasi []models.Notifikasi
	for rows.Next() {
		var n models.Notifikasi
		err := rows.Scan(
			&n.ID,
			&n.PeminjamanID,
			&n.PenerimaID,
			&n.JenisNotifikasi,
			&n.Pesan,
			&n.WaktuKirim,
			&n.Status,
		)
		if err != nil {
			return nil, err
		}
		notifikasi = append(notifikasi, n)
	}
	return notifikasi, nil
}

func (r *NotifikasiRepository) GetUnreadCount(penerimaID int) (int, error) {
	var count int
	query := `SELECT COUNT(*) FROM notifikasi WHERE penerima_id = $1 AND status = 'TERKIRIM'`
	err := r.DB.QueryRow(query, penerimaID).Scan(&count)
	return count, err
}

func (r *NotifikasiRepository) MarkAsRead(id int, penerimaID int) error {
	query := `UPDATE notifikasi SET status = 'DIBACA' WHERE id = $1 AND penerima_id = $2`
	_, err := r.DB.Exec(query, id, penerimaID)
	return err
}

