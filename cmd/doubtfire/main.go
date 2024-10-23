package main

import (
	"fmt"
	"os"

	"github.com/ildx/doubtfire/internal/app"
)

func main() {
	fmt.Print("\033[2J") // clear screen
	fmt.Print("\033[H")  // move cursor to top-left
	if err := app.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
