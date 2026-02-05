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
		case "o":
			// Open file directly with default application when download is complete
			if m.status != nil && m.status.IsComplete() {
				filePath := m.status.GetFilePath()
				
				if filePath != "" {
					if err := openFile(filePath); err != nil {
						m.err = fmt.Errorf("failed to open file: %v", err)
					}
					// Still allow user to quit after opening
					return m, nil
				} else {
					m.err = fmt.Errorf("could not determine file path")
				}
			}
			return m, nil
		case "d", "s":
			// Show file in directory (reveal in Finder/Explorer) when download is complete
			if m.status != nil && m.status.IsComplete() {
				filePath := m.status.GetFilePath()
				fileDir := m.status.GetFileDirectory()
				
				var err error
				if filePath != "" {
					// Reveal the file in Finder/Explorer (highlights the file)
					err = revealFileInDirectory(filePath)
				} else if fileDir != "" {
					// Fallback to opening the directory
					err = openDirectory(fileDir)
				} else {
					m.err = fmt.Errorf("could not determine download directory")
					return m, nil
				}
				
				if err != nil {
					m.err = fmt.Errorf("failed to open directory: %v", err)
				}
				// Still allow user to quit after opening
				return m, nil
			}
			return m, nil
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

		// Update speed history
		if m.status != nil && m.status.DownloadSpeed > 0 {
			m.speedHistory = append(m.speedHistory, m.status.DownloadSpeed)
			if len(m.speedHistory) > m.maxHistory {
				m.speedHistory = m.speedHistory[1:]
			}
		}

		// Track completion time
		if m.status != nil && m.status.IsComplete() && m.completionTime.IsZero() {
			m.completionTime = time.Now()
		}

		// Don't auto-quit on completion - let user press 'o' or 'q'
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
