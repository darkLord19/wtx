package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/darkLord19/wtx/internal/config"
	"github.com/darkLord19/wtx/internal/editor"
	"github.com/darkLord19/wtx/internal/git"
	"github.com/darkLord19/wtx/internal/metadata"
)

// Tab represents a tab in the TUI
type Tab int

const (
	TabWorktrees Tab = iota
	TabManage
	TabSettings
)

// managerModel is the main TUI model with tabs
type managerModel struct {
	// Core dependencies
	gitMgr     *git.Manager
	metaStore  *metadata.Store
	config     *config.Config
	edDetector *editor.Detector

	// UI state
	activeTab Tab
	width     int
	height    int
	quitting  bool

	// Worktrees tab (selector)
	worktreeList list.Model
	items        []WorktreeItem
	choice       *WorktreeItem

	// Manage tab
	manageModel *ManageModel

	// Settings tab
	settings      []SettingItem
	settingCursor int
	settingEdit   bool
	settingInput  textinput.Model
	optionCursor  int

	// Messages
	message Message

	// Help
	help *HelpPanel
}

// ManageMode represents the current mode in the manage tab
type ManageMode int

const (
	ManageModeList ManageMode = iota
	ManageModeCreate
	ManageModeDelete
	ManageModePrune
)

// WorktreeListMsg contains the list of worktrees fetched asynchronously
type WorktreeListMsg struct {
	Items []WorktreeItem
	Err   error
}

// NewManagerModel creates a new manager TUI model
func NewManagerModel(gitMgr *git.Manager, metaStore *metadata.Store, cfg *config.Config, edDetector *editor.Detector) (*managerModel, error) {
	// Load worktrees
	items, wtItems, err := LoadWorktreeItems(gitMgr, metaStore)
	if err != nil {
		return nil, fmt.Errorf("failed to list worktrees: %w", err)
	}

	// Create worktree list
	l := CreateListModel(items, "Select Worktree")

	// Initialize ManageModel
	manageModel, err := NewManageModel(gitMgr, metaStore)
	if err != nil {
		return nil, err
	}

	// Initialize help panel
	help := NewHelpPanel()
	help.AddSection("Global", GetGlobalHelp())
	help.AddSection("Navigation", GetNavigationHelp())
	help.AddSection("Worktrees", GetWorktreesHelp())
	help.AddSection("Manage", GetManageHelp())
	help.AddSection("Settings", GetSettingsHelp())

	// Settings input
	settingInput := textinput.New()
	settingInput.Placeholder = "Enter value..."
	settingInput.CharLimit = 256
	settingInput.Width = 40

	// Get available editors
	editors := edDetector.DetectAll()
	editorOptions := []string{"(auto-detect)"}
	for _, ed := range editors {
		editorOptions = append(editorOptions, strings.ToLower(ed.Name()))
	}
	editorOptions = append(editorOptions, "(custom)")

	settings := []SettingItem{
		{
			Name:        "Editor",
			Description: "Preferred code editor (select or enter custom)",
			Type:        SettingEditor,
			Value:       cfg.Editor,
			Options:     editorOptions,
		},
		{
			Name:        "Custom Editor Command",
			Description: "Custom editor command (e.g., 'subl', 'atom')",
			Type:        SettingCustomEditor,
			Value:       cfg.Editor,
			Options:     nil, // Text input
		},
		{
			Name:        "Reuse Window",
			Description: "Reuse existing editor window",
			Type:        SettingReuseWindow,
			Value:       boolToString(cfg.ReuseWindow),
			Options:     []string{"true", "false"},
		},
		{
			Name:        "Worktree Directory",
			Description: "Where to create worktrees",
			Type:        SettingWorktreeDir,
			Value:       cfg.WorktreeDir,
			Options:     nil,
		},
		{
			Name:        "Auto Start Dev",
			Description: "Auto-start dev server",
			Type:        SettingAutoStartDev,
			Value:       boolToString(cfg.AutoStartDev),
			Options:     []string{"true", "false"},
		},
	}

	return &managerModel{
		gitMgr:        gitMgr,
		metaStore:     metaStore,
		config:        cfg,
		edDetector:    edDetector,
		activeTab:     TabWorktrees,
		worktreeList:  l,
		items:         wtItems,
		manageModel:   manageModel,
		settingInput:  settingInput,
		settings:      settings,
		help:          help,
	}, nil
}

