package tui

import (
	"fmt"

	"github.com/darkLord19/wtx/internal/git"
	"github.com/darkLord19/wtx/internal/metadata"
)

// WorktreeItem represents a worktree in the TUI list
type WorktreeItem struct {
	Name     string
	Path     string
	Branch   string
	Status   *git.Status
	Metadata *metadata.WorktreeMetadata
	IsMain   bool
}

// Title returns the primary display text
func (w WorktreeItem) Title() string {
	return w.Name
}

// Description returns the secondary display text with status
func (w WorktreeItem) Description() string {
	desc := ""

	// Branch name
	if w.Branch != "" {
		desc += fmt.Sprintf("%s ", w.Branch)
	}

	// Status indicators
	if w.Status != nil {
		if w.Status.Clean {
			desc += cleanStyle.Render("● clean")
		} else {
			desc += dirtyStyle.Render("✗ dirty")
		}

		if w.Status.Ahead > 0 {
			desc += fmt.Sprintf(" ↑%d", w.Status.Ahead)
		}
		if w.Status.Behind > 0 {
			desc += fmt.Sprintf(" ↓%d", w.Status.Behind)
		}
	}

	// Ports
	if w.Metadata != nil && len(w.Metadata.Ports) > 0 {
		for _, port := range w.Metadata.Ports {
			desc += portStyle.Render(fmt.Sprintf(" :%d", port))
		}
	}

	return desc
}

// FilterValue returns the value used for filtering
func (w WorktreeItem) FilterValue() string {
	return w.Name
}
