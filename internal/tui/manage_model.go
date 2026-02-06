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
	"github.com/darkLord19/wtx/internal/validation"
)

// ManageModel handles the logic for listing, creating, deleting, and pruning worktrees.
// It is designed to be embedded in other models (like the main manager) or run standalone.
type ManageModel struct {
	// Dependencies
	gitMgr    *git.Manager
	metaStore *metadata.Store

	// State
	List      list.Model
	Items     []WorktreeItem
	Mode      ManageMode
	Width     int
	Height    int

	// Create Form
	Inputs    [3]textinput.Model // 0: Name, 1: Branch, 2: Base
	Focus     int

	// Delete Confirmation
	DeleteTarget *WorktreeItem
	ForceDelete  bool

	// Prune Mode
	StaleItems    []WorktreeItem
	PruneCursor   int
	PruneSelected map[int]bool
	StaleDays     int

	// UI
	Message Message
	Help    *HelpPanel
}

// NewManageModel creates a new ManageModel
func NewManageModel(gitMgr *git.Manager, metaStore *metadata.Store) (*ManageModel, error) {
	items, wtItems, err := LoadWorktreeItems(gitMgr, metaStore)
	if err != nil {
		return nil, err
	}

	l := CreateListModel(items, "Worktree Manager")

	// Initialize inputs
	var inputs [3]textinput.Model
	inputs[0] = textinput.New()
	inputs[0].Placeholder = "worktree-name"
	inputs[0].CharLimit = 64
	inputs[0].Width = 30

	inputs[1] = textinput.New()
	inputs[1].Placeholder = "branch-name (optional)"
	inputs[1].CharLimit = 64
	inputs[1].Width = 40

	inputs[2] = textinput.New()
	inputs[2].Placeholder = "main"
	inputs[2].CharLimit = 64
	inputs[2].Width = 30
	inputs[2].SetValue("main")

	// Initialize Help
	help := NewHelpPanel()
	help.AddSection("Global", GetGlobalHelp())
	help.AddSection("Manage", GetManageHelp())

	return &ManageModel{
		gitMgr:        gitMgr,
		metaStore:     metaStore,
		List:          l,
		Items:         wtItems,
		Mode:          ManageModeList,
		Inputs:        inputs,
		PruneSelected: make(map[int]bool),
		StaleDays:     30,
		Help:          help,
	}, nil
}

func (m *ManageModel) Init() tea.Cmd {
	return nil
}

func (m *ManageModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.Width = msg.Width
		m.Height = msg.Height
		m.List.SetWidth(msg.Width)
		m.List.SetHeight(msg.Height - 6) // Adjust for help/message
		return m, nil

	case WorktreeListMsg:
		if msg.Err != nil {
			m.SetMessage(fmt.Sprintf("Failed to refresh: %v", msg.Err), true)
			return m, nil
		}

		items := make([]list.Item, 0, len(msg.Items))
		m.Items = make([]WorktreeItem, 0, len(msg.Items))
		for _, item := range msg.Items {
			items = append(items, item)
			m.Items = append(m.Items, item)
		}
		m.List.SetItems(items)

		if m.Message.Text() == "Refreshing..." {
			m.Message.Clear()
		}
		return m, nil
	}

	switch m.Mode {
	case ManageModeList:
		return m.updateList(msg)
	case ManageModeCreate:
		return m.updateCreate(msg)
	case ManageModeDelete:
		return m.updateDelete(msg)
	case ManageModePrune:
		return m.updatePrune(msg)
	}

	return m, nil
}

func (m *ManageModel) updateList(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		// Clear message on keypress
		if !m.Message.IsEmpty() {
			m.Message.Clear()
		}

		switch msg.String() {
		case "c", "n":
			m.Mode = ManageModeCreate
			m.Focus = 0
			m.Inputs[0].SetValue("")
			m.Inputs[1].SetValue("")
			m.Inputs[2].SetValue("main")
			m.Inputs[0].Focus()
			return m, nil

		case "d", "x":
			if i, ok := m.List.SelectedItem().(WorktreeItem); ok {
				if i.IsMain {
					m.SetMessage("Cannot delete main worktree", true)
					return m, nil
				}
				m.Mode = ManageModeDelete
				m.DeleteTarget = &i
				m.ForceDelete = false
			}
			return m, nil

		case "p":
			m.enterPruneMode()
			return m, nil

		case "r":
			return m.RefreshList()

		case "?":
			m.Help.Toggle()
			return m, nil
		}
	}

	var cmd tea.Cmd
	m.List, cmd = m.List.Update(msg)
	return m, cmd
}

