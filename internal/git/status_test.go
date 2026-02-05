package git

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

func TestGetStatuses(t *testing.T) {
	// Setup repo
	rootDir := t.TempDir()
	repoDir := filepath.Join(rootDir, "repo")
	if err := os.Mkdir(repoDir, 0755); err != nil {
		t.Fatal(err)
	}

	runGit(t, repoDir, "init")
	runGit(t, repoDir, "config", "user.email", "test@example.com")
	runGit(t, repoDir, "config", "user.name", "Test")
	runGit(t, repoDir, "commit", "--allow-empty", "-m", "Initial commit")

	worktreesDir := filepath.Join(rootDir, "worktrees")
	if err := os.Mkdir(worktreesDir, 0755); err != nil {
		t.Fatal(err)
	}

	paths := []string{}

	// WT1: Clean
	runGit(t, repoDir, "branch", "wt1", "HEAD")
	wt1Path := filepath.Join(worktreesDir, "wt1")
	runGit(t, repoDir, "worktree", "add", wt1Path, "wt1")
	paths = append(paths, wt1Path)

	// WT2: Dirty
	runGit(t, repoDir, "branch", "wt2", "HEAD")
	wt2Path := filepath.Join(worktreesDir, "wt2")
	runGit(t, repoDir, "worktree", "add", wt2Path, "wt2")
	// Make it dirty by creating an untracked file.
	// Note: 'git status --porcelain' shows untracked files by default, making it dirty.
	if err := os.WriteFile(filepath.Join(wt2Path, "dirty.txt"), []byte("dirty"), 0644); err != nil {
		t.Fatal(err)
	}
	paths = append(paths, wt2Path)

	repo := &Repository{
		Path:   repoDir,
		GitDir: filepath.Join(repoDir, ".git"),
	}
	mgr := NewManager(repo)

	statuses := mgr.GetStatuses(paths)

	if len(statuses) != 2 {
		t.Errorf("Expected 2 statuses, got %d", len(statuses))
	}

	if s, ok := statuses[wt1Path]; !ok {
		t.Error("Missing status for wt1")
	} else if !s.Clean {
		t.Error("wt1 should be clean")
	}

	if s, ok := statuses[wt2Path]; !ok {
		t.Error("Missing status for wt2")
	} else if s.Clean {
		t.Error("wt2 should be dirty")
	}
}

func runGit(t *testing.T, dir string, args ...string) {
	cmd := exec.Command("git", args...)
	cmd.Dir = dir
	if output, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("git %v failed: %v\nOutput: %s", args, err, output)
	}
}
