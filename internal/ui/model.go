package ui

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/syouta-yamaguchi-rv/lazy-finder/internal/fsops"
)

type mode int

const (
	modeNormal mode = iota
	modeInput
	modeConfirm
	modeHelp
)

type clipboardOp int

const (
	opNone clipboardOp = iota
	opCopy
	opCut
)

type inputAction int

const (
	actNone inputAction = iota
	actRename
	actMkdir
)

// Model is the root Bubble Tea model holding all TUI state.
type Model struct {
	cwd     string
	entries []fsops.Entry
	cursor  int
	offset  int // first visible row in the current column

	parentEntries []fsops.Entry
	parentCursor  int

	marked      map[string]bool
	clipboard   []string
	clipboardOp clipboardOp

	showHidden bool

	mode        mode
	input       textinput.Model
	inputAction inputAction
	inputTarget string // path being renamed (actRename)

	confirmPrompt string
	confirmAction func() (string, error)

	width, height int

	status  string
	isError bool

	keys   keyMap
	styles styles
	help   help.Model
}

// New builds a Model rooted at startDir (falling back to the home directory).
func New(startDir string) (Model, error) {
	abs, err := filepath.Abs(startDir)
	if err != nil {
		return Model{}, err
	}
	if info, err := os.Stat(abs); err != nil || !info.IsDir() {
		abs = fsops.HomeDir()
	}

	ti := textinput.New()
	ti.Prompt = "› "
	ti.CharLimit = 255

	m := Model{
		cwd:    abs,
		marked: map[string]bool{},
		keys:   defaultKeys(),
		styles: newStyles(),
		help:   help.New(),
		input:  ti,
	}
	if err := m.load(); err != nil {
		return Model{}, err
	}
	return m, nil
}

// Init implements tea.Model.
func (m Model) Init() tea.Cmd { return nil }

// load refreshes the current and parent directory listings, clamping the cursor.
func (m *Model) load() error {
	entries, err := fsops.ReadDir(m.cwd, m.showHidden)
	if err != nil {
		return err
	}
	m.entries = entries
	m.clampCursor()

	parent := filepath.Dir(m.cwd)
	if parent != m.cwd {
		if pe, err := fsops.ReadDir(parent, m.showHidden); err == nil {
			m.parentEntries = pe
			m.parentCursor = indexOf(pe, m.cwd)
		}
	} else {
		m.parentEntries = nil
		m.parentCursor = 0
	}
	return nil
}

func (m *Model) clampCursor() {
	if m.cursor >= len(m.entries) {
		m.cursor = len(m.entries) - 1
	}
	if m.cursor < 0 {
		m.cursor = 0
	}
}

func (m *Model) selected() (fsops.Entry, bool) {
	if m.cursor >= 0 && m.cursor < len(m.entries) {
		return m.entries[m.cursor], true
	}
	return fsops.Entry{}, false
}

// targets returns the marked entries, or the entry under the cursor when nothing
// is marked. This drives copy/cut/delete.
func (m *Model) targets() []string {
	if len(m.marked) > 0 {
		paths := make([]string, 0, len(m.marked))
		for _, e := range m.entries {
			if m.marked[e.Path] {
				paths = append(paths, e.Path)
			}
		}
		return paths
	}
	if e, ok := m.selected(); ok {
		return []string{e.Path}
	}
	return nil
}

func (m *Model) setStatus(s string) { m.status, m.isError = s, false }
func (m *Model) setError(s string)  { m.status, m.isError = s, true }
func (m *Model) setErr(err error)   { m.setError(err.Error()) }

// enterDir descends into dir and resets the cursor to the top.
func (m *Model) enterDir(dir string) {
	m.cwd = dir
	m.cursor, m.offset = 0, 0
	if err := m.load(); err != nil {
		m.setErr(err)
	}
}

// goParent ascends to the parent directory, keeping the cursor on the directory
// we just left.
func (m *Model) goParent() {
	parent := filepath.Dir(m.cwd)
	if parent == m.cwd {
		return
	}
	prev := m.cwd
	m.cwd = parent
	m.offset = 0
	if err := m.load(); err != nil {
		m.setErr(err)
		return
	}
	m.cursor = indexOf(m.entries, prev)
	m.clampCursor()
}

