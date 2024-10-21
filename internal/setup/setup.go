package setup

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/ildx/doubtfire/internal/config"
	"github.com/ildx/doubtfire/internal/tui"
	"github.com/ildx/doubtfire/internal/utils"
)

func ValidateAndSetDirectory(cfg *config.Config, reader *bufio.Reader) error {
	for {
		fmt.Print("Enter the destination directory: ")
		destDir, _ := reader.ReadString('\n')
		destDir = strings.TrimSpace(destDir)

		if err := validateDirectoryName(destDir); err != nil {
			fmt.Println(err)
			continue
		}

		destDir = expandPath(destDir)

		fmt.Println("Destination directory:", destDir)

		if err := createDirectory(destDir); err != nil {
			fmt.Println(err)
			continue
		}

		cfg.DestinationDirectory = destDir
		if err := config.SaveConfig(cfg); err != nil {
			fmt.Println("Error saving configuration:", err)
			continue
		}

		break
	}
	return nil
}

func ValidateAndChangeDirectory(cfg *config.Config) error {
	for {
		newDir, err := tui.New()
		if err != nil {
			return fmt.Errorf("error running TUI: %v", err)
		}

		if err := validateDirectoryName(newDir); err != nil {
			fmt.Println(err)
			continue
		}

		newDir = expandPath(newDir)

		fmt.Println("New destination directory:", newDir)

		if err := createDirectory(newDir); err != nil {
			fmt.Println(err)
			continue
		}

		if cfg.DestinationDirectory != "" && cfg.DestinationDirectory != newDir {
			if err := utils.CopyDir(cfg.DestinationDirectory, newDir); err != nil {
				fmt.Println("Error copying existing destination directory:", err)
				continue
			}

			if err := os.RemoveAll(cfg.DestinationDirectory); err != nil {
				fmt.Println("Error deleting old destination directory:", err)
				continue
			}
		}

		cfg.DestinationDirectory = newDir
		if err := config.SaveConfig(cfg); err != nil {
			fmt.Println("Error saving configuration:", err)
			continue
		}

		fmt.Println("New destination directory is set to:", cfg.DestinationDirectory)
		break
	}
	return nil
}

func validateDirectoryName(dir string) error {
	homeDir, _ := os.UserHomeDir()
	if dir == homeDir {
		return fmt.Errorf("the destination directory cannot be the home directory")
	} else if dir == "" {
		return fmt.Errorf("the destination directory cannot be empty")
	}
	return nil
}

func expandPath(path string) string {
	homeDir, _ := os.UserHomeDir()
	if strings.HasPrefix(path, "~") {
		return filepath.Join(homeDir, path[1:])
	} else if !filepath.IsAbs(path) {
		return filepath.Join(homeDir, path)
	}
	return path
}

func createDirectory(dir string) error {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if err := os.MkdirAll(dir, os.ModePerm); err != nil {
			return fmt.Errorf("error creating destination directory: %v", err)
		}
	}
	return nil
}
