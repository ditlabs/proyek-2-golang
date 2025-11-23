# Panduan Jalankan Frontend Secara Lokal

Anda sudah berhasil menjalankan backend di `http://localhost:8000` dengan status `ok`. Sekarang saatnya menjalankan frontend (`frontend-sarpras`) yang sudah Anda push ke repository GitHub terpisah.

---

## Prasyarat

Pastikan sudah install:
- **Git** — untuk clone repository
- **VS Code** — text editor & development tools
- **Live Server extension** di VS Code — untuk menjalankan static server

Jika belum install Live Server di VS Code:
1. Buka VS Code
2. Tekan `Ctrl+Shift+X` (Extensions Marketplace)
3. Cari `Live Server`
4. Klik **Install** (publisher: Ritwick Dey)

---

## 1. Clone Frontend Repository

Pertama, Anda perlu clone `frontend-sarpras` repository dari GitHub Anda.

**Buka PowerShell / Command Prompt:**

```powershell
# Pilih lokasi untuk menyimpan frontend (sebaiknya terpisah dari backend)
cd 'd:\Kuliah\Proyek 2\'

# Clone repository frontend
git clone https://github.com/<username>/<frontend-repo-name>.git

# Masuk ke folder frontend
cd frontend-sarpras

# Verifikasi struktur file
ls -R
```

**Output yang diharapkan:**
```
frontend-sarpras/
├── index.html
├── sarpras.html
├── dashboard-mahasiswa.html
├── dashboard-sarpras.html
├── dashboard-security.html
├── jadwal-ruangan.html
├── kelola-barang.html
├── kelola-ruangan.html
├── laporan-peminjaman.html
├── pengajuan-peminjaman.html
├── register.html
├── riwayat-peminjaman.html
├── rooms.html
├── verifikasi-peminjaman.html
├── assets/
│   ├── js/
│   │   ├── config.js
│   │   ├── api.js
│   │   ├── api.legacy.js
│   │   ├── auth.js
│   │   └── dashboard-sarpras.js
│   └── styles/ (jika ada custom CSS)
└── README.md
```

---

## 2. Buka Frontend di VS Code

```powershell
# Buka folder frontend dengan VS Code
code .
```

Atau manual:
1. Buka VS Code
2. `File → Open Folder`
3. Pilih folder `frontend-sarpras`

---

## 3. Jalankan Frontend dengan Live Server

### Cara A: Klik Kanan pada `index.html`

1. Di file explorer VS Code, klik kanan file `index.html`
2. Pilih **"Open with Live Server"**

Maka browser akan terbuka otomatis di `http://127.0.0.1:5500` atau `http://localhost:5500`

### Cara B: Shortcut Keyboard

1. Buka file `index.html`
2. Tekan `Alt+L` kemudian `Alt+O` (shortcut Live Server)
3. Atau klik tombol "Go Live" di bottom-right VS Code

### Output yang Diharapkan

**Browser terbuka di:**
```
http://127.0.0.1:5500
```

**VS Code output (Bottom panel):**
```
[Live Server] Server is started at 127.0.0.1:5500
[Live Server] Open http://127.0.0.1:5500
```

---

## 4. Verifikasi Frontend Terhubung ke Backend

### Langkah 1: Cek Konfigurasi API_BASE_URL

Buka file `assets/js/config.js` dan pastikan konfigurasi sudah benar:

```javascript
const getApiBaseUrl = () => {
  // Jika ada override dari window
  if (window.API_BASE_URL) {
    return window.API_BASE_URL;
  }
  // Development lokal — pastikan ini sesuai backend Anda
  if (location.hostname === 'localhost' || location.hostname === '127.0.0.1') {
    return 'http://localhost:8000/api';
  }
  // Production
  return 'https://api.sarpras.example.com/api';
};

export const API_BASE_URL = getApiBaseUrl();
```

**Catatan:**
- Pastikan URL backend `http://localhost:8000/api` sesuai dengan port backend Anda
- Jika backend di port lain, ubah `8000` ke port yang benar

### Langkah 2: Buka Browser Console

1. Di browser (tampilan frontend), tekan `F12` (DevTools)
2. Tab **Console**
3. Jalankan perintah untuk test koneksi:

```javascript
fetch('http://localhost:8000/api/health')
  .then(res => res.json())
  .then(data => console.log('Backend response:', data))
  .catch(err => console.error('Backend error:', err))
```

