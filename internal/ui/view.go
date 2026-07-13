package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"

	"github.com/syouta-yamaguchi-rv/lazy-finder/internal/fsops"
)

// View implements tea.Model.
func (m Model) View() string {
	if m.width == 0 || m.height == 0 {
		return "読み込み中…"
	}

	header := m.renderHeader()
	bar := m.renderStatusBar()

	// Rows available for the three-column panel area.
	areaH := m.height - lipgloss.Height(header) - lipgloss.Height(bar)
	if areaH < 3 {
		areaH = 3
	}

	body := m.renderColumns(areaH)

	if m.mode == modeHelp {
		body = m.overlay(body, m.renderHelp(), areaH)
	}

	return lipgloss.JoinVertical(lipgloss.Left, header, body, bar)
}

func (m Model) renderHeader() string {
	path := m.cwd
	home := fsops.HomeDir()
	if strings.HasPrefix(path, home) {
		path = "~" + strings.TrimPrefix(path, home)
	}
	title := m.styles.titleActive.Render(" lazy-finder ")
	loc := m.styles.muted.Render(path)
	return m.styles.header.Width(m.width).MaxWidth(m.width).Render(title + " " + loc)
}

// renderColumns lays out parent | current | preview.
func (m Model) renderColumns(height int) string {
	// Inner widths (the +2 per panel accounts for left/right borders).
	parentW := clamp(m.width/6, 14, 26)
	previewW := (m.width - 6) * 4 / 10
	currentW := m.width - 6 - parentW - previewW
	if currentW < 10 {
		currentW = 10
	}

	parent := m.panel("親", m.renderList(m.parentEntries, m.parentCursor, false, false, parentW, height), false, parentW, height)
	current := m.panel(baseName(m.cwd), m.renderList(m.entries, m.cursor, true, true, currentW, height), true, currentW, height)
	preview := m.panel("プレビュー", m.renderPreview(previewW, height), false, previewW, height)

	return lipgloss.JoinHorizontal(lipgloss.Top, parent, current, preview)
}

// panel wraps content in a titled, rounded border box.
func (m Model) panel(title, content string, active bool, innerW, totalH int) string {
	box := m.styles.panelInactive
	titleStyle := m.styles.title
	if active {
		box = m.styles.panelActive
		titleStyle = m.styles.titleActive
	}
	innerH := totalH - 2 // borders
	if innerH < 1 {
		innerH = 1
	}

	titleLine := fitContent(titleStyle.Render(truncate(title, innerW)), innerW)
	bodyH := innerH - 1
	body := fitBlock(content, innerW, bodyH)

	return box.Width(innerW).Height(innerH).Render(titleLine + "\n" + body)
}

// renderList renders entries as fixed-width rows with selection highlighting.
func (m Model) renderList(entries []fsops.Entry, cursor int, focused, isCurrent bool, innerW, totalH int) string {
	capacity := totalH - 3 // borders + title
	if capacity < 1 {
		capacity = 1
	}
	if len(entries) == 0 {
		return m.styles.muted.Render("(空)")
	}

	start := windowStart(cursor, len(entries), capacity)
	end := min(start+capacity, len(entries))

	var b strings.Builder
	for i := start; i < end; i++ {
		e := entries[i]
		b.WriteString(m.renderRow(e, i == cursor, focused, isCurrent, innerW))
		if i < end-1 {
			b.WriteByte('\n')
		}
	}
	return b.String()
}

func (m Model) renderRow(e fsops.Entry, onCursor, focused, isCurrent bool, innerW int) string {
	// Pick a single background for the whole row so every segment sits on it.
	rowBG := colBG
	switch {
	case onCursor && focused:
		rowBG = colBGSel
	case onCursor: // selected directory in the parent column
		rowBG = colBGMark
	}

	nameFG := colFG
	name := e.Name
	switch {
	case e.IsLink:
		nameFG = colCyan
		name += " →"
	case e.IsDir:
		nameFG = colBlue
		name += "/"
	}

	markerFG, marker := colFG, " "
	if isCurrent && m.marked[e.Path] {
		markerFG, marker = colAccent, "●"
	}

	// Reserve room for the marker and a leading space.
	textW := innerW - 2
	if textW < 1 {
		textW = 1
	}
	seg := lipgloss.NewStyle().Background(rowBG)
	bold := e.IsDir || (onCursor && focused)
	row := seg.Foreground(markerFG).Bold(true).Render(marker) +
		seg.Render(" ") +
		seg.Foreground(nameFG).Bold(bold).Render(truncate(name, textW))

	return seg.Width(innerW).MaxWidth(innerW).Render(row)
}

