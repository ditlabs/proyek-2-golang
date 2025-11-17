package repositories

import (
	"backend-sarpras/models"
	"database/sql"
	"time"
)

type PeminjamanRepository struct {
	DB *sql.DB
}

func NewPeminjamanRepository(db *sql.DB) *PeminjamanRepository {
	return &PeminjamanRepository{DB: db}
}

func (r *PeminjamanRepository) Create(peminjaman *models.Peminjaman) error {
	query := `
		INSERT INTO peminjaman (peminjam_id, ruangan_id, tanggal_mulai, tanggal_selesai, keperluan, status, surat_digital_url)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, created_at
	`
	err := r.DB.QueryRow(
		query,
		peminjaman.PeminjamID,
		peminjaman.RuanganID,
		peminjaman.TanggalMulai,
		peminjaman.TanggalSelesai,
		peminjaman.Keperluan,
		peminjaman.Status,
		peminjaman.SuratDigitalURL,
	).Scan(&peminjaman.ID, &peminjaman.CreatedAt)
	return err
}

func (r *PeminjamanRepository) CreatePeminjamanBarang(peminjamanID int, barangID int, jumlah int) error {
	query := `INSERT INTO peminjaman_barang (peminjaman_id, barang_id, jumlah) VALUES ($1, $2, $3)`
	_, err := r.DB.Exec(query, peminjamanID, barangID, jumlah)
	return err
}

