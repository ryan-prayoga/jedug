# 🚧 JEDUG (Jendela Edukasi & Data Umum Gangguan)

**Lapor, Pantau, Viralkan. Biar Nggak Ada yang Kejedug Lagi!**

JEDUG adalah platform pelaporan kerusakan jalan berbasis partisipasi publik yang dirancang untuk menciptakan transparansi dan memberikan tekanan sosial (_social pressure_) agar pihak berwenang bergerak lebih cepat. Melalui visualisasi data yang transparan, JEDUG mengubah keluhan warga menjadi statistik yang tidak bisa diabaikan.

---

## 🚀 Misi & Visi

Project ini berawal dari keresahan melihat kondisi jalanan di **Tangerang** yang seringkali luput dari perhatian pemerintah setempat. Meskipun saat ini difokuskan untuk wilayah Tangerang, JEDUG dirancang dengan sistem yang scalable untuk dapat diimplementasikan di seluruh wilayah Indonesia di masa mendatang.

## ✨ Fitur Utama (Core Features)

- **Peta Interaktif (Health Bar Maps):** Visualisasi kondisi jalan berdasarkan jalur (Hijau, Kuning, Merah, hingga Hitam untuk Zona Bahaya).
- **Sistem Lapor Anti-Hoax:** Mewajibkan pengambilan foto langsung di lokasi (GPS Radius Lock) untuk menjamin validitas data.
- **Viral Engine:**
  - **Ulang Tahun Lubang:** Icon kue ulang tahun 🎂 untuk jalan yang tidak diperbaiki lebih dari 30 hari.
  - **Rupiah Counter:** Estimasi kerugian materi warga (bensin & sparepart) yang terus bertambah selama jalan rusak dicuekin.
- **Emergency Button:** Pelaporan khusus kecelakaan dengan icon Tengkorak 💀 untuk menandai zona maut.
- **Poster Generator:** Membuat grafis berita otomatis untuk dibagikan ke WhatsApp Story atau Instagram.

## 🛠️ Tech Stack

Project ini dibangun dengan arsitektur **Monorepo** yang mengutamakan efisiensi biaya dan performa tinggi (High Concurrency):

- **Backend:** [Go](https://go.dev/) + [Fiber](https://gofiber.io/) (High Performance API).
- **Frontend:** [SvelteKit](https://kit.svelte.dev/) (Fast & Reactive UI).
- **Database:** [PostgreSQL](https://www.postgresql.org/) + [PostGIS](https://postgis.net/) (Spatial Queries & Zero-cost Geocoding).
- **Admin & Auth:** [PocketBase](https://pocketbase.io/) (Simplified Backend for Admin Management).
- **Storage:** [Cloudflare R2](https://www.cloudflare.com/developer-platform/r2/) (Zero Egress Image Storage).

## 📁 Struktur Repositori

```text
jedug/
├── backend/      # API Service (Go Fiber)
├── frontend/     # Web Application (SvelteKit)
├── docs/         # Masterplan & Database Schema
└── pocketbase/   # Admin Dashboard & Auth
```
