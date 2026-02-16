<script lang="ts">
  import { page } from '$app/stores';
  import ThemeToggle from './ThemeToggle.svelte';
  import { Map, Camera, Trophy, User, Menu, X } from 'lucide-svelte';

  let mobileMenuOpen = false;

  const navItems = [
    { href: '/', label: 'Peta', icon: Map },
    { href: '/lapor', label: 'Lapor', icon: Camera },
    { href: '/ranking', label: 'Ranking', icon: Trophy },
    { href: '/profil', label: 'Profil', icon: User }
  ];

  $: currentPath = $page.url.pathname;

  function closeMobileMenu() {
    mobileMenuOpen = false;
  }
</script>

<header class="top-header">
  <div class="header-inner">
    <a href="/" class="logo" on:click={closeMobileMenu}>
      <img src="/android-chrome-192x192.png" alt="JEDUG" class="logo-img" />
      <span class="logo-text">JEDUG</span>
    </a>

    <!-- Desktop Nav -->
    <nav class="desktop-nav hide-mobile">
      {#each navItems as item}
        <a
          href={item.href}
          class="desktop-nav-item"
          class:active={item.href === '/' ? currentPath === '/' : currentPath.startsWith(item.href)}
        >
          <svelte:component this={item.icon} size={18} />
          <span>{item.label}</span>
        </a>
      {/each}
    </nav>

    <div class="header-actions">
      <ThemeToggle />
    </div>
  </div>
</header>

<style>
  .top-header {
    position: fixed;
    top: 0;
    left: 0;
    right: 0;
    height: var(--nav-height);
    background: var(--bg-nav);
    backdrop-filter: blur(12px);
    -webkit-backdrop-filter: blur(12px);
    border-bottom: 1px solid var(--border-color);
    z-index: var(--z-nav);
  }

  .header-inner {
    display: flex;
    align-items: center;
    justify-content: space-between;
    height: 100%;
    max-width: var(--max-content-width);
    margin: 0 auto;
    padding: 0 var(--space-lg);
  }

  .logo {
    display: flex;
    align-items: center;
    gap: var(--space-sm);
    text-decoration: none;
    color: var(--text-primary);
  }
  .logo:hover {
    text-decoration: none;
  }
  .logo-img {
    width: 32px;
    height: 32px;
    border-radius: var(--radius-sm);
    object-fit: contain;
  }
  .logo-text {
    font-size: var(--text-xl);
    font-weight: var(--font-extrabold);
    letter-spacing: -0.5px;
    color: var(--color-primary);
  }

  /* Desktop nav */
  .desktop-nav {
    display: flex;
    align-items: center;
    gap: var(--space-xs);
  }
  .desktop-nav-item {
    display: flex;
    align-items: center;
    gap: var(--space-sm);
    padding: var(--space-sm) var(--space-md);
    border-radius: var(--radius-lg);
    color: var(--text-secondary);
    font-size: var(--text-sm);
    font-weight: var(--font-medium);
    text-decoration: none;
    transition: all var(--transition-fast);
  }
  .desktop-nav-item:hover {
    background: var(--bg-tertiary);
    color: var(--text-primary);
    text-decoration: none;
  }
  .desktop-nav-item.active {
    background: var(--color-primary-light);
    color: var(--color-primary);
    font-weight: var(--font-semibold);
  }

  .header-actions {
    display: flex;
    align-items: center;
    gap: var(--space-sm);
  }
</style>
