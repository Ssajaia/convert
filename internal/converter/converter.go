package converter

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/convert/internal/filesystem"
	"github.com/convert/internal/images"
	"github.com/convert/internal/logger"
	"github.com/convert/internal/pandoc"
	"github.com/convert/pkg/types"
)

// Options configures a single file conversion.
type Options struct {
	InputFile  string
	OutputFile string
	// ImagesRootDir is the top-level images directory (e.g. "images/").
	// Individual files get a sub-directory: images/<basename>/
	ImagesRootDir string
	Format        string
	Overwrite     bool
}

// ConvertFile converts a single .docx or .pdf file to Markdown.
// It runs pandoc, normalises extracted images, and rewrites image paths.
func ConvertFile(opts Options) (types.ConvertResult, error) {
	result := types.ConvertResult{
		InputFile:  opts.InputFile,
		OutputFile: opts.OutputFile,
	}

	// --- Validation ---
	if !filesystem.IsSupportedFile(opts.InputFile) {
		return result, fmt.Errorf("unsupported file type: %s", filepath.Ext(opts.InputFile))
	}
	if !filesystem.FileExists(opts.InputFile) {
		return result, fmt.Errorf("input file not found: %s", opts.InputFile)
	}
	if !opts.Overwrite && filesystem.FileExists(opts.OutputFile) {
		return result, fmt.Errorf("output file already exists (use --overwrite): %s", opts.OutputFile)
	}

	// --- Prepare directories ---
	baseName := filesystem.BaseName(opts.InputFile)
	imgTargetDir := filepath.Join(opts.ImagesRootDir, baseName)

	// pandoc writes extracted media into a temp working subdir so we can
	// distinguish new files from pre-existing ones.
	imgExtractDir := filepath.Join(opts.ImagesRootDir, "_extract_"+baseName)

	if err := filesystem.EnsureDir(filepath.Dir(opts.OutputFile)); err != nil {
		return result, fmt.Errorf("create output dir: %w", err)
	}

	// --- Run pandoc ---
	logger.Info("converting: %s → %s", opts.InputFile, opts.OutputFile)

	pandocOpts := pandoc.ConvertOptions{
		InputFile:    opts.InputFile,
		OutputFile:   opts.OutputFile,
		ExtractMedia: imgExtractDir,
		Format:       opts.Format,
	}
	if err := pandoc.Run(pandocOpts); err != nil {
		return result, fmt.Errorf("pandoc: %w", err)
	}

	// --- Normalise images ---
	norm := &images.Normalizer{
		SourceDir: imgExtractDir,
		TargetDir: imgTargetDir,
	}
	replacements, err := norm.Normalize()
	if err != nil {
		// Non-fatal: log and continue with un-normalised paths.
		logger.Warn("image normalization: %v", err)
	} else {
		result.ImageCount = len(replacements)
		result.ImagesDir = imgTargetDir
	}

	// --- Rewrite Markdown image paths ---
	if len(replacements) > 0 {
		raw, err := filesystem.ReadFile(opts.OutputFile)
		if err != nil {
			return result, fmt.Errorf("read output md: %w", err)
		}

		updated := images.UpdateMarkdownPaths(string(raw), replacements)

		// Also clean up any residual extract dir prefix that pandoc may have left.
		updated = cleanExtractPrefix(updated, imgExtractDir, imgTargetDir)

		if err := filesystem.WriteFile(opts.OutputFile, []byte(updated)); err != nil {
			return result, fmt.Errorf("write updated md: %w", err)
		}
		logger.Info("updated %d image path(s)", len(replacements))
	}

	// Clean up the extract staging directory (should be empty after move).
	filesystem.RemoveAll(imgExtractDir)

	logger.Info("done: %s (%d images)", opts.OutputFile, result.ImageCount)
	return result, nil
}

// cleanExtractPrefix replaces any remaining references to the staging extract
// dir with the final images dir in the Markdown content.
func cleanExtractPrefix(content, extractDir, targetDir string) string {
	extractSlash := filepath.ToSlash(extractDir)
	targetSlash := filepath.ToSlash(targetDir)

	// Strip leading "./" if present for consistent matching.
	extractSlash = strings.TrimPrefix(extractSlash, "./")
	targetSlash = strings.TrimPrefix(targetSlash, "./")

	return strings.ReplaceAll(content, extractSlash, targetSlash)
}
