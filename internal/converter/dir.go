package converter

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/convert/internal/filesystem"
	"github.com/convert/internal/logger"
	"github.com/convert/pkg/types"
)

// DirOptions configures directory-mode conversion.
type DirOptions struct {
	InputDir      string
	OutputDir     string
	ImagesRootDir string
	Format        string
	Recursive     bool
	Overwrite     bool
	Flatten       bool
}

// ConvertDir converts all supported files in InputDir.
// It returns a slice of results (including failures) so callers can report
// a summary without losing partial progress.
func ConvertDir(opts DirOptions) ([]types.ConvertResult, error) {
	files, err := filesystem.WalkSupportedFiles(opts.InputDir, opts.Recursive)
	if err != nil {
		return nil, err
	}

	if len(files) == 0 {
		logger.Warn("no supported files found in %s", opts.InputDir)
		return nil, nil
	}

	var results []types.ConvertResult

	for _, inputFile := range files {
		outputFile, err := resolveOutputPath(inputFile, opts.InputDir, opts.OutputDir, opts.Flatten)
		if err != nil {
			logger.Error("resolve output path for %s: %v", inputFile, err)
			results = append(results, types.ConvertResult{
				InputFile: inputFile,
				Error:     err,
			})
			continue
		}

		res, err := ConvertFile(Options{
			InputFile:     inputFile,
			OutputFile:    outputFile,
			ImagesRootDir: opts.ImagesRootDir,
			Format:        opts.Format,
			Overwrite:     opts.Overwrite,
		})
		if err != nil {
			logger.Error("convert %s: %v", inputFile, err)
			res.Error = err
		}
		results = append(results, res)
	}

	return results, nil
}

// resolveOutputPath computes the .md output path for an input file.
// If flatten is true, the directory hierarchy is collapsed to OutputDir.
func resolveOutputPath(inputFile, inputDir, outputDir string, flatten bool) (string, error) {
	baseName := filesystem.BaseName(inputFile)
	mdName := baseName + ".md"

	if flatten {
		return filepath.Join(outputDir, mdName), nil
	}

	// Preserve relative directory structure.
	rel, err := filepath.Rel(inputDir, inputFile)
	if err != nil {
		return "", err
	}
	relDir := filepath.Dir(rel)
	if relDir == "." {
		relDir = ""
	}

	return filepath.Join(outputDir, relDir, mdName), nil
}

// SummaryLine returns a human-readable summary for a set of results.
func SummaryLine(results []types.ConvertResult) string {
	ok, fail := 0, 0
	for _, r := range results {
		if r.Error != nil {
			fail++
		} else {
			ok++
		}
	}
	sb := strings.Builder{}
	sb.WriteString(fmt.Sprintf("conversion complete: %d succeeded, %d failed", ok, fail))
	return sb.String()
}
