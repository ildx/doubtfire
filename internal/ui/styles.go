package ui

import (
	"github.com/charmbracelet/lipgloss"
)

var (
	errorStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("9")).Render
	successStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("10")).Render
	debugStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("241")).Render

	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("39")).                           // Light blue color
			Border(lipgloss.NormalBorder(), false, false, true, false). // Bottom border only
			MarginBottom(1).
			Padding(0, 1). // Add some padding
			Render
)