func (m *managerModel) Init() tea.Cmd {
	return nil
}

func (m *managerModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.worktreeList.SetWidth(msg.Width)
		m.worktreeList.SetHeight(msg.Height - 8)

		// Propagate resize to manage model
		m.manageModel.Update(msg)
		return m, nil

	case WorktreeListMsg:
		if msg.Err != nil {
			m.setMessage(fmt.Sprintf("Failed to refresh: %v", msg.Err), true)
			return m, nil
		}

		items := make([]list.Item, 0, len(msg.Items))
		m.items = make([]WorktreeItem, 0, len(msg.Items))

		for _, item := range msg.Items {
			items = append(items, item)
			m.items = append(m.items, item)
		}

		m.worktreeList.SetItems(items)

		// Update manage model list as well (sync)
		m.manageModel.Update(msg)

		// Only clear message if it was "Refreshing..."
		if m.message.Text() == "Refreshing..." {
			m.message.Clear()
		}
		return m, nil

	case tea.KeyMsg:
		// Clear message on any key
		if !m.message.IsEmpty() {
			m.message.Clear()
		}

		// Global keys
		switch msg.String() {
		case "?":
			m.help.Toggle()
			return m, nil

		case "ctrl+c":
			m.quitting = true
			return m, tea.Quit

		case "1":
			if !m.isInSubMode() {
				m.activeTab = TabWorktrees
				return m, nil
			}

		case "2":
			if !m.isInSubMode() {
				m.activeTab = TabManage
				return m, nil
			}

		case "3":
			if !m.isInSubMode() {
				m.activeTab = TabSettings
				return m, nil
			}
		}

		// Tab-specific handling
		switch m.activeTab {
		case TabWorktrees:
			return m.updateWorktreesTab(msg)
		case TabManage:
			return m.updateManageTab(msg)
		case TabSettings:
			return m.updateSettingsTab(msg)
		}
	}

	return m, nil
}

func (m *managerModel) isInSubMode() bool {
	if m.activeTab == TabManage && m.manageModel.IsInSubMode() {
		return true
	}
	if m.activeTab == TabSettings && m.settingEdit {
		return true
	}
	return false
}

// Worktrees tab handlers
func (m *managerModel) updateWorktreesTab(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "q", "esc":
		m.quitting = true
		return m, tea.Quit

	case "enter":
		if i, ok := m.worktreeList.SelectedItem().(WorktreeItem); ok {
			m.choice = &i
			return m, tea.Quit
		}
	}

	var cmd tea.Cmd
	m.worktreeList, cmd = m.worktreeList.Update(msg)
	return m, cmd
}

// Manage tab handlers
func (m *managerModel) updateManageTab(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	// Let the manage model handle it
	newModel, cmd := m.manageModel.Update(msg)
	m.manageModel = newModel.(*ManageModel)

	// Check for quit signal from submodel if it were to emit one (not currently implemented, but safe)
	return m, cmd
}

// Settings tab handlers
func (m *managerModel) updateSettingsTab(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	if m.settingEdit {
		return m.updateSettingEdit(msg)
	}

	switch msg.String() {
	case "q", "esc":
		m.quitting = true
		return m, tea.Quit

	case "up", "k":
		if m.settingCursor > 0 {
			m.settingCursor--
		}

	case "down", "j":
		if m.settingCursor < len(m.settings)-1 {
			m.settingCursor++
		}

	case "enter", " ":
		m.startSettingEdit()

	case "s":
		m.saveSettings()
	}

	return m, nil
}

