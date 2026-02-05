package tui

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/darkLord19/wtx/internal/git"
	"github.com/darkLord19/wtx/internal/metadata"
)

func setupBenchmarkRepo(t *testing.B, worktreeCount int) (string, func()) {
	// Create temp dir
	dir, err := os.MkdirTemp("", "wtx-bench")
	if err != nil {
		t.Fatal(err)
	}

	repoDir := filepath.Join(dir, "repo")
	if err := os.Mkdir(repoDir, 0755); err != nil {
		t.Fatal(err)
	}

	// Init git repo
	runGit(t, repoDir, "init", "-b", "main")
	runGit(t, repoDir, "config", "user.name", "Test User")
	runGit(t, repoDir, "config", "user.email", "test@example.com")

	// Create initial commit
	if err := os.WriteFile(filepath.Join(repoDir, "README.md"), []byte("# Test Repo"), 0644); err != nil {
		t.Fatal(err)
	}
	runGit(t, repoDir, "add", ".")
	runGit(t, repoDir, "commit", "-m", "Initial commit")

	// Create worktrees
	worktreesDir := filepath.Join(dir, "worktrees")
	if err := os.Mkdir(worktreesDir, 0755); err != nil {
		t.Fatal(err)
	}

	for i := 0; i < worktreeCount; i++ {
		branchName := fmt.Sprintf("branch-%d", i)
		wtPath := filepath.Join(worktreesDir, fmt.Sprintf("wt-%d", i))
		runGit(t, repoDir, "worktree", "add", "-b", branchName, wtPath, "main")
	}

	return repoDir, func() {
		os.RemoveAll(dir)
	}
}

func runGit(t *testing.B, dir string, args ...string) {
	cmd := exec.Command("git", args...)
	cmd.Dir = dir
	if output, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("git command failed: %v\nOutput: %s", err, output)
	}
}

func BenchmarkWorktreeManagerInit(b *testing.B) {
	// Setup repo with 20 worktrees
	repoPath, cleanup := setupBenchmarkRepo(b, 20)
	defer cleanup()

	// Initialize dependencies
	repo, err := git.FindRepo(repoPath)
	if err != nil {
		b.Fatal(err)
	}
	gitMgr := git.NewManager(repo)
	metaStore := metadata.NewStore(repoPath)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := NewWorktreeManagerModel(gitMgr, metaStore)
		if err != nil {
			b.Fatal(err)
		}
	}
}
