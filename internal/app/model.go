package app

import "os"

// Model represents the application state
type Model struct {
	Files          []string // All music files
	Err            error
	Cursor         int              // Current cursor position in the song list
	Choices        []string         // Songs currently visible/selectable in the main list
	Selected       map[int]struct{} // Placeholder for multi-select logic
	ProcessPid     *os.Process      // Pointer to the currently playing process
	IsPaused       bool             // Indicates if the current song is paused or not
	CurrentPlaying int              // Index of the currently playing song in the queue (-1 if none)
	Width          int              // Terminal width (from tea.WindowSizeMsg)
	Height         int              // Terminal height (from tea.WindowSizeMsg)
	CurrentSong    string           // Name of the currently playing song file
	Queue          []string         // User-managed playback queue
	QueueCursor    int              // Current cursor position within the queue panel
	ActivePanel    int              // 0: song list panel active, 1: queue panel active
	ScrollOffset   int
}

// FilesMsg is a message carrying a list of discovered files.
type FilesMsg []string

// ErrMsg wraps an error into a tea.Msg for Bubbletea's message system.
type ErrMsg struct{ error }
