package metadata

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// Load reads metadata from disk, creating a new store if the file doesn't exist
func Load(repoPath string) (*Store, error) {
	path := metadataPath(repoPath)

	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return NewStore(repoPath), nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to read metadata: %w", err)
	}

	var store Store
	if err := json.Unmarshal(data, &store); err != nil {
		return nil, fmt.Errorf("failed to parse metadata: %w", err)
	}

	return &store, nil
}

// Save writes metadata to disk
func (s *Store) Save() error {
	s.UpdatedAt = time.Now()

	path := metadataPath(s.RepoPath)

	// Ensure directory exists
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return fmt.Errorf("failed to create metadata directory: %w", err)
	}

	data, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal metadata: %w", err)
	}

	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("failed to write metadata: %w", err)
	}

	return nil
}

// metadataPath returns the path to the metadata file
func metadataPath(repoPath string) string {
	gitDir := filepath.Join(repoPath, ".git")
	return filepath.Join(gitDir, "wtx-meta.json")
}
