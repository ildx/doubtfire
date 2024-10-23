package ui

import (
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type ModelState int

const (
	StateWelcome ModelState = iota
	StateSetup
	StateCleanupSummary
)

type WelcomeModel struct {
	list      list.Model
	title     string
	quitting  bool
	setup     setupState
	lastMsg   tea.Msg
	debug     string
	showDebug bool
	cleanup   CleanupModel
	state     ModelState
}

type setupState struct {
	active    bool
	textInput textinput.Model
	err       error
}

type CleanupModel struct {
	progress    progress.Model
	message     string
	statusLog   string
	files       []string
	currentFile string
	finished    bool
	summary     CleanupSummary
}

type CleanupSummary struct {
	sourceDir       string
	destDir         string
	successfulMoves []string
	failedMoves     []string
}

type listItem struct {
	title, desc string
}

func (i listItem) Title() string       { return i.title }
func (i listItem) Description() string { return i.desc }
func (i listItem) FilterValue() string { return i.title }

type errorMsg string
type successMsg string
type cleanupMsg struct {
	log string
	err error
}
