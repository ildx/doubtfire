package utils

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// resolveFileNameConflict appends a running number to the file name if a file with the same name already exists
func ResolveFileNameConflict(destPath string) string {
	base := destPath
	ext := filepath.Ext(destPath)
	name := strings.TrimSuffix(base, ext)
	counter := 1

	for {
		if _, err := os.Stat(destPath); os.IsNotExist(err) {
			break
		}
		destPath = fmt.Sprintf("%s(%d)%s", name, counter, ext)
		fmt.Println("Conflict detected, new destination path:", destPath) // Debugging output
		counter++
	}

	return destPath
}

// CopyDir recursively copies a directory from src to dst.
func CopyDir(src, dst string) error {
	// Read the source directory
	entries, err := os.ReadDir(src)
	if err != nil {
		return err
	}

	// Create the destination directory
	err = os.MkdirAll(dst, os.ModePerm)
	if err != nil {
		return err
	}

	// Loop through each entry in the source directory
	for _, entry := range entries {
		srcPath := filepath.Join(src, entry.Name())
		dstPath := filepath.Join(dst, entry.Name())

		if entry.IsDir() {
			// If the entry is a directory, recursively copy it
			err = CopyDir(srcPath, dstPath)
			if err != nil {
				return err
			}
		} else {
			// If the entry is a file, copy the file
			err = CopyFile(srcPath, dstPath)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// CopyFile copies a single file from src to dst.
func CopyFile(src, dst string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	destFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, sourceFile)
	return err
}
