export interface Report {
  id: string;
  title: string;
  description: string;
  lat: number;
  lng: number;
  latitude: number;
  longitude: number;
  location: string;
  kecamatan: string;
  kelurahan: string;
  severity: "ringan" | "sedang" | "berat" | "korban";
  status: "baru" | "diverifikasi" | "proses" | "selesai";
  photos: string[];
  reporterName: string;
  reporterLevel: string;
  reactions: number;
  comments: number;
  views: number;
  createdAt: string;
  daysOld: number;
  estimatedLoss: number;
  // API fields
  fingerprint_hash?: string;
  user_id?: string;
  district_id?: number;
  district_name?: string;
  image_url?: string;
  road_type?: string;
  view_count?: number;
  reaction_count?: number;
  estimated_loss?: number;
  days_old?: number;
  created_at?: string;
  updated_at?: string;
}

export interface Interaction {
  id: number;
  report_id: string;
  fingerprint_hash: string;
  type: string;
  value?: string;
  created_at: string;
}

export interface KecamatanRank {
  rank: number;
  name: string;
  totalReports: number;
  totalLoss: number;
  percentFixed: number;
  trend: "up" | "down" | "stable";
  topSeverity: "ringan" | "sedang" | "berat" | "korban";
}

export interface KecamatanRankAPI {
  rank: number;
  name: string;
  total_reports: number;
  total_loss: number;
  percent_fixed: number;
  top_severity: number;
}

export interface UserProfile {
  name: string;
  level: string;
  levelIcon: string;
  points: number;
  totalReports: number;
  totalVerified: number;
  totalFixed: number;
  joinDate: string;
  badges: Badge[];
}

export interface Badge {
  id: string;
  name: string;
  icon: string;
  description: string;
  earned: boolean;
}

export type SeverityColor = {
  [key in Report["severity"]]: string;
};

export const SEVERITY_LABELS: Record<Report["severity"], string> = {
  ringan: "Rusak Ringan",
  sedang: "Rusak Sedang",
  berat: "Rusak Berat",
  korban: "Ada Korban",
};

export const SEVERITY_COLORS: SeverityColor = {
  ringan: "#ECC94B",
  sedang: "#ED8936",
  berat: "#E53E3E",
  korban: "#1A202C",
};

export const STATUS_LABELS: Record<Report["status"], string> = {
  baru: "Baru Dilaporkan",
  diverifikasi: "Terverifikasi",
  proses: "Sedang Diperbaiki",
  selesai: "Selesai",
};

export const LEVEL_PROGRESSION = [
  { name: "Warga Biasa", icon: "👤", minPoints: 0 },
  { name: "Warga Peduli", icon: "👁️", minPoints: 50 },
  { name: "Pak RT", icon: "🏘️", minPoints: 200 },
  { name: "Pak Lurah", icon: "🏛️", minPoints: 500 },
  { name: "Camat Tangguh", icon: "⭐", minPoints: 1000 },
  { name: "Menteri PU Swasta", icon: "🏗️", minPoints: 2500 },
];

// Severity number to label mapping (from API)
export function severityFromNumber(n: number): Report["severity"] {
  switch (n) {
    case 1:
      return "ringan";
    case 2:
      return "sedang";
    case 3:
      return "berat";
    case 4:
    case 5:
      return "korban";
    default:
      return "ringan";
  }
}

export function severityToNumber(s: Report["severity"]): number {
  switch (s) {
    case "ringan":
      return 1;
    case "sedang":
      return 2;
    case "berat":
      return 3;
    case "korban":
      return 5;
    default:
      return 1;
  }
}
