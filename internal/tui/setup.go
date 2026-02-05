package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/darkLord19/wtx/internal/config"
	"github.com/darkLord19/wtx/internal/editor"
)

// SetupStep represents a step in the setup wizard
type SetupStep int

const (
	StepWelcome SetupStep = iota
	StepEditor
	StepCustomEditor
	StepWorktreeDir
	StepReuseWindow
	StepComplete
)

// setupModel is the TUI model for first-run setup
type setupModel struct {
	step       SetupStep
	config     *config.Config
	edDetector *editor.Detector
	width      int
	height     int
	quitting   bool

	// Editor selection
	editorOptions []string
	editorCursor  int

	// Text inputs
	customEditorInput textinput.Model
	worktreeDirInput  textinput.Model

	// Boolean selections
	reuseWindowValue bool
}

// NewSetupModel creates a new setup wizard model
func NewSetupModel(cfg *config.Config, edDetector *editor.Detector) *setupModel {
	// Get available editors
	editors := edDetector.DetectAll()
	editorOptions := []string{"(auto-detect)"}
	for _, ed := range editors {
		editorOptions = append(editorOptions, strings.ToLower(ed.Name()))
	}
	editorOptions = append(editorOptions, "(custom)")

	// Custom editor input
	customEditorInput := textinput.New()
	customEditorInput.Placeholder = "e.g., subl, atom, emacs"
	customEditorInput.CharLimit = 64
	customEditorInput.Width = 40

	// Worktree directory input
	worktreeDirInput := textinput.New()
	worktreeDirInput.Placeholder = "../worktrees"
	worktreeDirInput.CharLimit = 128
	worktreeDirInput.Width = 40
	worktreeDirInput.SetValue("../worktrees")

	return &setupModel{
		step:              StepWelcome,
		config:            cfg,
		edDetector:        edDetector,
		editorOptions:     editorOptions,
		customEditorInput: customEditorInput,
		worktreeDirInput:  worktreeDirInput,
		reuseWindowValue:  true,
	}
}

func (m *setupModel) Init() tea.Cmd {
	return nil
}

func (m *setupModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case tea.KeyMsg:
		switch m.step {
		case StepWelcome:
			return m.handleWelcomeKeys(msg)
		case StepEditor:
			return m.handleEditorKeys(msg)
		case StepCustomEditor:
			return m.handleCustomEditorKeys(msg)
		case StepWorktreeDir:
			return m.handleWorktreeDirKeys(msg)
		case StepReuseWindow:
			return m.handleReuseWindowKeys(msg)
		case StepComplete:
			return m.handleCompleteKeys(msg)
		}
	}

	// Update text inputs
	var cmd tea.Cmd
	switch m.step {
	case StepCustomEditor:
		m.customEditorInput, cmd = m.customEditorInput.Update(msg)
	case StepWorktreeDir:
		m.worktreeDirInput, cmd = m.worktreeDirInput.Update(msg)
	}

	return m, cmd
}

func (m *setupModel) handleWelcomeKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c", "q":
		m.quitting = true
		return m, tea.Quit
	case "enter", " ":
		m.step = StepEditor
	}
	return m, nil
}

func (m *setupModel) handleEditorKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c", "q":
		m.quitting = true
		return m, tea.Quit
	case "up", "k":
		if m.editorCursor > 0 {
			m.editorCursor--
		}
	case "down", "j":
		if m.editorCursor < len(m.editorOptions)-1 {
			m.editorCursor++
		}
	case "enter", " ":
		selected := m.editorOptions[m.editorCursor]
		if selected == "(custom)" {
			m.step = StepCustomEditor
			m.customEditorInput.Focus()
		} else {
			if selected == "(auto-detect)" {
				m.config.Editor = ""
			} else {
				m.config.Editor = selected
			}
			m.step = StepWorktreeDir
			m.worktreeDirInput.Focus()
		}
	}
	return m, nil
}