func (m *ManageModel) updateCreate(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			m.Mode = ManageModeList
			m.blurInputs()
			return m, nil

		case "tab", "down":
			m.Focus = (m.Focus + 1) % 3
			m.updateFocus()
			return m, nil

		case "shift+tab", "up":
			m.Focus = (m.Focus + 2) % 3
			m.updateFocus()
			return m, nil

		case "enter":
			if m.Focus < 2 {
				m.Focus++
				m.updateFocus()
				return m, nil
			}
			return m.createWorktree()

		case "ctrl+s":
			return m.createWorktree()
		}
	}

	var cmd tea.Cmd
	m.Inputs[m.Focus], cmd = m.Inputs[m.Focus].Update(msg)
	return m, cmd
}

func (m *ManageModel) updateDelete(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc", "n":
			m.Mode = ManageModeList
			m.DeleteTarget = nil
			return m, nil

		case "y":
			return m.deleteWorktree(false)

		case "f":
			return m.deleteWorktree(true)
		}
	}
	return m, nil
}

func (m *ManageModel) updatePrune(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc", "q":
			m.Mode = ManageModeList
			m.StaleItems = nil
			return m, nil

		case "up", "k":
			if m.PruneCursor > 0 {
				m.PruneCursor--
			}

		case "down", "j":
			if m.PruneCursor < len(m.StaleItems)-1 {
				m.PruneCursor++
			}

		case " ":
			m.PruneSelected[m.PruneCursor] = !m.PruneSelected[m.PruneCursor]

		case "a":
			for i := range m.StaleItems {
				m.PruneSelected[i] = true
			}

		case "n":
			for i := range m.StaleItems {
				m.PruneSelected[i] = false
			}

		case "enter":
			return m.executePrune()
		}
	}
	return m, nil
}

// Actions

func (m *ManageModel) RefreshList() (tea.Model, tea.Cmd) {
	m.Message = NewInfoMessage("Refreshing...")
	return m, fetchWorktreesCmd(m.gitMgr, m.metaStore)
}

func (m *ManageModel) createWorktree() (tea.Model, tea.Cmd) {
	name := strings.TrimSpace(m.Inputs[0].Value())
	branch := strings.TrimSpace(m.Inputs[1].Value())
	base := strings.TrimSpace(m.Inputs[2].Value())

	if branch == "" {
		branch = name
	}
	if base == "" {
		base = "main"
	}

	// Validation
	validator := validation.NewWorktreeValidator()
	if err := validator.ValidateName(name); err != nil {
		m.SetMessage(fmt.Sprintf("Invalid name: %v", err), true)
		return m, nil
	}
	if err := validator.ValidateBranchName(branch); err != nil {
		m.SetMessage(fmt.Sprintf("Invalid branch: %v", err), true)
		return m, nil
	}

	path, err := m.gitMgr.Add(name, branch, base)
	if err != nil {
		m.SetMessage(fmt.Sprintf("Failed to create: %v", err), true)
		return m, nil
	}

	meta := &metadata.WorktreeMetadata{
		Name:       name,
		Path:       path,
		Branch:     branch,
		CreatedAt:  time.Now(),
		LastOpened: time.Now(),
	}
	m.metaStore.Add(meta)
	if err := m.metaStore.Save(); err != nil {
		m.SetMessage(fmt.Sprintf("Warning: metadata save failed: %v", err), true)
	}

	m.Mode = ManageModeList
	m.blurInputs()
	m.SetMessage(fmt.Sprintf("✓ Created worktree: %s", name), false)

	return m.RefreshList()
}

