package app

import (
	"github.com/charmbracelet/lipgloss"
)

// View
func (m *Model) View() string {
	mainHeight := m.Height - (m.Height / 10)
	mainWidth := m.Width - (m.Height / 10)
	// leftWidth := (mainWidth / 2)
	rightWidth := mainWidth / 3

	contentHeight := m.Height - 1
	maxVisibleLines := contentHeight - 3
	leftBorderColor, rightBorderColor := "0", "0"
	if m.ActivePanel == 0 {
		leftBorderColor = "#cba6f7"
	} else {
		rightBorderColor = "#cba6f7"
	}

	leftPanel := lipgloss.NewStyle().
		Border(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color(leftBorderColor)).
		Padding(0, 0, 0, 1).
		Width(mainWidth - rightWidth).
		Height(mainHeight - 1).
		Render(RenderSongList(m, maxVisibleLines))

	rightPanel := lipgloss.NewStyle().
		Border(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color(rightBorderColor)).
		Padding(0, 0, 0, 1).
		Width(rightWidth).
		Height(mainHeight - 1).
		Render(RenderQueue(m, maxVisibleLines))

	// statusBar := lipgloss.NewStyle().
	// 	Width(mainWidth).
	// 	Foreground(lipgloss.Color("#FFFFFF")).
	// 	Align(lipgloss.Center).
	// 	Render(HelpView())

	panelView := lipgloss.JoinHorizontal(lipgloss.Left, leftPanel, rightPanel)

	return lipgloss.JoinVertical(
		lipgloss.Left,
		panelView,
		// statusBar,
	)
}

func GetMaxVisibleLines(totalHeight int, paddingTopBottom int, borderTopBottom int) int {
	return totalHeight - (paddingTopBottom + borderTopBottom)
}
