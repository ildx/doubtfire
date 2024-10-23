package ui

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/ildx/doubtfire/internal/utils"
)

func CleanUp() cleanupResult {
	var result cleanupResult

	config, err := utils.LoadConfig()
	if err != nil {
		result.failedFiles = append(result.failedFiles, FailedFile{
			file: "config",
			err:  fmt.Errorf("failed to load config: %w", err),
		})
		return result
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		result.failedFiles = append(result.failedFiles, FailedFile{
			file: "home directory",
			err:  fmt.Errorf("failed to get home directory: %w", err),
		})
		return result
	}

	desktopPath, err := utils.GetDesktopPath()
	if err != nil {
		result.failedFiles = append(result.failedFiles, FailedFile{
			file: "desktop",
			err:  fmt.Errorf("failed to get desktop path: %w", err),
		})
		return result
	}

	now := time.Now()
	destPath := filepath.Join(homeDir, config.DestinationDir, fmt.Sprintf("%d", now.Year()), fmt.Sprintf("%02d", now.Month()))

	// Create the destination directory
	err = os.MkdirAll(destPath, os.ModePerm)
	if err != nil {
		result.failedFiles = append(result.failedFiles, FailedFile{
			file: destPath,
			err:  fmt.Errorf("failed to create destination directory: %w", err),
		})
		return result
	}

	// First pass: move files and create directories
	err = filepath.Walk(desktopPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			result.failedFiles = append(result.failedFiles, FailedFile{
				file: path,
				err:  fmt.Errorf("error accessing path: %w", err),
			})
			return nil // Continue walking despite error
		}

		// Skip the desktop directory itself
		if path == desktopPath {
			return nil
		}

		relPath, err := filepath.Rel(desktopPath, path)
		if err != nil {
			result.failedFiles = append(result.failedFiles, FailedFile{
				file: path,
				err:  fmt.Errorf("failed to get relative path: %w", err),
			})
			return nil
		}

		newPath := filepath.Join(destPath, relPath)

		if info.IsDir() {
			uniquePath := utils.GetUniquePath(newPath)
			if err := os.MkdirAll(uniquePath, os.ModePerm); err != nil {
				result.failedFiles = append(result.failedFiles, FailedFile{
					file: path,
					err:  fmt.Errorf("failed to create directory: %w", err),
				})
			}
			return nil
		}

		if err := utils.MoveFile(path, newPath); err != nil {
			result.failedFiles = append(result.failedFiles, FailedFile{
				file: path,
				err:  fmt.Errorf("failed to move file: %w", err),
			})
		} else {
			result.movedFiles = append(result.movedFiles, path)
		}

		return nil
	})

	// Second pass: remove empty directories
	_ = filepath.Walk(desktopPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil // Skip problematic paths
		}

		// Skip the desktop directory itself
		if path == desktopPath {
			return nil
		}

		if info.IsDir() {
			empty, err := utils.IsDirEmpty(path)
			if err != nil {
				result.failedFiles = append(result.failedFiles, FailedFile{
					file: path,
					err:  fmt.Errorf("failed to check if directory is empty: %w", err),
				})
				return nil
			}
			if empty {
				if err := os.Remove(path); err != nil {
					result.failedFiles = append(result.failedFiles, FailedFile{
						file: path,
						err:  fmt.Errorf("failed to remove empty directory: %w", err),
					})
				}
			}
		}

		return nil
	})

	return result
}
