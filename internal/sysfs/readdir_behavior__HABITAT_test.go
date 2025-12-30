package sysfs

import (
	"os"
	"path/filepath"
	"testing"
)

func TestReaddirReopensDirectory__HABITAT(t *testing.T) {
	dir := t.TempDir()

	firstFile := filepath.Join(dir, "first.txt")
	if err := os.WriteFile(firstFile, []byte("one"), 0o644); err != nil {
		t.Fatalf("failed to create first file: %v", err)
	}

	f, err := os.Open(dir)
	if err != nil {
		t.Fatalf("failed to open directory: %v", err)
	}
	defer f.Close()

	_, err = f.Readdirnames(-1)
	if err != nil {
		t.Fatalf("first readdir failed: %v", err)
	}

	secondFile := filepath.Join(dir, "second.txt")
	if err := os.WriteFile(secondFile, []byte("two"), 0o644); err != nil {
		t.Fatalf("failed to create second file: %v", err)
	}

	names, err := f.Readdirnames(-1)
	if err != nil {
		t.Fatalf("second readdir failed: %v", err)
	}

	found := false
	for _, n := range names {
		if n == "second.txt" {
			found = true
		}
	}

	if !found {
		t.Fatalf("expected readdir to include file created after open")
	}
}
