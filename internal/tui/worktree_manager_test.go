package tui

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
	"time"

	"github.com/darkLord19/wtx/internal/git"
	"github.com/darkLord19/wtx/internal/metadata"
)

func TestRefreshListPerformance(t *testing.T) {
	// 1. Setup
	rootDir, err := os.MkdirTemp("", "wtx-test-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(rootDir)

	repoPath := filepath.Join(rootDir, "repo")
	if err := os.MkdirAll(repoPath, 0755); err != nil {
		t.Fatal(err)
	}

	// Helper to run git commands
	runGit := func(args ...string) {
		cmd := exec.Command("git", args...)
		cmd.Dir = repoPath
		// Suppress output unless error
		if out, err := cmd.CombinedOutput(); err != nil {
			t.Fatalf("git %v failed: %v\n%s", args, err, out)
		}
	}

	// Initialize git repo
	runGit("init")
	runGit("config", "user.email", "test@example.com")
	runGit("config", "user.name", "Test User")
	runGit("commit", "--allow-empty", "-m", "initial commit")
	runGit("branch", "-m", "main")

	// Create repository object
	// We need to construct Repository manually as FindRepo might fail or look elsewhere
	repo := &git.Repository{
		Path:   repoPath,
		GitDir: filepath.Join(repoPath, ".git"),
	}
	mgr := git.NewManager(repo)

	store, err := metadata.Load(repoPath)
	if err != nil {
		t.Fatal(err)
	}

	// Create worktrees to make the operation slower
	// We create 5 worktrees
	for i := 0; i < 5; i++ {
		name := fmt.Sprintf("wt-%d", i)
		branch := fmt.Sprintf("branch-%d", i)
		if _, err := mgr.Add(name, branch, "main"); err != nil {
			t.Logf("Failed to create worktree %s: %v", name, err)
		}
	}

	model, err := NewWorktreeManagerModel(mgr, store)
	if err != nil {
		t.Fatalf("Failed to create model: %v", err)
	}

	// 2. Measure
	start := time.Now()
	_, cmd := model.refreshList()
	duration := time.Since(start)

	t.Logf("refreshList took %v", duration)

	// 3. Verify cmd
	// Before optimization, cmd is nil. After optimization, cmd should be non-nil.
	if cmd != nil {
		t.Logf("cmd is not nil (optimized)")

		// 4. Verify command execution
		msg := cmd()

		// Since worktreeListMsg is unexported, we can't assert type directly if we were outside package,
		// but we are in package tui, so we can.
		listMsg, ok := msg.(worktreeListMsg)
		if !ok {
			t.Fatalf("Expected worktreeListMsg, got %T", msg)
		}

		// 5 created + 1 main = 6
		if len(listMsg) != 6 {
			t.Errorf("Expected 6 items, got %d", len(listMsg))
		}
	} else {
		t.Logf("cmd is nil (blocking)")
	}
}
