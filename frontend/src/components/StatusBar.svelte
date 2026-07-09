<script lang="ts">
  import { createEventDispatcher } from 'svelte';

  export let status: string = 'idle';
  export let activeMode: string = 'normal';

  const dispatch = createEventDispatcher<{ restart: void; modeChange: string }>();

  const statusClasses: Record<string, string> = {
    idle: 'idle',
    thinking: 'thinking',
    responding: 'responding',
    stuck: 'stuck',
    error: 'error'
  };

  const statusText: Record<string, string> = {
    idle: 'idle',
    running: 'running',
    responding: 'responding',
    stuck: 'stuck',
    stopped: 'stopped',
    echo: 'echo',
    error: 'error',
    thinking: 'thinking'
  };

  const modeLabels: Record<string, string> = {
    normal:  'Normal',
    plan:    'Plan',
    explore: 'Explore',
    ask:     'Ask',
    build:   'Build',
  };

  const modes = ['normal', 'plan', 'explore', 'ask', 'build'];

  $: dotClass = statusClasses[status] || 'idle';
  $: displayStatus = statusText[status] || status;
  $: showRestart = status === 'stuck' || status === 'error';

  function handleRestart() {
    dispatch('restart');
  }

  function handleModeChange(event: Event) {
    const target = event.target as HTMLSelectElement;
    dispatch('modeChange', target.value);
  }
</script>

<div class="status-bar">
  <div class="status-left">
    <span>CCNext</span>
    <span class="status-dot {dotClass}"></span>
    <span class="status-text">{displayStatus}</span>
  </div>
  <div class="status-right">
    <select
      class="mode-select"
      value={activeMode}
      on:change={handleModeChange}
      aria-label="Switch mode"
    >
      {#each modes as mode}
        <option value={mode}>{modeLabels[mode]}</option>
      {/each}
    </select>
    {#if showRestart}
      <button class="text-button" on:click={handleRestart}>Restart Claude</button>
    {/if}
  </div>
</div>

<style>
  .mode-select {
    background-color: var(--bg-input);
    border: 1px solid var(--border);
    border-radius: var(--radius-sm);
    padding: 0.25rem 0.5rem;
    color: var(--text-primary);
    font-family: inherit;
    font-size: 0.8rem;
    outline: none;
    cursor: pointer;
    transition: border-color var(--transition-fast);
  }
  .mode-select:focus {
    border-color: var(--accent);
    box-shadow: 0 0 0 2px rgba(74, 108, 247, 0.15);
  }
  .mode-select option {
    background-color: var(--bg-secondary);
    color: var(--text-primary);
  }
</style>
