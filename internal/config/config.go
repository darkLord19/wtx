package config

// Config holds user configuration
type Config struct {
	Editor         string            `mapstructure:"editor"`
	ReuseWindow    bool              `mapstructure:"reuse_window"`
	WorktreeDir    string            `mapstructure:"worktree_dir"`
	AutoStartDev   bool              `mapstructure:"auto_start_dev"`
	CustomCommands map[string]string `mapstructure:"custom_commands"`
}

// Default returns the default configuration
func Default() *Config {
	return &Config{
		Editor:         "", // Auto-detect
		ReuseWindow:    true,
		WorktreeDir:    "../worktrees",
		AutoStartDev:   false,
		CustomCommands: make(map[string]string),
	}
}
