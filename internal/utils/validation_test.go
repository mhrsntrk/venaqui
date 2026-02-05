package utils

import (
	"os"
	"path/filepath"
	"testing"
)

func TestValidateURL(t *testing.T) {
	tests := []struct {
		name    string
		url     string
		wantErr bool
	}{
		{
			name:    "valid HTTP URL",
			url:     "http://example.com/file.zip",
			wantErr: false,
		},
		{
			name:    "valid HTTPS URL",
			url:     "https://example.com/file.zip",
			wantErr: false,
		},
		{
			name:    "empty URL",
			url:     "",
			wantErr: true,
		},
		{
			name:    "URL without scheme",
			url:     "example.com/file.zip",
			wantErr: true,
		},
		{
			name:    "URL without host",
			url:     "http:///file.zip",
			wantErr: true,
		},
		{
			name:    "invalid URL format",
			url:     "not a url",
			wantErr: true,
		},
		{
			name:    "URL with path",
			url:     "https://example.com/path/to/file.zip",
			wantErr: false,
		},
		{
			name:    "URL with query parameters",
			url:     "https://example.com/file.zip?token=abc123",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateURL(tt.url)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateURL() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestNormalizePath(t *testing.T) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		t.Fatalf("Failed to get home directory: %v", err)
	}

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "path with tilde",
			input:    "~/Downloads/test",
			expected: filepath.Join(homeDir, "Downloads", "test"),
		},
		{
			name:     "absolute path",
			input:    "/tmp/test",
			expected: "/tmp/test",
		},
		{
			name:     "path with dots",
			input:    "/tmp/../tmp/test",
			expected: "/tmp/test",
		},
		{
			name:     "path with multiple slashes",
			input:    "/tmp//test",
			expected: "/tmp/test",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := NormalizePath(tt.input)
			if result != tt.expected {
				t.Errorf("NormalizePath() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestValidatePath(t *testing.T) {
	// Create a temporary directory for testing
	tmpDir, err := os.MkdirTemp("", "venaqui-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create a writable subdirectory
	writableDir := filepath.Join(tmpDir, "writable")
	err = os.MkdirAll(writableDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create writable dir: %v", err)
	}

	// Create a non-writable directory (if possible)
	nonWritableDir := filepath.Join(tmpDir, "nonwritable")
	err = os.MkdirAll(nonWritableDir, 0555)
	if err != nil {
		t.Fatalf("Failed to create non-writable dir: %v", err)
	}

	tests := []struct {
		name    string
		path    string
		wantErr bool
	}{
		{
			name:    "empty path",
			path:    "",
			wantErr: true,
		},
		{
			name:    "relative path",
			path:    "relative/path",
			wantErr: true,
		},
		{
			name:    "existing writable directory",
			path:    writableDir,
			wantErr: false,
		},
		{
			name:    "non-existent path with writable parent",
			path:    filepath.Join(writableDir, "newdir"),
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidatePath(tt.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidatePath() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
