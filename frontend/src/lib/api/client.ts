const API_BASE = "/api/v1";

async function request<T>(endpoint: string, options?: RequestInit): Promise<T> {
  const url = `${API_BASE}${endpoint}`;
  const res = await fetch(url, {
    headers: {
      "Content-Type": "application/json",
      ...options?.headers,
    },
    ...options,
  });

  if (!res.ok) {
    const error = await res.json().catch(() => ({ message: res.statusText }));
    throw new Error(error.message || `API Error: ${res.status}`);
  }

  return res.json();
}

async function uploadFile(
  file: File | Blob,
): Promise<{ url: string; filename: string }> {
  const formData = new FormData();
  formData.append("image", file);

  const res = await fetch(`${API_BASE}/upload`, {
    method: "POST",
    body: formData,
  });

  if (!res.ok) {
    const error = await res.json().catch(() => ({ message: res.statusText }));
    throw new Error(error.message || "Upload failed");
  }

  return res.json();
}

export const api = {
  get: <T>(endpoint: string) => request<T>(endpoint),

  post: <T>(endpoint: string, body: unknown) =>
    request<T>(endpoint, {
      method: "POST",
      body: JSON.stringify(body),
    }),

  put: <T>(endpoint: string, body: unknown) =>
    request<T>(endpoint, {
      method: "PUT",
      body: JSON.stringify(body),
    }),

  delete: <T>(endpoint: string) => request<T>(endpoint, { method: "DELETE" }),

  health: () =>
    request<{ status: string; message: string; db: string }>("/health"),

  // === Reports ===
  listReports: (params?: {
    severity?: string;
    status?: string;
    district_id?: number;
    limit?: number;
    offset?: number;
  }) => {
    const searchParams = new URLSearchParams();
    if (params?.severity) searchParams.set("severity", params.severity);
    if (params?.status) searchParams.set("status", params.status);
    if (params?.district_id)
      searchParams.set("district_id", String(params.district_id));
    if (params?.limit) searchParams.set("limit", String(params.limit));
    if (params?.offset) searchParams.set("offset", String(params.offset));
    const qs = searchParams.toString();
    return request<{
      data: import("$lib/types").Report[];
      total: number;
      limit: number;
      offset: number;
      has_more: boolean;
    }>(`/reports${qs ? `?${qs}` : ""}`);
  },

  getReport: (id: string) =>
    request<{ data: import("$lib/types").Report }>(`/reports/${id}`),

  createReport: (data: {
    fingerprint_hash: string;
    latitude: number;
    longitude: number;
    image_url: string;
    severity: number;
    description?: string;
  }) =>
    request<{ data: import("$lib/types").Report; message: string }>(
      "/reports",
      data,
    ),

  getNearbyReports: (lat: number, lng: number, radius?: number) => {
    const params = new URLSearchParams({
      lat: String(lat),
      lng: String(lng),
      radius: String(radius || 1000),
    });
    return request<{ data: import("$lib/types").Report[] }>(
      `/reports/nearby?${params}`,
    );
  },

  // === Interactions ===
  addReaction: (reportId: string, fingerprint: string, type: string) =>
    request<{ message: string }>(`/reports/${reportId}/reactions`, {
      method: "POST",
      body: JSON.stringify({ fingerprint_hash: fingerprint, type }),
    }),

  getComments: (reportId: string) =>
    request<{ data: import("$lib/types").Interaction[] }>(
      `/reports/${reportId}/comments`,
    ),

  addComment: (reportId: string, fingerprint: string, text: string) =>
    request<{ data: import("$lib/types").Interaction }>(
      `/reports/${reportId}/comments`,
      {
        method: "POST",
        body: JSON.stringify({ fingerprint_hash: fingerprint, text }),
      },
    ),

  // === Ranking ===
  getRanking: (limit?: number) =>
    request<{ data: import("$lib/types").KecamatanRankAPI[] }>(
      `/ranking${limit ? `?limit=${limit}` : ""}`,
    ),

  // === Upload ===
  uploadImage: (file: File | Blob) => uploadFile(file),
};
