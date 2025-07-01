package app

import (
	tea "github.com/charmbracelet/bubbletea"
)

// Init Method for Initialization
func (m *Model) Init() tea.Cmd {
	// homeDir, _ := os.UserHomeDir()
	// fullPath := filepath.Join(homeDir + "Music")
	// log.Print(fullPath)
	return ReadFilesCmd("/home/arcadian/Music")
}
