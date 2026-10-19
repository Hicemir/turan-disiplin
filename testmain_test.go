package main

import (
	"os"
	"testing"
)

// TestMain isolates tests from the user's real config/database: integration
// tests call SaveConfig/ApplyTheme and open the SQLite store, which would
// otherwise overwrite the real TuranDisiplin settings (e.g. the saved theme).
func TestMain(m *testing.M) {
	tmp, err := os.MkdirTemp("", "turandisiplin-test-")
	if err != nil {
		os.Exit(1)
	}
	os.Setenv("XDG_CONFIG_HOME", tmp)
	os.Setenv("APPDATA", tmp)
	os.Setenv("HOME", tmp)
	code := m.Run()
	os.RemoveAll(tmp)
	os.Exit(code)
}
