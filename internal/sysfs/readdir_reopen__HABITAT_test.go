package sysfs

import (
	"io"
	"os"
	"path/filepath"
	"testing"

	experimentalsys "github.com/tetratelabs/wazero/experimental/sys"
)

func TestReaddirReopensDirectory__HABITAT(t *testing.T) {
	dir := t.TempDir()

	// Create file before opening directory
	if err := os.WriteFile(filepath.Join(dir, "first.txt"), []byte("one"), 0o644); err != nil {
		t.Fatalf("failed to create first file: %v", err)
	}

	fs := DirFS(dir)

	// Open directory ONCE
	f, errno := fs.OpenFile(".", experimentalsys.O_RDONLY, 0)
	if errno != 0 {
		t.Fatalf("failed to open directory: errno=%d", errno)
	}
	defer f.Close()

	// First read
	if _, errno := f.Readdir(-1); errno != 0 {
		t.Fatalf("first readdir failed: errno=%d", errno)
	}

	// Create file AFTER first read
	if err := os.WriteFile(filepath.Join(dir, "second.txt"), []byte("two"), 0o644); err != nil {
		t.Fatalf("failed to create second file: %v", err)
	}

	// Reset position (golden patch reopens internally)
	if _, err := f.Seek(0, io.SeekStart); err != nil {
		t.Fatalf("failed to seek directory: %v", err)
	}

	// Second read on SAME handle
	dirents, errno := f.Readdir(-1)
	if errno != 0 {
		t.Fatalf("second readdir failed: errno=%d", errno)
	}

	found := false
	for _, d := range dirents {
		if d.Name == "second.txt" {
			found = true
			break
		}
	}

	if !found {
		t.Fatalf("expected readdir to include file created after first read on same handle")
	}
}
