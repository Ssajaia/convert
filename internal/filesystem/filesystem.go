package filesystem

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/convert/pkg/types"
)

// EnsureDir creates a directory and all parents if they don't exist.
func EnsureDir(path string) error {
	return os.MkdirAll(path, 0755)
}

// FileExists returns true if path exists and is a regular file.
func FileExists(path string) bool {
	info, err := os.Stat(path)
	return err == nil && info.Mode().IsRegular()
}

// DirExists returns true if path exists and is a directory.
func DirExists(path string) bool {
	info, err := os.Stat(path)
	return err == nil && info.IsDir()
}

// IsSupportedFile returns true if the file has a supported extension.
func IsSupportedFile(path string) bool {
	ext := strings.ToLower(filepath.Ext(path))
	return types.SupportedExtensions[ext]
}

// WalkSupportedFiles walks root (recursively if recursive=true) and returns
// all supported input files.
func WalkSupportedFiles(root string, recursive bool) ([]string, error) {
	var files []string

	err := filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			if path != root && !recursive {
				return filepath.SkipDir
			}
			return nil
		}
		if IsSupportedFile(path) {
			files = append(files, path)
		}
		return nil
	})
	return files, err
}

// CopyFile copies src to dst, creating dst's parent directories as needed.
func CopyFile(src, dst string) error {
	if err := EnsureDir(filepath.Dir(dst)); err != nil {
		return err
	}

	in, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("open src: %w", err)
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return fmt.Errorf("create dst: %w", err)
	}
	defer out.Close()

	if _, err := io.Copy(out, in); err != nil {
		return fmt.Errorf("copy: %w", err)
	}
	return out.Sync()
}

// MoveFile moves src to dst.
func MoveFile(src, dst string) error {
	if err := EnsureDir(filepath.Dir(dst)); err != nil {
		return err
	}
	// Try atomic rename first.
	if err := os.Rename(src, dst); err == nil {
		return nil
	}
	// Fallback: copy then delete.
	if err := CopyFile(src, dst); err != nil {
		return err
	}
	return os.Remove(src)
}

// SafeTempDir creates a temporary directory with the given prefix.
func SafeTempDir(prefix string) (string, error) {
	return os.MkdirTemp("", prefix)
}

// RemoveAll removes a directory tree, ignoring errors.
func RemoveAll(path string) {
	_ = os.RemoveAll(path)
}

// BaseName returns the filename without extension.
func BaseName(path string) string {
	base := filepath.Base(path)
	ext := filepath.Ext(base)
	return strings.TrimSuffix(base, ext)
}

// ReadFile reads a file and returns its contents.
func ReadFile(path string) ([]byte, error) {
	return os.ReadFile(path)
}

// WriteFile writes data to path, creating parent directories as needed.
func WriteFile(path string, data []byte) error {
	if err := EnsureDir(filepath.Dir(path)); err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}
