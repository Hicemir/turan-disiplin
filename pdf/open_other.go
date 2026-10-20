//go:build !windows

package pdf

import (
	"os/exec"
	"runtime"
)

// openWithSystem opens a file or folder with the default system application.
func openWithSystem(path string) error {
	var cmd *exec.Cmd
	if runtime.GOOS == "darwin" {
		cmd = exec.Command("open", path)
	} else {
		cmd = exec.Command("xdg-open", path)
	}
	if err := cmd.Start(); err != nil {
		return err
	}
	go func() { _ = cmd.Wait() }() // reap child, avoid zombies
	return nil
}
