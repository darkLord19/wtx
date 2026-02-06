package metadata

import (
	"sort"
	"time"
)

// Touch updates the last opened time and increments open count
func (s *Store) Touch(name string) {
	if wt, exists := s.Worktrees[name]; exists {
		wt.LastOpened = time.Now()
		wt.OpenCount++
		s.UpdatedAt = time.Now()
	}
}

// GetFrequent returns the most frequently opened worktrees
func (s *Store) GetFrequent(limit int) []*WorktreeMetadata {
	type wtWithCount struct {
		meta  *WorktreeMetadata
		count int
	}

	items := make([]wtWithCount, 0, len(s.Worktrees))
	for _, wt := range s.Worktrees {
		items = append(items, wtWithCount{meta: wt, count: wt.OpenCount})
	}

	// Sort by open count descending
	sort.Slice(items, func(i, j int) bool {
		if items[i].count == items[j].count {
			// If counts are equal, sort by most recently opened
			return items[i].meta.LastOpened.After(items[j].meta.LastOpened)
		}
		return items[i].count > items[j].count
	})

	result := make([]*WorktreeMetadata, 0, limit)
	for i := 0; i < limit && i < len(items); i++ {
		result = append(result, items[i].meta)
	}
	return result
}

// GetRecent returns the most recently opened worktrees
func (s *Store) GetRecent(limit int) []*WorktreeMetadata {
	items := make([]*WorktreeMetadata, 0, len(s.Worktrees))
	for _, wt := range s.Worktrees {
		items = append(items, wt)
	}

	// Sort by last opened descending
	sort.Slice(items, func(i, j int) bool {
		return items[i].LastOpened.After(items[j].LastOpened)
	})

	result := make([]*WorktreeMetadata, 0, limit)
	for i := 0; i < limit && i < len(items); i++ {
		result = append(result, items[i])
	}
	return result
}

// GetByAge returns worktrees sorted by creation date
func (s *Store) GetByAge(newest bool) []*WorktreeMetadata {
	items := make([]*WorktreeMetadata, 0, len(s.Worktrees))
	for _, wt := range s.Worktrees {
		items = append(items, wt)
	}

	sort.Slice(items, func(i, j int) bool {
		if newest {
			return items[i].CreatedAt.After(items[j].CreatedAt)
		}
		return items[i].CreatedAt.Before(items[j].CreatedAt)
	})

	return items
}
