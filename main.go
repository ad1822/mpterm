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
			Foreground(lipgloss.Color("#A78BFA")).
			Bold(true).
			MarginBottom(1)

	cursorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#f5e0dc")).
			Bold(true)

	queueCursorStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#f2cdcd")).
				Bold(true)

	normalStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#94A3B8"))

	playingStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#4ADE80")).
			Bold(true)

	pausedStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FACC15")).
			Bold(true)

	footerStyle = lipgloss.NewStyle()

	border = lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("228")).
		BorderBackground(lipgloss.Color("63")).
		BorderTop(true).
		BorderLeft(true)
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
	width          int
	height         int
	currentSong    string
	queue          []string
	queueCursor    int
	activePanel    int // 0 for main list, 1 for queue
}

type filesMsg []string
type errMsg struct{ error }

func getMusicDir() (string,error) {
	cmd := exec.Command("xdg-user-dir","MUSIC") // use xdg-user-dir for getting music dir
	output, err := cmd.Output()
	if err != nil{
		return "",err
	}
	musicDir := strings.TrimSpace(string(output)) + "/"
	return musicDir,nil
}



func (m model) Init() tea.Cmd {
	musicDir,err := getMusicDir()
	if err != nil{
		fmt.Printf("%s\n",err)
		return readFilesCmd("/home/arcadian/Music")
	}
	return readFilesCmd(musicDir)
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
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil
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
			if m.processPid != nil {
				_ = m.processPid.Kill()
				_ = m.processPid.Release()
				m.processPid = nil
			}
			return m, tea.Quit
		case "tab":
			// Switch between panels
			m.activePanel = (m.activePanel + 1) % 2
		case "up", "k":
			if m.activePanel == 0 && m.cursor > 0 {
				m.cursor--
			} else if m.activePanel == 1 && m.queueCursor > 0 {
				m.queueCursor--
			}
		case "down", "j":
			if m.activePanel == 0 && m.cursor < len(m.choices)-1 {
				m.cursor++
			} else if m.activePanel == 1 && m.queueCursor < len(m.queue)-1 {
				m.queueCursor++
			}
		case "a":
			// Add to queue
			if len(m.files) > 0 {
				m.queue = append(m.queue, m.files[m.cursor])
			}
		case "d":
			// Remove from queue
			if len(m.queue) > 0 && m.queueCursor >= 0 && m.queueCursor < len(m.queue) {
				// If we're removing the currently playing song, stop it
				if m.currentPlaying >= 0 && m.queueCursor == m.currentPlaying {
					if m.processPid != nil {
						_ = m.processPid.Kill()
						_ = m.processPid.Release()
						m.processPid = nil
					}
					m.currentPlaying = -1
					m.currentSong = ""
				}
				m.queue = append(m.queue[:m.queueCursor], m.queue[m.queueCursor+1:]...)
				if m.queueCursor >= len(m.queue) && len(m.queue) > 0 {
					m.queueCursor = len(m.queue) - 1
				}
				if len(m.queue) == 0 {
					m.queueCursor = 0
				}
				// Adjust currentPlaying if it was after the removed item
				if m.currentPlaying > m.queueCursor {
					m.currentPlaying--
				}
			}
		case "enter":
			if m.activePanel == 0 {
				// Play from main list
				m.playSong(m.files[m.cursor], -1)
			} else if m.activePanel == 1 && len(m.queue) > 0 {
				// Play from queue
				m.playSong(m.queue[m.queueCursor], m.queueCursor)
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
		case "l":
			// Play next in queue
			if len(m.queue) > 0 {
				nextIndex := m.currentPlaying + 1
				if nextIndex < len(m.queue) {
					m.playSong(m.queue[nextIndex], nextIndex)
					m.queueCursor = nextIndex
				}
			}
		case "h":
			// Play previous in queue
			if len(m.queue) > 0 {
				prevIndex := m.currentPlaying - 1
				if prevIndex >= 0 {
					m.playSong(m.queue[prevIndex], prevIndex)
					m.queueCursor = prevIndex
				}
			}
		case "s":
			// Stop playback
			if m.processPid != nil {
				_ = m.processPid.Kill()
				_ = m.processPid.Release()
				m.processPid = nil
				m.currentSong = ""
				m.currentPlaying = -1
				m.isPaused = false
			}
		}
	}

	return m, nil
}

func (m *model) playSong(filename string, queueIndex int) {
	if m.processPid != nil {
		_ = m.processPid.Kill()
		_ = m.processPid.Release()
		m.processPid = nil
	}
	musicDir,err := getMusicDir()
	if err != nil{
		fmt.Printf("%s\n",err)
		musicDir = "/home/arcadian/Music/"
	}
	cmd := exec.Command("pw-play", musicDir+filename)

	if err := cmd.Start(); err != nil {
		// Error handling would go here
	} else {
		m.currentSong = filename
		m.processPid = cmd.Process
		m.isPaused = false
		m.currentPlaying = queueIndex

		go func() {
			err := cmd.Wait()
			if err != nil {
				// Error handling would go here
			}
			m.processPid = nil
			m.currentPlaying = -1
		}()
	}
}

