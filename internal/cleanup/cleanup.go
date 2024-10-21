package cleanup

import (
	"os"
	"path/filepath"
	"time"

	"github.com/charmbracelet/log"
	"github.com/ildx/doubtfire/internal/config"
	"github.com/ildx/doubtfire/internal/errors"
	"github.com/ildx/doubtfire/internal/utils"
)

func PerformCleanup(cfg *config.Config, manual bool) {
	// Check if today's cleanup has already been performed
	today := time.Now().Format("2006-01-02")
	lastCleanup := cfg.LastCleanupDate.Format("2006-01-02")
	if !manual && today == lastCleanup {
		log.Info("Today's cleanup has already been performed.")
		return
	}

	// Perform cleanup
	desktopPath := filepath.Join(os.Getenv("HOME"), "Desktop")
	files, err := os.ReadDir(desktopPath)
	if err != nil {
		log.Error(errors.ErrReadDir, err)
		return
	}

	// Create subfolders based on the current year and month
	year := time.Now().Format("2006")
	month := time.Now().Format("01")
	destDir := filepath.Join(cfg.DestinationDirectory, year, month)
	if err := utils.CreateDirectory(destDir); err != nil {
		log.Error(errors.ErrCreateDir, err)
		return
	}

	// Initialize counters for total files moved and total size cleaned
	totalFilesMoved, totalSizeCleaned := moveFiles(files, desktopPath, destDir)

	// Update last cleanup date
	cfg.LastCleanupDate = time.Now()
	err = config.SaveConfig(cfg)
	if err != nil {
		log.Error(errors.ErrUpdateCleanupDate, err)
		return
	}

	log.Infof("Cleanup completed successfully.\nTotal files moved: %d\nTotal size cleaned: %d bytes\n", totalFilesMoved, totalSizeCleaned)
}

func moveFiles(files []os.DirEntry, srcDir, destDir string) (int, int64) {
	totalFilesMoved := 0
	totalSizeCleaned := int64(0)

	for _, file := range files {
		srcPath := filepath.Join(srcDir, file.Name())
		destPath := filepath.Join(destDir, file.Name())

		// Handle file name conflicts
		destPath = utils.ResolveFileNameConflict(destPath)

		// Print source and destination paths
		log.Infof("Moving file from: %s to: %s", srcPath, destPath)

		err := os.Rename(srcPath, destPath)
		if err != nil {
			log.Error(errors.ErrMoveFile, file.Name(), err)
			continue
		}

		// Update total files moved and total size cleaned
		totalFilesMoved++
		fileInfo, err := os.Stat(destPath)
		if err == nil {
			totalSizeCleaned += fileInfo.Size()
		}
	}

	return totalFilesMoved, totalSizeCleaned
}
