package tui

import (
	"github.com/charmbracelet/lipgloss"
)

var (
	// titleStyle styles the application title
	titleStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#00D9FF")).
		MarginBottom(1)

	// filenameStyle styles the filename display
	filenameStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FFEB3B"))

	// statusStyle styles the status display
	statusStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#4CAF50"))

	// progressStyle styles the progress bar
	progressStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#00D9FF"))

	// successStyle styles success messages
	successStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#4CAF50"))

	// errorStyle styles error messages
	errorStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#F44336"))

	// helpStyle styles help text
	helpStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#757575"))
)
