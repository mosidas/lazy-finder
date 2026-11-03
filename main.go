// Command lazy-finder is a lazygit-style terminal file manager modelled on the
// macOS Finder column view.
package main

import (
	"flag"
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/mosidas/lazy-finder/internal/ui"
)

func main() {
	flag.Parse()

	start := "."
	if flag.NArg() > 0 {
		start = flag.Arg(0)
	}

	theme, err := ui.LookupTheme(ui.DefaultThemeName)
	if err != nil {
		fmt.Fprintln(os.Stderr, "lazy-finder:", err)
		os.Exit(1)
	}

	model, err := ui.New(start, theme)
	if err != nil {
		fmt.Fprintln(os.Stderr, "lazy-finder:", err)
		os.Exit(1)
	}

	p := tea.NewProgram(model, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Fprintln(os.Stderr, "lazy-finder:", err)
		os.Exit(1)
	}
}
