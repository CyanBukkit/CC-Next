<script lang="ts">
  import { isSupported, startListening, stopListening, setInterimResultCallback } from '../lib/voice';

  export let onVoiceResult: (text: string) => void;
  export let disabled: boolean = false;

  let isRecording = false;
  let hasError = false;
  let interimText = '';

  function showErrorBriefly() {
    hasError = true;
    setTimeout(() => {
      hasError = false;
    }, 2000);
  }

  function handleClick() {
    if (isRecording) {
      stopListening();
      isRecording = false;
      return;
    }

    if (!isSupported()) {
      showErrorBriefly();
      return;
    }

    hasError = false;
    interimText = '';
    setInterimResultCallback((text: string) => {
      interimText = text;
    });

    isRecording = true;
    startListening()
      .then((text: string) => {
        onVoiceResult(text || interimText);
      })
      .catch((error) => {
        console.error('Voice recognition failed:', error);
        showErrorBriefly();
      })
      .finally(() => {
        isRecording = false;
        interimText = '';
      });
  }
</script>

<button
  class="icon-button voice-button"
  class:recording={isRecording}
  class:error={hasError}
  title={isRecording ? '停止聆听' : '语音输入'}
  aria-label={isRecording ? '停止聆听' : '语音输入'}
  disabled={disabled}
  on:click={handleClick}
>
  <svg viewBox="0 0 24 24" width="20" height="20" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true">
    {#if isRecording}
      <rect x="6" y="6" width="12" height="12" rx="2" />
    {:else}
      <path d="M12 2a3 3 0 0 1 3 3v7a3 3 0 0 1-6 0V5a3 3 0 0 1 3-3z" />
      <path d="M19 10v2a7 7 0 0 1-14 0v-2" />
      <line x1="12" y1="19" x2="12" y2="22" />
      <line x1="8" y1="22" x2="16" y2="22" />
    {/if}
  </svg>
</button>
