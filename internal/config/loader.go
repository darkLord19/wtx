package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

// LoadResult contains the loaded config and whether it's a first run
type LoadResult struct {
	Config     *Config
	IsFirstRun bool
}

// Load reads configuration from disk or creates default
func Load() (*Config, error) {
	result, err := LoadWithFirstRunCheck()
	if err != nil {
		return nil, err
	}
	return result.Config, nil
}

// LoadWithFirstRunCheck reads configuration and indicates if this is first run
func LoadWithFirstRunCheck() (*LoadResult, error) {
	v := viper.New()
	isFirstRun := false

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
			// Config doesn't exist - this is first run
			isFirstRun = true
			if err := os.MkdirAll(wtxConfigDir, 0755); err != nil {
				return nil, fmt.Errorf("failed to create config directory: %w", err)
			}
			// Don't write config yet - let setup wizard handle it
		} else {
			return nil, fmt.Errorf("failed to read config: %w", err)
		}
	}

	var config Config
	if err := v.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return &LoadResult{
		Config:     &config,
		IsFirstRun: isFirstRun,
	}, nil
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
