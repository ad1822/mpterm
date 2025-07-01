package app

import (
	"strings"

	"github.com/ad1822/musicplayer/internal/style"
	"github.com/charmbracelet/lipgloss"
)

// Help Panel in bottom that displays keybinding
func HelpView() string {
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
		parts := strings.SplitN(entry, ":", 2)
		if len(parts) == 2 {
			key := parts[0]  // e.g., "j/k"
			desc := parts[1] // e.g., " Navigate"

			boldKey := lipgloss.NewStyle().Bold(true).Render(key)
			helpText = append(helpText, style.FooterStyle.Render(boldKey+":"+lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#94A3B8")).Render(desc)))
		} else {
			helpText = append(helpText, style.FooterStyle.Render(entry))
		}
	}
	return strings.Join(helpText, " | ")
}
