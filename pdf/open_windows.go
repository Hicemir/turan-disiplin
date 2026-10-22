//go:build windows

package pdf

import (
	"fmt"
	"os/exec"
	"runtime"
	"syscall"
	"unsafe"
)

var (
	shell32            = syscall.NewLazyDLL("shell32.dll")
	procShellExecuteW  = shell32.NewProc("ShellExecuteW")
	ole32              = syscall.NewLazyDLL("ole32.dll")
	procCoInitializeEx = ole32.NewProc("CoInitializeEx")
	procCoUninitialize = ole32.NewProc("CoUninitialize")
)

const (
	coinitApartmentThreaded = 0x2
	coinitDisableOleDde     = 0x4
	swShowNormal            = 1
	createNoWindow          = 0x08000000
)

// shellExecute calls ShellExecuteW on a locked OS thread with COM initialised
// (required by many shell handlers, otherwise it fails on some machines).
func shellExecute(verb, file, params string) error {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	hr, _, _ := procCoInitializeEx.Call(0, coinitApartmentThreaded|coinitDisableOleDde)
	// S_OK (0) or S_FALSE (1) both require a matching CoUninitialize.
	if hr == 0 || hr == 1 {
		defer procCoUninitialize.Call()
	}

	var verbPtr, paramsPtr *uint16
	var err error
	if verb != "" {
		if verbPtr, err = syscall.UTF16PtrFromString(verb); err != nil {
			return err
		}
	}
	filePtr, err := syscall.UTF16PtrFromString(file)
	if err != nil {
		return err
	}
	if params != "" {
		if paramsPtr, err = syscall.UTF16PtrFromString(params); err != nil {
			return err
		}
	}
	ret, _, callErr := procShellExecuteW.Call(
		0,
		uintptr(unsafe.Pointer(verbPtr)),
		uintptr(unsafe.Pointer(filePtr)),
		uintptr(unsafe.Pointer(paramsPtr)),
		0,
		swShowNormal,
	)
	// ShellExecute returns a value > 32 on success.
	if ret <= 32 {
		return fmt.Errorf("ShellExecute hata kodu %d: %v", ret, callErr)
	}
	return nil
}

func runHidden(name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true, CreationFlags: createNoWindow}
	return cmd.Start()
}

// openWithSystem opens a file or folder with the default Windows handler.
// Several strategies are tried in order because individual machines differ
// (missing PDF association, locked-down shell, non-ASCII profile paths...).
func openWithSystem(path string) error {
	err1 := shellExecute("open", path, "")
	if err1 == nil {
		return nil
	}
	// Second attempt: no explicit verb (uses the default verb of the file type).
	if err := shellExecute("", path, ""); err == nil {
		return nil
	}
	// Third attempt: rundll32 FileProtocolHandler.
	if err := runHidden("rundll32.exe", "url.dll,FileProtocolHandler", path); err == nil {
		return nil
	}
	// Last resort: show the file in Explorer so the user can still reach it
	// (e.g. no PDF reader installed).
	if err := shellExecute("open", "explorer.exe", "/select,\""+path+"\""); err == nil {
		return fmt.Errorf("bilgisayarda PDF okuyucu bulunamadı veya açılamadı; dosya klasörde gösterildi (%v)", err1)
	}
	return err1
}
