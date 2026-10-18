//go:build windows

package picker

import (
	"fmt"
	"strings"
	"syscall"
	"unicode/utf16"
	"unsafe"

	"github.com/TheTitanrain/w32"
)

func utf16ptr(s []uint16) *uint16 {
	if len(s) == 0 {
		return nil
	}
	return unsafe.SliceData(s)
}

func decodeFilename(buf []uint16) string {
	i := 0
	for i < len(buf) && buf[i] != 0 {
		i++
	}
	return string(utf16.Decode(buf[:i]))
}

func buildFilter(filters []FileFilter) []uint16 {
	var f []uint16
	for _, filt := range filters {
		f = append(f, utf16.Encode([]rune(filt.Name))...)
		f = append(f, 0)
		for _, ext := range filt.Extensions {
			clean := strings.TrimPrefix(ext, ".")
			s := fmt.Sprintf("*.%s;", clean)
			f = append(f, utf16.Encode([]rune(s))...)
		}
		f = append(f, 0)
	}
	if len(f) > 0 {
		f = append(f, 0, 0)
	}
	return f
}

// OpenFile displays the native Windows file open dialog using Win32 API
func OpenFile(title string, filters []FileFilter) (string, error) {
	buf := make([]uint16, w32.MAX_PATH)
	filterBytes := buildFilter(filters)

	opf := &w32.OPENFILENAME{
		File:    utf16ptr(buf),
		MaxFile: uint32(len(buf)),
		Flags:   w32.OFN_FILEMUSTEXIST | w32.OFN_NOCHANGEDIR,
	}
	opf.StructSize = uint32(unsafe.Sizeof(*opf))
	if title != "" {
		opf.Title, _ = syscall.UTF16PtrFromString(title)
	}
	if len(filterBytes) > 0 {
		opf.Filter = utf16ptr(filterBytes)
	}

	if w32.GetOpenFileName(opf) {
		res := decodeFilename(buf)
		if res != "" {
			return res, nil
		}
	}

	errCode := w32.CommDlgExtendedError()
	if errCode == 0 {
		return "", ErrCancelled
	}
	return "", fmt.Errorf("Windows OpenFile error: %#x", int(errCode))
}

// SaveFile displays the native Windows file save dialog using Win32 API
func SaveFile(title string, defaultFilename string, filters []FileFilter) (string, error) {
	buf := make([]uint16, w32.MAX_PATH)
	if defaultFilename != "" {
		initialName, _ := syscall.UTF16FromString(defaultFilename)
		for i := 0; i < len(initialName) && i < w32.MAX_PATH; i++ {
			buf[i] = initialName[i]
		}
	}
	filterBytes := buildFilter(filters)

	opf := &w32.OPENFILENAME{
		File:    utf16ptr(buf),
		MaxFile: uint32(len(buf)),
		Flags:   w32.OFN_OVERWRITEPROMPT | w32.OFN_NOCHANGEDIR,
	}
	opf.StructSize = uint32(unsafe.Sizeof(*opf))
	if title != "" {
		opf.Title, _ = syscall.UTF16PtrFromString(title)
	}
	if len(filterBytes) > 0 {
		opf.Filter = utf16ptr(filterBytes)
	}

	if w32.GetSaveFileName(opf) {
		res := decodeFilename(buf)
		if res != "" {
			return res, nil
		}
	}

	errCode := w32.CommDlgExtendedError()
	if errCode == 0 {
		return "", ErrCancelled
	}
	return "", fmt.Errorf("Windows SaveFile error: %#x", int(errCode))
}
