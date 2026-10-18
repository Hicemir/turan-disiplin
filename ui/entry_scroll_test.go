package ui

import (
	"testing"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/test"
	"fyne.io/fyne/v2/widget"
)

func TestSmartMultiLineEntryFocusAndScroll(t *testing.T) {
	test.NewApp()

	parentContent := widget.NewLabel("Long content")
	parentScroll := container.NewVScroll(parentContent)

	entry := NewSmartMultiLineEntry(parentScroll)
	if entry.isFocused {
		t.Fatalf("expected entry to start unfocused")
	}

	// When unfocused, scrolling should forward to parentScroll without error
	ev := &fyne.ScrollEvent{
		Scrolled: fyne.NewDelta(0, -50),
	}
	entry.Scrolled(ev)

	entry.FocusGained()
	if !entry.isFocused {
		t.Fatalf("expected entry to be focused after FocusGained")
	}

	// When focused, scrolling handled internally
	entry.Scrolled(ev)

	entry.FocusLost()
	if entry.isFocused {
		t.Fatalf("expected entry to be unfocused after FocusLost")
	}
}
