package app

import (
	"os"
	"strings"

	"github.com/ad1822/musicplayer/internal/style"
	tea "github.com/charmbracelet/bubbletea"
)

// Read Path
func ReadFilesCmd(path string) tea.Cmd {
	return func() tea.Msg {
		entries, err := os.ReadDir(path)
		if err != nil {
			return ErrMsg{err}
		}

		var names []string
		for _, entry := range entries {
			names = append(names, entry.Name())
		}

		return FilesMsg(names)
	}
}

// Render Song from Path at initialization
func RenderSongList(m *Model) string {
	if len(m.Files) == 0 {
		return "No songs found"
	}

	var b strings.Builder
	for i, file := range m.Files {
		var line string

		if i == m.Cursor && m.ActivePanel == 0 {
			line = style.CursorStyle.Render("  " + file)
		} else {
			line = style.NormalStyle.Render("  " + file)
		}

		if file == m.CurrentSong {
			if m.IsPaused {
				line = style.PausedStyle.Render("⏸ " + file)
			} else {
				line = style.PlayingStyle.Render("▶ " + file)
			}
		}

		b.WriteString(line + "\n")
	}
	return b.String()
}
