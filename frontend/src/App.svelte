<script lang="ts">
  import { onMount } from 'svelte';
  import StatusBar from './components/StatusBar.svelte';
  import TerminalView from './components/TerminalView.svelte';
  import StuckBanner from './components/StuckBanner.svelte';
  import InputBar from './components/InputBar.svelte';
  import SettingsPanel from './components/SettingsPanel.svelte';
  import {
    sendMessage,
    getSettings,
    updateSettings,
    getStatus,
    restartClaude,
    clearHistory,
    newSession,
    resizeTerminal,
    sendPTYKey,
    Keys,
    EventsOn,
    EventsOff,
  } from './lib/wails-runtime';
  import type { Settings } from './lib/types';

  let terminal: TerminalView;
  let status: string = 'idle';
  let settings: Settings | null = null;
  let stuckVisible: boolean = false;
  let settingsOpen: boolean = false;
  let inputText: string = '';

  $: activeMode = settings?.active_mode ?? 'normal';

  onMount(() => {
    loadInitialState();

    // Raw PTY output → xterm.js
    EventsOn('pty-output', (data: any) => {
      if (terminal && data) {
        if (typeof data === 'string') {
          terminal.write(data);
        } else if (typeof data.text === 'string') {
          terminal.write(data.text);
        }
      }
    });

    // Passphrase detected → task complete
    EventsOn('claude-done', () => {
      status = 'idle';
    });

    EventsOn('claude-stuck', () => {
      status = 'stuck';
      stuckVisible = true;
    });

    EventsOn('claude-status', (data: any) => {
      if (data && typeof data.status === 'string') {
        status = data.status;
        if (status !== 'stuck') stuckVisible = false;
      }
    });

    EventsOn('claude-error', (data: any) => {
      status = 'error';
      if (terminal && data && typeof data === 'string') {
        terminal.writeln('\r\n\x1b[31m[Error] ' + data + '\x1b[0m');
      }
    });

    // Global shortcut keys for PTY control
    window.addEventListener('keydown', handleKeydown);

    return () => {
      window.removeEventListener('keydown', handleKeydown);
      EventsOff('pty-output');
      EventsOff('claude-done');
      EventsOff('claude-stuck');
      EventsOff('claude-status');
      EventsOff('claude-error');
    };
  });

  async function loadInitialState() {
    try {
      settings = await getSettings();
      applyDarkMode(settings.dark_mode);
      const s = await getStatus();
      status = s || 'idle';
    } catch (e) {
      console.error('Failed to init:', e);
    }
  }

  function applyDarkMode(dark: boolean) {
    document.body.classList.toggle('light', !dark);
  }

  async function handleSend(text: string) {
    if (!text.trim()) return;
    inputText = '';
    status = 'responding';
    try {
      await sendMessage(text);
    } catch (e) {
      console.error('Send failed:', e);
      status = 'error';
    }
  }

  function handleContinue() {
    stuckVisible = false;
    handleSend('please continue');
  }

  function handleDismiss() {
    stuckVisible = false;
    if (status === 'stuck') status = 'idle';
  }

  async function handleToggleMode() {
    if (!settings) return;
    const nextMode = settings.active_mode === 'custom_provider' ? 'normal' : 'custom_provider';
    const customProviderMode = nextMode === 'custom_provider';
    const newSettings = {
      ...settings,
      active_mode: nextMode,
      custom_provider_mode: customProviderMode,
    };
    try {
      await updateSettings(newSettings);
      settings = newSettings;
      terminal?.clear();
      await restartClaude();
    } catch (e) {
      console.error('Toggle custom provider mode failed:', e);
      status = 'error';
    }
  }

  async function handleRestart() {
    try {
      status = 'idle';
      stuckVisible = false;
      terminal?.clear();
      await restartClaude();
    } catch (e) {
      console.error('Restart failed:', e);
      status = 'error';
    }
  }

  async function handleClearHistory() {
    try {
      await clearHistory();
      terminal?.clear();
      stuckVisible = false;
      status = 'idle';
    } catch (e) {
      console.error('Clear failed:', e);
    }
  }

  async function handleNewSession() {
    try {
      await newSession('');
      terminal?.clear();
      stuckVisible = false;
      status = 'idle';
    } catch (e) {
      console.error('New session failed:', e);
    }
  }

  async function handleSaveSettings(newSettings: Settings) {
    try {
      await updateSettings(newSettings);
      settings = newSettings;
      applyDarkMode(newSettings.dark_mode);
      settingsOpen = false;
    } catch (e) {
      console.error('Save settings failed:', e);
    }
  }

  async function handleTerminalResize(event: CustomEvent<{cols: number; rows: number}>) {
    try {
      const { cols, rows } = event.detail;
      await resizeTerminal(cols, rows);
    } catch {}
  }

  function handleShiftTab() { sendPTYKey(Keys.shiftTab); }
  function handleCtrlO() { sendPTYKey(Keys.ctrlO); }
  function handleUp() { sendPTYKey(Keys.up); }
  function handleDown() { sendPTYKey(Keys.down); }

  function handleKeydown(event: KeyboardEvent) {
    if (event.altKey) {
      switch (event.key) {
        case '1':
          event.preventDefault();
          handleShiftTab();
          break;
        case '2':
          event.preventDefault();
          handleCtrlO();
          break;
        case '3':
          event.preventDefault();
          handleUp();
          break;
        case '4':
          event.preventDefault();
          handleDown();
          break;
      }
    }
  }
</script>

<StatusBar
  {status}
  {activeMode}
  on:restart={handleRestart}
  on:toggleMode={handleToggleMode}
/>
<TerminalView bind:this={terminal} on:resize={handleTerminalResize} />
<StuckBanner visible={stuckVisible} onContinue={handleContinue} onDismiss={handleDismiss} />
<InputBar
  bind:value={inputText}
  disabled={status === 'responding'}
  on:send={(e) => handleSend(e.detail)}
  on:settings={() => (settingsOpen = true)}
/>

{#if settingsOpen}
  <SettingsPanel {settings} onSave={handleSaveSettings} onClose={() => (settingsOpen = false)} />
{/if}