func (m *setupModel) handleCustomEditorKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c":
		m.quitting = true
		return m, tea.Quit
	case "esc":
		m.step = StepEditor
		m.customEditorInput.Blur()
		return m, nil
	case "enter":
		value := strings.TrimSpace(m.customEditorInput.Value())
		if value != "" {
			m.config.Editor = value
		}
		m.customEditorInput.Blur()
		m.step = StepWorktreeDir
		m.worktreeDirInput.Focus()
		return m, nil
	}

	var cmd tea.Cmd
	m.customEditorInput, cmd = m.customEditorInput.Update(msg)
	return m, cmd
}

func (m *setupModel) handleWorktreeDirKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c":
		m.quitting = true
		return m, tea.Quit
	case "esc":
		m.step = StepEditor
		m.worktreeDirInput.Blur()
		return m, nil
	case "enter":
		value := strings.TrimSpace(m.worktreeDirInput.Value())
		if value != "" {
			m.config.WorktreeDir = value
		}
		m.worktreeDirInput.Blur()
		m.step = StepReuseWindow
		return m, nil
	}

	var cmd tea.Cmd
	m.worktreeDirInput, cmd = m.worktreeDirInput.Update(msg)
	return m, cmd
}

func (m *setupModel) handleReuseWindowKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c", "q":
		m.quitting = true
		return m, tea.Quit
	case "left", "right", "h", "l", " ":
		m.reuseWindowValue = !m.reuseWindowValue
	case "enter":
		m.config.ReuseWindow = m.reuseWindowValue
		m.step = StepComplete
		// Save config
		m.config.Save()
	}
	return m, nil
}

func (m *setupModel) handleCompleteKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c", "q", "enter", " ":
		m.quitting = true
		return m, tea.Quit
	}
	return m, nil
}

func (m *setupModel) View() string {
	if m.quitting && m.step != StepComplete {
		return ""
	}

	var b strings.Builder

	// Header
	header := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#7D56F4")).
		Padding(1, 0).
		Render("🚀 wtx Setup Wizard")
	b.WriteString(header)
	b.WriteString("\n\n")

	// Progress indicator
	progress := m.renderProgress()
	b.WriteString(progress)
	b.WriteString("\n\n")

	// Step content
	switch m.step {
	case StepWelcome:
		b.WriteString(m.viewWelcome())
	case StepEditor:
		b.WriteString(m.viewEditor())
	case StepCustomEditor:
		b.WriteString(m.viewCustomEditor())
	case StepWorktreeDir:
		b.WriteString(m.viewWorktreeDir())
	case StepReuseWindow:
		b.WriteString(m.viewReuseWindow())
	case StepComplete:
		b.WriteString(m.viewComplete())
	}

	return b.String()
}

func (m *setupModel) renderProgress() string {
	steps := []string{"Welcome", "Editor", "Directory", "Window", "Done"}
	currentStep := int(m.step)
	if m.step == StepCustomEditor {
		currentStep = 1 // Custom editor is part of editor step
	}

	var parts []string
	for i, step := range steps {
		var rendered string
		if i < currentStep {
			rendered = lipgloss.NewStyle().Foreground(lipgloss.Color("#04B575")).Render("✓ " + step)
		} else if i == currentStep {
			rendered = lipgloss.NewStyle().Foreground(lipgloss.Color("#7D56F4")).Bold(true).Render("● " + step)
		} else {
			rendered = lipgloss.NewStyle().Foreground(lipgloss.Color("#626262")).Render("○ " + step)
		}
		parts = append(parts, rendered)
	}

	return strings.Join(parts, "  →  ")
}

func (m *setupModel) viewWelcome() string {
	var b strings.Builder

	b.WriteString(lipgloss.NewStyle().Bold(true).Render("Welcome to wtx!"))
	b.WriteString("\n\n")
	b.WriteString("wtx makes Git worktrees feel like instant \"workspace tabs\" across editors.\n\n")
	b.WriteString("Let's configure a few settings to get you started.\n\n")
	b.WriteString(helpStyle.Render("Press enter to continue • q to quit"))

	return b.String()
}

func (m *setupModel) viewEditor() string {
	var b strings.Builder

	b.WriteString(lipgloss.NewStyle().Bold(true).Render("Select your preferred editor:"))
	b.WriteString("\n\n")

	for i, opt := range m.editorOptions {
		cursor := "  "
		if i == m.editorCursor {
			cursor = "▸ "
		}

		style := lipgloss.NewStyle()
		if i == m.editorCursor {
			style = style.Foreground(lipgloss.Color("#7D56F4"))
		}

		b.WriteString(fmt.Sprintf("%s%s\n", cursor, style.Render(opt)))
	}

	b.WriteString("\n")
	b.WriteString(helpStyle.Render("↑/↓ navigate • enter select • q quit"))

	return b.String()
}

