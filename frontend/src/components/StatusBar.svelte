<script lang="ts">
  import { createEventDispatcher } from 'svelte';

  export let status: string = 'idle';
  export let activeMode: string = 'normal';
  export let activeTab: string = 'terminal';

  const dispatch = createEventDispatcher<{ restart: void; toggleMode: void; tabChange: string; modeChange: string }>();

  const tabs = ['terminal', 'log'] as const;
  const tabLabels: Record<string, string> = { terminal: 'Terminal', log: 'Log' };

  const modes = ['normal', 'plan', 'explore', 'ask', 'build', 'hack'];
  const modeLabels: Record<string, string> = {
    normal: 'Normal', plan: 'Plan', explore: 'Explore', ask: 'Ask', build: 'Build',
    hack: 'Hack',
  };

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
  function handleModeClick(mode: string) { dispatch('modeChange', mode); }
</script>

<div class="status-bar">
  <div class="tabs-row">
    {#each tabs as tab}
      <button class="tab-btn" class:active={activeTab === tab}
              on:click={() => handleTabClick(tab)}>{tabLabels[tab]}</button>
    {/each}
    <button class="hack-btn" on:click={() => handleModeClick("hack")}>HACK</button>
    <span class="tabs-filler"></span>
  </div>
  <div class="status-right">
    <span class="status-dot {dotClass}"></span>
    {#each modes as m}
      <button class="mode-btn" class:active={activeMode === m} class:glow={m === "hack"}
              on:click={() => handleModeClick(m)}>{modeLabels[m]}</button>
    {/each}
    <button class="text-button" class:active={isCustomProviderMode}
            on:click={handleToggleMode}>{isCustomProviderMode ? 'Custom' : 'Default'}</button>
    {#if showRestart}
      <button class="text-button" on:click={handleRestart}>R</button>
    {/if}
  </div>
</div>

<style>
  .status-bar {
    padding-top: 0.2rem;
    padding-bottom: 0;
    border-bottom: none;
    background: var(--bg-secondary);
  }
  .tabs-row { display: flex; gap: 0; align-self: flex-end; }
  .tabs-filler { flex: 1; min-width: 12px; border-bottom: 1px solid var(--border); }
  .tab-btn {
    background: transparent;
    border: 1px solid transparent;
    border-bottom: 1px solid var(--border);
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
    border-bottom-color: #0f0f1a;
  }
  .hack-btn {
    background: transparent;
    border: none;
    border-bottom: 1px solid #7c3aed;
    color: #a855f7;
    font-family: inherit;
    font-size: 0.72rem;
    font-weight: 700;
    padding: 2px 8px;
    cursor: pointer;
    letter-spacing: 1px;
    text-shadow: 0 0 6px rgba(168, 85, 247, 0.4);
    transition: all var(--transition-fast);
  }
  .hack-btn:hover {
    color: #c084fc;
    text-shadow: 0 0 10px rgba(168, 85, 247, 0.6);
  }
  .mode-btn {
    background: transparent;
    border: none;
    border-bottom: 2px solid transparent;
    color: var(--text-secondary);
    font-family: inherit;
    font-size: 0.68rem;
    padding: 1px 5px;
    cursor: pointer;
    transition: all var(--transition-fast);
  }
  .mode-btn:hover { color: var(--text-primary); }
  .mode-btn.active {
    color: var(--accent);
    border-bottom-color: var(--accent);
  }
  .mode-btn.glow {
    color: #c084fc;
    text-shadow: 0 0 8px rgba(168, 85, 247, 0.6);
  }
  .mode-btn.glow:hover { text-shadow: 0 0 14px rgba(168, 85, 247, 0.8); }
  .text-button { font-size: 0.7rem; padding: 1px 6px; border-radius: var(--radius-sm); }
  .text-button.default-mode { background-color: rgba(74, 108, 247, 0.10); border: 1px solid rgba(74, 108, 247, 0.25); }
  .text-button.active { color: var(--accent); background-color: rgba(74, 108, 247, 0.12); }
</style>
