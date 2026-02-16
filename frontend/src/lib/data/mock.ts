import type { Report, KecamatanRank, UserProfile } from "$lib/types";

export const dummyReports: Report[] = [
  {
    id: "1",
    title: "Lubang Besar di Jl. Sudirman",
    description:
      "Lubang berukuran sekitar 1x1 meter dengan kedalaman 30cm. Sudah makan banyak korban ban bocor. Sangat berbahaya saat malam hari karena tidak ada penerangan.",
    lat: -6.2088,
    lng: 106.8456,
    location: "Jl. Jend. Sudirman No. 45",
    kecamatan: "Tanah Abang",
    kelurahan: "Bendungan Hilir",
    severity: "berat",
    status: "diverifikasi",
    photos: [
      "https://images.unsplash.com/photo-1515162816999-a0c47dc192f7?w=400&h=300&fit=crop",
      "https://images.unsplash.com/photo-1590496793929-36417d3117de?w=400&h=300&fit=crop",
    ],
    reporterName: "Budi Santoso",
    reporterLevel: "Pak RT",
    reactions: 234,
    comments: 45,
    views: 1520,
    createdAt: "2026-01-10",
    daysOld: 37,
    estimatedLoss: 4460000,
  },
  {
    id: "2",
    title: "Jalan Amblas Dekat Pasar Minggu",
    description:
      "Jalan amblas sepanjang 3 meter, lebar 2 meter. Kendaraan roda empat tidak bisa lewat. Warga sudah pasang tanda peringatan dari bambu.",
    lat: -6.2835,
    lng: 106.843,
    location: "Jl. Raya Pasar Minggu KM 5",
    kecamatan: "Pasar Minggu",
    kelurahan: "Pejaten Timur",
    severity: "berat",
    status: "baru",
    photos: [
      "https://images.unsplash.com/photo-1590496793929-36417d3117de?w=400&h=300&fit=crop",
    ],
    reporterName: "Siti Rahmawati",
    reporterLevel: "Warga Peduli",
    reactions: 567,
    comments: 89,
    views: 3200,
    createdAt: "2025-11-05",
    daysOld: 103,
    estimatedLoss: 11900000,
  },
  {
    id: "3",
    title: "Retakan Aspal di Jl. Gatot Subroto",
    description:
      "Retakan memanjang sekitar 5 meter di bahu jalan. Belum terlalu dalam tapi berbahaya untuk motor.",
    lat: -6.235,
    lng: 106.827,
    location: "Jl. Gatot Subroto Kav. 22",
    kecamatan: "Setiabudi",
    kelurahan: "Karet Kuningan",
    severity: "ringan",
    status: "proses",
    photos: [
      "https://images.unsplash.com/photo-1515162816999-a0c47dc192f7?w=400&h=300&fit=crop",
    ],
    reporterName: "Ahmad Fauzi",
    reporterLevel: "Warga Biasa",
    reactions: 45,
    comments: 12,
    views: 380,
    createdAt: "2026-02-01",
    daysOld: 15,
    estimatedLoss: 1690000,
  },
  {
    id: "4",
    title: "Kecelakaan Akibat Jalan Rusak",
    description:
      "Pengendara motor jatuh karena lubang yang tidak terlihat saat hujan. Korban mengalami luka di tangan dan kaki. Lubang ini sudah dilaporkan 2 bulan lalu tapi belum diperbaiki.",
    lat: -6.1751,
    lng: 106.865,
    location: "Jl. Pemuda No. 10",
    kecamatan: "Pulo Gadung",
    kelurahan: "Jati",
    severity: "korban",
    status: "diverifikasi",
    photos: [
      "https://images.unsplash.com/photo-1590496793929-36417d3117de?w=400&h=300&fit=crop",
      "https://images.unsplash.com/photo-1515162816999-a0c47dc192f7?w=400&h=300&fit=crop",
    ],
    reporterName: "Devi Lestari",
    reporterLevel: "Camat Tangguh",
    reactions: 1205,
    comments: 234,
    views: 8900,
    createdAt: "2026-02-14",
    daysOld: 2,
    estimatedLoss: 4650000,
  },
  {
    id: "5",
    title: "Aspal Mengelupas di Jl. Kemang",
    description:
      "Aspal mengelupas di beberapa titik sepanjang 100 meter. Debu dan kerikil berterbangan saat kendaraan lewat.",
    lat: -6.2607,
    lng: 106.8137,
    location: "Jl. Kemang Raya No. 88",
    kecamatan: "Mampang Prapatan",
    kelurahan: "Bangka",
    severity: "sedang",
    status: "diverifikasi",
    photos: [
      "https://images.unsplash.com/photo-1515162816999-a0c47dc192f7?w=400&h=300&fit=crop",
    ],
    reporterName: "Rudi Hartono",
    reporterLevel: "Pak Lurah",
    reactions: 189,
    comments: 33,
    views: 1100,
    createdAt: "2026-01-25",
    daysOld: 22,
    estimatedLoss: 2750000,
  },
  {
    id: "6",
    title: "Jalan Berlubang Berjajar",
    description:
      "Lima lubang berjajar dalam jarak 50 meter. Pengendara harus slalom untuk menghindar. Sudah pernah ditambal tapi jebol lagi.",
    lat: -6.1944,
    lng: 106.8229,
    location: "Jl. Salemba Raya No. 30",
    kecamatan: "Senen",
    kelurahan: "Paseban",
    severity: "sedang",
    status: "baru",
    photos: [
      "https://images.unsplash.com/photo-1590496793929-36417d3117de?w=400&h=300&fit=crop",
    ],
    reporterName: "Nina Kusuma",
    reporterLevel: "Warga Peduli",
    reactions: 312,
    comments: 56,
    views: 2100,
    createdAt: "2026-01-18",
    daysOld: 29,
    estimatedLoss: 3950000,
  },
  {
    id: "7",
    title: "Trotoar Rusak & Jalan Retak",
    description:
      "Trotoar pecah dan jalan retak di depan sekolah. Bahaya untuk anak-anak sekolah yang lewat setiap hari.",
    lat: -6.22,
    lng: 106.85,
    location: "Jl. Casablanca No. 15",
    kecamatan: "Tebet",
    kelurahan: "Menteng Dalam",
    severity: "sedang",
    status: "selesai",
    photos: [
      "https://images.unsplash.com/photo-1515162816999-a0c47dc192f7?w=400&h=300&fit=crop",
    ],
    reporterName: "Hendra Wijaya",
    reporterLevel: "Pak RT",
    reactions: 156,
    comments: 28,
    views: 950,
    createdAt: "2025-12-20",
    daysOld: 58,
    estimatedLoss: 0,
  },
];

