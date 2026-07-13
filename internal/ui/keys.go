package ui

import "github.com/charmbracelet/bubbles/key"

// keyMap defines every binding. Bindings follow vim/lazygit conventions.
type keyMap struct {
	Up      key.Binding
	Down    key.Binding
	Top     key.Binding
	Bottom  key.Binding
	Enter   key.Binding
	Parent  key.Binding
	Home    key.Binding
	Mark    key.Binding
	Yank    key.Binding
	Cut     key.Binding
	Paste   key.Binding
	Delete  key.Binding
	Rename  key.Binding
	NewDir  key.Binding
	Open    key.Binding
	VSCode  key.Binding
	Hidden  key.Binding
	Refresh key.Binding
	Help    key.Binding
	Quit    key.Binding
	Confirm key.Binding
	Cancel  key.Binding
}

func defaultKeys() keyMap {
	return keyMap{
		Up:      key.NewBinding(key.WithKeys("up", "k"), key.WithHelp("↑/k", "上へ")),
		Down:    key.NewBinding(key.WithKeys("down", "j"), key.WithHelp("↓/j", "下へ")),
		Top:     key.NewBinding(key.WithKeys("g", "home"), key.WithHelp("g", "先頭")),
		Bottom:  key.NewBinding(key.WithKeys("G", "end"), key.WithHelp("G", "末尾")),
		Enter:   key.NewBinding(key.WithKeys("enter", "l", "right"), key.WithHelp("enter/l", "開く/入る")),
		Parent:  key.NewBinding(key.WithKeys("h", "left", "backspace"), key.WithHelp("h", "親へ")),
		Home:    key.NewBinding(key.WithKeys("~"), key.WithHelp("~", "ホーム")),
		Mark:    key.NewBinding(key.WithKeys(" "), key.WithHelp("space", "選択")),
		Yank:    key.NewBinding(key.WithKeys("y"), key.WithHelp("y", "コピー")),
		Cut:     key.NewBinding(key.WithKeys("x"), key.WithHelp("x", "切り取り")),
		Paste:   key.NewBinding(key.WithKeys("p"), key.WithHelp("p", "貼り付け")),
		Delete:  key.NewBinding(key.WithKeys("d"), key.WithHelp("d", "ゴミ箱へ")),
		Rename:  key.NewBinding(key.WithKeys("r"), key.WithHelp("r", "名前変更")),
		NewDir:  key.NewBinding(key.WithKeys("n"), key.WithHelp("n", "新規フォルダ")),
		Open:    key.NewBinding(key.WithKeys("o"), key.WithHelp("o", "標準アプリ")),
		VSCode:  key.NewBinding(key.WithKeys("e"), key.WithHelp("e", "VSCodeで開く")),
		Hidden:  key.NewBinding(key.WithKeys("."), key.WithHelp(".", "隠しファイル")),
		Refresh: key.NewBinding(key.WithKeys("R", "ctrl+r"), key.WithHelp("R", "再読込")),
		Help:    key.NewBinding(key.WithKeys("?"), key.WithHelp("?", "ヘルプ")),
		Quit:    key.NewBinding(key.WithKeys("q", "ctrl+c"), key.WithHelp("q", "終了")),
		Confirm: key.NewBinding(key.WithKeys("y", "enter"), key.WithHelp("y", "確定")),
		Cancel:  key.NewBinding(key.WithKeys("esc", "n"), key.WithHelp("esc", "取消")),
	}
}

// shortHelp is the one-line hint shown in the status bar (normal mode).
func (k keyMap) shortHelp() []key.Binding {
	return []key.Binding{k.Enter, k.Parent, k.Mark, k.Yank, k.Paste, k.Delete, k.Rename, k.VSCode, k.Help, k.Quit}
}

// fullHelp is the grouped help shown in the help overlay.
func (k keyMap) fullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Up, k.Down, k.Top, k.Bottom},
		{k.Enter, k.Parent, k.Home, k.Hidden, k.Refresh},
		{k.Mark, k.Yank, k.Cut, k.Paste},
		{k.Delete, k.Rename, k.NewDir},
		{k.Open, k.VSCode, k.Help, k.Quit},
	}
}
