<script lang="ts">
  import { createEventDispatcher } from 'svelte';
  import VoiceButton from './VoiceButton.svelte';
  import { Keys, sendPTYKey, pickWorkDir } from '../lib/wails-runtime';

  export let value: string = '';
  export let disabled: boolean = false;

  const dispatch = createEventDispatcher<{ send: string; settings: void }>();
  let textarea: HTMLTextAreaElement;

  $: isBusy = disabled;
  $: canSend = value.trim().length > 0 && !isBusy;
  $: rows = Math.min(6, Math.max(1, value.split('\n').length));

  function handleKeydown(event: KeyboardEvent) {
    if (event.key === 'Enter' && !event.shiftKey) {
      event.preventDefault();
      handleSend();
    }
  }

  function handleSend() {
    if (canSend) {
      dispatch('send', value);
    }
  }

  function handleVoiceResult(text: string) {
    value = value ? `${value} ${text}` : text;
  }
</script>

<div class="input-bar">
  <!-- Shortcut keys row -->
  <div class="shortcut-bar">
    <button class="shortcut-btn" title="返回 (ESC)" on:click={() => sendPTYKey(Keys.esc)}>
      <svg viewBox="0 0 24 24" width="12" height="12" fill="none" stroke="currentColor" stroke-width="2.5">
        <polyline points="11 7 6 12 11 17"/>
        <polyline points="6 12 18 12"/>
      </svg>
      <span>返回</span>
    </button>
    <div class="shortcut-sep"></div>
    <button class="shortcut-btn" title="回车确认 (Enter)" on:click={() => sendPTYKey(Keys.enter)}>
      <svg viewBox="0 0 24 24" width="12" height="12" fill="none" stroke="currentColor" stroke-width="2.5">
        <polyline points="9 17 4 12 9 7"/>
        <path d="M4 12h12c2.2 0 4 1.8 4 4v4"/>
      </svg>
      <span>回车</span>
    </button>
    <button class="shortcut-btn" title="空格选中 (Space)" on:click={() => sendPTYKey(Keys.space)}>
      <svg viewBox="0 0 24 24" width="12" height="12" fill="none" stroke="currentColor" stroke-width="2.5">
        <line x1="4" y1="12" x2="20" y2="12"/>
      </svg>
      <span>空格</span>
    </button>
    <div class="shortcut-sep"></div>
    <button class="shortcut-btn" title="后台运行 (Ctrl+B)" on:click={() => sendPTYKey(Keys.ctrlB)}>
      <svg viewBox="0 0 24 24" width="12" height="12" fill="none" stroke="currentColor" stroke-width="2.5">
        <rect x="2" y="3" width="20" height="14" rx="2"/>
        <line x1="8" y1="21" x2="16" y2="21"/>
        <line x1="12" y1="17" x2="12" y2="21"/>
      </svg>
      <span>后台</span>
    </button>
    <button class="shortcut-btn" title="切换模式 (Shift+Tab)" on:click={() => sendPTYKey(Keys.shiftTab)}>
      <svg viewBox="0 0 24 24" width="12" height="12" fill="none" stroke="currentColor" stroke-width="2.5">
        <polyline points="7 7 3 12 7 17"/>
        <line x1="21" y1="12" x2="3" y2="12"/>
      </svg>
      <span>模式</span>
    </button>
    <button class="shortcut-btn" title="查看思考 (Ctrl+O)" on:click={() => sendPTYKey(Keys.ctrlO)}>
      <svg viewBox="0 0 24 24" width="12" height="12" fill="none" stroke="currentColor" stroke-width="2.5">
        <circle cx="12" cy="12" r="10"/>
        <circle cx="12" cy="12" r="3"/>
        <line x1="12" y1="2" x2="12" y2="6"/>
        <line x1="12" y1="18" x2="12" y2="22"/>
        <line x1="2" y1="12" x2="6" y2="12"/>
        <line x1="18" y1="12" x2="22" y2="12"/>
      </svg>
      <span>思考</span>
    </button>
    <button class="shortcut-btn" title="上 (↑)" on:click={() => sendPTYKey(Keys.up)}>
      <svg viewBox="0 0 24 24" width="12" height="12" fill="none" stroke="currentColor" stroke-width="2.5">
        <polyline points="7 14 12 9 17 14"/>
      </svg>
      <span>上</span>
    </button>
    <button class="shortcut-btn" title="下 (↓)" on:click={() => sendPTYKey(Keys.down)}>
      <svg viewBox="0 0 24 24" width="12" height="12" fill="none" stroke="currentColor" stroke-width="2.5">
        <polyline points="7 10 12 15 17 10"/>
      </svg>
      <span>下</span>
    </button>
    <div class="shortcut-sep"></div>
    <button class="shortcut-btn" title="切换目录" on:click={pickWorkDir}>
      <svg viewBox="0 0 24 24" width="12" height="12" fill="none" stroke="currentColor" stroke-width="2.5">
        <path d="M2 6a2 2 0 0 1 2-2h5l2 2h9a2 2 0 0 1 2 2v10a2 2 0 0 1-2 2H4a2 2 0 0 1-2-2V6z"/>
      </svg>
      <span>目录</span>
    </button>
    <button class="shortcut-btn" title="刷新终端 (Ctrl+L)" on:click={() => sendPTYKey(Keys.ctrlL)}>
      <svg viewBox="0 0 24 24" width="12" height="12" fill="none" stroke="currentColor" stroke-width="2.5">
        <polyline points="1 4 1 10 7 10"/>
        <path d="M3.5 16.5A9 9 0 1 0 2 12"/>
      </svg>
      <span>刷新</span>
    </button>
  </div>

  <!-- Text input row -->
  <div class="input-row">
    <div class="input-wrapper">
      <textarea
        bind:this={textarea}
        bind:value
        class="input-textarea"
        {rows}
        placeholder="输入消息..."
        disabled={isBusy}
        on:keydown={handleKeydown}
      ></textarea>
      <div class="input-actions">
        <span class="char-count">{value.length}</span>
        <VoiceButton onVoiceResult={handleVoiceResult} disabled={isBusy} />
        <button class="icon-button" title="设置" on:click={() => dispatch('settings')}>
          <svg viewBox="0 0 24 24" width="20" height="20" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true">
            <circle cx="12" cy="12" r="3"></circle>
            <path d="M19.4 15a1.65 1.65 0 0 0 .33 1.82l.06.06a2 2 0 0 1 0 2.83 2 2 0 0 1-2.83 0l-.06-.06a1.65 1.65 0 0 0-1.82-.33 1.65 1.65 0 0 0-1 1.51V21a2 2 0 0 1-2 2 2 2 0 0 1-2-2v-.09A1.65 1.65 0 0 0 9 19.4a1.65 1.65 0 0 0-1.82.33l-.06.06a2 2 0 0 1-2.83 0 2 2 0 0 1 0-2.83l.06-.06a1.65 1.65 0 0 0 .33-1.82 1.65 1.65 0 0 0-1.51-1H3a2 2 0 0 1-2-2 2 2 0 0 1 2-2h.09A1.65 1.65 0 0 0 4.6 9a1.65 1.65 0 0 0-.33-1.82l-.06-.06a2 2 0 0 1 0-2.83 2 2 0 0 1 2.83 0l.06.06a1.65 1.65 0 0 0 1.82.33H9a1.65 1.65 0 0 0 1-1.51V3a2 2 0 0 1 2-2 2 2 0 0 1 2 2v.09a1.65 1.65 0 0 0 1 1.51 1.65 1.65 0 0 0 1.82-.33l.06-.06a2 2 0 0 1 2.83 0 2 2 0 0 1 0 2.83l-.06.06a1.65 1.65 0 0 0-.33 1.82V9a1.65 1.65 0 0 0 1.51 1H21a2 2 0 0 1 2 2 2 2 0 0 1-2 2h-.09a1.65 1.65 0 0 0-1.51 1z"></path>
          </svg>
        </button>
      </div>
    </div>
    <button class="send-button" title="发送" disabled={!canSend} on:click={handleSend}>
      <svg viewBox="0 0 24 24" width="20" height="20" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true">
        <line x1="22" y1="2" x2="11" y2="13"></line>
        <polygon points="22 2 15 22 11 13 2 9 22 2"></polygon>
      </svg>
    </button>
  </div>
