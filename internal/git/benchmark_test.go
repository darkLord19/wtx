package git

import (
	"testing"
)

func BenchmarkGetStatusSequential(b *testing.B) {
	// Setup repo with 10 worktrees
	repoPath, cleanup, err := setupTestRepo(nil, 10)
	if err != nil {
		b.Fatalf("Failed to setup test repo: %v", err)
	}
	defer cleanup()

	repo := &Repository{Path: repoPath}
	mgr := NewManager(repo)

	worktrees, err := mgr.List()
	if err != nil {
		b.Fatalf("Failed to list worktrees: %v", err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// Mimic the sequential loop found in manager.go
		for _, wt := range worktrees {
			_, err := mgr.GetStatus(wt.Path)
			if err != nil {
				b.Fatalf("GetStatus failed: %v", err)
			}
		}
	}
}

func BenchmarkGetStatusesParallel(b *testing.B) {
	// Setup repo with 10 worktrees
	repoPath, cleanup, err := setupTestRepo(nil, 10)
	if err != nil {
		b.Fatalf("Failed to setup test repo: %v", err)
	}
	defer cleanup()

	repo := &Repository{Path: repoPath}
	mgr := NewManager(repo)

	worktrees, err := mgr.List()
	if err != nil {
		b.Fatalf("Failed to list worktrees: %v", err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		mgr.GetStatuses(worktrees)
	}
}
