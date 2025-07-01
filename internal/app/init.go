package app

import (
	tea "github.com/charmbracelet/bubbletea"
	"os"
)

// Init Method for Initialization
func (m *Model) Init() tea.Cmd {
	// homeDir, _ := os.UserHomeDir()
	// fullPath := filepath.Join(homeDir + "Music")
	// log.Print(fullPath)
	cmd := exec("xdg-user-dir","MUSIC") // use xdg-user-dir for getting music dir
	output, err := cmd.Output()
	if err != nil{
		return "",err
	}
	musicDir := strings.TrimSpace(string(output)) + "/"
	return ReadFilesCmd(musicDir)
}
