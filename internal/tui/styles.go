package tui

import "github.com/charmbracelet/lipgloss"

var (
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#7D56F4")).
			Padding(0, 1)

	helpStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#626262"))

	cleanStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#04B575"))

	dirtyStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FF6B6B"))

	portStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#7D56F4"))
)