func indexOf(entries []fsops.Entry, path string) int {
	for i, e := range entries {
		if e.Path == path {
			return i
		}
	}
	return 0
}

// Update implements tea.Model.
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height
		m.help.Width = msg.Width
		return m, nil
	case tea.KeyMsg:
		switch m.mode {
		case modeInput:
			return m.updateInput(msg)
		case modeConfirm:
			return m.updateConfirm(msg)
		case modeHelp:
			return m.updateHelp(msg)
		default:
			return m.updateNormal(msg)
		}
	}
	return m, nil
}

func (m Model) updateNormal(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	k := m.keys
	switch {
	case key.Matches(msg, k.Quit):
		return m, tea.Quit
	case key.Matches(msg, k.Up):
		m.moveCursor(-1)
	case key.Matches(msg, k.Down):
		m.moveCursor(1)
	case key.Matches(msg, k.Top):
		m.cursor = 0
	case key.Matches(msg, k.Bottom):
		m.cursor = len(m.entries) - 1
		m.clampCursor()
	case key.Matches(msg, k.Enter):
		m.activate()
	case key.Matches(msg, k.Parent):
		m.goParent()
	case key.Matches(msg, k.Home):
		m.enterDir(fsops.HomeDir())
	case key.Matches(msg, k.Hidden):
		m.showHidden = !m.showHidden
		if err := m.load(); err != nil {
			m.setErr(err)
		}
		m.setStatus(fmt.Sprintf("隠しファイル: %v", m.showHidden))
	case key.Matches(msg, k.Refresh):
		if err := m.load(); err != nil {
			m.setErr(err)
		} else {
			m.setStatus("再読込しました")
		}
	case key.Matches(msg, k.Mark):
		m.toggleMark()
	case key.Matches(msg, k.Yank):
		m.yank(opCopy)
	case key.Matches(msg, k.Cut):
		m.yank(opCut)
	case key.Matches(msg, k.Paste):
		m.paste()
	case key.Matches(msg, k.Delete):
		m.askDelete()
	case key.Matches(msg, k.Rename):
		m.startRename()
	case key.Matches(msg, k.NewDir):
		m.startMkdir()
	case key.Matches(msg, k.VSCode):
		m.openVSCode()
	case key.Matches(msg, k.Open):
		m.openDefault()
	case key.Matches(msg, k.Help):
		m.mode = modeHelp
	}
	return m, nil
}

func (m *Model) moveCursor(delta int) {
	m.cursor += delta
	m.clampCursor()
}

// activate enters a directory or opens a file with the default application.
func (m *Model) activate() {
	e, ok := m.selected()
	if !ok {
		return
	}
	if e.IsDir {
		m.enterDir(e.Path)
		return
	}
	if err := fsops.OpenWithDefault(e.Path); err != nil {
		m.setErr(err)
	} else {
		m.setStatus("開きました: " + e.Name)
	}
}

func (m *Model) toggleMark() {
	if e, ok := m.selected(); ok {
		if m.marked[e.Path] {
			delete(m.marked, e.Path)
		} else {
			m.marked[e.Path] = true
		}
		m.moveCursor(1)
	}
}

func (m *Model) yank(op clipboardOp) {
	t := m.targets()
	if len(t) == 0 {
		return
	}
	m.clipboard, m.clipboardOp = t, op
	verb := "コピー"
	if op == opCut {
		verb = "切り取り"
	}
	m.setStatus(fmt.Sprintf("%d 件を%sしました", len(t), verb))
}

func (m *Model) paste() {
	if len(m.clipboard) == 0 {
		m.setError("クリップボードが空です")
		return
	}
	var done, failed int
	for _, src := range m.clipboard {
		var err error
		if m.clipboardOp == opCut {
			err = fsops.MovePath(src, m.cwd)
		} else {
			err = fsops.CopyPath(src, m.cwd)
		}
		if err != nil {
			failed++
		} else {
			done++
		}
	}
	if m.clipboardOp == opCut {
		m.clipboard, m.clipboardOp = nil, opNone
	}
	m.marked = map[string]bool{}
	if err := m.load(); err != nil {
		m.setErr(err)
		return
	}
	if failed > 0 {
		m.setError(fmt.Sprintf("%d 件成功, %d 件失敗", done, failed))
	} else {
		m.setStatus(fmt.Sprintf("%d 件を貼り付けました", done))
	}
}

