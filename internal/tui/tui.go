package tui

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
	input string
	done  bool
}

func initialModel() model {
	return model{}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			m.done = true
			return m, tea.Quit
		case "ctrl+c":
			return m, tea.Quit
		default:
			m.input += msg.String()
		}
	}
	return m, nil
}

func (m model) View() string {
	if m.done {
		return fmt.Sprintf("New destination directory: %s\n", m.input)
	}
	return fmt.Sprintf("Enter the new destination directory: %s", m.input)
}

func New() (string, error) {
	p := tea.NewProgram(initialModel())
	m, err := p.Run()
	if err != nil {
		return "", err
	}
	return m.(model).input, nil
}
