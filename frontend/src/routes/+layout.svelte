<script lang="ts">
  import '$lib/styles/global.css';
  import TopHeader from '$lib/components/TopHeader.svelte';
  import BottomNav from '$lib/components/BottomNav.svelte';
</script>

<div class="app-shell">
  <TopHeader />
  <main class="main-content">
    <slot />
  </main>
  <BottomNav />
</div>

<style>
  .app-shell {
    display: flex;
    flex-direction: column;
    min-height: 100dvh;
  }

  .main-content {
    flex: 1;
    margin-top: var(--nav-height);
    padding-bottom: var(--bottom-nav-height);
  }

  @media (min-width: 769px) {
    .main-content {
      padding-bottom: 0;
    }
  }
</style>
