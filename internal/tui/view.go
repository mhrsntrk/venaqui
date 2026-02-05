package tui

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
	"github.com/dustin/go-humanize"
)

// View renders the UI
func (m Model) View() string {
	if m.err != nil {
		return errorStyle.Render(fmt.Sprintf("✗ Error: %v", m.err)) + "\n\n" +
			helpStyle.Render("Press 'q' to quit") + "\n"
	}

	if m.quitting {
		if m.status != nil && m.status.IsComplete() {
			return successStyle.Render("✓ Download complete!") + "\n\n" +
				helpStyle.Render("Press 'q' to quit") + "\n"
		}
		return "Exiting...\n"
	}

	// Show completion message but don't quit automatically
	if m.status != nil && m.status.IsComplete() {
		var s strings.Builder
		s.WriteString(titleStyle.Render("Venaqui - Download Manager"))
		s.WriteString("\n\n")
		s.WriteString(successStyle.Render("✓ Download complete!"))
		s.WriteString("\n\n")
		
		filePath := m.status.GetFilePath()
		fileDir := m.status.GetFileDirectory()
		if filePath != "" || fileDir != "" {
			pathToShow := filePath
			if pathToShow == "" {
				pathToShow = fileDir
			}
			s.WriteString(fmt.Sprintf("%s %s\n\n",
				statLabelStyle.Render("Location:"),
				filenameStyle.Render(pathToShow),
			))
		}
		
		s.WriteString(helpStyle.Render("Press 'o' to open file | 'd' to show in directory | 'q' to quit"))
		s.WriteString("\n")
		return s.String()
	}

	if m.status == nil {
		return titleStyle.Render("Venaqui - Download Manager") + "\n\n" +
			"Initializing download...\n"
	}

	// Build the view
	var s strings.Builder

	// Header
	s.WriteString(titleStyle.Render("Venaqui - Download Manager"))
	s.WriteString("\n\n")

	// File info box
	fileBox := boxStyle.Render(
		fmt.Sprintf("%s %s\n%s %s",
			statLabelStyle.Render("File:"),
			filenameStyle.Render(m.filename),
			statLabelStyle.Render("Status:"),
			m.getStatusStyle().Render(m.getStatusText()),
		),
	)
	s.WriteString(fileBox)
	s.WriteString("\n\n")

	// Progress section
	progress := m.status.GetProgress()
	progressBox := boxStyle.Render(
		fmt.Sprintf("%s\n%s\n%s",
			statLabelStyle.Render("Progress:"),
			m.renderProgressBar(progress),
			m.renderProgressText(progress),
		),
	)
	s.WriteString(progressBox)
	s.WriteString("\n\n")

	// Statistics section - two columns
	statsLeft := m.renderStatsLeft()
	statsRight := m.renderStatsRight()
	
	statsBox := boxStyle.Render(
		lipgloss.JoinHorizontal(lipgloss.Left,
			statsLeft,
			strings.Repeat(" ", 4),
			statsRight,
		),
	)
	s.WriteString(statsBox)
	s.WriteString("\n\n")

	// Speed graph
	if len(m.speedHistory) > 0 {
		graphBox := boxStyle.Render(
			fmt.Sprintf("%s\n%s",
				statLabelStyle.Render("Speed History:"),
				m.renderSpeedGraph(),
			),
		)
		s.WriteString(graphBox)
		s.WriteString("\n\n")
	}

	// Help text
	if m.status != nil && m.status.IsComplete() {
		s.WriteString(helpStyle.Render("Press 'o' to open file | 'd' to show in directory | 'q' to quit"))
	} else {
		s.WriteString(helpStyle.Render("Press 'q' or 'Esc' to quit"))
	}
	s.WriteString("\n")

	return s.String()
}

// getStatusStyle returns the appropriate style for the status
func (m Model) getStatusStyle() lipgloss.Style {
	switch m.status.Status {
	case "active":
		return statusActiveStyle
	case "complete":
		return statusCompleteStyle
	case "error":
		return statusErrorStyle
	default:
		return statusStyle
	}
}

// getStatusText returns a formatted status text
func (m Model) getStatusText() string {
	switch m.status.Status {
	case "active":
		return "● Active"
	case "complete":
		return "✓ Complete"
	case "error":
		return "✗ Error"
	case "waiting":
		return "⏳ Waiting"
	case "paused":
		return "⏸ Paused"
	default:
		// Capitalize first letter
		if len(m.status.Status) > 0 {
			return strings.ToUpper(string(m.status.Status[0])) + strings.ToLower(m.status.Status[1:])
		}
		return m.status.Status
	}
}

// renderProgressBar creates a visual progress bar
func (m Model) renderProgressBar(percent float64) string {
	width := 60
	filled := int(percent / 100 * float64(width))
	if filled > width {
		filled = width
	}

	var bar strings.Builder
	for i := 0; i < width; i++ {
		if i < filled {
			// Use gradient-like effect
			if i < filled-2 {
				bar.WriteString(progressBarStyle.Render("█"))
			} else {
				bar.WriteString(progressBarStyle.Render("▊"))
			}
		} else {
			bar.WriteString(progressBarBgStyle.Render("░"))
		}
	}

	return bar.String()
}

