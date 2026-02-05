package tui

import (
	"github.com/charmbracelet/lipgloss"
)

var (
	// Color palette - harmonized around primary red (#EF4444)
	primaryColor    = lipgloss.Color("#EF4444") // Red primary
	secondaryColor  = lipgloss.Color("#F87171") // Lighter red/pink accent
	successColor   = lipgloss.Color("#10B981") // Emerald green (complements red)
	warningColor   = lipgloss.Color("#F59E0B") // Amber (works with red)
	errorColor     = lipgloss.Color("#DC2626") // Darker red for errors
	textColor      = lipgloss.Color("#F3F4F6") // Light gray text
	dimTextColor   = lipgloss.Color("#9CA3AF") // Medium gray for dim text
	borderColor    = lipgloss.Color("#4B5563") // Darker gray for borders
	backgroundColor = lipgloss.Color("#1F2937") // Dark background

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
