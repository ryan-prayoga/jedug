<script lang="ts">
  import { dummyUserProfile, dummyReports } from '$lib/data/mock';
  import ReportCard from '$lib/components/ReportCard.svelte';
  import { Settings, LogIn, Share2, Award, MapPin, CheckCircle, Star, ChevronRight, Bell, Shield, HelpCircle, LogOut, Edit3 } from 'lucide-svelte';
  import { LEVEL_PROGRESSION } from '$lib/types';

  const profile = dummyUserProfile;
  const userReports = dummyReports.filter(r => r.reporterName === profile.name);

  let activeTab: 'laporan' | 'badge' = 'laporan';

  // Calculate level progress
  $: currentLevelIdx = LEVEL_PROGRESSION.findIndex(l => l.name === profile.level);
  $: nextLevel = currentLevelIdx < LEVEL_PROGRESSION.length - 1 ? LEVEL_PROGRESSION[currentLevelIdx + 1] : null;
  $: progressPercent = nextLevel
    ? Math.round(((profile.points - LEVEL_PROGRESSION[currentLevelIdx].minPoints) / (nextLevel.minPoints - LEVEL_PROGRESSION[currentLevelIdx].minPoints)) * 100)
    : 100;
</script>

<svelte:head>
  <title>JEDUG - Profil Saya</title>
</svelte:head>

