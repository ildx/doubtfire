package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type cleanupResult struct {
	movedFiles  []string
	failedFiles []FailedFile
}

func NewWelcomeModel() *WelcomeModel {
	items := []list.Item{
		listItem{title: "Setup", desc: "Set up the destination directory"},
		listItem{title: "Relocate", desc: "Move destination directory to a new location"},
		listItem{title: "Clean up", desc: "Move files to destination"},
		listItem{title: "Debug", desc: "Turn on/off helpful logs for solving problems"},
		listItem{title: "Quit", desc: "Exit the application"},
	}

	l := list.New(items, list.NewDefaultDelegate(), 0, 0)
	l.SetShowTitle(false)
	l.SetShowStatusBar(false)

	ti := textinput.New()
	ti.Placeholder = "Enter destination directory"
	ti.Focus()

	return &WelcomeModel{
		list:  l,
		title: "Doubtfire - The lovable housekeeper",
		setup: setupState{
			active:    false,
			textInput: ti,
			err:       nil,
		},
		showDebug: false,
	}
}

func NewWelcomeModelWithSize(width, height int) *WelcomeModel {
	m := NewWelcomeModel()
	m.list.SetWidth(width)
	m.list.SetHeight(height - 4) // Reserve space for title and command line
	return m
}

func (m WelcomeModel) Init() tea.Cmd {
	return nil
}

func (m WelcomeModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	m.debug = fmt.Sprintf("Received message type: %T\n", msg)

	if m.setup.active {
		var cmd tea.Cmd
		m, cmd = m.updateSetup(msg)
		return m, cmd
	}

	switch msg := msg.(type) {
	case *SummaryModel:
		return msg, nil
	case cleanupMsg:
		if msg.err != nil {
			m.lastMsg = errorMsg(fmt.Sprintf("Cleanup failed: %v\n\nLog:\n%s", msg.err, msg.log))
		} else {
			m.lastMsg = successMsg(fmt.Sprintf("Cleanup completed successfully\n\nLog:\n%s", msg.log))
		}
		return m, nil
	case errorMsg, successMsg:
		m.lastMsg = msg
		return m, nil
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			m.quitting = true
			return m, tea.Quit
		case "enter":
			i, ok := m.list.SelectedItem().(listItem)
			if ok {
				switch i.title {
				case "Quit":
					m.quitting = true
					return m, tea.Quit
				case "Setup":
					m.setup.active = true
					m.setup.textInput.Focus()
					return m, textinput.Blink
				case "Clean up":
					result := CleanUp()
					return NewSummaryModel(result.movedFiles, result.failedFiles), nil
				case "Debug":
					m.showDebug = !m.showDebug
					debugStatus := "on"
					if !m.showDebug {
						debugStatus = "off"
					}
					m.lastMsg = successMsg(fmt.Sprintf("Debug mode turned %s", debugStatus))
					return m, nil
				}
			}
		}
	case tea.WindowSizeMsg:
		m.list.SetWidth(msg.Width)
		listHeight := msg.Height - 4
		if listHeight < 0 {
			listHeight = 0
		}
		m.list.SetHeight(listHeight)
		return m, nil
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m WelcomeModel) View() string {
	if m.quitting {
		return "Thanks for using Doubtfire! Goodbye.\n"
	}

	if m.setup.active {
		errMsg := ""
		if m.setup.err != nil {
			errMsg = errorStyle(m.setup.err.Error())
		}
		return "Enter destination directory:\n" + m.setup.textInput.View() + "\n" + "Press Enter to confirm, Esc to cancel\n" + errMsg
	}

	var sb strings.Builder
	width, height := m.list.Width(), m.list.Height()

	// Render the title
	styledTitle := titleStyle(m.title)
	sb.WriteString(styledTitle + "\n")

	// Calculate remaining height for the list
	titleHeight := lipgloss.Height(styledTitle)
	remainingHeight := height - titleHeight - 4 // 2 for newlines, 2 for command line
	if remainingHeight < 0 {
		remainingHeight = 0
	}

	// Render list items
	listView := m.list.View()
	listLines := strings.Split(listView, "\n")
	if len(listLines) > remainingHeight {
		listLines = listLines[:remainingHeight]
	}
	if len(listLines) > 0 {
		sb.WriteString(strings.Join(listLines, "\n") + "\n")
	}

	// debug section
	if m.showDebug {
		debugContent := debugStyle(fmt.Sprintf("Debug:\nWidth: %d, Height: %d\n%s", width, height, m.debug))
		sb.WriteString(debugContent)
	} else {
		sb.WriteString("\n\n\n")
	}

	// status message handling
	var statusMsg string
	switch msg := m.lastMsg.(type) {
	case errorMsg:
		statusMsg = errorStyle(fmt.Sprintf("Error: %s", string(msg)))
	case successMsg:
		statusMsg = successStyle(fmt.Sprintf("Success: %s", string(msg)))
	}
	if statusMsg != "" {
		sb.WriteString("\n" + statusMsg + "\n\n")
	}

	if !m.showDebug && statusMsg == "" {
		sb.WriteString("\n\n\n")
	}

	// Fill remaining space
	currentHeight := lipgloss.Height(sb.String())
	if currentHeight < height-2 { // -2 for command line
		sb.WriteString(strings.Repeat("\n", height-2-currentHeight))
	}

	// Command line always at the bottom
	sb.WriteString("↑/k up • ↓/j down • / filter • q quit • ? more")

	return sb.String()
}
