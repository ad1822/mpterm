package app

import (
	"os/exec"
	"path/filepath"
	"syscall"
)

// For Playing Song when selected
func (m *Model) PlaySong(filename string, queueIndex int) {
	if m.ProcessPid != nil {
		_ = m.ProcessPid.Kill()
		_ = m.ProcessPid.Release()
		m.ProcessPid = nil
	}

	var fullPath string
	if filepath.IsAbs(filename) {
		fullPath = filename
	} else {
		fullPath = filepath.Join(m.CurrentPath, filename)
	}
	// path := filepath.Join(GetFullPath())
	cmd := exec.Command("mpv", "--input-ipc-server=/tmp/mpvsock", fullPath)
	if err := cmd.Start(); err != nil {
		return
	}

	m.CurrentSong = filename
	m.ProcessPid = cmd.Process
	m.IsPaused = false
	m.CurrentPlaying = queueIndex

	go func(currentIndex int, fromQueue bool) {
		err := cmd.Wait()
		if err != nil {
			return
		}

		m.ProcessPid = nil
		m.CurrentPlaying = -1

		if fromQueue && m.QueueCursor+1 < len(m.Queue) {
			m.QueueCursor++
			m.PlaySong(m.Queue[m.QueueCursor], m.QueueCursor)
		} else if !fromQueue {
			nextIndex := currentIndex + 1
			if nextIndex < len(m.Files) {
				m.Cursor = nextIndex
				m.PlaySong(m.Files[nextIndex], nextIndex)
			}
		}
	}(queueIndex, m.PlayingFromQueue)

}

// Stop playing song
func (m *Model) stopPlayback() {
	if m.ProcessPid != nil {
		_ = m.ProcessPid.Kill()
		_ = m.ProcessPid.Release()
		m.ProcessPid = nil
	}
	// m.CurrentPlaying = -1
	m.CurrentSong = ""
	m.IsPaused = false
}

// Toggle Song
func (m *Model) togglePause() {
	if m.ProcessPid == nil {
		return
	}
	var err error
	if m.IsPaused {
		err = m.ProcessPid.Signal(syscall.SIGCONT)
		if err == nil {
			m.IsPaused = false
		}
	} else {
		err = m.ProcessPid.Signal(syscall.SIGSTOP)
		if err == nil {
			m.IsPaused = true
		}
	}
}

// Remove Song from Queue
func (m *Model) removeFromQueue() {
	if len(m.Queue) == 0 || m.QueueCursor < 0 || m.QueueCursor >= len(m.Queue) {
		return
	}

	if m.CurrentPlaying >= 0 && m.QueueCursor == m.CurrentPlaying {
		m.stopPlayback()
	}

	m.Queue = append(m.Queue[:m.QueueCursor], m.Queue[m.QueueCursor+1:]...)

	if m.QueueCursor >= len(m.Queue) && len(m.Queue) > 0 {
		m.QueueCursor = len(m.Queue) - 1
	}
	if len(m.Queue) == 0 {
		m.QueueCursor = 0
	}

	if m.CurrentPlaying > m.QueueCursor {
		m.CurrentPlaying--
	}
}
