<script lang="ts">
  import { page } from '$app/stores';
  import { dummyReports, formatRupiah, getAgeLabel, getSeverityBadge, getStatusBadge } from '$lib/data/mock';
  import { ArrowLeft, Heart, Share2, MessageCircle, Eye, MapPin, Clock, Camera, AlertTriangle, Flag, ExternalLink, ChevronLeft, ChevronRight } from 'lucide-svelte';
  import type { Report } from '$lib/types';

  $: reportId = $page.params.id;
  $: report = dummyReports.find(r => r.id === reportId) as Report | undefined;

  $: severityBadge = report ? getSeverityBadge(report.severity) : { class: '', label: '' };
  $: statusBadge = report ? getStatusBadge(report.status) : { class: '', label: '' };
  $: ageInfo = report ? getAgeLabel(report.daysOld) : { label: '', icon: '' };

  let liked = false;
  let currentPhotoIdx = 0;
  let commentText = '';

  const dummyComments = [
    { name: 'Andi Prasetyo', time: '2 jam lalu', text: 'Ini udah parah banget, kemarin motor saya nyaris jatuh.', level: '👤 Warga' },
    { name: 'Rina Sari', time: '5 jam lalu', text: 'Mana nih pemerintahnya? Sudah sebulan lebih ga ada action.', level: '👁️ Warga Peduli' },
    { name: 'Joko Widodo Jr.', time: '1 hari lalu', text: 'Saya sudah lapor ke kelurahan tapi belum ada tindak lanjut 😤', level: '🏘️ Pak RT' },
  ];

  function toggleLike() {
    liked = !liked;
  }

  function shareReport() {
    if (navigator.share && report) {
      navigator.share({
        title: `JEDUG: ${report.title}`,
        text: `${report.title} - ${report.location}. Sudah ${report.daysOld} hari! Estimasi kerugian ${formatRupiah(report.estimatedLoss)}. #JEDUG`,
        url: window.location.href
      }).catch(() => {});
    } else {
      alert('Link disalin! (Demo Mode)');
    }
  }

  function nextPhoto() {
    if (report && currentPhotoIdx < report.photos.length - 1) {
      currentPhotoIdx++;
    }
  }
  function prevPhoto() {
    if (currentPhotoIdx > 0) {
      currentPhotoIdx--;
    }
  }
</script>

<svelte:head>
  <title>{report ? report.title : 'Laporan'} - JEDUG</title>
</svelte:head>

