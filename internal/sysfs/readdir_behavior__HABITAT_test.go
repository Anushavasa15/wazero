package sysfs

import (
	"context"
	"os"
	"path/filepath"
	"testing"
)

func TestReaddirReopensDirectory__HABITAT(t *testing.T) {
	ctx := context.Background()
	dir := t.TempDir()

	// Create initial file
	if err := os.WriteFile(filepath.Join(dir, "first.txt"), []byte("one"), 0o644); err != nil {
		t.Fatalf("failed to create first file: %v", err)
	}

	// Use sysfs filesystem (this is critical)
	fs := NewDirFS(dir)

	f, err := fs.OpenFile(ctx, ".", os.O_RDONLY, 0)
	if err != nil {
		t.Fatalf("failed to open directory via sysfs: %v", err)
	}
	defer f.Close(ctx)

	// First readdir
	if _, err := f.Readdirnames(-1); err != nil {
		t.Fatalf("first readdir failed: %v", err)
	}

	// Create file AFTER directory open
	if err := os.WriteFile(filepath.Join(dir, "second.txt"), []byte("two"), 0o644); err != nil {
		t.Fatalf("failed to create second file: %v", err)
	}

	// Second readdir — this is what PR #2355 fixes
	names, err := f.Readdirnames(-1)
	if err != nil {
		t.Fatalf("second readdir failed: %v", err)
	}

	found := false
	for _, n := range names {
		if n == "second.txt" {
			found = true
			break
		}
	}

	if !found {
		t.Fatalf("expected readdir to include file created after open")
	}
}