**Output yang diharapkan di Console:**
```
Backend response: {status: 'ok'}
```

Jika error `CORS error` atau `Failed to fetch`, berarti:
- Backend belum running → jalankan backend di terminal lain
- CORS_ALLOWED_ORIGIN tidak cocok → update `.env` backend

### Langkah 3: Update CORS Backend untuk Frontend Lokal

Di backend (folder `proyek-2-golang`), pastikan file `.env` sudah set:

```
CORS_ALLOWED_ORIGIN=http://127.0.0.1:5500
```

Jika sudah diubah, restart backend:
```powershell
# Terminal backend (berhenti dengan Ctrl+C)
# Lalu jalankan ulang:
go run ./cmd/server/main.go
```

---

## 5. Test Login & Navigasi Frontend

### Cek Apakah User Sudah Ada

Sebelum login, pastikan user sudah ada di database. Jika belum, insert user test:

**PowerShell:**
```powershell
psql -U postgres -d sarpras_db -c "
INSERT INTO organisasi (nama, jenis) VALUES ('Test Org', 'UKM');
INSERT INTO users (nama, email, password_hash, role, organisasi_id) 
VALUES ('Test Mahasiswa', 'mahasiswa@test.com', 'hashed', 'MAHASISWA', 1);
INSERT INTO users (nama, email, password_hash, role, organisasi_id) 
VALUES ('Test Sarpras', 'sarpras@test.com', 'hashed', 'SARPRAS', 1);
"
```

### Coba Login

Di halaman frontend:

1. **Email**: `mahasiswa@test.com`
2. **Password**: Apa pun (tergantung password yang di-hash di database)

**Catatan:** Untuk password yang benar, Anda perlu:
- Registrasi user baru melalui endpoint `/api/auth/register`, atau
- Gunakan user yang sudah ada dan coba password yang sesuai

Jika login berhasil:
- Token JWT tersimpan di `localStorage`
- Redirect ke dashboard sesuai role user
- Breadcrumb/navbar menampilkan nama user

---

## 6. Troubleshooting Frontend

| Masalah | Solusi |
|---------|--------|
| **Live Server tidak ada** | Install extension "Live Server" di VS Code |
| **Port 5500 sudah dipakai** | (a) Tutup app lain yang pakai 5500, atau (b) Ubah port di Live Server settings |
| **CORS error di Console** | Pastikan backend `.env` punya `CORS_ALLOWED_ORIGIN=http://127.0.0.1:5500` dan restart backend |
| **Backend unreachable** | (a) Pastikan backend running di `http://localhost:8000`, (b) Cek API_BASE_URL di config.js |
| **Login gagal 401** | (a) User belum ada di database, insert user test, (b) Password salah, (c) Backend error log check |
| **CSS/JS tidak ter-load** | Hard refresh: `Ctrl+Shift+R` (delete cache) |
| **Console error "api.legacy.js tidak ditemukan"** | Pastikan file `assets/js/api.legacy.js` ada di folder |

---

## 7. Workflow Frontend Development

### Edit & Hot Reload

Live Server akan **otomatis reload** browser saat Anda save file:

1. Edit HTML/CSS/JS di VS Code
2. Save file (`Ctrl+S`)
3. Browser otomatis refresh
4. Lihat perubahan langsung

**Tip:** Jika cache browser lama, gunakan `Ctrl+Shift+R` (hard refresh)

### Debugging dengan DevTools

Tekan `F12` untuk buka DevTools:

- **Console**: Lihat error JavaScript
- **Network**: Lihat request/response API
- **Elements**: Inspect HTML structure
- **Application**: Lihat localStorage (token JWT)

---

## 8. Struktur Folder Frontend

