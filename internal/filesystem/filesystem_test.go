package filesystem_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/convert/internal/filesystem"
)

func TestBaseName(t *testing.T) {
	cases := []struct {
		input    string
		expected string
	}{
		{"file.docx", "file"},
		{"path/to/file.pdf", "file"},
		{"/abs/path/document.docx", "document"},
		{"no_ext", "no_ext"},
		{"dots.in.name.docx", "dots.in.name"},
	}
	for _, c := range cases {
		got := filesystem.BaseName(c.input)
		if got != c.expected {
			t.Errorf("BaseName(%q) = %q, want %q", c.input, got, c.expected)
		}
	}
}

func TestIsSupportedFile(t *testing.T) {
	cases := []struct {
		path      string
		supported bool
	}{
		{"file.docx", true},
		{"file.DOCX", true},
		{"file.pdf", true},
		{"file.PDF", true},
		{"file.txt", false},
		{"file.xlsx", false},
		{"file", false},
	}
	for _, c := range cases {
		got := filesystem.IsSupportedFile(c.path)
		if got != c.supported {
			t.Errorf("IsSupportedFile(%q) = %v, want %v", c.path, got, c.supported)
		}
	}
}

func TestEnsureDir(t *testing.T) {
	tmp := t.TempDir()
	nested := filepath.Join(tmp, "a", "b", "c")
	if err := filesystem.EnsureDir(nested); err != nil {
		t.Fatalf("EnsureDir: %v", err)
	}
	if !filesystem.DirExists(nested) {
		t.Fatal("expected directory to exist after EnsureDir")
	}
}

func TestCopyFile(t *testing.T) {
	tmp := t.TempDir()
	src := filepath.Join(tmp, "src.txt")
	dst := filepath.Join(tmp, "sub", "dst.txt")

	if err := os.WriteFile(src, []byte("hello"), 0644); err != nil {
		t.Fatal(err)
	}

	if err := filesystem.CopyFile(src, dst); err != nil {
		t.Fatalf("CopyFile: %v", err)
	}

	data, err := os.ReadFile(dst)
	if err != nil {
		t.Fatalf("read dst: %v", err)
	}
	if string(data) != "hello" {
		t.Errorf("got %q, want %q", string(data), "hello")
	}
	// Source must still exist.
	if !filesystem.FileExists(src) {
		t.Error("source file should still exist after copy")
	}
}

func TestMoveFile(t *testing.T) {
	tmp := t.TempDir()
	src := filepath.Join(tmp, "src.txt")
	dst := filepath.Join(tmp, "moved.txt")

	if err := os.WriteFile(src, []byte("move me"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := filesystem.MoveFile(src, dst); err != nil {
		t.Fatalf("MoveFile: %v", err)
	}
	if filesystem.FileExists(src) {
		t.Error("source should be gone after move")
	}
	if !filesystem.FileExists(dst) {
		t.Error("destination should exist after move")
	}
}

func TestWalkSupportedFiles(t *testing.T) {
	tmp := t.TempDir()

	create := func(path string) {
		full := filepath.Join(tmp, path)
		_ = os.MkdirAll(filepath.Dir(full), 0755)
		_ = os.WriteFile(full, nil, 0644)
	}

	create("a.docx")
	create("b.pdf")
	create("c.txt")
	create("sub/d.docx")
	create("sub/e.xlsx")

	// Non-recursive: should find a.docx and b.pdf only.
	files, err := filesystem.WalkSupportedFiles(tmp, false)
	if err != nil {
		t.Fatal(err)
	}
	if len(files) != 2 {
		t.Errorf("non-recursive: expected 2 files, got %d: %v", len(files), files)
	}

	// Recursive: should find a.docx, b.pdf, sub/d.docx.
	files, err = filesystem.WalkSupportedFiles(tmp, true)
	if err != nil {
		t.Fatal(err)
	}
	if len(files) != 3 {
		t.Errorf("recursive: expected 3 files, got %d: %v", len(files), files)
	}
}
