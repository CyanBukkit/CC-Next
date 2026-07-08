<script lang="ts">
  import { createEventDispatcher } from 'svelte';
  export let status: string = 'idle';

  const dispatch = createEventDispatcher<{ restart: void }>();

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

  function handleRestart() {
    dispatch('restart');
  }
</script>

<div class="status-bar">
  <div class="status-left">
    <span>CCNext</span>
    <span class="status-dot {dotClass}"></span>
    <span class="status-text">{displayStatus}</span>
  </div>
  <div class="status-right">
    {#if showRestart}
      <button class="text-button" on:click={handleRestart}>重启 Claude</button>
    {/if}
  </div>
</div>
