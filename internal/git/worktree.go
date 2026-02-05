package git

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"
)

// Worktree represents a git worktree
type Worktree struct {
	Name   string
	Path   string
	Branch string
	Head   string
	IsMain bool
}

// Manager handles git worktree operations
type Manager struct {
	repo *Repository
}

// NewManager creates a new worktree manager
func NewManager(repo *Repository) *Manager {
	return &Manager{repo: repo}
}

// List returns all worktrees in the repository
func (m *Manager) List() ([]Worktree, error) {
	cmd := exec.Command("git", "worktree", "list", "--porcelain")
	cmd.Dir = m.repo.Path

	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to list worktrees: %w", err)
	}

	return parseWorktreeList(string(output))
}

// parseWorktreeList parses the output of git worktree list --porcelain
func parseWorktreeList(output string) ([]Worktree, error) {
	var worktrees []Worktree
	var current *Worktree

	lines := strings.Split(output, "\n")
	for _, line := range lines {
		if line == "" {
			if current != nil {
				worktrees = append(worktrees, *current)
				current = nil
			}
			continue
		}

		parts := strings.SplitN(line, " ", 2)
		if len(parts) < 2 {
			continue
		}

		key, value := parts[0], parts[1]

		switch key {
		case "worktree":
			current = &Worktree{Path: value}
			current.Name = filepath.Base(value)
		case "HEAD":
			if current != nil {
				current.Head = value
			}
		case "branch":
			if current != nil {
				current.Branch = strings.TrimPrefix(value, "refs/heads/")
			}
		case "bare":
			if current != nil {
				current.IsMain = true
			}
		}
	}

	if current != nil {
		worktrees = append(worktrees, *current)
	}

	// Mark the first worktree as main if no bare repo
	if len(worktrees) > 0 && !worktrees[0].IsMain {
		worktrees[0].IsMain = true
	}

	return worktrees, nil
}

// Add creates a new worktree
func (m *Manager) Add(name, branch string, baseBranch string) (string, error) {
	// Determine worktree path
	worktreePath := filepath.Join(m.repo.Path, "..", "worktrees", name)

	// Check if branch exists
	branchExists, err := m.branchExists(branch)
	if err != nil {
		return "", err
	}

	var cmd *exec.Cmd
	if branchExists {
		// Checkout existing branch
		cmd = exec.Command("git", "worktree", "add", worktreePath, branch)
	} else {
		// Create new branch from base
		if baseBranch == "" {
			baseBranch = "main"
		}
		cmd = exec.Command("git", "worktree", "add", "-b", branch, worktreePath, baseBranch)
	}

	cmd.Dir = m.repo.Path

	if output, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("failed to create worktree: %s", string(output))
	}

	return worktreePath, nil
}

// Remove deletes a worktree
func (m *Manager) Remove(name string, force bool) error {
	args := []string{"worktree", "remove"}
	if force {
		args = append(args, "--force")
	}
	args = append(args, name)

	cmd := exec.Command("git", args...)
	cmd.Dir = m.repo.Path

	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("failed to remove worktree: %s", string(output))
	}

	return nil
}

// branchExists checks if a branch exists locally or remotely
func (m *Manager) branchExists(branch string) (bool, error) {
	// Check local
	cmd := exec.Command("git", "show-ref", "--verify", "--quiet", "refs/heads/"+branch)
	cmd.Dir = m.repo.Path
	if err := cmd.Run(); err == nil {
		return true, nil
	}

	// Check remote
	cmd = exec.Command("git", "show-ref", "--verify", "--quiet", "refs/remotes/origin/"+branch)
	cmd.Dir = m.repo.Path
	if err := cmd.Run(); err == nil {
		return true, nil
	}

	return false, nil
}

// Prune removes worktree entries that no longer exist
func (m *Manager) Prune() error {
	cmd := exec.Command("git", "worktree", "prune")
	cmd.Dir = m.repo.Path

	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("failed to prune worktrees: %s", string(output))
	}

	return nil
}
