package models

// DownloadState represents the state of a download
type DownloadState string

const (
	DownloadStateActive   DownloadState = "active"
	DownloadStateWaiting  DownloadState = "waiting"
	DownloadStatePaused   DownloadState = "paused"
	DownloadStateComplete DownloadState = "complete"
	DownloadStateError    DownloadState = "error"
	DownloadStateRemoved  DownloadState = "removed"
)