// renderProgressText renders progress percentage and size info
func (m Model) renderProgressText(progress float64) string {
	completed := humanize.Bytes(uint64(m.status.CompletedLength))
	total := humanize.Bytes(uint64(m.status.TotalLength))
	return fmt.Sprintf("%s / %s (%.1f%%)",
		statValueStyle.Render(completed),
		statValueStyle.Render(total),
		progress,
	)
}

// renderStatsLeft renders left column of statistics
func (m Model) renderStatsLeft() string {
	var s strings.Builder
	
	// Download speed
	speed := humanize.Bytes(uint64(m.status.DownloadSpeed))
	s.WriteString(fmt.Sprintf("%s %s/s\n",
		statLabelStyle.Render("Speed:"),
		statValueHighlightStyle.Render(speed),
	))

	// Upload speed
	uploadSpeed := humanize.Bytes(uint64(m.status.UploadSpeed))
	s.WriteString(fmt.Sprintf("%s %s/s\n",
		statLabelStyle.Render("Upload:"),
		statValueStyle.Render(uploadSpeed),
	))

	// Connections
	s.WriteString(fmt.Sprintf("%s %s\n",
		statLabelStyle.Render("Connections:"),
		statValueStyle.Render(fmt.Sprintf("%d", m.status.Connections)),
	))

	return s.String()
}

// renderStatsRight renders right column of statistics
func (m Model) renderStatsRight() string {
	var s strings.Builder

	// Elapsed time
	elapsed := time.Since(m.startTime)
	elapsedStr := formatDuration(elapsed)
	s.WriteString(fmt.Sprintf("%s %s\n",
		statLabelStyle.Render("Elapsed:"),
		statValueStyle.Render(elapsedStr),
	))

	// ETA
	eta := m.status.GetETA()
	if eta > 0 {
		etaDuration := time.Duration(eta) * time.Second
		etaStr := formatDuration(etaDuration)
		s.WriteString(fmt.Sprintf("%s %s\n",
			statLabelStyle.Render("ETA:"),
			statValueStyle.Render(etaStr),
		))
	} else {
		s.WriteString(fmt.Sprintf("%s %s\n",
			statLabelStyle.Render("ETA:"),
			statValueStyle.Render("Calculating..."),
		))
	}

	// Remaining
	remaining := m.status.TotalLength - m.status.CompletedLength
	if remaining > 0 {
		remainingStr := humanize.Bytes(uint64(remaining))
		s.WriteString(fmt.Sprintf("%s %s\n",
			statLabelStyle.Render("Remaining:"),
			statValueStyle.Render(remainingStr),
		))
	} else {
		s.WriteString(fmt.Sprintf("%s %s\n",
			statLabelStyle.Render("Remaining:"),
			statValueStyle.Render("0 B"),
		))
	}

	return s.String()
}

// renderSpeedGraph renders a simple ASCII speed graph
func (m Model) renderSpeedGraph() string {
	if len(m.speedHistory) == 0 {
		return "No data yet..."
	}

	width := 60
	height := 8

	// Find max speed for scaling
	maxSpeed := int64(1)
	for _, speed := range m.speedHistory {
		if speed > maxSpeed {
			maxSpeed = speed
		}
	}
	if maxSpeed == 0 {
		maxSpeed = 1
	}

	// Create graph grid
	graph := make([][]bool, height)
	for i := range graph {
		graph[i] = make([]bool, width)
	}

	// Plot points
	step := float64(len(m.speedHistory)) / float64(width)
	for x := 0; x < width; x++ {
		idx := int(float64(x) * step)
		if idx >= len(m.speedHistory) {
			idx = len(m.speedHistory) - 1
		}
		if idx < 0 {
			idx = 0
		}

		speed := m.speedHistory[idx]
		y := int(float64(speed) / float64(maxSpeed) * float64(height-1))
		if y >= height {
			y = height - 1
		}
		if y < 0 {
			y = 0
		}

		// Draw line from bottom
		for lineY := height - 1; lineY >= height-1-y; lineY-- {
			if lineY >= 0 && lineY < height {
				graph[lineY][x] = true
			}
		}
	}

	// Render graph
	var result strings.Builder
	for y := height - 1; y >= 0; y-- {
		for x := 0; x < width; x++ {
			if graph[y][x] {
				result.WriteString(graphStyle.Render("▁"))
			} else {
				result.WriteString(" ")
			}
		}
		if y == height-1 {
			result.WriteString(fmt.Sprintf(" %s/s", humanize.Bytes(uint64(maxSpeed))))
		}
		result.WriteString("\n")
	}

	return result.String()
}

// formatDuration formats a duration in a human-readable way
func formatDuration(d time.Duration) string {
	if d < time.Minute {
		return fmt.Sprintf("%ds", int(d.Seconds()))
	}
	if d < time.Hour {
		minutes := int(d.Minutes())
		seconds := int(d.Seconds()) % 60
		return fmt.Sprintf("%dm %ds", minutes, seconds)
	}
	hours := int(d.Hours())
	minutes := int(d.Minutes()) % 60
	return fmt.Sprintf("%dh %dm", hours, minutes)
}
