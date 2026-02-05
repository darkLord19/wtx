package metadata

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestNewStore(t *testing.T) {
	store := NewStore("/test/repo")

	if store.RepoPath != "/test/repo" {
		t.Errorf("Expected RepoPath to be '/test/repo', got '%s'", store.RepoPath)
	}

	if store.Worktrees == nil {
		t.Error("Expected Worktrees map to be initialized")
	}

	if len(store.Worktrees) != 0 {
		t.Errorf("Expected empty Worktrees map, got %d items", len(store.Worktrees))
	}
}

func TestAddWorktree(t *testing.T) {
	store := NewStore("/test/repo")

	wt := &WorktreeMetadata{
		Name:       "test-feature",
		Path:       "/test/path",
		Branch:     "feature/test",
		CreatedAt:  time.Now(),
		LastOpened: time.Now(),
	}

	store.Add(wt)

	if len(store.Worktrees) != 1 {
		t.Errorf("Expected 1 worktree, got %d", len(store.Worktrees))
	}

	retrieved, exists := store.Get("test-feature")
	if !exists {
		t.Error("Expected worktree to exist")
	}

	if retrieved.Name != "test-feature" {
		t.Errorf("Expected name 'test-feature', got '%s'", retrieved.Name)
	}
}

func TestRemoveWorktree(t *testing.T) {
	store := NewStore("/test/repo")

	wt := &WorktreeMetadata{
		Name:   "test-feature",
		Path:   "/test/path",
		Branch: "feature/test",
	}

	store.Add(wt)
	store.Remove("test-feature")

	if len(store.Worktrees) != 0 {
		t.Errorf("Expected 0 worktrees, got %d", len(store.Worktrees))
	}

	_, exists := store.Get("test-feature")
	if exists {
		t.Error("Expected worktree to not exist")
	}
}

func TestTouchWorktree(t *testing.T) {
	store := NewStore("/test/repo")

	wt := &WorktreeMetadata{
		Name:       "test-feature",
		Path:       "/test/path",
		Branch:     "feature/test",
		LastOpened: time.Now().Add(-24 * time.Hour),
	}

	store.Add(wt)
	originalTime := wt.LastOpened

	time.Sleep(10 * time.Millisecond)
	store.Touch("test-feature")

	retrieved, _ := store.Get("test-feature")
	if !retrieved.LastOpened.After(originalTime) {
		t.Error("Expected LastOpened to be updated")
	}
}

func TestGetStale(t *testing.T) {
	store := NewStore("/test/repo")

	// Add fresh worktree
	fresh := &WorktreeMetadata{
		Name:       "fresh",
		Path:       "/test/fresh",
		Branch:     "fresh",
		LastOpened: time.Now(),
	}
	store.Add(fresh)

	// Add stale worktree
	stale := &WorktreeMetadata{
		Name:       "stale",
		Path:       "/test/stale",
		Branch:     "stale",
		LastOpened: time.Now().Add(-40 * 24 * time.Hour),
	}
	store.Add(stale)

	staleWorktrees := store.GetStale(30)

	if len(staleWorktrees) != 1 {
		t.Errorf("Expected 1 stale worktree, got %d", len(staleWorktrees))
	}

	if staleWorktrees[0] != "stale" {
		t.Errorf("Expected 'stale', got '%s'", staleWorktrees[0])
	}
}

func TestSaveAndLoad(t *testing.T) {
	// Create temporary directory
	tmpDir := t.TempDir()
	gitDir := filepath.Join(tmpDir, ".git")
	if err := os.MkdirAll(gitDir, 0755); err != nil {
		t.Fatal(err)
	}

	// Create and save store
	store := NewStore(tmpDir)
	wt := &WorktreeMetadata{
		Name:       "test-feature",
		Path:       "/test/path",
		Branch:     "feature/test",
		CreatedAt:  time.Now(),
		LastOpened: time.Now(),
	}
	store.Add(wt)

	if err := store.Save(); err != nil {
		t.Fatalf("Failed to save store: %v", err)
	}

	// Load store
	loaded, err := Load(tmpDir)
	if err != nil {
		t.Fatalf("Failed to load store: %v", err)
	}

	if len(loaded.Worktrees) != 1 {
		t.Errorf("Expected 1 worktree in loaded store, got %d", len(loaded.Worktrees))
	}

	retrieved, exists := loaded.Get("test-feature")
	if !exists {
		t.Error("Expected worktree to exist in loaded store")
	}

	if retrieved.Name != "test-feature" {
		t.Errorf("Expected name 'test-feature', got '%s'", retrieved.Name)
	}
}

func TestLoadNonExistent(t *testing.T) {
	tmpDir := t.TempDir()
	gitDir := filepath.Join(tmpDir, ".git")
	if err := os.MkdirAll(gitDir, 0755); err != nil {
		t.Fatal(err)
	}

	// Load from non-existent file should create new store
	store, err := Load(tmpDir)
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if store == nil {
		t.Fatal("Expected store to be created")
	}

	if len(store.Worktrees) != 0 {
		t.Errorf("Expected empty store, got %d worktrees", len(store.Worktrees))
	}
}
