package style

import "github.com/charmbracelet/lipgloss"

var (
	TitleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#A78BFA")).
			Bold(true).
			MarginBottom(1)

	CursorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#f5e0dc")).
			Bold(true)

	QueueCursorStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#f2cdcd")).
				Bold(true)

	NormalStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#94A3B8"))

	PlayingStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#4ADE80")).
			Bold(true)

	PausedStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FACC15")).
			Bold(true)

	FooterStyle = lipgloss.NewStyle()

	Border = lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("228")).
		BorderBackground(lipgloss.Color("63")).
		BorderTop(true).
		BorderLeft(true)
)
