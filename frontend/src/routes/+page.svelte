<script lang="ts">
  import { onMount, onDestroy } from 'svelte';
  import { browser } from '$app/environment';
  import ReportCard from '$lib/components/ReportCard.svelte';
  import { dummyReports } from '$lib/data/mock';
  import { MapPin, Filter, List, X, ChevronDown, Layers, AlertTriangle } from 'lucide-svelte';
  import type { Report } from '$lib/types';

  let mapContainer: HTMLDivElement;
  let map: any;
  let L: any;
  let mapLoaded = false;
  let showSidebar = false;
  let selectedFilter: string = 'semua';
  let selectedReport: Report | null = null;
  let districtLayer: any = null;
  let roadLayer: any = null;
  let markerLayer: any = null;

  const filters = [
    { value: 'semua', label: 'Semua' },
    { value: 'berat', label: '🔴 Berat' },
    { value: 'sedang', label: '🟠 Sedang' },
    { value: 'ringan', label: '🟡 Ringan' },
    { value: 'korban', label: '💀 Korban' }
  ];

  $: filteredReports = selectedFilter === 'semua'
    ? dummyReports
    : dummyReports.filter(r => r.severity === selectedFilter);

  onMount(async () => {
    if (!browser) return;

    L = await import('leaflet');
    await import('leaflet/dist/leaflet.css');

    map = L.map(mapContainer, {
      center: [-6.1781, 106.6297],
      zoom: 12,
      zoomControl: false,
      attributionControl: false
    });

    // Dark OSM tile layer (CartoDB Dark Matter)
    L.tileLayer('https://{s}.basemaps.cartocdn.com/dark_all/{z}/{x}/{y}{r}.png', {
      maxZoom: 19,
      subdomains: 'abcd'
    }).addTo(map);

    // Zoom control kanan atas
    L.control.zoom({ position: 'topright' }).addTo(map);

    mapLoaded = true;
    if (browser) (window as any).__map = map;

    loadDistrictBoundaries();
    addRoadSegments();
    addMarkers();
  });

  async function loadDistrictBoundaries() {
    try {
      const res = await fetch('/api/v1/districts?city=kota-tangerang');
      if (!res.ok) throw new Error(`Failed to fetch districts: ${res.status}`);
      const data = await res.json();
      const geojson = data.geojson ?? data;

      // Hapus layer lama jika ada (HMR reload)
      if (districtLayer) {
        map.removeLayer(districtLayer);
      }

      districtLayer = L.geoJSON(geojson, {
        style: () => ({
          fillColor: '#4299E1',
          fillOpacity: 0.15,
          color: '#63B3ED',
          weight: 2,
          opacity: 0.7
        }),
        onEachFeature: (feature: any, layer: any) => {
          if (feature.properties?.name) {
            layer.bindTooltip(feature.properties.name.toUpperCase(), {
              permanent: true,
              direction: 'center',
              className: 'district-label'
            });
          }
        }
      }).addTo(map);
    } catch (err) {
      console.error('Failed to load district boundaries:', err);
    }
  }

  const roadSegments = [
    { id: 'road-sudirman-1', color: '#38A169', coords: [
        [106.8019, -6.2270], [106.8022, -6.2250], [106.8025, -6.2230],
        [106.8029, -6.2210], [106.8033, -6.2190], [106.8038, -6.2170],
        [106.8042, -6.2150], [106.8046, -6.2130], [106.8050, -6.2110],
        [106.8054, -6.2090], [106.8058, -6.2070], [106.8062, -6.2050],
        [106.8066, -6.2030], [106.8069, -6.2010]
      ]
    },
    { id: 'road-thamrin', color: '#ECC94B', coords: [
        [106.8069, -6.2010], [106.8072, -6.1990], [106.8075, -6.1970],
        [106.8078, -6.1950], [106.8080, -6.1935], [106.8082, -6.1920],
        [106.8084, -6.1900], [106.8085, -6.1880], [106.8086, -6.1860],
        [106.8087, -6.1845], [106.8088, -6.1830]
      ]
    },
    { id: 'road-gatsu-barat', color: '#E53E3E', coords: [
        [106.7930, -6.1960], [106.7950, -6.1975], [106.7970, -6.1990],
        [106.7990, -6.2005], [106.8010, -6.2020], [106.8030, -6.2040],
        [106.8045, -6.2055], [106.8060, -6.2065], [106.8069, -6.2070]
      ]
    },
    { id: 'road-gatsu-timur', color: '#ED8936', coords: [
        [106.8069, -6.2070], [106.8090, -6.2085], [106.8120, -6.2100],
        [106.8150, -6.2115], [106.8180, -6.2130], [106.8210, -6.2148],
        [106.8240, -6.2165], [106.8270, -6.2180], [106.8300, -6.2200],
        [106.8330, -6.2220], [106.8350, -6.2240], [106.8370, -6.2260],
        [106.8385, -6.2290], [106.8395, -6.2320], [106.8400, -6.2350]
      ]
    },
    { id: 'road-rasuna', color: '#38A169', coords: [
        [106.8260, -6.2040], [106.8270, -6.2060], [106.8275, -6.2080],
        [106.8278, -6.2100], [106.8280, -6.2120], [106.8282, -6.2140],
        [106.8284, -6.2160], [106.8286, -6.2180], [106.8290, -6.2200],
        [106.8295, -6.2220], [106.8300, -6.2240], [106.8310, -6.2260]
      ]
    },
    { id: 'road-casablanca', color: '#E53E3E', coords: [
        [106.8260, -6.2210], [106.8240, -6.2220], [106.8220, -6.2230],
        [106.8200, -6.2245], [106.8180, -6.2255], [106.8160, -6.2265],
        [106.8140, -6.2275], [106.8120, -6.2280]
      ]
    },
    { id: 'road-sparman', color: '#38A169', coords: [
        [106.7930, -6.1960], [106.7910, -6.1945], [106.7890, -6.1930],
        [106.7870, -6.1915], [106.7850, -6.1900], [106.7830, -6.1885],
        [106.7810, -6.1870], [106.7790, -6.1855]
      ]
    },
    { id: 'road-monas-barat', color: '#ECC94B', coords: [
        [106.8220, -6.1770], [106.8220, -6.1790], [106.8220, -6.1810],
        [106.8220, -6.1830], [106.8220, -6.1850], [106.8220, -6.1870]
      ]
    },
    { id: 'road-monas-timur', color: '#38A169', coords: [
        [106.8330, -6.1770], [106.8330, -6.1790], [106.8330, -6.1810],
        [106.8330, -6.1830], [106.8330, -6.1850], [106.8330, -6.1870]
      ]
    },
    { id: 'road-monas-utara', color: '#38A169', coords: [
        [106.8220, -6.1770], [106.8240, -6.1770], [106.8260, -6.1770],
        [106.8280, -6.1770], [106.8300, -6.1770], [106.8330, -6.1770]
      ]
    },
    { id: 'road-monas-selatan', color: '#ED8936', coords: [
        [106.8220, -6.1870], [106.8240, -6.1870], [106.8260, -6.1870],
        [106.8280, -6.1870], [106.8300, -6.1870], [106.8330, -6.1870]
      ]
    },
    { id: 'road-mtharyono', color: '#1A202C', coords: [
        [106.8400, -6.2350], [106.8420, -6.2370], [106.8440, -6.2390],
        [106.8460, -6.2410], [106.8480, -6.2430], [106.8500, -6.2450],
        [106.8520, -6.2470], [106.8540, -6.2490]
      ]
    },
    { id: 'road-sudirman-2', color: '#ED8936', coords: [
        [106.8019, -6.2270], [106.8015, -6.2290], [106.8010, -6.2310],
        [106.8005, -6.2330], [106.8000, -6.2350], [106.7995, -6.2370],
        [106.7990, -6.2390], [106.7985, -6.2410]
      ]
    },
    { id: 'road-asiaafrika', color: '#38A169', coords: [
        [106.8019, -6.2270], [106.8040, -6.2265], [106.8060, -6.2260],
        [106.8080, -6.2255], [106.8100, -6.2250], [106.8120, -6.2250],
        [106.8140, -6.2250]
      ]
    },
    { id: 'road-satrio', color: '#E53E3E', coords: [
        [106.8210, -6.2148], [106.8230, -6.2130], [106.8250, -6.2110],
        [106.8260, -6.2090], [106.8260, -6.2070], [106.8260, -6.2040]
      ]
    }
  ];

  function addRoadSegments() {
    if (roadLayer) map.removeLayer(roadLayer);

    roadLayer = L.layerGroup();

    roadSegments.forEach(segment => {
      // Leaflet coords = [lat, lng] bukan [lng, lat]
      const latLngs = segment.coords.map(c => [c[1], c[0]] as [number, number]);

      // Glow
      L.polyline(latLngs, {
        color: segment.color,
        weight: 16,
        opacity: 0.3
      }).addTo(roadLayer);

      // Main line
      L.polyline(latLngs, {
        color: segment.color,
        weight: 6,
        opacity: 0.9,
        lineCap: 'round',
        lineJoin: 'round'
      }).addTo(roadLayer);
    });

    roadLayer.addTo(map);
  }

  function addMarkers() {
    if (markerLayer) map.removeLayer(markerLayer);
    markerLayer = L.layerGroup();

    dummyReports.forEach(report => {
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

  onDestroy(() => {
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

    <!-- Map Controls (Mobile) -->
    <div class="map-controls hide-desktop">
      <button class="map-btn" on:click={() => showSidebar = !showSidebar}>
        {#if showSidebar}
          <X size={20} />
        {:else}
          <List size={20} />
        {/if}
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
      <div class="stat-pill danger">
        <AlertTriangle size={14} />
        <span>{dummyReports.filter(r => r.severity === 'berat' || r.severity === 'korban').length} Parah</span>
      </div>
      <div class="stat-pill">
        <MapPin size={14} />
        <span>{dummyReports.length} Laporan</span>
      </div>
    </div>

    <!-- Legend -->
    <div class="map-legend">
      <div class="legend-title">Kondisi Jalan</div>
      <div class="legend-items">
        <div class="legend-item">
          <span class="legend-line" style="background: #38A169;"></span>
          <span>Aman</span>
        </div>
        <div class="legend-item">
          <span class="legend-line" style="background: #ECC94B;"></span>
          <span>Rusak Ringan</span>
        </div>
        <div class="legend-item">
          <span class="legend-line" style="background: #ED8936;"></span>
          <span>Rusak Sedang</span>
        </div>
        <div class="legend-item">
          <span class="legend-line" style="background: #E53E3E;"></span>
          <span>Rusak Berat</span>
        </div>
        <div class="legend-item">
          <span class="legend-line" style="background: #1A202C;"></span>
          <span>Ada Korban</span>
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
          <p>Tidak ada laporan dengan filter ini</p>
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
  :global(.district-label) {
    background: transparent !important;
    border: none !important;
    box-shadow: none !important;
    color: #90CDF4 !important;
    font-size: 11px !important;
    font-weight: 700 !important;
    text-transform: uppercase !important;
    letter-spacing: 0.06em !important;
    text-shadow: 0 0 4px rgba(26, 32, 44, 0.9), 0 0 8px rgba(26, 32, 44, 0.7) !important;
    white-space: nowrap !important;
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

  /* Map Controls */
  .map-controls {
    position: absolute;
    top: var(--space-md);
    left: var(--space-md);
    z-index: 10;
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
  }
  .map-btn:active {
    transform: scale(0.95);
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
    z-index: 10;
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
    z-index: 10;
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
    z-index: 10;
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
  .legend-line {
    width: 20px;
    height: 4px;
    border-radius: 2px;
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
      max-width: 130px;
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

