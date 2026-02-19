# 🚧 JEDUG (Jendela Edukasi & Data Umum Gangguan)

**Lapor, Pantau, Viralkan. Biar Nggak Ada yang Kejedug Lagi!**

JEDUG adalah platform pelaporan kerusakan jalan berbasis partisipasi publik yang dirancang untuk menciptakan transparansi dan memberikan tekanan sosial (_social pressure_) agar pihak berwenang bergerak lebih cepat. Melalui visualisasi data yang transparan, JEDUG mengubah keluhan warga menjadi statistik yang tidak bisa diabaikan.

---

## 🚀 Misi & Visi

Project ini berawal dari keresahan melihat kondisi jalanan di **Tangerang** yang seringkali luput dari perhatian pemerintah setempat. Meskipun saat ini difokuskan untuk wilayah Tangerang, JEDUG dirancang dengan sistem yang scalable untuk dapat diimplementasikan di seluruh wilayah Indonesia di masa mendatang.

## ✨ Fitur Utama (Core Features)

- **Peta Interaktif (Health Bar Maps):** Visualisasi kondisi jalan berdasarkan jalur (Hijau, Kuning, Merah, hingga Hitam untuk Zona Bahaya).
- **Sistem Lapor Anti-Hoax:** Mewajibkan pengambilan foto langsung di lokasi (GPS Radius Lock) untuk menjamin validitas data.
- **Client-Side Image Compression:** Foto otomatis dikompresi ke WebP (max 1080px, 80% quality) sebelum upload.
- **Device Fingerprinting:** No-login-first user identification via browser fingerprint.
- **Smart Grouping:** Laporan baru di radius 10m dari laporan lama = masuk galeri foto (bukan pin baru).
- **Zero-Cost Geocoding:** Reverse geocoding gratis menggunakan PostGIS `ST_Contains` query.
- **Viral Engine:**
  - **Ulang Tahun Lubang:** Icon kue ulang tahun 🎂 untuk jalan yang tidak diperbaiki lebih dari 30 hari.
  - **Rupiah Counter:** Estimasi kerugian materi formula: `(views x Rp500) + (hari_rusak x Rp100.000)`.
- **Emergency Button:** Pelaporan khusus kecelakaan dengan icon Tengkorak 💀 untuk menandai zona maut.
- **Rate Limiting:** Proteksi dari spam (100 req/min API, 10 req/min upload, 5 req/min report).

## 🛠️ Tech Stack

Project ini dibangun dengan arsitektur **Monorepo** yang mengutamakan efisiensi biaya dan performa tinggi (High Concurrency):

- **Backend:** [Go](https://go.dev/) + [Fiber](https://gofiber.io/) (High Performance API).
- **Frontend:** [SvelteKit](https://kit.svelte.dev/) (Fast & Reactive UI).
- **Database:** [PostgreSQL](https://www.postgresql.org/) + [PostGIS](https://postgis.net/) (Spatial Queries & Zero-cost Geocoding).
- **Admin & Auth:** [PocketBase](https://pocketbase.io/) (Simplified Backend for Admin Management).
- **Storage:** [Cloudflare R2](https://www.cloudflare.com/developer-platform/r2/) (Zero Egress Image Storage).
- **Maps:** [Leaflet](https://leafletjs.com/) + [CartoDB](https://carto.com/) Dark Matter tiles.

## 📁 Struktur Repositori

```text
jedug/
├── backend/              # API Service (Go Fiber)
│   ├── cmd/server/       # Entry point
│   ├── internal/
│   │   ├── config/       # Configuration
│   │   ├── database/     # PostgreSQL connection
│   │   ├── handlers/     # HTTP handlers (report, district, upload)
│   │   ├── middleware/    # Rate limiting
│   │   ├── models/       # Data models
│   │   └── repository/   # Database queries (report, district, interaction)
│   └── uploads/          # Uploaded images (gitignored)
├── frontend/             # Web Application (SvelteKit)
│   └── src/
│       ├── lib/
│       │   ├── api/      # API client
│       │   ├── components/  # Svelte components
│       │   ├── data/     # Mock data (fallback)
│       │   ├── stores/   # Svelte stores
│       │   ├── styles/   # CSS design system
│       │   └── utils/    # Image compression, fingerprint
│       └── routes/       # Pages (peta, lapor, ranking, profil, detail)
├── docs/                 # Masterplan & Database Schema
└── pocketbase/           # Admin Dashboard & Auth
```

## 🚦 API Endpoints

| Method | Path                            | Description                          |
| ------ | ------------------------------- | ------------------------------------ |
| `GET`  | `/api/v1/health`                | Health check                         |
| `GET`  | `/api/v1/districts`             | Get district boundaries (GeoJSON)    |
| `GET`  | `/api/v1/reports`               | List reports (paginated, filterable) |
| `GET`  | `/api/v1/reports/nearby`        | Get nearby reports (spatial query)   |
| `GET`  | `/api/v1/reports/:id`           | Get report detail                    |
| `POST` | `/api/v1/reports`               | Create new report                    |
| `POST` | `/api/v1/reports/:id/reactions` | Add reaction                         |
| `GET`  | `/api/v1/reports/:id/comments`  | Get comments                         |
| `POST` | `/api/v1/reports/:id/comments`  | Add comment                          |
| `GET`  | `/api/v1/ranking`               | Get kecamatan ranking                |
| `POST` | `/api/v1/upload`                | Upload image                         |

## 🏃 Quick Start

```bash
# 1. Start database
docker compose up -d db

# 2. Start backend
cd backend
cp .env.example .env  # Edit with your DB credentials
go run ./cmd/server/

# 3. Start frontend
cd frontend
npm install
npm run dev
```
