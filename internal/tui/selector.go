package tui

import (
	"fmt"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/darkLord19/wtx/internal/git"
	"github.com/darkLord19/wtx/internal/metadata"
)

type model struct {
	list      list.Model
	items     []WorktreeItem
	choice    *WorktreeItem
	quitting  bool
	gitMgr    *git.Manager
	metaStore *metadata.Store
}

// NewSelector creates a new TUI selector
func NewSelector(gitMgr *git.Manager, metaStore *metadata.Store) (*model, error) {
	// Load worktrees
	worktrees, err := gitMgr.List()
	if err != nil {
		return nil, fmt.Errorf("failed to list worktrees: %w", err)
	}

	// Build items
	items := make([]list.Item, 0, len(worktrees))
	wtItems := make([]WorktreeItem, 0, len(worktrees))

	for _, wt := range worktrees {
		status, _ := gitMgr.GetStatus(wt.Path)

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

	// Create list
	delegate := list.NewDefaultDelegate()
	l := list.New(items, delegate, 0, 0)
	l.Title = "Workspace Manager"
	l.Styles.Title = titleStyle
	l.SetFilteringEnabled(true)

	return &model{
		list:      l,
		items:     wtItems,
		gitMgr:    gitMgr,
		metaStore: metaStore,
	}, nil
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.list.SetWidth(msg.Width)
		m.list.SetHeight(msg.Height - 4)
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q", "esc":
			m.quitting = true
			return m, tea.Quit

		case "enter":
			i, ok := m.list.SelectedItem().(WorktreeItem)
			if ok {
				m.choice = &i
				return m, tea.Quit
			}
		}
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m model) View() string {
	if m.quitting {
		return ""
	}

	helpText := helpStyle.Render("\nPress enter to open • q/esc to quit")
	return lipgloss.JoinVertical(
		lipgloss.Left,
		m.list.View(),
		helpText,
	)
}

// Run starts the TUI and returns the selected worktree
func Run(gitMgr *git.Manager, metaStore *metadata.Store) (*WorktreeItem, error) {
	m, err := NewSelector(gitMgr, metaStore)
	if err != nil {
		return nil, err
	}

	p := tea.NewProgram(m, tea.WithAltScreen())

	finalModel, err := p.Run()
	if err != nil {
		return nil, err
	}

	if m, ok := finalModel.(model); ok {
		return m.choice, nil
	}

	return nil, nil
}
