<script lang="ts">
  export let visible: boolean = false;
  export let onContinue: () => void;
  export let onDismiss: () => void;
</script>

{#if visible}
  <div class="stuck-banner" role="alert">
    <div class="stuck-banner-text">
      <span>⚠️</span>
      <span>Claude 似乎停止响应了</span>
    </div>
    <div class="stuck-banner-actions">
      <button class="settings-button primary" on:click={onContinue}>继续</button>
      <button class="settings-button secondary" on:click={onDismiss}>忽略</button>
    </div>
  </div>
{/if}