func (m *ManageModel) deleteWorktree(force bool) (tea.Model, tea.Cmd) {
	if m.DeleteTarget == nil {
		return m, nil
	}

	name := m.DeleteTarget.Name

	if !force {
		clean, err := m.gitMgr.IsClean(m.DeleteTarget.Path)
		if err != nil {
			m.SetMessage(fmt.Sprintf("Failed to check status: %v", err), true)
			m.Mode = ManageModeList
			m.DeleteTarget = nil
			return m, nil
		}

		if !clean {
			m.ForceDelete = true
			return m, nil
		}
	}

	if err := m.gitMgr.Remove(name, force); err != nil {
		m.SetMessage(fmt.Sprintf("Failed to remove: %v", err), true)
		m.Mode = ManageModeList
		m.DeleteTarget = nil
		return m, nil
	}

	m.metaStore.Remove(name)
	if err := m.metaStore.Save(); err != nil {
		m.SetMessage(fmt.Sprintf("Warning: metadata save failed: %v", err), true)
	}

	m.Mode = ManageModeList
	m.DeleteTarget = nil
	m.SetMessage(fmt.Sprintf("✓ Removed worktree: %s", name), false)

	return m.RefreshList()
}

func (m *ManageModel) enterPruneMode() {
	staleNames := m.metaStore.GetStale(m.StaleDays)

	if len(staleNames) == 0 {
		m.SetMessage(fmt.Sprintf("No stale worktrees found (>%d days)", m.StaleDays), false)
		return
	}

	m.StaleItems = []WorktreeItem{}
	itemMap := make(map[string]WorktreeItem)
	for _, item := range m.Items {
		itemMap[item.Name] = item
	}

	for _, name := range staleNames {
		if item, ok := itemMap[name]; ok {
			if !item.IsMain {
				clean, _ := m.gitMgr.IsClean(item.Path)
				if clean {
					m.StaleItems = append(m.StaleItems, item)
				}
			}
		}
	}

	if len(m.StaleItems) == 0 {
		m.SetMessage(fmt.Sprintf("No clean stale worktrees found (>%d days)", m.StaleDays), false)
		return
	}

	m.Mode = ManageModePrune
	m.PruneCursor = 0
	m.PruneSelected = make(map[int]bool)
	for i := range m.StaleItems {
		m.PruneSelected[i] = true
	}
}

func (m *ManageModel) executePrune() (tea.Model, tea.Cmd) {
	removed := 0
	for i, item := range m.StaleItems {
		if !m.PruneSelected[i] {
			continue
		}

		if err := m.gitMgr.Remove(item.Name, false); err != nil {
			continue
		}
		m.metaStore.Remove(item.Name)
		removed++
	}

	if err := m.metaStore.Save(); err != nil {
		m.SetMessage(fmt.Sprintf("Warning: metadata save failed: %v", err), true)
	}

	m.Mode = ManageModeList
	m.StaleItems = nil
	m.SetMessage(fmt.Sprintf("✓ Removed %d worktree(s)", removed), false)

	return m.RefreshList()
}

// Helpers

func (m *ManageModel) updateFocus() {
	m.blurInputs()
	m.Inputs[m.Focus].Focus()
}

func (m *ManageModel) blurInputs() {
	for i := range m.Inputs {
		m.Inputs[i].Blur()
	}
}

func (m *ManageModel) SetMessage(msg string, isError bool) {
	if isError {
		m.Message = NewErrorMessage(msg)
	} else {
		m.Message = NewSuccessMessage(msg)
	}
}

func (m *ManageModel) IsInSubMode() bool {
	return m.Mode != ManageModeList
}

// View

