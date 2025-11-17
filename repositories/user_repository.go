package repositories

import (
	"database/sql"
	"backend-sarpras/models"
)

type UserRepository struct {
	DB *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{DB: db}
}

func (r *UserRepository) Create(user *models.User) error {
	query := `
		INSERT INTO users (nama, email, password_hash, role, organisasi_id)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at
	`
	err := r.DB.QueryRow(
		query,
		user.Nama,
		user.Email,
		user.PasswordHash,
		user.Role,
		user.OrganisasiID,
	).Scan(&user.ID, &user.CreatedAt)
	return err
}

func (r *UserRepository) GetByEmail(email string) (*models.User, error) {
	user := &models.User{}
	query := `
		SELECT id, nama, email, password_hash, role, organisasi_id, created_at
		FROM users
		WHERE email = $1
	`
	err := r.DB.QueryRow(query, email).Scan(
		&user.ID,
		&user.Nama,
		&user.Email,
		&user.PasswordHash,
		&user.Role,
		&user.OrganisasiID,
		&user.CreatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return user, err
}

func (r *UserRepository) GetByID(id int) (*models.User, error) {
	user := &models.User{}
	query := `
		SELECT id, nama, email, password_hash, role, organisasi_id, created_at
		FROM users
		WHERE id = $1
	`
	err := r.DB.QueryRow(query, id).Scan(
		&user.ID,
		&user.Nama,
		&user.Email,
		&user.PasswordHash,
		&user.Role,
		&user.OrganisasiID,
		&user.CreatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return user, err
}

func (r *UserRepository) GetByRole(role string) ([]models.User, error) {
	query := `
		SELECT id, nama, email, role, organisasi_id, created_at
		FROM users
		WHERE role = $1
		ORDER BY nama
	`
	rows, err := r.DB.Query(query, role)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var u models.User
		err := rows.Scan(
			&u.ID,
			&u.Nama,
			&u.Email,
			&u.Role,
			&u.OrganisasiID,
			&u.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	return users, nil
}