func (m *managerModel) updateSettingEdit(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	setting := &m.settings[m.settingCursor]

	switch msg.String() {
	case "esc":
		m.settingEdit = false
		m.settingInput.Blur()
		return m, nil

	case "enter":
		if setting.Options == nil {
			setting.Value = m.settingInput.Value()
		}
		m.settingEdit = false
		m.settingInput.Blur()
		return m, nil

	case "left", "h":
		if setting.Options != nil {
			m.optionCursor--
			if m.optionCursor < 0 {
				m.optionCursor = len(setting.Options) - 1
			}
			setting.Value = setting.Options[m.optionCursor]
			if setting.Value == "(auto-detect)" {
				setting.Value = ""
			}
		}

	case "right", "l":
		if setting.Options != nil {
			m.optionCursor++
			if m.optionCursor >= len(setting.Options) {
				m.optionCursor = 0
			}
			setting.Value = setting.Options[m.optionCursor]
			if setting.Value == "(auto-detect)" {
				setting.Value = ""
			}
		}

	default:
		if setting.Options == nil {
			var cmd tea.Cmd
			m.settingInput, cmd = m.settingInput.Update(msg)
			return m, cmd
		}
	}

	return m, nil
}

// Helper methods

func fetchWorktreesCmd(gitMgr *git.Manager, metaStore *metadata.Store) tea.Cmd {
	return func() tea.Msg {
		_, wtItems, err := LoadWorktreeItems(gitMgr, metaStore)
		if err != nil {
			return WorktreeListMsg{Err: err}
		}
		return WorktreeListMsg{Items: wtItems}
	}
}

func (m *managerModel) startSettingEdit() {
	m.settingEdit = true
	setting := m.settings[m.settingCursor]

	if setting.Options != nil {
		m.optionCursor = 0
		currentVal := setting.Value
		if currentVal == "" && setting.Type == SettingEditor {
			currentVal = "(auto-detect)"
		}
		for i, opt := range setting.Options {
			if opt == currentVal {
				m.optionCursor = i
				break
			}
		}
	} else {
		m.settingInput.SetValue(setting.Value)
		m.settingInput.Focus()
	}
}

func (m *managerModel) saveSettings() {
	var editorValue string
	var customEditorValue string

	for _, setting := range m.settings {
		switch setting.Type {
		case SettingEditor:
			editorValue = setting.Value
		case SettingCustomEditor:
			customEditorValue = setting.Value
		case SettingReuseWindow:
			m.config.ReuseWindow = setting.Value == "true"
		case SettingWorktreeDir:
			m.config.WorktreeDir = setting.Value
		case SettingAutoStartDev:
			m.config.AutoStartDev = setting.Value == "true"
		}
	}

	// Handle editor selection - use custom if "(custom)" is selected
	if editorValue == "(custom)" && customEditorValue != "" {
		m.config.Editor = customEditorValue
	} else if editorValue != "(custom)" {
		m.config.Editor = editorValue
	}

	if err := m.config.Save(); err != nil {
		m.setMessage(fmt.Sprintf("Failed to save: %v", err), true)
	} else {
		m.setMessage("✓ Settings saved!", false)
	}
}

func (m *managerModel) setMessage(msg string, isError bool) {
	if isError {
		m.message = NewErrorMessage(msg)
	} else {
		m.message = NewSuccessMessage(msg)
	}
}

func (m *managerModel) View() string {
	if m.quitting {
		return ""
	}

	// If help is visible, show it overlaying everything
	if m.help.IsVisible() {
		return m.help.Render(m.width)
	}

	var b strings.Builder

	// Breadcrumb
	breadcrumb := NewBreadcrumb("wtx")
	switch m.activeTab {
	case TabWorktrees:
		breadcrumb.Add("Worktrees")
	case TabManage:
		breadcrumb.Add("Manage")
		if m.manageModel.Mode == ManageModeCreate {
			breadcrumb.Add("Create")
		} else if m.manageModel.Mode == ManageModeDelete {
			breadcrumb.Add("Delete")
		} else if m.manageModel.Mode == ManageModePrune {
			breadcrumb.Add("Prune")
		}
	case TabSettings:
		breadcrumb.Add("Settings")
		if m.settingEdit {
			breadcrumb.Add("Edit")
		}
	}
	b.WriteString(breadcrumb.Render())
	b.WriteString("\n\n")

	// Tab bar
	b.WriteString(m.renderTabBar())
	b.WriteString("\n\n")

	// Tab content
	switch m.activeTab {
	case TabWorktrees:
		b.WriteString(m.viewWorktreesTab())
	case TabManage:
		b.WriteString(m.viewManageTab())
	case TabSettings:
		b.WriteString(m.viewSettingsTab())
	}

	// Message
	if !m.message.IsEmpty() {
		b.WriteString("\n")
		b.WriteString(m.message.Render())
	}

	return b.String()
}

