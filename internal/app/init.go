package app

import (
	"log"

	tea "github.com/charmbracelet/bubbletea"
)

// Init Method for Initialization
func (m *Model) Init() tea.Cmd {
	m.CurrentPath = GetFullPath()
	log.Print(m.CurrentPath)
	return ReadFilesCmd(m.CurrentPath)
}
