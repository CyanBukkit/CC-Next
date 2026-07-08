<script lang="ts">
  import { onMount, onDestroy, createEventDispatcher } from 'svelte';
  import { Terminal } from '@xterm/xterm';
  import { FitAddon } from '@xterm/addon-fit';
  import '@xterm/xterm/css/xterm.css';

  const dispatch = createEventDispatcher();

  let container: HTMLDivElement;
  let term: Terminal;
  let fitAddon: FitAddon;
  let resizeTimer: ReturnType<typeof setTimeout>;

  onMount(() => {
    term = new Terminal({
      fontSize: 14,
      fontFamily: '"Cascadia Code", "Fira Code", "JetBrains Mono", Consolas, "Noto Sans Mono", monospace',
      theme: {
        background: '#0f0f1a',
        foreground: '#d4d4e8',
        cursor: '#5a7cff',
        cursorAccent: '#0f0f1a',
        selectionBackground: '#5a7cff55',
        selectionForeground: '#ffffff',
        black:   '#1a1a30',
        red:     '#ff6b6b',
        green:   '#51cf66',
        yellow:  '#ffd43b',
        blue:    '#5a7cff',
        magenta: '#cc5de8',
        cyan:    '#20c997',
        white:   '#d4d4e8',
        brightBlack:   '#3d3d5c',
        brightRed:     '#ff8787',
        brightGreen:   '#69db7c',
        brightYellow:  '#ffe066',
        brightBlue:    '#7c9bff',
        brightMagenta: '#da77f2',
        brightCyan:    '#38d9a9',
        brightWhite:   '#ffffff',
      },
      cursorBlink: true,
      cursorStyle: 'bar',
      allowProposedApi: true,
      scrollback: 10000,
      disableStdin: true,
      smoothScrollDuration: 50,
      convertEol: true,
      allowTransparency: false,
      cols: 100,
      rows: 30,
    });

    fitAddon = new FitAddon();
    term.loadAddon(fitAddon);
    term.open(container);

    // Initial fit
    requestAnimationFrame(() => {
      fitAddon.fit();
      notifySize();
    });

    // Resize → fit + notify Go backend
    const observer = new ResizeObserver(() => {
      clearTimeout(resizeTimer);
      resizeTimer = setTimeout(() => {
        try {
          fitAddon.fit();
          notifySize();
        } catch {}
      }, 100);
    });
    observer.observe(container);

    return () => {
      observer.disconnect();
      clearTimeout(resizeTimer);
    };
  });

  onDestroy(() => {
    if (term) term.dispose();
  });

  function notifySize() {
    dispatch('resize', {
      cols: term.cols,
      rows: term.rows,
    });
  }

  export function write(data: string) {
    term?.write(data);
  }

  export function writeln(data: string) {
    term?.writeln(data);
  }

  export function clear() {
    term?.clear();
  }

  export function refresh() {
    if (!term) return;
    term.clear();
    try { fitAddon?.fit(); } catch {}
    notifySize();
    // Send SIGWINCH-like resize escape to force PTY redraw
    dispatch('resize', { cols: term.cols, rows: term.rows });
  }

  export function getSize() {
    return term ? { cols: term.cols, rows: term.rows } : { cols: 0, rows: 0 };
  }
</script>

<div bind:this={container} class="terminal-container"></div>

<style>
  .terminal-container {
    width: 100%;
    height: 100%;
    min-height: 0;
    overflow: hidden;
    background: #0f0f1a;
    display: flex;
  }

  :global(.terminal-container .xterm) {
    flex: 1;
    min-height: 0;
    padding: 6px 8px;
  }

  :global(.terminal-container .xterm-viewport) {
    scrollbar-width: thin;
    scrollbar-color: #2a2a45 #0f0f1a;
  }

  :global(.terminal-container .xterm-viewport::-webkit-scrollbar) {
    width: 5px;
  }

  :global(.terminal-container .xterm-viewport::-webkit-scrollbar-track) {
    background: #0f0f1a;
  }

  :global(.terminal-container .xterm-viewport::-webkit-scrollbar-thumb) {
    background: #2a2a45;
    border-radius: 3px;
  }

  :global(.terminal-container .xterm-viewport::-webkit-scrollbar-thumb:hover) {
    background: #3d3d5c;
  }

  /* Hide xterm cursor since input is handled externally */
  :global(.terminal-container .xterm-cursor-layer) {
    display: none;
  }
</style>
