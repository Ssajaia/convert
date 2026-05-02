package images

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/convert/internal/filesystem"
	"github.com/convert/internal/logger"
)

// NewNormalizer creates a Normalizer.
func NewNormalizer(sourceDir, targetDir string) *Normalizer {
	return &Normalizer{SourceDir: sourceDir, TargetDir: targetDir}
}

// Normalizer moves extracted images into a canonical location and returns
// a mapping from old relative path → new relative path.
type Normalizer struct {
	// SourceDir is the directory pandoc wrote images into (--extract-media target).
	SourceDir string
	// TargetDir is the canonical images/<filename>/ directory.
	TargetDir string
}

// Normalize discovers all images under SourceDir, moves them to TargetDir
// with sequentially numbered names, and returns a path-replacement map.
// Keys and values are slash-separated relative paths suitable for Markdown.
func (n *Normalizer) Normalize() (map[string]string, error) {
	replacements := make(map[string]string)

	if !filesystem.DirExists(n.SourceDir) {
		// No images extracted – nothing to do.
		return replacements, nil
	}

	files, err := collectImageFiles(n.SourceDir)
	if err != nil {
		return nil, fmt.Errorf("collect images: %w", err)
	}

	if len(files) == 0 {
		return replacements, nil
	}

	if err := filesystem.EnsureDir(n.TargetDir); err != nil {
		return nil, fmt.Errorf("create target dir: %w", err)
	}

	for i, src := range files {
		ext := strings.ToLower(filepath.Ext(src))
		newName := fmt.Sprintf("image%d%s", i+1, ext)
		dst := filepath.Join(n.TargetDir, newName)

		// Handle collisions.
		dst = resolveCollision(dst)

		if err := filesystem.MoveFile(src, dst); err != nil {
			logger.Warn("failed to move image %s → %s: %v", src, dst, err)
			continue
		}

		// Build old key: path relative to SourceDir's parent (pandoc writes
		// paths like "<extractMediaArg>/media/imageN.ext" in the Markdown).
		// We store several key variants to maximise match coverage.
		relSrc, _ := filepath.Rel(filepath.Dir(n.SourceDir), src)
		relSrc = filepath.ToSlash(relSrc)

		relDst, _ := filepath.Rel(filepath.Dir(n.TargetDir), dst)
		relDst = filepath.ToSlash(relDst)

		replacements[relSrc] = relDst

		// Also map by just the filename as a fallback key.
		replacements[filepath.Base(src)] = relDst

		logger.Debug("image: %s → %s", relSrc, relDst)
	}

	// Clean up the source directory if it's now empty.
	_ = removeEmptyDirs(n.SourceDir)

	return replacements, nil
}

// collectImageFiles returns all regular files under root, sorted for determinism.
func collectImageFiles(root string) ([]string, error) {
	var files []string
	err := filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() {
			files = append(files, path)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	sort.Strings(files)
	return files, nil
}

// resolveCollision appends a counter suffix if dst already exists.
func resolveCollision(dst string) string {
	if !filesystem.FileExists(dst) {
		return dst
	}
	ext := filepath.Ext(dst)
	base := strings.TrimSuffix(dst, ext)
	for i := 2; ; i++ {
		candidate := fmt.Sprintf("%s_%d%s", base, i, ext)
		if !filesystem.FileExists(candidate) {
			return candidate
		}
	}
}

// removeEmptyDirs removes empty directories bottom-up.
func removeEmptyDirs(root string) error {
	return filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
		if err != nil || !d.IsDir() || path == root {
			return nil
		}
		entries, _ := os.ReadDir(path)
		if len(entries) == 0 {
			_ = os.Remove(path)
		}
		return nil
	})
}

// UpdateMarkdownPaths replaces image paths in markdown content using the
// replacements map produced by Normalizer.Normalize().
func UpdateMarkdownPaths(content string, replacements map[string]string) string {
	for old, newPath := range replacements {
		// Markdown image syntax: ![alt](path) or ![alt](path "title")
		content = replaceImagePath(content, old, newPath)
	}
	return content
}

func replaceImagePath(content, old, newPath string) string {
	// Normalise separators for matching.
	oldSlash := filepath.ToSlash(old)

	// We do a simple string replacement on the path segment inside Markdown
	// image references. This handles the common pandoc output patterns.
	result := strings.ReplaceAll(content, "]("+oldSlash, "]("+newPath)
	result = strings.ReplaceAll(result, "]("+old, "]("+newPath)

	// Also handle encoded spaces (%20) in filenames.
	encoded := strings.ReplaceAll(oldSlash, " ", "%20")
	result = strings.ReplaceAll(result, "]("+encoded, "]("+newPath)

	return result
}
