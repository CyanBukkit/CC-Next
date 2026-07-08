package message

import (
	"ccnext/internal/passphrase"
	"sync"
)

// Role identifies the author of a message.
type Role string

const (
	RoleUser   Role = "user"
	RoleClaude Role = "claude"
	RoleSystem Role = "system"
)

// Message represents a single entry in the conversation history.
type Message struct {
	Role    Role   `json:"role"`
	Content string `json:"content"`
}

// Router prepares outgoing prompts and filters incoming Claude output.
type Router struct {
	mu            sync.RWMutex
	store         *Store
	sessionID     string
	history       []Message
	monitor       *passphrase.Monitor
	currentPhrase string
}

// NewRouter creates a new Router with the given store and session.
func NewRouter(store *Store, sessionID string) *Router {
	r := &Router{
		store:     store,
		sessionID: sessionID,
		history:   make([]Message, 0),
	}

	if store != nil && sessionID != "" {
		msgs, err := store.LoadMessages(sessionID, 0)
		if err == nil {
			r.history = msgs
		}
	}

	return r
}

// SessionID returns the current session identifier.
func (r *Router) SessionID() string {
	return r.sessionID
}

// SetSession switches the active session and loads its history.
func (r *Router) SetSession(sessionID string) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.sessionID = sessionID
	r.history = make([]Message, 0)

	if r.store != nil && sessionID != "" {
		msgs, err := r.store.LoadMessages(sessionID, 0)
		if err == nil {
			r.history = msgs
		}
	}
}

// PrepareMessage injects the hidden instruction using the given template.
func (r *Router) PrepareMessage(userText string, template string) string {
	r.currentPhrase = passphrase.Generate(10)
	r.monitor = passphrase.NewMonitor(r.currentPhrase)
	return userText + "\n\n" + passphrase.BuildInstruction(template, r.currentPhrase)
}

// FeedClaudeOutput processes Claude's streaming output through the passphrase monitor.
func (r *Router) FeedClaudeOutput(chunk string) (filtered string, passphraseFound bool) {
	if r.monitor == nil {
		return chunk, false
	}
	return r.monitor.Feed(chunk)
}

// AddToHistory adds a message to conversation history and persists it.
func (r *Router) AddToHistory(role Role, content string) {
	r.mu.Lock()
	msg := Message{Role: role, Content: content}
	r.history = append(r.history, msg)
	store := r.store
	sessionID := r.sessionID
	r.mu.Unlock()

	if store != nil && sessionID != "" && role != RoleSystem {
		_ = store.SaveMessage(sessionID, Message{Role: role, Content: content})
	}
}

// GetHistory returns the in-memory conversation history.
func (r *Router) GetHistory() []Message {
	r.mu.RLock()
	defer r.mu.RUnlock()

	result := make([]Message, len(r.history))
	copy(result, r.history)
	return result
}

// CurrentPassphrase returns the active passphrase (for debugging).
func (r *Router) CurrentPassphrase() string {
	return r.currentPhrase
}

// Found returns true if the passphrase was detected in Claude's output.
func (r *Router) Found() bool {
	if r.monitor == nil {
		return false
	}
	return r.monitor.Found()
}

// ClearMonitor resets the passphrase monitor (for slash commands that skip detection).
func (r *Router) ClearMonitor() {
	r.monitor = nil
	r.currentPhrase = ""
}

// Flush drains any remaining buffered output from the passphrase monitor.
func (r *Router) Flush() string {
	if r.monitor == nil {
		return ""
	}
	return r.monitor.Flush()
}

// TrimHistory removes old messages while keeping the most recent 'keep' entries.
func (r *Router) TrimHistory(keep int) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if keep <= 0 || len(r.history) <= keep {
		return
	}

	r.history = r.history[len(r.history)-keep:]

	if r.store != nil && r.sessionID != "" {
		_ = r.store.DeleteOldMessages(r.sessionID, keep)
	}
}
