package tui

import (
	"github.com/charmbracelet/lipgloss"
)

var (
	// Color palette
	primaryColor    = lipgloss.Color("#00D9FF")
	secondaryColor  = lipgloss.Color("#FF6B9D")
	successColor   = lipgloss.Color("#4CAF50")
	warningColor   = lipgloss.Color("#FFC107")
	errorColor     = lipgloss.Color("#F44336")
	textColor      = lipgloss.Color("#E0E0E0")
	dimTextColor   = lipgloss.Color("#757575")
	borderColor    = lipgloss.Color("#424242")
	backgroundColor = lipgloss.Color("#1E1E1E")

	// Title style
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(primaryColor).
			Background(backgroundColor).
			Padding(0, 1).
			MarginBottom(1)

	// Box/border styles
	boxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(borderColor).
			Padding(1, 2).
			Margin(0, 1)

	sectionStyle = lipgloss.NewStyle().
			Border(lipgloss.NormalBorder(), false, false, false, true).
			BorderForeground(borderColor).
			PaddingLeft(1).
			MarginLeft(1)

	// Text styles
	filenameStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(textColor).
			Background(backgroundColor)

	statusStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(successColor)

	statusActiveStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(primaryColor)

	statusCompleteStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(successColor)

	statusErrorStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(errorColor)

	// Progress bar styles
	progressBarStyle = lipgloss.NewStyle().
				Foreground(primaryColor)

	progressBarBgStyle = lipgloss.NewStyle().
				Foreground(dimTextColor)

	// Stat styles
	statLabelStyle = lipgloss.NewStyle().
			Foreground(dimTextColor).
			Width(12)

	statValueStyle = lipgloss.NewStyle().
			Foreground(textColor).
			Bold(true)

	statValueHighlightStyle = lipgloss.NewStyle().
				Foreground(primaryColor).
				Bold(true)

	// Success/Error styles
	successStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(successColor).
			Background(backgroundColor).
			Padding(1, 2)

	errorStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(errorColor).
			Background(backgroundColor).
			Padding(1, 2)

	// Help text style
	helpStyle = lipgloss.NewStyle().
			Foreground(dimTextColor).
			Italic(true)

	// Graph style
	graphStyle = lipgloss.NewStyle().
			Foreground(primaryColor)
)
