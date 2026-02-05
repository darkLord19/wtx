package tui

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/darkLord19/wtx/internal/config"
	"github.com/darkLord19/wtx/internal/editor"
	"github.com/darkLord19/wtx/internal/git"
	"github.com/darkLord19/wtx/internal/metadata"
)

func setupTestRepo(t *testing.T) (string, *git.Manager, *metadata.Store) {
	t.Helper()

	// Create temp dir
	dir, err := os.MkdirTemp("", "wtx-test-*")
	if err != nil {
		t.Fatal(err)
	}

	// Initialize git repo
	if err := exec.Command("git", "init", dir).Run(); err != nil {
		t.Fatal(err)
	}

	// Configure git user for commits
	exec.Command("git", "-C", dir, "config", "user.email", "you@example.com").Run()
	exec.Command("git", "-C", dir, "config", "user.name", "Your Name").Run()

	// Commit initial file
	if err := os.WriteFile(filepath.Join(dir, "README.md"), []byte("test"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := exec.Command("git", "-C", dir, "add", ".").Run(); err != nil {
		t.Fatal(err)
	}
	if err := exec.Command("git", "-C", dir, "commit", "-m", "initial").Run(); err != nil {
		t.Fatal(err)
	}

	// Setup dependencies
	repo := &git.Repository{Path: dir, GitDir: filepath.Join(dir, ".git")}
	gitMgr := git.NewManager(repo)
	metaStore := metadata.NewStore(dir)

	return dir, gitMgr, metaStore
}

func TestRefreshList_Async(t *testing.T) {
	dir, gitMgr, metaStore := setupTestRepo(t)
	defer os.RemoveAll(dir)

	cfg := &config.Config{}
	edDetector := editor.NewDetector(cfg)

	m, err := NewManagerModel(gitMgr, metaStore, cfg, edDetector)
	if err != nil {
		t.Fatalf("Failed to create model: %v", err)
	}

	// Switch to Manage tab ("2")
	msgSwitch := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'2'}, Alt: false}
	m.Update(msgSwitch)

	// Simulate pressing 'r'
	msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'r'}, Alt: false}

	start := time.Now()
	_, cmd := m.Update(msg)
	duration := time.Since(start)

	// We expect cmd to be NON-nil (async refresh)
	if cmd == nil {
		t.Errorf("Expected non-nil command for async refresh, got nil")
	}

	t.Logf("Update took %v", duration)

	// Execute the command to verify it returns the correct message type
	if cmd != nil {
		msgFromCmd := cmd()

		listMsg, ok := msgFromCmd.(WorktreeListMsg)
		if !ok {
			t.Errorf("Expected WorktreeListMsg, got %T", msgFromCmd)
		} else {
			if listMsg.Err != nil {
				t.Errorf("Refresh command failed: %v", listMsg.Err)
			}
			if len(listMsg.Items) == 0 {
				t.Error("Expected items in list, got 0")
			}
			t.Logf("Got %d items", len(listMsg.Items))
		}

		// Feed it back to Update
		m.Update(msgFromCmd)
	}
}
