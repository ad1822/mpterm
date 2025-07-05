package app

import (
	"encoding/json"
	"log"
	"net"

	tea "github.com/charmbracelet/bubbletea"
)

func forwardSong(seconds int) error {
	conn, err := net.Dial("unix", "/tmp/mpvsock")
	if err != nil {
		return err
	}
	defer conn.Close()
	msg := map[string]interface{}{

		"command": []interface{}{"seek", seconds, "relative"},
	}

	data, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	_, err = conn.Write(append(data, '\n'))
	return err
}

// Update
func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		if msg.Width != m.Width || msg.Height != m.Height {
			m.Width = msg.Width
			m.Height = msg.Height
			log.Print(m.Width, " ", m.Height)
		}
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
				if m.Cursor < m.ScrollOffset {
					m.ScrollOffset--
				}
			} else if m.ActivePanel == 1 && m.QueueCursor > 0 {
				m.QueueCursor--
			}

		// Down
		case "down", "j":
			if m.ActivePanel == 0 && m.Cursor < len(m.Choices)-1 {
				m.Cursor++
				if m.Cursor >= m.ScrollOffset+m.Height-5 {
					m.ScrollOffset++
				}
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
				m.PlaySong(m.Files[m.Cursor], m.Cursor)
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

		// Right Arrow for rewind next 5 seconds
		case "right":
			_ = forwardSong(5)

		// Left Arrow for skip previous 5 seconds
		case "left":
			_ = forwardSong(-5)

		// Stop playing Song
		case "s":
			m.stopPlayback()
		}
	}

	return m, nil
}