export const dummyKecamatanRanks: KecamatanRank[] = [
  {
    rank: 1,
    name: "Pasar Minggu",
    totalReports: 45,
    totalLoss: 58500000,
    percentFixed: 18,
    trend: "up",
    topSeverity: "berat",
  },
  {
    rank: 2,
    name: "Pulo Gadung",
    totalReports: 38,
    totalLoss: 42300000,
    percentFixed: 25,
    trend: "up",
    topSeverity: "korban",
  },
  {
    rank: 3,
    name: "Tanah Abang",
    totalReports: 34,
    totalLoss: 38900000,
    percentFixed: 32,
    trend: "stable",
    topSeverity: "berat",
  },
  {
    rank: 4,
    name: "Senen",
    totalReports: 31,
    totalLoss: 35200000,
    percentFixed: 22,
    trend: "down",
    topSeverity: "sedang",
  },
  {
    rank: 5,
    name: "Mampang Prapatan",
    totalReports: 28,
    totalLoss: 31400000,
    percentFixed: 40,
    trend: "down",
    topSeverity: "sedang",
  },
  {
    rank: 6,
    name: "Tebet",
    totalReports: 25,
    totalLoss: 28100000,
    percentFixed: 52,
    trend: "down",
    topSeverity: "sedang",
  },
  {
    rank: 7,
    name: "Setiabudi",
    totalReports: 22,
    totalLoss: 24800000,
    percentFixed: 45,
    trend: "stable",
    topSeverity: "ringan",
  },
  {
    rank: 8,
    name: "Kebayoran Baru",
    totalReports: 19,
    totalLoss: 21200000,
    percentFixed: 58,
    trend: "down",
    topSeverity: "ringan",
  },
  {
    rank: 9,
    name: "Menteng",
    totalReports: 15,
    totalLoss: 17500000,
    percentFixed: 65,
    trend: "down",
    topSeverity: "ringan",
  },
  {
    rank: 10,
    name: "Cempaka Putih",
    totalReports: 12,
    totalLoss: 13800000,
    percentFixed: 70,
    trend: "stable",
    topSeverity: "ringan",
  },
];

