package repositories

import (
	"database/sql"
	"backend-sarpras/models"
)

type LogAktivitasRepository struct {
	DB *sql.DB
}

func NewLogAktivitasRepository(db *sql.DB) *LogAktivitasRepository {
	return &LogAktivitasRepository{DB: db}
}

func (r *LogAktivitasRepository) Create(log *models.LogAktivitas) error {
	query := `
		INSERT INTO log_aktivitas (user_id, peminjaman_id, aksi, keterangan)
		VALUES ($1, $2, $3, $4)
		RETURNING id, waktu
	`
	err := r.DB.QueryRow(
		query,
		log.UserID,
		log.PeminjamanID,
		log.Aksi,
		log.Keterangan,
	).Scan(&log.ID, &log.Waktu)
	return err
}

func (r *LogAktivitasRepository) GetAll(filter string) ([]models.LogAktivitas, error) {
	query := `
		SELECT id, user_id, peminjaman_id, aksi, keterangan, waktu
		FROM log_aktivitas
	`
	if filter != "" {
		query += " WHERE aksi = $1"
	}
	query += " ORDER BY waktu DESC LIMIT 100"

	var rows *sql.Rows
	var err error
	if filter != "" {
		rows, err = r.DB.Query(query, filter)
	} else {
		rows, err = r.DB.Query(query)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var logs []models.LogAktivitas
	for rows.Next() {
		var log models.LogAktivitas
		err := rows.Scan(
			&log.ID,
			&log.UserID,
			&log.PeminjamanID,
			&log.Aksi,
			&log.Keterangan,
			&log.Waktu,
		)
		if err != nil {
			return nil, err
		}
		logs = append(logs, log)
	}
	return logs, nil
}