// renderPreview shows the child listing for a directory, or file contents.
func (m Model) renderPreview(innerW, totalH int) string {
	e, ok := m.selected()
	if !ok {
		return m.styles.muted.Render("(選択なし)")
	}
	if e.IsDir {
		children, err := fsops.ReadDir(e.Path, m.showHidden)
		if err != nil {
			return m.styles.errorText.Render(err.Error())
		}
		return m.renderList(children, -1, false, false, innerW, totalH)
	}

	capacity := totalH - 3
	if capacity < 1 {
		capacity = 1
	}
	info := m.styles.muted.Render(fmt.Sprintf("%s  %s", fsops.HumanSize(e.Size), e.ModTime.Format("2006-01-02 15:04")))
	content := fsops.FilePreview(e.Path, capacity-2)

	var b strings.Builder
	b.WriteString(fitContent(info, innerW) + "\n")
	for _, line := range strings.Split(content, "\n") {
		b.WriteString(fitContent(truncate(line, innerW), innerW) + "\n")
	}
	return strings.TrimRight(b.String(), "\n")
}

func (m Model) renderStatusBar() string {
	var line string
	switch m.mode {
	case modeInput:
		prompt := "名前変更"
		if m.inputAction == actMkdir {
			prompt = "新規フォルダ"
		}
		line = m.styles.statusKey.Render(" "+prompt+" ") + " " + m.input.View()
	case modeConfirm:
		line = m.styles.statusKey.Render(" 確認 ") + " " + m.confirmPrompt + m.styles.muted.Render("  [y] はい  [n/esc] いいえ")
	default:
		if m.status != "" {
			if m.isError {
				line = m.styles.errorText.Render("✗ " + m.status)
			} else {
				line = m.styles.okText.Render("✓ " + m.status)
			}
		}
		if clip := m.clipboardLabel(); clip != "" {
			line = strings.TrimSpace(line + "  " + m.styles.muted.Render(clip))
		}
	}

	statusLine := m.styles.status.Width(m.width).MaxWidth(m.width).Render(line)
	helpLine := m.styles.help.Width(m.width).MaxWidth(m.width).Render(m.renderShortHelp())
	return statusLine + "\n" + helpLine
}

func (m Model) clipboardLabel() string {
	if len(m.clipboard) == 0 {
		return ""
	}
	verb := "copy"
	if m.clipboardOp == opCut {
		verb = "cut"
	}
	return fmt.Sprintf("[%s %d 件]", verb, len(m.clipboard))
}

func (m Model) renderShortHelp() string {
	var parts []string
	for _, b := range m.keys.shortHelp() {
		h := b.Help()
		parts = append(parts, fmt.Sprintf("%s:%s", h.Key, h.Desc))
	}
	return truncate(strings.Join(parts, "  "), m.width-1)
}

func (m Model) renderHelp() string {
	var b strings.Builder
	b.WriteString(m.styles.titleActive.Render("キーバインド") + "\n\n")
	for _, group := range m.keys.fullHelp() {
		for _, bind := range group {
			h := bind.Help()
			b.WriteString(fmt.Sprintf("  %-12s %s\n", h.Key, h.Desc))
		}
		b.WriteString("\n")
	}
	b.WriteString(m.styles.muted.Render("  ? / esc で閉じる"))
	return m.styles.panelActive.Padding(1, 2).Render(b.String())
}

// overlay centres box over the body area, filling the surround with the theme bg.
func (m Model) overlay(body, box string, areaH int) string {
	return lipgloss.Place(m.width, areaH, lipgloss.Center, lipgloss.Center, box,
		lipgloss.WithWhitespaceBackground(colBG))
}

// --- helpers -------------------------------------------------------------

func windowStart(cursor, count, capacity int) int {
	if count <= capacity || cursor < 0 {
		return 0
	}
	start := cursor - capacity/2
	if start < 0 {
		start = 0
	}
	if start > count-capacity {
		start = count - capacity
	}
	return start
}

// fitContent pads/truncates a line to exactly w columns, filling with the
// theme's light background and default foreground.
func fitContent(s string, w int) string {
	if w < 1 {
		w = 1
	}
	return lipgloss.NewStyle().Foreground(colFG).Background(colBG).Width(w).MaxWidth(w).Render(s)
}

// fitBlock forces content to exactly w×h, padding with blank (themed) lines.
func fitBlock(s string, w, h int) string {
	lines := strings.Split(s, "\n")
	for len(lines) < h {
		lines = append(lines, "")
	}
	if len(lines) > h {
		lines = lines[:h]
	}
	for i := range lines {
		lines[i] = fitContent(lines[i], w)
	}
	return strings.Join(lines, "\n")
}

func truncate(s string, w int) string {
	if w < 1 {
		return ""
	}
	if lipgloss.Width(s) <= w {
		return s
	}
	if w <= 1 {
		return "…"
	}
	return lipgloss.NewStyle().MaxWidth(w-1).Render(s) + "…"
}

func baseName(p string) string {
	if p == "/" {
		return "/"
	}
	for i := len(p) - 1; i >= 0; i-- {
		if p[i] == '/' {
			return p[i+1:]
		}
	}
	return p
}

func clamp(v, lo, hi int) int {
	if v < lo {
		return lo
	}
	if v > hi {
		return hi
	}
	return v
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
