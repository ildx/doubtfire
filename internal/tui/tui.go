package tui

import (
	"fmt"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type model struct {
	textInput textinput.Model
	done      bool
}

func initialModel() model {
	ti := textinput.New()
	ti.Placeholder = "Enter the new destination directory"
	ti.Focus()
	ti.CharLimit = 156
	ti.Width = 20

	return model{
		textInput: ti,
		done:      false,
	}
}

func (m model) Init() tea.Cmd {
	return textinput.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			m.done = true
			return m, tea.Quit
		case "esc", "ctrl+c":
			m.done = true
			return m, tea.Quit
		}
	}

	m.textInput, cmd = m.textInput.Update(msg)
	return m, cmd
}

func (m model) View() string {
	if m.done {
		return fmt.Sprintf("New destination directory: %s\n", m.textInput.Value())
	}

	titleStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("205")).
		Background(lipgloss.Color("63")).
		Padding(1, 2).
		MarginBottom(1).
		Bold(true)

	inputStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("63")).
		Background(lipgloss.Color("205")).
		Padding(1, 2).
		Width(50)

	return lipgloss.JoinVertical(lipgloss.Left,
		titleStyle.Render("Change Destination Directory"),
		inputStyle.Render(m.textInput.View()),
	)
}

func New() (string, error) {
	p := tea.NewProgram(initialModel())
	m, err := p.Run()
	if err != nil {
		return "", err
	}
	return m.(model).textInput.Value(), nil
}
