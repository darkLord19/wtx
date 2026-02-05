package tui

import (
	"fmt"
	"strings"
	"time"

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
	manageMode    ManageMode
	createInputs  [3]textinput.Model // name, branch, base
	createFocus   int
	deleteTarget  *WorktreeItem
	forceDelete   bool
	staleItems    []WorktreeItem
	pruneCursor   int
	pruneSelected map[int]bool
	staleDays     int

	// Settings tab
	settings      []SettingItem
	settingCursor int
	settingEdit   bool
	settingInput  textinput.Model
	optionCursor  int

	// Messages
	message      string
	messageStyle lipgloss.Style
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

	// Create worktree list
	delegate := list.NewDefaultDelegate()
	l := list.New(items, delegate, 0, 0)
	l.Title = "Select Worktree"
	l.Styles.Title = titleStyle
	l.SetFilteringEnabled(true)

	// Create text inputs for create form
	var createInputs [3]textinput.Model
	createInputs[0] = textinput.New()
	createInputs[0].Placeholder = "worktree-name"
	createInputs[0].CharLimit = 64
	createInputs[0].Width = 30

	createInputs[1] = textinput.New()
	createInputs[1].Placeholder = "branch-name (optional)"
	createInputs[1].CharLimit = 64
	createInputs[1].Width = 40

	createInputs[2] = textinput.New()
	createInputs[2].Placeholder = "main"
	createInputs[2].CharLimit = 64
	createInputs[2].Width = 30
	createInputs[2].SetValue("main")

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
		createInputs:  createInputs,
		settingInput:  settingInput,
		settings:      settings,
		pruneSelected: make(map[int]bool),
		staleDays:     30,
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
		return m, nil

	case WorktreeListMsg:
		if msg.Err != nil {
			m.setMessage(fmt.Sprintf("Failed to refresh: %v", msg.Err), true)
			return m, nil
		}

		items := make([]list.Item, 0, len(msg.Items))
		m.items = make([]WorktreeItem, 0, len(msg.Items))

		for _, item := range msg.Items {
			if meta, ok := m.metaStore.Get(item.Name); ok {
				item.Metadata = meta
			}
			items = append(items, item)
			m.items = append(m.items, item)
		}

		m.worktreeList.SetItems(items)
		// Only clear message if it was "Refreshing..."
		if m.message == "Refreshing..." {
			m.message = ""
		}
		return m, nil

	case tea.KeyMsg:
		// Clear message on any key
		m.message = ""

		// Global keys
		switch msg.String() {
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
	if m.activeTab == TabManage && m.manageMode != ManageModeList {
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
	switch m.manageMode {
	case ManageModeList:
		return m.updateManageList(msg)
	case ManageModeCreate:
		return m.updateManageCreate(msg)
	case ManageModeDelete:
		return m.updateManageDelete(msg)
	case ManageModePrune:
		return m.updateManagePrune(msg)
	}
	return m, nil
}

func (m *managerModel) updateManageList(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "q", "esc":
		m.quitting = true
		return m, tea.Quit

	case "c", "n":
		m.manageMode = ManageModeCreate
		m.createFocus = 0
		m.createInputs[0].Focus()
		m.createInputs[0].SetValue("")
		m.createInputs[1].SetValue("")
		m.createInputs[2].SetValue("main")
		return m, nil

	case "d", "x":
		if i, ok := m.worktreeList.SelectedItem().(WorktreeItem); ok {
			if i.IsMain {
				m.setMessage("Cannot delete main worktree", true)
				return m, nil
			}
			m.manageMode = ManageModeDelete
			m.deleteTarget = &i
			m.forceDelete = false
		}
		return m, nil

	case "p":
		m.enterPruneMode()
		return m, nil

	case "r":
		return m.refreshList()
	}

	var cmd tea.Cmd
	m.worktreeList, cmd = m.worktreeList.Update(msg)
	return m, cmd
}

func (m *managerModel) updateManageCreate(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc":
		m.manageMode = ManageModeList
		m.blurCreateInputs()
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
		return m.createWorktree()

	case "ctrl+s":
		return m.createWorktree()
	}

	var cmd tea.Cmd
	m.createInputs[m.createFocus], cmd = m.createInputs[m.createFocus].Update(msg)
	return m, cmd
}

