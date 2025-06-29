package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"

	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
	files    []string
	err      error
	cursor   int
	choices  []string
	selected map[int]struct{}
}
type filesMsg []string
type errMsg struct{ error }

func (m model) Init() tea.Cmd {
	return readFilesCmd("/home/arcadian/Music")
	// return model{
	// 	choices: []string{m.files[]...},

	// 	selected: make(map[int]struct{}),
	// }
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
		case "enter", " ":
			if _, ok := m.selected[m.cursor]; ok {
				file := m.files[m.cursor]
				cmd := exec.Command("pw-play", "/home/arcadian/Music/"+file)
				if err := cmd.Run(); err != nil {
					log.Printf("Error playing file: %v", err)
				}
				delete(m.selected, m.cursor)
			} else {
				m.selected[m.cursor] = struct{}{}
			}
		}
	}

	return m, nil
}

func (m model) View() string {
	if m.err != nil {
		return fmt.Sprintf("Error: %v\n", m.err)
	}

	if len(m.files) == 0 {
		return "Loading files...\n"
	}

	s := "Files in current directory:\n\n"
	for i, file := range m.files {
		cursor := " "
		if m.cursor == i {
			cursor = ">"
		}

		checked := " "
		if _, ok := m.selected[i]; ok {
			checked = "x"
		}

		s += fmt.Sprintf("%s [%s] %s\n", cursor, checked, file)
	}

	return s
}

func main() {
	p := tea.NewProgram(model{})
	if err := p.Start(); err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
}
