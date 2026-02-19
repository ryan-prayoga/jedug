<script lang="ts">
  import { onMount, onDestroy } from 'svelte';
  import { browser } from '$app/environment';
  import ReportCard from '$lib/components/ReportCard.svelte';
  import { dummyReports } from '$lib/data/mock';
  import { MapPin, List, X, AlertTriangle, Crosshair, Loader2, Navigation } from 'lucide-svelte';
  import type { Report } from '$lib/types';
  import { api } from '$lib/api/client';
  import { severityFromNumber } from '$lib/types';

  let mapContainer: HTMLDivElement;
  let map: any;
  let L: any;
  let mapLoaded = false;
  let showSidebar = false;
  let selectedFilter: string = 'semua';
  let selectedReport: Report | null = null;
  let markerLayer: any = null;
  let userMarker: any = null;
  let reports: Report[] = [];
  let useAPI = false;

  // Location state
  let locationState: 'idle' | 'requesting' | 'granted' | 'denied' | 'error' = 'idle';
  let showLocationPrompt = true;
  let userLat = 0;
  let userLng = 0;
  let locationError = '';

  // Map move loading
  let loadingReports = false;
  let moveDebounceTimer: any = null;
  const DEFAULT_RADIUS = 20000; // 20km default radius

  const filters = [
    { value: 'semua', label: 'Semua' },
    { value: 'berat', label: '🔴 Berat' },
    { value: 'sedang', label: '🟠 Sedang' },
    { value: 'ringan', label: '🟡 Ringan' },
    { value: 'korban', label: '💀 Korban' }
  ];

  $: filteredReports = selectedFilter === 'semua'
    ? reports
    : reports.filter(r => r.severity === selectedFilter);

  onMount(async () => {
    if (!browser) return;

    L = await import('leaflet');
    await import('leaflet/dist/leaflet.css');

    // Initialize map with default center (Tangerang)
    map = L.map(mapContainer, {
      center: [-6.1781, 106.6297],
      zoom: 13,
      zoomControl: false,
      attributionControl: false
    });

    // Dark OSM tile layer (CartoDB Dark Matter)
    L.tileLayer('https://{s}.basemaps.cartocdn.com/dark_all/{z}/{x}/{y}{r}.png', {
      maxZoom: 19,
      subdomains: 'abcd'
    }).addTo(map);

    L.control.zoom({ position: 'topright' }).addTo(map);

    mapLoaded = true;

    // Listen for map moves to load reports in that area
    map.on('moveend', handleMapMove);

    // Show location prompt — don't auto-request, let user see the instructions first
    // locationState stays 'idle', showLocationPrompt stays true
  });

  function requestLocation() {
    if (!navigator.geolocation) {
      locationState = 'error';
      locationError = 'Browser kamu tidak mendukung geolokasi.';
      showLocationPrompt = false;
      // Load reports for default area
      loadReportsForArea(-6.1781, 106.6297, DEFAULT_RADIUS);
      return;
    }

    locationState = 'requesting';
    showLocationPrompt = false;

    navigator.geolocation.getCurrentPosition(
      (position) => {
        userLat = position.coords.latitude;
        userLng = position.coords.longitude;
        locationState = 'granted';

        // Center map on user
        if (map) {
          map.setView([userLat, userLng], 14, { animate: true, duration: 1 });
          addUserMarker(userLat, userLng);
        }

        // Load reports near user
        loadReportsForArea(userLat, userLng, DEFAULT_RADIUS);
      },
      (error) => {
        switch (error.code) {
          case error.PERMISSION_DENIED:
            locationState = 'denied';
            locationError = 'Akses lokasi ditolak.';
            break;
          case error.POSITION_UNAVAILABLE:
            locationState = 'error';
            locationError = 'Lokasi tidak tersedia.';
            break;
          case error.TIMEOUT:
            locationState = 'error';
            locationError = 'Request timeout.';
            break;
          default:
            locationState = 'error';
            locationError = 'Gagal mendapatkan lokasi.';
        }

        // Load reports for default area anyway
        loadReportsForArea(-6.1781, 106.6297, DEFAULT_RADIUS);
      },
      {
        enableHighAccuracy: true,
        timeout: 10000,
        maximumAge: 60000
      }
    );
  }

  function addUserMarker(lat: number, lng: number) {
    if (userMarker) map.removeLayer(userMarker);

    const icon = L.divIcon({
      className: 'user-location-marker',
      html: `
        <div class="user-dot-wrapper">
          <div class="user-dot-pulse"></div>
          <div class="user-dot"></div>
        </div>
      `,
      iconSize: [24, 24],
      iconAnchor: [12, 12]
    });

    userMarker = L.marker([lat, lng], { icon, zIndexOffset: 1000 }).addTo(map);
  }

  function handleMapMove() {
    if (!map || loadingReports) return;

    // Debounce to avoid too many requests
    clearTimeout(moveDebounceTimer);
    moveDebounceTimer = setTimeout(() => {
      const center = map.getCenter();
      const bounds = map.getBounds();
      
      // Calculate visible radius from bounds
      const ne = bounds.getNorthEast();
      const sw = bounds.getSouthWest();
      const latDiff = Math.abs(ne.lat - sw.lat);
      const lngDiff = Math.abs(ne.lng - sw.lng);
      const maxDiff = Math.max(latDiff, lngDiff);
      // Rough km per degree ~111km at equator
      const radiusKm = (maxDiff / 2) * 111;
      const radiusMeters = Math.min(radiusKm * 1000, 50000); // Max 50km

      loadReportsForArea(center.lat, center.lng, radiusMeters);
    }, 700);
  }

  async function loadReportsForArea(lat: number, lng: number, radius: number) {
    loadingReports = true;
    
    try {
      const result = await api.getNearbyReports(lat, lng, radius);
      if (result.data && result.data.length > 0) {
        reports = result.data.map(mapAPIReportToLocal);
        useAPI = true;
        addMarkers();
      } else if (!useAPI) {
        // First load with no API data — show mock data
        reports = dummyReports;
        addMarkers();
      }
    } catch {
      if (!useAPI && reports.length === 0) {
        // API failed, use mock data as fallback
        reports = dummyReports;
        addMarkers();
        console.log('Using mock data (API not available)');
      }
    }
    
    loadingReports = false;
  }

  function mapAPIReportToLocal(d: any): Report {
    return {
      ...d,
      id: d.id,
      title: d.description ? d.description.substring(0, 60) : `Laporan #${d.id.substring(0, 8)}`,
      description: d.description || '',
      lat: d.latitude || d.lat,
      lng: d.longitude || d.lng,
      location: d.district_name || 'Lokasi',
      kecamatan: d.district_name || '',
      kelurahan: '',
      severity: severityFromNumber(d.severity),
      status: d.status === 'open' ? 'baru' : d.status === 'fixed' ? 'selesai' : d.status,
      photos: d.photos || [d.image_url],
      reporterName: 'Warga',
      reporterLevel: 'Warga Biasa',
      reactions: d.reaction_count || 0,
      comments: 0,
      views: d.view_count || 0,
      createdAt: d.created_at,
      daysOld: d.days_old || Math.floor((Date.now() - new Date(d.created_at).getTime()) / 86400000),
      estimatedLoss: d.estimated_loss || 0,
    };
  }

  function addMarkers() {
    if (!L || !map) return;
    if (markerLayer) map.removeLayer(markerLayer);
    markerLayer = L.layerGroup();

    reports.forEach(report => {
      const color = getSeverityColor(report.severity);
      const icon = L.divIcon({
        className: 'custom-marker',
        html: `
          <div style="
            width: 36px; height: 36px;
            background: ${color};
            border: 3px solid white;
            border-radius: 50%;
            box-shadow: 0 2px 8px rgba(0,0,0,0.3);
            display: flex; align-items: center; justify-content: center;
            cursor: pointer;
            font-size: 14px;
          ">
            ${report.severity === 'korban' ? '💀' : '🕳️'}
          </div>
        `,
        iconSize: [36, 36],
        iconAnchor: [18, 18],
        popupAnchor: [0, -20]
      });

      const marker = L.marker([report.lat, report.lng], { icon })
        .bindPopup(`
          <div style="padding:8px;max-width:220px;font-family:Inter,sans-serif;">
            <strong style="font-size:13px;display:block;margin-bottom:4px;">${report.title}</strong>
            <span style="font-size:11px;color:#666;">${report.location}</span>
            <div style="margin-top:6px;font-size:11px;color:#999;">
              👁️ ${report.views} · ❤️ ${report.reactions} · 💬 ${report.comments}
            </div>
          </div>
        `, { closeButton: false, offset: [0, -10] });

      marker.on('click', () => {
        selectedReport = report;
      });

      marker.addTo(markerLayer);
    });

    markerLayer.addTo(map);
  }

  function getSeverityColor(severity: string): string {
    switch (severity) {
      case 'ringan': return '#ECC94B';
      case 'sedang': return '#ED8936';
      case 'berat': return '#E53E3E';
      case 'korban': return '#1A202C';
      default: return '#3182CE';
    }
  }

  function selectReport(report: Report) {
    selectedReport = report;
    if (map) {
      map.setView([report.lat, report.lng], 15, { animate: true, duration: 0.8 });
    }
  }

  function goToMyLocation() {
    if (userLat && userLng && map) {
      map.setView([userLat, userLng], 14, { animate: true, duration: 0.8 });
    } else {
      requestLocation();
    }
  }

  onDestroy(() => {
    clearTimeout(moveDebounceTimer);
    if (map) map.remove();
  });
