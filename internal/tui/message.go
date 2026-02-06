package tui

import "github.com/charmbracelet/lipgloss"

// Message represents a styled message to display to the user
type Message struct {
	text  string
	style lipgloss.Style
}

// NewSuccessMessage creates a green success message
func NewSuccessMessage(text string) Message {
	return Message{
		text:  text,
		style: lipgloss.NewStyle().Foreground(lipgloss.Color("#04B575")),
	}
}

// NewErrorMessage creates a red error message
func NewErrorMessage(text string) Message {
	return Message{
		text:  text,
		style: lipgloss.NewStyle().Foreground(lipgloss.Color("#FF6B6B")),
	}
}

// NewInfoMessage creates a purple info message
func NewInfoMessage(text string) Message {
	return Message{
		text:  text,
		style: lipgloss.NewStyle().Foreground(lipgloss.Color("#7D56F4")),
	}
}

// NewWarningMessage creates an orange warning message
func NewWarningMessage(text string) Message {
	return Message{
		text:  text,
		style: lipgloss.NewStyle().Foreground(lipgloss.Color("#FFA500")),
	}
}

// Render returns the styled message string
func (m Message) Render() string {
	if m.text == "" {
		return ""
	}
	return m.style.Render(m.text)
}

// Text returns the raw message text
func (m Message) Text() string {
	return m.text
}

// IsEmpty returns true if the message is empty
func (m Message) IsEmpty() bool {
	return m.text == ""
}

// SetText updates the message text
func (m *Message) SetText(text string) {
	m.text = text
}

// Clear clears the message
func (m *Message) Clear() {
	m.text = ""
}
