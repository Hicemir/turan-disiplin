package ui

import (
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

// systemPrefersDark reports whether the OS is set to a dark color scheme.
// Fyne's own variant is trusted when it says "dark"; otherwise (Fyne can report
// light on Linux desktops where the portal is unavailable) a few well-known
// Linux desktop settings are probed as a fallback.
func systemPrefersDark(a fyne.App) bool {
	if a != nil && a.Settings().ThemeVariant() == theme.VariantDark {
		return true
	}
	if runtime.GOOS != "linux" {
		return false
	}

	// GNOME / freedesktop color-scheme preference
	if out, err := exec.Command("gsettings", "get", "org.gnome.desktop.interface", "color-scheme").Output(); err == nil {
		if strings.Contains(strings.ToLower(string(out)), "dark") {
			return true
		}
	}
	// GTK theme name
	if strings.Contains(strings.ToLower(os.Getenv("GTK_THEME")), "dark") {
		return true
	}
	if out, err := exec.Command("gsettings", "get", "org.gnome.desktop.interface", "gtk-theme").Output(); err == nil {
		if strings.Contains(strings.ToLower(string(out)), "dark") {
			return true
		}
	}
	// KDE Plasma: ColorScheme in kdeglobals (e.g. BreezeDark)
	if home, err := os.UserHomeDir(); err == nil {
		if b, err := os.ReadFile(filepath.Join(home, ".config", "kdeglobals")); err == nil {
			for _, line := range strings.Split(string(b), "\n") {
				if strings.HasPrefix(line, "ColorScheme=") {
					return strings.Contains(strings.ToLower(line), "dark")
				}
			}
		}
	}
	return false
}
