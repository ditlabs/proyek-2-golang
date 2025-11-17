package repositories

import (
	"database/sql"
	"backend-sarpras/models"
)

type OrganisasiRepository struct {
	DB *sql.DB
}

func NewOrganisasiRepository(db *sql.DB) *OrganisasiRepository {
	return &OrganisasiRepository{DB: db}
}

func (r *OrganisasiRepository) GetAll() ([]models.Organisasi, error) {
	query := `SELECT id, nama, jenis, kontak, created_at FROM organisasi ORDER BY nama`
	rows, err := r.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orgs []models.Organisasi
	for rows.Next() {
		var o models.Organisasi
		err := rows.Scan(&o.ID, &o.Nama, &o.Jenis, &o.Kontak, &o.CreatedAt)
		if err != nil {
			return nil, err
		}
		orgs = append(orgs, o)
	}
	return orgs, nil
}

func (r *OrganisasiRepository) GetByID(id int) (*models.Organisasi, error) {
	org := &models.Organisasi{}
	query := `SELECT id, nama, jenis, kontak, created_at FROM organisasi WHERE id = $1`
	err := r.DB.QueryRow(query, id).Scan(&org.ID, &org.Nama, &org.Jenis, &org.Kontak, &org.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return org, err
}

func (r *OrganisasiRepository) Create(org *models.Organisasi) error {
	query := `INSERT INTO organisasi (nama, jenis, kontak) VALUES ($1, $2, $3) RETURNING id, created_at`
	return r.DB.QueryRow(query, org.Nama, org.Jenis, org.Kontak).Scan(&org.ID, &org.CreatedAt)
}

