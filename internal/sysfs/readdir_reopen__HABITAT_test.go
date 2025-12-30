package sysfs

import (
	"os"
	"path/filepath"
	"testing"

	experimentalsys "github.com/tetratelabs/wazero/experimental/sys"
)

func TestReaddirReopensDirectory__HABITAT(t *testing.T) {
	dir := t.TempDir()

	// Create a file before opening the directory
	if err := os.WriteFile(filepath.Join(dir, "first.txt"), []byte("one"), 0o644); err != nil {
		t.Fatalf("failed to create first file: %v", err)
	}

	fs := DirFS(dir)

	// Open directory
	f, errno := fs.OpenFile(".", experimentalsys.O_RDONLY, 0)
	if errno != 0 {
		t.Fatalf("failed to open directory: errno=%d", errno)
	}

	// First read consumes existing entries
	if _, errno := f.Readdir(-1); errno != 0 {
		t.Fatalf("first readdir failed: errno=%d", errno)
	}

	// Create a file AFTER the directory has been read
	if err := os.WriteFile(filepath.Join(dir, "second.txt"), []byte("two"), 0o644); err != nil {
		t.Fatalf("failed to create second file: %v", err)
	}

	// IMPORTANT: close and reopen directory to observe new state
	if err := f.Close(); err != nil {
		t.Fatalf("failed to close directory: %v", err)
	}

	f, errno = fs.OpenFile(".", experimentalsys.O_RDONLY, 0)
	if errno != 0 {
		t.Fatalf("failed to reopen directory: errno=%d", errno)
	}
	defer f.Close()

	// Read again after reopen
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
		t.Fatalf("expected readdir to include file created after reopen")
	}
}
