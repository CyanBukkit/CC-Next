<script lang="ts">
  import { onMount } from 'svelte';
  import type { Settings } from '../lib/types';
  import { defaultSettings } from '../lib/types';

  export let settings: Settings | null = null;
  export let onSave: (settings: Settings) => void;
  export let onClose: () => void;

  let formSettings: Settings = { ...defaultSettings };
  let initialized = false;

  onMount(() => {
    if (settings) {
      formSettings = { ...settings };
      initialized = true;
    }
  });

  $: if (settings && !initialized) {
    formSettings = { ...settings };
    initialized = true;
  }

  $: customProviderMode = formSettings.active_mode === 'custom_provider';

  function handleSave() {
    onSave({ ...formSettings });
  }

  function handleOverlayClick(event: MouseEvent) {
    if (event.target === event.currentTarget) {
      onClose();
    }
  }
</script>

<div
	class="settings-overlay"
	role="button"
	tabindex="0"
	aria-label="Close settings"
	on:click={handleOverlayClick}
	on:keydown={(event) => event.key === 'Enter' && onClose()}
>
  <div class="settings-panel" role="dialog" aria-modal="true" aria-labelledby="settings-title">
    <div class="settings-header">
      <h2 id="settings-title" class="settings-title">设置</h2>
      <button class="settings-close" aria-label="关闭设置" on:click={onClose}>✕</button>
    </div>

    <div class="settings-form">
      {#if !settings}
        <p class="settings-description">正在加载设置…</p>
      {:else}
        <div class="settings-field">
          <label class="settings-label" for="stuck-timeout">卡住超时</label>
          <input
            id="stuck-timeout"
            class="settings-range"
            type="range"
            min="30"
            max="1800"
            step="5"
            bind:value={formSettings.stuck_timeout_seconds}
          />
          <span class="settings-description">{formSettings.stuck_timeout_seconds}秒 ({Math.floor(formSettings.stuck_timeout_seconds / 60)}分{formSettings.stuck_timeout_seconds % 60}秒)</span>
        </div>

        <label class="settings-field settings-toggle">
          <span class="settings-label">自动续接模式</span>
          <input class="toggle-input" type="checkbox" bind:checked={formSettings.auto_continue_mode} />
          <span class="toggle-switch"><span class="toggle-slider"></span></span>
        </label>

        <div class="settings-field">
          <label class="settings-label" for="voice-provider">语音提供商</label>
          <select id="voice-provider" class="settings-select" bind:value={formSettings.voice_provider}>
            <option value="Web Speech">Web Speech</option>
            <option value="Azure">Azure</option>
          </select>
        </div>

        <label class="settings-field settings-toggle">
          <span class="settings-label">深色模式</span>
          <input class="toggle-input" type="checkbox" bind:checked={formSettings.dark_mode} />
          <span class="toggle-switch"><span class="toggle-slider"></span></span>
        </label>

        <div class="settings-divider"></div>

        <label class="settings-field settings-toggle">
          <span class="settings-label">自研模式</span>
          <input
            class="toggle-input"
            type="checkbox"
            checked={customProviderMode}
            on:change={() => {
              formSettings.active_mode = customProviderMode ? 'normal' : 'custom_provider';
              formSettings.custom_provider_mode = !customProviderMode;
            }}
          />
          <span class="toggle-switch"><span class="toggle-slider"></span></span>
        </label>

        <div class="settings-field" class:disabled={!customProviderMode}>
          <label class="settings-label" for="provider-base-url">自定义 Provider Base URL</label>
          <input
            id="provider-base-url"
            class="settings-input"
            type="text"
            placeholder="https://api.openai.com/v1"
            disabled={!customProviderMode}
            bind:value={formSettings.provider_base_url}
          />
        </div>

        <div class="settings-field" class:disabled={!customProviderMode}>
          <label class="settings-label" for="provider-model">自定义 Provider Model</label>
          <input
            id="provider-model"
            class="settings-input"
            type="text"
            placeholder="gpt-4o"
            disabled={!customProviderMode}
            bind:value={formSettings.provider_model}
          />
        </div>

        <div class="settings-field" class:disabled={!customProviderMode}>
          <label class="settings-label" for="provider-auth-token">自定义 Provider Auth Token</label>
          <input
            id="provider-auth-token"
            class="settings-input"
            type="password"
            placeholder="sk-..."
            disabled={!customProviderMode}
            bind:value={formSettings.provider_auth_token}
          />
        </div>

        <div class="settings-field">
          <label class="settings-label" for="max-history">最大历史消息数</label>
          <input
            id="max-history"
            class="settings-input"
            type="number"
            min="1"
            bind:value={formSettings.max_history_messages}
          />
        </div>

        <div class="settings-field">
          <label class="settings-label" for="passphrase-length">口令长度</label>
          <input
            id="passphrase-length"
            class="settings-input"
            type="number"
            min="1"
            bind:value={formSettings.passphrase_length}
          />
        </div>

        <div class="settings-field">
          <label class="settings-label" for="passphrase-template">隐藏指令模板</label>
          <p class="settings-description">%random% 会被替换为随机生成的口令</p>
          <textarea
            id="passphrase-template"
            class="settings-input settings-textarea"
            rows="3"
            bind:value={formSettings.passphrase_template}
          ></textarea>
        </div>

        <div class="settings-field">
          <label class="settings-label" for="claude-model">Claude 模型</label>
          <input
            id="claude-model"
            class="settings-input"
            type="text"
            bind:value={formSettings.claude_model}
          />
        </div>
      {/if}
    </div>

    <div class="settings-footer">
      <button class="settings-button primary" disabled={!settings} on:click={handleSave}>保存</button>
      <button class="settings-button secondary" on:click={onClose}>取消</button>
    </div>
  </div>
</div>
