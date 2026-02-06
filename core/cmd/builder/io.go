package main

import (
	"io"
	"os"
	"path/filepath"
)

// CopyDir recursively copies a directory tree from src to dst.
func CopyDir(src string, dst string) error {
	return filepath.WalkDir(src, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// Calculate the destination path
		relPath, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}
		targetPath := filepath.Join(dst, relPath)

		if d.IsDir() {
			// Create directory with same permissions
			info, err := d.Info()
			if err != nil {
				return err
			}
			return os.MkdirAll(targetPath, info.Mode())
		}

		// If it's a file, copy it
		return CopyFile(path, targetPath)
	})
}

// CopyFile copies a single file from src to dst, overwriting if it exists.
func CopyFile(src, dst string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	// Get source file mode for the new file
	info, err := sourceFile.Stat()
	if err != nil {
		return err
	}

	// Create/Overwrite destination file
	destFile, err := os.OpenFile(dst, os.O_RDWR|os.O_CREATE|os.O_TRUNC, info.Mode())
	if err != nil {
		return err
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, sourceFile)
	return err
}
