package app

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// Left panel: Songs
func LeftPanel(m *Model, leftWidth int, mainHeight int, leftColor string, maxVisibleLines int) string {
	leftTitle := " Songs "
	leftTitlePrefix := "╭─" + leftTitle
	leftDashes := leftWidth - lipgloss.Width(leftTitlePrefix) - 1
	if leftDashes < 0 {
		leftDashes = 0
	}
	leftTopLine := leftTitlePrefix + strings.Repeat("─", leftDashes) + "╮"
	styledLeftTop := lipgloss.NewStyle().
		Width(leftWidth).
		Foreground(lipgloss.Color(leftColor)).
		Render(leftTopLine)

	leftPanel := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderTop(false).
		Padding(1, 0, 0, 1).
		BorderForeground(lipgloss.Color(leftColor)).
		Width(leftWidth - 2).
		Height(mainHeight - 1).
		Render(RenderSongList(m, maxVisibleLines))

	return lipgloss.JoinVertical(lipgloss.Left, styledLeftTop, leftPanel)
}

func RightPanel(m *Model, rightWidth int, mainHeight int, rightColor string, maxVisibleLines int) string {
	// Right panel: Queue
	rightTitle := " Queue "
	rightTitlePrefix := "╭─" + rightTitle
	rightDashes := rightWidth - lipgloss.Width(rightTitlePrefix) - 1
	if rightDashes < 0 {
		rightDashes = 0
	}
	rightTopLine := rightTitlePrefix + strings.Repeat("─", rightDashes) + "╮"
	styledRightTop := lipgloss.NewStyle().
		Width(rightWidth).
		Foreground(lipgloss.Color(rightColor)).
		Render(rightTopLine)

	rightPanel := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderTop(false).
		Padding(1, 0, 0, 1).
		BorderForeground(lipgloss.Color(rightColor)).
		Width(rightWidth - 2).
		Height(mainHeight - 1).
		Render(RenderQueue(m, maxVisibleLines))

	return lipgloss.JoinVertical(lipgloss.Left, styledRightTop, rightPanel)

}

// View
func (m *Model) View() string {
	mainHeight := m.Height - (m.Height / 20)
	mainWidth := m.Width - (m.Height / 20)

	rightPanelRatio := 0.4
	rightWidth := int(float64(mainWidth) * rightPanelRatio)
	leftWidth := mainWidth - rightWidth

	contentHeight := mainHeight - 1
	maxVisibleLines := contentHeight - 3

	leftColor, rightColor := "0", "0"
	if m.ActivePanel == 0 {
		leftColor = "#cba6f7"
	} else {
		rightColor = "#cba6f7"
	}

	// Combine panels
	panelLayout := lipgloss.JoinHorizontal(lipgloss.Left, LeftPanel(m, leftWidth, mainHeight, leftColor, maxVisibleLines), RightPanel(m, rightWidth, mainHeight, rightColor, maxVisibleLines))

	return lipgloss.JoinVertical(
		lipgloss.Left,
		panelLayout,
	)
}

func GetMaxVisibleLines(totalHeight int, paddingTopBottom int, borderTopBottom int) int {
	return totalHeight - (paddingTopBottom + borderTopBottom)
}
