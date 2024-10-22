package ui

import (
	"fmt"
	"os"
	"path/filepath"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/ildx/doubtfire/internal/utils"
)

func (m WelcomeModel) updateSetup(msg tea.Msg) (WelcomeModel, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			dir := m.setup.textInput.Value()
			homeDir, err := os.UserHomeDir()
			if err != nil {
				m.setup.err = fmt.Errorf("failed to get home directory: %w", err)
				return m, nil
			}
			fullPath := filepath.Join(homeDir, dir)
			err = os.MkdirAll(fullPath, os.ModePerm)
			if err != nil {
				m.setup.err = fmt.Errorf("failed to create destination directory: %w", err)
				return m, nil
			}
			config := &utils.Config{DestinationDir: dir}
			if err := utils.SaveConfig(config); err != nil {
				m.setup.err = fmt.Errorf("failed to save config: %w", err)
				return m, nil
			}
			m.setup.active = false
			return m, nil
		case tea.KeyCtrlC, tea.KeyEsc:
			m.setup.active = false
			return m, nil
		}
	}

	m.setup.textInput, cmd = m.setup.textInput.Update(msg)
	return m, cmd
}
