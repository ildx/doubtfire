package ui

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/ildx/doubtfire/internal/utils"
)

type errorMsg string
type successMsg string

type listItem struct {
	title, desc string
}

type WelcomeModel struct {
	list     list.Model
	quitting bool
	setup    setupState
	lastMsg  tea.Msg
	debug    string
	cleanup  CleanupModel
	state    ModelState
}

type ModelState int

type setupState struct {
	active    bool
	textInput textinput.Model
	err       error
}

type cleanupMsg struct {
	log string
	err error
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

const (
	StateWelcome ModelState = iota
	StateSetup
	StateCleanupSummary
)

func NewCleanupModel() CleanupModel {
	return CleanupModel{
		progress: progress.New(progress.WithDefaultGradient()),
		message:  "Clean up in progress. Moving files...",
	}
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

func (i listItem) Title() string       { return i.title }
func (i listItem) Description() string { return i.desc }
func (i listItem) FilterValue() string { return i.title }

func CleanUp() (string, error) {
	var logBuffer strings.Builder
	fmt.Fprintln(&logBuffer, "Starting cleanup process")

	config, err := utils.LoadConfig()
	if err != nil {
		return logBuffer.String(), fmt.Errorf("failed to load config: %w", err)
	}
	fmt.Fprintf(&logBuffer, "Loaded config, destination directory: %s\n", config.DestinationDir)

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return logBuffer.String(), fmt.Errorf("failed to get home directory: %w", err)
	}

	desktopPath, err := getDesktopPath()
	if err != nil {
		return logBuffer.String(), fmt.Errorf("failed to get desktop path: %w", err)
	}
	fmt.Fprintf(&logBuffer, "Desktop path: %s\n", desktopPath)

	now := time.Now()
	destPath := filepath.Join(homeDir, config.DestinationDir, fmt.Sprintf("%d", now.Year()), fmt.Sprintf("%02d", now.Month()))
	fmt.Fprintf(&logBuffer, "Destination path: %s\n", destPath)

	// Create the destination directory
	err = os.MkdirAll(destPath, os.ModePerm)
	if err != nil {
		return logBuffer.String(), fmt.Errorf("failed to create destination directory: %w", err)
	}

	// First pass: move files and create directories
	err = filepath.Walk(desktopPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Fprintf(&logBuffer, "Error accessing path %s: %v\n", path, err)
			return err
		}

		// Skip the desktop directory itself
		if path == desktopPath {
			return nil
		}

		relPath, err := filepath.Rel(desktopPath, path)
		if err != nil {
			fmt.Fprintf(&logBuffer, "Error getting relative path for %s: %v\n", path, err)
			return fmt.Errorf("failed to get relative path: %w", err)
		}

		newPath := filepath.Join(destPath, relPath)

		if info.IsDir() {
			uniquePath := getUniquePath(newPath)
			fmt.Fprintf(&logBuffer, "Creating directory: %s\n", uniquePath)
			return os.MkdirAll(uniquePath, os.ModePerm)
		}

		fmt.Fprintf(&logBuffer, "Moving file: %s to %s\n", path, newPath)
		return moveFile(path, newPath)
	})

	if err != nil {
		fmt.Fprintf(&logBuffer, "Error during first pass: %v\n", err)
		return logBuffer.String(), err
	}

	// Second pass: remove empty directories
	err = filepath.Walk(desktopPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Fprintf(&logBuffer, "Error accessing path %s: %v\n", path, err)
			return err
		}

		// Skip the desktop directory itself
		if path == desktopPath {
			return nil
		}

		if info.IsDir() {
			empty, err := isDirEmpty(path)
			if err != nil {
				fmt.Fprintf(&logBuffer, "Error checking if directory is empty %s: %v\n", path, err)
				return fmt.Errorf("failed to check if directory is empty: %w", err)
			}
			if empty {
				fmt.Fprintf(&logBuffer, "Removing empty directory: %s\n", path)
				if err := os.Remove(path); err != nil {
					fmt.Fprintf(&logBuffer, "Error removing empty directory %s: %v\n", path, err)
					return fmt.Errorf("failed to remove empty directory: %w", err)
				}
			}
		}

		return nil
	})

	if err != nil {
		fmt.Fprintf(&logBuffer, "Error during second pass: %v\n", err)
		return logBuffer.String(), err
	}

	fmt.Fprintln(&logBuffer, "Cleanup process completed successfully")
	return logBuffer.String(), nil
}

func isDirEmpty(path string) (bool, error) {
	f, err := os.Open(path)
	if err != nil {
		return false, err
	}
	defer f.Close()

	_, err = f.Readdirnames(1)
	if err == io.EOF {
		return true, nil
	}
	return false, err
}

func getDesktopPath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(homeDir, "Desktop"), nil
}

func moveFile(src, dest string) error {
	destDir := filepath.Dir(dest)
	if err := os.MkdirAll(destDir, os.ModePerm); err != nil {
		return fmt.Errorf("failed to create destination directory: %w", err)
	}

	uniqueDest := getUniquePath(dest)

	srcFile, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("failed to open source file: %w", err)
	}
	defer srcFile.Close()

	destFile, err := os.Create(uniqueDest)
	if err != nil {
		return fmt.Errorf("failed to create destination file: %w", err)
	}
	defer destFile.Close()

	_, err = srcFile.Seek(0, 0)
	if err != nil {
		return fmt.Errorf("failed to seek to beginning of source file: %w", err)
	}

	_, err = destFile.Seek(0, 0)
	if err != nil {
		return fmt.Errorf("failed to seek to beginning of destination file: %w", err)
	}

	_, err = io.Copy(destFile, srcFile)
	if err != nil {
		return fmt.Errorf("failed to copy file: %w", err)
	}

	err = os.Remove(src)
	if err != nil {
		return fmt.Errorf("failed to remove source file: %w", err)
	}

	return nil
}

func getUniquePath(path string) string {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return path
	}

	dir, file := filepath.Split(path)
	ext := filepath.Ext(file)
	name := strings.TrimSuffix(file, ext)

	for i := 1; ; i++ {
		var newPath string
		if i == 1 {
			newPath = filepath.Join(dir, name+" copy"+ext)
		} else {
			newPath = filepath.Join(dir, fmt.Sprintf("%s copy %d%s", name, i, ext))
		}
		if _, err := os.Stat(newPath); os.IsNotExist(err) {
			return newPath
		}
	}
}
