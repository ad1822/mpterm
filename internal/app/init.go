package app

import (
	"log"

	tea "github.com/charmbracelet/bubbletea"
)

// Init Method for Initialization
func (m *Model) Init() tea.Cmd {
	path := GetFullPath()
	log.Print(path)
	return ReadFilesCmd(path)
}
