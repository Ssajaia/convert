package images_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/convert/internal/images"
)

func TestUpdateMarkdownPaths(t *testing.T) {
	cases := []struct {
		name         string
		content      string
		replacements map[string]string
		expected     string
	}{
		{
			name:         "single image",
			content:      "![alt](media/image1.png)",
			replacements: map[string]string{"media/image1.png": "images/doc/image1.png"},
			expected:     "![alt](images/doc/image1.png)",
		},
		{
			name:         "multiple images",
			content:      "![a](media/img1.png)\n![b](media/img2.jpg)",
			replacements: map[string]string{"media/img1.png": "images/doc/image1.png", "media/img2.jpg": "images/doc/image2.jpg"},
			expected:     "![a](images/doc/image1.png)\n![b](images/doc/image2.jpg)",
		},
		{
			name:         "no match",
			content:      "![alt](other/image.png)",
			replacements: map[string]string{"media/image1.png": "images/doc/image1.png"},
			expected:     "![alt](other/image.png)",
		},
		{
			name:         "empty replacements",
			content:      "![alt](media/image1.png)",
			replacements: map[string]string{},
			expected:     "![alt](media/image1.png)",
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := images.UpdateMarkdownPaths(c.content, c.replacements)
			if got != c.expected {
				t.Errorf("got:\n%s\nwant:\n%s", got, c.expected)
			}
		})
	}
}

func TestNormalizer_NoImages(t *testing.T) {
	tmp := t.TempDir()
	norm := images.NewNormalizer(filepath.Join(tmp, "nonexistent"), filepath.Join(tmp, "target"))
	replacements, err := norm.Normalize()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(replacements) != 0 {
		t.Errorf("expected empty replacements, got %v", replacements)
	}
}

func TestNormalizer_MovesImages(t *testing.T) {
	tmp := t.TempDir()
	srcDir := filepath.Join(tmp, "media")
	dstDir := filepath.Join(tmp, "images", "doc")

	if err := os.MkdirAll(srcDir, 0755); err != nil {
		t.Fatal(err)
	}

	// Create fake image files.
	for _, name := range []string{"img1.png", "img2.jpg"} {
		if err := os.WriteFile(filepath.Join(srcDir, name), []byte("fake"), 0644); err != nil {
			t.Fatal(err)
		}
	}

	norm := images.NewNormalizer(srcDir, dstDir)

	replacements, err := norm.Normalize()
	if err != nil {
		t.Fatalf("Normalize: %v", err)
	}

	if len(replacements) < 2 {
		t.Fatalf("expected at least 2 replacements, got %d: %v", len(replacements), replacements)
	}

	// Verify destination files exist.
	entries, _ := os.ReadDir(dstDir)
	if len(entries) != 2 {
		t.Errorf("expected 2 files in target dir, got %d", len(entries))
	}

	// Verify sequential naming.
	names := make(map[string]bool)
	for _, e := range entries {
		names[e.Name()] = true
	}
	if !names["image1.png"] && !names["image1.jpg"] {
		t.Error("expected image1.ext in target dir")
	}
}
