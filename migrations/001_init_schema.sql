-- Migration: Inisialisasi Schema Database
-- Sistem Informasi Peminjaman Sarana dan Prasarana Kampus

-- Tabel organisasi
CREATE TABLE IF NOT EXISTS organisasi (
    id SERIAL PRIMARY KEY,
    nama VARCHAR(255) NOT NULL,
    jenis VARCHAR(100) NOT NULL, -- Himpunan, UKM, BEM, dll
    kontak VARCHAR(255),
    created_at TIMESTAMP DEFAULT NOW()
);

-- Tabel users
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    nama VARCHAR(255) NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    role VARCHAR(50) NOT NULL, -- MAHASISWA, SARPRAS, SECURITY, ADMIN
    organisasi_id INTEGER REFERENCES organisasi(id) ON DELETE SET NULL,
    created_at TIMESTAMP DEFAULT NOW()
);

-- Tabel ruangan
CREATE TABLE IF NOT EXISTS ruangan (
    id SERIAL PRIMARY KEY,
    kode_ruangan VARCHAR(100) UNIQUE NOT NULL,
    nama_ruangan VARCHAR(255) NOT NULL,
    lokasi VARCHAR(255),
    kapasitas INTEGER,
    deskripsi TEXT
);

-- Tabel barang
CREATE TABLE IF NOT EXISTS barang (
    id SERIAL PRIMARY KEY,
    kode_barang VARCHAR(100) UNIQUE NOT NULL,
    nama_barang VARCHAR(255) NOT NULL,
    deskripsi TEXT,
    jumlah_total INTEGER NOT NULL DEFAULT 0,
    jumlah_tersedia INTEGER NOT NULL DEFAULT 0,
    ruangan_id INTEGER REFERENCES ruangan(id) ON DELETE SET NULL
);

-- Tabel peminjaman
CREATE TABLE IF NOT EXISTS peminjaman (
    id SERIAL PRIMARY KEY,
    peminjam_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    ruangan_id INTEGER REFERENCES ruangan(id) ON DELETE SET NULL,
    tanggal_mulai TIMESTAMP NOT NULL,
    tanggal_selesai TIMESTAMP NOT NULL,
    keperluan TEXT,
    status VARCHAR(50) NOT NULL DEFAULT 'PENDING', -- PENDING, APPROVED, REJECTED, CANCELLED, FINISHED
    surat_digital_url TEXT NOT NULL,
    verified_by INTEGER REFERENCES users(id) ON DELETE SET NULL,
    verified_at TIMESTAMP,
    catatan_verifikasi TEXT,
    created_at TIMESTAMP DEFAULT NOW()
);

-- Tabel peminjaman_barang
CREATE TABLE IF NOT EXISTS peminjaman_barang (
    id SERIAL PRIMARY KEY,
    peminjaman_id INTEGER NOT NULL REFERENCES peminjaman(id) ON DELETE CASCADE,
    barang_id INTEGER NOT NULL REFERENCES barang(id) ON DELETE CASCADE,
    jumlah INTEGER NOT NULL
);

-- Tabel kehadiran_peminjam
CREATE TABLE IF NOT EXISTS kehadiran_peminjam (
    id SERIAL PRIMARY KEY,
    peminjaman_id INTEGER NOT NULL REFERENCES peminjaman(id) ON DELETE CASCADE,
    security_id INTEGER REFERENCES users(id) ON DELETE SET NULL,
    status_kehadiran VARCHAR(50) NOT NULL, -- HADIR, TIDAK_HADIR, TERLAMBAT
    waktu_verifikasi TIMESTAMP DEFAULT NOW(),
    catatan TEXT
);

-- Tabel notifikasi
CREATE TABLE IF NOT EXISTS notifikasi (
    id SERIAL PRIMARY KEY,
    peminjaman_id INTEGER REFERENCES peminjaman(id) ON DELETE SET NULL,
    penerima_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    jenis_notifikasi VARCHAR(100) NOT NULL, -- PENGAJUAN_DIBUAT, STATUS_APPROVED, STATUS_REJECTED, REMINDER_KEHADIRAN
    pesan TEXT NOT NULL,
    waktu_kirim TIMESTAMP DEFAULT NOW(),
    status VARCHAR(50) NOT NULL DEFAULT 'TERKIRIM' -- TERKIRIM, DIBACA
);

-- Tabel log_aktivitas
CREATE TABLE IF NOT EXISTS log_aktivitas (
    id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(id) ON DELETE SET NULL,
    peminjaman_id INTEGER REFERENCES peminjaman(id) ON DELETE SET NULL,
    aksi VARCHAR(100) NOT NULL, -- CREATE_PEMINJAMAN, UPDATE_STATUS, UPDATE_KEHADIRAN, dll
    keterangan TEXT,
    waktu TIMESTAMP DEFAULT NOW()
);

-- Index untuk performa query
CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);
CREATE INDEX IF NOT EXISTS idx_users_role ON users(role);
CREATE INDEX IF NOT EXISTS idx_peminjaman_status ON peminjaman(status);
CREATE INDEX IF NOT EXISTS idx_peminjaman_peminjam ON peminjaman(peminjam_id);
CREATE INDEX IF NOT EXISTS idx_peminjaman_tanggal ON peminjaman(tanggal_mulai, tanggal_selesai);
CREATE INDEX IF NOT EXISTS idx_notifikasi_penerima ON notifikasi(penerima_id);
CREATE INDEX IF NOT EXISTS idx_notifikasi_status ON notifikasi(status);
CREATE INDEX IF NOT EXISTS idx_log_aktivitas_waktu ON log_aktivitas(waktu);

