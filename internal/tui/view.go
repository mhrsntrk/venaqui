package tui

import (
	"fmt"

	"github.com/dustin/go-humanize"
)

// View renders the UI
func (m Model) View() string {
	if m.err != nil {
		return errorStyle.Render(fmt.Sprintf("Error: %v\n", m.err)) + "\n"
	}

	if m.quitting {
		if m.status != nil && m.status.IsComplete() {
			return successStyle.Render("✓ Download complete!\n") + "\n"
		}
		return "Exiting...\n"
	}

	if m.status == nil {
		return "Initializing download...\n"
	}

	// Calculate progress
	progress := m.status.GetProgress()

	// Format sizes
	completed := humanize.Bytes(uint64(m.status.CompletedLength))
	total := humanize.Bytes(uint64(m.status.TotalLength))
	speed := humanize.Bytes(uint64(m.status.DownloadSpeed))

	// Build view
	s := titleStyle.Render("Venaqui - Download Manager") + "\n\n"
	s += fmt.Sprintf("File: %s\n", filenameStyle.Render(m.filename))
	s += fmt.Sprintf("Status: %s\n", statusStyle.Render(m.status.Status))
	s += fmt.Sprintf("Progress: %s\n", progressBar(progress))
	s += fmt.Sprintf("Size: %s / %s (%.2f%%)\n", completed, total, progress)
	s += fmt.Sprintf("Speed: %s/s\n", speed)
	s += "\n" + helpStyle.Render("Press 'q' to quit")

	return s
}

// progressBar creates a visual progress bar
func progressBar(percent float64) string {
	width := 50
	filled := int(percent / 100 * float64(width))

	bar := ""
	for i := 0; i < width; i++ {
		if i < filled {
			bar += "█"
		} else {
			bar += "░"
		}
	}

	return progressStyle.Render(bar)
}
