package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/ildx/doubtfire/internal/cleanup"
	"github.com/ildx/doubtfire/internal/config"
	"github.com/ildx/doubtfire/internal/tui"
	"github.com/ildx/doubtfire/internal/utils"
)

func main() {
	// Define a command-line flag for manual cleanup
	manual := flag.Bool("manual", false, "Manually trigger the cleanup")
	changeDir := flag.Bool("change-dir", false, "Change the destination directory")
	flag.Parse()

	if *changeDir {
		newDir, err := tui.New()
		if err != nil {
			fmt.Printf("Error running TUI: %v\n", err)
			return
		}

		// Load cfg
		cfg, err := config.LoadConfig()
		if err != nil {
			fmt.Println("Error loading configuration:", err)
			return
		}

		// Expand ~ to the full path of the home directory
		if strings.HasPrefix(newDir, "~") {
			homeDir, _ := os.UserHomeDir()
			newDir = filepath.Join(homeDir, newDir[1:])
		} else if !filepath.IsAbs(newDir) {
			// handle relative paths
			homeDir, _ := os.UserHomeDir()
			newDir = filepath.Join(homeDir, newDir)
		}

		// Print the new destination directory path
		fmt.Println("New destination directory:", newDir)

		// Create the new destination directory if it does not exist
		if _, err := os.Stat(newDir); os.IsNotExist(err) {
			err := os.MkdirAll(newDir, os.ModePerm)
			if err != nil {
				fmt.Println("Error creating new destination directory:", err)
				return
			}
		}

		// Copy the contents of the old destination directory to the new destination directory
		if cfg.DestinationDirectory != "" && cfg.DestinationDirectory != newDir {
			err := utils.CopyDir(cfg.DestinationDirectory, newDir)
			if err != nil {
				fmt.Println("Error copying existing destination directory:", err)
				return
			}

			// Delete the old destination directory
			err = os.RemoveAll(cfg.DestinationDirectory)
			if err != nil {
				fmt.Println("Error deleting old destination directory:", err)
				return
			}
		}

		// Update the new destination directory in the JSON configuration file
		cfg.DestinationDirectory = newDir
		err = config.SaveConfig(cfg)
		if err != nil {
			fmt.Println("Error saving configuration:", err)
			return
		}

		fmt.Println("New destination directory is set to:", cfg.DestinationDirectory)
		return
	}

	// Load cfg
	cfg, err := config.LoadConfig()
	if err != nil {
		fmt.Println("Error loading configuration:", err)
		return
	}

	// Check if destination directory is already set
	if cfg.DestinationDirectory == "" {
		// Prompt the user for the destination directory
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Enter the destination directory: ")
		destDir, _ := reader.ReadString('\n')
		destDir = strings.TrimSpace(destDir)

		// Expand ~ to the full path of the home directory
		if strings.HasPrefix(destDir, "~") {
			homeDir, _ := os.UserHomeDir()
			destDir = filepath.Join(homeDir, destDir[1:])
		} else if !filepath.IsAbs(destDir) {
			// handle relative paths
			homeDir, _ := os.UserHomeDir()
			destDir = filepath.Join(homeDir, destDir)
		}

		fmt.Println("Destination directory:", destDir)

		// Create the destination directory if it does not exist
		if _, err := os.Stat(destDir); os.IsNotExist(err) {
			err := os.MkdirAll(destDir, os.ModePerm)
			if err != nil {
				fmt.Println("Error creating destination directory:", err)
				return
			}
		}

		// Save the destination directory
		cfg.DestinationDirectory = destDir
		err = config.SaveConfig(cfg)
		if err != nil {
			fmt.Println("Error saving configuration:", err)
			return
		}
	}

	cleanup.PerformCleanup(cfg, *manual)
}
