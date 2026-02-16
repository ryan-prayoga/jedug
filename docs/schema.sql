-- 1. EXTENSION WAJIB (Jalankan sekali di awal)
CREATE EXTENSION IF NOT EXISTS postgis;
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- 2. TABEL WILAYAH (Data Statis)
-- Import GeoJSON batas kecamatan/kelurahan Indonesia ke sini
CREATE TABLE districts (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100),
    level VARCHAR(50), -- 'kecamatan' or 'kelurahan'
    geom GEOMETRY(POLYGON, 4326) -- Polygon Batas Wilayah
);
-- Index Spatial supaya query "Saya ada di kecamatan mana?" cepat
CREATE INDEX idx_districts_geom ON districts USING GIST(geom);


-- 3. TABEL LAPORAN (Main Data)
CREATE TABLE reports (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    
    -- Identitas Pelapor
    fingerprint_hash VARCHAR(64) NOT NULL, -- ID User Tanpa Login
    user_id UUID NULL, -- ID User Kalau Login (Optional)
    
    -- Lokasi & Visual
    location GEOGRAPHY(POINT, 4326) NOT NULL,
    district_id INT REFERENCES districts(id), -- Auto-fill by PostGIS Trigger/Query
    image_url VARCHAR(255) NOT NULL,
    
    -- Status Jalan
    status VARCHAR(20) DEFAULT 'open', -- 'open', 'fixed', 'archived', 'rejected'
    severity INT DEFAULT 1, -- 1 (Kuning) - 5 (Hitam/Neraka)
    road_type VARCHAR(50), -- 'national', 'provincial', 'village' (Dari OSM)
    
    -- Statistik Viral
    view_count INT DEFAULT 0,
    reaction_count INT DEFAULT 0,
    estimated_loss DECIMAL(15, 2) DEFAULT 0, -- Rupiah Counter
    
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);
-- Index Spatial supaya query "Cari lubang radius 500m" cepat
CREATE INDEX idx_reports_loc ON reports USING GIST(location);


-- 4. TABEL BUKTI TAMBAHAN (Galeri Warga)
CREATE TABLE report_proofs (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    report_id UUID REFERENCES reports(id) ON DELETE CASCADE,
    fingerprint_hash VARCHAR(64) NOT NULL,
    image_url VARCHAR(255) NOT NULL,
    taken_at TIMESTAMPTZ DEFAULT NOW()
);


-- 5. TABEL INTERAKSI (Reaction & Comment)
CREATE TABLE interactions (
    id SERIAL PRIMARY KEY,
    report_id UUID REFERENCES reports(id) ON DELETE CASCADE,
    fingerprint_hash VARCHAR(64) NOT NULL,
    
    -- Tipe Interaksi
    type VARCHAR(20) NOT NULL, -- 'angry', 'danger', 'comment', 'upvote'
    value TEXT, -- Isi komentar jika type='comment'
    
    created_at TIMESTAMPTZ DEFAULT NOW(),
    
    -- Constraint: 1 Device cuma boleh 1 reaction per laporan (Anti Spam)
    UNIQUE(report_id, fingerprint_hash, type)
);


-- 6. TABEL USER (Untuk Gamification - Optional)
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    email VARCHAR(255) UNIQUE NOT NULL,
    name VARCHAR(100),
    avatar_url VARCHAR(255),
    xp_points INT DEFAULT 0,
    rank_title VARCHAR(50) DEFAULT 'Warga Biasa',
    is_verified BOOLEAN DEFAULT FALSE, -- Centang Biru
    created_at TIMESTAMPTZ DEFAULT NOW()
);