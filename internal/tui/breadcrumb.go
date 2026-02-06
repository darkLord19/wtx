package tui

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// Breadcrumb represents a navigation breadcrumb
type Breadcrumb struct {
	items []string
}

// NewBreadcrumb creates a new breadcrumb
func NewBreadcrumb(items ...string) Breadcrumb {
	return Breadcrumb{items: items}
}

// Add adds an item to the breadcrumb
func (b *Breadcrumb) Add(item string) {
	b.items = append(b.items, item)
}

// Pop removes the last item from the breadcrumb
func (b *Breadcrumb) Pop() {
	if len(b.items) > 0 {
		b.items = b.items[:len(b.items)-1]
	}
}

// Render renders the breadcrumb
func (b Breadcrumb) Render() string {
	if len(b.items) == 0 {
		return ""
	}

	separator := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#626262")).
		Render(" > ")

	var parts []string
	for i, item := range b.items {
		style := lipgloss.NewStyle().Foreground(lipgloss.Color("#626262"))
		if i == len(b.items)-1 {
			// Last item is highlighted
			style = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#7D56F4")).
				Bold(true)
		}
		parts = append(parts, style.Render(item))
	}

	return strings.Join(parts, separator)
}

// String returns the breadcrumb as a plain string
func (b Breadcrumb) String() string {
	return strings.Join(b.items, " > ")
}

// Clear clears all breadcrumb items
func (b *Breadcrumb) Clear() {
	b.items = []string{}
}

// Set replaces all items with the given items
func (b *Breadcrumb) Set(items ...string) {
	b.items = items
}
