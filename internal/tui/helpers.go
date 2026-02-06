package tui

import (
	"github.com/charmbracelet/bubbles/list"
	"github.com/darkLord19/wtx/internal/git"
	"github.com/darkLord19/wtx/internal/metadata"
)

// LoadWorktreeItems loads worktrees and their statuses in parallel
// Returns list items and WorktreeItem slice for use in TUI models
func LoadWorktreeItems(gitMgr *git.Manager, metaStore *metadata.Store) ([]list.Item, []WorktreeItem, error) {
	worktrees, err := gitMgr.List()
	if err != nil {
		return nil, nil, err
	}

	// Fetch statuses in parallel
	statuses := gitMgr.GetStatuses(worktrees)

	items := make([]list.Item, 0, len(worktrees))
	wtItems := make([]WorktreeItem, 0, len(worktrees))

	for _, wt := range worktrees {
		status := statuses[wt.Path]
		var meta *metadata.WorktreeMetadata
		if m, ok := metaStore.Get(wt.Name); ok {
			meta = m
		}

		item := WorktreeItem{
			Name:     wt.Name,
			Path:     wt.Path,
			Branch:   wt.Branch,
			Status:   status,
			Metadata: meta,
			IsMain:   wt.IsMain,
		}

		items = append(items, item)
		wtItems = append(wtItems, item)
	}

	return items, wtItems, nil
}

// CreateListModel creates a configured list.Model with standard settings
func CreateListModel(items []list.Item, title string) list.Model {
	delegate := list.NewDefaultDelegate()
	l := list.New(items, delegate, 0, 0)
	l.Title = title
	l.Styles.Title = titleStyle
	l.SetFilteringEnabled(true)
	return l
}
