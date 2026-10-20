//go:build windows

package pdf

import (
	"fmt"
	"syscall"
	"unsafe"
)

var (
	shell32           = syscall.NewLazyDLL("shell32.dll")
	procShellExecuteW = shell32.NewProc("ShellExecuteW")
)

// openWithSystem opens a file or folder with the default Windows handler via
// ShellExecuteW. This works reliably from a GUI (-H=windowsgui) process and
// with non-ASCII paths, unlike spawning rundll32 as a child process.
func openWithSystem(path string) error {
	verb, err := syscall.UTF16PtrFromString("open")
	if err != nil {
		return err
	}
	file, err := syscall.UTF16PtrFromString(path)
	if err != nil {
		return err
	}
	const swShowNormal = 1
	ret, _, callErr := procShellExecuteW.Call(
		0,
		uintptr(unsafe.Pointer(verb)),
		uintptr(unsafe.Pointer(file)),
		0,
		0,
		swShowNormal,
	)
	// ShellExecute returns a value > 32 on success.
	if ret <= 32 {
		return fmt.Errorf("ShellExecute hata kodu %d: %v", ret, callErr)
	}
	return nil
}
