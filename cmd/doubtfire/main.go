package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"

	"github.com/ildx/doubtfire/internal/cleanup"
	"github.com/ildx/doubtfire/internal/config"
	"github.com/ildx/doubtfire/internal/setup"
)

func main() {
	// Define a command-line flag for manual cleanup
	manual := flag.Bool("manual", false, "Manually trigger the cleanup")
	changeDir := flag.Bool("change-dir", false, "Change the destination directory")
	flag.Parse()

	// Load cfg
	cfg, err := config.LoadConfig()
	if err != nil {
		fmt.Println("Error loading configuration:", err)
		return
	}

	if *changeDir {
		err := setup.ValidateAndChangeDirectory(cfg)
		if err != nil {
			fmt.Println(err)
		}
		return
	}

	// Check if destination directory is already set
	if cfg.DestinationDirectory == "" {
		// Prompt the user for the destination directory
		reader := bufio.NewReader(os.Stdin)
		err := setup.ValidateAndSetDirectory(cfg, reader)
		if err != nil {
			fmt.Println(err)
		}
	}

	cleanup.PerformCleanup(cfg, *manual)
}
