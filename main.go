package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"syscall"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	titleStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("#7D56F4")).Bold(true)
	cursorStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("#ffffff")).Background(lipgloss.Color("#5A56E0")).Bold(true)
	normalStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("#888888"))
	playingStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#00ff00")).Bold(true)
	pausedStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("#ffcc00")).Italic(true)
	footerStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("#444444")).Italic(true)
)

type model struct {
	files      []string
	err        error
	cursor     int
	choices    []string
	selected   map[int]struct{}
	processPid *os.Process
	isPaused   bool
}
type filesMsg []string
type errMsg struct{ error }

func (m model) Init() tea.Cmd {
	return readFilesCmd("/home/arcadian/Music")
}

func readFilesCmd(path string) tea.Cmd {
	return func() tea.Msg {
		entries, err := os.ReadDir(path)
		if err != nil {
			return errMsg{err}
		}

		var names []string
		for _, entry := range entries {
			names = append(names, entry.Name())
		}

		return filesMsg(names)
	}
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case filesMsg:
		m.files = msg
		m.choices = msg
		m.selected = make(map[int]struct{})
		return m, nil

	case errMsg:
		m.err = msg
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < len(m.choices)-1 {
				m.cursor++
			}
		case "enter":
			if m.processPid != nil {
				err := m.processPid.Kill()
				if err != nil {
					// log.Println("Failed to kill process:", err)
				} else {
					// log.Println("Killed process:", m.processPid.Pid)
				}
				m.processPid.Release()
				m.processPid = nil
			}

			file := m.files[m.cursor]
			cmd := exec.Command("pw-play", "/home/arcadian/Music/"+file)

			if err := cmd.Start(); err != nil {
				log.Printf("Error playing file: %v", err)
			} else {
				m.processPid = cmd.Process
				m.isPaused = false

				go func() {
					err := cmd.Wait()
					if err != nil {
						log.Println("Process exited with error:", err)
					} else {
						log.Println("Process finished:", cmd.Process.Pid)
					}
				}()
				// log.Println("Started process:", m.processPid.Pid)
			}

		case " ":
			if m.processPid != nil {
				var err error
				if m.isPaused {
					err = m.processPid.Signal(syscall.SIGCONT)
					if err != nil {
					} else {
						// log.Println("Resumed process:", m.processPid.Pid)
						m.isPaused = false
					}
				} else {
					err = m.processPid.Signal(syscall.SIGSTOP)
					if err != nil {
						// log.Println("Failed to pause process:", err)
					} else {
						// log.Println("Paused process:", m.processPid.Pid)
						m.isPaused = true
					}
				}
			}

		}

	}

	return m, nil
}

func (m model) View() string {
	var b strings.Builder

	// Title
	fmt.Fprintln(&b, titleStyle.Render("🎵 Music Player"))
	fmt.Fprintln(&b, "")

	if m.err != nil {
		return fmt.Sprintf("Error: %v\n", m.err)
	}

	if len(m.files) == 0 {
		return "Loading songs..."
	}

	// File List
	for i, file := range m.files {
		var line string

		switch {
		case i == m.cursor && m.isPaused:
			line = pausedStyle.Render("⏸ " + file)
		case i == m.cursor && m.processPid != nil:
			line = playingStyle.Render("▶ " + file)
		case i == m.cursor:
			line = cursorStyle.Render("> " + file)
		default:
			line = normalStyle.Render("  " + file)
		}

		fmt.Fprintln(&b, line)
	}

	// Footer
	fmt.Fprintln(&b, "")
	fmt.Fprintln(&b, footerStyle.Render("↑/↓ to move · Enter to play · Space to pause/resume · q to quit"))

	return b.String()
}

func main() {
	p := tea.NewProgram(model{})

	if err := p.Start(); err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
}
