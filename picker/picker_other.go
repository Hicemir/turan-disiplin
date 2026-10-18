//go:build !linux && !windows

package picker

// OpenFile fallback for other platforms
func OpenFile(title string, filters []FileFilter) (string, error) {
	return "", ErrFallbackRequired
}

// SaveFile fallback for other platforms
func SaveFile(title string, defaultFilename string, filters []FileFilter) (string, error) {
	return "", ErrFallbackRequired
}
