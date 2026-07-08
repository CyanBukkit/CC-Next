<script lang="ts">
  import { marked } from 'marked';
  import type { Message } from '../lib/types';

  export let message: Message;
  export let isStreaming: boolean = false;

  let renderedContent: string = '';

  async function render(text: string): Promise<void> {
    renderedContent = await marked.parse(text);
  }

  $: render(message.content);
</script>

<div class="message-row {message.role}">
  <div class="message-bubble {message.role}">
    {@html renderedContent}
    {#if isStreaming}
      <span class="streaming-cursor"></span>
    {/if}
  </div>
</div>
