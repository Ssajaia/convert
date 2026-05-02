package cli

import (
	"flag"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/convert/internal/converter"
	"github.com/convert/internal/filesystem"
	"github.com/convert/internal/logger"
	"github.com/convert/internal/pandoc"
	zippkg "github.com/convert/internal/zip"
)

// runDirCmd handles: convert dir <input_dir> <output_dir_or_zip> <images_dir> [--zip]
func runDirCmd(args []string) error {
	fs := flag.NewFlagSet("dir", flag.ContinueOnError)
	zipOutput := fs.Bool("zip", false, "package output as a zip archive")

	if err := fs.Parse(args); err != nil {
		return err
	}
	rest := fs.Args()
	if len(rest) != 3 {
		return fmt.Errorf("dir mode requires: <input_dir> <output_dir_or_zip> <images_dir>\ngot: %v", rest)
	}

	inputDir := rest[0]
	outputTarget := rest[1]
	imagesDir := rest[2]

	if err := logger.Init(gf.logFile, logger.INFO); err != nil {
		return err
	}
	if err := pandoc.CheckAvailable(); err != nil {
		return err
	}
	if !filesystem.DirExists(inputDir) {
		return fmt.Errorf("input directory not found: %s", inputDir)
	}

	var mdOutputDir string
	var tempDir string

	if *zipOutput {
		var err error
		tempDir, err = filesystem.SafeTempDir("convert-")
		if err != nil {
			return fmt.Errorf("create temp dir: %w", err)
		}
		defer filesystem.RemoveAll(tempDir)

		mdOutputDir = filepath.Join(tempDir, "output")
		imagesDir = filepath.Join(tempDir, "images")
	} else {
		mdOutputDir = outputTarget
	}

	results, err := converter.ConvertDir(converter.DirOptions{
		InputDir:      inputDir,
		OutputDir:     mdOutputDir,
		ImagesRootDir: imagesDir,
		Format:        gf.format,
		Recursive:     gf.recursive,
		Overwrite:     gf.overwrite,
		Flatten:       gf.flatten,
	})
	if err != nil {
		return err
	}

	for _, r := range results {
		if r.Error != nil {
			fmt.Printf("FAIL %s: %v\n", r.InputFile, r.Error)
		} else {
			fmt.Printf("OK   %s → %s (%d image(s))\n", r.InputFile, r.OutputFile, r.ImageCount)
		}
	}
	fmt.Println(converter.SummaryLine(results))

	if *zipOutput {
		destZip := outputTarget
		if !strings.HasSuffix(strings.ToLower(destZip), ".zip") {
			destZip += ".zip"
		}
		logger.Info("creating zip: %s", destZip)
		if err := zippkg.CreateFromDir(tempDir, destZip); err != nil {
			return fmt.Errorf("create zip: %w", err)
		}
		fmt.Printf("zip created: %s\n", destZip)
	}

	return nil
}
