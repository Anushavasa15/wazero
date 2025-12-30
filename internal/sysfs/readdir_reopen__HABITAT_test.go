package sysfs

import (
	"os"
	"path/filepath"
	"testing"
)

func TestReaddirReopensDirectory__HABITAT(t *testing.T) {
	dir := t.TempDir()

	// Create initial file
	if err := os.WriteFile(filepath.Join(dir, "first.txt"), []byte("one"), 0o644); err != nil {
		t.Fatalf("failed to create first file: %v", err)
	}

	// Open directory via sysfs
	fs := DirFS(dir)
	f, errno := fs.OpenFile(".", O_RDONLY, 0)
	if errno != 0 {
		t.Fatalf("failed to open directory: errno=%d", errno)
	}
	defer f.Close()

	// First read
	if _, errno := f.Readdir(-1); errno != 0 {
		t.Fatalf("first readdir failed: errno=%d", errno)
	}

	// Create file AFTER directory open
	if err := os.WriteFile(filepath.Join(dir, "second.txt"), []byte("two"), 0o644); err != nil {
		t.Fatalf("failed to create second file: %v", err)
	}

	// Second read must include new file
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
		t.Fatalf("expected readdir to include file created after open")
	}
}
