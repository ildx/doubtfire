package ui

import (
	"fmt"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

func NewWelcomeModel() *WelcomeModel {
	items := []list.Item{
		listItem{title: "Setup", desc: "Set up the destination directory"},
		listItem{title: "Relocate", desc: "Move destination directory to a new location"},
		listItem{title: "Clean up", desc: "Move files to destination"},
		listItem{title: "Quit", desc: "Exit the application"},
	}

	l := list.New(items, list.NewDefaultDelegate(), 0, 0)
	l.Title = "Doubtfire - The lovable housekeeping"
	l.SetShowStatusBar(false)

	ti := textinput.New()
	ti.Placeholder = "Enter destination directory"
	ti.Focus()

	return &WelcomeModel{
		list: l,
		setup: setupState{
			active:    false,
			textInput: ti,
			err:       nil,
		},
	}
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
					return m, func() tea.Msg {
						log, err := CleanUp()
						return cleanupMsg{log: log, err: err}
					}
				}
			}
		}
	case tea.WindowSizeMsg:
		m.list.SetWidth(msg.Width)
		m.list.SetHeight(msg.Height - 2)
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
			errMsg = m.setup.err.Error()
		}
		return "Enter destination directory:\n" + m.setup.textInput.View() + "\n" + "Press Enter to confirm, Esc to cancel\n" + errMsg
	}
	var statusMsg string
	switch msg := m.lastMsg.(type) {
	case errorMsg:
		statusMsg = fmt.Sprintf("Error: %s\n", string(msg))
	case successMsg:
		statusMsg = fmt.Sprintf("Success: %s\n", string(msg))
	}
	return m.list.View() + "\n" + statusMsg + "\n\nDebug:\n" + m.debug
}
