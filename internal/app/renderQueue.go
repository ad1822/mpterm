package app

import (
	"strings"

	"github.com/ad1822/musicplayer/internal/style"
)

// Render Queue's song
func RenderQueue(m *Model, maxHeight int) string {
	if len(m.Queue) == 0 {
		return "Queue is empty (press 'a' to add)"
	}

	lineShow := 0
	var b strings.Builder
	for i, entry := range m.Queue {
		if lineShow >= maxHeight {
			break
		}
		var line string

		if i == m.QueueCursor && m.ActivePanel == 1 {
			line = style.QueueCursorStyle.Render("  " + entry)
		} else {
			line = style.NormalStyle.Render("  " + entry)
		}

		if i == m.CurrentPlaying {
			if m.IsPaused {
				line = style.PausedStyle.Render("⏸ " + entry)
			} else {
				line = style.PlayingStyle.Render("▶ " + entry)
			}
		}

		b.WriteString(line + "\n")
		lineShow++
	}
	return b.String()
}
