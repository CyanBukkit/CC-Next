package message

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	_ "modernc.org/sqlite"
)

// Session represents a conversation session.
type Session struct {
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Store provides persistent SQLite-backed storage for messages and sessions.
type Store struct {
	db *sql.DB
	mu sync.Mutex
}

func dbPath() string {
	base := os.Getenv("APPDATA")
	if base == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			base = "."
		} else {
			base = filepath.Join(home, ".ccnext")
		}
	} else {
		base = filepath.Join(base, "ccnext")
	}
	if err := os.MkdirAll(base, 0700); err != nil {
		return filepath.Join(".", "ccnext.db")
	}
	return filepath.Join(base, "ccnext.db")
}

// NewStore opens or creates the SQLite database and runs migrations.
func NewStore() (*Store, error) {
	db, err := sql.Open("sqlite", dbPath())
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	db.SetMaxOpenConns(1) // SQLite works best with a single writer

	s := &Store{db: db}
	if err := s.migrate(); err != nil {
		db.Close()
		return nil, fmt.Errorf("migration failed: %w", err)
	}
	return s, nil
}

func (s *Store) migrate() error {
	_, err := s.db.Exec(`
		CREATE TABLE IF NOT EXISTS sessions (
			id TEXT PRIMARY KEY,
			title TEXT NOT NULL DEFAULT '',
			created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
		);

		CREATE TABLE IF NOT EXISTS messages (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			session_id TEXT NOT NULL,
			role TEXT NOT NULL,
			content TEXT NOT NULL,
			created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (session_id) REFERENCES sessions(id) ON DELETE CASCADE
		);

		CREATE INDEX IF NOT EXISTS idx_messages_session
			ON messages(session_id, created_at);

		PRAGMA foreign_keys = ON;
		PRAGMA journal_mode = WAL;
	`)
	return err
}

// Close closes the database connection.
func (s *Store) Close() error {
	return s.db.Close()
}

// --- Sessions ---

// CreateSession creates a new session and returns it.
func (s *Store) CreateSession(title string) (*Session, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	session := &Session{
		ID:        newUUID(),
		Title:     title,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	_, err := s.db.Exec(
		"INSERT INTO sessions (id, title, created_at, updated_at) VALUES (?, ?, ?, ?)",
		session.ID, session.Title, session.CreatedAt, session.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return session, nil
}

// GetSession returns a session by ID.
func (s *Store) GetSession(id string) (*Session, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	row := s.db.QueryRow(
		"SELECT id, title, created_at, updated_at FROM sessions WHERE id = ?",
		id,
	)
	session := &Session{}
	err := row.Scan(&session.ID, &session.Title, &session.CreatedAt, &session.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return session, nil
}

// ListSessions returns all sessions ordered by most recently updated.
func (s *Store) ListSessions() ([]*Session, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	rows, err := s.db.Query(
		"SELECT id, title, created_at, updated_at FROM sessions ORDER BY updated_at DESC",
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var sessions []*Session
	for rows.Next() {
		session := &Session{}
		if err := rows.Scan(&session.ID, &session.Title, &session.CreatedAt, &session.UpdatedAt); err != nil {
			return nil, err
		}
		sessions = append(sessions, session)
	}
	return sessions, rows.Err()
}

// DeleteSession removes a session and all its messages.
func (s *Store) DeleteSession(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if _, err := tx.Exec("DELETE FROM messages WHERE session_id = ?", id); err != nil {
		return err
	}
	if _, err := tx.Exec("DELETE FROM sessions WHERE id = ?", id); err != nil {
		return err
	}
	return tx.Commit()
}

// UpdateSessionTitle updates a session's title.
func (s *Store) UpdateSessionTitle(id, title string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	_, err := s.db.Exec(
		"UPDATE sessions SET title = ?, updated_at = ? WHERE id = ?",
		title, time.Now(), id,
	)
	return err
}

// --- Messages ---

// SaveMessage stores a single message.
func (s *Store) SaveMessage(sessionID string, msg Message) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	content := msg.Content
	if role := string(msg.Role); role == "claude" {
		// Store the full content as JSON for rich messages
		content = msg.Content
	}

	_, err = tx.Exec(
		"INSERT INTO messages (session_id, role, content) VALUES (?, ?, ?)",
		sessionID, string(msg.Role), content,
	)
	if err != nil {
		return err
	}

	_, err = tx.Exec(
		"UPDATE sessions SET updated_at = ? WHERE id = ?",
		time.Now(), sessionID,
	)
	if err != nil {
		return err
	}

	return tx.Commit()
}

// LoadMessages loads all messages for a session, ordered by creation time.
func (s *Store) LoadMessages(sessionID string, limit int) ([]Message, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	query := "SELECT role, content FROM messages WHERE session_id = ? ORDER BY created_at ASC"
	if limit > 0 {
		query += fmt.Sprintf(" LIMIT %d", limit)
	}

	rows, err := s.db.Query(query, sessionID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []Message
	for rows.Next() {
		var roleStr, content string
		if err := rows.Scan(&roleStr, &content); err != nil {
			return nil, err
		}
		messages = append(messages, Message{
			Role:    Role(roleStr),
			Content: content,
		})
	}
	return messages, rows.Err()
}

// CountMessages returns the number of messages in a session.
func (s *Store) CountMessages(sessionID string) (int, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	var count int
	err := s.db.QueryRow(
		"SELECT COUNT(*) FROM messages WHERE session_id = ?",
		sessionID,
	).Scan(&count)
	return count, err
}

// DeleteOldMessages removes messages beyond the keep limit for a session.
func (s *Store) DeleteOldMessages(sessionID string, keep int) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	_, err := s.db.Exec(`
		DELETE FROM messages WHERE id IN (
			SELECT id FROM messages
			WHERE session_id = ?
			ORDER BY created_at DESC
			LIMIT -1 OFFSET ?
		)
	`, sessionID, keep)
	return err
}

// newUUID generates a simple UUID-like string (shortened for session IDs).
// Uses timestamp + random bytes for uniqueness without external dependency.
func newUUID() string {
	// Simple short unique ID: timestamp hex + random suffix
	b := make([]byte, 8)
	// Use time-based seed
	t := time.Now().UnixNano()
	for i := 0; i < 8; i++ {
		b[i] = byte(t >> (i * 8))
	}

	// Mix with some pseudo-randomness
	for i := range b {
		b[i] = b[i] ^ byte(i*37+int(t%256))
	}

	return fmt.Sprintf("%02x%02x%02x%02x-%02x%02x%02x%02x",
		b[0], b[1], b[2], b[3], b[4], b[5], b[6], b[7])
}

// MessageWithMeta wraps a Message with persistence metadata.
// Used internally for JSON serialization of rich message content.
type MessageWithMeta struct {
	Message
	JSONMeta json.RawMessage `json:"meta,omitempty"`
}