</div>

<style>
  .input-bar {
    width: 100%;
    box-sizing: border-box;
    padding: 0.3rem 0.5rem 0.5rem;
    background-color: var(--bg-secondary);
    border-top: 1px solid var(--border);
  }

  .shortcut-bar {
    display: flex;
    flex-wrap: wrap;
    gap: 0.25rem;
    margin-bottom: 0.2rem;
    justify-content: center;
  }

  .input-row {
    display: flex;
    align-items: flex-end;
    gap: 0.6rem;
    width: 100%;
  }

  .shortcut-btn {
    display: inline-flex;
    align-items: center;
    gap: 4px;
    padding: 1px 5px;
    border-radius: 4px;
    border: 1px solid var(--border);
    background: var(--bg-tertiary);
    color: var(--text-secondary);
    font-size: 0.7rem;
    font-family: inherit;
    cursor: pointer;
    transition: all var(--transition-fast);
  }

  .shortcut-btn:hover {
    background: var(--accent);
    color: #fff;
    border-color: var(--accent);
  }

  .shortcut-btn:active {
    transform: scale(0.95);
  }

  .shortcut-sep {
    width: 1px;
    background: var(--border);
    margin: 0 2px;
    align-self: stretch;
  }

  .input-wrapper {
    flex: 1;
    display: flex;
    align-items: flex-end;
    background-color: var(--bg-input);
    border: 1px solid var(--border);
    border-radius: var(--radius-lg);
    padding: 0.25rem 0.5rem;
    transition: border-color var(--transition-fast), box-shadow var(--transition-fast);
  }

  .input-wrapper:focus-within {
    border-color: var(--accent);
    box-shadow: 0 0 0 3px rgba(74, 108, 247, 0.15);
  }

  .input-textarea {
    flex: 1;
    resize: none;
    border: none;
    outline: none;
    background: transparent;
    color: var(--text-primary);
    font-family: inherit;
    font-size: 0.82rem;
    line-height: 1.4;
    max-height: 100px;
    min-height: 22px;
    padding: 0;
    margin-right: 0.5rem;
  }

  .input-textarea::placeholder {
    color: var(--text-secondary);
  }

  .input-actions {
    display: flex;
    align-items: center;
    gap: 0.3rem;
  }

  .char-count {
    font-size: 0.7rem;
    color: var(--text-secondary);
    user-select: none;
    min-width: 20px;
    text-align: right;
  }

  .icon-button {
    display: inline-flex;
    align-items: center;
    justify-content: center;
    width: 34px;
    height: 34px;
    border-radius: 50%;
    border: none;
    background-color: transparent;
    color: var(--text-secondary);
    cursor: pointer;
    transition: background-color var(--transition-fast), color var(--transition-fast);
  }

  .icon-button:hover {
    background-color: rgba(255, 255, 255, 0.08);
    color: var(--text-primary);
  }

  .send-button {
    display: inline-flex;
    align-items: center;
    justify-content: center;
    width: 40px;
    height: 40px;
    border-radius: 50%;
    border: none;
    background-color: var(--accent);
    color: #ffffff;
    cursor: pointer;
    transition: background-color var(--transition-fast), transform var(--transition-fast);
    flex-shrink: 0;
  }

  .send-button:hover:not(:disabled) {
    background-color: var(--accent-hover);
  }

  .send-button:active:not(:disabled) {
    transform: scale(0.95);
  }

  .send-button:disabled {
    background-color: var(--bg-tertiary);
    color: var(--text-secondary);
    cursor: not-allowed;
  }
</style>
