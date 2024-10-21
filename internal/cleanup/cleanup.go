package cleanup

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/ildx/doubtfire/internal/config"
	"github.com/ildx/doubtfire/internal/utils"
)

func PerformCleanup(cfg *config.Config, manual bool) {
	// Check if today's cleanup has already been performed
	today := time.Now().Format("2006-01-02")
	lastCleanup := cfg.LastCleanupDate.Format("2006-01-02")
	if !manual && today == lastCleanup {
		fmt.Println("Today's cleanup has already been performed.")
		return
	}

	// Perform cleanup
	desktopPath := filepath.Join(os.Getenv("HOME"), "Desktop")
	files, err := os.ReadDir(desktopPath)
	if err != nil {
		fmt.Println("Error reading desktop directory:", err)
		return
	}

	// Create subfolders based on the current year and month
	year := time.Now().Format("2006")
	month := time.Now().Format("01")
	destDir := filepath.Join(cfg.DestinationDirectory, year, month)
	if err := utils.CreateDirectory(destDir); err != nil {
		fmt.Println(err)
		return
	}

	// Initialize counters for total files moved and total size cleaned
	totalFilesMoved, totalSizeCleaned := moveFiles(files, desktopPath, destDir)

	// Update last cleanup date
	cfg.LastCleanupDate = time.Now()
	err = config.SaveConfig(cfg)
	if err != nil {
		fmt.Println("Error updating last cleanup date:", err)
		return
	}

	fmt.Printf("Cleanup completed successfully.\nTotal files moved: %d\nTotal size cleaned: %d bytes\n", totalFilesMoved, totalSizeCleaned)
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
		fmt.Println("Moving file from:", srcPath, "to:", destPath)

		err := os.Rename(srcPath, destPath)
		if err != nil {
			fmt.Println("Error moving file:", file.Name(), err)
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
