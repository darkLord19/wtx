package metadata

import "time"

// WorktreeMetadata stores information about a worktree
type WorktreeMetadata struct {
	Name       string    `json:"name"`
	Path       string    `json:"path"`
	Branch     string    `json:"branch"`
	CreatedAt  time.Time `json:"created_at"`
	LastOpened time.Time `json:"last_opened"`
	DevCommand string    `json:"dev_command,omitempty"`
	Ports      []int     `json:"ports,omitempty"`
}

// Store holds all worktree metadata for a repository
type Store struct {
	RepoPath  string                       `json:"repo_path"`
	Worktrees map[string]*WorktreeMetadata `json:"worktrees"`
	UpdatedAt time.Time                    `json:"updated_at"`
}

// NewStore creates a new metadata store
func NewStore(repoPath string) *Store {
	return &Store{
		RepoPath:  repoPath,
		Worktrees: make(map[string]*WorktreeMetadata),
		UpdatedAt: time.Now(),
	}
}

// Add registers a new worktree in the store
func (s *Store) Add(wt *WorktreeMetadata) {
	s.Worktrees[wt.Name] = wt
	s.UpdatedAt = time.Now()
}

// Remove deletes a worktree from the store
func (s *Store) Remove(name string) {
	delete(s.Worktrees, name)
	s.UpdatedAt = time.Now()
}

// Touch updates the last opened time for a worktree
func (s *Store) Touch(name string) {
	if wt, exists := s.Worktrees[name]; exists {
		wt.LastOpened = time.Now()
		s.UpdatedAt = time.Now()
	}
}

// Get retrieves metadata for a worktree
func (s *Store) Get(name string) (*WorktreeMetadata, bool) {
	wt, exists := s.Worktrees[name]
	return wt, exists
}

// GetStale returns worktrees not opened in the specified number of days
func (s *Store) GetStale(days int) []string {
	cutoff := time.Now().AddDate(0, 0, -days)
	var stale []string

	for name, wt := range s.Worktrees {
		if wt.LastOpened.Before(cutoff) {
			stale = append(stale, name)
		}
	}

	return stale
}
