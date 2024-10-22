package ui

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/ildx/doubtfire/internal/utils"
)

func CleanUp() (string, error) {
	var logBuffer strings.Builder
	fmt.Fprintln(&logBuffer, "Starting cleanup process")

	config, err := utils.LoadConfig()
	if err != nil {
		return logBuffer.String(), fmt.Errorf("failed to load config: %w", err)
	}
	fmt.Fprintf(&logBuffer, "Loaded config, destination directory: %s\n", config.DestinationDir)

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return logBuffer.String(), fmt.Errorf("failed to get home directory: %w", err)
	}

	desktopPath, err := utils.GetDesktopPath()
	if err != nil {
		return logBuffer.String(), fmt.Errorf("failed to get desktop path: %w", err)
	}
	fmt.Fprintf(&logBuffer, "Desktop path: %s\n", desktopPath)

	now := time.Now()
	destPath := filepath.Join(homeDir, config.DestinationDir, fmt.Sprintf("%d", now.Year()), fmt.Sprintf("%02d", now.Month()))
	fmt.Fprintf(&logBuffer, "Destination path: %s\n", destPath)

	// Create the destination directory
	err = os.MkdirAll(destPath, os.ModePerm)
	if err != nil {
		return logBuffer.String(), fmt.Errorf("failed to create destination directory: %w", err)
	}

	// First pass: move files and create directories
	err = filepath.Walk(desktopPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Fprintf(&logBuffer, "Error accessing path %s: %v\n", path, err)
			return err
		}

		// Skip the desktop directory itself
		if path == desktopPath {
			return nil
		}

		relPath, err := filepath.Rel(desktopPath, path)
		if err != nil {
			fmt.Fprintf(&logBuffer, "Error getting relative path for %s: %v\n", path, err)
			return fmt.Errorf("failed to get relative path: %w", err)
		}

		newPath := filepath.Join(destPath, relPath)

		if info.IsDir() {
			uniquePath := utils.GetUniquePath(newPath)
			fmt.Fprintf(&logBuffer, "Creating directory: %s\n", uniquePath)
			return os.MkdirAll(uniquePath, os.ModePerm)
		}

		fmt.Fprintf(&logBuffer, "Moving file: %s to %s\n", path, newPath)
		return utils.MoveFile(path, newPath)
	})

	if err != nil {
		fmt.Fprintf(&logBuffer, "Error during first pass: %v\n", err)
		return logBuffer.String(), err
	}

	// Second pass: remove empty directories
	err = filepath.Walk(desktopPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Fprintf(&logBuffer, "Error accessing path %s: %v\n", path, err)
			return err
		}

		// Skip the desktop directory itself
		if path == desktopPath {
			return nil
		}

		if info.IsDir() {
			empty, err := utils.IsDirEmpty(path)
			if err != nil {
				fmt.Fprintf(&logBuffer, "Error checking if directory is empty %s: %v\n", path, err)
				return fmt.Errorf("failed to check if directory is empty: %w", err)
			}
			if empty {
				fmt.Fprintf(&logBuffer, "Removing empty directory: %s\n", path)
				if err := os.Remove(path); err != nil {
					fmt.Fprintf(&logBuffer, "Error removing empty directory %s: %v\n", path, err)
					return fmt.Errorf("failed to remove empty directory: %w", err)
				}
			}
		}

		return nil
	})

	if err != nil {
		fmt.Fprintf(&logBuffer, "Error during second pass: %v\n", err)
		return logBuffer.String(), err
	}

	fmt.Fprintln(&logBuffer, "Cleanup process completed successfully")
	return logBuffer.String(), nil
}