</script>

<svelte:head>
  <title>JEDUG - Peta Jalan Rusak</title>
</svelte:head>

<div class="map-page">
  <!-- Map Container -->
  <div class="map-wrapper">
    <div class="map-container" bind:this={mapContainer}></div>

    {#if !mapLoaded}
      <div class="map-loading">
        <div class="loading-spinner"></div>
        <p>Memuat peta...</p>
      </div>
    {/if}

    <!-- Location Permission Prompt -->
    {#if showLocationPrompt && locationState === 'idle'}
      <div class="location-prompt">
        <div class="prompt-card">
          <div class="prompt-icon">📍</div>
          <h3>Izinkan Akses Lokasi</h3>
          <p>JEDUG membutuhkan lokasi untuk menampilkan <strong>lubang jalan di sekitar kamu</strong> (radius 20km).</p>
          <div class="prompt-steps">
            <div class="step-item">
              <span class="step-num">1</span>
              <span>Tekan tombol di bawah</span>
            </div>
            <div class="step-item">
              <span class="step-num">2</span>
              <span>Browser akan minta izin — tekan <strong>"Allow"</strong> atau <strong>"Izinkan"</strong></span>
            </div>
            <div class="step-item">
              <span class="step-num">3</span>
              <span>Peta otomatis pindah ke lokasi kamu</span>
            </div>
          </div>
          <button class="prompt-btn" on:click={requestLocation}>
            <Navigation size={18} />
            Gunakan Lokasi Saya
          </button>
          <button class="prompt-skip" on:click={() => { showLocationPrompt = false; loadReportsForArea(-6.1781, 106.6297, DEFAULT_RADIUS); }}>
            Lewati, lihat peta Tangerang
          </button>
        </div>
      </div>
    {/if}

    <!-- Location Requesting Overlay -->
    {#if locationState === 'requesting'}
      <div class="location-requesting">
        <div class="requesting-card">
          <div class="requesting-spinner"></div>
          <p>Mendeteksi lokasi kamu...</p>
          <span class="requesting-hint">Jika muncul popup, tekan <strong>"Allow"</strong></span>
        </div>
      </div>
    {/if}

    <!-- Location Denied Banner -->
    {#if locationState === 'denied'}
      <div class="location-denied-banner">
        <div class="denied-content">
          <span>📍 Lokasi ditolak</span>
          <p>Buka Settings → Privacy → Location, lalu izinkan browser ini mengakses lokasi.</p>
        </div>
        <button class="denied-retry" on:click={requestLocation}>Coba Lagi</button>
      </div>
    {/if}

    <!-- Map Controls (Mobile) -->
    <div class="map-controls">
      <button class="map-btn" on:click={() => showSidebar = !showSidebar} aria-label="Toggle daftar laporan">
        {#if showSidebar}
          <X size={20} />
        {:else}
          <List size={20} />
        {/if}
      </button>

      <!-- My Location Button -->
      <button class="map-btn location-btn-map" on:click={goToMyLocation} aria-label="Ke lokasi saya" class:active={locationState === 'granted'}>
        <Crosshair size={20} />
      </button>
    </div>

    <!-- Filter Bar -->
    <div class="filter-bar">
      {#each filters as f}
        <button
          class="filter-chip"
          class:active={selectedFilter === f.value}
          on:click={() => selectedFilter = f.value}
        >
          {f.label}
        </button>
      {/each}
    </div>

    <!-- Stats Overlay -->
    <div class="stats-overlay">
      {#if loadingReports}
        <div class="stat-pill loading">
          <Loader2 size={14} class="spin-icon" />
          <span>Memuat...</span>
        </div>
      {/if}
      <div class="stat-pill danger">
        <AlertTriangle size={14} />
        <span>{reports.filter(r => r.severity === 'berat' || r.severity === 'korban').length} Parah</span>
      </div>
      <div class="stat-pill">
        <MapPin size={14} />
        <span>{reports.length} Laporan</span>
      </div>
    </div>

    <!-- Legend -->
    <div class="map-legend">
      <div class="legend-title">Tingkat Kerusakan</div>
      <div class="legend-items">
        <div class="legend-item">
          <span class="legend-dot" style="background: #ECC94B;"></span>
          <span>Ringan</span>
        </div>
        <div class="legend-item">
          <span class="legend-dot" style="background: #ED8936;"></span>
          <span>Sedang</span>
        </div>
        <div class="legend-item">
          <span class="legend-dot" style="background: #E53E3E;"></span>
          <span>Berat</span>
        </div>
        <div class="legend-item">
          <span class="legend-dot" style="background: #1A202C; border: 1px solid #666;"></span>
          <span>Korban</span>
        </div>
      </div>
    </div>
  </div>

  <!-- Sidebar (Desktop always visible, Mobile toggle) -->
  <aside class="sidebar" class:open={showSidebar}>
    <div class="sidebar-header">
      <h2>Laporan Terbaru</h2>
      <span class="report-count">{filteredReports.length} laporan</span>
      <button class="close-sidebar hide-desktop" on:click={() => showSidebar = false}>
        <X size={20} />
      </button>
    </div>
    <div class="sidebar-list">
      {#each filteredReports as report (report.id)}
        <div on:click={() => selectReport(report)} on:keydown={() => selectReport(report)} role="button" tabindex="0">
          <ReportCard {report} compact />
        </div>
      {/each}
      {#if filteredReports.length === 0}
        <div class="empty-state">
          <span class="empty-icon">🔍</span>
          {#if loadingReports}
            <p>Memuat laporan sekitar...</p>
          {:else}
            <p>Tidak ada laporan di area ini</p>
            <p class="empty-hint">Coba geser peta ke area lain</p>
          {/if}
        </div>
      {/if}
    </div>
  </aside>
</div>

<!-- Mobile Bottom Sheet (selected report) -->
{#if selectedReport}
  <div class="mobile-sheet hide-desktop" on:click={() => selectedReport = null} on:keydown={() => selectedReport = null} role="button" tabindex="0">
    <div class="sheet-content" on:click|stopPropagation on:keydown|stopPropagation role="presentation">
      <div class="sheet-handle"></div>
      <ReportCard report={selectedReport} />
    </div>
  </div>
{/if}

<style>
  .map-page {
    display: flex;
    height: calc(100dvh - var(--nav-height));
    position: relative;
    margin-bottom: calc(-1 * var(--bottom-nav-height));
  }

  .map-wrapper {
    flex: 1;
    position: relative;
    min-height: 0;
  }

  .map-container {
    width: 100%;
    height: 100%;
  }

  :global(.leaflet-popup-content-wrapper) {
    border-radius: 12px !important;
    box-shadow: 0 4px 12px rgba(0,0,0,0.15) !important;
    padding: 0 !important;
  }
  :global(.leaflet-popup-tip) {
    border-top-color: white !important;
  }
  :global(.leaflet-control-zoom) {
    border: 1px solid var(--border-color) !important;
    border-radius: var(--radius-lg) !important;
    overflow: hidden;
  }
  :global(.leaflet-control-zoom a) {
    background: var(--bg-card) !important;
    color: var(--text-primary) !important;
    border-bottom-color: var(--border-color) !important;
  }

  /* User Location Marker */
  :global(.user-dot-wrapper) {
    position: relative;
    width: 24px;
    height: 24px;
  }
  :global(.user-dot) {
    position: absolute;
    top: 50%;
    left: 50%;
    transform: translate(-50%, -50%);
    width: 14px;
    height: 14px;
    background: #4285F4;
    border: 3px solid white;
    border-radius: 50%;
    box-shadow: 0 2px 6px rgba(66, 133, 244, 0.5);
    z-index: 2;
  }
  :global(.user-dot-pulse) {
    position: absolute;
    top: 50%;
    left: 50%;
    transform: translate(-50%, -50%);
    width: 40px;
    height: 40px;
    background: rgba(66, 133, 244, 0.2);
    border-radius: 50%;
    animation: pulse 2s ease-out infinite;
    z-index: 1;
  }
  @keyframes pulse {
    0% { transform: translate(-50%, -50%) scale(0.5); opacity: 1; }
    100% { transform: translate(-50%, -50%) scale(1.5); opacity: 0; }
  }

  /* Map Loading */
  .map-loading {
    position: absolute;
    inset: 0;
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    background: var(--bg-secondary);
    gap: var(--space-md);
    color: var(--text-secondary);
    z-index: 5;
  }
  .loading-spinner {
    width: 40px;
    height: 40px;
    border: 3px solid var(--border-color);
    border-top-color: var(--color-primary);
    border-radius: 50%;
    animation: spin 0.8s linear infinite;
  }
  @keyframes spin {
    to { transform: rotate(360deg); }
  }

  /* Location Prompt */
  .location-prompt {
    position: absolute;
    inset: 0;
    z-index: 1100;
    background: rgba(0,0,0,0.5);
    display: flex;
    align-items: center;
    justify-content: center;
    padding: var(--space-lg);
  }
  .prompt-card {
    background: var(--bg-card);
    border-radius: var(--radius-2xl);
    padding: var(--space-2xl);
    max-width: 400px;
    width: 100%;
    text-align: center;
    box-shadow: var(--shadow-xl);
    border: 1px solid var(--border-color);
    animation: slideUp 0.4s ease;
  }
  @keyframes slideUp {
    from { transform: translateY(20px); opacity: 0; }
    to { transform: translateY(0); opacity: 1; }
  }
  .prompt-icon {
    font-size: 3rem;
    margin-bottom: var(--space-md);
  }
  .prompt-card h3 {
    font-size: var(--text-xl);
    font-weight: var(--font-bold);
    margin-bottom: var(--space-sm);
    color: var(--text-primary);
  }
  .prompt-card > p {
    font-size: var(--text-sm);
    color: var(--text-secondary);
    line-height: var(--leading-relaxed);
    margin-bottom: var(--space-xl);
  }
  .prompt-steps {
    display: flex;
    flex-direction: column;
    gap: var(--space-sm);
    margin-bottom: var(--space-xl);
    text-align: left;
  }
  .step-item {
    display: flex;
    align-items: center;
    gap: var(--space-md);
    font-size: var(--text-sm);
    color: var(--text-secondary);
  }
  .step-num {
    width: 28px;
    height: 28px;
    border-radius: 50%;
    background: var(--color-primary);
    color: white;
    display: flex;
    align-items: center;
    justify-content: center;
    font-size: var(--text-xs);
    font-weight: var(--font-bold);
    flex-shrink: 0;
  }
  .prompt-btn {
    width: 100%;
    padding: 0.875rem 1.5rem;
    background: var(--color-primary);
    color: white;
    border: none;
    border-radius: var(--radius-xl);
    font-size: var(--text-base);
    font-weight: var(--font-semibold);
    cursor: pointer;
    display: flex;
    align-items: center;
    justify-content: center;
    gap: var(--space-sm);
    transition: all var(--transition-fast);
    margin-bottom: var(--space-md);
  }
  .prompt-btn:hover {
    background: var(--color-primary-hover);
  }
  .prompt-btn:active {
    transform: scale(0.98);
  }
  .prompt-skip {
    background: none;
    border: none;
    color: var(--text-tertiary);
    font-size: var(--text-sm);
    cursor: pointer;
    text-decoration: underline;
    padding: var(--space-xs);
  }
  .prompt-skip:hover {
    color: var(--text-secondary);
  }

  /* Location Requesting */
  .location-requesting {
    position: absolute;
    top: var(--space-lg);
    left: 50%;
    transform: translateX(-50%);
    z-index: 1100;
    animation: slideDown 0.3s ease;
  }
  @keyframes slideDown {
    from { transform: translateX(-50%) translateY(-10px); opacity: 0; }
    to { transform: translateX(-50%) translateY(0); opacity: 1; }
  }
  .requesting-card {
    display: flex;
    align-items: center;
    gap: var(--space-md);
    padding: var(--space-md) var(--space-lg);
    background: var(--bg-card);
    border-radius: var(--radius-xl);
    box-shadow: var(--shadow-lg);
    border: 1px solid var(--border-color);
    font-size: var(--text-sm);
    color: var(--text-primary);
  }
  .requesting-spinner {
    width: 20px;
    height: 20px;
    border: 2px solid var(--border-color);
    border-top-color: var(--color-primary);
    border-radius: 50%;
    animation: spin 0.8s linear infinite;
    flex-shrink: 0;
  }
  .requesting-hint {
    font-size: var(--text-xs);
    color: var(--text-tertiary);
    margin-left: var(--space-sm);
  }

  /* Location Denied Banner */
  .location-denied-banner {
    position: absolute;
    top: var(--space-md);
    left: var(--space-md);
    right: var(--space-md);
    z-index: 1050;
    display: flex;
    align-items: center;
    justify-content: space-between;
    gap: var(--space-md);
    padding: var(--space-md) var(--space-lg);
    background: var(--color-warning-light);
    border: 1px solid var(--color-warning);
    border-radius: var(--radius-xl);
    box-shadow: var(--shadow-md);
    animation: slideDown 0.3s ease;
  }
  .denied-content span {
    font-size: var(--text-sm);
    font-weight: var(--font-semibold);
    color: var(--text-primary);
  }
  .denied-content p {
    font-size: var(--text-xs);
    color: var(--text-secondary);
    margin-top: 2px;
    line-height: var(--leading-relaxed);
  }
  .denied-retry {
    padding: 0.5rem 1rem;
    background: var(--color-primary);
    color: white;
    border: none;
    border-radius: var(--radius-lg);
    font-size: var(--text-xs);
    font-weight: var(--font-semibold);
    cursor: pointer;
    white-space: nowrap;
    flex-shrink: 0;
  }

  /* Map Controls */
  .map-controls {
    position: absolute;
    top: var(--space-md);
    left: var(--space-md);
    z-index: 1000;
    display: flex;
    flex-direction: column;
    gap: var(--space-sm);
  }
  .map-btn {
    width: 44px;
    height: 44px;
    border-radius: var(--radius-lg);
    background: var(--bg-card);
    color: var(--text-primary);
    display: flex;
    align-items: center;
    justify-content: center;
    box-shadow: var(--shadow-md);
    border: 1px solid var(--border-color);
    cursor: pointer;
    transition: all var(--transition-fast);
  }
  .map-btn:active {
    transform: scale(0.95);
  }
  .map-btn.location-btn-map.active {
    color: #4285F4;
    border-color: #4285F4;
  }

  /* Filter Bar */
  .filter-bar {
    position: absolute;
    bottom: var(--space-lg);
    left: 50%;
    transform: translateX(-50%);
    display: flex;
    gap: var(--space-xs);
    padding: var(--space-xs);
    background: var(--bg-nav);
    backdrop-filter: blur(12px);
    border-radius: var(--radius-full);
    box-shadow: var(--shadow-lg);
    border: 1px solid var(--border-color);
    z-index: 1000;
    max-width: calc(100% - 2rem);
    overflow-x: auto;
    scrollbar-width: none;
  }
  .filter-bar::-webkit-scrollbar { display: none; }

  .filter-chip {
    padding: 0.375rem 0.875rem;
    border-radius: var(--radius-full);
    font-size: var(--text-xs);
    font-weight: var(--font-medium);
    white-space: nowrap;
    color: var(--text-secondary);
    background: transparent;
    transition: all var(--transition-fast);
    cursor: pointer;
    border: none;
  }
  .filter-chip:hover {
    background: var(--bg-tertiary);
  }
  .filter-chip.active {
    background: var(--color-primary);
    color: white;
  }

  /* Stats Overlay */
  .stats-overlay {
    position: absolute;
    top: var(--space-md);
    right: 80px;
    display: flex;
    gap: var(--space-xs);
    z-index: 1000;
  }
  .stat-pill {
    display: flex;
    align-items: center;
    gap: 4px;
    padding: 0.375rem 0.75rem;
    background: var(--bg-nav);
    backdrop-filter: blur(12px);
    border-radius: var(--radius-full);
    font-size: var(--text-xs);
    font-weight: var(--font-semibold);
    color: var(--text-secondary);
    box-shadow: var(--shadow-sm);
    border: 1px solid var(--border-color);
  }
  .stat-pill.danger {
    background: var(--color-danger-light);
    color: var(--color-danger);
    border-color: transparent;
  }
  .stat-pill.loading {
    background: var(--color-primary-light);
    color: var(--color-primary);
    border-color: transparent;
  }
  :global(.spin-icon) {
    animation: spin 1s linear infinite;
  }

  /* Legend */
  .map-legend {
    position: absolute;
    bottom: var(--space-lg);
    left: var(--space-md);
    background: var(--bg-nav);
    backdrop-filter: blur(12px);
    border-radius: var(--radius-lg);
    padding: var(--space-sm) var(--space-md);
    box-shadow: var(--shadow-md);
    border: 1px solid var(--border-color);
    z-index: 1000;
    font-size: var(--text-xs);
  }
  .legend-title {
    font-weight: var(--font-semibold);
    color: var(--text-primary);
    margin-bottom: var(--space-xs);
    font-size: var(--text-xs);
  }
  .legend-items {
    display: flex;
    flex-direction: column;
    gap: 4px;
  }
  .legend-item {
    display: flex;
    align-items: center;
    gap: var(--space-sm);
    color: var(--text-secondary);
  }
  .legend-dot {
    width: 12px;
    height: 12px;
    border-radius: 50%;
    flex-shrink: 0;
  }

  @media (max-width: 768px) {
    .map-page {
      height: calc(100dvh - var(--nav-height) - var(--bottom-nav-height));
    }

    .map-legend {
      bottom: 4.5rem;
      left: var(--space-sm);
      padding: var(--space-xs) var(--space-sm);
      font-size: 0.65rem;
      max-width: 120px;
    }

    .filter-bar {
      bottom: var(--space-sm);
      padding: 3px;
      gap: 2px;
    }
    .filter-chip {
      padding: 0.25rem 0.625rem;
      font-size: 0.65rem;
    }

    .stats-overlay {
      right: var(--space-md);
      left: auto;
      top: auto;
      bottom: 3.5rem;
      flex-direction: row;
      gap: 4px;
    }
    .stat-pill {
      font-size: 0.65rem;
      padding: 0.25rem 0.5rem;
    }

    .location-denied-banner {
      flex-direction: column;
      align-items: flex-start;
    }
  }

  /* Sidebar */
  .sidebar {
    width: var(--sidebar-width);
    background: var(--bg-secondary);
    border-left: 1px solid var(--border-color);
    display: flex;
    flex-direction: column;
    overflow: hidden;
  }
  .sidebar-header {
    display: flex;
    align-items: center;
    gap: var(--space-md);
    padding: var(--space-lg);
    border-bottom: 1px solid var(--border-color);
    background: var(--bg-card);
  }
  .sidebar-header h2 {
    font-size: var(--text-lg);
    font-weight: var(--font-bold);
    flex: 1;
  }
  .report-count {
    font-size: var(--text-xs);
    color: var(--text-tertiary);
    background: var(--bg-tertiary);
    padding: 0.25rem 0.625rem;
    border-radius: var(--radius-full);
  }
  .close-sidebar {
    color: var(--text-secondary);
    cursor: pointer;
    border: none;
    background: none;
    padding: var(--space-xs);
  }

  .sidebar-list {
    flex: 1;
    overflow-y: auto;
    padding: var(--space-md);
    display: flex;
    flex-direction: column;
    gap: var(--space-md);
  }

  .empty-state {
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    padding: var(--space-3xl) var(--space-lg);
    color: var(--text-tertiary);
    text-align: center;
  }
  .empty-icon {
    font-size: 2.5rem;
    margin-bottom: var(--space-md);
  }
  .empty-hint {
    font-size: var(--text-xs);
    margin-top: var(--space-xs);
    color: var(--text-tertiary);
    opacity: 0.7;
  }

  /* Mobile Sidebar */
  @media (max-width: 768px) {
    .sidebar {
      position: fixed;
      top: var(--nav-height);
      right: 0;
      bottom: 0;
      width: 100%;
      max-width: 380px;
      z-index: var(--z-fab);
      transform: translateX(100%);
      transition: transform var(--transition-base);
      box-shadow: var(--shadow-xl);
      padding-bottom: var(--bottom-nav-height);
    }
    .sidebar.open {
      transform: translateX(0);
    }
  }

  /* Mobile Bottom Sheet */
  .mobile-sheet {
    position: fixed;
    inset: 0;
    z-index: var(--z-modal);
    background: var(--bg-overlay);
    display: flex;
    align-items: flex-end;
  }
  .sheet-content {
    background: var(--bg-primary);
    border-radius: var(--radius-2xl) var(--radius-2xl) 0 0;
    padding: var(--space-md);
    padding-bottom: calc(var(--bottom-nav-height) + var(--space-md));
    width: 100%;
    max-height: 70vh;
    overflow-y: auto;
  }
  .sheet-handle {
    width: 40px;
    height: 4px;
    background: var(--border-color-strong);
    border-radius: var(--radius-full);
    margin: 0 auto var(--space-md);
  }
</style>
