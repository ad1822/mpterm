package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"syscall"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	titleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#A78BFA")). // Purple-400
			Bold(true).
			MarginBottom(1)

	cursorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#1E1B4B")). // Dark purple
			Background(lipgloss.Color("#C4B5FD")). // Purple-200
			Bold(true)

	normalStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#94A3B8")) // Slate-400

	playingStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#4ADE80")). // Green-400
			Bold(true)

	pausedStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FACC15")). // Yellow-400
			Italic(true)

	footerStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#64748B")). // Slate-500
			Italic(true).
			MarginTop(1)
)

type model struct {
	files          []string
	err            error
	cursor         int
	choices        []string
	selected       map[int]struct{}
	processPid     *os.Process
	isPaused       bool
	currentPlaying int // Track currently playing song index
	// width          int
	// height         int
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
				}
				m.processPid.Release()
				m.processPid = nil
				m.currentPlaying = -1
			}

			file := m.files[m.cursor]
			cmd := exec.Command("pw-play", "/home/arcadian/Music/"+file)

			if err := cmd.Start(); err != nil {
				// log.Printf("Error playing file: %v", err)
			} else {
				m.processPid = cmd.Process
				m.isPaused = false
				m.currentPlaying = m.cursor

				go func() {
					err := cmd.Wait()
					if err != nil {
						// log.Println("Process exited with error:", err)
					}
					m.processPid = nil
					m.currentPlaying = -1
				}()
			}

		case " ":
			if m.processPid != nil {
				var err error
				if m.isPaused {
					err = m.processPid.Signal(syscall.SIGCONT)
					if err == nil {
						m.isPaused = false
					}
				} else {
					err = m.processPid.Signal(syscall.SIGSTOP)
					if err == nil {
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
	fmt.Fprintf(&b, "%s\n", titleStyle.Render("🎵 Music Player"))

	if m.err != nil {
		return fmt.Sprintf("Error: %v\n", m.err)
	}

	if len(m.files) == 0 {
		return "Loading songs..."
	}

	// File List
	for i, file := range m.files {
		var line string

		// Apply cursor style first (background)
		if i == m.cursor {
			line = cursorStyle.Render("  " + file + "|")

			// Then apply playing/paused style if needed (text color)
			if i == m.currentPlaying {
				if m.isPaused {
					line = pausedStyle.Render(line)
				} else {
					line = playingStyle.Render(line)
				}
			}
		} else if i == m.currentPlaying {
			// Currently playing song but not selected
			if m.isPaused {
				line = pausedStyle.Render("  " + file)
			} else {
				line = playingStyle.Render("  " + file)
			}
		} else {
			line = normalStyle.Render("  " + file)
		}

		fmt.Fprintln(&b, line)
	}

	// Footer
	// fmt.Fprintf(&b, "\n%s", footerStyle.Render(
	// 	"↑/↓ to move • Enter to play • Space to pause/resume • q to quit",
	// ))

	return b.String()
}

func main() {
	p := tea.NewProgram(model{}, tea.WithAltScreen())

	if err := p.Start(); err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
}
