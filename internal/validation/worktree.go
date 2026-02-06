package validation

import (
	"fmt"
	"regexp"
	"strings"
)

var (
	// nameRegex allows letters, numbers, hyphens, underscores, and forward slashes
	nameRegex = regexp.MustCompile(`^[a-zA-Z0-9_/-]+$`)
)

// WorktreeValidator validates worktree-related inputs
type WorktreeValidator struct{}

// NewWorktreeValidator creates a new validator
func NewWorktreeValidator() *WorktreeValidator {
	return &WorktreeValidator{}
}

// ValidateName validates a worktree name
func (v *WorktreeValidator) ValidateName(name string) error {
	name = strings.TrimSpace(name)
	
	if name == "" {
		return fmt.Errorf("name cannot be empty")
	}
	
	if len(name) > 128 {
		return fmt.Errorf("name too long (max 128 characters)")
	}

	// Prevent problematic names
	if name == "." || name == ".." {
		return fmt.Errorf("invalid name: cannot use '.' or '..'")
	}
	
	if !nameRegex.MatchString(name) {
		return fmt.Errorf("name can only contain letters, numbers, hyphens, underscores, and forward slashes")
	}

	return nil
}

// ValidateBranchName validates a git branch name
func (v *WorktreeValidator) ValidateBranchName(branch string) error {
	branch = strings.TrimSpace(branch)
	
	if branch == "" {
		return nil // Empty is allowed (will use name)
	}
	
	if len(branch) > 255 {
		return fmt.Errorf("branch name too long (max 255 characters)")
	}
	
	// Git branch name restrictions
	forbidden := []string{"..", "~", "^", ":", "?", "*", "[", "\\", " "}
	for _, f := range forbidden {
		if strings.Contains(branch, f) {
			return fmt.Errorf("branch name cannot contain '%s'", f)
		}
	}
	
	if strings.HasPrefix(branch, "/") || strings.HasSuffix(branch, "/") {
		return fmt.Errorf("branch name cannot start or end with '/'")
	}
	
	if strings.HasSuffix(branch, ".lock") {
		return fmt.Errorf("branch name cannot end with '.lock'")
	}
	
	return nil
}

// ValidatePath validates a worktree directory path
func (v *WorktreeValidator) ValidatePath(path string) error {
	path = strings.TrimSpace(path)
	
	if path == "" {
		return fmt.Errorf("path cannot be empty")
	}
	
	if len(path) > 4096 {
		return fmt.Errorf("path too long (max 4096 characters)")
	}
	
	return nil
}

// SanitizeName sanitizes a worktree name
func (v *WorktreeValidator) SanitizeName(name string) string {
	name = strings.TrimSpace(name)
	// Replace spaces with hyphens
	name = strings.ReplaceAll(name, " ", "-")
	// Remove consecutive hyphens
	for strings.Contains(name, "--") {
		name = strings.ReplaceAll(name, "--", "-")
	}
	return name
}
