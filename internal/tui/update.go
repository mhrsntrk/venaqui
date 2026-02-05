package tui

import (
	"fmt"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

// Update handles messages and updates the model
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c", "esc":
			m.quitting = true
			return m, tea.Quit
		}

	case tickMsg:
		// Update last update time
		m.lastUpdate = time.Time(msg)
		// Fetch status and schedule next tick
		return m, tea.Batch(
			tickCmd(),
			m.fetchStatus,
		)

	case statusMsg:
		m.status = msg

		// Check if download is complete
		if m.status.IsComplete() {
			m.quitting = true
			return m, tea.Quit
		}

		// Check for errors
		if m.status.IsError() {
			if m.status.ErrorMessage != "" {
				m.err = fmt.Errorf("download error: %s", m.status.ErrorMessage)
			} else {
				m.err = fmt.Errorf("download error")
			}
			m.quitting = true
			return m, tea.Quit
		}

		return m, nil

	case errMsg:
		m.err = msg
		m.quitting = true
		return m, tea.Quit
	}

	return m, nil
}
