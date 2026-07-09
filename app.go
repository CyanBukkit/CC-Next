package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"regexp"
	"strings"
	"time"

	"ccnext/internal/claude"
	"ccnext/internal/config"
	"ccnext/internal/message"
	"ccnext/internal/stuck"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// App struct
type App struct {
	ctx            context.Context
	cfg            config.Settings
	cfgStore       *config.Store
	msgStore       *message.Store
	claude         *claude.Manager
	router         *message.Router
	detector       *stuck.Detector
	recovery       *stuck.RecoveryHandler
	currentSession *message.Session
	responding     bool
	skipDetection  bool // slash commands skip stuck/passphrase detection
}

// NewApp creates a new App application struct.
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts.
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx

	// Load settings
	a.cfgStore = config.NewStore()
	settings, err := a.cfgStore.Load()
	if err != nil {
		settings = config.DefaultSettings()
	}
	a.cfg = settings

	// Initialize database
	store, err := message.NewStore()
	if err != nil {
		log.Printf("database: %v", err)
	}
	a.msgStore = store

	// Resume or create session
	sessions, _ := store.ListSessions()
	var sess *message.Session
	if len(sessions) == 0 {
		sess, _ = store.CreateSession("New Chat")
	} else {
		sess = sessions[0]
	}
	a.currentSession = sess

	if sess != nil {
		a.router = message.NewRouter(a.msgStore, sess.ID)
	} else {
		a.router = message.NewRouter(a.msgStore, "")
	}

	// Setup stuck detection
	action := stuck.ActionAlertUser
	if settings.AutoContinueMode {
		action = stuck.ActionAutoContinue
	}
	a.recovery = stuck.NewRecoveryHandler(
		action,
		func() { a.handleAutoContinue() },
		func() { a.handleAlertUser() },
	)
	a.detector = stuck.NewDetector(
		time.Duration(settings.StuckTimeoutSeconds)*time.Second,
		func() { a.recovery.Handle() },
	)

	// CC Switch handles API config globally — just check if claude is available
	if !claude.CheckAvailability() {
		log.Printf("claude CLI not found — echo mode")
		return
	}

	// Start PTY session in configured working directory
	a.claude = claude.NewManager(
		settings.WorkDir,
		settings.CustomProviderMode,
		settings.ProviderBaseURL,
		settings.ProviderAuthToken,
		settings.ProviderModel,
	)
	if err := a.claude.Start(); err != nil {
		log.Printf("PTY start failed: %v — falling back to echo mode", err)
		a.claude = nil
		return
	}

	// Start background PTY output reader
	go a.ptyOutputLoop()

	log.Printf("CCNext ready — Claude Code via PTY")
}

// shutdown is called when the app is shutting down.
func (a *App) shutdown(ctx context.Context) {
	if a.detector != nil {
		a.detector.Stop()
	}
	if a.claude != nil {
		a.claude.Stop()
	}
	if a.msgStore != nil {
		a.msgStore.Close()
	}
}

// ptyOutputLoop continuously reads raw PTY output:
// - Sends raw bytes to frontend xterm.js via "pty-output" events
// - Also feeds cleaned text to passphrase/stuck detector
func (a *App) ptyOutputLoop() {
	reader := a.claude.OutputReader()
	if reader == nil {
		return
	}

	buf := make([]byte, 4096)
	for {
		n, err := reader.Read(buf)
		if err != nil {
			if err == io.EOF {
				return
			}
			time.Sleep(100 * time.Millisecond)
			continue
		}
		if n == 0 {
			continue
		}

		raw := string(buf[:n])

		// Send raw bytes to frontend for xterm.js rendering
		runtime.EventsEmit(a.ctx, "pty-output", raw)

		// Detect agent progress to keep timeout alive even during long silence.
		// Claude Code shows Running/N agents, tool uses, token counts, etc.
		clean := claude.StripANSI(raw)
		if clean != "" && hasAgentProgress(clean) {
			a.detector.OnTokenReceived()
		}

		// Feed CLEANED text to passphrase/stuck detector (skip for slash commands)
		if !a.skipDetection && clean != "" {
			a.detector.OnTokenReceived()
			_, found := a.router.FeedClaudeOutput(clean)
			if found {
				a.detector.OnResponseEnd()
				runtime.EventsEmit(a.ctx, "claude-done", true)
			}
		}
	}
}

