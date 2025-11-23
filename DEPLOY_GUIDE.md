# Panduan Deploy Backend Golang ke Production

Panduan lengkap untuk deploy backend API Golang (`backend-sarpras`) ke production dengan berbagai pilihan hosting provider.

---

## 📋 Daftar Isi

1. [Persiapan Pre-Deploy](#persiapan-pre-deploy)
2. [Deploy ke Render](#deploy-ke-render) - **Recommended (Gratis + Mudah)**
3. [Deploy ke Railway](#deploy-ke-railway)
4. [Deploy ke Heroku](#deploy-ke-heroku) - **Bayar**
5. [Deploy ke VPS (Cloud)**](#deploy-ke-vps-cloud)
6. [Setup Domain & SSL](#setup-domain--ssl)
7. [Monitoring & Logging](#monitoring--logging)

---

## Persiapan Pre-Deploy

### 1. Update Environment Variables untuk Production

Sebelum deploy, pastikan `main.go` dan middleware menggunakan env variables dengan baik:

**File: `internal/config/config.go`** ✅ (sudah benar)
- `DATABASE_URL` → dari environment
- `PORT` → dari environment (default 8000)
- `JWT_SECRET` → dari environment
- `CORS_ALLOWED_ORIGIN` → dari environment

### 2. Update JWT_SECRET untuk Production

**⚠️ PENTING:** Ganti hardcoded JWT_SECRET dengan environment variable.

**File: `middleware/auth.go`**
```go
// ❌ SEBELUMNYA (TIDAK AMAN):
var jwtSecret = []byte("your-secret-key-change-in-production")

// ✅ SESUDAHNYA (AMAN):
var jwtSecret = []byte(os.Getenv("JWT_SECRET"))
```

**File: `services/auth_service.go`**
```go
// ✅ Sama seperti di atas
var jwtSecret = []byte(os.Getenv("JWT_SECRET"))
```

### 3. Pastikan Database Connection Aman

✅ **Sudah benar di `config.go`:**
```go
dbURL := os.Getenv("DATABASE_URL")
if dbURL == "" {
    log.Fatal("DATABASE_URL is required")
}
```

### 4. Build Test Lokal

Sebelum push ke production, test build:

```powershell
# Bersihkan sebelumnya
go clean

# Build aplikasi
go build -o backend-sarpras.exe ./cmd/server

# Test jalankan (perlu set env vars)
$env:DATABASE_URL = "postgresql://..."
$env:JWT_SECRET = "your-secret-key"
$env:PORT = "8000"
$env:CORS_ALLOWED_ORIGIN = "https://frontend.example.com"
.\backend-sarpras.exe
```

### 5. Commit & Push ke GitHub

```powershell
cd 'd:\Kuliah\Proyek 2\proyek-2-golang'
git add -A
git commit -m "chore: prepare for production deployment"
git push origin refactor/api-only
```

---

## Deploy ke Render

### ✨ Kenapa Render?

- **Gratis** dengan tier free (jam pertama gratis, lalu $0.50/hour)
- **Setup mudah** via GitHub integration
- **Auto-deploy** saat push ke branch
- **Built-in PostgreSQL** (optional)
- **Environment variables** management via dashboard
- **SSL/HTTPS** gratis

### 📝 Langkah-Langkah Deploy ke Render

#### **Step 1: Login/Register di Render**

1. Buka [render.com](https://render.com)
2. Klik **"Sign Up"** → pilih **GitHub**
3. Authorize akses ke GitHub repository

#### **Step 2: Buat Web Service Baru**

1. Di dashboard Render, klik **"New +"** → **"Web Service"**
2. Pilih repository `proyek-2-golang` (atau `backend-sarpras` jika sudah di-rename)
3. Pilih branch: `refactor/api-only` (atau branch production Anda)

#### **Step 3: Konfigurasi Build & Deploy**

**Form Configuration:**

| Field | Value |
|-------|-------|
| **Name** | `backend-sarpras` (atau nama unik lainnya) |
| **Environment** | `Go` |
| **Build Command** | `go build -o backend-sarpras ./cmd/server` |
| **Start Command** | `./backend-sarpras` |
| **Instance Type** | `Free` (atau Starter untuk production) |

**Advanced Settings:**
- Auto-deploy: ✅ (auto-deploy on push)
- Persistent Disk: Optional (jika perlu store files)

#### **Step 4: Atur Environment Variables**

Scroll ke bagian **"Environment"** dan tambahkan:

```
DATABASE_URL=postgresql://user:password@db.xxx.supabase.co:5432/postgres?sslmode=require
JWT_SECRET=your-very-secure-random-string-min-32-chars-long-change-this
PORT=10000
CORS_ALLOWED_ORIGIN=https://username.github.io/frontend-sarpras
```

**Catatan:**
- Render akan assign PORT otomatis, tapi kita set untuk consistency
- `CORS_ALLOWED_ORIGIN` → ganti dengan URL frontend GitHub Pages Anda

#### **Step 5: Deploy

1. Klik **"Create Web Service"**
2. Render akan otomatis:
   - Clone repository
   - Install dependencies (`go mod download`)
   - Build aplikasi
   - Deploy & jalankan

3. Monitor progress di tab **"Logs"**

#### **Step 6: Dapatkan URL Production**

Setelah deploy sukses, Anda akan dapat URL seperti:
```
https://backend-sarpras.onrender.com
```

**Gunakan URL ini untuk:**
- Frontend API base: `https://backend-sarpras.onrender.com/api`
- Update `CORS_ALLOWED_ORIGIN` jika frontend di domain lain

### 🔄 Auto-Deploy di Render

Setiap kali Anda push ke branch `refactor/api-only`, Render akan:
1. Trigger build otomatis
2. Jalankan build command
3. Deploy versi baru
4. Restart service

---

## Deploy ke Railway

### ✨ Kenapa Railway?

- **UI lebih modern** dari Render
- **Free tier** $5/bulan gratis credit
- **Simple & cepat** untuk setup
- **PostgreSQL built-in** bisa langsung gunakan

### 📝 Langkah-Langkah Deploy ke Railway

#### **Step 1: Login di Railway**

1. Buka [railway.app](https://railway.app)
2. Klik **"Login"** → **"Login with GitHub"**
3. Authorize Railway access

#### **Step 2: Import Project dari GitHub**

1. Dashboard → **"New Project"** → **"Deploy from GitHub repo"**
2. Pilih repository `proyek-2-golang`
3. Pilih branch `refactor/api-only`

#### **Step 3: Konfigurasi Environment**

1. Di Railway dashboard, klik **"Variables"**
2. Tambahkan environment variables:

```
DATABASE_URL=postgresql://...@db.xxx.supabase.co:5432/postgres
JWT_SECRET=your-secure-key-here
CORS_ALLOWED_ORIGIN=https://frontend-url.github.io
PORT=5000
```

#### **Step 4: Setup Build Configuration

1. Klik **"Settings"**
2. Build Command: `go build -o backend ./cmd/server`
3. Start Command: `./backend`

#### **Step 5: Deploy

1. Railway akan otomatis detect Go project
2. Build & deploy otomatis
3. Monitor di tab **"Deployment"**

#### **Step 6: Dapatkan URL**

URL deployment akan muncul di dashboard, contoh:
```
https://backend-sarpras-production.up.railway.app
```

---

## Deploy ke Heroku (Bayar)

### ⚠️ Note: Heroku Free Tier Sudah Dihapus (Shutdown Nov 2022)

Alternatif bayar:
- **Heroku Eco Dyno**: $5/bulan
- **Heroku Standard Dyno**: $7-50/bulan

Jika ingin tetap pakai Heroku:

### 📝 Langkah-Langkah (Jika Tetap Pakai Heroku)

#### **Step 1: Install Heroku CLI**

```powershell
# Download & install dari heroku.com/download
# Atau via chocolatey:
choco install heroku-cli

# Verifikasi
heroku --version
```

#### **Step 2: Login Heroku**

```powershell
heroku login
```

#### **Step 3: Buat Heroku App**

```powershell
cd 'd:\Kuliah\Proyek 2\proyek-2-golang'
heroku create backend-sarpras
```

#### **Step 4: Add PostgreSQL Addon**

```powershell
heroku addons:create heroku-postgresql:standard-0
```

#### **Step 5: Set Environment Variables**

```powershell
heroku config:set JWT_SECRET="your-secret-key"
heroku config:set CORS_ALLOWED_ORIGIN="https://frontend-url.github.io"
```

#### **Step 6: Deploy**

```powershell
git push heroku refactor/api-only:main
```

#### **Step 7: Monitor**

```powershell
heroku logs --tail
heroku open
```

---

## Deploy ke VPS (Cloud)

Pilihan untuk kontrol penuh & biaya lebih murah.

### Provider VPS Populer:
- **DigitalOcean** - $4-6/bulan
- **Linode** - $5/bulan
- **AWS EC2** - $10-50/bulan
- **Google Cloud** - Gratis tier + bayar sesuai pakai
- **Azure** - Gratis tier + bayar sesuai pakai
- **Hetzner** - €3/bulan (murah!)

### 📝 Setup Dasar VPS (Ubuntu 22.04)

#### **Step 1: SSH ke VPS**

```powershell
ssh root@your_vps_ip
```

#### **Step 2: Update System**

```bash
apt update && apt upgrade -y
```

#### **Step 3: Install Go**

```bash
wget https://go.dev/dl/go1.25.3.linux-amd64.tar.gz
tar -C /usr/local -xzf go1.25.3.linux-amd64.tar.gz
echo 'export PATH=$PATH:/usr/local/go/bin' >> /etc/profile
source /etc/profile
go version
```

#### **Step 4: Install Git & PostgreSQL Client**

```bash
apt install -y git postgresql-client
```

#### **Step 5: Clone Repository**

```bash
cd /home
git clone https://github.com/ditlabs/proyek-2-golang.git
cd proyek-2-golang
git checkout refactor/api-only
```

#### **Step 6: Setup Environment Variables**

```bash
cat > .env << EOF
DATABASE_URL=postgresql://user:pass@db.xxx.supabase.co:5432/postgres?sslmode=require
JWT_SECRET=your-very-secure-key-here
PORT=8000
CORS_ALLOWED_ORIGIN=https://frontend-url.github.io
EOF
```

#### **Step 7: Build & Test**

```bash
go mod download
go build -o backend ./cmd/server
./backend
```

#### **Step 8: Setup Systemd Service**

Buat file service agar backend auto-start saat VPS reboot:

```bash
cat > /etc/systemd/system/backend-sarpras.service << EOF
[Unit]
Description=Backend Sarpras API
After=network.target

[Service]
User=root
WorkingDirectory=/home/proyek-2-golang
EnvironmentFile=/home/proyek-2-golang/.env
ExecStart=/home/proyek-2-golang/backend
Restart=always
RestartSec=10

[Install]
WantedBy=multi-user.target
EOF
```

#### **Step 9: Enable & Start Service**

```bash
systemctl daemon-reload
systemctl enable backend-sarpras
systemctl start backend-sarpras
systemctl status backend-sarpras
```

#### **Step 10: Setup Nginx Reverse Proxy**

```bash
apt install -y nginx

# Buat config
cat > /etc/nginx/sites-available/backend-sarpras << EOF
server {
    listen 80;
    server_name api.sarpras.example.com;

    location / {
        proxy_pass http://localhost:8000;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
}
EOF

# Enable site
ln -s /etc/nginx/sites-available/backend-sarpras /etc/nginx/sites-enabled/
nginx -t
systemctl restart nginx
```

#### **Step 11: Setup SSL (Let's Encrypt)**

```bash
apt install -y certbot python3-certbot-nginx
certbot --nginx -d api.sarpras.example.com
```

---

## Setup Domain & SSL

### 1. Ganti Domain ke Custom Domain

#### **Untuk Render:**
1. Settings → Custom Domain
2. Input domain: `api.sarpras.example.com`
3. Render akan generate CNAME record
4. Update DNS provider

#### **Untuk Railway:**
1. Settings → Custom Domain
2. Add domain & follow instruction

#### **Untuk VPS:**
1. Update DNS A record ke IP VPS
2. Setup SSL via Certbot (seperti di atas)

### 2. Update Frontend Configuration

Setelah domain production siap, update `frontend-sarpras/assets/js/config.js`:

```javascript
const getApiBaseUrl = () => {
  if (window.API_BASE_URL) {
    return window.API_BASE_URL;
  }
  
  // Production
  if (location.hostname === 'username.github.io') {
    return 'https://api.sarpras.example.com/api';
  }
  
  // Development
  return 'http://localhost:8000/api';
};

export const API_BASE_URL = getApiBaseUrl();
```

---

## Monitoring & Logging

### 1. Render Monitoring

```
Dashboard → "Logs" tab → real-time logs
Dashboard → "Metrics" → CPU, Memory usage
```

### 2. Railway Monitoring

```
Dashboard → "Deployments" → logs
"Observability" → metrics
```

### 3. VPS Monitoring

```bash
# Check status
systemctl status backend-sarpras

# View logs
journalctl -u backend-sarpras -f

# Check resource usage
top
htop
```

### 4. Add External Monitoring (Optional)

- **UptimeRobot** - Monitor endpoint & alert jika down
- **DataDog** - Application performance monitoring
- **Sentry** - Error tracking
- **LogRocket** - Session replay & debugging

---

## ✅ Checklist Pre-Deployment

- [ ] JWT_SECRET ter-set dari environment variable (bukan hardcoded)
- [ ] DATABASE_URL ter-set dari environment variable
- [ ] CORS_ALLOWED_ORIGIN sesuai dengan frontend URL
- [ ] Build test lokal sukses (`go build`)
- [ ] Jalankan lokal & test endpoints
- [ ] All tests pass (jika ada)
- [ ] Commit & push ke GitHub
- [ ] Database migration sudah di-run
- [ ] Frontend URL ter-whitelist di CORS
- [ ] Monitoring & alerting sudah setup

---

## 🆘 Troubleshooting Deploy

| Masalah | Solusi |
|---------|--------|
| Build gagal | Cek logs, pastikan Go version kompatibel, run `go mod tidy` |
| Connection refused | Pastikan port accessible, firewall settings, DATABASE_URL benar |
| CORS error dari frontend | Update `CORS_ALLOWED_ORIGIN` di env vars, restart service |
| 502 Bad Gateway | Backend crash/down, cek logs di provider |
| Database connection timeout | Cek DATABASE_URL, whitelist IP di database |
| JWT error | Pastikan JWT_SECRET sama dengan lokal, token sudah expire |

---

## 📚 Quick Start Deploy Summary

### **Render (Recommended)**
```powershell
# 1. Commit code
git push origin refactor/api-only

# 2. Go to render.com → New Web Service
# 3. Select GitHub repo & branch
# 4. Set environment variables
# 5. Deploy — Done! ✅
```

### **VPS (Full Control)**
```bash
# 1. SSH ke VPS
# 2. Install Go, Git, clone repo
# 3. Setup .env file
# 4. Build: go build -o backend ./cmd/server
# 5. Setup systemd service
# 6. Setup Nginx + SSL
# 7. Deploy — Done! ✅
```

---

## 🎉 Setelah Deploy Sukses

1. **Test endpoints production:**
```powershell
curl https://api.sarpras.example.com/api/health
```

2. **Update frontend config** ke production API

3. **Deploy frontend ke GitHub Pages** (jika belum)

4. **Test end-to-end:**
   - Login via frontend production
   - Fetch data dari backend production
   - Verify CORS working

5. **Monitor & maintain:**
   - Setup uptime monitoring
   - Regular backup database
   - Update dependencies
   - Monitor logs

---

**Pertanyaan? Tanya langsung!** 🚀
