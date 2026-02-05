package tui

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

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
	list         list.Model
	items        []WorktreeItem
	gitMgr       *git.Manager
	metaStore    *metadata.Store
	width        int
	height       int
	quitting     bool
	message      string
	messageStyle lipgloss.Style

	// Create form
	createMode  bool
	nameInput   textinput.Model
	branchInput textinput.Model
	baseInput   textinput.Model
	createFocus int // 0=name, 1=branch, 2=base

	// Delete confirmation
	deleteMode   bool
	deleteTarget *WorktreeItem
	forceDelete  bool

	// Prune mode
	pruneMode     bool
	staleItems    []WorktreeItem
	pruneCursor   int
	pruneSelected map[int]bool
	staleDays     int
}

// NewWorktreeManagerModel creates a new worktree manager TUI model
func NewWorktreeManagerModel(gitMgr *git.Manager, metaStore *metadata.Store) (*worktreeManagerModel, error) {
	// Load worktrees
	worktrees, err := gitMgr.List()
	if err != nil {
		return nil, fmt.Errorf("failed to list worktrees: %w", err)
	}

	// Fetch statuses in parallel
	statuses := gitMgr.GetStatuses(worktrees)

	// Build items
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

	// Create list
	delegate := list.NewDefaultDelegate()
	l := list.New(items, delegate, 0, 0)
	l.Title = "Worktree Manager"
	l.Styles.Title = titleStyle
	l.SetFilteringEnabled(true)

	// Create text inputs
	nameInput := textinput.New()
	nameInput.Placeholder = "worktree-name"
	nameInput.CharLimit = 64
	nameInput.Width = 30

	branchInput := textinput.New()
	branchInput.Placeholder = "branch-name (optional, defaults to name)"
	branchInput.CharLimit = 64
	branchInput.Width = 40

	baseInput := textinput.New()
	baseInput.Placeholder = "main"
	baseInput.CharLimit = 64
	baseInput.Width = 30
	baseInput.SetValue("main")

	return &worktreeManagerModel{
		list:          l,
		items:         wtItems,
		gitMgr:        gitMgr,
		metaStore:     metaStore,
		nameInput:     nameInput,
		branchInput:   branchInput,
		baseInput:     baseInput,
		pruneSelected: make(map[int]bool),
		staleDays:     30,
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
		m.list.SetWidth(msg.Width)
		m.list.SetHeight(msg.Height - 6)
		return m, nil

	case tea.KeyMsg:
		// Clear message on any key
		m.message = ""

		if m.createMode {
			return m.handleCreateKeys(msg)
		}
		if m.deleteMode {
			return m.handleDeleteKeys(msg)
		}
		if m.pruneMode {
			return m.handlePruneKeys(msg)
		}
		return m.handleMainKeys(msg)
	}

	// Update text inputs if in create mode
	if m.createMode {
		var cmd tea.Cmd
		switch m.createFocus {
		case 0:
			m.nameInput, cmd = m.nameInput.Update(msg)
		case 1:
			m.branchInput, cmd = m.branchInput.Update(msg)
		case 2:
			m.baseInput, cmd = m.baseInput.Update(msg)
		}
		return m, cmd
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m *worktreeManagerModel) handleMainKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c", "q", "esc":
		m.quitting = true
		return m, tea.Quit

	case "c", "n":
		// Create new worktree
		m.createMode = true
		m.createFocus = 0
		m.nameInput.Focus()
		m.nameInput.SetValue("")
		m.branchInput.SetValue("")
		m.baseInput.SetValue("main")
		return m, nil

	case "d", "x":
		// Delete selected worktree
		if i, ok := m.list.SelectedItem().(WorktreeItem); ok {
			if i.IsMain {
				m.setMessage("Cannot delete main worktree", true)
				return m, nil
			}
			m.deleteMode = true
			m.deleteTarget = &i
			m.forceDelete = false
		}
		return m, nil

	case "p":
		// Prune stale worktrees
		m.enterPruneMode()
		return m, nil

	case "r":
		// Refresh list
		return m.refreshList()
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m *worktreeManagerModel) handleCreateKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc":
		m.createMode = false
		m.nameInput.Blur()
		m.branchInput.Blur()
		m.baseInput.Blur()
		return m, nil

	case "tab", "down":
		m.createFocus = (m.createFocus + 1) % 3
		m.updateCreateFocus()
		return m, nil

	case "shift+tab", "up":
		m.createFocus = (m.createFocus + 2) % 3
		m.updateCreateFocus()
		return m, nil

	case "enter":
		if m.createFocus < 2 {
			m.createFocus++
			m.updateCreateFocus()
			return m, nil
		}
		// Submit
		return m.createWorktree()

	case "ctrl+s":
		// Submit from any field
		return m.createWorktree()
	}

	var cmd tea.Cmd
	switch m.createFocus {
	case 0:
		m.nameInput, cmd = m.nameInput.Update(msg)
	case 1:
		m.branchInput, cmd = m.branchInput.Update(msg)
	case 2:
		m.baseInput, cmd = m.baseInput.Update(msg)
	}
	return m, cmd
}

