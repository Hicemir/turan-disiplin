//go:build linux

package picker

import (
	"bytes"
	"os/exec"
	"strings"
)

// OpenFile displays the native Linux file open dialog using zenity or kdialog
func OpenFile(title string, filters []FileFilter) (string, error) {
	// 1. Try zenity (GNOME, Cinnamon, XFCE, etc.)
	if zenityPath, err := exec.LookPath("zenity"); err == nil {
		args := []string{"--file-selection"}
		if title != "" {
			args = append(args, "--title="+title)
		}
		for _, f := range filters {
			var patterns []string
			for _, ext := range f.Extensions {
				clean := strings.TrimPrefix(ext, ".")
				patterns = append(patterns, "*."+clean)
			}
			filterArg := f.Name + " | " + strings.Join(patterns, " ")
			args = append(args, "--file-filter="+filterArg)
		}

		cmd := exec.Command(zenityPath, args...)
		var out bytes.Buffer
		cmd.Stdout = &out
		err := cmd.Run()
		if err == nil {
			path := strings.TrimSpace(out.String())
			if path != "" {
				return path, nil
			}
		}
		if cmd.ProcessState != nil && cmd.ProcessState.ExitCode() == 1 {
			return "", ErrCancelled
		}
	}

	// 2. Try kdialog (KDE Plasma)
	if kdialogPath, err := exec.LookPath("kdialog"); err == nil {
		var patterns []string
		for _, f := range filters {
			for _, ext := range f.Extensions {
				clean := strings.TrimPrefix(ext, ".")
				patterns = append(patterns, "*."+clean)
			}
		}
		filterArg := strings.Join(patterns, " ")
		args := []string{"--getopenfilename", ".", filterArg}
		if title != "" {
			args = append([]string{"--title", title}, args...)
		}

		cmd := exec.Command(kdialogPath, args...)
		var out bytes.Buffer
		cmd.Stdout = &out
		err := cmd.Run()
		if err == nil {
			path := strings.TrimSpace(out.String())
			if path != "" {
				return path, nil
			}
		}
		if cmd.ProcessState != nil && cmd.ProcessState.ExitCode() == 1 {
			return "", ErrCancelled
		}
	}

	// 3. Fallback required if CLI tools are missing
	return "", ErrFallbackRequired
}

// SaveFile displays the native Linux file save dialog using zenity or kdialog
func SaveFile(title string, defaultFilename string, filters []FileFilter) (string, error) {
	// 1. Try zenity
	if zenityPath, err := exec.LookPath("zenity"); err == nil {
		args := []string{"--file-selection", "--save"}
		if title != "" {
			args = append(args, "--title="+title)
		}
		if defaultFilename != "" {
			args = append(args, "--filename="+defaultFilename)
		}
		for _, f := range filters {
			var patterns []string
			for _, ext := range f.Extensions {
				clean := strings.TrimPrefix(ext, ".")
				patterns = append(patterns, "*."+clean)
			}
			filterArg := f.Name + " | " + strings.Join(patterns, " ")
			args = append(args, "--file-filter="+filterArg)
		}

		cmd := exec.Command(zenityPath, args...)
		var out bytes.Buffer
		cmd.Stdout = &out
		err := cmd.Run()
		if err == nil {
			path := strings.TrimSpace(out.String())
			if path != "" {
				return path, nil
			}
		}
		if cmd.ProcessState != nil && cmd.ProcessState.ExitCode() == 1 {
			return "", ErrCancelled
		}
	}

	// 2. Try kdialog
	if kdialogPath, err := exec.LookPath("kdialog"); err == nil {
		var patterns []string
		for _, f := range filters {
			for _, ext := range f.Extensions {
				clean := strings.TrimPrefix(ext, ".")
				patterns = append(patterns, "*."+clean)
			}
		}
		filterArg := strings.Join(patterns, " ")
		startFile := "."
		if defaultFilename != "" {
			startFile = "./" + defaultFilename
		}
		args := []string{"--getsavefilename", startFile, filterArg}
		if title != "" {
			args = append([]string{"--title", title}, args...)
		}

		cmd := exec.Command(kdialogPath, args...)
		var out bytes.Buffer
		cmd.Stdout = &out
		err := cmd.Run()
		if err == nil {
			path := strings.TrimSpace(out.String())
			if path != "" {
				return path, nil
			}
		}
		if cmd.ProcessState != nil && cmd.ProcessState.ExitCode() == 1 {
			return "", ErrCancelled
		}
	}

	// 3. Fallback required if CLI tools are missing
	return "", ErrFallbackRequired
}
