package utils

import (
	"fmt"

	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
)

type progressModel struct {
	progress progress.Model
	current  int
	total    int
}

func (m progressModel) Init() tea.Cmd {
	return nil
}

func (m progressModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.Type == tea.KeyCtrlC {
			return m, tea.Quit
		}
	case int:
		m.current = msg
		percent := float64(m.current) / float64(m.total)
		m.progress.SetPercent(percent)
		if m.current >= m.total {
			return m, tea.Quit
		}
	}
	var cmd tea.Cmd
	newProgress, cmd := m.progress.Update(msg)
	m.progress = newProgress.(progress.Model)
	return m, cmd
}

func (m progressModel) View() string {
	return fmt.Sprintf("\nCopying files: %s\n", m.progress.View())
}
