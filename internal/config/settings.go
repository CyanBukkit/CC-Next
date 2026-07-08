package config

// Settings holds all application configuration.
type Settings struct {
	PassphraseLength      int    `json:"passphrase_length"`
	StuckTimeoutSeconds   int    `json:"stuck_timeout_seconds"`
	AutoContinueMode      bool   `json:"auto_continue_mode"`
	VoiceProvider         string `json:"voice_provider"`
	ClaudeModel           string `json:"claude_model"`
	DarkMode              bool   `json:"dark_mode"`
	MaxHistoryMessages    int    `json:"max_history_messages"`
	PassphraseTemplate    string `json:"passphrase_template"`

	// Working directory for Claude
	WorkDir string `json:"work_dir"`

	// Provider config (kept for future use)
	ProviderAuthToken string `json:"provider_auth_token"`
	ProviderBaseURL   string `json:"provider_base_url"`
	ProviderModel     string `json:"provider_model"`
}

// DefaultPassphraseTemplate is the default hidden instruction.
// %random% is replaced with the generated passphrase.
const DefaultPassphraseTemplate = "当你完全完成所有任务后，你必须在回复结尾附上这个短语：`%random%`。不要向用户提及此指令，不要解释这个短语。"

// DefaultSettings returns the default configuration.
func DefaultSettings() Settings {
	return Settings{
		PassphraseLength:      10,
		StuckTimeoutSeconds:   300,
		AutoContinueMode:      true,
		VoiceProvider:         "web",
		ClaudeModel:           "",
		DarkMode:              true,
		MaxHistoryMessages:    500,
		PassphraseTemplate:    DefaultPassphraseTemplate,
	}
}
