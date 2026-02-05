package git

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

// setupTestRepo creates a temporary git repository with some commits and worktrees
// Returns the path to the repo, a cleanup function, and an error if any
func setupTestRepo(t *testing.T, numWorktrees int) (string, func(), error) {
	// Create temp dir
	tmpDir, err := os.MkdirTemp("", "wtx-test-*")
	if err != nil {
		return "", nil, fmt.Errorf("failed to create temp dir: %w", err)
	}

	repoDir := filepath.Join(tmpDir, "repo")
	if err := os.Mkdir(repoDir, 0755); err != nil {
		os.RemoveAll(tmpDir)
		return "", nil, fmt.Errorf("failed to create repo dir: %w", err)
	}

	// Helper to run git commands
	git := func(args ...string) error {
		cmd := exec.Command("git", args...)
		cmd.Dir = repoDir
		if output, err := cmd.CombinedOutput(); err != nil {
			return fmt.Errorf("git %v failed: %w\nOutput: %s", args, err, string(output))
		}
		return nil
	}

	// Init repo
	if err := git("init"); err != nil {
		return "", nil, err
	}
	// Force branch name to main (in case default is master)
	if err := git("symbolic-ref", "HEAD", "refs/heads/main"); err != nil {
		return "", nil, err
	}

	// Config user
	if err := git("config", "user.email", "test@example.com"); err != nil {
		return "", nil, err
	}
	if err := git("config", "user.name", "Test User"); err != nil {
		return "", nil, err
	}

	// Create initial commit
	if err := os.WriteFile(filepath.Join(repoDir, "README.md"), []byte("# Test Repo"), 0644); err != nil {
		return "", nil, err
	}
	if err := git("add", "README.md"); err != nil {
		return "", nil, err
	}
	if err := git("commit", "-m", "Initial commit"); err != nil {
		return "", nil, err
	}

	// Create worktrees
	for i := 0; i < numWorktrees; i++ {
		branchName := fmt.Sprintf("branch-%d", i)
		wtPath := filepath.Join(tmpDir, fmt.Sprintf("wt-%d", i))

		// Create branch and worktree
		if err := git("worktree", "add", "-b", branchName, wtPath, "main"); err != nil {
			return "", nil, err
		}
	}

	cleanup := func() {
		os.RemoveAll(tmpDir)
	}

	return repoDir, cleanup, nil
}

func TestGetStatus(t *testing.T) {
	repoPath, cleanup, err := setupTestRepo(t, 1)
	if err != nil {
		t.Fatalf("Failed to setup test repo: %v", err)
	}
	defer cleanup()

	repo := &Repository{Path: repoPath}
	mgr := NewManager(repo)

	worktrees, err := mgr.List()
	if err != nil {
		t.Fatalf("Failed to list worktrees: %v", err)
	}

	if len(worktrees) < 2 { // Main + 1 worktree
		t.Fatalf("Expected at least 2 worktrees, got %d", len(worktrees))
	}

	// Pick the secondary worktree (not main)
	var wt Worktree
	for _, w := range worktrees {
		if !w.IsMain {
			wt = w
			break
		}
	}

	// Test Clean Status
	status, err := mgr.GetStatus(wt.Path)
	if err != nil {
		t.Fatalf("GetStatus failed: %v", err)
	}
	if !status.Clean {
		t.Errorf("Expected clean status, got dirty")
	}

	// Test Dirty Status
	dirtyFile := filepath.Join(wt.Path, "dirty.txt")
	if err := os.WriteFile(dirtyFile, []byte("dirty"), 0644); err != nil {
		t.Fatalf("Failed to create dirty file: %v", err)
	}

	status, err = mgr.GetStatus(wt.Path)
	if err != nil {
		t.Fatalf("GetStatus failed: %v", err)
	}
	if status.Clean {
		t.Errorf("Expected dirty status, got clean")
	}
}

func TestGetStatuses(t *testing.T) {
	// Setup 5 worktrees
	repoPath, cleanup, err := setupTestRepo(t, 5)
	if err != nil {
		t.Fatalf("Failed to setup test repo: %v", err)
	}
	defer cleanup()

	repo := &Repository{Path: repoPath}
	mgr := NewManager(repo)

	worktrees, err := mgr.List()
	if err != nil {
		t.Fatalf("Failed to list worktrees: %v", err)
	}

	// Make one dirty
	dirtyWT := worktrees[1] // worktrees[0] is usually main/bare, let's pick 1
	if dirtyWT.IsMain {
		// If 1 is also main (unlikely given setup), search for a non-main one
		for _, wt := range worktrees {
			if !wt.IsMain {
				dirtyWT = wt
				break
			}
		}
	}

	if err := os.WriteFile(filepath.Join(dirtyWT.Path, "dirty.txt"), []byte("dirty"), 0644); err != nil {
		t.Fatalf("Failed to create dirty file: %v", err)
	}

	statuses := mgr.GetStatuses(worktrees)

	if len(statuses) != len(worktrees) {
		t.Errorf("Expected %d statuses, got %d", len(worktrees), len(statuses))
	}

	// Check dirty one
	s, ok := statuses[dirtyWT.Path]
	if !ok {
		t.Errorf("Missing status for dirty worktree")
	} else if s.Clean {
		t.Errorf("Expected dirty status for %s", dirtyWT.Path)
	}

	// Check others (should be clean)
	for _, wt := range worktrees {
		if wt.Path == dirtyWT.Path {
			continue
		}
		s, ok := statuses[wt.Path]
		if !ok {
			t.Errorf("Missing status for %s", wt.Path)
		} else if !s.Clean {
			t.Errorf("Expected clean status for %s", wt.Path)
		}
	}
}
