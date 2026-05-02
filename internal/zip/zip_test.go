package zip_test

import (
	"archive/zip"
	"os"
	"path/filepath"
	"testing"

	zippkg "github.com/convert/internal/zip"
)

func TestCreateFromDir(t *testing.T) {
	tmp := t.TempDir()

	// Create a small directory tree.
	files := map[string]string{
		"output/file1.md":        "# File 1",
		"output/file2.md":        "# File 2",
		"images/doc/image1.png":  "fakepng",
	}
	for rel, content := range files {
		full := filepath.Join(tmp, rel)
		_ = os.MkdirAll(filepath.Dir(full), 0755)
		_ = os.WriteFile(full, []byte(content), 0644)
	}

	destZip := filepath.Join(t.TempDir(), "archive.zip")
	if err := zippkg.CreateFromDir(tmp, destZip); err != nil {
		t.Fatalf("CreateFromDir: %v", err)
	}

	// Verify the zip contains all files.
	r, err := zip.OpenReader(destZip)
	if err != nil {
		t.Fatalf("open zip: %v", err)
	}
	defer r.Close()

	found := make(map[string]bool)
	for _, f := range r.File {
		found[f.Name] = true
	}

	for rel := range files {
		key := filepath.ToSlash(rel)
		if !found[key] {
			t.Errorf("zip missing entry: %s", key)
		}
	}
}