func (m *managerModel) updateManageDelete(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc", "n":
		m.manageMode = ManageModeList
		m.deleteTarget = nil
		return m, nil

	case "y":
		return m.deleteWorktree(false)

	case "f":
		return m.deleteWorktree(true)
	}

	return m, nil
}

func (m *managerModel) updateManagePrune(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc", "q":
		m.manageMode = ManageModeList
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
		m.pruneSelected[m.pruneCursor] = !m.pruneSelected[m.pruneCursor]

	case "a":
		for i := range m.staleItems {
			m.pruneSelected[i] = true
		}

	case "n":
		for i := range m.staleItems {
			m.pruneSelected[i] = false
		}

	case "enter":
		return m.executePrune()
	}

	return m, nil
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
func (m *managerModel) updateCreateFocus() {
	m.blurCreateInputs()
	m.createInputs[m.createFocus].Focus()
}

func (m *managerModel) blurCreateInputs() {
	for i := range m.createInputs {
		m.createInputs[i].Blur()
	}
}

func (m *managerModel) createWorktree() (tea.Model, tea.Cmd) {
	name := strings.TrimSpace(m.createInputs[0].Value())
	if name == "" {
		m.setMessage("Name is required", true)
		return m, nil
	}

	branch := strings.TrimSpace(m.createInputs[1].Value())
	if branch == "" {
		branch = name
	}

	base := strings.TrimSpace(m.createInputs[2].Value())
	if base == "" {
		base = "main"
	}

	path, err := m.gitMgr.Add(name, branch, base)
	if err != nil {
		m.setMessage(fmt.Sprintf("Failed to create: %v", err), true)
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
		m.setMessage(fmt.Sprintf("Warning: metadata save failed: %v", err), true)
	}

	m.manageMode = ManageModeList
	m.blurCreateInputs()
	m.setMessage(fmt.Sprintf("✓ Created worktree: %s", name), false)

	return m.refreshList()
}

func (m *managerModel) deleteWorktree(force bool) (tea.Model, tea.Cmd) {
	if m.deleteTarget == nil {
		return m, nil
	}

	name := m.deleteTarget.Name

	if !force {
		clean, err := m.gitMgr.IsClean(m.deleteTarget.Path)
		if err != nil {
			m.setMessage(fmt.Sprintf("Failed to check status: %v", err), true)
			m.manageMode = ManageModeList
			m.deleteTarget = nil
			return m, nil
		}

		if !clean {
			m.forceDelete = true
			return m, nil
		}
	}

	if err := m.gitMgr.Remove(name, force); err != nil {
		m.setMessage(fmt.Sprintf("Failed to remove: %v", err), true)
		m.manageMode = ManageModeList
		m.deleteTarget = nil
		return m, nil
	}

	m.metaStore.Remove(name)
	if err := m.metaStore.Save(); err != nil {
		m.setMessage(fmt.Sprintf("Warning: metadata save failed: %v", err), true)
	}

	m.manageMode = ManageModeList
	m.deleteTarget = nil
	m.setMessage(fmt.Sprintf("✓ Removed worktree: %s", name), false)

	return m.refreshList()
}

func (m *managerModel) enterPruneMode() {
	staleNames := m.metaStore.GetStale(m.staleDays)

	if len(staleNames) == 0 {
		m.setMessage(fmt.Sprintf("No stale worktrees found (>%d days)", m.staleDays), false)
		return
	}

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

	m.manageMode = ManageModePrune
	m.pruneCursor = 0
	m.pruneSelected = make(map[int]bool)
	for i := range m.staleItems {
		m.pruneSelected[i] = true
	}
}

