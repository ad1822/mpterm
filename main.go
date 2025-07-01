package main

import (
	"fmt"
	"os"

	"github.com/ad1822/musicplayer/internal/app"
	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	p := tea.NewProgram(&app.Model{
		CurrentPlaying: -1,
		QueueCursor:    0,
		ActivePanel:    0,
	}, tea.WithAltScreen())

	if err := p.Start(); err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
}
