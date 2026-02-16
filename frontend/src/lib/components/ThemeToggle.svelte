<script lang="ts">
  import { theme } from '$lib/stores/theme';
  import { Sun, Moon } from 'lucide-svelte';

  function toggle() {
    theme.toggle();
  }
</script>

<button
  class="theme-toggle"
  on:click={toggle}
  aria-label="Toggle theme"
  title={$theme === 'light' ? 'Aktifkan Mode Gelap' : 'Aktifkan Mode Terang'}
>
  {#if $theme === 'light'}
    <Moon size={20} />
  {:else}
    <Sun size={20} />
  {/if}
</button>

<style>
  .theme-toggle {
    display: flex;
    align-items: center;
    justify-content: center;
    width: 40px;
    height: 40px;
    border-radius: var(--radius-full);
    background: var(--bg-tertiary);
    color: var(--text-secondary);
    transition: all var(--transition-fast);
    cursor: pointer;
    border: none;
  }
  .theme-toggle:hover {
    background: var(--border-color-strong);
    color: var(--text-primary);
  }
  .theme-toggle:active {
    transform: scale(0.92);
  }
</style>
