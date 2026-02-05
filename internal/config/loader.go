package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

// Load reads configuration from disk or creates default
func Load() (*Config, error) {
	v := viper.New()

	// Get config directory
	configDir, err := os.UserConfigDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get config directory: %w", err)
	}

	wtxConfigDir := filepath.Join(configDir, "wtx")
	v.AddConfigPath(wtxConfigDir)
	v.SetConfigName("config")
	v.SetConfigType("json")

	// Set defaults
	cfg := Default()
	v.SetDefault("editor", cfg.Editor)
	v.SetDefault("reuse_window", cfg.ReuseWindow)
	v.SetDefault("worktree_dir", cfg.WorktreeDir)
	v.SetDefault("auto_start_dev", cfg.AutoStartDev)

	// Read config
	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config doesn't exist, create it
			if err := os.MkdirAll(wtxConfigDir, 0755); err != nil {
				return nil, fmt.Errorf("failed to create config directory: %w", err)
			}
			if err := v.SafeWriteConfig(); err != nil {
				return nil, fmt.Errorf("failed to write config: %w", err)
			}
		} else {
			return nil, fmt.Errorf("failed to read config: %w", err)
		}
	}

	var config Config
	if err := v.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return &config, nil
}

// Save writes configuration to disk
func (c *Config) Save() error {
	v := viper.New()

	configDir, err := os.UserConfigDir()
	if err != nil {
		return fmt.Errorf("failed to get config directory: %w", err)
	}

	wtxConfigDir := filepath.Join(configDir, "wtx")
	v.AddConfigPath(wtxConfigDir)
	v.SetConfigName("config")
	v.SetConfigType("json")

	// Set values
	v.Set("editor", c.Editor)
	v.Set("reuse_window", c.ReuseWindow)
	v.Set("worktree_dir", c.WorktreeDir)
	v.Set("auto_start_dev", c.AutoStartDev)
	v.Set("custom_commands", c.CustomCommands)

	return v.WriteConfig()
}
