package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// KeyHelp represents a keyboard shortcut and its description
type KeyHelp struct {
	Key  string
	Desc string
}

// HelpPanel represents a help panel with keyboard shortcuts
type HelpPanel struct {
	sections map[string][]KeyHelp
	visible  bool
}

// NewHelpPanel creates a new help panel
func NewHelpPanel() *HelpPanel {
	return &HelpPanel{
		sections: make(map[string][]KeyHelp),
		visible:  false,
	}
}

// AddSection adds a section with keyboard shortcuts
func (h *HelpPanel) AddSection(name string, shortcuts []KeyHelp) {
	h.sections[name] = shortcuts
}

// Toggle toggles the visibility of the help panel
func (h *HelpPanel) Toggle() {
	h.visible = !h.visible
}

// Show shows the help panel
func (h *HelpPanel) Show() {
	h.visible = true
}

// Hide hides the help panel
func (h *HelpPanel) Hide() {
	h.visible = false
}

// IsVisible returns whether the help panel is visible
func (h *HelpPanel) IsVisible() bool {
	return h.visible
}

// Render renders the help panel
func (h *HelpPanel) Render(width int) string {
	if !h.visible {
		return ""
	}

	var b strings.Builder

	// Title
	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#7D56F4")).
		Padding(1, 0)

	b.WriteString(titleStyle.Render("⌨️  Keyboard Shortcuts"))
	b.WriteString("\n\n")

	// Sections
	sectionOrder := []string{"Global", "Navigation", "Worktrees", "Manage", "Settings"}
	
	for _, sectionName := range sectionOrder {
		shortcuts, ok := h.sections[sectionName]
		if !ok || len(shortcuts) == 0 {
			continue
		}

		// Section title
		sectionStyle := lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#04B575"))

		b.WriteString(sectionStyle.Render(sectionName))
		b.WriteString("\n")

		// Shortcuts
		for _, shortcut := range shortcuts {
			keyStyle := lipgloss.NewStyle().
				Foreground(lipgloss.Color("#7D56F4")).
				Bold(true).
				Width(15)

			descStyle := lipgloss.NewStyle().
				Foreground(lipgloss.Color("#FFFFFF"))

			b.WriteString(fmt.Sprintf("  %s %s\n",
				keyStyle.Render(shortcut.Key),
				descStyle.Render(shortcut.Desc)))
		}
		b.WriteString("\n")
	}

	// Footer
	footerStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#626262")).
		Italic(true)

	b.WriteString(footerStyle.Render("Press ? to close help"))

	// Wrap in a box
	boxStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#7D56F4")).
		Padding(1, 2).
		Width(width - 4)

	return boxStyle.Render(b.String())
}

// GetGlobalHelp returns common global shortcuts
func GetGlobalHelp() []KeyHelp {
	return []KeyHelp{
		{"?", "Toggle help"},
		{"q / esc", "Quit / Cancel"},
		{"ctrl+c", "Force quit"},
	}
}

// GetNavigationHelp returns navigation shortcuts
func GetNavigationHelp() []KeyHelp {
	return []KeyHelp{
		{"↑ / k", "Move up"},
		{"↓ / j", "Move down"},
		{"enter", "Select"},
		{"/", "Search / Filter"},
	}
}

// GetWorktreesHelp returns worktree-specific shortcuts
func GetWorktreesHelp() []KeyHelp {
	return []KeyHelp{
		{"enter", "Open worktree"},
		{"1-3", "Switch tabs"},
		{"r", "Refresh list"},
	}
}

// GetManageHelp returns manage tab shortcuts
func GetManageHelp() []KeyHelp {
	return []KeyHelp{
		{"c / n", "Create worktree"},
		{"d / x", "Delete worktree"},
		{"p", "Prune stale"},
		{"r", "Refresh list"},
	}
}

// GetSettingsHelp returns settings tab shortcuts
func GetSettingsHelp() []KeyHelp {
	return []KeyHelp{
		{"enter", "Edit setting"},
		{"←/→", "Change option"},
		{"s", "Save settings"},
	}
}