export const dummyUserProfile: UserProfile = {
  name: "Budi Santoso",
  level: "Pak RT",
  levelIcon: "🏘️",
  points: 350,
  totalReports: 23,
  totalVerified: 18,
  totalFixed: 7,
  joinDate: "2025-10-15",
  badges: [
    {
      id: "1",
      name: "Pelopor",
      icon: "🏆",
      description: "Melaporkan 1 jalan rusak pertama",
      earned: true,
    },
    {
      id: "2",
      name: "Pantang Mundur",
      icon: "🔥",
      description: "Lapor 7 hari berturut-turut",
      earned: true,
    },
    {
      id: "3",
      name: "Mata Elang",
      icon: "🦅",
      description: "10 laporan terverifikasi",
      earned: true,
    },
    {
      id: "4",
      name: "Viral King",
      icon: "👑",
      description: "Laporan di-share 100 kali",
      earned: false,
    },
    {
      id: "5",
      name: "Fixer",
      icon: "🔧",
      description: "5 laporan selesai diperbaiki",
      earned: true,
    },
    {
      id: "6",
      name: "Menteri PU Swasta",
      icon: "🏗️",
      description: "Mencapai level tertinggi",
      earned: false,
    },
  ],
};

export function formatRupiah(amount: number): string {
  if (amount >= 1000000000) {
    return `Rp ${(amount / 1000000000).toFixed(1)} M`;
  }
  if (amount >= 1000000) {
    return `Rp ${(amount / 1000000).toFixed(1)} Jt`;
  }
  if (amount >= 1000) {
    return `Rp ${(amount / 1000).toFixed(0)} Rb`;
  }
  return `Rp ${amount.toLocaleString("id-ID")}`;
}

export function getAgeLabel(days: number): { label: string; icon: string } {
  if (days > 90) return { label: `${days} hari (Fosil!)`, icon: "🦖" };
  if (days > 30) return { label: `${days} hari (Ultah Lubang!)`, icon: "🎂" };
  if (days > 7) return { label: `${days} hari`, icon: "⏰" };
  return { label: `${days} hari`, icon: "🆕" };
}

export function getSeverityBadge(severity: string): {
  class: string;
  label: string;
} {
  switch (severity) {
    case "ringan":
      return { class: "badge-warning", label: "Rusak Ringan" };
    case "sedang":
      return { class: "badge-info", label: "Rusak Sedang" };
    case "berat":
      return { class: "badge-danger", label: "Rusak Berat" };
    case "korban":
      return { class: "badge-danger", label: "💀 Ada Korban" };
    default:
      return { class: "badge-info", label: severity };
  }
}

export function getStatusBadge(status: string): {
  class: string;
  label: string;
} {
  switch (status) {
    case "baru":
      return { class: "badge-info", label: "Baru" };
    case "diverifikasi":
      return { class: "badge-warning", label: "Terverifikasi" };
    case "proses":
      return { class: "badge-success", label: "Diperbaiki" };
    case "selesai":
      return { class: "badge-success", label: "✅ Selesai" };
    default:
      return { class: "badge-info", label: status };
  }
}
