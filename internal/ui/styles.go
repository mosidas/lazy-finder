package ui

import "github.com/charmbracelet/lipgloss"

// Ayu Light palette. The UI paints an explicit light background so it keeps the
// ayu look regardless of the terminal's own background colour.
const (
	hexBG     = "#FAFAFA" // editor background
	hexBGSel  = "#D6E5F3" // cursor row — light blue selection
	hexBGMark = "#E6E8EB" // parent selection / status chip
	hexFG     = "#5C6166" // primary text
	hexAccent = "#FA8D3E" // orange accent (focused border, marks)
	hexBlue   = "#399EE6" // directories, titles
	hexCyan   = "#55B4D4" // symlinks
	hexGreen  = "#86B300" // success
	hexRed    = "#E65050" // errors
	hexMuted  = "#8A9199" // comments / muted text
	hexBorder = "#C5C9CE" // inactive border
)

var (
	colBG     = lipgloss.Color(hexBG)
	colBGSel  = lipgloss.Color(hexBGSel)
	colBGMark = lipgloss.Color(hexBGMark)
	colFG     = lipgloss.Color(hexFG)
	colAccent = lipgloss.Color(hexAccent)
	colBlue   = lipgloss.Color(hexBlue)
	colCyan   = lipgloss.Color(hexCyan)
	colGreen  = lipgloss.Color(hexGreen)
	colRed    = lipgloss.Color(hexRed)
	colMuted  = lipgloss.Color(hexMuted)
	colBorder = lipgloss.Color(hexBorder)
)

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
