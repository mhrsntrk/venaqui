package models

import "time"

// Download represents a download task
type Download struct {
	GID             string
	Filename        string
	Status          DownloadState
	TotalLength     int64
	CompletedLength int64
	DownloadSpeed   int64
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

// Progress returns the download progress as a percentage (0-100)
func (d *Download) Progress() float64 {
	if d.TotalLength == 0 {
		return 0
	}
	return float64(d.CompletedLength) / float64(d.TotalLength) * 100
}

// IsComplete returns true if the download is complete
func (d *Download) IsComplete() bool {
	return d.Status == DownloadStateComplete
}

// IsError returns true if the download has an error
func (d *Download) IsError() bool {
	return d.Status == DownloadStateError
}
