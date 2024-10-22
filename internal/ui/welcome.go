package ui

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

type WelcomeModel struct {
	list     list.Model
	quitting bool
}

func NewWelcomeModel() *WelcomeModel {
	items := []list.Item{
		listItem{title: "Setup", desc: "Set up the destination directory"},
		listItem{title: "Relocate", desc: "Move destination directory to a new location"},
		listItem{title: "Clean up", desc: "Move files to destination"},
		listItem{title: "Reset", desc: "Reset app the default state"},
		listItem{title: "Quit", desc: "Exit the application"},
	}

	l := list.New(items, list.NewDefaultDelegate(), 0, 0)
	l.Title = "Doubtfire - The lovable housekeeping"
	l.SetShowStatusBar(false)

	return &WelcomeModel{list: l}
}

func (m WelcomeModel) Init() tea.Cmd {
	return tea.Batch(
		tea.ClearScreen,
	)
}

func (m WelcomeModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			m.quitting = true
			return m, tea.Quit
		case "enter":
			i, ok := m.list.SelectedItem().(listItem)
			if ok && i.title == "Quit" {
				m.quitting = true
				return m, tea.Quit
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
	return m.list.View() + "\n"
}

type listItem struct {
	title, desc string
}

func (i listItem) Title() string       { return i.title }
func (i listItem) Description() string { return i.desc }
func (i listItem) FilterValue() string { return i.title }
