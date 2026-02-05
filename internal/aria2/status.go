package aria2

import (
	"fmt"
	"path/filepath"

	"github.com/siku2/arigo"
)

// DownloadStatus represents the status of a download
type DownloadStatus struct {
	GID             string
	Status          string
	TotalLength     int64
	CompletedLength int64
	DownloadSpeed   int64
	UploadSpeed     int64
	Connections     int
	NumPieces       int
	PieceLength     int64
	Dir             string
	Files           []arigo.File
	ErrorMessage    string
}

// GetStatus retrieves the status of a download by GID
func (c *Client) GetStatus(gid string) (*DownloadStatus, error) {
	status, err := c.rpc.TellStatus(gid,
		"gid", "status", "totalLength", "completedLength",
		"downloadSpeed", "uploadSpeed", "connections", "numPieces",
		"pieceLength", "dir", "files", "errorMessage")
	if err != nil {
		return nil, fmt.Errorf("failed to get status: %w", err)
	}

	ds := &DownloadStatus{
		GID:             status.GID,
		Status:          string(status.Status),
		TotalLength:     int64(status.TotalLength),
		CompletedLength: int64(status.CompletedLength),
		DownloadSpeed:   int64(status.DownloadSpeed),
		UploadSpeed:     int64(status.UploadSpeed),
		Connections:     int(status.Connections),
		NumPieces:       int(status.NumPieces),
		PieceLength:     int64(status.PieceLength),
		Dir:             status.Dir,
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

// GetETA returns the estimated time to completion in seconds
func (ds *DownloadStatus) GetETA() int64 {
	if ds.DownloadSpeed == 0 || ds.TotalLength == 0 {
		return -1
	}
	remaining := ds.TotalLength - ds.CompletedLength
	if remaining <= 0 {
		return 0
	}
	return remaining / ds.DownloadSpeed
}

// GetFilePath returns the full path to the downloaded file
func (ds *DownloadStatus) GetFilePath() string {
	if len(ds.Files) == 0 {
		return ""
	}
	// The first file's path is usually the main file
	filePath := ds.Files[0].Path
	if filePath == "" && len(ds.Files) > 0 {
		// Fallback to URI if path is empty
		if len(ds.Files[0].URIs) > 0 {
			filePath = ds.Files[0].URIs[0].URI
		}
	}
	return filePath
}

// GetFileDirectory returns the directory containing the downloaded file
func (ds *DownloadStatus) GetFileDirectory() string {
	if ds.Dir != "" {
		return ds.Dir
	}
	// Try to extract directory from file path
	filePath := ds.GetFilePath()
	if filePath != "" {
		return filepath.Dir(filePath)
	}
	return ""
}
