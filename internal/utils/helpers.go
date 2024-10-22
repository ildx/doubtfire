package utils

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func DirectoryExists(path string) (bool, error) {
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return info.IsDir(), nil
}

func IsDirEmpty(path string) (bool, error) {
	f, err := os.Open(path)
	if err != nil {
		return false, err
	}
	defer f.Close()

	_, err = f.Readdirnames(1)
	if err == io.EOF {
		return true, nil
	}
	return false, err
}

func GetDesktopPath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(homeDir, "Desktop"), nil
}

func MoveFile(src, dest string) error {
	destDir := filepath.Dir(dest)
	if err := os.MkdirAll(destDir, os.ModePerm); err != nil {
		return fmt.Errorf("failed to create destination directory: %w", err)
	}

	uniqueDest := GetUniquePath(dest)

	srcFile, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("failed to open source file: %w", err)
	}
	defer srcFile.Close()

	destFile, err := os.Create(uniqueDest)
	if err != nil {
		return fmt.Errorf("failed to create destination file: %w", err)
	}
	defer destFile.Close()

	_, err = srcFile.Seek(0, 0)
	if err != nil {
		return fmt.Errorf("failed to seek to beginning of source file: %w", err)
	}

	_, err = destFile.Seek(0, 0)
	if err != nil {
		return fmt.Errorf("failed to seek to beginning of destination file: %w", err)
	}

	_, err = io.Copy(destFile, srcFile)
	if err != nil {
		return fmt.Errorf("failed to copy file: %w", err)
	}

	err = os.Remove(src)
	if err != nil {
		return fmt.Errorf("failed to remove source file: %w", err)
	}

	return nil
}

func GetUniquePath(path string) string {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return path
	}

	dir, file := filepath.Split(path)
	ext := filepath.Ext(file)
	name := strings.TrimSuffix(file, ext)

	for i := 1; ; i++ {
		var newPath string
		if i == 1 {
			newPath = filepath.Join(dir, name+" copy"+ext)
		} else {
			newPath = filepath.Join(dir, fmt.Sprintf("%s copy %d%s", name, i, ext))
		}
		if _, err := os.Stat(newPath); os.IsNotExist(err) {
			return newPath
		}
	}
}
