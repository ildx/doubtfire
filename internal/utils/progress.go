package utils

import (
	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type progressModel struct {
	progress progress.Model
	spinner  spinner.Model
	percent  float64
	total    int
	done     bool
}

func initialModel(total int) progressModel {
	return progressModel{
		progress: progress.New(progress.WithDefaultGradient()),
		spinner:  spinner.New(),
		total:    total,
	}
}

func (m progressModel) Init() tea.Cmd {
	return tea.Batch(m.spinner.Tick)
}

func (m progressModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "q" {
			return m, tea.Quit
		}
	case struct{}: // This case is for progress updates from the channel
		m.percent += 1.0 / float64(m.total)
		if m.percent >= 1.0 {
			m.done = true
			return m, tea.Quit
		}
		return m, m.progress.IncrPercent(1.0 / float64(m.total))
	}

	// Update spinner tick on each frame
	return m, tea.Batch(m.spinner.Tick)
}

func (m progressModel) View() string {
	if m.done {
		return ""
	}
	return lipgloss.JoinVertical(lipgloss.Left,
		m.spinner.View(),
		m.progress.ViewAs(m.percent),
	)
}
