# Panduan Setup & Jalankan Backend API-Only Secara Lokal

## Prasyarat

Pastikan Anda sudah install:
- **Go 1.25+** — download dari [golang.org](https://golang.org/dl)
- **PostgreSQL 12+** atau **Supabase account** untuk database remote
- **Git** untuk version control
- **Postman** atau **cURL** untuk testing API (opsional)

Verifikasi instalasi:
```powershell
go version
psql --version
git --version
```

---

## 1. Clone & Setup Project Lokal

```powershell
cd 'd:\Kuliah\Proyek 2\'
git clone https://github.com/ditlabs/proyek-2-golang.git
cd proyek-2-golang
git checkout refactor/api-only
```

---

## 2. Download Dependencies Go

```powershell
go mod download
go mod tidy
```

Verifikasi dependencies terpasang:
```powershell
go mod graph
```

---

## 3. Siapkan Database

### Opsi A: PostgreSQL Lokal

**Windows:**
1. Download & install PostgreSQL dari [postgresql.org](https://www.postgresql.org/download/windows/)
2. Saat instalasi, ingat password superuser (`postgres`)
3. Buka **pgAdmin** (UI manager PostgreSQL) atau gunakan terminal `psql`
4. Buat database baru:

```sql
-- Buka terminal psql
psql -U postgres

-- Jalankan di dalam psql:
CREATE DATABASE sarpras_db;
\q
```

5. Catat connection string lokal:
```
postgresql://postgres:PASSWORD@localhost:5432/sarpras_db?sslmode=disable
```

Ganti `PASSWORD` dengan password superuser Anda.

### Opsi B: Supabase (Cloud PostgreSQL)

1. Daftar di [supabase.com](https://supabase.com)
2. Buat project baru
3. Di tab **Settings → Database**, salin connection string PostgreSQL:
```
postgresql://postgres.XXXXX:PASSWORD@db.XXXXX.supabase.co:5432/postgres?sslmode=require
```

---

## 4. Jalankan Migrations

Gunakan migration SQL untuk menginisialisasi schema:

**Windows PowerShell:**
```powershell
# Ganti DATABASE_URL dengan connection string Anda
$env:DATABASE_URL = "postgresql://postgres:PASSWORD@localhost:5432/sarpras_db?sslmode=disable"

# Jalankan migration menggunakan psql
psql $env:DATABASE_URL -f migrations/001_init_schema.sql
```

**Atau jalankan SQL langsung via psql:**
```powershell
psql -U postgres -d sarpras_db -f 'migrations/001_init_schema.sql'
```

Verifikasi tabel berhasil dibuat:
```powershell
psql -U postgres -d sarpras_db -c "\dt"
```

---

## 5. Setup Environment Variables

Buat file `.env` di root folder project:

```powershell
# di folder: d:\Kuliah\Proyek 2\proyek-2-golang\.env

New-Item -Path '.\.env' -ItemType File -Force | Out-Null

# Isi dengan:
# DATABASE_URL=postgresql://postgres:PASSWORD@localhost:5432/sarpras_db?sslmode=disable
# PORT=8000
# JWT_SECRET=your-super-secret-key-min-32-chars-long!
# CORS_ALLOWED_ORIGIN=http://127.0.0.1:5500
```

**Atau edit file dengan text editor (notepad/VS Code):**
```
DATABASE_URL=postgresql://postgres:PASSWORD@localhost:5432/sarpras_db?sslmode=disable
PORT=8000
JWT_SECRET=your-super-secret-key-min-32-chars-long!
CORS_ALLOWED_ORIGIN=http://127.0.0.1:5500
```

**Catatan:**
- `DATABASE_URL`: Sesuaikan dengan connection string Anda (lokal atau Supabase)
- `PORT`: Default 8000, bisa diganti sesuai kebutuhan
- `JWT_SECRET`: Key untuk enkripsi JWT token, gunakan string panjang & random
- `CORS_ALLOWED_ORIGIN`: URL frontend (lokal: `http://127.0.0.1:5500` atau production: `https://username.github.io`)

---

## 6. Jalankan Server Backend

```powershell
cd 'd:\Kuliah\Proyek 2\proyek-2-golang'

# Cara 1: go run (development)
go run ./cmd/server/main.go

# Cara 2: build dulu, lalu jalankan (lebih cepat untuk reload banyak)
go build -o sarpras-api.exe ./cmd/server
./sarpras-api.exe
```

**Output yang diharapkan:**
```
database connected
Server running on http://localhost:8000
```

Jika ada error, cek:
- ✅ `.env` file ada dan env vars ter-load
- ✅ Database connection string benar
- ✅ Database sudah di-migrate dengan SQL
- ✅ Port 8000 tidak digunakan aplikasi lain

---

## 7. Test API Backend Lokal

**Buka terminal baru** (jangan tutup terminal server) dan jalankan test:

### Test 1: Health Check (jika ada endpoint `/api/health`)

```powershell
curl http://localhost:8000/api/health
```

Respons yang diharapkan: `{"status":"ok"}` atau similar

### Test 2: Login (POST /api/auth/login)

```powershell
# Pertama, pastikan ada user di database atau insert manual:
# Jika belum ada user, insert satu:

psql -U postgres -d sarpras_db -c "
INSERT INTO organisasi (nama, jenis) VALUES ('Tes Org', 'UKM');
INSERT INTO users (nama, email, password_hash, role, organisasi_id) 
VALUES ('Test User', 'test@example.com', 'hashed_pwd', 'MAHASISWA', 1);
"

# Kemudian, coba login:
$loginPayload = @{
    email = "test@example.com"
    password = "password123"
} | ConvertTo-Json

curl -X POST http://localhost:8000/api/auth/login `
  -Headers @{'Content-Type'='application/json'} `
  -Body $loginPayload
```

**Respons yang diharapkan:**
```json
{
  "token": "eyJhbGciOiJIUzI1NiIs...",
  "user": {
    "id": 1,
    "nama": "Test User",
    "email": "test@example.com",
    "role": "MAHASISWA"
  }
}
```

### Test 3: Fetch Protected Endpoint (GET /api/ruangan) dengan Token

```powershell
# Gunakan token dari response login di atas
$token = "eyJhbGciOiJIUzI1NiIs..."

curl -X GET http://localhost:8000/api/ruangan `
  -Headers @{
    'Authorization' = "Bearer $token"
    'Content-Type' = 'application/json'
  }
```

**Respons yang diharapkan:**
```json
[
  {
    "id": 1,
    "kode_ruangan": "R-001",
    "nama_ruangan": "Ruang Pertemuan A",
    "lokasi": "Gedung A",
    "kapasitas": 30
  }
]
```

---

## 8. Development Workflow (Hot Reload)

Agar server auto-reload saat file Go diubah, gunakan **Air**:

```powershell
# Install air
go install github.com/cosmtrek/air@latest

# Jalankan dengan air (auto-reload)
cd 'd:\Kuliah\Proyek 2\proyek-2-golang'
air
```

Air akan otomatis rebuild & restart server saat `.go` file berubah.

---

## 9. Connect Frontend Lokal ke Backend

Jika Anda sudah setup frontend (`frontend-sarpras`), pastikan:

1. **Frontend folder terpisah** (bukan dalam `proyek-2-golang`)
2. **Jalankan frontend dengan Live Server** di VS Code:
   - Klik kanan `index.html` → "Open with Live Server"
   - Frontend biasanya berjalan di `http://127.0.0.1:5500`

3. **Di `frontend-sarpras/assets/js/config.js`**, pastikan `API_BASE_URL` mengarah ke backend lokal:

```javascript
// Contoh config.js
const getApiBaseUrl = () => {
  // Jika ada override dari window
  if (window.API_BASE_URL) {
    return window.API_BASE_URL;
  }
  // Development lokal
  if (location.hostname === 'localhost' || location.hostname === '127.0.0.1') {
    return 'http://localhost:8000/api';
  }
  // Production
  return 'https://api.sarpras.example.com/api';
};

export const API_BASE_URL = getApiBaseUrl();
```

4. **Backend CORS** harus allow origin frontend:
   - File: `middleware/cors.go`
   - Env var: `CORS_ALLOWED_ORIGIN=http://127.0.0.1:5500`

---

## 10. Troubleshooting

| Error | Solusi |
|-------|--------|
| `DATABASE_URL is required` | Pastikan `.env` ada dan `DATABASE_URL` terisi |
| `failed to open db` | Cek connection string, pastikan database server running |
| `connection refused on :8000` | Port 8000 sudah digunakan. Ganti `PORT` di `.env` atau tutup app lain |
| `Unauthorized (401)` | Token expired atau tidak valid. Login ulang dan gunakan token terbaru |
| `CORS error` di frontend | Periksa `CORS_ALLOWED_ORIGIN` di `.env` backend sesuai dengan URL frontend |
| `Module not found` | Jalankan `go mod tidy` & `go mod download` |

---

## 11. Struktur Folder

```
proyek-2-golang/
├── cmd/server/main.go          ← Entry point aplikasi
├── internal/
│   ├── config/config.go        ← Load env vars
│   ├── db/db.go                ← Database connection
│   └── router/router.go        ← Register routes & handlers
├── middleware/
│   ├── auth.go                 ← JWT authentication
│   └── cors.go                 ← CORS middleware
├── handlers/                   ← HTTP handlers
├── services/                   ← Business logic
├── repositories/               ← Database queries
├── models/                     ← Data models
├── migrations/001_init_schema.sql
├── go.mod & go.sum
├── .env                        ← Environment variables (jangan push ke git!)
└── README.md
```

---

## 12. Commit & Push ke Remote (Optional)

Jika sudah siap, commit perubahan lokal ke branch `refactor/api-only`:

```powershell
cd 'd:\Kuliah\Proyek 2\proyek-2-golang'
git add -A
git commit -m "docs: add local setup guide"
git push origin refactor/api-only
```

---

## Checklist Verifikasi

- [ ] Go 1.25+ ter-install
- [ ] PostgreSQL atau Supabase siap
- [ ] `.env` file dibuat dengan env vars yang benar
- [ ] Database migrated (schema table ada)
- [ ] `go mod download` sukses
- [ ] `go run ./cmd/server/main.go` berjalan tanpa error
- [ ] `curl http://localhost:8000/api/health` (atau endpoint lain) merespons
- [ ] Token JWT bisa didapat dari login endpoint
- [ ] Protected endpoint bisa diakses dengan token
- [ ] Frontend lokal terhubung ke backend tanpa CORS error

---

## Next Steps

1. **Local verification**: Pastikan semua test di atas berhasil
2. **Deploy frontend**: Push `frontend-sarpras` repo ke GitHub, enable GitHub Pages
3. **Deploy backend**: Gunakan Render/Railway/VPS dan set production env vars
4. **Setup CI/CD**: GitHub Actions untuk auto-build & auto-deploy

Pertanyaan? Tanya saja! 🚀