func (m *Model) askDelete() {
	t := m.targets()
	if len(t) == 0 {
		return
	}
	m.confirmPrompt = fmt.Sprintf("%d 件をゴミ箱へ移動しますか?", len(t))
	m.confirmAction = func() (string, error) {
		var failed int
		for _, p := range t {
			if err := fsops.Trash(p); err != nil {
				failed++
			}
		}
		m.marked = map[string]bool{}
		if err := m.load(); err != nil {
			return "", err
		}
		if failed > 0 {
			return "", fmt.Errorf("%d 件の削除に失敗しました", failed)
		}
		return fmt.Sprintf("%d 件をゴミ箱へ移動しました", len(t)), nil
	}
	m.mode = modeConfirm
}

func (m *Model) startRename() {
	e, ok := m.selected()
	if !ok {
		return
	}
	m.inputAction = actRename
	m.inputTarget = e.Path
	m.input.Placeholder = "新しい名前"
	m.input.SetValue(e.Name)
	m.input.CursorEnd()
	m.input.Focus()
	m.mode = modeInput
}

func (m *Model) startMkdir() {
	m.inputAction = actMkdir
	m.input.Placeholder = "フォルダ名"
	m.input.SetValue("")
	m.input.Focus()
	m.mode = modeInput
}

func (m *Model) openVSCode() {
	target := m.cwd
	if e, ok := m.selected(); ok && e.IsDir {
		target = e.Path
	}
	if err := fsops.OpenInVSCode(target); err != nil {
		m.setErr(err)
	} else {
		m.setStatus("VSCode で開きました: " + filepath.Base(target))
	}
}

func (m *Model) openDefault() {
	if e, ok := m.selected(); ok {
		if err := fsops.OpenWithDefault(e.Path); err != nil {
			m.setErr(err)
		} else {
			m.setStatus("開きました: " + e.Name)
		}
	}
}

func (m Model) updateInput(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.Type {
	case tea.KeyEsc:
		m.mode = modeNormal
		m.input.Blur()
		return m, nil
	case tea.KeyEnter:
		val := m.input.Value()
		var err error
		switch m.inputAction {
		case actRename:
			err = fsops.Rename(m.inputTarget, val)
		case actMkdir:
			err = fsops.Mkdir(m.cwd, val)
		}
		m.mode = modeNormal
		m.input.Blur()
		if err != nil {
			m.setErr(err)
			return m, nil
		}
		if lerr := m.load(); lerr != nil {
			m.setErr(lerr)
			return m, nil
		}
		if m.inputAction == actRename {
			m.cursor = indexOf(m.entries, filepath.Join(filepath.Dir(m.inputTarget), val))
			m.setStatus("名前を変更しました")
		} else {
			m.cursor = indexOf(m.entries, filepath.Join(m.cwd, val))
			m.setStatus("フォルダを作成しました")
		}
		m.clampCursor()
		return m, nil
	}
	var cmd tea.Cmd
	m.input, cmd = m.input.Update(msg)
	return m, cmd
}

func (m Model) updateConfirm(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch {
	case key.Matches(msg, m.keys.Confirm):
		m.mode = modeNormal
		if m.confirmAction != nil {
			status, err := m.confirmAction()
			if err != nil {
				m.setErr(err)
			} else {
				m.setStatus(status)
			}
		}
		m.confirmAction = nil
	case key.Matches(msg, m.keys.Cancel):
		m.mode = modeNormal
		m.confirmAction = nil
		m.setStatus("取消しました")
	}
	return m, nil
}

func (m Model) updateHelp(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	if key.Matches(msg, m.keys.Help) || key.Matches(msg, m.keys.Cancel) || key.Matches(msg, m.keys.Quit) {
		m.mode = modeNormal
	}
	return m, nil
}