func (m *worktreeManagerModel) updateCreateFocus() {
	m.nameInput.Blur()
	m.branchInput.Blur()
	m.baseInput.Blur()

	switch m.createFocus {
	case 0:
		m.nameInput.Focus()
	case 1:
		m.branchInput.Focus()
	case 2:
		m.baseInput.Focus()
	}
}

func (m *worktreeManagerModel) createWorktree() (tea.Model, tea.Cmd) {
	name := strings.TrimSpace(m.nameInput.Value())
	if name == "" {
		m.setMessage("Name is required", true)
		return m, nil
	}

	branch := strings.TrimSpace(m.branchInput.Value())
	if branch == "" {
		branch = name
	}

	base := strings.TrimSpace(m.baseInput.Value())
	if base == "" {
		base = "main"
	}

	// Create worktree
	path, err := m.gitMgr.Add(name, branch, base)
	if err != nil {
		m.setMessage(fmt.Sprintf("Failed to create: %v", err), true)
		return m, nil
	}

	// Save metadata
	meta := &metadata.WorktreeMetadata{
		Name:       name,
		Path:       path,
		Branch:     branch,
		CreatedAt:  time.Now(),
		LastOpened: time.Now(),
	}
	m.metaStore.Add(meta)
	if err := m.metaStore.Save(); err != nil {
		m.setMessage(fmt.Sprintf("Warning: failed to save metadata: %v", err), true)
	}

	m.createMode = false
	m.nameInput.Blur()
	m.branchInput.Blur()
	m.baseInput.Blur()
	m.setMessage(fmt.Sprintf("✓ Created worktree: %s", name), false)

	// Refresh list
	return m.refreshList()
}

func (m *worktreeManagerModel) handleDeleteKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc", "n":
		m.deleteMode = false
		m.deleteTarget = nil
		return m, nil

	case "y":
		return m.deleteWorktree(false)

	case "f":
		return m.deleteWorktree(true)
	}

	return m, nil
}

func (m *worktreeManagerModel) deleteWorktree(force bool) (tea.Model, tea.Cmd) {
	if m.deleteTarget == nil {
		return m, nil
	}

	name := m.deleteTarget.Name

	// Check if clean
	if !force {
		clean, err := m.gitMgr.IsClean(m.deleteTarget.Path)
		if err != nil {
			m.setMessage(fmt.Sprintf("Failed to check status: %v", err), true)
			m.deleteMode = false
			m.deleteTarget = nil
			return m, nil
		}

		if !clean {
			m.forceDelete = true
			return m, nil
		}
	}

	// Remove worktree
	if err := m.gitMgr.Remove(name, force); err != nil {
		m.setMessage(fmt.Sprintf("Failed to remove: %v", err), true)
		m.deleteMode = false
		m.deleteTarget = nil
		return m, nil
	}

	// Remove from metadata
	m.metaStore.Remove(name)
	if err := m.metaStore.Save(); err != nil {
		m.setMessage(fmt.Sprintf("Warning: failed to update metadata: %v", err), true)
	}

	m.deleteMode = false
	m.deleteTarget = nil
	m.setMessage(fmt.Sprintf("✓ Removed worktree: %s", name), false)

	// Refresh list
	return m.refreshList()
}

