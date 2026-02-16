<script lang="ts">
  import { dummyKecamatanRanks, dummyReports, formatRupiah } from '$lib/data/mock';
  import { Trophy, TrendingUp, TrendingDown, Minus, AlertTriangle, MapPin, DollarSign, BarChart3 } from 'lucide-svelte';

  let activeTab: 'terparah' | 'kerugian' = 'terparah';

  $: totalReports = dummyReports.length;
  $: totalLoss = dummyKecamatanRanks.reduce((sum, k) => sum + k.totalLoss, 0);
  $: avgFixed = Math.round(dummyKecamatanRanks.reduce((sum, k) => sum + k.percentFixed, 0) / dummyKecamatanRanks.length);

  $: sortedRanks = activeTab === 'terparah'
    ? [...dummyKecamatanRanks].sort((a, b) => b.totalReports - a.totalReports)
    : [...dummyKecamatanRanks].sort((a, b) => b.totalLoss - a.totalLoss);

  function getMedalEmoji(rank: number): string {
    if (rank === 1) return '🥇';
    if (rank === 2) return '🥈';
    if (rank === 3) return '🥉';
    return `#${rank}`;
  }

  function getSeverityIcon(severity: string): string {
    switch (severity) {
      case 'korban': return '💀';
      case 'berat': return '🔴';
      case 'sedang': return '🟠';
      case 'ringan': return '🟡';
      default: return '⚪';
    }
  }
</script>

<svelte:head>
  <title>JEDUG - Ranking Kecamatan Terparah</title>
</svelte:head>

