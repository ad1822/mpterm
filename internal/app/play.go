package app

import (
	"os/exec"
	"syscall"
)



// For Playing Song when selected
func (m *Model) PlaySong(filename string, queueIndex int) {
	if m.ProcessPid != nil {
		_ = m.ProcessPid.Kill()
		_ = m.ProcessPid.Release()
		m.ProcessPid = nil
	}

	cmd := exec.Command("mpv","--input-ipc-server=/tmp/mpvsock", "/home/arcadian/Music/"+filename)

	if err := cmd.Start(); err != nil {
		return
	} else {
		m.CurrentSong = filename
		m.ProcessPid = cmd.Process
		m.IsPaused = false
		m.CurrentPlaying = queueIndex

		go func() {
			err := cmd.Wait()
			if err != nil {
				return
			}
			m.ProcessPid = nil
			m.CurrentPlaying = -1
		}()
	}
}

// Stop playing song
func (m *Model) stopPlayback() {
	if m.ProcessPid != nil {
		_ = m.ProcessPid.Kill()
		_ = m.ProcessPid.Release()
		m.ProcessPid = nil
	}
	m.CurrentPlaying = -1
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
