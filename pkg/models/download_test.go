package models

import (
	"testing"
	"time"
)

func TestDownload_Progress(t *testing.T) {
	tests := []struct {
		name           string
		totalLength    int64
		completedLength int64
		expected       float64
	}{
		{
			name:           "zero total length",
			totalLength:    0,
			completedLength: 100,
			expected:       0,
		},
		{
			name:           "half complete",
			totalLength:    1000,
			completedLength: 500,
			expected:       50.0,
		},
		{
			name:           "fully complete",
			totalLength:    1000,
			completedLength: 1000,
			expected:       100.0,
		},
		{
			name:           "not started",
			totalLength:    1000,
			completedLength: 0,
			expected:       0.0,
		},
		{
			name:           "partial progress",
			totalLength:    1000,
			completedLength: 250,
			expected:       25.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &Download{
				TotalLength:     tt.totalLength,
				CompletedLength: tt.completedLength,
			}
			result := d.Progress()
			if result != tt.expected {
				t.Errorf("Download.Progress() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestDownload_IsComplete(t *testing.T) {
	tests := []struct {
		name     string
		status   DownloadState
		expected bool
	}{
		{
			name:     "complete status",
			status:   DownloadStateComplete,
			expected: true,
		},
		{
			name:     "active status",
			status:   DownloadStateActive,
			expected: false,
		},
		{
			name:     "error status",
			status:   DownloadStateError,
			expected: false,
		},
		{
			name:     "waiting status",
			status:   DownloadStateWaiting,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &Download{
				Status: tt.status,
			}
			result := d.IsComplete()
			if result != tt.expected {
				t.Errorf("Download.IsComplete() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestDownload_IsError(t *testing.T) {
	tests := []struct {
		name     string
		status   DownloadState
		expected bool
	}{
		{
			name:     "error status",
			status:   DownloadStateError,
			expected: true,
		},
		{
			name:     "complete status",
			status:   DownloadStateComplete,
			expected: false,
		},
		{
			name:     "active status",
			status:   DownloadStateActive,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &Download{
				Status: tt.status,
			}
			result := d.IsError()
			if result != tt.expected {
				t.Errorf("Download.IsError() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestDownload_Fields(t *testing.T) {
	now := time.Now()
	d := &Download{
		GID:             "test-gid",
		Filename:        "test.zip",
		Status:          DownloadStateActive,
		TotalLength:     1000,
		CompletedLength: 500,
		DownloadSpeed:   100,
		CreatedAt:       now,
		UpdatedAt:       now,
	}

	if d.GID != "test-gid" {
		t.Errorf("Download.GID = %v, want test-gid", d.GID)
	}
	if d.Filename != "test.zip" {
		t.Errorf("Download.Filename = %v, want test.zip", d.Filename)
	}
	if d.Status != DownloadStateActive {
		t.Errorf("Download.Status = %v, want DownloadStateActive", d.Status)
	}
	if d.TotalLength != 1000 {
		t.Errorf("Download.TotalLength = %v, want 1000", d.TotalLength)
	}
	if d.CompletedLength != 500 {
		t.Errorf("Download.CompletedLength = %v, want 500", d.CompletedLength)
	}
	if d.DownloadSpeed != 100 {
		t.Errorf("Download.DownloadSpeed = %v, want 100", d.DownloadSpeed)
	}
}
