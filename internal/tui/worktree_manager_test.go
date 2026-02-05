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

func setupBenchmarkRepo(t testing.TB, worktreeCount int) (string, func()) {
	// Create temp directory
	tmpDir, err := os.MkdirTemp("", "wtx-bench-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}

	cleanup := func() {
		os.RemoveAll(tmpDir)
	}

	// Initialize main repo
	mainRepo := filepath.Join(tmpDir, "main")
	if err := os.Mkdir(mainRepo, 0755); err != nil {
		cleanup()
		t.Fatalf("failed to create main repo dir: %v", err)
	}

	runGit := func(dir string, args ...string) {
		cmd := exec.Command("git", args...)
		cmd.Dir = dir
		if out, err := cmd.CombinedOutput(); err != nil {
			cleanup()
			t.Fatalf("git %v failed: %v\nOutput: %s", args, err, out)
		}
	}

	runGit(mainRepo, "init")
	runGit(mainRepo, "branch", "-M", "main")
	runGit(mainRepo, "config", "user.email", "bench@example.com")
	runGit(mainRepo, "config", "user.name", "Benchmark")

	// Create initial commit
	if err := os.WriteFile(filepath.Join(mainRepo, "README.md"), []byte("# Benchmark"), 0644); err != nil {
		cleanup()
		t.Fatalf("failed to write file: %v", err)
	}
	runGit(mainRepo, "add", "README.md")
	runGit(mainRepo, "commit", "-m", "Initial commit")

	// Create worktrees
	for i := 0; i < worktreeCount; i++ {
		branchName := fmt.Sprintf("branch-%d", i)
		runGit(mainRepo, "worktree", "add", "-b", branchName, fmt.Sprintf("../wt-%d", i), "main")
	}

	return mainRepo, cleanup
}

func TestWorktreeManagerInit(t *testing.T) {
	worktreeCount := 5
	repoPath, cleanup := setupBenchmarkRepo(t, worktreeCount)
	defer cleanup()

	repo, err := git.FindRepo(repoPath)
	if err != nil {
		t.Fatalf("failed to find repo: %v", err)
	}
	gitMgr := git.NewManager(repo)

	metaDir := filepath.Join(filepath.Dir(repoPath), ".wtx")
	os.Mkdir(metaDir, 0755)
	metaStore := metadata.NewStore(filepath.Join(metaDir, "metadata.json"))

	// Create manager
	m, err := NewWorktreeManagerModel(gitMgr, metaStore)
	if err != nil {
		t.Fatalf("failed to create manager: %v", err)
	}

	// Check if all worktrees are loaded (main + 5 created = 6)
	expectedCount := worktreeCount + 1
	if len(m.items) != expectedCount {
		t.Errorf("expected %d items, got %d", expectedCount, len(m.items))
	}
}

func BenchmarkWorktreeManagerInit(b *testing.B) {
	// Setup repo with worktrees
	worktreeCount := 10
	repoPath, cleanup := setupBenchmarkRepo(b, worktreeCount)
	defer cleanup()

	// Initialize dependencies
	repo, err := git.FindRepo(repoPath)
	if err != nil {
		b.Fatalf("failed to find repo: %v", err)
	}
	gitMgr := git.NewManager(repo)

	// Mock metadata store (empty is fine)
	metaDir := filepath.Join(filepath.Dir(repoPath), ".wtx")
	os.Mkdir(metaDir, 0755)
	metaStore := metadata.NewStore(filepath.Join(metaDir, "metadata.json"))

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// We are benchmarking the initialization logic specifically
		if _, err := NewWorktreeManagerModel(gitMgr, metaStore); err != nil {
			b.Fatalf("failed to create manager: %v", err)
		}
	}
}
