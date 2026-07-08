<script lang="ts">
  import { afterUpdate } from 'svelte';
  import MessageBubble from './MessageBubble.svelte';
  import type { Message } from '../lib/types';

  export let messages: Message[] = [];
  export let streamingText: string = '';
  export let status: string = 'idle';

  let scrollContainer: HTMLDivElement;

  afterUpdate(() => {
    if (scrollContainer) {
      scrollContainer.scrollTop = scrollContainer.scrollHeight;
    }
  });
</script>

<div class="chat-scroll chat-view" bind:this={scrollContainer}>
  <div class="chat-container">
    {#if messages.length === 0 && !streamingText}
      <div class="empty-state">
        <svg class="empty-logo" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true">
          <path d="M12 2a10 10 0 0 1 10 10c0 5.523-4.477 10-10 10S2 17.523 2 12 6.477 2 12 2z"/>
          <path d="M8 14s1.5 2 4 2 4-2 4-2"/>
          <path d="M9 9h.01"/>
          <path d="M15 9h.01"/>
        </svg>
        <h2 class="empty-title">Welcome to CCNext</h2>
        <p class="empty-hint">Type a message below or tap the microphone to start talking with Claude.</p>
      </div>
    {:else}
      {#each messages as message, index (index)}
        <MessageBubble {message} />
      {/each}

      {#if streamingText}
        <MessageBubble message={{ role: 'claude', content: streamingText }} isStreaming={true} />
      {:else if status === 'thinking'}
        <div class="thinking-indicator">
          <span>Claude is thinking</span>
          <div class="thinking-dots">
            <span></span>
            <span></span>
            <span></span>
          </div>
        </div>
      {/if}
    {/if}
  </div>
</div>

<style>
  .chat-view {
    height: 100%;
  }
</style>
