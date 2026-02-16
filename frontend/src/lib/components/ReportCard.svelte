<script lang="ts">
  import type { Report } from '$lib/types';
  import { formatRupiah, getAgeLabel, getSeverityBadge, getStatusBadge } from '$lib/data/mock';
  import { Heart, MessageCircle, Eye, MapPin, Clock, ChevronRight } from 'lucide-svelte';

  export let report: Report;
  export let compact: boolean = false;

  $: severityBadge = getSeverityBadge(report.severity);
  $: statusBadge = getStatusBadge(report.status);
  $: ageInfo = getAgeLabel(report.daysOld);
</script>

<a href="/laporan/{report.id}" class="report-card" class:compact>
  <!-- Photo -->
  <div class="card-photo">
    <img src={report.photos[0]} alt={report.title} loading="lazy" />
    <div class="card-badges">
      <span class="badge {severityBadge.class}">{severityBadge.label}</span>
    </div>
    {#if report.daysOld > 30}
      <div class="age-sticker">
        <span>{ageInfo.icon}</span>
      </div>
    {/if}
  </div>

  <!-- Content -->
  <div class="card-content">
    <div class="card-header">
      <h3 class="card-title">{report.title}</h3>
      <span class="badge {statusBadge.class} status-badge">{statusBadge.label}</span>
    </div>

    <div class="card-location">
      <MapPin size={14} />
      <span>{report.location}</span>
    </div>

    {#if !compact}
      <p class="card-desc">{report.description}</p>
    {/if}

    <div class="card-meta">
      <div class="meta-left">
        <span class="meta-item">
          <Clock size={13} />
          {ageInfo.icon} {ageInfo.label}
        </span>
        {#if report.estimatedLoss > 0}
          <span class="meta-item loss">
            💸 {formatRupiah(report.estimatedLoss)}
          </span>
        {/if}
      </div>
    </div>

    <div class="card-stats">
      <span class="stat">
        <Heart size={14} />
        {report.reactions}
      </span>
      <span class="stat">
        <MessageCircle size={14} />
        {report.comments}
      </span>
      <span class="stat">
        <Eye size={14} />
        {report.views}
      </span>
      <span class="card-arrow">
        <ChevronRight size={16} />
      </span>
    </div>
  </div>
</a>

<style>
  .report-card {
    display: flex;
    flex-direction: column;
    background: var(--bg-card);
    border-radius: var(--radius-xl);
    overflow: hidden;
    border: 1px solid var(--border-color);
    transition: all var(--transition-fast);
    text-decoration: none;
    color: inherit;
    cursor: pointer;
  }
  .report-card:hover {
    box-shadow: var(--shadow-lg);
    border-color: var(--border-color-strong);
    text-decoration: none;
    transform: translateY(-2px);
  }

  .compact {
    flex-direction: row;
    border-radius: var(--radius-lg);
  }
  .compact .card-photo {
    width: 120px;
    min-height: 100px;
    border-radius: var(--radius-lg) 0 0 var(--radius-lg);
  }
  .compact .card-photo img {
    height: 100%;
    object-fit: cover;
  }

  .card-photo {
    position: relative;
    width: 100%;
    height: 180px;
    overflow: hidden;
  }
  .card-photo img {
    width: 100%;
    height: 100%;
    object-fit: cover;
    transition: transform var(--transition-slow);
  }
  .report-card:hover .card-photo img {
    transform: scale(1.05);
  }

  .card-badges {
    position: absolute;
    top: var(--space-sm);
    left: var(--space-sm);
    display: flex;
    gap: var(--space-xs);
  }

  .age-sticker {
    position: absolute;
    top: var(--space-sm);
    right: var(--space-sm);
    font-size: 1.5rem;
    filter: drop-shadow(0 2px 4px rgba(0,0,0,0.3));
    animation: bounce 2s infinite;
  }

  @keyframes bounce {
    0%, 100% { transform: translateY(0); }
    50% { transform: translateY(-4px); }
  }

  .card-content {
    padding: var(--space-md);
    display: flex;
    flex-direction: column;
    gap: var(--space-sm);
    flex: 1;
  }

  .card-header {
    display: flex;
    align-items: flex-start;
    justify-content: space-between;
    gap: var(--space-sm);
  }

  .card-title {
    font-size: var(--text-base);
    font-weight: var(--font-semibold);
    color: var(--text-primary);
    line-height: var(--leading-tight);
    display: -webkit-box;
    -webkit-line-clamp: 2;
    -webkit-box-orient: vertical;
    overflow: hidden;
  }
  .compact .card-title {
    font-size: var(--text-sm);
    -webkit-line-clamp: 1;
  }

  .status-badge {
    flex-shrink: 0;
    font-size: 0.65rem;
  }

  .card-location {
    display: flex;
    align-items: center;
    gap: var(--space-xs);
    color: var(--text-secondary);
    font-size: var(--text-sm);
  }

  .card-desc {
    font-size: var(--text-sm);
    color: var(--text-secondary);
    line-height: var(--leading-relaxed);
    display: -webkit-box;
    -webkit-line-clamp: 2;
    -webkit-box-orient: vertical;
    overflow: hidden;
  }

  .card-meta {
    display: flex;
    align-items: center;
    justify-content: space-between;
    flex-wrap: wrap;
    gap: var(--space-xs);
  }
  .meta-left {
    display: flex;
    align-items: center;
    gap: var(--space-md);
    flex-wrap: wrap;
  }
  .meta-item {
    display: flex;
    align-items: center;
    gap: 4px;
    font-size: var(--text-xs);
    color: var(--text-tertiary);
  }
  .meta-item.loss {
    color: var(--color-primary);
    font-weight: var(--font-semibold);
  }

  .card-stats {
    display: flex;
    align-items: center;
    gap: var(--space-md);
    padding-top: var(--space-sm);
    border-top: 1px solid var(--border-color);
  }
  .stat {
    display: flex;
    align-items: center;
    gap: 4px;
    font-size: var(--text-xs);
    color: var(--text-tertiary);
    font-weight: var(--font-medium);
  }
  .card-arrow {
    margin-left: auto;
    color: var(--text-tertiary);
  }
</style>
