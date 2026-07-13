// Package fsops provides filesystem operations for lazy-finder.
// All paths are absolute. Operations are designed to work on macOS (the
// primary target) while remaining runnable on Linux for development.
package fsops

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
	"time"
)

// Entry is a single item in a directory listing.
type Entry struct {
	Name    string
	Path    string
	IsDir   bool
	IsLink  bool
	Size    int64
	Mode    os.FileMode
	ModTime time.Time
}

// ReadDir lists the entries of dir. Hidden (dot) files are excluded unless
// showHidden is true. Directories sort before files, then by name (case
// insensitive), matching Finder's default folder-first behaviour.
func ReadDir(dir string, showHidden bool) ([]Entry, error) {
	des, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	entries := make([]Entry, 0, len(des))
	for _, de := range des {
		name := de.Name()
		if !showHidden && strings.HasPrefix(name, ".") {
			continue
		}
		full := filepath.Join(dir, name)

		info, err := de.Info()
		if err != nil {
			// The entry vanished or is unreadable; skip it rather than abort.
			continue
		}

		isLink := info.Mode()&os.ModeSymlink != 0
		isDir := de.IsDir()
		if isLink {
			// Resolve symlinks so directory links behave like directories.
			if target, err := os.Stat(full); err == nil {
				isDir = target.IsDir()
			}
		}

		entries = append(entries, Entry{
			Name:    name,
			Path:    full,
			IsDir:   isDir,
			IsLink:  isLink,
			Size:    info.Size(),
			Mode:    info.Mode(),
			ModTime: info.ModTime(),
		})
	}

	sort.SliceStable(entries, func(i, j int) bool {
		if entries[i].IsDir != entries[j].IsDir {
			return entries[i].IsDir
		}
		return strings.ToLower(entries[i].Name) < strings.ToLower(entries[j].Name)
	})

	return entries, nil
}

// Mkdir creates a new directory named name inside parent.
func Mkdir(parent, name string) error {
	name = strings.TrimSpace(name)
	if name == "" {
		return fmt.Errorf("名前が空です")
	}
	target := filepath.Join(parent, name)
	if _, err := os.Lstat(target); err == nil {
		return fmt.Errorf("%q は既に存在します", name)
	}
	return os.Mkdir(target, 0o755)
}

// Rename renames path to a sibling called newName.
func Rename(path, newName string) error {
	newName = strings.TrimSpace(newName)
	if newName == "" {
		return fmt.Errorf("名前が空です")
	}
	if strings.ContainsRune(newName, os.PathSeparator) {
		return fmt.Errorf("名前にパス区切り文字は使えません")
	}
	target := filepath.Join(filepath.Dir(path), newName)
	if target == path {
		return nil
	}
	if _, err := os.Lstat(target); err == nil {
		return fmt.Errorf("%q は既に存在します", newName)
	}
	return os.Rename(path, target)
}

// CopyPath copies src into the directory dstDir, preserving the base name.
// Directories are copied recursively. If a name collision occurs a numbered
// suffix is appended (e.g. "file 2.txt").
func CopyPath(src, dstDir string) error {
	dst := uniqueDest(dstDir, filepath.Base(src))
	return copyRecursive(src, dst)
}

// MovePath moves src into dstDir. It first attempts a rename and falls back to
// a copy + delete when the source and destination are on different volumes.
func MovePath(src, dstDir string) error {
	dst := uniqueDest(dstDir, filepath.Base(src))
	if filepath.Dir(src) == dstDir && filepath.Base(src) == filepath.Base(dst) {
		return nil // moving onto itself
	}
	if err := os.Rename(src, dst); err == nil {
		return nil
	}
	// Cross-device or other rename failure: copy then remove.
	if err := copyRecursive(src, dst); err != nil {
		return err
	}
	return os.RemoveAll(src)
}

// Trash moves path to the platform trash directory rather than deleting it
// permanently. On macOS this is ~/.Trash; elsewhere the freedesktop location.
func Trash(path string) error {
	dir, err := trashDir()
	if err != nil {
		return err
	}
	if err := os.MkdirAll(dir, 0o700); err != nil {
		return err
	}
	dst := uniqueDest(dir, filepath.Base(path))
	if err := os.Rename(path, dst); err == nil {
		return nil
	}
	if err := copyRecursive(path, dst); err != nil {
		return err
	}
	return os.RemoveAll(path)
}

// OpenInVSCode opens path with the VS Code CLI (`code`).
func OpenInVSCode(path string) error {
	return runDetached("code", path)
}

// OpenWithDefault opens path with the OS default application.
func OpenWithDefault(path string) error {
	switch runtime.GOOS {
	case "darwin":
		return runDetached("open", path)
	default:
		return runDetached("xdg-open", path)
	}
}

// HomeDir returns the user's home directory, or "/" if it cannot be resolved.
func HomeDir() string {
	if h, err := os.UserHomeDir(); err == nil {
		return h
	}
	return "/"
}

func trashDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	if runtime.GOOS == "darwin" {
		return filepath.Join(home, ".Trash"), nil
	}
	return filepath.Join(home, ".local", "share", "Trash", "files"), nil
}

// uniqueDest returns a non-existing path inside dir for the given base name,
// inserting a numbered suffix before the extension on collision.
func uniqueDest(dir, base string) string {
	candidate := filepath.Join(dir, base)
	if _, err := os.Lstat(candidate); os.IsNotExist(err) {
		return candidate
	}
	ext := filepath.Ext(base)
	stem := strings.TrimSuffix(base, ext)
	for i := 2; ; i++ {
		name := fmt.Sprintf("%s %d%s", stem, i, ext)
		candidate = filepath.Join(dir, name)
		if _, err := os.Lstat(candidate); os.IsNotExist(err) {
			return candidate
		}
	}
}

func copyRecursive(src, dst string) error {
	info, err := os.Lstat(src)
	if err != nil {
		return err
	}

	switch {
	case info.Mode()&os.ModeSymlink != 0:
		target, err := os.Readlink(src)
		if err != nil {
			return err
		}
		return os.Symlink(target, dst)
	case info.IsDir():
		if err := os.MkdirAll(dst, info.Mode().Perm()); err != nil {
			return err
		}
		des, err := os.ReadDir(src)
		if err != nil {
			return err
		}
		for _, de := range des {
			if err := copyRecursive(filepath.Join(src, de.Name()), filepath.Join(dst, de.Name())); err != nil {
				return err
			}
		}
		return nil
	default:
		return copyFile(src, dst, info.Mode().Perm())
	}
}

func copyFile(src, dst string, perm os.FileMode) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.OpenFile(dst, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, perm)
	if err != nil {
		return err
	}
	if _, err := io.Copy(out, in); err != nil {
		out.Close()
		return err
	}
	return out.Close()
}
