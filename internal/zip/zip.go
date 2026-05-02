package zip

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// CreateFromDir creates a zip archive at destZip containing all files rooted
// at srcDir. The files are stored with paths relative to srcDir.
func CreateFromDir(srcDir, destZip string) error {
	// Ensure the destination parent directory exists.
	if err := os.MkdirAll(filepath.Dir(destZip), 0755); err != nil {
		return fmt.Errorf("create zip parent dir: %w", err)
	}

	zf, err := os.Create(destZip)
	if err != nil {
		return fmt.Errorf("create zip file: %w", err)
	}
	defer zf.Close()

	w := zip.NewWriter(zf)
	defer w.Close()

	err = filepath.WalkDir(srcDir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}

		rel, err := filepath.Rel(srcDir, path)
		if err != nil {
			return fmt.Errorf("relative path: %w", err)
		}
		// Always use forward slashes inside the archive.
		archivePath := filepath.ToSlash(rel)

		if err := addFileToZip(w, path, archivePath); err != nil {
			return fmt.Errorf("add %s to zip: %w", path, err)
		}
		return nil
	})
	return err
}

func addFileToZip(w *zip.Writer, fsPath, archivePath string) error {
	src, err := os.Open(fsPath)
	if err != nil {
		return err
	}
	defer src.Close()

	info, err := src.Stat()
	if err != nil {
		return err
	}

	header, err := zip.FileInfoHeader(info)
	if err != nil {
		return err
	}
	header.Name = archivePath
	header.Method = zip.Deflate

	dst, err := w.CreateHeader(header)
	if err != nil {
		return err
	}

	_, err = io.Copy(dst, src)
	return err
}