<div class="profile-page">
  <div class="profile-container">
    <!-- Profile Header Card -->
    <div class="profile-card">
      <div class="profile-top">
        <div class="avatar">
          <span class="avatar-emoji">{profile.levelIcon}</span>
        </div>
        <div class="profile-info">
          <h1 class="profile-name">{profile.name}</h1>
          <div class="profile-level">
            <span class="level-badge">{profile.levelIcon} {profile.level}</span>
            <span class="level-points">⭐ {profile.points} poin</span>
          </div>
        </div>
        <button class="edit-btn">
          <Edit3 size={18} />
        </button>
      </div>

      <!-- Level Progress -->
      {#if nextLevel}
        <div class="level-progress">
          <div class="progress-header">
            <span class="progress-label">Menuju {nextLevel.icon} {nextLevel.name}</span>
            <span class="progress-value">{profile.points}/{nextLevel.minPoints}</span>
          </div>
          <div class="progress-track">
            <div class="progress-fill" style="width: {progressPercent}%"></div>
          </div>
        </div>
      {/if}

      <!-- Stats Row -->
      <div class="profile-stats">
        <div class="p-stat">
          <span class="p-stat-value">{profile.totalReports}</span>
          <span class="p-stat-label">Laporan</span>
        </div>
        <div class="p-stat-divider"></div>
        <div class="p-stat">
          <span class="p-stat-value">{profile.totalVerified}</span>
          <span class="p-stat-label">Terverifikasi</span>
        </div>
        <div class="p-stat-divider"></div>
        <div class="p-stat">
          <span class="p-stat-value">{profile.totalFixed}</span>
          <span class="p-stat-label">Diperbaiki</span>
        </div>
      </div>
    </div>

    <!-- Tabs -->
    <div class="tabs">
      <button class="tab" class:active={activeTab === 'laporan'} on:click={() => activeTab = 'laporan'}>
        <MapPin size={16} />
        Laporan Saya
      </button>
      <button class="tab" class:active={activeTab === 'badge'} on:click={() => activeTab = 'badge'}>
        <Award size={16} />
        Lencana
      </button>
    </div>

    <!-- Tab Content -->
    {#if activeTab === 'laporan'}
      <div class="report-list">
        {#if userReports.length > 0}
          {#each userReports as report (report.id)}
            <ReportCard {report} compact />
          {/each}
        {:else}
          <div class="empty-state">
            <span class="empty-icon">📝</span>
            <h3>Belum ada laporan</h3>
            <p>Mulai laporkan jalan rusak di sekitarmu!</p>
            <a href="/lapor" class="empty-cta">
              Buat Laporan Pertama
            </a>
          </div>
        {/if}
      </div>
    {:else}
      <div class="badge-grid">
        {#each profile.badges as badge (badge.id)}
          <div class="badge-card" class:earned={badge.earned} class:locked={!badge.earned}>
            <span class="badge-icon">{badge.icon}</span>
            <h4 class="badge-name">{badge.name}</h4>
            <p class="badge-desc">{badge.description}</p>
            {#if !badge.earned}
              <span class="badge-lock">🔒</span>
            {/if}
          </div>
        {/each}
      </div>
    {/if}

    <!-- Menu Section -->
    <div class="menu-section">
      <h3 class="menu-title">Pengaturan</h3>
      <div class="menu-list">
        <button class="menu-item">
          <Bell size={18} />
          <span>Notifikasi</span>
          <ChevronRight size={16} />
        </button>
        <button class="menu-item">
          <Shield size={18} />
          <span>Privasi & Keamanan</span>
          <ChevronRight size={16} />
        </button>
        <button class="menu-item">
          <Share2 size={18} />
          <span>Bagikan Aplikasi</span>
          <ChevronRight size={16} />
        </button>
        <button class="menu-item">
          <HelpCircle size={18} />
          <span>Bantuan & FAQ</span>
          <ChevronRight size={16} />
        </button>
        <button class="menu-item">
          <LogIn size={18} />
          <span>Login dengan Google</span>
          <ChevronRight size={16} />
        </button>
      </div>
    </div>

    <!-- App version -->
    <div class="app-version">
      <p>JEDUG v1.0.0 (Demo)</p>
      <p>Lapor, Pantau, Viralkan. 🕳️</p>
    </div>
  </div>
</div>

<style>
  .profile-page {
    min-height: calc(100dvh - var(--nav-height) - var(--bottom-nav-height));
    background: var(--bg-secondary);
    padding: var(--space-md);
  }
  .profile-container {
    max-width: 600px;
    margin: 0 auto;
  }

  /* Profile Card */
  .profile-card {
    background: var(--bg-card);
    border-radius: var(--radius-xl);
    padding: var(--space-xl);
    border: 1px solid var(--border-color);
    margin-bottom: var(--space-lg);
  }
  .profile-top {
    display: flex;
    align-items: center;
    gap: var(--space-md);
    margin-bottom: var(--space-lg);
  }
  .avatar {
    width: 64px;
    height: 64px;
    border-radius: var(--radius-full);
    background: var(--color-primary-light);
    display: flex;
    align-items: center;
    justify-content: center;
    flex-shrink: 0;
  }
  .avatar-emoji {
    font-size: 2rem;
  }
  .profile-info {
    flex: 1;
    min-width: 0;
  }
  .profile-name {
    font-size: var(--text-xl);
    font-weight: var(--font-bold);
    margin-bottom: 4px;
  }
  .profile-level {
    display: flex;
    align-items: center;
    gap: var(--space-sm);
    flex-wrap: wrap;
  }
  .level-badge {
    font-size: var(--text-sm);
    font-weight: var(--font-semibold);
    color: var(--color-primary);
    background: var(--color-primary-light);
    padding: 0.125rem 0.5rem;
    border-radius: var(--radius-full);
  }
  .level-points {
    font-size: var(--text-sm);
    color: var(--text-tertiary);
  }
  .edit-btn {
    width: 40px;
    height: 40px;
    border-radius: var(--radius-full);
    background: var(--bg-tertiary);
    color: var(--text-secondary);
    display: flex;
    align-items: center;
    justify-content: center;
    cursor: pointer;
    transition: all var(--transition-fast);
    border: none;
  }
  .edit-btn:hover {
    background: var(--border-color-strong);
    color: var(--text-primary);
  }

  /* Level Progress */
  .level-progress {
    margin-bottom: var(--space-lg);
  }
  .progress-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: var(--space-xs);
  }
  .progress-label {
    font-size: var(--text-xs);
    color: var(--text-secondary);
    font-weight: var(--font-medium);
  }
  .progress-value {
    font-size: var(--text-xs);
    color: var(--text-tertiary);
  }
  .progress-track {
    width: 100%;
    height: 8px;
    background: var(--bg-tertiary);
    border-radius: var(--radius-full);
    overflow: hidden;
  }
  .progress-fill {
    height: 100%;
    background: linear-gradient(90deg, var(--color-primary), var(--color-secondary));
    border-radius: var(--radius-full);
    transition: width var(--transition-slow);
  }

  /* Profile Stats */
  .profile-stats {
    display: flex;
    align-items: center;
    justify-content: center;
    gap: var(--space-lg);
    padding-top: var(--space-lg);
    border-top: 1px solid var(--border-color);
  }
  .p-stat {
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: 2px;
    flex: 1;
  }
  .p-stat-value {
    font-size: var(--text-2xl);
    font-weight: var(--font-bold);
    color: var(--text-primary);
  }
  .p-stat-label {
    font-size: var(--text-xs);
    color: var(--text-tertiary);
  }
  .p-stat-divider {
    width: 1px;
    height: 40px;
    background: var(--border-color);
  }

  /* Tabs */
  .tabs {
    display: flex;
    gap: var(--space-xs);
    margin-bottom: var(--space-lg);
    background: var(--bg-card);
    border-radius: var(--radius-xl);
    padding: var(--space-xs);
    border: 1px solid var(--border-color);
  }
  .tab {
    flex: 1;
    display: flex;
    align-items: center;
    justify-content: center;
    gap: var(--space-sm);
    padding: 0.625rem var(--space-md);
    border-radius: var(--radius-lg);
    font-size: var(--text-sm);
    font-weight: var(--font-medium);
    color: var(--text-secondary);
    background: transparent;
    transition: all var(--transition-fast);
    cursor: pointer;
    border: none;
  }
  .tab:hover {
    color: var(--text-primary);
    background: var(--bg-tertiary);
  }
  .tab.active {
    background: var(--color-primary);
    color: white;
    font-weight: var(--font-semibold);
  }

  /* Report List */
  .report-list {
    display: flex;
    flex-direction: column;
    gap: var(--space-md);
    margin-bottom: var(--space-xl);
  }

  .empty-state {
    display: flex;
    flex-direction: column;
    align-items: center;
    padding: var(--space-3xl) var(--space-xl);
    text-align: center;
    color: var(--text-secondary);
  }
  .empty-icon {
    font-size: 3rem;
    margin-bottom: var(--space-md);
  }
  .empty-state h3 {
    font-size: var(--text-lg);
    font-weight: var(--font-semibold);
    margin-bottom: var(--space-xs);
    color: var(--text-primary);
  }
  .empty-state p {
    font-size: var(--text-sm);
    margin-bottom: var(--space-lg);
  }
  .empty-cta {
    padding: 0.75rem 1.5rem;
    background: var(--color-primary);
    color: white;
    border-radius: var(--radius-xl);
    font-weight: var(--font-semibold);
    font-size: var(--text-sm);
    text-decoration: none;
    transition: all var(--transition-fast);
  }
  .empty-cta:hover {
    background: var(--color-primary-hover);
    text-decoration: none;
  }

  /* Badge Grid */
  .badge-grid {
    display: grid;
    grid-template-columns: repeat(auto-fill, minmax(150px, 1fr));
    gap: var(--space-md);
    margin-bottom: var(--space-xl);
  }
  .badge-card {
    background: var(--bg-card);
    border: 1px solid var(--border-color);
    border-radius: var(--radius-xl);
    padding: var(--space-lg);
    text-align: center;
    position: relative;
    transition: all var(--transition-fast);
  }
  .badge-card.earned {
    border-color: var(--color-primary);
    background: var(--color-primary-light);
  }
  .badge-card.locked {
    opacity: 0.6;
  }
  .badge-card:hover {
    transform: translateY(-2px);
    box-shadow: var(--shadow-md);
  }
  .badge-icon {
    font-size: 2.5rem;
    display: block;
    margin-bottom: var(--space-sm);
  }
  .badge-name {
    font-size: var(--text-sm);
    font-weight: var(--font-semibold);
    color: var(--text-primary);
    margin-bottom: 4px;
  }
  .badge-desc {
    font-size: var(--text-xs);
    color: var(--text-tertiary);
    line-height: var(--leading-relaxed);
  }
  .badge-lock {
    position: absolute;
    top: var(--space-sm);
    right: var(--space-sm);
    font-size: var(--text-sm);
  }

  /* Menu Section */
  .menu-section {
    margin-bottom: var(--space-xl);
  }
  .menu-title {
    font-size: var(--text-sm);
    font-weight: var(--font-semibold);
    color: var(--text-tertiary);
    text-transform: uppercase;
    letter-spacing: 0.5px;
    margin-bottom: var(--space-sm);
    padding-left: var(--space-md);
  }
  .menu-list {
    background: var(--bg-card);
    border-radius: var(--radius-xl);
    border: 1px solid var(--border-color);
    overflow: hidden;
  }
  .menu-item {
    display: flex;
    align-items: center;
    gap: var(--space-md);
    padding: var(--space-md) var(--space-lg);
    width: 100%;
    color: var(--text-primary);
    font-size: var(--text-sm);
    transition: background var(--transition-fast);
    cursor: pointer;
    border: none;
    background: none;
    border-bottom: 1px solid var(--border-color);
    text-align: left;
  }
  .menu-item:last-child {
    border-bottom: none;
  }
  .menu-item:hover {
    background: var(--bg-tertiary);
  }
  .menu-item span {
    flex: 1;
  }
  .menu-item :global(svg:last-child) {
    color: var(--text-tertiary);
  }

  /* App Version */
  .app-version {
    text-align: center;
    padding: var(--space-xl) 0 var(--space-3xl);
    color: var(--text-tertiary);
    font-size: var(--text-xs);
    line-height: var(--leading-relaxed);
  }

  @media (min-width: 769px) {
    .profile-page {
      padding: var(--space-xl);
    }
    .badge-grid {
      grid-template-columns: repeat(3, 1fr);
    }
  }
</style>
