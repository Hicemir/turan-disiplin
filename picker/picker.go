package picker

import "errors"

// ErrCancelled indicates the user cancelled or closed the file dialog
var ErrCancelled = errors.New("dialog cancelled by user")

// ErrFallbackRequired indicates a native picker tool is not found, requesting GUI fallback
var ErrFallbackRequired = errors.New("native picker unavailable, fallback required")

// FileFilter defines a file extension filter for dialogs
type FileFilter struct {
	Name       string   // e.g. "Excel Dosyaları"
	Extensions []string // e.g. ["xlsx", "xls"]
}
