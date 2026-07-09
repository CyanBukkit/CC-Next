package claude

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"regexp"
	"sync"
	"strings"
	"time"

	"github.com/UserExistsError/conpty"
)

// Manager manages a persistent Claude Code CLI session via Windows ConPTY.
type Manager struct {
	cpty     *conpty.ConPty
	workDir  string
	lastCols int
	lastRows int
	running  bool
	mu       sync.Mutex

	customProviderMode bool
	providerBaseURL    string
	providerAuthToken  string
	providerModel      string
}

// NewManager creates a Manager.
func NewManager(workDir string, customProviderMode bool, providerBaseURL, providerAuthToken, providerModel string) *Manager {
	return &Manager{
		workDir:            workDir,
		customProviderMode: customProviderMode,
		providerBaseURL:    providerBaseURL,
		providerAuthToken:  providerAuthToken,
		providerModel:      providerModel,
	}
}

// Start launches claude in the configured working directory.
func (m *Manager) Start() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.running {
		return fmt.Errorf("claude pty already running")
	}

	// Temporarily switch the process working directory so claude inherits it.
	// Restore afterwards — CCNext should not be affected.
	var prevDir string
	if m.workDir != "" {
		prevDir, _ = os.Getwd()
		os.Chdir(m.workDir)
	}

	c, err := conpty.Start("claude", conpty.ConPtyEnv(m.buildEnv()))

	// Restore CCNext's working directory
	if prevDir != "" {
		os.Chdir(prevDir)
	}

	if err != nil {
		return fmt.Errorf("start conpty: %w", err)
	}

	// Apply last known terminal size or reasonable default
	cols, rows := m.lastCols, m.lastRows
	if cols == 0 {
		cols, rows = 140, 40
	}
	c.Resize(cols, rows)
	m.cpty = c
	m.running = true

	go func() {
		c.Wait(context.Background())
		m.mu.Lock()
		m.running = false
		m.mu.Unlock()
	}()

	// Wait for Claude to initialize, then resize+refresh to fix any layout drift
	time.Sleep(2500 * time.Millisecond)
	c.Resize(cols, rows)
	time.Sleep(300 * time.Millisecond)
	// Send Ctrl+L to force Claude Code to redraw the TUI
	c.Write([]byte{'\x0c'})

	return nil
}

// Stop terminates the ConPTY session.
func (m *Manager) Stop() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if !m.running || m.cpty == nil {
		return nil
	}
	m.running = false
	return m.cpty.Close()
}

// IsRunning reports whether the PTY session is active.
func (m *Manager) IsRunning() bool {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.running
}

// SendInput sends text + Enter to Claude Code via the PTY.
func (m *Manager) SendInput(text string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if !m.running || m.cpty == nil {
		return fmt.Errorf("pty not running")
	}

	if _, err := m.cpty.Write([]byte(text)); err != nil {
		return err
	}
	time.Sleep(20 * time.Millisecond)
	_, err := m.cpty.Write([]byte{'\r'})
	return err
}

// SendKey sends a raw key sequence to the PTY.
func (m *Manager) SendKey(seq []byte) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if !m.running || m.cpty == nil {
		return fmt.Errorf("pty not running")
	}
	_, err := m.cpty.Write(seq)
	return err
}

// OutputReader returns an io.Reader for reading PTY output.
func (m *Manager) OutputReader() io.Reader {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.cpty == nil {
		return nil
	}
	return &ptyReader{cpty: m.cpty}
}

// Resize changes the terminal dimensions and remembers them for future restarts.
func (m *Manager) Resize(cols, rows int) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.lastCols = cols
	m.lastRows = rows
	if m.cpty == nil {
		return fmt.Errorf("pty not running")
	}
	return m.cpty.Resize(cols, rows)
}

// CheckAvailability returns true if claude is on PATH.
func CheckAvailability() bool {
	_, err := exec.LookPath("claude")
	return err == nil
}

// buildEnv returns the environment variables for the Claude process.
// When custom provider mode is enabled, it injects the provider configuration.
func (m *Manager) buildEnv() []string {
	env := os.Environ()

	if !m.customProviderMode {
		// Ensure no leftover custom provider variables leak into the child process.
		return filterEnv(env, []string{
			"ANTHROPIC_BASE_URL",
			"ANTHROPIC_API_KEY",
			"ANTHROPIC_MODEL",
		})
	}

	if m.providerBaseURL != "" {
		env = append(env, fmt.Sprintf("ANTHROPIC_BASE_URL=%s", m.providerBaseURL))
	}
	if m.providerAuthToken != "" {
		env = append(env, fmt.Sprintf("ANTHROPIC_API_KEY=%s", m.providerAuthToken))
	}
	if m.providerModel != "" {
		env = append(env, fmt.Sprintf("ANTHROPIC_MODEL=%s", m.providerModel))
	}
	return env
}

// filterEnv returns a copy of env with the given keys removed.
func filterEnv(env []string, keys []string) []string {
	keySet := make(map[string]struct{}, len(keys))
	for _, k := range keys {
		keySet[k] = struct{}{}
	}
	out := make([]string, 0, len(env))
nextVar:
	for _, e := range env {
		for k := range keySet {
			if strings.HasPrefix(e, k+"=") {
				continue nextVar
			}
		}
		out = append(out, e)
	}
	return out
}

// ---------------------------------------------------------------------------
// PTY reader wrapper
// ---------------------------------------------------------------------------

type ptyReader struct{ cpty *conpty.ConPty }

func (r *ptyReader) Read(p []byte) (int, error) { return r.cpty.Read(p) }

// ---------------------------------------------------------------------------
// ANSI cleaner
// ---------------------------------------------------------------------------

var ansiRe = regexp.MustCompile(`\x1b\[[?>!]?[0-9;]*[a-zA-Z]|\x1b\].*?\x07|\r`)

// StripANSI removes ANSI escape sequences and carriage returns.
func StripANSI(s string) string {
	return ansiRe.ReplaceAllString(s, "")
}
