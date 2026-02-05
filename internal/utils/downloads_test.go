package utils

import (
	"os"
	"path/filepath"
	"testing"
)

func TestGetDownloadFolder(t *testing.T) {
	folder, err := GetDownloadFolder()
	if err != nil {
		t.Fatalf("GetDownloadFolder() error = %v", err)
	}

	if folder == "" {
		t.Error("GetDownloadFolder() returned empty string")
	}

	// Should contain "Downloads"
	if filepath.Base(folder) != "Downloads" {
		t.Errorf("GetDownloadFolder() = %v, expected path ending with Downloads", folder)
	}

	// Should be absolute path
	if !filepath.IsAbs(folder) {
		t.Errorf("GetDownloadFolder() = %v, expected absolute path", folder)
	}
}

func TestEnsureDirExists(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "venaqui-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	tests := []struct {
		name    string
		dir     string
		wantErr bool
	}{
		{
			name:    "create new directory",
			dir:     filepath.Join(tmpDir, "newdir"),
			wantErr: false,
		},
		{
			name:    "create nested directories",
			dir:     filepath.Join(tmpDir, "nested", "deep", "dir"),
			wantErr: false,
		},
		{
			name:    "existing directory",
			dir:     tmpDir,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := EnsureDirExists(tt.dir)
			if (err != nil) != tt.wantErr {
				t.Errorf("EnsureDirExists() error = %v, wantErr %v", err, tt.wantErr)
			}

			// Verify directory was created
			if !tt.wantErr {
				info, err := os.Stat(tt.dir)
				if err != nil {
					t.Errorf("EnsureDirExists() directory not created: %v", err)
				}
				if !info.IsDir() {
					t.Errorf("EnsureDirExists() created path is not a directory")
				}
			}
		})
	}
}
