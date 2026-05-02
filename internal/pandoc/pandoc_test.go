package pandoc_test

import (
	"testing"

	"github.com/convert/internal/pandoc"
)

func TestCheckAvailable(t *testing.T) {
	// This test will pass if pandoc is installed, skip otherwise.
	if err := pandoc.CheckAvailable(); err != nil {
		t.Skipf("pandoc not available: %v", err)
	}
}

func TestVersion(t *testing.T) {
	if err := pandoc.CheckAvailable(); err != nil {
		t.Skip("pandoc not available")
	}
	v, err := pandoc.Version()
	if err != nil {
		t.Fatalf("Version: %v", err)
	}
	if v == "" {
		t.Error("expected non-empty version string")
	}
	t.Logf("pandoc version: %s", v)
}

func TestRun_InvalidInput(t *testing.T) {
	if err := pandoc.CheckAvailable(); err != nil {
		t.Skip("pandoc not available")
	}

	err := pandoc.Run(pandoc.ConvertOptions{
		InputFile:  "/nonexistent/file.docx",
		OutputFile: "/tmp/out.md",
	})
	if err == nil {
		t.Error("expected error for nonexistent input file")
	}
}
