# PROJECT BLUEPRINT: JEDUG

**Tagline:** "Lapor, Pantau, Viralkan. Biar Nggak Ada yang Kejedug Lagi."
**Core Philosophy:** _High Impact, Low Cost, High Viral Potential, Mobile-First._
**Status:** Ready for Development.

---

## 1. Executive Summary

JEDUG adalah platform pelaporan kerusakan jalan berbasis partisipasi publik ("Waze khusus Jalan Rusak"). Fokus utamanya adalah memvisualisasikan data kerusakan menjadi tekanan sosial (_social pressure_) agar pemerintah bertindak.

**Key Design Principles:**

1.  **Mobile-First (Prioritas Utama):** UI/UX dirancang seperti _Native App_. Tombol besar, navigasi di bawah (thumb-zone friendly), dan ringan untuk sinyal tidak stabil.
2.  **Responsive Desktop:** Tampilan desktop bukan sekadar versi mobile yang ditarik, melainkan mode "Command Center/Dashboard" untuk pemantauan wilayah yang luas.
3.  **Cost Efficiency:** Arsitektur seminimal mungkin menggunakan layanan berbayar (Zero API Cost Strategy).

---

## 2. Tech Stack (The "Zero-Cost" Architecture)

- **Backend:** **Go (Golang) + Fiber**.
  - _Reason:_ High performance, rendah penggunaan RAM/CPU. Bisa jalan di VPS $5/bulan untuk ribuan user.
- **Frontend:** **SvelteKit** (SSR + SPA).
  - _Reason:_ Bundle size sangat kecil, loading cepat di jaringan buruk. Wajib SSR untuk Link Preview di WA/IG.
- **Database:** **PostgreSQL + PostGIS**.
  - _Critical:_ Digunakan untuk _Spatial Queries_ DAN _Reverse Geocoding_ gratis (Mapping koordinat ke Kecamatan tanpa API Google).
- **Auth & Admin:** **PocketBase**.
  - _Reason:_ Single binary (Go), sudah include Auth, Admin UI, dan Database management. Sangat hemat resource.
- **Storage:** **Cloudflare R2**.
  - _Reason:_ S3 Compatible, **No Egress Fees** (Gratis bandwidth download).
- **Maps (Frontend):** **MapLibre GL JS**.
  - _Reason:_ Open-source, gratis, ringan, support vector tiles.
- **Maps (Data Source):** **OpenStreetMap (OSM)** via PMTiles atau Free Tile Providers.
  - _Strict Rule:_ **DILARANG** menggunakan Google Maps API / Mapbox Paid API.

---

## 3. UI/UX Strategy (Responsive & Adaptive)

### A. Mobile View (The Reporter Interface)

Target: User di jalan, pakai satu tangan.

- **Bottom Navigation Bar:** Menu utama (Peta, Lapor, Ranking, Profil) di bawah layar.
- **Floating Action Button (FAB):** Tombol "Lapor / Kamera" besar di tengah bawah.
- **Card-Based UI:** Info jalan rusak muncul sebagai kartu yang bisa di-swipe (mirip Tinder/Instagram), bukan tabel.
- **Touch Targets:** Semua tombol minimal 48x48px (mudah ditekan jempol).

### B. Desktop View (The Monitoring Dashboard)

Target: Netizen di rumah, Admin, Pemerintah.

- **Split Screen:** Kiri Peta Besar (Full Height), Kanan Sidebar List Laporan & Statistik.
- **Hover Interactions:** Arahkan mouse ke peta langsung muncul preview foto & status (Tooltip canggih).
- **Wide Leaderboard:** Tampilan klasemen kecamatan & kerugian rupiah lebih detail dan lebar.

---

## 4. Cost Efficiency Protocols (Wajib Diterapkan)

### A. Zero-Cost Geocoding (Logic PostGIS)

Jangan gunakan API berbayar untuk mendapatkan nama lokasi (Kecamatan/Kelurahan).

- **Metode:** Import data batas wilayah (Shapefile/GeoJSON) Indonesia ke tabel PostGIS.
- **Logic:** Query `ST_Contains` ke tabel wilayah untuk mendapatkan nama Kecamatan/Kelurahan dari Lat/Long.
- **Biaya:** **Rp 0,-** (CPU only).

### B. Client-Side Image Compression

Jangan terima file mentah dari user.

- **Logic Svelte:** Sebelum upload ke server, browser **WAJIB**:
  1.  Resize gambar ke max-width **1080px**.
  2.  Compress ke format **WebP** (Quality 80%).
  3.  Strip metadata (EXIF) kecuali tanggal/lokasi.
- **Hasil:** File 10MB menjadi ~150KB. Hemat storage R2 & bandwidth.

---

## 5. Core Features & Business Logic

### A. Input Data (The Scout)

1.  **Radius Lock:** GPS User vs Titik Lapor < 20 meter.
2.  **Camera Enforcement:** Wajib kamera langsung (`capture="environment"`).
3.  **Smart Grouping:** Laporan baru di radius 10m dari laporan lama = Masuk Galeri Foto (Bukan Pin Baru).

### B. Map Visualization (Social Pressure)

- **Jalur Berwarna:** Hijau (Aman), Kuning (Rusak Ringan), Merah (Rusak Berat), Hitam (Ada Korban).
- **Auto-Detect Road Type:** Gunakan metadata OSM (Way ID) untuk membedakan Jalan Nasional vs Kampung.

### C. The "Viral Engine" (Sarkasme Data)

1.  **Ulang Tahun Lubang:** Icon Kue 🎂 (>30 Hari), Icon Fosil 🦖 (>3 Bulan).
2.  **Rupiah Counter (Estimasi Kerugian):** Formula `(Viewers x Rp 500) + (Hari Rusak x Rp 100.000)`.
3.  **Emergency Button:** Tombol "Ada Korban" -> Icon Tengkorak 💀 -> Trigger Poster Generator.

### D. User System (Gamification)

1.  **No-Login First:** User bisa lapor/reaction langsung. Identifikasi via **Device Fingerprint**.
2.  **Optional Login:** Login Google untuk simpan history & naik level (Warga -> Pak RT -> Menteri PU Swasta).
3.  **Viral Sharing:** Generate gambar poster otomatis ("Berita Duka Jalan Rusak") untuk status WA.

---

## 6. Database Structure

_Lihat file `schema.sql` terlampir untuk struktur tabel PostgreSQL dan PostGIS._

---

## 7. Admin Panel & Safety

- **Moderasi:** Queue laporan yang di-flag "Spam/Hoax" (5 Flag = Auto Hide).
- **Ban Hammer:** Blokir Device Fingerprint user iseng.
- **Legal:** Checkbox S&K wajib dicentang user sebelum lapor ("Saya bertanggung jawab atas data ini").
- **Rate Limit:** Batasi API request per IP untuk mencegah serangan bot.

---

_Dokumen ini adalah acuan final pengembangan JEDUG. Prioritas utama adalah Mobile Experience, Stabilitas Server, dan Efisiensi Biaya._
