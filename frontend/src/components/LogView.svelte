<script lang="ts">
  export let text: string = '';

  let container: HTMLDivElement;
  $: if (container) {
    container.scrollTop = container.scrollHeight;
  }
</script>

<div bind:this={container} class="log-view">
  <pre>{text || 'Waiting for output...'}</pre>
</div>

<style>
  .log-view {
    width: 100%;
    height: 100%;
    min-height: 0;
    overflow-y: auto;
    background: #0f0f1a;
    color: #d4d4e8;
    font-family: var(--font-mono);
    font-size: 13px;
    line-height: 1.45;
    padding: 6px 8px;
  }
  .log-view pre {
    margin: 0;
    white-space: pre-wrap;
    word-break: break-all;
  }
  .log-view::-webkit-scrollbar {
    width: 5px;
  }
  .log-view::-webkit-scrollbar-track {
    background: #0f0f1a;
  }
  .log-view::-webkit-scrollbar-thumb {
    background: #2a2a45;
    border-radius: 3px;
  }
</style>
