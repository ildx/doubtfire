package ui

import (
	"fmt"
	"golang.org/x/term"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

type SummaryModel struct {
	title       string
	movedFiles  []string
	failedFiles []FailedFile
	quitting    bool
}

type FailedFile struct {
	file string
	err  error
}

func NewSummaryModel(moved []string, failed []FailedFile) *SummaryModel {
	return &SummaryModel{
		title:       "Doubtfile - The lovable housekeeper",
		movedFiles:  moved,
		failedFiles: failed,
	}
}

func (m SummaryModel) Init() tea.Cmd {
	return nil
}

func (m SummaryModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			m.quitting = true
			return m, tea.Quit
		case "b", "esc":
			return NewWelcomeModel(), tea.Batch(
				tea.EnterAltScreen,
				func() tea.Msg {
					width, height, _ := term.GetSize(int(os.Stdout.Fd()))
					return tea.WindowSizeMsg{
						Width:  width,
						Height: height,
					}
				},
			)
		}
	}
	return m, nil
}

func (m SummaryModel) View() string {
	var sb strings.Builder

	sb.WriteString(titleStyle("Doubtfire - The lovable housekeeper") + "\n\n")

	if len(m.movedFiles) > 0 {
		sb.WriteString(successStyle("Successfully moved files:\n"))
		for _, f := range m.movedFiles {
			sb.WriteString(fmt.Sprintf("  ✓ %s\n", f))
		}
		sb.WriteString("\n")
	}

	if len(m.failedFiles) > 0 {
		sb.WriteString(errorStyle("Failed to move files:") + "\n")
		for _, f := range m.failedFiles {
			sb.WriteString(fmt.Sprintf("  ✗ %s: %v\n", f.file, f.err))
		}
		sb.WriteString("\n")
	}

	if len(m.movedFiles) == 0 && len(m.failedFiles) == 0 {
		sb.WriteString("No files were processed\n\n")
	}

	// Command line
	sb.WriteString("\n\nb back • q quit")

	return sb.String()
}
