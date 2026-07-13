package fsops

import (
	"bytes"
	"fmt"
	"os"
	"strings"
	"unicode/utf8"
)

const previewByteLimit = 64 * 1024 // read at most 64KiB for previewing

// FilePreview returns a human-readable preview of a regular file: its text
// content when it looks like text, otherwise a short "binary file" notice.
// maxLines caps the number of returned lines so callers can avoid reading huge
// files into the view.
func FilePreview(path string, maxLines int) string {
	f, err := os.Open(path)
	if err != nil {
		return "プレビューを開けません: " + err.Error()
	}
	defer f.Close()

	buf := make([]byte, previewByteLimit)
	n, _ := f.Read(buf)
	buf = buf[:n]

	if n == 0 {
		return "(空のファイル)"
	}
	if isBinary(buf) {
		return "(バイナリファイル)"
	}

	text := strings.ReplaceAll(string(buf), "\r\n", "\n")
	lines := strings.Split(text, "\n")
	if maxLines > 0 && len(lines) > maxLines {
		lines = lines[:maxLines]
		lines = append(lines, fmt.Sprintf("… (%d 行以降は省略)", maxLines))
	}
	return strings.Join(lines, "\n")
}

// isBinary reports whether buf appears to be non-text data. A NUL byte or
// invalid UTF-8 is treated as binary.
func isBinary(buf []byte) bool {
	if bytes.IndexByte(buf, 0) != -1 {
		return true
	}
	if !utf8.Valid(buf) {
		// Trailing bytes may be a truncated rune; tolerate the final few.
		trimmed := buf
		for i := 0; i < 3 && len(trimmed) > 0 && !utf8.Valid(trimmed); i++ {
			trimmed = trimmed[:len(trimmed)-1]
		}
		return !utf8.Valid(trimmed)
	}
	return false
}

// HumanSize formats a byte count as a compact human readable string.
func HumanSize(n int64) string {
	const unit = 1024
	if n < unit {
		return fmt.Sprintf("%dB", n)
	}
	div, exp := int64(unit), 0
	for v := n / unit; v >= unit; v /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f%cB", float64(n)/float64(div), "KMGTPE"[exp])
}
