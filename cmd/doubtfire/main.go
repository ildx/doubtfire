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
		for {
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

			// Validate the destination directory name
			homeDir, _ := os.UserHomeDir()
			if newDir == homeDir {
				fmt.Println("Warning: The destination directory cannot be the home directory.")
				continue
			} else if newDir == "" {
				fmt.Println("Warning: The destination directory cannot be empty.")
				continue
			}

			// Expand ~ to the full path of the home directory
			if strings.HasPrefix(newDir, "~") {
				newDir = filepath.Join(homeDir, newDir[1:])
			} else if !filepath.IsAbs(newDir) {
				// handle relative paths
				newDir = filepath.Join(homeDir, newDir)
			}

			// Print the new destination directory path
			fmt.Println("New destination directory:", newDir)

			// Create the new destination directory if it does not exist
			if _, err := os.Stat(newDir); os.IsNotExist(err) {
				err := os.MkdirAll(newDir, os.ModePerm)
				if err != nil {
					fmt.Println("Error creating new destination directory:", err)
					continue
				}
			}

			// Copy the contents of the old destination directory to the new destination directory
			if cfg.DestinationDirectory != "" && cfg.DestinationDirectory != newDir {
				err := utils.CopyDir(cfg.DestinationDirectory, newDir)
				if err != nil {
					fmt.Println("Error copying existing destination directory:", err)
					continue
				}

				// Delete the old destination directory
				err = os.RemoveAll(cfg.DestinationDirectory)
				if err != nil {
					fmt.Println("Error deleting old destination directory:", err)
					continue
				}
			}

			// Update the new destination directory in the JSON configuration file
			cfg.DestinationDirectory = newDir
			err = config.SaveConfig(cfg)
			if err != nil {
				fmt.Println("Error saving configuration:", err)
				continue
			}

			fmt.Println("New destination directory is set to:", cfg.DestinationDirectory)
			break
		}
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

		for {
			fmt.Print("Enter the destination directory: ")
			destDir, _ := reader.ReadString('\n')
			destDir = strings.TrimSpace(destDir)

			// Validate the destination directory name
			homeDir, _ := os.UserHomeDir()
			if destDir == homeDir {
				fmt.Println("Warning: The destination directory cannot be the home directory.")
				continue
			} else if destDir == "" {
				fmt.Println("Warning: The destination directory cannot be empty.")
				continue
			}

			// Expand ~ to the full path of the home directory
			if strings.HasPrefix(destDir, "~") {
				destDir = filepath.Join(homeDir, destDir[1:])
			} else if !filepath.IsAbs(destDir) {
				// handle relative paths
				destDir = filepath.Join(homeDir, destDir)
			}

			fmt.Println("Destination directory:", destDir)

			// Create the destination directory if it does not exist
			if _, err := os.Stat(destDir); os.IsNotExist(err) {
				err := os.MkdirAll(destDir, os.ModePerm)
				if err != nil {
					fmt.Println("Error creating destination directory:", err)
					continue
				}
			}

			// Save the destination directory
			cfg.DestinationDirectory = destDir
			err = config.SaveConfig(cfg)
			if err != nil {
				fmt.Println("Error saving configuration:", err)
				continue
			}

			break
		}
	}

	cleanup.PerformCleanup(cfg, *manual)
}
