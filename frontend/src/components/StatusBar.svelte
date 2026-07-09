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

  const statusText: Record<string, string> = {
    idle: '', running: '', responding: '',
    stuck: '', stopped: '', echo: '', error: '', thinking: ''
  };

  $: dotClass = statusClasses[status] || 'idle';
  $: displayStatus = statusText[status] || status;
  $: showRestart = status === 'stuck' || status === 'error';
  $: isCustomProviderMode = activeMode === 'custom_provider';

  function handleRestart() { dispatch('restart'); }
  function handleToggleMode() { dispatch('toggleMode'); }
  function handleTabClick(tab: string) { dispatch('tabChange', tab); }
</script>

<div class="status-bar">
  <div class="status-left">
    <span class="status-dot {dotClass}"></span>
    {#each tabs as tab}
      <button
        class="tab-btn"
        class:active={activeTab === tab}
        on:click={() => handleTabClick(tab)}
      >{tabLabels[tab]}</button>
    {/each}
  </div>
  <div class="status-right">
    <button
      class="text-button"
      class:active={isCustomProviderMode}
      title={isCustomProviderMode ? 'Custom' : 'Default'}
      on:click={handleToggleMode}
    >{isCustomProviderMode ? 'Custom' : 'Default'}</button>
    {#if showRestart}
      <button class="text-button" on:click={handleRestart}>R</button>
    {/if}
  </div>
</div>

<style>
  .tab-btn {
    background: transparent;
    border: none;
    color: var(--text-secondary);
    font-family: inherit;
    font-size: 0.7rem;
    padding: 1px 6px;
    cursor: pointer;
    border-radius: 3px;
    transition: all var(--transition-fast);
  }
  .tab-btn:hover { color: var(--text-primary); }
  .tab-btn.active {
    color: var(--text-primary);
    background: var(--bg-tertiary);
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
