package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDefault(t *testing.T) {
	cfg := Default()
	if cfg == nil {
		t.Fatal("Default() returned nil")
	}
	if cfg.WorktreeDir != "../worktrees" {
		t.Errorf("Default worktree dir = %s, want ../worktrees", cfg.WorktreeDir)
	}
}

func setupConfigDir(t *testing.T) string {
	t.Helper()
	tmpDir := t.TempDir()

	// Mock UserConfigDir by setting environment variables
	// os.UserConfigDir checks XDG_CONFIG_HOME on Unix
	os.Setenv("XDG_CONFIG_HOME", tmpDir)
	// On macOS it uses HOME/Library/Application Support, so we set HOME too just in case
	os.Setenv("HOME", tmpDir)

	return tmpDir
}

func TestLoadSave(t *testing.T) {
	// Save original env
	origXDG := os.Getenv("XDG_CONFIG_HOME")
	origHome := os.Getenv("HOME")
	defer func() {
		os.Setenv("XDG_CONFIG_HOME", origXDG)
		os.Setenv("HOME", origHome)
	}()

	tmpDir := setupConfigDir(t)
	expectedConfigDir := filepath.Join(tmpDir, "wtx")

	t.Run("First Run", func(t *testing.T) {
		res, err := LoadWithFirstRunCheck()
		if err != nil {
			t.Fatalf("LoadWithFirstRunCheck() error = %v", err)
		}
		if !res.IsFirstRun {
			t.Error("Expected IsFirstRun to be true")
		}
		if res.Config == nil {
			t.Fatal("Expected Config to be non-nil")
		}

		// Verify config dir was created
		if _, err := os.Stat(expectedConfigDir); os.IsNotExist(err) {
			t.Error("Config directory was not created")
		}
	})

	t.Run("Save and Load", func(t *testing.T) {
		cfg := Default()
		cfg.Editor = "nvim"
		cfg.WorktreeDir = "/custom/path"

		if err := cfg.Save(); err != nil {
			t.Fatalf("Save() error = %v", err)
		}

		// Verify file exists
		configFile := filepath.Join(expectedConfigDir, "config.json")
		if _, err := os.Stat(configFile); os.IsNotExist(err) {
			t.Error("Config file was not created")
		}

		// Load again
		res, err := LoadWithFirstRunCheck()
		if err != nil {
			t.Fatalf("LoadWithFirstRunCheck() error = %v", err)
		}
		if res.IsFirstRun {
			t.Error("Expected IsFirstRun to be false")
		}
		if res.Config.Editor != "nvim" {
			t.Errorf("Loaded Editor = %s, want nvim", res.Config.Editor)
		}
		if res.Config.WorktreeDir != "/custom/path" {
			t.Errorf("Loaded WorktreeDir = %s, want /custom/path", res.Config.WorktreeDir)
		}
	})

	t.Run("Load simple", func(t *testing.T) {
		cfg, err := Load()
		if err != nil {
			t.Fatalf("Load() error = %v", err)
		}
		if cfg.Editor != "nvim" {
			t.Errorf("Load() returned Editor = %s, want nvim", cfg.Editor)
		}
	})
}
