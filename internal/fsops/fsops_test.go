package fsops

import (
	"os"
	"path/filepath"
	"testing"
)

func TestReadDirSortsFoldersFirst(t *testing.T) {
	dir := t.TempDir()
	mustWrite(t, filepath.Join(dir, "b.txt"), "x")
	mustWrite(t, filepath.Join(dir, "a.txt"), "x")
	if err := os.Mkdir(filepath.Join(dir, "zsub"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, ".hidden"), []byte("x"), 0o644); err != nil {
		t.Fatal(err)
	}

	entries, err := ReadDir(dir, false)
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) != 3 {
		t.Fatalf("expected 3 visible entries, got %d", len(entries))
	}
	if !entries[0].IsDir || entries[0].Name != "zsub" {
		t.Errorf("expected directory first, got %+v", entries[0])
	}
	if entries[1].Name != "a.txt" || entries[2].Name != "b.txt" {
		t.Errorf("files not sorted: %v %v", entries[1].Name, entries[2].Name)
	}

	all, _ := ReadDir(dir, true)
	if len(all) != 4 {
		t.Errorf("expected 4 entries with hidden, got %d", len(all))
	}
}

func TestCopyRenameMkdir(t *testing.T) {
	dir := t.TempDir()
	src := filepath.Join(dir, "src.txt")
	mustWrite(t, src, "hello")

	dst := filepath.Join(dir, "dest")
	if err := Mkdir(dir, "dest"); err != nil {
		t.Fatal(err)
	}
	if err := CopyPath(src, dst); err != nil {
		t.Fatal(err)
	}
	if !exists(filepath.Join(dst, "src.txt")) {
		t.Fatal("copy missing")
	}
	// Copy again -> numbered suffix.
	if err := CopyPath(src, dst); err != nil {
		t.Fatal(err)
	}
	if !exists(filepath.Join(dst, "src 2.txt")) {
		t.Fatal("expected numbered duplicate")
	}

	if err := Rename(src, "renamed.txt"); err != nil {
		t.Fatal(err)
	}
	if exists(src) || !exists(filepath.Join(dir, "renamed.txt")) {
		t.Fatal("rename failed")
	}
}

func TestMovePath(t *testing.T) {
	dir := t.TempDir()
	src := filepath.Join(dir, "m.txt")
	mustWrite(t, src, "data")
	sub := filepath.Join(dir, "sub")
	os.Mkdir(sub, 0o755)

	if err := MovePath(src, sub); err != nil {
		t.Fatal(err)
	}
	if exists(src) || !exists(filepath.Join(sub, "m.txt")) {
		t.Fatal("move failed")
	}
}

func TestFilePreview(t *testing.T) {
	dir := t.TempDir()
	txt := filepath.Join(dir, "a.txt")
	mustWrite(t, txt, "line1\nline2\nline3")
	if got := FilePreview(txt, 2); got == "" {
		t.Fatal("empty preview")
	}

	bin := filepath.Join(dir, "b.bin")
	os.WriteFile(bin, []byte{0x00, 0x01, 0x02}, 0o644)
	if got := FilePreview(bin, 10); got != "(バイナリファイル)" {
		t.Errorf("expected binary notice, got %q", got)
	}
}

func mustWrite(t *testing.T, path, content string) {
	t.Helper()
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}

func exists(p string) bool {
	_, err := os.Lstat(p)
	return err == nil
}
