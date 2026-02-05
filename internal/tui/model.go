package tui

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/mhrsntrk/venaqui/internal/aria2"
)

// Model represents the TUI application state
type Model struct {
	aria2Client  *aria2.Client
	gid          string
	status       *aria2.DownloadStatus
	filename     string
	err          error
	quitting     bool
	startTime    time.Time
	completionTime time.Time // Time when download completed
	lastUpdate   time.Time
	speedHistory []int64 // Speed history for graph (last 50 samples)
	maxHistory   int
}

// tickMsg is sent periodically to update the UI
type tickMsg time.Time

// statusMsg wraps a download status update
type statusMsg *aria2.DownloadStatus

// errMsg wraps an error
type errMsg error

// InitialModel creates a new model with initial state
func InitialModel(aria2Client *aria2.Client, gid, filename string) Model {
	now := time.Now()
	return Model{
		aria2Client:  aria2Client,
		gid:          gid,
		filename:     filename,
		startTime:    now,
		lastUpdate:   now,
		speedHistory: make([]int64, 0),
		maxHistory:   50,
	}
}

// Init initializes the model and returns initial commands
func (m Model) Init() tea.Cmd {
	return tea.Batch(
		tickCmd(),
		m.fetchStatus,
	)
}

// tickCmd returns a command that sends a tick message after 1 second
func tickCmd() tea.Cmd {
	return tea.Tick(time.Second, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

// fetchStatus fetches the current download status
func (m Model) fetchStatus() tea.Msg {
	status, err := m.aria2Client.GetStatus(m.gid)
	if err != nil {
		return errMsg(err)
	}
	return statusMsg(status)
}