func (m *ManageModel) View() string {
	if m.Help.IsVisible() {
		return m.Help.Render(m.Width)
	}

	var b strings.Builder

	// Breadcrumb
	breadcrumb := NewBreadcrumb("wtx", "Manage")
	switch m.Mode {
	case ManageModeCreate:
		breadcrumb.Add("Create")
	case ManageModeDelete:
		breadcrumb.Add("Delete")
	case ManageModePrune:
		breadcrumb.Add("Prune")
	}

	// Only show breadcrumb if we are in submode (optional, but consistent)
	// Or we could let the parent handle breadcrumbs.
	// For standalone reuse, we might want to render it here, or make it optional.
	// For now, let's keep it here.
	b.WriteString(breadcrumb.Render())
	b.WriteString("\n\n")

	switch m.Mode {
	case ManageModeCreate:
		b.WriteString(m.viewCreateForm())
	case ManageModeDelete:
		b.WriteString(m.viewDeleteConfirm())
	case ManageModePrune:
		b.WriteString(m.viewPruneMode())
	default:
		b.WriteString(m.List.View())
		b.WriteString("\n")
		b.WriteString(helpStyle.Render("c create • d delete • p prune • r refresh • q quit"))
	}

	if !m.Message.IsEmpty() {
		b.WriteString("\n")
		b.WriteString(m.Message.Render())
	}

	return b.String()
}

func (m *ManageModel) viewCreateForm() string {
	var b strings.Builder

	title := titleStyle.Render("📁 Create New Worktree")
	b.WriteString(title)
	b.WriteString("\n\n")

	labels := []string{"Name:", "Branch (optional):", "Base branch:"}
	for i, label := range labels {
		labelStyle := lipgloss.NewStyle()
		if i == m.Focus {
			labelStyle = labelStyle.Foreground(lipgloss.Color("#7D56F4")).Bold(true)
			label = "▸ " + label
		} else {
			label = "  " + label
		}
		b.WriteString(fmt.Sprintf("%s\n  %s\n\n", labelStyle.Render(label), m.Inputs[i].View()))
	}

	b.WriteString(helpStyle.Render("tab next • ctrl+s create • esc cancel"))
	return b.String()
}

func (m *ManageModel) viewDeleteConfirm() string {
	var b strings.Builder

	title := titleStyle.Render("🗑  Delete Worktree")
	b.WriteString(title)
	b.WriteString("\n\n")

	if m.DeleteTarget != nil {
		b.WriteString(fmt.Sprintf("Worktree: %s\n", lipgloss.NewStyle().Bold(true).Render(m.DeleteTarget.Name)))
		b.WriteString(fmt.Sprintf("Path:     %s\n", m.DeleteTarget.Path))
		b.WriteString(fmt.Sprintf("Branch:   %s\n\n", m.DeleteTarget.Branch))

		if m.ForceDelete {
			b.WriteString(lipgloss.NewStyle().
				Foreground(lipgloss.Color("#FF6B6B")).
				Bold(true).
				Render("⚠  This worktree has uncommitted changes!"))
			b.WriteString("\n\n")
			b.WriteString("f force delete • n/esc cancel\n")
		} else {
			b.WriteString("Delete this worktree?\n\n")
			b.WriteString("y confirm • n/esc cancel\n")
		}
	}

	return b.String()
}

func (m *ManageModel) viewPruneMode() string {
	var b strings.Builder

	title := titleStyle.Render(fmt.Sprintf("🧹 Prune Stale Worktrees (>%d days)", m.StaleDays))
	b.WriteString(title)
	b.WriteString("\n\n")

	for i, item := range m.StaleItems {
		cursor := "  "
		if i == m.PruneCursor {
			cursor = "▸ "
		}

		checkbox := "[ ]"
		if m.PruneSelected[i] {
			checkbox = "[✓]"
		}

		name := item.Name
		if i == m.PruneCursor {
			name = lipgloss.NewStyle().Foreground(lipgloss.Color("#7D56F4")).Render(name)
		}

		lastOpened := ""
		if item.Metadata != nil {
			lastOpened = fmt.Sprintf(" (last: %s)", item.Metadata.LastOpened.Format("2006-01-02"))
		}

		b.WriteString(fmt.Sprintf("%s%s %s%s\n", cursor, checkbox, name,
			lipgloss.NewStyle().Foreground(lipgloss.Color("#626262")).Render(lastOpened)))
	}

	selected := 0
	for _, v := range m.PruneSelected {
		if v {
			selected++
		}
	}

	b.WriteString(fmt.Sprintf("\n%d of %d selected\n", selected, len(m.StaleItems)))
	b.WriteString(helpStyle.Render("\n↑/↓ navigate • space toggle • a all • n none • enter delete • esc cancel"))

	return b.String()
}
