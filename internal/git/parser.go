package git

import (
	"fmt"
	"path/filepath"
	"strings"
)

// parseWorktreeList parses the output of git worktree list --porcelain
// Improved version with cleaner structure
func parseWorktreeList(output string) ([]Worktree, error) {
	output = strings.TrimSpace(output)
	if output == "" {
		return []Worktree{}, nil
	}

	// Split by double newline to get worktree blocks
	blocks := strings.Split(output, "\n\n")
	worktrees := make([]Worktree, 0, len(blocks))

	for _, block := range blocks {
		if block == "" {
			continue
		}

		wt := parseWorktreeBlock(block)
		if wt.Path != "" {
			worktrees = append(worktrees, wt)
		}
	}

	// Mark the first worktree as main if no bare repo was found
	if len(worktrees) > 0 && !worktrees[0].IsMain {
		worktrees[0].IsMain = true
	}

	return worktrees, nil
}

// parseWorktreeBlock parses a single worktree block
func parseWorktreeBlock(block string) Worktree {
	wt := Worktree{}
	lines := strings.Split(block, "\n")

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		// Handle special case: "bare" on its own line
		if line == "bare" {
			wt.IsMain = true
			continue
		}

		// Split by first space only
		parts := strings.SplitN(line, " ", 2)
		if len(parts) != 2 {
			continue
		}

		key := parts[0]
		value := parts[1]

		switch key {
		case "worktree":
			wt.Path = value
			wt.Name = filepath.Base(value)
		case "HEAD":
			wt.Head = value
		case "branch":
			wt.Branch = strings.TrimPrefix(value, "refs/heads/")
		}
	}

	return wt
}
