package cmd

import (
	"fmt"
	"os"

	"github.com/ad1822/mpterm/internal/app"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var (
	versionFlag     bool
	keybindingsFlag bool
)

var Version = "dev"

var rootCmd = &cobra.Command{
	Use:   "mpterm",
	Short: "TUI for music player",
	Run: func(cmd *cobra.Command, args []string) {
		if versionFlag {
			// fmt.Printf(" %-10s mpterm version: %s\n", Version)
			fmt.Println(color.New(color.FgHiGreen).Sprint("mpterm Version : ", Version))
			return
		}
		if keybindingsFlag {
			key := color.New(color.FgCyan, color.Bold).SprintFunc()
			desc := color.New(color.FgWhite).SprintFunc()

			fmt.Println()
			fmt.Println(color.New(color.FgMagenta).Sprint("Keybindings:"))
			fmt.Printf("  %-12s : %s\n", key("q / Ctrl+C"), desc("Quit the player"))
			fmt.Printf("  %-12s : %s\n", key("Tab"), desc("Switch between song list & queue"))
			fmt.Printf("  %-12s : %s\n", key("↑ / k"), desc("Move up"))
			fmt.Printf("  %-12s : %s\n", key("↓ / j"), desc("Move down"))
			fmt.Printf("  %-12s : %s\n", key("a"), desc("Add selected song to queue"))
			fmt.Printf("  %-12s : %s\n", key("d"), desc("Remove song from queue"))
			fmt.Printf("  %-12s : %s\n", key("Enter"), desc("Play selected song"))
			fmt.Printf("  %-12s : %s\n", key("Space"), desc("Pause/Resume current song"))
			fmt.Printf("  %-12s : %s\n", key("← / →"), desc("Seek backward/forward 5s"))
			fmt.Printf("  %-12s : %s\n", key("h / l"), desc("Play previous/next song"))
			fmt.Printf("  %-12s : %s\n", key("s"), desc("Stop current song"))

			return
		}
		app.InitSQLite(app.GetQueueDBPath())
		p := tea.NewProgram(&app.Model{
			CurrentPlaying: -1,
			QueueCursor:    0,
			ActivePanel:    0,
		}, tea.WithAltScreen())

		if err := p.Start(); err != nil {
			fmt.Println("Error:", err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.PersistentFlags().BoolVarP(&versionFlag, "version", "v", false, "Print Version")
	rootCmd.PersistentFlags().BoolVarP(&keybindingsFlag, "keybindings", "k", false, "Print Keybindings")

}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
