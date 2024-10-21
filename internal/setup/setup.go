package setup

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/log"
	"github.com/ildx/doubtfire/internal/config"
	"github.com/ildx/doubtfire/internal/errors"
	"github.com/ildx/doubtfire/internal/tui"
	"github.com/ildx/doubtfire/internal/utils"
)

func ValidateAndSetDirectory(cfg *config.Config, reader *bufio.Reader) error {
	for {
		log.Info("Enter the destination directory: ")
		destDir, _ := reader.ReadString('\n')
		destDir = strings.TrimSpace(destDir)

		if err := validateDirectoryName(destDir); err != nil {
			log.Warn(err.Error())
			continue
		}

		destDir = expandPath(destDir)

		log.Info("Destination directory:", "path", destDir)

		if err := createDirectory(destDir); err != nil {
			log.Error(err.Error())
			continue
		}

		cfg.DestinationDirectory = destDir
		if err := config.SaveConfig(cfg); err != nil {
			log.Error(errors.ErrSaveConfig, err)
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
			log.Error(errors.ErrRunTUI, err)
			return err
		}

		if err := validateDirectoryName(newDir); err != nil {
			log.Warn(err.Error())
			continue
		}

		newDir = expandPath(newDir)

		log.Info("New destination directory:", "path", newDir)

		if err := createDirectory(newDir); err != nil {
			log.Error(err.Error())
			continue
		}

		if cfg.DestinationDirectory != "" && cfg.DestinationDirectory != newDir {
			if err := utils.CopyDir(cfg.DestinationDirectory, newDir); err != nil {
				log.Error(errors.ErrCopyDir, err)
				continue
			}

			if err := os.RemoveAll(cfg.DestinationDirectory); err != nil {
				log.Error(errors.ErrDeleteOldDir, err)
				continue
			}
		}

		cfg.DestinationDirectory = newDir
		if err := config.SaveConfig(cfg); err != nil {
			log.Error(errors.ErrSaveConfig, err)
			continue
		}

		log.Info("New destination directory is set to:", "path", cfg.DestinationDirectory)
		break
	}
	return nil
}

func validateDirectoryName(dir string) error {
	homeDir, _ := os.UserHomeDir()
	if dir == homeDir {
		log.Error(errors.ErrHomeDir)
		return fmt.Errorf(errors.ErrHomeDir)
	} else if dir == "" {
		log.Error(errors.ErrEmptyDir)
		return fmt.Errorf(errors.ErrEmptyDir)
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
			log.Error(errors.ErrCreateDir, err)
			return fmt.Errorf("%s: %v", errors.ErrCreateDir, err)
		}
	}
	return nil
}