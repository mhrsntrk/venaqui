package aria2

import (
	"fmt"

	"github.com/siku2/arigo"
)

// DownloadStatus represents the status of a download
type DownloadStatus struct {
	GID             string
	Status          string
	TotalLength     int64
	CompletedLength int64
	DownloadSpeed   int64
	Files           []arigo.File
	ErrorMessage    string
}

// GetStatus retrieves the status of a download by GID
func (c *Client) GetStatus(gid string) (*DownloadStatus, error) {
	status, err := c.rpc.TellStatus(gid,
		"gid", "status", "totalLength", "completedLength",
		"downloadSpeed", "files", "errorMessage")
	if err != nil {
		return nil, fmt.Errorf("failed to get status: %w", err)
	}

	ds := &DownloadStatus{
		GID:             status.GID,
		Status:          string(status.Status),
		TotalLength:     int64(status.TotalLength),
		CompletedLength: int64(status.CompletedLength),
		DownloadSpeed:   int64(status.DownloadSpeed),
		Files:           status.Files,
		ErrorMessage:    status.ErrorMessage,
	}

	return ds, nil
}

// GetProgress returns the download progress as a percentage (0-100)
func (ds *DownloadStatus) GetProgress() float64 {
	if ds.TotalLength == 0 {
		return 0
	}
	return float64(ds.CompletedLength) / float64(ds.TotalLength) * 100
}

// IsComplete returns true if the download is complete
func (ds *DownloadStatus) IsComplete() bool {
	return ds.Status == "complete"
}

// IsError returns true if the download has an error
func (ds *DownloadStatus) IsError() bool {
	return ds.Status == "error"
}

// IsActive returns true if the download is active
func (ds *DownloadStatus) IsActive() bool {
	return ds.Status == "active"
}
