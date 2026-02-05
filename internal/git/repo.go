package git

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// Repository represents a git repository
type Repository struct {
	Path   string
	GitDir string
}

// FindRepo walks up directories to find a git repository
func FindRepo(startPath string) (*Repository, error) {
	path, err := filepath.Abs(startPath)
	if err != nil {
		return nil, err
	}

	for {
		gitDir := filepath.Join(path, ".git")
		if info, err := os.Stat(gitDir); err == nil {
			if info.IsDir() {
				return &Repository{
					Path:   path,
					GitDir: gitDir,
				}, nil
			}
		}

		parent := filepath.Dir(path)
		if parent == path {
			return nil, fmt.Errorf("not a git repository")
		}
		path = parent
	}
}

// GetRootPath returns the repository root using git command
func GetRootPath() (string, error) {
	cmd := exec.Command("git", "rev-parse", "--show-toplevel")
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("not in a git repository")
	}
	return strings.TrimSpace(string(output)), nil
}

// IsGitInstalled checks if git is available
func IsGitInstalled() bool {
	_, err := exec.LookPath("git")
	return err == nil
}