func (m *managerModel) renderTabBar() string {
	tabs := []string{"[1] Worktrees", "[2] Manage", "[3] Settings"}

	var rendered []string
	for i, tab := range tabs {
		style := lipgloss.NewStyle().Padding(0, 2)
		if Tab(i) == m.activeTab {
			style = style.
				Background(lipgloss.Color("#7D56F4")).
				Foreground(lipgloss.Color("#FFFFFF")).
				Bold(true)
		} else {
			style = style.
				Foreground(lipgloss.Color("#626262"))
		}
		rendered = append(rendered, style.Render(tab))
	}

	return lipgloss.JoinHorizontal(lipgloss.Center, rendered...)
}

func (m *managerModel) viewWorktreesTab() string {
	var b strings.Builder
	b.WriteString(m.worktreeList.View())
	b.WriteString("\n")
	b.WriteString(helpStyle.Render("enter open • q quit"))
	return b.String()
}

func (m *managerModel) viewManageTab() string {
	// Delegate to the shared component
	return m.manageModel.View()
}

func (m *managerModel) viewSettingsTab() string {
	var b strings.Builder

	title := titleStyle.Render("⚙  Settings")
	b.WriteString(title)
	b.WriteString("\n\n")

	for i, setting := range m.settings {
		cursor := "  "
		if i == m.settingCursor {
			cursor = "▸ "
		}

		nameStyle := lipgloss.NewStyle().Bold(true)
		if i == m.settingCursor {
			nameStyle = nameStyle.Foreground(lipgloss.Color("#7D56F4"))
		}

		name := nameStyle.Render(setting.Name)

		var valueDisplay string
		if m.settingEdit && i == m.settingCursor {
			if setting.Options != nil {
				valueDisplay = m.renderOptionSelector(setting)
			} else {
				valueDisplay = m.settingInput.View()
			}
		} else {
			displayVal := setting.Value
			if displayVal == "" && setting.Type == SettingEditor {
				displayVal = "(auto-detect)"
			}
			valueDisplay = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#04B575")).
				Render(displayVal)
		}

		desc := lipgloss.NewStyle().
			Foreground(lipgloss.Color("#626262")).
			Render(setting.Description)

		b.WriteString(fmt.Sprintf("%s%s: %s\n", cursor, name, valueDisplay))
		b.WriteString(fmt.Sprintf("   %s\n\n", desc))
	}

	var helpText string
	if m.settingEdit {
		if m.settings[m.settingCursor].Options != nil {
			helpText = "←/→ change • enter confirm • esc cancel"
		} else {
			helpText = "type to edit • enter confirm • esc cancel"
		}
	} else {
		helpText = "↑/↓ navigate • enter edit • s save • q quit"
	}
	b.WriteString(helpStyle.Render(helpText))

	return b.String()
}

func (m *managerModel) renderOptionSelector(setting SettingItem) string {
	var parts []string
	for i, opt := range setting.Options {
		style := lipgloss.NewStyle()
		if i == m.optionCursor {
			style = style.
				Background(lipgloss.Color("#7D56F4")).
				Foreground(lipgloss.Color("#FFFFFF")).
				Padding(0, 1)
		} else {
			style = style.
				Foreground(lipgloss.Color("#626262")).
				Padding(0, 1)
		}
		parts = append(parts, style.Render(opt))
	}
	return lipgloss.JoinHorizontal(lipgloss.Center, parts...)
}

// GetChoice returns the selected worktree (if any)
func (m *managerModel) GetChoice() *WorktreeItem {
	return m.choice
}

// RunManager starts the main manager TUI and returns the selected worktree
func RunManager(gitMgr *git.Manager, metaStore *metadata.Store, cfg *config.Config, edDetector *editor.Detector) (*WorktreeItem, error) {
	m, err := NewManagerModel(gitMgr, metaStore, cfg, edDetector)
	if err != nil {
		return nil, err
	}

	p := tea.NewProgram(m, tea.WithAltScreen())

	finalModel, err := p.Run()
	if err != nil {
		return nil, err
	}

	if fm, ok := finalModel.(*managerModel); ok {
		return fm.GetChoice(), nil
	}

	return nil, nil
}
