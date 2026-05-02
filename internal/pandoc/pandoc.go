package pandoc

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
)

// CheckAvailable returns an error if pandoc is not found in PATH.
func CheckAvailable() error {
	_, err := exec.LookPath("pandoc")
	if err != nil {
		return fmt.Errorf("pandoc not found in PATH: install from https://pandoc.org/installing.html")
	}
	return nil
}

// Version returns the pandoc version string.
func Version() (string, error) {
	out, err := exec.Command("pandoc", "--version").Output()
	if err != nil {
		return "", err
	}
	lines := strings.SplitN(string(out), "\n", 2)
	if len(lines) > 0 {
		return strings.TrimSpace(lines[0]), nil
	}
	return "", nil
}

// ConvertOptions configures a single pandoc invocation.
type ConvertOptions struct {
	InputFile    string
	OutputFile   string
	ExtractMedia string // directory for --extract-media
	Format       string // output format, e.g. "markdown" or "gfm"
	ExtraArgs    []string
}

// Run executes pandoc with the given options.
// It returns the combined stdout+stderr on failure for diagnostic purposes.
func Run(opts ConvertOptions) error {
	if err := CheckAvailable(); err != nil {
		return err
	}

	format := opts.Format
	if format == "" {
		format = "markdown"
	}

	args := []string{
		opts.InputFile,
		"-o", opts.OutputFile,
		"-t", format,
		"--wrap=none",
	}

	if opts.ExtractMedia != "" {
		args = append(args, fmt.Sprintf("--extract-media=%s", opts.ExtractMedia))
	}

	args = append(args, opts.ExtraArgs...)

	cmd := exec.Command("pandoc", args...)

	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("pandoc failed: %w\nstderr: %s", err, stderr.String())
	}
	return nil
}