func (m *worktreeManagerModel) enterPruneMode() {
	// Get stale worktrees
	staleNames := m.metaStore.GetStale(m.staleDays)

	if len(staleNames) == 0 {
		m.setMessage(fmt.Sprintf("No stale worktrees found (>%d days)", m.staleDays), false)
		return
	}

	// Filter to only clean worktrees
	m.staleItems = []WorktreeItem{}
	for _, name := range staleNames {
		for _, item := range m.items {
			if item.Name == name && !item.IsMain {
				clean, _ := m.gitMgr.IsClean(item.Path)
				if clean {
					m.staleItems = append(m.staleItems, item)
				}
				break
			}
		}
	}

	if len(m.staleItems) == 0 {
		m.setMessage(fmt.Sprintf("No clean stale worktrees found (>%d days)", m.staleDays), false)
		return
	}

	m.pruneMode = true
	m.pruneCursor = 0
	m.pruneSelected = make(map[int]bool)
	// Select all by default
	for i := range m.staleItems {
		m.pruneSelected[i] = true
	}
}

func (m *worktreeManagerModel) handlePruneKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc", "q":
		m.pruneMode = false
		m.staleItems = nil
		return m, nil

	case "up", "k":
		if m.pruneCursor > 0 {
			m.pruneCursor--
		}

	case "down", "j":
		if m.pruneCursor < len(m.staleItems)-1 {
			m.pruneCursor++
		}

	case " ":
		// Toggle selection
		m.pruneSelected[m.pruneCursor] = !m.pruneSelected[m.pruneCursor]

	case "a":
		// Select all
		for i := range m.staleItems {
			m.pruneSelected[i] = true
		}

	case "n":
		// Select none
		for i := range m.staleItems {
			m.pruneSelected[i] = false
		}

	case "enter":
		return m.executePrune()
	}

	return m, nil
}

func (m *worktreeManagerModel) executePrune() (tea.Model, tea.Cmd) {
	removed := 0
	for i, item := range m.staleItems {
		if !m.pruneSelected[i] {
			continue
		}

		if err := m.gitMgr.Remove(item.Name, false); err != nil {
			continue
		}
		m.metaStore.Remove(item.Name)
		removed++
	}

	if err := m.metaStore.Save(); err != nil {
		m.setMessage(fmt.Sprintf("Warning: failed to update metadata: %v", err), true)
	}

	m.pruneMode = false
	m.staleItems = nil
	m.setMessage(fmt.Sprintf("✓ Removed %d worktree(s)", removed), false)

	return m.refreshList()
}

