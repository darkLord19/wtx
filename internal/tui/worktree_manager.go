package tui

import (
	tea "github.com/charmbracelet/bubbletea"

	"github.com/darkLord19/wtx/internal/git"
	"github.com/darkLord19/wtx/internal/metadata"
)

// WorktreeAction represents the action to perform
type WorktreeAction int

const (
	ActionNone WorktreeAction = iota
	ActionCreate
	ActionDelete
	ActionPrune
)

// worktreeManagerModel is the TUI model for worktree management
type worktreeManagerModel struct {
	manageModel *ManageModel
	quitting    bool
	width       int
	height      int
}

// NewWorktreeManagerModel creates a new worktree manager TUI model
func NewWorktreeManagerModel(gitMgr *git.Manager, metaStore *metadata.Store) (*worktreeManagerModel, error) {
	manageModel, err := NewManageModel(gitMgr, metaStore)
	if err != nil {
		return nil, err
	}

	return &worktreeManagerModel{
		manageModel: manageModel,
	}, nil
}

func (m *worktreeManagerModel) Init() tea.Cmd {
	return nil
}

func (m *worktreeManagerModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.manageModel.Update(msg)
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			// Only quit if not in a submode, or if 'q' is pressed in list mode
			if !m.manageModel.IsInSubMode() {
				m.quitting = true
				return m, tea.Quit
			}
			// If in submode, let manageModel handle 'q' (e.g. prune mode) or ignore it
		}
	}

	newModel, cmd := m.manageModel.Update(msg)
	m.manageModel = newModel.(*ManageModel)
	return m, cmd
}

func (m *worktreeManagerModel) View() string {
	if m.quitting {
		return ""
	}

	// Delegate mostly to manageModel
	// We might want to add a top-level breadcrumb "wtx > Manage" here if standalone
	// But manageModel already adds "Manage" breadcrumb.

	return m.manageModel.View()
}

// RunWorktreeManager starts the worktree manager TUI
func RunWorktreeManager(gitMgr *git.Manager, metaStore *metadata.Store) error {
	m, err := NewWorktreeManagerModel(gitMgr, metaStore)
	if err != nil {
		return err
	}

	p := tea.NewProgram(m, tea.WithAltScreen())
	_, err = p.Run()
	return err
}
