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

// SettingType represents the type of setting
type SettingType int

const (
	SettingEditor SettingType = iota
	SettingCustomEditor
	SettingReuseWindow
	SettingWorktreeDir
	SettingAutoStartDev
)

// SettingItem represents a configurable setting
type SettingItem struct {
	Name        string
	Description string
	Type        SettingType
	Value       string
	Options     []string // For selection-based settings
}

// settingsModel is the TUI model for settings configuration
type settingsModel struct {
	settings     []SettingItem
	cursor       int
	editing      bool
	textInput    textinput.Model
	config       *config.Config
	edDetector   *editor.Detector
	saved        bool
	err          error
	width        int
	height       int
	quitting     bool
	optionCursor int // For cycling through options
}

// NewSettingsModel creates a new settings TUI model
func NewSettingsModel(cfg *config.Config, edDetector *editor.Detector) *settingsModel {
	ti := textinput.New()
	ti.Placeholder = "Enter value..."
	ti.CharLimit = 256
	ti.Width = 40

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
			Description: "Preferred code editor (select or choose custom)",
			Type:        SettingEditor,
			Value:       cfg.Editor,
			Options:     editorOptions,
		},
		{
			Name:        "Custom Editor Command",
			Description: "Custom editor command (e.g., 'subl', 'atom', 'emacs')",
			Type:        SettingCustomEditor,
			Value:       cfg.Editor,
			Options:     nil, // Text input
		},
		{
			Name:        "Reuse Window",
			Description: "Reuse existing editor window when opening worktrees",
			Type:        SettingReuseWindow,
			Value:       boolToString(cfg.ReuseWindow),
			Options:     []string{"true", "false"},
		},
		{
			Name:        "Worktree Directory",
			Description: "Directory where new worktrees are created (relative to repo)",
			Type:        SettingWorktreeDir,
			Value:       cfg.WorktreeDir,
			Options:     nil, // Text input
		},
		{
			Name:        "Auto Start Dev",
			Description: "Automatically start dev server when opening worktree",
			Type:        SettingAutoStartDev,
			Value:       boolToString(cfg.AutoStartDev),
			Options:     []string{"true", "false"},
		},
	}

	return &settingsModel{
		settings:   settings,
		config:     cfg,
		edDetector: edDetector,
		textInput:  ti,
	}
}

func (m *settingsModel) Init() tea.Cmd {
	return nil
}

func (m *settingsModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case tea.KeyMsg:
		if m.editing {
			return m.handleEditingKeys(msg)
		}
		return m.handleNavigationKeys(msg)
	}

	if m.editing && m.settings[m.cursor].Options == nil {
		var cmd tea.Cmd
		m.textInput, cmd = m.textInput.Update(msg)
		return m, cmd
	}

	return m, nil
}

func (m *settingsModel) handleNavigationKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c", "q", "esc":
		m.quitting = true
		return m, tea.Quit

	case "up", "k":
		if m.cursor > 0 {
			m.cursor--
		}

	case "down", "j":
		if m.cursor < len(m.settings)-1 {
			m.cursor++
		}

	case "enter", " ":
		m.startEditing()

	case "s":
		// Save settings
		m.saveSettings()
	}

	return m, nil
}

func (m *settingsModel) handleEditingKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	setting := &m.settings[m.cursor]

	switch msg.String() {
	case "esc":
		m.editing = false
		m.textInput.Blur()
		return m, nil

	case "enter":
		if setting.Options == nil {
			// Text input - save the value
			setting.Value = m.textInput.Value()
		}
		m.editing = false
		m.textInput.Blur()
		return m, nil

	case "left", "h":
		if setting.Options != nil {
			// Cycle through options
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
			// Cycle through options
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
			m.textInput, cmd = m.textInput.Update(msg)
			return m, cmd
		}
	}

	return m, nil
}

func (m *settingsModel) startEditing() {
	m.editing = true
	setting := m.settings[m.cursor]

	if setting.Options != nil {
		// Find current option index
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
		// Text input
		m.textInput.SetValue(setting.Value)
		m.textInput.Focus()
	}
}

func (m *settingsModel) saveSettings() {
	var editorValue string
	var customEditorValue string

	// Update config from settings
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

	// Save to disk
	if err := m.config.Save(); err != nil {
		m.err = err
		m.saved = false
	} else {
		m.saved = true
		m.err = nil
	}
}

func (m *settingsModel) View() string {
	if m.quitting {
		return ""
	}

	var b strings.Builder

	// Title
	title := titleStyle.Render("⚙  Settings Configuration")
	b.WriteString(title)
	b.WriteString("\n\n")

	// Settings list
	for i, setting := range m.settings {
		cursor := "  "
		if i == m.cursor {
			cursor = "▸ "
		}

		// Setting name
		nameStyle := lipgloss.NewStyle().Bold(true)
		if i == m.cursor {
			nameStyle = nameStyle.Foreground(lipgloss.Color("#7D56F4"))
		}

		name := nameStyle.Render(setting.Name)

		// Value display
		var valueDisplay string
		if m.editing && i == m.cursor {
			if setting.Options != nil {
				// Show option selector
				valueDisplay = m.renderOptionSelector(setting)
			} else {
				// Show text input
				valueDisplay = m.textInput.View()
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

		// Description
		desc := lipgloss.NewStyle().
			Foreground(lipgloss.Color("#626262")).
			Render(setting.Description)

		b.WriteString(fmt.Sprintf("%s%s: %s\n", cursor, name, valueDisplay))
		b.WriteString(fmt.Sprintf("   %s\n\n", desc))
	}

	// Status message
	if m.saved {
		b.WriteString(lipgloss.NewStyle().
			Foreground(lipgloss.Color("#04B575")).
			Render("✓ Settings saved successfully!"))
		b.WriteString("\n")
	}
	if m.err != nil {
		b.WriteString(lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FF6B6B")).
			Render(fmt.Sprintf("✗ Error: %v", m.err)))
		b.WriteString("\n")
	}

	// Help
	var helpText string
	if m.editing {
		if m.settings[m.cursor].Options != nil {
			helpText = "←/→ change • enter confirm • esc cancel"
		} else {
			helpText = "type to edit • enter confirm • esc cancel"
		}
	} else {
		helpText = "↑/↓ navigate • enter edit • s save • q quit"
	}
	b.WriteString("\n")
	b.WriteString(helpStyle.Render(helpText))

	return b.String()
}

func (m *settingsModel) renderOptionSelector(setting SettingItem) string {
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

// RunSettings starts the settings TUI
func RunSettings(cfg *config.Config, edDetector *editor.Detector) error {
	m := NewSettingsModel(cfg, edDetector)
	p := tea.NewProgram(m, tea.WithAltScreen())

	_, err := p.Run()
	return err
}

func boolToString(b bool) string {
	if b {
		return "true"
	}
	return "false"
}
