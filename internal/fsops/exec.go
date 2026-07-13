package fsops

import (
	"fmt"
	"os/exec"
)

// runDetached starts cmd with args without waiting for it to finish. It is used
// to hand a path off to an external GUI application (VS Code, the system opener)
// and return control to the TUI immediately.
func runDetached(name string, args ...string) error {
	if _, err := exec.LookPath(name); err != nil {
		return fmt.Errorf("%s が見つかりません (PATH を確認してください)", name)
	}
	cmd := exec.Command(name, args...)
	return cmd.Start()
}
