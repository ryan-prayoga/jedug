<script lang="ts">
  import { onMount } from 'svelte';
  import { goto } from '$app/navigation';
  import { Camera, MapPin, AlertTriangle, Send, X, ChevronDown, Info, Crosshair, Loader2, CheckCircle } from 'lucide-svelte';
  import { compressImage, blobToDataURL, formatFileSize } from '$lib/utils/imageCompression';
  import { getDeviceFingerprint } from '$lib/utils/fingerprint';
  import { api } from '$lib/api/client';
  import { severityToNumber } from '$lib/types';

  let photos: { preview: string; blob: Blob; originalSize: number; compressedSize: number }[] = [];
  let severity: string = '';
  let description: string = '';
  let agreedTos = false;
  let isLocating = false;
  let locationText = '';
  let step = 1; // 1=foto, 2=detail, 3=konfirmasi
  let isSubmitting = false;
  let submitSuccess = false;
  let submitError = '';
  let userLat = 0;
  let userLng = 0;
  let fingerprint = '';

  const severityOptions = [
    { value: 'ringan', label: 'Rusak Ringan', icon: '🟡', desc: 'Retakan kecil, masih bisa dilewati' },
    { value: 'sedang', label: 'Rusak Sedang', icon: '🟠', desc: 'Lubang terlihat, perlu hati-hati' },
    { value: 'berat', label: 'Rusak Berat', icon: '🔴', desc: 'Lubang besar, berbahaya' },
    { value: 'korban', label: 'Ada Korban', icon: '💀', desc: 'Sudah ada korban kecelakaan' }
  ];

  onMount(async () => {
    fingerprint = await getDeviceFingerprint();
  });

  async function handleFileSelect(e: Event) {
    const input = e.target as HTMLInputElement;
    if (!input.files) return;
    const files = Array.from(input.files);
    
    for (const file of files) {
      if (photos.length >= 3) break;
      
      try {
        const originalSize = file.size;
        const compressed = await compressImage(file);
        const preview = await blobToDataURL(compressed);
        
        photos = [...photos, {
          preview,
          blob: compressed,
          originalSize,
          compressedSize: compressed.size
        }];
      } catch (err) {
        console.error('Failed to compress image:', err);
      }
    }
    
    // Reset input so same file can be selected again
    input.value = '';
  }

  function removePhoto(index: number) {
    photos = photos.filter((_, i) => i !== index);
  }

  function getLocation() {
    if (!navigator.geolocation) {
      locationText = 'Geolokasi tidak didukung di perangkat ini';
      return;
    }

    isLocating = true;
    navigator.geolocation.getCurrentPosition(
      async (position) => {
        userLat = position.coords.latitude;
        userLng = position.coords.longitude;
        
        // Try reverse geocode via Nominatim (free)
        try {
          const res = await fetch(
            `https://nominatim.openstreetmap.org/reverse?lat=${userLat}&lon=${userLng}&format=json&accept-language=id`
          );
          const data = await res.json();
          locationText = data.display_name || `${userLat.toFixed(6)}, ${userLng.toFixed(6)}`;
        } catch {
          locationText = `${userLat.toFixed(6)}, ${userLng.toFixed(6)}`;
        }
        
        isLocating = false;
      },
      (error) => {
        isLocating = false;
        switch (error.code) {
          case error.PERMISSION_DENIED:
            locationText = 'Akses lokasi ditolak. Aktifkan GPS.';
            break;
          case error.POSITION_UNAVAILABLE:
            locationText = 'Lokasi tidak tersedia.';
            break;
          case error.TIMEOUT:
            locationText = 'Timeout. Coba lagi.';
            break;
          default:
            locationText = 'Gagal mendapatkan lokasi.';
        }
      },
      {
        enableHighAccuracy: true,
        timeout: 10000,
        maximumAge: 0
      }
    );
  }

  function nextStep() {
    if (step < 3) step++;
  }

  function prevStep() {
    if (step > 1) step--;
  }

  async function submitReport() {
    if (isSubmitting) return;
    
    isSubmitting = true;
    submitError = '';
    
    try {
      // 1. Upload photos
      let imageUrl = '';
      if (photos.length > 0) {
        const uploadResult = await api.uploadImage(photos[0].blob);
        imageUrl = uploadResult.url;
      }

      // 2. Upload additional photos as proofs (if any)
      // For now, we only use the first photo as main image

      // 3. Create report via API
      const result = await api.createReport({
        fingerprint_hash: fingerprint,
        latitude: userLat,
        longitude: userLng,
        image_url: imageUrl,
        severity: severityToNumber(severity as any),
        description: description || undefined,
      });

      submitSuccess = true;
      
      // Redirect after 2 seconds
      setTimeout(() => {
        goto('/');
      }, 2000);
      
    } catch (err: any) {
      submitError = err.message || 'Gagal mengirim laporan. Coba lagi.';
    } finally {
      isSubmitting = false;
    }
  }

  $: totalSaved = photos.reduce((sum, p) => sum + (p.originalSize - p.compressedSize), 0);