func (r *PeminjamanRepository) GetByID(id int) (*models.Peminjaman, error) {
	p := &models.Peminjaman{}
	query := `
		SELECT id, peminjam_id, ruangan_id, tanggal_mulai, tanggal_selesai, keperluan, status,
		       surat_digital_url, verified_by, verified_at, COALESCE(catatan_verifikasi, ''), created_at
		FROM peminjaman
		WHERE id = $1
	`
	err := r.DB.QueryRow(query, id).Scan(
		&p.ID,
		&p.PeminjamID,
		&p.RuanganID,
		&p.TanggalMulai,
		&p.TanggalSelesai,
		&p.Keperluan,
		&p.Status,
		&p.SuratDigitalURL,
		&p.VerifiedBy,
		&p.VerifiedAt,
		&p.CatatanVerifikasi,
		&p.CreatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return p, err
}

func (r *PeminjamanRepository) GetByPeminjamID(peminjamID int) ([]models.Peminjaman, error) {
	query := `
		SELECT id, peminjam_id, ruangan_id, tanggal_mulai, tanggal_selesai, keperluan, status,
		       surat_digital_url, verified_by, verified_at, COALESCE(catatan_verifikasi, ''), created_at
		FROM peminjaman
		WHERE peminjam_id = $1
		ORDER BY created_at DESC
	`
	rows, err := r.DB.Query(query, peminjamID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var peminjaman []models.Peminjaman
	for rows.Next() {
		var p models.Peminjaman
		err := rows.Scan(
			&p.ID,
			&p.PeminjamID,
			&p.RuanganID,
			&p.TanggalMulai,
			&p.TanggalSelesai,
			&p.Keperluan,
			&p.Status,
			&p.SuratDigitalURL,
			&p.VerifiedBy,
			&p.VerifiedAt,
			&p.CatatanVerifikasi,
			&p.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		peminjaman = append(peminjaman, p)
	}
	return peminjaman, nil
}

func (r *PeminjamanRepository) GetPending() ([]models.Peminjaman, error) {
	query := `
		SELECT id, peminjam_id, ruangan_id, tanggal_mulai, tanggal_selesai, keperluan, status,
		       surat_digital_url, verified_by, verified_at, COALESCE(catatan_verifikasi, ''), created_at
		FROM peminjaman
		WHERE status = 'PENDING'
		ORDER BY created_at ASC
	`
	rows, err := r.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var peminjaman []models.Peminjaman
	for rows.Next() {
		var p models.Peminjaman
		err := rows.Scan(
			&p.ID,
			&p.PeminjamID,
			&p.RuanganID,
			&p.TanggalMulai,
			&p.TanggalSelesai,
			&p.Keperluan,
			&p.Status,
			&p.SuratDigitalURL,
			&p.VerifiedBy,
			&p.VerifiedAt,
			&p.CatatanVerifikasi,
			&p.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		peminjaman = append(peminjaman, p)
	}
	return peminjaman, nil
}

func (r *PeminjamanRepository) GetJadwalRuangan(start, end time.Time) ([]models.JadwalRuanganResponse, error) {
	query := `
		SELECT p.ruangan_id, r.kode_ruangan, r.nama_ruangan, p.tanggal_mulai, p.tanggal_selesai, p.status, u.nama
		FROM peminjaman p
		JOIN ruangan r ON p.ruangan_id = r.id
		JOIN users u ON p.peminjam_id = u.id
		WHERE p.ruangan_id IS NOT NULL
		  AND p.status = 'APPROVED'
		  AND p.tanggal_mulai >= $1
		  AND p.tanggal_selesai <= $2
		ORDER BY p.tanggal_mulai
	`
	rows, err := r.DB.Query(query, start, end)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var jadwal []models.JadwalRuanganResponse
	for rows.Next() {
		var j models.JadwalRuanganResponse
		err := rows.Scan(
			&j.RuanganID,
			&j.KodeRuangan,
			&j.NamaRuangan,
			&j.TanggalMulai,
			&j.TanggalSelesai,
			&j.Status,
			&j.Peminjam,
		)
		if err != nil {
			return nil, err
		}
		jadwal = append(jadwal, j)
	}
	return jadwal, nil
}

func (r *PeminjamanRepository) GetJadwalAktif(start, end time.Time) ([]models.Peminjaman, error) {
	query := `
		SELECT id, peminjam_id, ruangan_id, tanggal_mulai, tanggal_selesai, keperluan, status,
		       surat_digital_url, verified_by, verified_at, COALESCE(catatan_verifikasi, ''), created_at
		FROM peminjaman
		WHERE status = 'APPROVED'
		  AND ruangan_id IS NOT NULL
		  AND tanggal_mulai >= $1
		  AND tanggal_selesai <= $2
		ORDER BY tanggal_mulai
	`
	rows, err := r.DB.Query(query, start, end)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var peminjaman []models.Peminjaman
	for rows.Next() {
		var p models.Peminjaman
		err := rows.Scan(
			&p.ID,
			&p.PeminjamID,
			&p.RuanganID,
			&p.TanggalMulai,
			&p.TanggalSelesai,
			&p.Keperluan,
			&p.Status,
			&p.SuratDigitalURL,
			&p.VerifiedBy,
			&p.VerifiedAt,
			&p.CatatanVerifikasi,
			&p.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		peminjaman = append(peminjaman, p)
	}
	return peminjaman, nil
}

func (r *PeminjamanRepository) UpdateStatus(id int, status string, verifiedBy *int, catatan string) error {
	query := `
		UPDATE peminjaman
		SET status = $1, verified_by = $2, verified_at = NOW(), catatan_verifikasi = $3
		WHERE id = $4
	`
	_, err := r.DB.Exec(query, status, verifiedBy, catatan, id)
	return err
}

func (r *PeminjamanRepository) GetPeminjamanBarang(peminjamanID int) ([]models.PeminjamanBarangDetail, error) {
	query := `
		SELECT pb.id, pb.barang_id, pb.jumlah,
		       b.kode_barang, b.nama_barang, b.deskripsi, b.jumlah_total, b.jumlah_tersedia, b.ruangan_id
		FROM peminjaman_barang pb
		JOIN barang b ON pb.barang_id = b.id
		WHERE pb.peminjaman_id = $1
	`
	rows, err := r.DB.Query(query, peminjamanID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []models.PeminjamanBarangDetail
	for rows.Next() {
		var item models.PeminjamanBarangDetail
		var barang models.Barang
		err := rows.Scan(
			&item.ID,
			&item.BarangID,
			&item.Jumlah,
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
		item.Barang = barang
		items = append(items, item)
	}
	return items, nil
}

func (r *PeminjamanRepository) GetLaporan(start, end time.Time, status string) ([]models.Peminjaman, error) {
	query := `
		SELECT id, peminjam_id, ruangan_id, tanggal_mulai, tanggal_selesai, keperluan, status,
		       surat_digital_url, verified_by, verified_at, COALESCE(catatan_verifikasi, ''), created_at
		FROM peminjaman
		WHERE created_at >= $1 AND created_at <= $2
	`
	args := []interface{}{start, end}
	if status != "" {
		query += " AND status = $3"
		args = append(args, status)
	}
	query += " ORDER BY created_at DESC"

	rows, err := r.DB.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var peminjaman []models.Peminjaman
	for rows.Next() {
		var p models.Peminjaman
		err := rows.Scan(
			&p.ID,
			&p.PeminjamID,
			&p.RuanganID,
			&p.TanggalMulai,
			&p.TanggalSelesai,
			&p.Keperluan,
			&p.Status,
			&p.SuratDigitalURL,
			&p.VerifiedBy,
			&p.VerifiedAt,
			&p.CatatanVerifikasi,
			&p.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		peminjaman = append(peminjaman, p)
	}
	return peminjaman, nil
}
