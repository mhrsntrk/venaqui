package aria2

import (
	"testing"
)

func TestDownloadStatus_GetProgress(t *testing.T) {
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
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ds := &DownloadStatus{
				TotalLength:     tt.totalLength,
				CompletedLength: tt.completedLength,
			}
			result := ds.GetProgress()
			if result != tt.expected {
				t.Errorf("DownloadStatus.GetProgress() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestDownloadStatus_IsComplete(t *testing.T) {
	tests := []struct {
		name     string
		status   string
		expected bool
	}{
		{
			name:     "complete status",
			status:   "complete",
			expected: true,
		},
		{
			name:     "active status",
			status:   "active",
			expected: false,
		},
		{
			name:     "error status",
			status:   "error",
			expected: false,
		},
		{
			name:     "waiting status",
			status:   "waiting",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ds := &DownloadStatus{
				Status: tt.status,
			}
			result := ds.IsComplete()
			if result != tt.expected {
				t.Errorf("DownloadStatus.IsComplete() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestDownloadStatus_IsError(t *testing.T) {
	tests := []struct {
		name     string
		status   string
		expected bool
	}{
		{
			name:     "error status",
			status:   "error",
			expected: true,
		},
		{
			name:     "complete status",
			status:   "complete",
			expected: false,
		},
		{
			name:     "active status",
			status:   "active",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ds := &DownloadStatus{
				Status: tt.status,
			}
			result := ds.IsError()
			if result != tt.expected {
				t.Errorf("DownloadStatus.IsError() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestDownloadStatus_IsActive(t *testing.T) {
	tests := []struct {
		name     string
		status   string
		expected bool
	}{
		{
			name:     "active status",
			status:   "active",
			expected: true,
		},
		{
			name:     "complete status",
			status:   "complete",
			expected: false,
		},
		{
			name:     "waiting status",
			status:   "waiting",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ds := &DownloadStatus{
				Status: tt.status,
			}
			result := ds.IsActive()
			if result != tt.expected {
				t.Errorf("DownloadStatus.IsActive() = %v, want %v", result, tt.expected)
			}
		})
	}
}
