<script lang="ts">
  import { page } from '$app/stores';
  import { Map, Camera, Trophy, User } from 'lucide-svelte';

  const navItems = [
    { href: '/', label: 'Peta', icon: Map },
    { href: '/ranking', label: 'Ranking', icon: Trophy },
    { href: '/profil', label: 'Profil', icon: User }
  ];

  $: currentPath = $page.url.pathname;
</script>

<nav class="bottom-nav">
  <div class="nav-items">
    {#each navItems as item, i}
      {#if i === 1}
        <!-- FAB in the middle -->
        <a href="/lapor" class="fab" class:active={currentPath === '/lapor'} aria-label="Lapor Jalan Rusak">
          <div class="fab-inner">
            <Camera size={24} strokeWidth={2.5} />
          </div>
          <span class="fab-label">Lapor</span>
        </a>
      {/if}
      <a
        href={item.href}
        class="nav-item"
        class:active={item.href === '/' ? currentPath === '/' : currentPath.startsWith(item.href)}
        aria-label={item.label}
      >
        <svelte:component this={item.icon} size={22} />
        <span>{item.label}</span>
      </a>
    {/each}
  </div>
</nav>

<style>
  .bottom-nav {
    position: fixed;
    bottom: 0;
    left: 0;
    right: 0;
    height: var(--bottom-nav-height);
    background: var(--bg-nav);
    backdrop-filter: blur(12px);
    -webkit-backdrop-filter: blur(12px);
    border-top: 1px solid var(--border-color);
    box-shadow: var(--shadow-nav);
    z-index: var(--z-nav);
    padding-bottom: env(safe-area-inset-bottom, 0);
  }

  .nav-items {
    display: flex;
    align-items: flex-end;
    justify-content: space-around;
    height: 100%;
    max-width: 500px;
    margin: 0 auto;
    padding: 0 var(--space-sm);
  }

  .nav-item {
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    gap: 2px;
    padding: var(--space-sm) var(--space-md);
    color: var(--text-tertiary);
    text-decoration: none;
    font-size: var(--text-xs);
    font-weight: var(--font-medium);
    transition: color var(--transition-fast);
    min-width: 64px;
    height: 100%;
  }
  .nav-item:hover {
    color: var(--text-secondary);
    text-decoration: none;
  }
  .nav-item.active {
    color: var(--color-primary);
  }
  .nav-item span {
    margin-top: 2px;
  }

  /* FAB */
  .fab {
    display: flex;
    flex-direction: column;
    align-items: center;
    text-decoration: none;
    position: relative;
    margin-top: -1.25rem;
  }
  .fab:hover {
    text-decoration: none;
  }
  .fab-inner {
    width: var(--fab-size);
    height: var(--fab-size);
    border-radius: var(--radius-full);
    background: var(--color-primary);
    color: white;
    display: flex;
    align-items: center;
    justify-content: center;
    box-shadow: var(--shadow-fab);
    transition: all var(--transition-fast);
  }
  .fab:hover .fab-inner {
    background: var(--color-primary-hover);
    transform: scale(1.08);
  }
  .fab:active .fab-inner {
    transform: scale(0.95);
  }
  .fab.active .fab-inner {
    background: var(--color-primary-dark);
  }
  .fab-label {
    font-size: var(--text-xs);
    font-weight: var(--font-semibold);
    color: var(--color-primary);
    margin-top: 4px;
  }

  @media (min-width: 769px) {
    .bottom-nav {
      display: none;
    }
  }
</style>