// --- User-facing bindings ---

// SendMessage sends a user message to Claude via the PTY.
func (a *App) SendMessage(text string) error {
	if a.router == nil {
		return fmt.Errorf("router not initialised")
	}
	if a.responding {
		return fmt.Errorf("already responding")
	}

	a.router.AddToHistory(message.RoleUser, text)

	// Slash commands (e.g. /model, /doctor) are sent verbatim — no hidden instruction.
	var isSlashCmd bool
	var enhanced string
	if strings.HasPrefix(text, "/") {
		enhanced = text
		isSlashCmd = true
		a.skipDetection = true
		a.router.ClearMonitor()
	} else {
		enhanced = a.router.PrepareMessage(text, a.cfg.PassphraseTemplate)
		a.skipDetection = false
	}

	// Echo fallback
	if a.claude == nil {
		go a.simulateEchoResponse(text)
		return nil
	}

	a.responding = true
	if !isSlashCmd {
		a.detector.OnResponseStart()
	}

	go func() {
		defer func() { a.responding = false }()
		if isSlashCmd {
			defer a.detector.OnResponseEnd()
		}

		// Send the message via PTY
		if err := a.claude.SendInput(enhanced); err != nil {
			log.Printf("PTY send failed: %v", err)
			runtime.EventsEmit(a.ctx, "claude-error", err.Error())
			if !isSlashCmd {
				a.detector.OnResponseEnd()
			}
			return
		}

		// PTY output is handled by ptyOutputLoop goroutine
		// The done/stuck detection happens there
	}()

	return nil
}

// GetSettings returns current settings.
func (a *App) GetSettings() config.Settings {
	return a.cfg
}

// UpdateSettings updates and persists settings.
func (a *App) UpdateSettings(s config.Settings) error {
	// Keep active_mode and custom_provider_mode in sync.
	if s.ActiveMode == "custom_provider" {
		s.CustomProviderMode = true
	} else {
		s.CustomProviderMode = false
		if s.ActiveMode == "" {
			s.ActiveMode = "normal"
		}
	}

	a.cfg = s

	if a.recovery != nil {
		action := stuck.ActionAlertUser
		if s.AutoContinueMode {
			action = stuck.ActionAutoContinue
		}
		a.recovery.SetAction(action)
	}

	a.detector.Stop()
	a.detector = stuck.NewDetector(
		time.Duration(s.StuckTimeoutSeconds)*time.Second,
		func() { a.recovery.Handle() },
	)

	return a.cfgStore.Save(s)
}

// GetStatus returns Claude process status.
func (a *App) GetStatus() string {
	if a.claude != nil && a.claude.IsRunning() {
		return "running"
	}
	if a.claude != nil {
		return "stopped"
	}
	return "echo"
}

// GetActiveMode returns the currently selected Claude Code mode.
func (a *App) GetActiveMode() string {
	return a.cfg.ActiveMode
}

// SetActiveMode switches the mode and sends Shift+Tab to cycle the Claude Code TUI.
func (a *App) SetActiveMode(mode string) error {
	a.cfg.ActiveMode = mode
	_ = a.cfgStore.Save(a.cfg)
	if a.claude == nil || !a.claude.IsRunning() {
		return nil
	}
	order := []string{"normal", "plan", "explore", "ask", "build"}
	targetIdx := -1
	for i, m := range order {
		if m == mode {
			targetIdx = i
			break
		}
	}
	if targetIdx < 0 {
		return fmt.Errorf("unknown mode: %s", mode)
	}
	for i := 0; i < len(order); i++ {
		a.claude.SendKey([]byte("\x1b[Z"))
		time.Sleep(120 * time.Millisecond)
	}
	for i := 0; i < targetIdx; i++ {
		a.claude.SendKey([]byte("\x1b[Z"))
		time.Sleep(120 * time.Millisecond)
	}
	return nil
}

