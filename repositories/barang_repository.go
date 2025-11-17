package repositories

import (
	"database/sql"
	"backend-sarpras/models"
)

type BarangRepository struct {
	DB *sql.DB
}

func NewBarangRepository(db *sql.DB) *BarangRepository {
	return &BarangRepository{DB: db}
}

func (r *BarangRepository) GetAll() ([]models.Barang, error) {
	query := `
		SELECT id, kode_barang, nama_barang, deskripsi, jumlah_total, jumlah_tersedia, ruangan_id
		FROM barang
		ORDER BY kode_barang
	`
	rows, err := r.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var barangs []models.Barang
	for rows.Next() {
		var barang models.Barang
		err := rows.Scan(
			&barang.ID,
			&barang.KodeBarang,
			&barang.NamaBarang,
			&barang.Deskripsi,
			&barang.JumlahTotal,
			&barang.JumlahTersedia,
			&barang.RuanganID,
		)
		if err != nil {
			return nil, err
		}
		barangs = append(barangs, barang)
	}
	return barangs, nil
}

func (r *BarangRepository) GetByID(id int) (*models.Barang, error) {
	barang := &models.Barang{}
	query := `
		SELECT id, kode_barang, nama_barang, deskripsi, jumlah_total, jumlah_tersedia, ruangan_id
		FROM barang
		WHERE id = $1
	`
	err := r.DB.QueryRow(query, id).Scan(
		&barang.ID,
		&barang.KodeBarang,
		&barang.NamaBarang,
		&barang.Deskripsi,
		&barang.JumlahTotal,
		&barang.JumlahTersedia,
		&barang.RuanganID,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return barang, err
}

func (r *BarangRepository) Create(barang *models.Barang) error {
	query := `
		INSERT INTO barang (kode_barang, nama_barang, deskripsi, jumlah_total, jumlah_tersedia, ruangan_id)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id
	`
	// Saat create, jumlah_tersedia = jumlah_total
	barang.JumlahTersedia = barang.JumlahTotal
	return r.DB.QueryRow(
		query,
		barang.KodeBarang,
		barang.NamaBarang,
		barang.Deskripsi,
		barang.JumlahTotal,
		barang.JumlahTersedia,
		barang.RuanganID,
	).Scan(&barang.ID)
}

func (r *BarangRepository) Update(barang *models.Barang) error {
	query := `
		UPDATE barang
		SET nama_barang = $1, deskripsi = $2, jumlah_total = $3, ruangan_id = $4
		WHERE id = $5
	`
	_, err := r.DB.Exec(query, barang.NamaBarang, barang.Deskripsi, barang.JumlahTotal, barang.RuanganID, barang.ID)
	return err
}

func (r *BarangRepository) UpdateJumlahTersedia(id int, jumlah int) error {
	query := `UPDATE barang SET jumlah_tersedia = $1 WHERE id = $2`
	_, err := r.DB.Exec(query, jumlah, id)
	return err
}

func (r *BarangRepository) Delete(id int) error {
	query := `DELETE FROM barang WHERE id = $1`
	_, err := r.DB.Exec(query, id)
	return err
}

