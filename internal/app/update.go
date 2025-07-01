package app

import (
	tea "github.com/charmbracelet/bubbletea"
)

// Update
func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		m.Width = msg.Width
		m.Height = msg.Height
		return m, nil

	case FilesMsg:
		m.Files = msg
		m.Choices = msg
		m.Selected = make(map[int]struct{})
		return m, nil

	case ErrMsg:
		m.Err = msg
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {

		// Quit
		case "ctrl+c", "q":
			m.stopPlayback()
			return m, tea.Quit

		// Switch List (Queue / Song list)
		case "tab":
			m.ActivePanel = (m.ActivePanel + 1) % 2

		// Up
		case "up", "k":
			if m.ActivePanel == 0 && m.Cursor > 0 {
				m.Cursor--
			} else if m.ActivePanel == 1 && m.QueueCursor > 0 {
				m.QueueCursor--
			}

		// Down
		case "down", "j":
			if m.ActivePanel == 0 && m.Cursor < len(m.Choices)-1 {
				m.Cursor++
			} else if m.ActivePanel == 1 && m.QueueCursor < len(m.Queue)-1 {
				m.QueueCursor++
			}

		// Add song in queue
		case "a":
			if len(m.Files) > 0 {
				m.Queue = append(m.Queue, m.Files[m.Cursor])
			}

		// Remove song from queue
		case "d":
			m.removeFromQueue()

		// Start Playing song
		case "enter":
			if m.ActivePanel == 0 && len(m.Files) > 0 {
				m.PlaySong(m.Files[m.Cursor], -1)
			} else if m.ActivePanel == 1 && len(m.Queue) > 0 {
				m.PlaySong(m.Queue[m.QueueCursor], m.QueueCursor)
			}

		// Toggle Song (Play/ Pause)
		case " ":
			m.togglePause()

		// Play next song from queue
		case "l":
			if len(m.Queue) > 0 {
				next := m.CurrentPlaying + 1
				if next < len(m.Queue) {
					m.PlaySong(m.Queue[next], next)
					m.QueueCursor = next
				}
			}

		// Play previous song from queue
		case "h":
			if len(m.Queue) > 0 {
				prev := m.CurrentPlaying - 1
				if prev >= 0 {
					m.PlaySong(m.Queue[prev], prev)
					m.QueueCursor = prev
				}
			}

		// Stop playing Song
		case "s":
			m.stopPlayback()
		}
	}

	return m, nil
}