```
frontend-sarpras/
├── index.html                      # Halaman login
├── sarpras.html                    # Dashboard SARPRAS (fallback)
├── dashboard-mahasiswa.html        # Dashboard Mahasiswa
├── dashboard-sarpras.html          # Dashboard SARPRAS
├── dashboard-security.html         # Dashboard Security
├── jadwal-ruangan.html             # Jadwal ruangan
├── kelola-barang.html              # Kelola barang (SARPRAS)
├── kelola-ruangan.html             # Kelola ruangan (SARPRAS)
├── laporan-peminjaman.html         # Laporan peminjaman (SARPRAS)
├── pengajuan-peminjaman.html       # Form pengajuan peminjaman
├── register.html                   # Halaman registrasi
├── riwayat-peminjaman.html         # Riwayat peminjaman user
├── rooms.html                      # List ruangan (public)
├── verifikasi-peminjaman.html      # Verifikasi peminjaman (SARPRAS)
├── assets/
│   └── js/
│       ├── config.js               # ← API_BASE_URL & config
│       ├── api.js                  # Modern ES module untuk fetch
│       ├── api.legacy.js           # ← Compatibility globals (apiCall, checkAuth, dll)
│       ├── auth.js                 # Login form handler
│       └── dashboard-sarpras.js    # Dashboard logic
└── README.md
```

---

## 9. Integrasi Frontend ↔ Backend

### Frontend → Backend Flow

```
1. User masukkan email & password → form
2. JavaScript fetch ke POST /api/auth/login
3. Backend return token + user data
4. Frontend simpan token di localStorage
5. Setiap request ke endpoint protected, attach token di header Authorization: Bearer <token>
6. Backend verifikasi token via middleware
7. Response data ditampilkan di HTML
```

### Contoh API Call dari Frontend

**File: `assets/js/auth.js` (login)**

```javascript
import { API_BASE_URL } from './config.js';

async function login(email, password) {
  const response = await fetch(`${API_BASE_URL}/auth/login`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ email, password })
  });
  
  const data = await response.json();
  if (response.ok) {
    localStorage.setItem('token', data.token);
    localStorage.setItem('user', JSON.stringify(data.user));
    window.location.href = '/dashboard-mahasiswa.html'; // redirect
  } else {
    alert('Login gagal: ' + data.error);
  }
}
```

**File: `assets/js/api.legacy.js` (global function untuk legacy HTML)**

```javascript
// Global function untuk backward compatibility
window.apiCall = async function(endpoint, method = 'GET', body = null) {
  const token = localStorage.getItem('token');
  const headers = {
    'Content-Type': 'application/json',
    'Authorization': `Bearer ${token}`
  };
  
  const response = await fetch(`http://localhost:8000/api${endpoint}`, {
    method,
    headers,
    body: body ? JSON.stringify(body) : null
  });
  
  return await response.json();
};
```

---

## 10. Checklist Frontend Setup

- [ ] `frontend-sarpras` folder sudah di-clone dari GitHub
- [ ] File struktur lengkap (index.html, assets/js/*, dll)
- [ ] Live Server ter-install di VS Code
- [ ] Frontend berjalan di `http://127.0.0.1:5500`
- [ ] Backend berjalan di `http://localhost:8000`
- [ ] `config.js` menunjuk ke `http://localhost:8000/api`
- [ ] CORS_ALLOWED_ORIGIN di backend = `http://127.0.0.1:5500`
- [ ] Fetch ke `/api/health` success (testing di DevTools Console)
- [ ] User test sudah ada di database
- [ ] Login berhasil dan token tersimpan
- [ ] Dashboard sesuai role user ditampilkan

---

## 11. Next: Deploy ke GitHub Pages (Optional)

Jika sudah siap production:

1. Push `frontend-sarpras` ke GitHub repo Anda
2. Di repository Settings → Pages
3. Pilih branch `main` (atau `gh-pages`) dan deploy
4. Frontend akan online di `https://<username>.github.io/<repo-name>`
5. Update backend `.env` CORS_ALLOWED_ORIGIN ke GitHub Pages URL
6. Update `config.js` API_BASE_URL ke production API

---

## Catatan Penting

**Lokal Development:**
- Frontend: `http://127.0.0.1:5500` (Live Server)
- Backend: `http://localhost:8000`
- API Base: `http://localhost:8000/api`

**Production:**
- Frontend: `https://<username>.github.io/<repo>`
- Backend: `https://api.sarpras.example.com` (Render/Railway/VPS)
- API Base: `https://api.sarpras.example.com/api`
- CORS Origin: `https://<username>.github.io`

---

## Hubungi & Tanya

Jika ada yang tidak jelas, tanyakan langsung dengan deskripsi error. Bagikan:
- Screenshot error di browser/console
- Pesan error lengkap
- URL yang diakses

Salam, dan selamat ngoding! 🚀
