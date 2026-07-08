package config

import (
	"encoding/json"
	"os"
	"path/filepath"
	"runtime"
)

// Store persists application settings to disk.
type Store struct {
	path string
}

// NewStore creates a Store using the platform-specific settings path.
func NewStore() *Store {
	return &Store{path: settingsPath()}
}

// Path returns the active settings file path.
func (s *Store) Path() string {
	return s.path
}

// Load reads settings from disk. If the file does not exist, it returns defaults.
func (s *Store) Load() (Settings, error) {
	settings := DefaultSettings()

	data, err := os.ReadFile(s.path)
	if err != nil {
		if os.IsNotExist(err) {
			return settings, nil
		}
		return settings, err
	}

	if err := json.Unmarshal(data, &settings); err != nil {
		return settings, err
	}
	return settings, nil
}

// Save writes settings atomically to disk (temp file + rename).
func (s *Store) Save(settings Settings) error {
	if err := os.MkdirAll(filepath.Dir(s.path), 0755); err != nil {
		return err
	}

	data, err := json.MarshalIndent(settings, "", "  ")
	if err != nil {
		return err
	}
	data = append(data, '\n')

	tmpPath := s.path + ".tmp"
	if err := os.WriteFile(tmpPath, data, 0644); err != nil {
		return err
	}

	return os.Rename(tmpPath, s.path)
}

// settingsPath returns the platform-specific settings file path.
func settingsPath() string {
	if runtime.GOOS == "windows" {
		appdata := os.Getenv("APPDATA")
		if appdata != "" {
			return filepath.Join(appdata, "ccnext", "settings.json")
		}
	}

	home, err := os.UserHomeDir()
	if err != nil || home == "" {
		return filepath.Join(".", ".ccnext", "settings.json")
	}
	return filepath.Join(home, ".ccnext", "settings.json")
}
