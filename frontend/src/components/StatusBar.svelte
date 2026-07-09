<script lang="ts">
  import { createEventDispatcher } from 'svelte';

  export let status: string = 'idle';
  export let activeMode: string = 'normal';

  const dispatch = createEventDispatcher<{ restart: void; toggleMode: void }>();

  const statusClasses: Record<string, string> = {
    idle: 'idle',
    thinking: 'thinking',
    responding: 'responding',
    stuck: 'stuck',
    error: 'error'
  };

  const statusText: Record<string, string> = {
    idle: '就绪',
    running: '运行中',
    responding: '响应中',
    stuck: '卡住',
    stopped: '已停止',
    echo: '回显模式',
    error: '出错',
    thinking: '思考中'
  };

  $: dotClass = statusClasses[status] || 'idle';
  $: displayStatus = statusText[status] || status;
  $: showRestart = status === 'stuck' || status === 'error';
  $: isCustomProviderMode = activeMode === 'custom_provider';

  function handleRestart() {
    dispatch('restart');
  }

  function handleToggleMode() {
    dispatch('toggleMode');
  }
</script>

<div class="status-bar">
  <div class="status-left">
    <span>CCNext</span>
    <span class="status-dot {dotClass}"></span>
    <span class="status-text">{displayStatus}</span>
  </div>
  <div class="status-right">
    <button
      class="text-button {isCustomProviderMode ? 'active' : ''}"
      title={isCustomProviderMode ? '自研模式已开启' : '自研模式已关闭'}
      on:click={handleToggleMode}
    >
      自研模式
    </button>
    {#if showRestart}
      <button class="text-button" on:click={handleRestart}>Restart Claude</button>
    {/if}
  </div>
</div>

<style>
  .text-button.active {
    color: var(--accent);
    background-color: rgba(74, 108, 247, 0.12);
  }
</style>