<div class="ranking-page">
  <div class="ranking-container">
    <!-- Header -->
    <div class="ranking-header">
      <div class="header-text">
        <h1>🏆 Klasemen Jalan Rusak</h1>
        <p>Kecamatan mana yang paling butuh perhatian?</p>
      </div>
    </div>

    <!-- Stats Summary -->
    <div class="stats-grid">
      <div class="stat-card">
        <div class="stat-icon">
          <MapPin size={20} />
        </div>
        <div class="stat-info">
          <span class="stat-value">{totalReports}</span>
          <span class="stat-label">Total Laporan</span>
        </div>
      </div>
      <div class="stat-card danger">
        <div class="stat-icon">
          <DollarSign size={20} />
        </div>
        <div class="stat-info">
          <span class="stat-value">{formatRupiah(totalLoss)}</span>
          <span class="stat-label">Est. Kerugian</span>
        </div>
      </div>
      <div class="stat-card success">
        <div class="stat-icon">
          <BarChart3 size={20} />
        </div>
        <div class="stat-info">
          <span class="stat-value">{avgFixed}%</span>
          <span class="stat-label">Rata-rata Diperbaiki</span>
        </div>
      </div>
    </div>

    <!-- Tabs -->
    <div class="tabs">
      <button class="tab" class:active={activeTab === 'terparah'} on:click={() => activeTab = 'terparah'}>
        <AlertTriangle size={16} />
        Terparah
      </button>
      <button class="tab" class:active={activeTab === 'kerugian'} on:click={() => activeTab = 'kerugian'}>
        <DollarSign size={16} />
        Kerugian Terbesar
      </button>
    </div>

    <!-- Leaderboard -->
    <div class="leaderboard">
      {#each sortedRanks as kec, i (kec.name)}
        <div class="rank-card" class:top3={i < 3}>
          <div class="rank-position">
            <span class="rank-medal">{getMedalEmoji(i + 1)}</span>
          </div>

          <div class="rank-info">
            <div class="rank-header-row">
              <h3 class="rank-name">{kec.name}</h3>
              <span class="severity-indicator">{getSeverityIcon(kec.topSeverity)}</span>
            </div>
            <div class="rank-stats-row">
              <span class="rank-stat">
                <MapPin size={12} />
                {kec.totalReports} laporan
              </span>
              <span class="rank-stat loss">
                💸 {formatRupiah(kec.totalLoss)}
              </span>
            </div>
          </div>

          <div class="rank-right">
            <div class="fixed-bar-container">
              <div class="fixed-bar" style="width: {kec.percentFixed}%"></div>
              <span class="fixed-label">{kec.percentFixed}%</span>
            </div>
            <div class="trend-indicator" class:up={kec.trend === 'up'} class:down={kec.trend === 'down'}>
              {#if kec.trend === 'up'}
                <TrendingUp size={14} />
              {:else if kec.trend === 'down'}
                <TrendingDown size={14} />
              {:else}
                <Minus size={14} />
              {/if}
            </div>
          </div>
        </div>
      {/each}
    </div>

    <!-- Sarcasm Footer -->
    <div class="sarcasm-footer">
      <div class="sarcasm-card">
        <span class="sarcasm-icon">🎂</span>
        <div>
          <strong>3 lubang</strong> sudah berulang tahun bulan ini
          <span class="sarcasm-sub">Selamat! Sudah 30+ hari tidak diperbaiki 🎉</span>
        </div>
      </div>
      <div class="sarcasm-card">
        <span class="sarcasm-icon">🦖</span>
        <div>
          <strong>1 lubang</strong> menjadi fosil
          <span class="sarcasm-sub">Sudah 3+ bulan. Mungkin perlu arkeolog? 🤔</span>
        </div>
      </div>
    </div>
  </div>
</div>

<style>
  .ranking-page {
    min-height: calc(100dvh - var(--nav-height) - var(--bottom-nav-height));
    background: var(--bg-secondary);
    padding: var(--space-md);
  }

  .ranking-container {
    max-width: 800px;
    margin: 0 auto;
  }

  /* Header */
  .ranking-header {
    padding: var(--space-lg) 0;
  }
  .header-text h1 {
    font-size: var(--text-2xl);
    font-weight: var(--font-extrabold);
    margin-bottom: var(--space-xs);
  }
  .header-text p {
    color: var(--text-secondary);
    font-size: var(--text-sm);
  }

  /* Stats Grid */
  .stats-grid {
    display: grid;
    grid-template-columns: repeat(3, 1fr);
    gap: var(--space-sm);
    margin-bottom: var(--space-lg);
  }
  .stat-card {
    background: var(--bg-card);
    border-radius: var(--radius-xl);
    padding: var(--space-md);
    border: 1px solid var(--border-color);
    display: flex;
    align-items: center;
    gap: var(--space-sm);
  }
  .stat-icon {
    width: 40px;
    height: 40px;
    border-radius: var(--radius-lg);
    background: var(--color-info-light);
    color: var(--color-info);
    display: flex;
    align-items: center;
    justify-content: center;
    flex-shrink: 0;
  }
  .stat-card.danger .stat-icon {
    background: var(--color-danger-light);
    color: var(--color-danger);
  }
  .stat-card.success .stat-icon {
    background: var(--color-success-light);
    color: var(--color-success);
  }
  .stat-info {
    display: flex;
    flex-direction: column;
    min-width: 0;
  }
  .stat-value {
    font-size: var(--text-lg);
    font-weight: var(--font-bold);
    color: var(--text-primary);
    line-height: 1.2;
  }
  .stat-label {
    font-size: var(--text-xs);
    color: var(--text-tertiary);
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

  /* Leaderboard */
  .leaderboard {
    display: flex;
    flex-direction: column;
    gap: var(--space-sm);
    margin-bottom: var(--space-xl);
  }
  .rank-card {
    display: flex;
    align-items: center;
    gap: var(--space-md);
    padding: var(--space-md) var(--space-lg);
    background: var(--bg-card);
    border-radius: var(--radius-xl);
    border: 1px solid var(--border-color);
    transition: all var(--transition-fast);
  }
  .rank-card:hover {
    box-shadow: var(--shadow-md);
    transform: translateX(4px);
  }
  .rank-card.top3 {
    border-left: 3px solid var(--color-primary);
  }

  .rank-position {
    min-width: 40px;
    text-align: center;
  }
  .rank-medal {
    font-size: var(--text-xl);
    font-weight: var(--font-bold);
    color: var(--text-tertiary);
  }
  .top3 .rank-medal {
    font-size: 1.5rem;
  }

  .rank-info {
    flex: 1;
    min-width: 0;
  }
  .rank-header-row {
    display: flex;
    align-items: center;
    gap: var(--space-sm);
    margin-bottom: 2px;
  }
  .rank-name {
    font-size: var(--text-base);
    font-weight: var(--font-semibold);
    color: var(--text-primary);
  }
  .severity-indicator {
    font-size: var(--text-sm);
  }
  .rank-stats-row {
    display: flex;
    align-items: center;
    gap: var(--space-md);
  }
  .rank-stat {
    display: flex;
    align-items: center;
    gap: 3px;
    font-size: var(--text-xs);
    color: var(--text-tertiary);
  }
  .rank-stat.loss {
    color: var(--color-primary);
    font-weight: var(--font-medium);
  }

  .rank-right {
    display: flex;
    align-items: center;
    gap: var(--space-md);
    flex-shrink: 0;
  }

  .fixed-bar-container {
    width: 80px;
    height: 8px;
    background: var(--bg-tertiary);
    border-radius: var(--radius-full);
    position: relative;
    overflow: hidden;
  }
  .fixed-bar {
    height: 100%;
    background: var(--color-success);
    border-radius: var(--radius-full);
    transition: width var(--transition-slow);
  }
  .fixed-label {
    position: absolute;
    right: -35px;
    top: 50%;
    transform: translateY(-50%);
    font-size: var(--text-xs);
    color: var(--text-tertiary);
    font-weight: var(--font-medium);
    white-space: nowrap;
  }

  .trend-indicator {
    width: 28px;
    height: 28px;
    border-radius: var(--radius-full);
    display: flex;
    align-items: center;
    justify-content: center;
    background: var(--bg-tertiary);
    color: var(--text-tertiary);
  }
  .trend-indicator.up {
    background: var(--color-danger-light);
    color: var(--color-danger);
  }
  .trend-indicator.down {
    background: var(--color-success-light);
    color: var(--color-success);
  }

  /* Sarcasm Footer */
  .sarcasm-footer {
    display: flex;
    flex-direction: column;
    gap: var(--space-sm);
    padding-bottom: var(--space-xl);
  }
  .sarcasm-card {
    display: flex;
    align-items: flex-start;
    gap: var(--space-md);
    padding: var(--space-lg);
    background: var(--bg-card);
    border-radius: var(--radius-xl);
    border: 1px solid var(--border-color);
    font-size: var(--text-sm);
    color: var(--text-primary);
    line-height: var(--leading-relaxed);
  }
  .sarcasm-icon {
    font-size: 2rem;
    flex-shrink: 0;
  }
  .sarcasm-sub {
    display: block;
    font-size: var(--text-xs);
    color: var(--text-tertiary);
    margin-top: 2px;
  }

  @media (max-width: 768px) {
    .stats-grid {
      grid-template-columns: 1fr;
    }
    .stat-card {
      padding: var(--space-sm) var(--space-md);
    }
    .fixed-bar-container {
      width: 50px;
    }
    .rank-card {
      padding: var(--space-md);
    }
    .rank-right {
      gap: var(--space-sm);
    }
  }

  @media (min-width: 769px) {
    .ranking-page {
      padding: var(--space-xl);
    }
  }
</style>
