package git

import (
	"fmt"
	"os/exec"
	"strconv"
	"strings"
	"sync"
)

// Status represents the git status of a worktree
type Status struct {
	Clean      bool
	Ahead      int
	Behind     int
	HasChanges bool
}

// GetStatus returns the git status for a worktree
func (m *Manager) GetStatus(worktreePath string) (*Status, error) {
	status := &Status{Clean: true}

	// Check for uncommitted changes
	cmd := exec.Command("git", "status", "--porcelain")
	cmd.Dir = worktreePath
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to get status: %w", err)
	}

	if len(strings.TrimSpace(string(output))) > 0 {
		status.Clean = false
		status.HasChanges = true
	}

	// Check ahead/behind
	cmd = exec.Command("git", "rev-list", "--left-right", "--count", "HEAD...@{upstream}")
	cmd.Dir = worktreePath
	output, err = cmd.Output()
	if err == nil {
		// Parse "X\tY" format
		parts := strings.Fields(string(output))
		if len(parts) == 2 {
			if ahead, err := strconv.Atoi(parts[0]); err == nil {
				status.Ahead = ahead
			}
			if behind, err := strconv.Atoi(parts[1]); err == nil {
				status.Behind = behind
			}
		}
	}
	// Ignore error if no upstream is set

	return status, nil
}

// GetStatuses returns the status for multiple worktrees concurrently
func (m *Manager) GetStatuses(worktrees []Worktree) map[string]*Status {
	results := make(map[string]*Status)
	var mu sync.Mutex
	var wg sync.WaitGroup

	// Limit concurrency to avoid resource exhaustion
	semaphore := make(chan struct{}, 20)

	for _, wt := range worktrees {
		wg.Add(1)
		go func(wt Worktree) {
			defer wg.Done()
			semaphore <- struct{}{}        // Acquire
			defer func() { <-semaphore }() // Release

			status, err := m.GetStatus(wt.Path)
			if err == nil {
				mu.Lock()
				results[wt.Path] = status
				mu.Unlock()
			}
		}(wt)
	}

	wg.Wait()
	return results
}

// IsClean checks if a worktree has no uncommitted changes
func (m *Manager) IsClean(worktreePath string) (bool, error) {
	status, err := m.GetStatus(worktreePath)
	if err != nil {
		return false, err
	}
	return status.Clean, nil
}
