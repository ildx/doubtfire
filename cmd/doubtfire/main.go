package main

import (
	"bufio"
	"flag"
	"os"

	"github.com/charmbracelet/log"
	"github.com/ildx/doubtfire/internal/cleanup"
	"github.com/ildx/doubtfire/internal/config"
	"github.com/ildx/doubtfire/internal/setup"
)

// Error messages
const (
	ErrLoadConfig = "Error loading configuration"
)

func main() {
	// Define a command-line flag for manual cleanup
	manual := flag.Bool("manual", false, "Manually trigger the cleanup")
	changeDir := flag.Bool("change-dir", false, "Change the destination directory")
	flag.Parse()

	// Load cfg
	cfg, err := config.LoadConfig(os.Getenv("HOME"))
	if err != nil {
		log.Error(ErrLoadConfig, err)
		return
	}

	if *changeDir {
		err := setup.ValidateAndChangeDirectory(cfg)
		if err != nil {
			log.Error(err)
		}
		return
	}

	// Check if destination directory is already set
	if cfg.DestinationDirectory == "" {
		// Prompt the user for the destination directory
		reader := bufio.NewReader(os.Stdin)
		err := setup.ValidateAndSetDirectory(cfg, reader)
		if err != nil {
			log.Error(err)
		}
	}

	cleanup.PerformCleanup(cfg, *manual)
}
