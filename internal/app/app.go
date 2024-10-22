package app

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/ildx/doubtfire/internal/ui"
)

func Run() error {
	p := tea.NewProgram(ui.NewWelcomeModel())
	_, err := p.Run()
	return err
}