func (m *worktreeManagerModel) refreshList() (tea.Model, tea.Cmd) {
	worktrees, err := m.gitMgr.List()
	if err != nil {
		m.setMessage(fmt.Sprintf("Failed to refresh: %v", err), true)
		return m, nil
	}

	items := make([]list.Item, 0, len(worktrees))
	m.items = make([]WorktreeItem, 0, len(worktrees))

	statuses := m.gitMgr.GetStatuses(worktrees)

	for _, wt := range worktrees {
		status := statuses[wt.Path]

		var meta *metadata.WorktreeMetadata
		if mt, ok := m.metaStore.Get(wt.Name); ok {
			meta = mt
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
		m.items = append(m.items, item)
	}

	m.list.SetItems(items)
	return m, nil
}

func (m *worktreeManagerModel) setMessage(msg string, isError bool) {
	m.message = msg
	if isError {
		m.messageStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#FF6B6B"))
	} else {
		m.messageStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#04B575"))
	}
}

func (m *worktreeManagerModel) View() string {
	if m.quitting {
		return ""
	}

	var b strings.Builder

	if m.createMode {
		return m.viewCreateForm()
	}

	if m.deleteMode {
		return m.viewDeleteConfirm()
	}

	if m.pruneMode {
		return m.viewPruneMode()
	}

	// Main list view
	b.WriteString(m.list.View())

	// Message
	if m.message != "" {
		b.WriteString("\n")
		b.WriteString(m.messageStyle.Render(m.message))
	}

	// Help
	helpText := helpStyle.Render("\nc create • d delete • p prune • r refresh • q quit")
	b.WriteString(helpText)

	return b.String()
}

func (m *worktreeManagerModel) viewCreateForm() string {
	var b strings.Builder

	title := titleStyle.Render("📁 Create New Worktree")
	b.WriteString(title)
	b.WriteString("\n\n")

	// Name field
	nameLabel := "Name:"
	if m.createFocus == 0 {
		nameLabel = lipgloss.NewStyle().Foreground(lipgloss.Color("#7D56F4")).Bold(true).Render("▸ Name:")
	}
	b.WriteString(fmt.Sprintf("%s\n  %s\n\n", nameLabel, m.nameInput.View()))

	// Branch field
	branchLabel := "Branch (optional):"
	if m.createFocus == 1 {
		branchLabel = lipgloss.NewStyle().Foreground(lipgloss.Color("#7D56F4")).Bold(true).Render("▸ Branch (optional):")
	}
	b.WriteString(fmt.Sprintf("%s\n  %s\n\n", branchLabel, m.branchInput.View()))

	// Base branch field
	baseLabel := "Base branch:"
	if m.createFocus == 2 {
		baseLabel = lipgloss.NewStyle().Foreground(lipgloss.Color("#7D56F4")).Bold(true).Render("▸ Base branch:")
	}
	b.WriteString(fmt.Sprintf("%s\n  %s\n\n", baseLabel, m.baseInput.View()))

	// Message
	if m.message != "" {
		b.WriteString(m.messageStyle.Render(m.message))
		b.WriteString("\n")
	}

	// Help
	helpText := helpStyle.Render("\ntab next field • ctrl+s create • esc cancel")
	b.WriteString(helpText)

	return b.String()
}

func (m *worktreeManagerModel) viewDeleteConfirm() string {
	var b strings.Builder

	title := titleStyle.Render("🗑  Delete Worktree")
	b.WriteString(title)
	b.WriteString("\n\n")

	if m.deleteTarget != nil {
		b.WriteString(fmt.Sprintf("Worktree: %s\n", lipgloss.NewStyle().Bold(true).Render(m.deleteTarget.Name)))
		b.WriteString(fmt.Sprintf("Path:     %s\n", m.deleteTarget.Path))
		b.WriteString(fmt.Sprintf("Branch:   %s\n\n", m.deleteTarget.Branch))

		if m.forceDelete {
			b.WriteString(lipgloss.NewStyle().
				Foreground(lipgloss.Color("#FF6B6B")).
				Bold(true).
				Render("⚠  This worktree has uncommitted changes!"))
			b.WriteString("\n\n")
			b.WriteString("Press 'f' to force delete (lose changes)\n")
			b.WriteString("Press 'n' or 'esc' to cancel\n")
		} else {
			b.WriteString("Are you sure you want to delete this worktree?\n\n")
			b.WriteString("Press 'y' to confirm\n")
			b.WriteString("Press 'n' or 'esc' to cancel\n")
		}
	}

	return b.String()
}

func (m *worktreeManagerModel) viewPruneMode() string {
	var b strings.Builder

	title := titleStyle.Render(fmt.Sprintf("🧹 Prune Stale Worktrees (>%d days)", m.staleDays))
	b.WriteString(title)
	b.WriteString("\n\n")

	for i, item := range m.staleItems {
		cursor := "  "
		if i == m.pruneCursor {
			cursor = "▸ "
		}

		checkbox := "[ ]"
		if m.pruneSelected[i] {
			checkbox = "[✓]"
		}

		name := item.Name
		if i == m.pruneCursor {
			name = lipgloss.NewStyle().Foreground(lipgloss.Color("#7D56F4")).Render(name)
		}

		lastOpened := ""
		if item.Metadata != nil {
			lastOpened = fmt.Sprintf(" (last: %s)", item.Metadata.LastOpened.Format("2006-01-02"))
		}

		b.WriteString(fmt.Sprintf("%s%s %s%s\n", cursor, checkbox, name,
			lipgloss.NewStyle().Foreground(lipgloss.Color("#626262")).Render(lastOpened)))
	}

	// Count selected
	selected := 0
	for _, v := range m.pruneSelected {
		if v {
			selected++
		}
	}

	b.WriteString(fmt.Sprintf("\n%d of %d selected\n", selected, len(m.staleItems)))

	// Help
	helpText := helpStyle.Render("\n↑/↓ navigate • space toggle • a all • n none • enter delete • esc cancel")
	b.WriteString(helpText)

	return b.String()
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