func (m *managerModel) executePrune() (tea.Model, tea.Cmd) {
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
		m.setMessage(fmt.Sprintf("Warning: metadata save failed: %v", err), true)
	}

	m.manageMode = ManageModeList
	m.staleItems = nil
	m.setMessage(fmt.Sprintf("✓ Removed %d worktree(s)", removed), false)

	return m.refreshList()
}

func fetchWorktreesCmd(gitMgr *git.Manager) tea.Cmd {
	return func() tea.Msg {
		worktrees, err := gitMgr.List()
		if err != nil {
			return WorktreeListMsg{Err: err}
		}

		items := make([]WorktreeItem, 0, len(worktrees))

		statuses := gitMgr.GetStatuses(worktrees)

		for _, wt := range worktrees {
		  status := statuses[wt.Path]

			item := WorktreeItem{
				Name:   wt.Name,
				Path:   wt.Path,
				Branch: wt.Branch,
				Status: status,
				IsMain: wt.IsMain,
				// Metadata will be attached in Update
			}

			items = append(items, item)
		}
		return WorktreeListMsg{Items: items}
	}
}

func (m *managerModel) refreshList() (tea.Model, tea.Cmd) {
	m.setMessage("Refreshing...", false)
	return m, fetchWorktreesCmd(m.gitMgr)
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
	m.message = msg
	if isError {
		m.messageStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#FF6B6B"))
	} else {
		m.messageStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#04B575"))
	}
}

func (m *managerModel) View() string {
	if m.quitting {
		return ""
	}

	var b strings.Builder

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
	if m.message != "" {
		b.WriteString("\n")
		b.WriteString(m.messageStyle.Render(m.message))
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
	switch m.manageMode {
	case ManageModeCreate:
		return m.viewCreateForm()
	case ManageModeDelete:
		return m.viewDeleteConfirm()
	case ManageModePrune:
		return m.viewPruneMode()
	default:
		return m.viewManageList()
	}
}

func (m *managerModel) viewManageList() string {
	var b strings.Builder
	b.WriteString(m.worktreeList.View())
	b.WriteString("\n")
	b.WriteString(helpStyle.Render("c create • d delete • p prune • r refresh • q quit"))
	return b.String()
}

func (m *managerModel) viewCreateForm() string {
	var b strings.Builder

	title := titleStyle.Render("📁 Create New Worktree")
	b.WriteString(title)
	b.WriteString("\n\n")

	labels := []string{"Name:", "Branch (optional):", "Base branch:"}
	for i, label := range labels {
		labelStyle := lipgloss.NewStyle()
		if i == m.createFocus {
			labelStyle = labelStyle.Foreground(lipgloss.Color("#7D56F4")).Bold(true)
			label = "▸ " + label
		} else {
			label = "  " + label
		}
		b.WriteString(fmt.Sprintf("%s\n  %s\n\n", labelStyle.Render(label), m.createInputs[i].View()))
	}

	b.WriteString(helpStyle.Render("tab next • ctrl+s create • esc cancel"))
	return b.String()
}

func (m *managerModel) viewDeleteConfirm() string {
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
			b.WriteString("f force delete • n/esc cancel\n")
		} else {
			b.WriteString("Delete this worktree?\n\n")
			b.WriteString("y confirm • n/esc cancel\n")
		}
	}

	return b.String()
}

func (m *managerModel) viewPruneMode() string {
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

	selected := 0
	for _, v := range m.pruneSelected {
		if v {
			selected++
		}
	}

	b.WriteString(fmt.Sprintf("\n%d of %d selected\n", selected, len(m.staleItems)))
	b.WriteString(helpStyle.Render("\n↑/↓ navigate • space toggle • a all • n none • enter delete • esc cancel"))

	return b.String()
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