// ResizePTY resizes the PTY terminal dimensions.
func (a *App) ResizePTY(cols, rows int) {
	if a.claude != nil {
		a.claude.Resize(cols, rows)
	}
}

// SendPTYKey sends a key sequence to the PTY. seq is the raw byte sequence.
func (a *App) SendPTYKey(seq string) {
	if a.claude == nil {
		return
	}
	a.claude.SendKey([]byte(seq))
}

// RestartClaude restarts the PTY session.
func (a *App) RestartClaude() error {
	if a.claude != nil {
		a.claude.Stop()
	}

	a.claude = claude.NewManager(
		a.cfg.WorkDir,
		a.cfg.CustomProviderMode,
		a.cfg.ProviderBaseURL,
		a.cfg.ProviderAuthToken,
		a.cfg.ProviderModel,
	)
	if err := a.claude.Start(); err != nil {
		a.claude = nil
		return err
	}

	go a.ptyOutputLoop()
	return nil
}

// SetWorkDir changes the working directory and restarts Claude.
func (a *App) SetWorkDir(dir string) error {
	a.cfg.WorkDir = dir
	_ = a.cfgStore.Save(a.cfg)
	return a.RestartClaude()
}

// GetWorkDir returns the current working directory.
func (a *App) GetWorkDir() string {
	return a.cfg.WorkDir
}

// PickWorkDir opens a folder picker dialog, sets work dir, and restarts Claude.
func (a *App) PickWorkDir() (string, error) {
	dir, err := runtime.OpenDirectoryDialog(a.ctx, runtime.OpenDialogOptions{
		Title: "选择 Claude 工作目录",
	})
	if err != nil {
		return "", err
	}
	if dir == "" {
		return a.cfg.WorkDir, nil // user cancelled
	}
	if err := a.SetWorkDir(dir); err != nil {
		return "", err
	}
	return dir, nil
}

// GetHistory returns message history for the current session.
func (a *App) GetHistory() []message.Message {
	if a.router == nil {
		return nil
	}
	return a.router.GetHistory()
}

// ClearHistory clears the current conversation and starts a new session.
func (a *App) ClearHistory() (string, error) {
	if a.msgStore == nil {
		return "", fmt.Errorf("store not initialised")
	}

	sess, err := a.msgStore.CreateSession("New Chat")
	if err != nil {
		return "", err
	}

	a.currentSession = sess
	a.router = message.NewRouter(a.msgStore, sess.ID)
	return sess.ID, nil
}

// --- Session management ---

