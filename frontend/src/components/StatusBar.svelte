<script lang="ts">
  import { createEventDispatcher } from 'svelte';

  export let status: string = 'idle';
  export let activeMode: string = 'normal';
  export let activeTab: string = 'terminal';

  const dispatch = createEventDispatcher<{ restart: void; toggleMode: void; tabChange: string }>();

  const tabs = ['terminal', 'log'] as const;
  const tabLabels: Record<string, string> = { terminal: 'Terminal', log: 'Log' };

  const statusClasses: Record<string, string> = {
    idle: 'idle', thinking: 'thinking', responding: 'responding',
    stuck: 'stuck', error: 'error'
  };

  $: dotClass = statusClasses[status] || 'idle';
  $: showRestart = status === 'stuck' || status === 'error';
  $: isCustomProviderMode = activeMode === 'custom_provider';

  function handleRestart() { dispatch('restart'); }
  function handleToggleMode() { dispatch('toggleMode'); }
  function handleTabClick(tab: string) { dispatch('tabChange', tab); }
</script>

<div class="status-bar">
  <div class="tabs-row">
    {#each tabs as tab}
      <button
        class="tab-btn"
        class:active={activeTab === tab}
        on:click={() => handleTabClick(tab)}
      >{tabLabels[tab]}</button>
    {/each}
    <span class="tabs-filler"></span>
  </div>
  <div class="status-right">
    <span class="status-dot {dotClass}"></span>
    <button
      class="text-button"
      class:active={isCustomProviderMode}
      on:click={handleToggleMode}
    >{isCustomProviderMode ? 'Custom' : 'Default'}</button>
    {#if showRestart}
      <button class="text-button" on:click={handleRestart}>R</button>
    {/if}
  </div>
</div>

<style>
  .status-bar {
    padding-top: 0.2rem;
    padding-bottom: 0;
  }
  .tabs-row {
    display: flex;
    gap: 0;
    align-self: flex-end;
  }
  .tabs-filler {
    flex: 1;
    border-bottom: 1px solid var(--border);
  }
  .tab-btn {
    background: transparent;
    border: 1px solid transparent;
    border-bottom: none;
    color: var(--text-secondary);
    font-family: inherit;
    font-size: 0.72rem;
    padding: 2px 10px;
    cursor: pointer;
    border-radius: 4px 4px 0 0;
    transition: all var(--transition-fast);
  }
  .tab-btn:hover { color: var(--text-primary); background: rgba(255,255,255,0.03); }
  .tab-btn.active {
    color: var(--text-primary);
    background: #0f0f1a;
    border-color: var(--border);
  }
  .text-button {
    font-size: 0.7rem;
    padding: 1px 6px;
  }
  .text-button.active {
    color: var(--accent);
    background-color: rgba(74, 108, 247, 0.12);
  }
</style>
