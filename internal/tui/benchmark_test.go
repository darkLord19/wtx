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

func setupBenchmarkRepo(b *testing.B, worktreeCount int) (*git.Manager, *metadata.Store, func()) {
	// Create root dir
	rootDir := b.TempDir()
	repoPath := filepath.Join(rootDir, "repo")
	if err := os.Mkdir(repoPath, 0755); err != nil {
		b.Fatalf("failed to create repo dir: %v", err)
	}

	// Initialize git repo
	cmd := exec.Command("git", "init", repoPath)
	if err := cmd.Run(); err != nil {
		b.Fatalf("failed to init git repo: %v", err)
	}

	// Configure git user
	configCmds := [][]string{
		{"config", "user.email", "bench@example.com"},
		{"config", "user.name", "Benchmark User"},
		{"checkout", "-b", "main"},
		{"commit", "--allow-empty", "-m", "Initial commit"},
	}

	for _, args := range configCmds {
		cmd := exec.Command("git", args...)
		cmd.Dir = repoPath
		if err := cmd.Run(); err != nil {
			b.Fatalf("failed to run git config: %v", err)
		}
	}

	repo := &git.Repository{
		Path:   repoPath,
		GitDir: filepath.Join(repoPath, ".git"),
	}
	manager := git.NewManager(repo)
	metaStore := metadata.NewStore(repoPath)

	for i := 0; i < worktreeCount; i++ {
		name := fmt.Sprintf("wt-%d", i)
		branch := fmt.Sprintf("branch-%d", i)

		_, err := manager.Add(name, branch, "main")
		if err != nil {
			b.Fatalf("failed to add worktree %s: %v", name, err)
		}

		meta := &metadata.WorktreeMetadata{
			Name:       name,
			Path:       filepath.Join(rootDir, "worktrees", name),
			Branch:     branch,
			CreatedAt:  time.Now(),
			LastOpened: time.Now(),
		}
		metaStore.Add(meta)
	}

	cleanup := func() {
		os.RemoveAll(rootDir)
	}

	return manager, metaStore, cleanup
}

func BenchmarkRefreshList(b *testing.B) {
	// Setup repo with 10 worktrees
	gitMgr, metaStore, cleanup := setupBenchmarkRepo(b, 10)
	defer cleanup()

	// Initialize model
	_, err := NewWorktreeManagerModel(gitMgr, metaStore)
	if err != nil {
		b.Fatalf("failed to create model: %v", err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// We are benchmarking the refreshList method logic
		// But refreshList is private, so we can send a key 'r' which triggers it,
		// or expose it, or just copy the logic.
		// Since we want to benchmark the specific bottleneck (sync GetStatus calls),
		// we can invoke the logic directly if possible or simulate the update loop.

		// Ideally we test the function causing the issue.
		// `NewWorktreeManagerModel` calls `gitMgr.List()` and then loops `gitMgr.GetStatus`.
		// `refreshList` does the same.

		// Let's create a new model each time to test the initialization cost which has the same loop
		_, err := NewWorktreeManagerModel(gitMgr, metaStore)
		if err != nil {
			b.Fatalf("iteration %d: failed to create model: %v", i, err)
		}
	}
}
