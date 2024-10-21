package utils

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/log"
	"github.com/ildx/doubtfire/internal/errors"
)

type progressMsg int

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
	totalItems, err := countItems(src)
	if err != nil {
		return err
	}

	prog := progress.New(progress.WithDefaultGradient())
	prog.Width = 40

	m := progressModel{
		progress: prog,
		total:    totalItems,
	}

	updateChan := make(chan progressMsg)
	doneChan := make(chan struct{})
	p := tea.NewProgram(m)

	go func() {
		if _, err := p.Run(); err != nil {
			fmt.Printf("Error running progress bar: %v\n", err)
		}
		close(doneChan)
	}()

	go func() {
		for msg := range updateChan {
			p.Send(msg)
		}
	}()

	err = copyDirRecursive(src, dst, updateChan)
	close(updateChan)

	if err != nil {
		p.Quit()
		<-doneChan
		return err
	}

	// Ensure progress reaches 100%
	p.Send(progressMsg(totalItems))
	p.Quit()
	<-doneChan

	return nil
}

func copyDirRecursive(src, dst string, updateChan chan<- progressMsg) error {
	entries, err := os.ReadDir(src)
	if err != nil {
		return err
	}

	err = os.MkdirAll(dst, os.ModePerm)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		srcPath := filepath.Join(src, entry.Name())
		dstPath := filepath.Join(dst, entry.Name())

		if entry.IsDir() {
			err = copyDirRecursive(srcPath, dstPath, updateChan)
			if err != nil {
				return err
			}
		} else {
			err = CopyFile(srcPath, dstPath)
			if err != nil {
				return err
			}
			select {
			case updateChan <- 1: // Increment progress
			default:
				// Channel is closed or full, do nothing
			}
		}
	}

	return nil
}

func countItems(dir string) (int, error) {
	count := 0
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			count++
		}
		return nil
	})
	return count, err
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
