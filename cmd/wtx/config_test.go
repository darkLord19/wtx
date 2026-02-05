package main

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/darkLord19/wtx/internal/config"
)

func TestConfigSet(t *testing.T) {
	// Setup temp home for config
	tmpDir := t.TempDir()

	oldHome := os.Getenv("HOME")
	defer os.Setenv("HOME", oldHome)
	os.Setenv("HOME", tmpDir)

	// Setup config dir
	// On Linux, UserConfigDir is usually $HOME/.config
	configDir := filepath.Join(tmpDir, ".config", "wtx")
	if err := os.MkdirAll(configDir, 0755); err != nil {
		t.Fatal(err)
	}

	// Create empty config file
	f, err := os.Create(filepath.Join(configDir, "config.json"))
	if err != nil {
		t.Fatal(err)
	}
	f.WriteString("{}")
	f.Close()

	// Initialize global cfg to a clean default
	cfg = config.Default()

	tests := []struct {
		name    string
		args    []string
		verify  func(t *testing.T)
		wantErr bool
	}{
		{
			name: "Set Editor",
			args: []string{"editor", "vim"},
			verify: func(t *testing.T) {
				if cfg.Editor != "vim" {
					t.Errorf("expected editor 'vim', got '%s'", cfg.Editor)
				}
			},
			wantErr: false,
		},
		{
			name: "Set ReuseWindow True",
			args: []string{"reuse_window", "true"},
			verify: func(t *testing.T) {
				if !cfg.ReuseWindow {
					t.Error("expected reuse_window true")
				}
			},
			wantErr: false,
		},
		{
			name: "Set ReuseWindow False",
			args: []string{"reuse_window", "false"},
			verify: func(t *testing.T) {
				if cfg.ReuseWindow {
					t.Error("expected reuse_window false")
				}
			},
			wantErr: false,
		},
		{
			name: "Set WorktreeDir",
			args: []string{"worktree_dir", "../foo"},
			verify: func(t *testing.T) {
				if cfg.WorktreeDir != "../foo" {
					t.Errorf("expected worktree_dir '../foo', got '%s'", cfg.WorktreeDir)
				}
			},
			wantErr: false,
		},
		{
			name: "Set AutoStartDev",
			args: []string{"auto_start_dev", "true"},
			verify: func(t *testing.T) {
				if !cfg.AutoStartDev {
					t.Error("expected auto_start_dev true")
				}
			},
			wantErr: false,
		},
		{
			name: "Set Editor Custom Command",
			args: []string{"editor", "custom_command", "myide", "ide --open"},
			verify: func(t *testing.T) {
				if val, ok := cfg.CustomCommands["myide"]; !ok || val != "ide --open" {
					t.Errorf("expected custom command 'ide --open', got '%s'", val)
				}
			},
			wantErr: false,
		},
		{
			name: "Invalid ReuseWindow",
			args: []string{"reuse_window", "foo"},
			verify: func(t *testing.T) {},
			wantErr: true,
		},
		{
			name: "Invalid Key",
			args: []string{"unknown_key", "value"},
			verify: func(t *testing.T) {},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := configCmd.RunE(configCmd, tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("RunE() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && tt.verify != nil {
				tt.verify(t)
			}
		})
	}
}
