package ui

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

// sized returns a model with a window size applied so View() renders fully.
func sized(t *testing.T, dir string) Model {
	t.Helper()
	m, err := New(dir)
	if err != nil {
		t.Fatal(err)
	}
	updated, _ := m.Update(tea.WindowSizeMsg{Width: 100, Height: 30})
	return updated.(Model)
}

func keyMsg(s string) tea.KeyMsg { return tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune(s)} }

func TestNavigationAndView(t *testing.T) {
	dir := t.TempDir()
	os.Mkdir(filepath.Join(dir, "child"), 0o755)
	os.WriteFile(filepath.Join(dir, "child", "f.txt"), []byte("hi"), 0o644)
	os.WriteFile(filepath.Join(dir, "top.txt"), []byte("hello world"), 0o644)

	m := sized(t, dir)

	// View must render without panicking and contain entry names.
	out := m.View()
	if !strings.Contains(out, "child") || !strings.Contains(out, "top.txt") {
		t.Fatalf("view missing entries:\n%s", out)
	}

	// Cursor starts on the directory "child"; enter it.
	upd, _ := m.Update(keyMsg("l"))
	m = upd.(Model)
	if filepath.Base(m.cwd) != "child" {
		t.Fatalf("expected to enter child, cwd=%s", m.cwd)
	}

	// Go back to parent; cursor should land on "child".
	upd, _ = m.Update(keyMsg("h"))
	m = upd.(Model)
	if filepath.Base(m.cwd) != filepath.Base(dir) {
		t.Fatalf("expected to return to parent, cwd=%s", m.cwd)
	}
	if sel, ok := m.selected(); !ok || sel.Name != "child" {
		t.Fatalf("cursor should be on child, got %+v", sel)
	}
}

func TestMkdirViaInput(t *testing.T) {
	dir := t.TempDir()
	m := sized(t, dir)

	upd, _ := m.Update(keyMsg("n")) // start new-folder input
	m = upd.(Model)
	if m.mode != modeInput {
		t.Fatal("expected input mode")
	}
	for _, r := range "newdir" {
		upd, _ = m.Update(keyMsg(string(r)))
		m = upd.(Model)
	}
	upd, _ = m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	m = upd.(Model)

	if _, err := os.Stat(filepath.Join(dir, "newdir")); err != nil {
		t.Fatalf("newdir not created: %v", err)
	}
	if m.mode != modeNormal {
		t.Fatal("should return to normal mode")
	}
}

func TestCopyPasteFlow(t *testing.T) {
	dir := t.TempDir()
	os.WriteFile(filepath.Join(dir, "a.txt"), []byte("x"), 0o644)
	os.Mkdir(filepath.Join(dir, "dst"), 0o755)
	m := sized(t, dir)

	// Cursor on "dst" (dir, sorts first). Move down to a.txt, then yank.
	upd, _ := m.Update(keyMsg("j"))
	m = upd.(Model)
	sel, _ := m.selected()
	if sel.Name != "a.txt" {
		t.Fatalf("expected a.txt under cursor, got %s", sel.Name)
	}
	upd, _ = m.Update(keyMsg("y")) // copy
	m = upd.(Model)
	if len(m.clipboard) != 1 {
		t.Fatal("clipboard should hold 1 item")
	}

	// Enter dst and paste.
	for filepath.Base(m.cwd) != "dst" {
		// move cursor up to dst then enter
		upd, _ = m.Update(keyMsg("k"))
		m = upd.(Model)
		upd, _ = m.Update(keyMsg("l"))
		m = upd.(Model)
	}
	upd, _ = m.Update(keyMsg("p"))
	m = upd.(Model)
	if _, err := os.Stat(filepath.Join(dir, "dst", "a.txt")); err != nil {
		t.Fatalf("paste failed: %v", err)
	}
}

func TestDeleteConfirm(t *testing.T) {
	dir := t.TempDir()
	target := filepath.Join(dir, "gone.txt")
	os.WriteFile(target, []byte("x"), 0o644)
	t.Setenv("HOME", dir) // trash goes under this HOME

	m := sized(t, dir)
	upd, _ := m.Update(keyMsg("d"))
	m = upd.(Model)
	if m.mode != modeConfirm {
		t.Fatal("expected confirm mode")
	}
	upd, _ = m.Update(keyMsg("y")) // confirm
	m = upd.(Model)
	if _, err := os.Stat(target); !os.IsNotExist(err) {
		t.Fatal("file should have been trashed")
	}
}
