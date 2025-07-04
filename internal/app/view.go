package app

import (
	"github.com/charmbracelet/lipgloss"
)

// View
func (m *Model) View() string {
	mainHeight := m.Height - (m.Height / 10)
	mainWidth := m.Width - (m.Height / 10)
	leftWidth := (mainWidth / 2)

	// Panel border colors (active panel highlight)
	leftBorderColor, rightBorderColor := "0", "0"
	if m.ActivePanel == 0 {
		leftBorderColor = "#cba6f7"
	} else {
		rightBorderColor = "#cba6f7"
	}

	// Left panel (with explicit borders)
	leftPanel := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color(leftBorderColor)).
		Padding(0, 0, 0, 1). // Reduced top padding
		Width(leftWidth).
		Height(mainHeight - 1). // Adjust height for title
		Render(RenderSongList(m))

	// Right panel
	rightPanel := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color(rightBorderColor)).
		Padding(0, 0, 0, 1).
		Width(leftWidth).
		Height(mainHeight - 1).
		Render(RenderQueue(m))

	// Status bar (help text)
	statusBar := lipgloss.NewStyle().
		Width(mainWidth).
		Foreground(lipgloss.Color("#FFFFFF")).
		Align(lipgloss.Center).
		Render(HelpView())

	// Combine panels horizontally
	panelView := lipgloss.JoinHorizontal(lipgloss.Left, leftPanel, rightPanel)

	// Final layout: Title → Panels → Status bar
	return lipgloss.JoinVertical(
		lipgloss.Left,
		// title,     // Title at the top
		panelView, // Panels below title
		statusBar, // Status bar at bottom
	)
}
