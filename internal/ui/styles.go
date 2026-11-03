package ui

import "github.com/charmbracelet/lipgloss"

var (
	colBG     lipgloss.Color
	colBGSel  lipgloss.Color
	colBGMark lipgloss.Color
	colFG     lipgloss.Color
	colAccent lipgloss.Color
	colBlue   lipgloss.Color
	colCyan   lipgloss.Color
	colGreen  lipgloss.Color
	colRed    lipgloss.Color
	colMuted  lipgloss.Color
	colBorder lipgloss.Color
)

// applyTheme sets the palette that newStyles and the view read.
func applyTheme(t Theme) {
	colBG = lipgloss.Color(t.BG)
	colBGSel = lipgloss.Color(t.BGSel)
	colBGMark = lipgloss.Color(t.BGMark)
	colFG = lipgloss.Color(t.FG)
	colAccent = lipgloss.Color(t.Accent)
	colBlue = lipgloss.Color(t.Blue)
	colCyan = lipgloss.Color(t.Cyan)
	colGreen = lipgloss.Color(t.Green)
	colRed = lipgloss.Color(t.Red)
	colMuted = lipgloss.Color(t.Muted)
	colBorder = lipgloss.Color(t.Border)
}

type styles struct {
	panelActive   lipgloss.Style
	panelInactive lipgloss.Style
	title         lipgloss.Style
	titleActive   lipgloss.Style
	muted         lipgloss.Style
	header        lipgloss.Style
	status        lipgloss.Style
	statusKey     lipgloss.Style
	help          lipgloss.Style
	errorText     lipgloss.Style
	okText        lipgloss.Style
}

func newStyles() styles {
	box := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		Background(colBG).
		Foreground(colFG).
		BorderBackground(colBG)

	onBG := lipgloss.NewStyle().Background(colBG)

	return styles{
		panelActive:   box.BorderForeground(colAccent),
		panelInactive: box.BorderForeground(colBorder),
		title:         onBG.Foreground(colBlue).Bold(true),
		titleActive:   onBG.Foreground(colAccent).Bold(true),
		muted:         onBG.Foreground(colMuted),
		header:        onBG.Foreground(colFG),
		status:        onBG.Foreground(colFG),
		statusKey:     lipgloss.NewStyle().Background(colBGMark).Foreground(colAccent).Bold(true),
		help:          onBG.Foreground(colMuted),
		errorText:     onBG.Foreground(colRed).Bold(true),
		okText:        onBG.Foreground(colGreen),
	}
}
