# Terminal Music Player

A fast, minimalist, and stylish terminal-based music player written in Go, powered by [Bubbletea](https://github.com/charmbracelet/bubbletea) and [Lipgloss](https://github.com/charmbracelet/lipgloss). Built for Unix-like systems using [PipeWire](https://pipewire.org/) (`pw-play`) for audio playback.

---

## ✨ Features

- Automatically loads songs from a predefined directory
- Play, pause, and stop audio directly from the terminal
- Add and manage a playback queue interactively
- Navigate between currently queued songs
- Stylish dual-pane interface with color highlights
- Mouse-free operation using Vim-style keybindings
- Real-time playback state tracking (playing/paused)
- PipeWire backend via `pw-play`

---

## Requirements

- **Go** 1.21+
- **PipeWire** installed with `pw-play` available in PATH
- Audio files in `~/Music` (or the configured directory)

---

## Installation

```bash
git clone https://github.com/ad1822/mpterm.git
cd mpterm
go build -o mpterm ./main.go
./mpterm
````

---

## ⌨️ Keybindings

| Key            | Action                           |
| -------------- | -------------------------------- |
| `j` / `k`      | Move up/down in the list         |
| `Tab`          | Switch between song list & queue |
| `Enter`        | Play selected song               |
| `Space`        | Pause/Resume current song        |
| `a`            | Add selected song to queue       |
| `d`            | Remove song from queue           |
| `h` / `l`      | Play previous/next in queue      |
| `s`            | Stop current song                |
| `q` / `Ctrl+C` | Quit the player                  |

---

## Project Structure

```
music-player/
├── main.go               # Entry point
├── internal/
│   ├── app/              # Model, update, init, play logic
│   └── style/            # Centralized Lipgloss styles
├── go.mod
└── README.md
```

---

## TODO

* [ ] UI animations and transitions
* [ ] Configurable music directory (`~/.config/music-player/config.yaml`)
* [ ] Volume control (`+`, `-`)
* [ ] Shuffle and repeat modes
* [ ] Search/filter songs
* [ ] Native MPV/Ffplay backend support
* [ ] Playlist persistence across sessions

---


## Contributing

PRs and feature ideas are welcome. Open an issue to discuss improvements or bugs.

```