{#if report}
  <div class="detail-page">
    <div class="detail-container">
      <!-- Back button -->
      <a href="/" class="back-btn">
        <ArrowLeft size={20} />
        <span>Kembali</span>
      </a>

      <!-- Photo Gallery -->
      <div class="photo-gallery">
        <img src={report.photos[currentPhotoIdx]} alt="{report.title} foto {currentPhotoIdx + 1}" class="gallery-img" />
        
        {#if report.photos.length > 1}
          <div class="gallery-controls">
            <button class="gallery-btn" on:click={prevPhoto} disabled={currentPhotoIdx === 0}>
              <ChevronLeft size={20} />
            </button>
            <span class="gallery-counter">{currentPhotoIdx + 1} / {report.photos.length}</span>
            <button class="gallery-btn" on:click={nextPhoto} disabled={currentPhotoIdx === report.photos.length - 1}>
              <ChevronRight size={20} />
            </button>
          </div>
        {/if}

        <div class="gallery-badges">
          <span class="badge {severityBadge.class}">{severityBadge.label}</span>
          <span class="badge {statusBadge.class}">{statusBadge.label}</span>
        </div>

        {#if report.daysOld > 30}
          <div class="age-sticker-large">
            <span>{ageInfo.icon}</span>
            <span class="age-text">{report.daysOld > 90 ? 'FOSIL!' : 'ULTAH!'}</span>
          </div>
        {/if}
      </div>

      <!-- Content -->
      <div class="detail-content">
        <!-- Title & Actions -->
        <div class="detail-header">
          <h1>{report.title}</h1>
          <div class="header-actions">
            <button class="action-btn" class:liked on:click={toggleLike}>
              <Heart size={20} fill={liked ? 'currentColor' : 'none'} />
              <span>{report.reactions + (liked ? 1 : 0)}</span>
            </button>
            <button class="action-btn" on:click={shareReport}>
              <Share2 size={20} />
              <span>Share</span>
            </button>
            <button class="action-btn flag-btn">
              <Flag size={18} />
            </button>
          </div>
        </div>

        <!-- Location -->
        <div class="info-row">
          <MapPin size={16} />
          <div>
            <span class="info-main">{report.location}</span>
            <span class="info-sub">{report.kelurahan}, {report.kecamatan}</span>
          </div>
        </div>

        <!-- Age & Loss -->
        <div class="highlight-cards">
          <div class="highlight-card age">
            <span class="highlight-icon">{ageInfo.icon}</span>
            <div>
              <span class="highlight-value">{ageInfo.label}</span>
              <span class="highlight-label">Umur Kerusakan</span>
            </div>
          </div>
          <div class="highlight-card loss">
            <span class="highlight-icon">💸</span>
            <div>
              <span class="highlight-value">{formatRupiah(report.estimatedLoss)}</span>
              <span class="highlight-label">Estimasi Kerugian</span>
            </div>
          </div>
        </div>

        <!-- Stats -->
        <div class="stats-row">
          <div class="detail-stat">
            <Eye size={16} />
            <span>{report.views.toLocaleString('id-ID')} dilihat</span>
          </div>
          <div class="detail-stat">
            <Heart size={16} />
            <span>{report.reactions} reaksi</span>
          </div>
          <div class="detail-stat">
            <MessageCircle size={16} />
            <span>{report.comments} komentar</span>
          </div>
        </div>

        <!-- Description -->
        <div class="desc-section">
          <h3>Deskripsi</h3>
          <p>{report.description}</p>
        </div>

        <!-- Reporter -->
        <div class="reporter-section">
          <div class="reporter-avatar">
            <span>👤</span>
          </div>
          <div class="reporter-info">
            <span class="reporter-name">{report.reporterName}</span>
            <span class="reporter-level">{report.reporterLevel}</span>
          </div>
          <span class="reporter-date">
            <Clock size={13} />
            {new Date(report.createdAt).toLocaleDateString('id-ID', { 
              day: 'numeric', month: 'long', year: 'numeric' 
            })}
          </span>
        </div>

        <!-- Share Banner -->
        <div class="viral-banner">
          <div class="viral-content">
            <h3>🔥 Viralkan Laporan Ini!</h3>
            <p>Semakin banyak yang lihat, semakin besar tekanan untuk diperbaiki.</p>
          </div>
          <button class="viral-btn" on:click={shareReport}>
            <Share2 size={18} />
            Bagikan
          </button>
        </div>

        <!-- Emergency Button (for korban severity) -->
        {#if report.severity === 'korban'}
          <div class="emergency-banner">
            <AlertTriangle size={24} />
            <div>
              <strong>⚠️ Laporan Ada Korban</strong>
              <p>Laporan ini melibatkan korban kecelakaan. Prioritas penanganan tinggi.</p>
            </div>
          </div>
        {/if}

        <!-- Comments Section -->
        <div class="comments-section">
          <h3>💬 Komentar ({dummyComments.length})</h3>
          
          <div class="comment-input-wrapper">
            <input 
              type="text" 
              class="comment-input" 
              placeholder="Tulis komentar..."
              bind:value={commentText}
            />
            <button class="comment-send" disabled={!commentText.trim()}>
              Kirim
            </button>
          </div>

          <div class="comments-list">
            {#each dummyComments as comment}
              <div class="comment-item">
                <div class="comment-header">
                  <span class="comment-name">{comment.name}</span>
                  <span class="comment-level">{comment.level}</span>
                  <span class="comment-time">{comment.time}</span>
                </div>
                <p class="comment-text">{comment.text}</p>
              </div>
            {/each}
          </div>
        </div>
      </div>
    </div>
  </div>
{:else}
  <div class="not-found">
    <span class="nf-icon">🕳️</span>
    <h2>Laporan Tidak Ditemukan</h2>
    <p>Sepertinya lubangnya sudah ditutup... atau belum dilaporkan.</p>
    <a href="/" class="nf-link">Kembali ke Peta</a>
  </div>
{/if}

<style>
  .detail-page {
    min-height: calc(100dvh - var(--nav-height) - var(--bottom-nav-height));
    background: var(--bg-secondary);
  }
  .detail-container {
    max-width: 700px;
    margin: 0 auto;
  }

  /* Back */
  .back-btn {
    display: inline-flex;
    align-items: center;
    gap: var(--space-sm);
    padding: var(--space-md) var(--space-lg);
    color: var(--text-secondary);
    font-size: var(--text-sm);
    font-weight: var(--font-medium);
    text-decoration: none;
    transition: color var(--transition-fast);
  }
  .back-btn:hover {
    color: var(--text-primary);
    text-decoration: none;
  }

  /* Gallery */
  .photo-gallery {
    position: relative;
    width: 100%;
    aspect-ratio: 16/10;
    background: var(--bg-tertiary);
    overflow: hidden;
  }
  .gallery-img {
    width: 100%;
    height: 100%;
    object-fit: cover;
  }
  .gallery-controls {
    position: absolute;
    bottom: var(--space-md);
    left: 50%;
    transform: translateX(-50%);
    display: flex;
    align-items: center;
    gap: var(--space-sm);
    background: rgba(0,0,0,0.6);
    border-radius: var(--radius-full);
    padding: var(--space-xs) var(--space-sm);
  }
  .gallery-btn {
    width: 32px;
    height: 32px;
    border-radius: 50%;
    background: transparent;
    border: none;
    color: white;
    display: flex;
    align-items: center;
    justify-content: center;
    cursor: pointer;
  }
  .gallery-btn:disabled {
    opacity: 0.3;
    cursor: not-allowed;
  }
  .gallery-counter {
    color: white;
    font-size: var(--text-xs);
    font-weight: var(--font-medium);
  }
  .gallery-badges {
    position: absolute;
    top: var(--space-md);
    left: var(--space-md);
    display: flex;
    gap: var(--space-xs);
  }
  .age-sticker-large {
    position: absolute;
    top: var(--space-md);
    right: var(--space-md);
    display: flex;
    flex-direction: column;
    align-items: center;
    background: rgba(0,0,0,0.6);
    border-radius: var(--radius-lg);
    padding: var(--space-sm) var(--space-md);
    animation: bounce 2s infinite;
  }
  .age-sticker-large span:first-child {
    font-size: 2rem;
  }
  .age-text {
    color: white;
    font-size: var(--text-xs);
    font-weight: var(--font-bold);
  }

  @keyframes bounce {
    0%, 100% { transform: translateY(0); }
    50% { transform: translateY(-4px); }
  }

  /* Content */
  .detail-content {
    padding: var(--space-lg);
  }

  .detail-header {
    margin-bottom: var(--space-lg);
  }
  .detail-header h1 {
    font-size: var(--text-2xl);
    font-weight: var(--font-bold);
    line-height: var(--leading-tight);
    margin-bottom: var(--space-md);
  }
  .header-actions {
    display: flex;
    gap: var(--space-sm);
  }
  .action-btn {
    display: flex;
    align-items: center;
    gap: var(--space-xs);
    padding: 0.5rem 1rem;
    border-radius: var(--radius-full);
    background: var(--bg-tertiary);
    color: var(--text-secondary);
    font-size: var(--text-sm);
    font-weight: var(--font-medium);
    transition: all var(--transition-fast);
    cursor: pointer;
    border: none;
  }
  .action-btn:hover {
    background: var(--border-color-strong);
    color: var(--text-primary);
  }
  .action-btn.liked {
    background: var(--color-danger-light);
    color: var(--color-danger);
  }
  .flag-btn {
    padding: 0.5rem;
  }

  /* Info row */
  .info-row {
    display: flex;
    align-items: flex-start;
    gap: var(--space-sm);
    padding: var(--space-md);
    background: var(--bg-card);
    border-radius: var(--radius-lg);
    border: 1px solid var(--border-color);
    margin-bottom: var(--space-md);
    color: var(--text-secondary);
  }
  .info-main {
    display: block;
    font-size: var(--text-sm);
    font-weight: var(--font-medium);
    color: var(--text-primary);
  }
  .info-sub {
    font-size: var(--text-xs);
    color: var(--text-tertiary);
  }

  /* Highlight Cards */
  .highlight-cards {
    display: grid;
    grid-template-columns: 1fr 1fr;
    gap: var(--space-sm);
    margin-bottom: var(--space-lg);
  }
  .highlight-card {
    display: flex;
    align-items: center;
    gap: var(--space-sm);
    padding: var(--space-md);
    border-radius: var(--radius-lg);
    background: var(--bg-card);
    border: 1px solid var(--border-color);
  }
  .highlight-card.age {
    border-left: 3px solid var(--color-warning);
  }
  .highlight-card.loss {
    border-left: 3px solid var(--color-danger);
  }
  .highlight-icon {
    font-size: 1.5rem;
    flex-shrink: 0;
  }
  .highlight-value {
    display: block;
    font-size: var(--text-sm);
    font-weight: var(--font-bold);
    color: var(--text-primary);
  }
  .highlight-label {
    font-size: var(--text-xs);
    color: var(--text-tertiary);
  }

  /* Stats */
  .stats-row {
    display: flex;
    gap: var(--space-lg);
    padding: var(--space-md) 0;
    margin-bottom: var(--space-lg);
    border-bottom: 1px solid var(--border-color);
  }
  .detail-stat {
    display: flex;
    align-items: center;
    gap: var(--space-xs);
    font-size: var(--text-sm);
    color: var(--text-secondary);
  }

  /* Description */
  .desc-section {
    margin-bottom: var(--space-lg);
  }
  .desc-section h3 {
    font-size: var(--text-base);
    font-weight: var(--font-semibold);
    margin-bottom: var(--space-sm);
  }
  .desc-section p {
    font-size: var(--text-sm);
    color: var(--text-secondary);
    line-height: var(--leading-relaxed);
  }

  /* Reporter */
  .reporter-section {
    display: flex;
    align-items: center;
    gap: var(--space-md);
    padding: var(--space-md);
    background: var(--bg-card);
    border-radius: var(--radius-lg);
    border: 1px solid var(--border-color);
    margin-bottom: var(--space-lg);
  }
  .reporter-avatar {
    width: 40px;
    height: 40px;
    border-radius: var(--radius-full);
    background: var(--bg-tertiary);
    display: flex;
    align-items: center;
    justify-content: center;
    font-size: 1.25rem;
  }
  .reporter-info {
    flex: 1;
  }
  .reporter-name {
    display: block;
    font-size: var(--text-sm);
    font-weight: var(--font-semibold);
    color: var(--text-primary);
  }
  .reporter-level {
    font-size: var(--text-xs);
    color: var(--text-tertiary);
  }
  .reporter-date {
    display: flex;
    align-items: center;
    gap: 4px;
    font-size: var(--text-xs);
    color: var(--text-tertiary);
  }

  /* Viral Banner */
  .viral-banner {
    display: flex;
    align-items: center;
    justify-content: space-between;
    padding: var(--space-lg);
    background: linear-gradient(135deg, #FF6B35, #E53E3E);
    border-radius: var(--radius-xl);
    margin-bottom: var(--space-lg);
    gap: var(--space-md);
  }
  .viral-content h3 {
    color: white;
    font-size: var(--text-base);
    margin-bottom: 2px;
  }
  .viral-content p {
    color: rgba(255,255,255,0.8);
    font-size: var(--text-xs);
  }
  .viral-btn {
    display: flex;
    align-items: center;
    gap: var(--space-xs);
    padding: 0.625rem 1.25rem;
    background: white;
    color: var(--color-primary);
    border-radius: var(--radius-full);
    font-weight: var(--font-semibold);
    font-size: var(--text-sm);
    cursor: pointer;
    transition: all var(--transition-fast);
    flex-shrink: 0;
    border: none;
  }
  .viral-btn:hover {
    transform: scale(1.05);
  }

  /* Emergency */
  .emergency-banner {
    display: flex;
    align-items: flex-start;
    gap: var(--space-md);
    padding: var(--space-lg);
    background: var(--color-danger-light);
    border: 2px solid var(--color-danger);
    border-radius: var(--radius-xl);
    margin-bottom: var(--space-lg);
    color: var(--color-danger);
  }
  .emergency-banner p {
    font-size: var(--text-sm);
    margin-top: 4px;
    opacity: 0.8;
  }

  /* Comments */
  .comments-section {
    margin-bottom: var(--space-xl);
  }
  .comments-section h3 {
    font-size: var(--text-base);
    font-weight: var(--font-semibold);
    margin-bottom: var(--space-md);
  }
  .comment-input-wrapper {
    display: flex;
    gap: var(--space-sm);
    margin-bottom: var(--space-lg);
  }
  .comment-input {
    flex: 1;
    padding: 0.75rem 1rem;
    border: 1.5px solid var(--border-color);
    border-radius: var(--radius-full);
    background: var(--bg-card);
    color: var(--text-primary);
    font-size: var(--text-sm);
  }
  .comment-input:focus {
    border-color: var(--color-primary);
    box-shadow: 0 0 0 3px var(--color-primary-light);
  }
  .comment-send {
    padding: 0.75rem 1.25rem;
    background: var(--color-primary);
    color: white;
    border-radius: var(--radius-full);
    font-size: var(--text-sm);
    font-weight: var(--font-semibold);
    cursor: pointer;
    transition: all var(--transition-fast);
    border: none;
  }
  .comment-send:hover {
    background: var(--color-primary-hover);
  }
  .comment-send:disabled {
    opacity: 0.5;
    cursor: not-allowed;
  }

  .comments-list {
    display: flex;
    flex-direction: column;
    gap: var(--space-md);
  }
  .comment-item {
    padding: var(--space-md);
    background: var(--bg-card);
    border-radius: var(--radius-lg);
    border: 1px solid var(--border-color);
  }
  .comment-header {
    display: flex;
    align-items: center;
    gap: var(--space-sm);
    margin-bottom: var(--space-xs);
    flex-wrap: wrap;
  }
  .comment-name {
    font-size: var(--text-sm);
    font-weight: var(--font-semibold);
    color: var(--text-primary);
  }
  .comment-level {
    font-size: var(--text-xs);
    color: var(--text-tertiary);
    background: var(--bg-tertiary);
    padding: 0.125rem 0.375rem;
    border-radius: var(--radius-sm);
  }
  .comment-time {
    font-size: var(--text-xs);
    color: var(--text-tertiary);
    margin-left: auto;
  }
  .comment-text {
    font-size: var(--text-sm);
    color: var(--text-secondary);
    line-height: var(--leading-relaxed);
  }

  /* Not Found */
  .not-found {
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    min-height: calc(100dvh - var(--nav-height) - var(--bottom-nav-height));
    padding: var(--space-xl);
    text-align: center;
  }
  .nf-icon {
    font-size: 4rem;
    margin-bottom: var(--space-lg);
  }
  .not-found h2 {
    font-size: var(--text-2xl);
    margin-bottom: var(--space-sm);
  }
  .not-found p {
    color: var(--text-secondary);
    font-size: var(--text-sm);
    margin-bottom: var(--space-xl);
  }
  .nf-link {
    padding: 0.75rem 1.5rem;
    background: var(--color-primary);
    color: white;
    border-radius: var(--radius-xl);
    font-weight: var(--font-semibold);
    text-decoration: none;
    transition: all var(--transition-fast);
  }
  .nf-link:hover {
    background: var(--color-primary-hover);
    text-decoration: none;
  }

  @media (min-width: 769px) {
    .detail-content {
      padding: var(--space-xl) var(--space-lg);
    }
    .photo-gallery {
      aspect-ratio: 2/1;
      border-radius: 0 0 var(--radius-xl) var(--radius-xl);
    }
  }
</style>
