package app

import (
	"fmt"
	"log"
	"path/filepath"
	"strings"

	"github.com/ad1822/mpterm/internal/style"
)

// Render Queue's song
func RenderQueue(m *Model, maxHeight int) string {

	songs, err := m.GetQueueSongs()
	if err != nil {
		return "err"
	}

	if len(songs) == 0 {
		return "Queue is empty (press 'a' to add)"
	}

	m.Queue = songs

	lineShow := 0
	var b strings.Builder
	for i, entry := range m.Queue {
		if lineShow >= maxHeight {
			break
		}
		var line string
		// Show only the base name of the file
		base := filepath.Base(entry)

		if i == m.QueueCursor && m.ActivePanel == 1 {
			line = style.QueueCursorStyle.Render("  " + base)

		} else {
			line = style.NormalStyle.Render("  " + base)
		}

		if i == m.CurrentPlaying && m.ActivePanel == 1 && m.PlayingFromQueue {
			if m.IsPaused {
				line = style.PausedStyle.Render("⏸ " + base)
			} else {
				line = style.PlayingStyle.Render("▶ " + base)
			}
		}

		b.WriteString(line + "\n")
		lineShow++
	}
	return b.String()
}

// Add Songs in DB
func (m *Model) AddSongInQueue(file string) error {
	var fullPath string

	// Check if the file is absolute or relative
	if filepath.IsAbs(file) {
		fullPath = file
	} else {
		fullPath = filepath.Join(m.CurrentPath, file)
	}

	stmt, err := DB.Prepare("INSERT INTO queues(song_name) VALUES(?)")
	if err != nil {
		return fmt.Errorf("failed to prepare statement: %v", err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(fullPath)
	if err != nil {
		return fmt.Errorf("failed to insert song: %v", err)
	}

	return nil
}

func (m *Model) DeleteSongFromQueue(file string) error {
	stmt, err := DB.Prepare("DELETE FROM queues WHERE song_name = ?")
	if err != nil {
		return fmt.Errorf("failed to prepare statement: %v", err)
	}

	defer stmt.Close()
	_, err = stmt.Exec(file)
	if err != nil {
		return fmt.Errorf("failed to insert song: %v", err)
	}

	return nil
}

// Get Songs from DB
func (m *Model) GetQueueSongs() ([]string, error) {
	const query = `SELECT song_name FROM queues ORDER BY id`

	rows, err := DB.Query(query)
	if err != nil {
		return nil, fmt.Errorf("database query error: %w", err)
	}
	defer rows.Close()

	var songs []string
	for rows.Next() {
		var songName string
		if err := rows.Scan(&songName); err != nil {
			log.Printf("Error scanning row: %v", err)
			continue
		}
		songs = append(songs, songName)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}

	return songs, nil
}
