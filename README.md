# Sistem Informasi Peminjaman Sarana dan Prasarana Kampus

Sistem web-based untuk mengelola peminjaman ruangan dan barang di kampus dengan role-based access control.

## Teknologi

- **Backend**: Golang native (tanpa framework besar)
- **Database**: PostgreSQL (Supabase)
- **Frontend**: HTML5, TailwindCSS, Vanilla JavaScript
- **Authentication**: JWT (JSON Web Token)

## Struktur Proyek

```
proyek-2-golang/
├── cmd/
│   └── server/
│       └── main.go              # Entry point aplikasi
├── internal/
│   ├── config/
│   │   └── config.go            # Konfigurasi (env variables)
│   ├── db/
│   │   └── db.go                # Koneksi database
│   └── router/
│       └── router.go           # Routing HTTP
├── models/                      # Struct domain models
├── repositories/                # Data access layer (CRUD)
├── services/                     # Business logic layer
├── handlers/                    # HTTP handlers (controllers)
├── middleware/                  # JWT auth & role-based access
├── migrations/                   # SQL migration files
└── view/                        # Frontend HTML/CSS/JS
    ├── index.html               # Halaman login
    ├── dashboard-mahasiswa.html
    ├── dashboard-sarpras.html
    ├── dashboard-security.html
    └── js/
        └── api.js               # Utility functions untuk API calls
```

## Setup

### 1. Install Dependencies

```bash
go mod download
```

### 2. Setup Database

Jalankan SQL migration di Supabase:

```bash
# Buka Supabase Dashboard > SQL Editor
# Copy dan jalankan isi file migrations/001_init_schema.sql
```

### 3. Environment Variables

Buat file `.env` di root project:

```env
DATABASE_URL=postgresql://user:password@host:port/database
PORT=8000
JWT_SECRET=your-secret-key-change-in-production
```

### 4. Run Server

```bash
go run cmd/server/main.go
```

Server akan berjalan di `http://localhost:8000`

## API Endpoints

### Authentication
- `POST /api/auth/login` - Login user

### Master Data (SARPRAS only)
- `GET /api/ruangan` - List semua ruangan
- `GET /api/ruangan/{id}` - Detail ruangan
- `POST /api/ruangan/create` - Tambah ruangan
- `PUT /api/ruangan/{id}` - Update ruangan
- `DELETE /api/ruangan/{id}` - Hapus ruangan

- `GET /api/barang` - List semua barang
- `GET /api/barang/{id}` - Detail barang
- `POST /api/barang/create` - Tambah barang
- `PUT /api/barang/{id}` - Update barang
- `DELETE /api/barang/{id}` - Hapus barang

### Peminjaman
- `POST /api/peminjaman` - Buat pengajuan peminjaman (MAHASISWA)
- `GET /api/peminjaman/me` - List peminjaman milik user (MAHASISWA)
- `GET /api/peminjaman/{id}` - Detail peminjaman
- `GET /api/peminjaman/pending` - List pengajuan pending (SARPRAS)
- `POST /api/peminjaman/{id}/verifikasi` - Verifikasi peminjaman (SARPRAS)
- `GET /api/jadwal-ruangan` - Jadwal ruangan (public)
- `GET /api/jadwal-aktif` - Jadwal aktif untuk Security (SECURITY)

### Kehadiran
- `POST /api/kehadiran` - Catat kehadiran peminjam (SECURITY)

### Notifikasi
- `GET /api/notifikasi/me` - List notifikasi user
- `GET /api/notifikasi/count` - Jumlah notifikasi belum dibaca
- `PATCH /api/notifikasi/{id}/dibaca` - Tandai notifikasi sebagai dibaca

### Laporan
- `GET /api/laporan/peminjaman` - Laporan peminjaman (SARPRAS)
- `GET /api/laporan/kehadiran` - Laporan kehadiran

## Role & Akses

1. **MAHASISWA**: 
   - Melihat jadwal ruangan
   - Mengajukan peminjaman
   - Melihat riwayat peminjaman

2. **SARPRAS**:
   - Semua akses MAHASISWA
   - Kelola master data (ruangan, barang)
   - Verifikasi pengajuan peminjaman
   - Melihat laporan

3. **SECURITY**:
   - Melihat jadwal peminjaman aktif
   - Mencatat kehadiran peminjam

4. **ADMIN**:
   - Semua akses

## Workflow

1. **Mahasiswa** mengajukan peminjaman dengan mengupload surat digital
2. Sistem memvalidasi data dan menyimpan dengan status `PENDING`
3. **Petugas Sarpras** melihat daftar pengajuan pending dan melakukan verifikasi
4. Jika disetujui, status menjadi `APPROVED` dan jadwal dikirim ke Security
5. **Petugas Security** melihat jadwal aktif dan mencatat kehadiran pada hari H

## Catatan Penting

- JWT secret harus diubah di production (file `middleware/auth.go` dan `services/auth_service.go`)
- URL surat digital harus diupload ke Supabase Storage terlebih dahulu
- Pastikan database connection string sudah benar di `.env`
- Frontend menggunakan TailwindCSS via CDN (untuk production, gunakan build process)

## Development

Untuk development, pastikan:
1. Database Supabase sudah setup dan migration sudah dijalankan
2. Environment variables sudah dikonfigurasi
3. Server backend berjalan di port 8000
4. Frontend diakses melalui browser (tidak perlu build process)

## License

Proyek 2 - Sistem Informasi Peminjaman Sarana dan Prasarana Kampus
