package utils

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/log"
	"github.com/ildx/doubtfire/internal/errors"
)

// ResolveFileNameConflict appends a running number to the file name if a file with the same name already exists
func ResolveFileNameConflict(destPath string) string {
	dir := filepath.Dir(destPath)
	base := filepath.Base(destPath)
	ext := filepath.Ext(destPath)
	name := strings.TrimSuffix(base, ext)
	counter := 1

	for {
		if _, err := os.Stat(destPath); os.IsNotExist(err) {
			break
		}
		if counter == 1 {
			destPath = filepath.Join(dir, fmt.Sprintf("%s copy%s", name, ext))
		} else {
			destPath = filepath.Join(dir, fmt.Sprintf("%s copy %d%s", name, counter, ext))
		}
		log.Info("Resolving conflict, new path:", "path", destPath)
		counter++
	}

	return destPath
}

func CopyDir(src, dst string) error {
	// Read the source directory
	entries, err := os.ReadDir(src)
	if err != nil {
		log.Error(errors.ErrReadDir, err)
		return err
	}

	// Create the destination directory
	err = os.MkdirAll(dst, os.ModePerm)
	if err != nil {
		log.Error(errors.ErrCreateDir, err)
		return err
	}

	// Initialize the progress model
	totalEntries := len(entries)
	p := tea.NewProgram(initialModel(totalEntries))

	// Create a channel to communicate progress updates
	progressChannel := make(chan struct{})

	// Start the progress model in a Goroutine
	go func() {
		if _, err := p.Run(); err != nil {
			log.Error("Error running progress program:", err)
		}
	}()

	// Create a Goroutine to handle copying
	go func() {
		defer close(progressChannel)
		for _, entry := range entries {
			srcPath := filepath.Join(src, entry.Name())
			dstPath := filepath.Join(dst, entry.Name())

			if entry.IsDir() {
				dstPath = ResolveFileNameConflict(dstPath)
				if err := CopyDir(srcPath, dstPath); err != nil {
					log.Warn("Skipping directory: error", "path", srcPath, "error", err)
					continue
				}
			} else {
				dstPath = ResolveFileNameConflict(dstPath)
				if err := CopyFile(srcPath, dstPath); err != nil {
					log.Warn("Skipping file: error", "path", srcPath, "error", err)
					continue
				}
			}

			// Send a progress update after each file
			progressChannel <- struct{}{}
		}
	}()

	// Handle progress updates
	go func() {
		for range progressChannel {
			p.Send(struct{}{}) // Send progress update to the progress model
		}
	}()

	// Wait for all file operations to complete
	<-progressChannel // Wait for the copying to complete
	log.Info("File copying completed successfully!")

	return nil
}

// CopyFile copies a single file from src to dst.
func CopyFile(src, dst string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		log.Error(errors.ErrCopyFile, err)
		return err
	}
	defer sourceFile.Close()

	destFile, err := os.Create(dst)
	if err != nil {
		log.Error(errors.ErrCopyFile, err)
		return err
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, sourceFile)
	if err != nil {
		log.Error(errors.ErrCopyFile, err)
		return err
	}

	log.Info("File copied successfully", "path", dst)

	return nil
}

// CreateDirectory creates the destination directory if it does not exist
func CreateDirectory(dir string) error {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if err := os.MkdirAll(dir, os.ModePerm); err != nil {
			log.Error(errors.ErrCreateDir, err)
			return err
		}
	}
	return nil
}