func (m *setupModel) viewCustomEditor() string {
	var b strings.Builder

	b.WriteString(lipgloss.NewStyle().Bold(true).Render("Enter your custom editor command:"))
	b.WriteString("\n\n")
	b.WriteString(m.customEditorInput.View())
	b.WriteString("\n\n")
	b.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color("#626262")).Render(
		"This should be the command to launch your editor (e.g., 'subl', 'atom', 'emacs')"))
	b.WriteString("\n\n")
	b.WriteString(helpStyle.Render("enter confirm • esc go back"))

	return b.String()
}

func (m *setupModel) viewWorktreeDir() string {
	var b strings.Builder

	b.WriteString(lipgloss.NewStyle().Bold(true).Render("Where should worktrees be created?"))
	b.WriteString("\n\n")
	b.WriteString(m.worktreeDirInput.View())
	b.WriteString("\n\n")
	b.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color("#626262")).Render(
		"This is relative to your repository root. Default: ../worktrees"))
	b.WriteString("\n\n")
	b.WriteString(helpStyle.Render("enter confirm • esc go back"))

	return b.String()
}

func (m *setupModel) viewReuseWindow() string {
	var b strings.Builder

	b.WriteString(lipgloss.NewStyle().Bold(true).Render("Reuse existing editor window?"))
	b.WriteString("\n\n")

	yesStyle := lipgloss.NewStyle().Padding(0, 2)
	noStyle := lipgloss.NewStyle().Padding(0, 2)

	if m.reuseWindowValue {
		yesStyle = yesStyle.Background(lipgloss.Color("#7D56F4")).Foreground(lipgloss.Color("#FFFFFF"))
	} else {
		noStyle = noStyle.Background(lipgloss.Color("#7D56F4")).Foreground(lipgloss.Color("#FFFFFF"))
	}

	b.WriteString(lipgloss.JoinHorizontal(lipgloss.Center,
		yesStyle.Render("Yes"),
		"  ",
		noStyle.Render("No"),
	))
	b.WriteString("\n\n")
	b.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color("#626262")).Render(
		"If yes, opening a worktree will reuse your current editor window.\nIf no, a new window will be opened each time."))
	b.WriteString("\n\n")
	b.WriteString(helpStyle.Render("←/→ toggle • enter confirm"))

	return b.String()
}

func (m *setupModel) viewComplete() string {
	var b strings.Builder

	b.WriteString(lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#04B575")).
		Render("✓ Setup Complete!"))
	b.WriteString("\n\n")

	b.WriteString("Your configuration:\n\n")

	editorDisplay := m.config.Editor
	if editorDisplay == "" {
		editorDisplay = "(auto-detect)"
	}

	b.WriteString(fmt.Sprintf("  Editor:         %s\n", editorDisplay))
	b.WriteString(fmt.Sprintf("  Worktree Dir:   %s\n", m.config.WorktreeDir))
	b.WriteString(fmt.Sprintf("  Reuse Window:   %v\n", m.config.ReuseWindow))
	b.WriteString("\n")

	b.WriteString("You can change these settings anytime with:\n")
	b.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color("#7D56F4")).Render("  wtx config --tui"))
	b.WriteString("\n\n")

	b.WriteString(helpStyle.Render("Press enter to start using wtx!"))

	return b.String()
}

// IsSetupComplete returns true if setup was completed
func (m *setupModel) IsSetupComplete() bool {
	return m.step == StepComplete
}

// RunSetup starts the setup wizard
func RunSetup(cfg *config.Config, edDetector *editor.Detector) (bool, error) {
	m := NewSetupModel(cfg, edDetector)
	p := tea.NewProgram(m, tea.WithAltScreen())

	finalModel, err := p.Run()
	if err != nil {
		return false, err
	}

	if fm, ok := finalModel.(*setupModel); ok {
		return fm.IsSetupComplete(), nil
	}

	return false, nil
}
