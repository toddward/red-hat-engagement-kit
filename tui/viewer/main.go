package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/toddward/red-hat-engagement-kit/tui/viewer/ui"
)

func main() {
	// stdin/stdout are connected to bash core via pipes.
	// Open /dev/tty for terminal I/O so bubbletea does not
	// interfere with the bash communication channel.
	tty, err := os.OpenFile("/dev/tty", os.O_RDWR, 0)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error opening terminal: %v\n", err)
		os.Exit(1)
	}
	defer tty.Close()

	// bash reads from our stdout; we read from bash's stdout (our stdin).
	app := ui.NewApp(os.Stdout, os.Stdin)

	p := tea.NewProgram(app,
		tea.WithInput(tty),
		tea.WithOutput(tty),
	)
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