</script>

<svelte:head>
  <title>JEDUG - Lapor Jalan Rusak</title>
</svelte:head>

{#if submitSuccess}
  <div class="report-page">
    <div class="success-screen">
      <div class="success-icon">
        <CheckCircle size={64} />
      </div>
      <h2>Laporan Berhasil Dikirim! 🎉</h2>
      <p>Terima kasih sudah jadi warga peduli. Laporanmu akan membantu menekan pemerintah untuk bertindak.</p>
      <p class="redirect-text">Mengalihkan ke peta...</p>
    </div>
  </div>
{:else}
  <div class="report-page">
    <div class="report-container">
      <!-- Progress Steps -->
      <div class="progress-bar">
        {#each [1, 2, 3] as s}
          <div class="progress-step" class:active={step >= s} class:current={step === s}>
            <div class="step-circle">{s}</div>
            <span class="step-label">
              {s === 1 ? 'Foto' : s === 2 ? 'Detail' : 'Kirim'}
            </span>
          </div>
          {#if s < 3}
            <div class="progress-line" class:active={step > s}></div>
          {/if}
        {/each}
      </div>

      <!-- Step 1: Photo -->
      {#if step === 1}
        <div class="step-content">
          <div class="step-header">
            <h2>📸 Ambil Foto</h2>
            <p>Foto jalan rusak langsung dari kamera (wajib minimal 1 foto)</p>
          </div>

          <div class="photo-grid">
            {#each photos as photo, i}
              <div class="photo-item">
                <img src={photo.preview} alt="Foto {i + 1}" />
                <button class="photo-remove" on:click={() => removePhoto(i)}>
                  <X size={14} />
                </button>
                <div class="photo-size">
                  {formatFileSize(photo.compressedSize)}
                </div>
              </div>
            {/each}

            {#if photos.length < 3}
              <label class="photo-add">
                <input
                  type="file"
                  accept="image/*"
                  capture="environment"
                  on:change={handleFileSelect}
                  hidden
                />
                <div class="add-icon">
                  <Camera size={32} />
                  <span>Ambil Foto</span>
                </div>
              </label>
            {/if}
          </div>

          {#if totalSaved > 0}
            <div class="compression-info">
              <CheckCircle size={16} />
              <span>Gambar otomatis dikompresi ke WebP. Hemat {formatFileSize(totalSaved)} data!</span>
            </div>
          {/if}

          <div class="info-box">
            <Info size={16} />
            <span>Foto akan dikompresi otomatis ke WebP (max 1080px, 80% quality) untuk menghemat data. Maksimal 3 foto.</span>
          </div>
        </div>
      {/if}

      <!-- Step 2: Detail -->
      {#if step === 2}
        <div class="step-content">
          <div class="step-header">
            <h2>📋 Detail Laporan</h2>
            <p>Lengkapi informasi kerusakan jalan</p>
          </div>

          <!-- Location -->
          <div class="form-group">
            <label class="form-label">
              <MapPin size={16} />
              Lokasi
            </label>
            <div class="location-input">
              <input
                type="text"
                class="form-input"
                placeholder="Deteksi lokasi otomatis..."
                bind:value={locationText}
                readonly
              />
              <button class="location-btn" on:click={getLocation} disabled={isLocating}>
                {#if isLocating}
                  <div class="mini-spinner"></div>
                {:else}
                  <Crosshair size={18} />
                {/if}
              </button>
            </div>
            <span class="form-hint">GPS harus aktif. Laporan harus dari lokasi (maks 20m)</span>
            {#if userLat && userLng}
              <span class="form-hint coordinate">📍 {userLat.toFixed(6)}, {userLng.toFixed(6)}</span>
            {/if}
          </div>

          <!-- Severity -->
          <div class="form-group">
            <label class="form-label">
              <AlertTriangle size={16} />
              Tingkat Kerusakan
            </label>
            <div class="severity-grid">
              {#each severityOptions as opt}
                <button
                  class="severity-card"
                  class:selected={severity === opt.value}
                  on:click={() => severity = opt.value}
                >
                  <span class="severity-icon">{opt.icon}</span>
                  <span class="severity-label">{opt.label}</span>
                  <span class="severity-desc">{opt.desc}</span>
                </button>
              {/each}
            </div>
          </div>

          <!-- Description -->
          <div class="form-group">
            <label class="form-label" for="desc-input">
              📝 Deskripsi
            </label>
            <textarea
              id="desc-input"
              class="form-textarea"
              placeholder="Ceritakan kondisi jalannya... (opsional tapi membantu)"
              rows="3"
              bind:value={description}
              maxlength="500"
            ></textarea>
            <span class="form-hint">{description.length}/500 karakter</span>
          </div>
        </div>
      {/if}

      <!-- Step 3: Confirm -->
      {#if step === 3}
        <div class="step-content">
          <div class="step-header">
            <h2>✅ Konfirmasi Laporan</h2>
            <p>Periksa kembali sebelum mengirim</p>
          </div>

          <div class="preview-card">
            <!-- Photos preview -->
            {#if photos.length > 0}
              <div class="preview-photos">
                {#each photos as photo, i}
                  <img src={photo.preview} alt="Preview {i + 1}" class="preview-img" />
                {/each}
              </div>
            {:else}
              <div class="preview-no-photo">
                <Camera size={48} />
                <p>Belum ada foto</p>
              </div>
            {/if}

            <div class="preview-info">
              <div class="preview-row">
                <span class="preview-label">📍 Lokasi</span>
                <span class="preview-value">{locationText || 'Belum dideteksi'}</span>
              </div>
              <div class="preview-row">
                <span class="preview-label">⚠️ Tingkat</span>
                <span class="preview-value">
                  {severityOptions.find(o => o.value === severity)?.label || 'Belum dipilih'}
                </span>
              </div>
              {#if description}
                <div class="preview-row">
                  <span class="preview-label">📝 Deskripsi</span>
                  <span class="preview-value">{description}</span>
                </div>
              {/if}
            </div>
          </div>

          {#if submitError}
            <div class="error-box">
              <AlertTriangle size={16} />
              <span>{submitError}</span>
            </div>
          {/if}

          <label class="tos-checkbox">
            <input type="checkbox" bind:checked={agreedTos} />
            <span class="checkmark"></span>
            <span>Saya bertanggung jawab atas kebenaran data yang dilaporkan dan setuju dengan S&K JEDUG.</span>
          </label>
        </div>
      {/if}

      <!-- Actions -->
      <div class="form-actions">
        {#if step > 1}
          <button class="btn btn-secondary" on:click={prevStep} disabled={isSubmitting}>
            Kembali
          </button>
        {/if}

        {#if step < 3}
          <button
            class="btn btn-primary"
            on:click={nextStep}
            disabled={step === 1 && photos.length === 0}
          >
            Lanjut
          </button>
        {:else}
          <button
            class="btn btn-primary btn-submit"
            on:click={submitReport}
            disabled={!agreedTos || isSubmitting || !userLat || !severity}
          >
            {#if isSubmitting}
              <Loader2 size={18} class="spin" />
              Mengirim...
            {:else}
              <Send size={18} />
              Kirim Laporan
            {/if}
          </button>
        {/if}
      </div>
    </div>
  </div>
{/if}

<style>
  .report-page {
    min-height: calc(100dvh - var(--nav-height) - var(--bottom-nav-height));
    background: var(--bg-secondary);
    padding: var(--space-md);
    display: flex;
    justify-content: center;
  }

  .report-container {
    width: 100%;
    max-width: 600px;
  }

  /* Success Screen */
  .success-screen {
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    min-height: 60vh;
    text-align: center;
    padding: var(--space-xl);
    gap: var(--space-md);
  }
  .success-icon {
    color: var(--color-success);
    animation: scaleIn 0.5s ease;
  }
  @keyframes scaleIn {
    from { transform: scale(0); opacity: 0; }
    to { transform: scale(1); opacity: 1; }
  }
  .success-screen h2 {
    font-size: var(--text-2xl);
    font-weight: var(--font-bold);
  }
  .success-screen p {
    color: var(--text-secondary);
    font-size: var(--text-sm);
    max-width: 400px;
  }
  .redirect-text {
    color: var(--text-tertiary) !important;
    font-size: var(--text-xs) !important;
    margin-top: var(--space-lg);
  }

  /* Progress Bar */
  .progress-bar {
    display: flex;
    align-items: center;
    justify-content: center;
    padding: var(--space-lg) 0;
    gap: 0;
  }
  .progress-step {
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: 4px;
    position: relative;
  }
  .step-circle {
    width: 36px;
    height: 36px;
    border-radius: 50%;
    background: var(--bg-tertiary);
    color: var(--text-tertiary);
    display: flex;
    align-items: center;
    justify-content: center;
    font-weight: var(--font-bold);
    font-size: var(--text-sm);
    transition: all var(--transition-base);
    border: 2px solid transparent;
  }
  .progress-step.active .step-circle {
    background: var(--color-primary);
    color: white;
  }
  .progress-step.current .step-circle {
    border-color: var(--color-primary);
    box-shadow: 0 0 0 4px var(--color-primary-light);
  }
  .step-label {
    font-size: var(--text-xs);
    color: var(--text-tertiary);
    font-weight: var(--font-medium);
  }
  .progress-step.active .step-label {
    color: var(--color-primary);
  }
  .progress-line {
    width: 60px;
    height: 2px;
    background: var(--bg-tertiary);
    margin: 0 var(--space-sm);
    margin-bottom: 20px;
    transition: background var(--transition-base);
  }
  .progress-line.active {
    background: var(--color-primary);
  }

  /* Step Content */
  .step-content {
    background: var(--bg-card);
    border-radius: var(--radius-xl);
    padding: var(--space-xl);
    border: 1px solid var(--border-color);
    margin-bottom: var(--space-lg);
  }

  .step-header {
    margin-bottom: var(--space-xl);
  }
  .step-header h2 {
    font-size: var(--text-xl);
    font-weight: var(--font-bold);
    margin-bottom: var(--space-xs);
  }
  .step-header p {
    color: var(--text-secondary);
    font-size: var(--text-sm);
  }

  /* Photo Grid */
  .photo-grid {
    display: grid;
    grid-template-columns: repeat(auto-fill, minmax(140px, 1fr));
    gap: var(--space-md);
    margin-bottom: var(--space-lg);
  }

  .photo-item {
    position: relative;
    aspect-ratio: 4/3;
    border-radius: var(--radius-lg);
    overflow: hidden;
    border: 2px solid var(--border-color);
  }
  .photo-item img {
    width: 100%;
    height: 100%;
    object-fit: cover;
  }
  .photo-remove {
    position: absolute;
    top: 6px;
    right: 6px;
    width: 24px;
    height: 24px;
    border-radius: 50%;
    background: rgba(0, 0, 0, 0.6);
    color: white;
    display: flex;
    align-items: center;
    justify-content: center;
    cursor: pointer;
    border: none;
  }
  .photo-size {
    position: absolute;
    bottom: 4px;
    left: 4px;
    font-size: 0.6rem;
    color: white;
    background: rgba(0,0,0,0.6);
    padding: 1px 6px;
    border-radius: var(--radius-sm);
  }

  .photo-add {
    aspect-ratio: 4/3;
    border: 2px dashed var(--border-color-strong);
    border-radius: var(--radius-lg);
    display: flex;
    align-items: center;
    justify-content: center;
    cursor: pointer;
    transition: all var(--transition-fast);
  }
  .photo-add:hover {
    border-color: var(--color-primary);
    background: var(--color-primary-light);
  }
  .add-icon {
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: var(--space-sm);
    color: var(--text-tertiary);
    font-size: var(--text-sm);
    font-weight: var(--font-medium);
  }
  .photo-add:hover .add-icon {
    color: var(--color-primary);
  }

  .compression-info {
    display: flex;
    align-items: center;
    gap: var(--space-sm);
    padding: var(--space-sm) var(--space-md);
    background: var(--color-success-light);
    border-radius: var(--radius-md);
    color: var(--color-success);
    font-size: var(--text-xs);
    font-weight: var(--font-medium);
    margin-bottom: var(--space-md);
  }

  .info-box {
    display: flex;
    align-items: flex-start;
    gap: var(--space-sm);
    padding: var(--space-md);
    background: var(--color-info-light);
    border-radius: var(--radius-md);
    color: var(--color-info);
    font-size: var(--text-sm);
    line-height: var(--leading-relaxed);
  }

  .error-box {
    display: flex;
    align-items: center;
    gap: var(--space-sm);
    padding: var(--space-md);
    background: var(--color-danger-light);
    border-radius: var(--radius-md);
    color: var(--color-danger);
    font-size: var(--text-sm);
    margin-bottom: var(--space-md);
  }

  /* Form */
  .form-group {
    margin-bottom: var(--space-xl);
  }
  .form-label {
    display: flex;
    align-items: center;
    gap: var(--space-sm);
    font-weight: var(--font-semibold);
    font-size: var(--text-sm);
    margin-bottom: var(--space-sm);
    color: var(--text-primary);
  }
  .form-input {
    width: 100%;
    padding: 0.75rem 1rem;
    border: 1.5px solid var(--border-color);
    border-radius: var(--radius-lg);
    background: var(--bg-secondary);
    color: var(--text-primary);
    font-size: var(--text-base);
    transition: border-color var(--transition-fast);
  }
  .form-input:focus {
    border-color: var(--color-primary);
    box-shadow: 0 0 0 3px var(--color-primary-light);
  }
  .form-textarea {
    width: 100%;
    padding: 0.75rem 1rem;
    border: 1.5px solid var(--border-color);
    border-radius: var(--radius-lg);
    background: var(--bg-secondary);
    color: var(--text-primary);
    font-size: var(--text-base);
    resize: vertical;
    min-height: 80px;
    transition: border-color var(--transition-fast);
    line-height: var(--leading-relaxed);
  }
  .form-textarea:focus {
    border-color: var(--color-primary);
    box-shadow: 0 0 0 3px var(--color-primary-light);
  }
  .form-hint {
    display: block;
    margin-top: var(--space-xs);
    font-size: var(--text-xs);
    color: var(--text-tertiary);
  }
  .form-hint.coordinate {
    color: var(--color-success);
    font-weight: var(--font-medium);
  }

  /* Location input */
  .location-input {
    display: flex;
    gap: var(--space-sm);
  }
  .location-input .form-input {
    flex: 1;
  }
  .location-btn {
    width: 48px;
    height: 48px;
    border-radius: var(--radius-lg);
    background: var(--color-primary);
    color: white;
    display: flex;
    align-items: center;
    justify-content: center;
    cursor: pointer;
    transition: all var(--transition-fast);
    flex-shrink: 0;
    border: none;
  }
  .location-btn:hover {
    background: var(--color-primary-hover);
  }
  .location-btn:disabled {
    opacity: 0.6;
    cursor: not-allowed;
  }
  .mini-spinner {
    width: 18px;
    height: 18px;
    border: 2px solid rgba(255,255,255,0.3);
    border-top-color: white;
    border-radius: 50%;
    animation: spin 0.8s linear infinite;
  }

  @keyframes spin {
    to { transform: rotate(360deg); }
  }

  /* Severity Grid */
  .severity-grid {
    display: grid;
    grid-template-columns: repeat(2, 1fr);
    gap: var(--space-sm);
  }
  .severity-card {
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: 4px;
    padding: var(--space-md);
    border: 2px solid var(--border-color);
    border-radius: var(--radius-lg);
    background: var(--bg-secondary);
    cursor: pointer;
    transition: all var(--transition-fast);
    text-align: center;
  }
  .severity-card:hover {
    border-color: var(--border-color-strong);
    background: var(--bg-tertiary);
  }
  .severity-card.selected {
    border-color: var(--color-primary);
    background: var(--color-primary-light);
  }
  .severity-icon {
    font-size: 1.5rem;
  }
  .severity-label {
    font-size: var(--text-sm);
    font-weight: var(--font-semibold);
    color: var(--text-primary);
  }
  .severity-desc {
    font-size: var(--text-xs);
    color: var(--text-tertiary);
  }

  /* Preview */
  .preview-card {
    border: 1px solid var(--border-color);
    border-radius: var(--radius-lg);
    overflow: hidden;
    margin-bottom: var(--space-lg);
  }
  .preview-photos {
    display: flex;
    gap: 2px;
    height: 160px;
    overflow: hidden;
  }
  .preview-img {
    flex: 1;
    object-fit: cover;
    min-width: 0;
  }
  .preview-no-photo {
    height: 120px;
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    gap: var(--space-sm);
    background: var(--bg-tertiary);
    color: var(--text-tertiary);
  }
  .preview-info {
    padding: var(--space-md);
    display: flex;
    flex-direction: column;
    gap: var(--space-md);
  }
  .preview-row {
    display: flex;
    gap: var(--space-md);
  }
  .preview-label {
    font-size: var(--text-sm);
    color: var(--text-tertiary);
    min-width: 80px;
    flex-shrink: 0;
  }
  .preview-value {
    font-size: var(--text-sm);
    color: var(--text-primary);
    font-weight: var(--font-medium);
  }

  /* TOS checkbox */
  .tos-checkbox {
    display: flex;
    align-items: flex-start;
    gap: var(--space-sm);
    font-size: var(--text-sm);
    color: var(--text-secondary);
    cursor: pointer;
    line-height: var(--leading-relaxed);
    padding: var(--space-md);
    background: var(--color-warning-light);
    border-radius: var(--radius-md);
  }
  .tos-checkbox input[type="checkbox"] {
    width: 20px;
    height: 20px;
    margin-top: 2px;
    accent-color: var(--color-primary);
    flex-shrink: 0;
  }

  /* Actions */
  .form-actions {
    display: flex;
    gap: var(--space-md);
    padding-bottom: var(--space-xl);
  }

  .btn {
    flex: 1;
    padding: 0.875rem 1.5rem;
    border-radius: var(--radius-xl);
    font-size: var(--text-base);
    font-weight: var(--font-semibold);
    display: flex;
    align-items: center;
    justify-content: center;
    gap: var(--space-sm);
    cursor: pointer;
    transition: all var(--transition-fast);
    min-height: 52px;
    border: none;
  }
  .btn-primary {
    background: var(--color-primary);
    color: white;
  }
  .btn-primary:hover {
    background: var(--color-primary-hover);
  }
  .btn-primary:active {
    transform: scale(0.98);
  }
  .btn-primary:disabled {
    opacity: 0.5;
    cursor: not-allowed;
  }

  .btn-secondary {
    background: var(--bg-tertiary);
    color: var(--text-primary);
  }
  .btn-secondary:hover {
    background: var(--border-color);
  }

  .btn-submit {
    background: var(--color-success);
  }
  .btn-submit:hover {
    background: #2F855A;
  }

  :global(.spin) {
    animation: spin 1s linear infinite;
  }

  @media (min-width: 769px) {
    .report-page {
      padding: var(--space-xl);
    }
    .severity-grid {
      grid-template-columns: repeat(4, 1fr);
    }
  }
</style>