type SessionInfo struct {
	ID        string `json:"id"`
	Title     string `json:"title"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

func (a *App) ListSessions() ([]SessionInfo, error) {
	if a.msgStore == nil {
		return nil, fmt.Errorf("store not initialised")
	}
	sessions, err := a.msgStore.ListSessions()
	if err != nil {
		return nil, err
	}
	result := make([]SessionInfo, 0, len(sessions))
	for _, s := range sessions {
		result = append(result, SessionInfo{
			ID:        s.ID,
			Title:     s.Title,
			CreatedAt: s.CreatedAt.Format(time.RFC3339),
			UpdatedAt: s.UpdatedAt.Format(time.RFC3339),
		})
	}
	return result, nil
}

func (a *App) SwitchSession(sessionID string) ([]message.Message, error) {
	if a.msgStore == nil {
		return nil, fmt.Errorf("store not initialised")
	}
	sess, err := a.msgStore.GetSession(sessionID)
	if err != nil {
		return nil, err
	}
	a.currentSession = sess
	a.router = message.NewRouter(a.msgStore, sessionID)
	return a.router.GetHistory(), nil
}

func (a *App) NewSession(title string) (SessionInfo, error) {
	if a.msgStore == nil {
		return SessionInfo{}, fmt.Errorf("store not initialised")
	}
	if title == "" {
		title = "New Chat"
	}
	sess, err := a.msgStore.CreateSession(title)
	if err != nil {
		return SessionInfo{}, err
	}
	a.currentSession = sess
	a.router = message.NewRouter(a.msgStore, sess.ID)
	return SessionInfo{
		ID:        sess.ID,
		Title:     sess.Title,
		CreatedAt: sess.CreatedAt.Format(time.RFC3339),
		UpdatedAt: sess.UpdatedAt.Format(time.RFC3339),
	}, nil
}

func (a *App) DeleteSession(sessionID string) error {
	if a.msgStore == nil {
		return fmt.Errorf("store not initialised")
	}
	if a.currentSession != nil && a.currentSession.ID == sessionID {
		sessions, _ := a.msgStore.ListSessions()
		for _, s := range sessions {
			if s.ID != sessionID {
				a.currentSession = s
				a.router = message.NewRouter(a.msgStore, s.ID)
				break
			}
		}
	}
	return a.msgStore.DeleteSession(sessionID)
}

func (a *App) GetCurrentSession() SessionInfo {
	if a.currentSession == nil {
		return SessionInfo{}
	}
	return SessionInfo{
		ID:        a.currentSession.ID,
		Title:     a.currentSession.Title,
		CreatedAt: a.currentSession.CreatedAt.Format(time.RFC3339),
		UpdatedAt: a.currentSession.UpdatedAt.Format(time.RFC3339),
	}
}

// --- Internal ---

func (a *App) simulateEchoResponse(text string) {
	a.responding = true
	defer func() { a.responding = false }()

	a.detector.OnResponseStart()
	defer a.detector.OnResponseEnd()

	reply := fmt.Sprintf("\r\n\x1b[32mEcho:\x1b[0m %s\r\n", text)
	runtime.EventsEmit(a.ctx, "pty-output", reply)

	// Check passphrase
	for _, r := range reply {
		a.detector.OnTokenReceived()
		a.router.FeedClaudeOutput(string(r))
		time.Sleep(5 * time.Millisecond)
	}
	a.router.Flush()
	a.router.AddToHistory(message.RoleClaude, reply)
	runtime.EventsEmit(a.ctx, "claude-done", a.router.Found())
}

func (a *App) handleAutoContinue() {
	runtime.EventsEmit(a.ctx, "claude-stuck", "auto-continue")
	if a.claude != nil {
		a.claude.SendInput("Please continue from where you left off.")
	}
}

func (a *App) handleAlertUser() {
	runtime.EventsEmit(a.ctx, "claude-stuck", "alert")
}

// agentProgressPatterns matches TUI indicators that Claude Code is still running agents.
var agentProgressPatterns = regexp.MustCompile(strings.Join([]string{
	`Running \d+ `,       // "Running 3 Explore agents…"
	`tool uses`,           // "37 tool uses"
	`[\d.]+[km] tokens`,  // "18.7k tokens", "148.4k tokens"
	`\(\d+m \d+s`,         // "(27m 56s"
	`\b(Searching|reading|Executing|Scanning|Running)\b`,
	`Running \d`,          // "Running 3"
	`(ctrl\+[bo])`,        // "(ctrl+o to expand)", "(ctrl+b to run in background)"
	`⎿`,                   // Agent status prefix
	`\bDone\b`,            // "Done"
}, "|"))

// hasAgentProgress returns true if the text contains indicators that Claude Code agents are still working.
func hasAgentProgress(s string) bool {
	return agentProgressPatterns.MatchString(s)
}