func renderQueue(m model) string {
	if len(m.queue) == 0 {
		return "Queue is empty (press 'a' to add)"
	}

	var b strings.Builder
	for i, entry := range m.queue {
		var line string
		// prefix := "   " // default spacing

		// Apply cursor style if this is the selected item in the queue
		if i == m.queueCursor && m.activePanel == 1 {
			line = queueCursorStyle.Render("  " + entry)
			// prefix = "  " // one less space for cursor
		} else {
			line = normalStyle.Render("  " + entry)
		}

		// Apply playing indicator if this is the currently playing track
		if i == m.currentPlaying {
			if m.isPaused {
				line = pausedStyle.Render("⏸ " + entry)
			} else {
				line = playingStyle.Render("▶ " + entry)
			}
		}

		b.WriteString(line + "\n")
	}
	return b.String()
}

func renderSongList(m model) string {
	if len(m.files) == 0 {
		return "No songs found"
	}

	var b strings.Builder
	for i, file := range m.files {
		var line string
		// prefix := "   " // default spacing

		// Apply cursor style
		if i == m.cursor && m.activePanel == 0 {
			line = cursorStyle.Render("  " + file)
			// prefix := "  " // one less space for cursor
		} else {
			line = normalStyle.Render("  " + file)
		}

		// Check if this is the currently playing song
		if file == m.currentSong {
			if m.isPaused {
				line = pausedStyle.Render("⏸ " + file)
			} else {
				line = playingStyle.Render("▶ " + file)
			}
		}

		b.WriteString(line + "\n")
	}
	return b.String()
}

func (m model) View() string {
	mainHeight := m.height - 2 // Adjust for title + status bar
	mainWidth := m.width - 4
	leftWidth := (mainWidth / 2)

	// Panel border colors (active panel highlight)
	leftBorderColor, rightBorderColor := "0", "0"
	if m.activePanel == 0 {
		leftBorderColor = "#cba6f7"
	} else {
		rightBorderColor = "#cba6f7"
	}

	// Title (centered at the top)
	title := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#cba6f7")). // Purple accent
		Align(lipgloss.Center).
		Width(mainWidth).
		Render("🎵 Music Player 🎵")

	// Left panel (with explicit borders)
	leftPanel := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color(leftBorderColor)).
		Padding(1, 0, 0, 2). // Reduced top padding
		Width(leftWidth).
		Height(mainHeight - 1). // Adjust height for title
		Render(renderSongList(m))

	// Right panel
	rightPanel := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color(rightBorderColor)).
		Padding(1, 0, 0, 2).
		Width(leftWidth).
		Height(mainHeight - 1).
		Render(renderQueue(m))

	// Status bar (help text)
	statusBar := lipgloss.NewStyle().
		Width(mainWidth).
		Foreground(lipgloss.Color("#FFFFFF")).
		Align(lipgloss.Center).
		Render(helpView())

	// Combine panels horizontally
	panelView := lipgloss.JoinHorizontal(lipgloss.Left, leftPanel, rightPanel)

	// Final layout: Title → Panels → Status bar
	return lipgloss.JoinVertical(
		lipgloss.Left,
		title,     // Title at the top
		panelView, // Panels below title
		statusBar, // Status bar at bottom
	)
}

func helpView() string {
	helpKeys := []string{
		"j/k: Navigate",
		"Tab: Switch panels",
		"Enter: Play",
		"Space: Pause",
		"a: Add to queue",
		"d: Remove from queue",
		"h/l: Prev/Next in queue",
		"s: Stop",
		"q/ctrl+c: Quit",
	}

	var helpText []string
	for _, entry := range helpKeys {
		// Split keybinding and description
		parts := strings.SplitN(entry, ":", 2)
		if len(parts) == 2 {
			key := parts[0]  // e.g., "j/k"
			desc := parts[1] // e.g., " Navigate"

			// Apply bold to the key, normal style to description
			boldKey := lipgloss.NewStyle().Bold(true).Render(key)
			helpText = append(helpText, footerStyle.Render(boldKey+":"+lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#94A3B8")).Render(desc)))
		} else {
			// Fallback if no ":" is found
			helpText = append(helpText, footerStyle.Render(entry))
		}
	}
	return strings.Join(helpText, " | ")
}

func main() {
	p := tea.NewProgram(model{
		currentPlaying: -1,
		queueCursor:    0,
		activePanel:    0,
	}, tea.WithAltScreen())

	if err := p.Start(); err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
}
