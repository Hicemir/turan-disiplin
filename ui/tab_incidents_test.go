package ui

import (
	"testing"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/test"
	"fyne.io/fyne/v2/widget"
)

func TestIncidentBorderLayoutIndices(t *testing.T) {
	test.NewApp()

	codeLbl := widget.NewLabelWithStyle("#000", fyne.TextAlignLeading, fyne.TextStyle{Bold: true})
	dateLbl := widget.NewLabel("00.00.0000")
	titleLbl := widget.NewLabel("Olay Adı / Başlığı")

	row := container.NewBorder(nil, nil, codeLbl, dateLbl, titleLbl)

	if row.Objects[0] != titleLbl {
		t.Fatalf("expected objects[0] to be titleLbl")
	}
	if row.Objects[1] != codeLbl {
		t.Fatalf("expected objects[1] to be codeLbl")
	}
	if row.Objects[2] != dateLbl {
		t.Fatalf("expected objects[2] to be dateLbl")
	}
}

func TestFixedWidthLayout(t *testing.T) {
	test.NewApp()

	w := widget.NewLabel("Test")
	c := container.New(&fixedWidthLayout{width: 300}, w)

	minSize := c.MinSize()
	if minSize.Width != 300 {
		t.Fatalf("expected min width 300, got %f", minSize.Width)
	}

	c.Resize(fyne.NewSize(300, 500))
	if w.Size().Width != 300 || w.Size().Height != 500 {
		t.Fatalf("expected child to be resized to (300, 500), got %v", w.Size())
	}
}
