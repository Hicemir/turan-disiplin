package ui

import (
	"fmt"
	"image/color"
	"math"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	fyneDialog "fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"

	"okulevrak/config"
)

// ShowAdvancedColorPicker displays a rich, interactive color picker modal with
// RGB sliders, live hex input, curated palette swatches, and high-contrast preview.
func ShowAdvancedColorPicker(title string, initialColor color.Color, onChosen func(color.Color), parent fyne.Window) {
	if initialColor == nil {
		initialColor = color.NRGBA{R: 30, G: 30, B: 46, A: 255}
	}

	rVal, gVal, bVal, _ := initialColor.RGBA()
	curR := float64(uint8(rVal >> 8))
	curG := float64(uint8(gVal >> 8))
	curB := float64(uint8(bVal >> 8))

	currentColor := color.NRGBA{R: uint8(curR), G: uint8(curG), B: uint8(curB), A: 255}

	// Live Preview Swatch (with high-contrast double border so it's always visible)
	previewRect := canvas.NewRectangle(currentColor)
	previewRect.SetMinSize(fyne.NewSize(120, 50))
	previewRect.CornerRadius = 8.0
	previewRect.StrokeColor = color.NRGBA{R: 255, G: 255, B: 255, A: 230}
	previewRect.StrokeWidth = 2.0

	previewOuterFrame := canvas.NewRectangle(color.NRGBA{R: 35, G: 35, B: 35, A: 220})
	previewOuterFrame.SetMinSize(fyne.NewSize(126, 56))
	previewOuterFrame.CornerRadius = 10.0
	previewOuterFrame.StrokeColor = color.NRGBA{R: 120, G: 120, B: 120, A: 200}
	previewOuterFrame.StrokeWidth = 1.0

	previewStack := container.NewStack(previewOuterFrame, container.NewCenter(previewRect))

	hexLabel := widget.NewLabelWithStyle(config.ColorToHex(currentColor), fyne.TextAlignCenter, fyne.TextStyle{Bold: true, Monospace: true})

	// Hex Input Entry
	hexEntry := widget.NewEntry()
	hexEntry.SetText(config.ColorToHex(currentColor))
	hexEntry.SetPlaceHolder("#RRGGBB")

	// Sliders & Value Labels
	rLabel := widget.NewLabel(fmt.Sprintf("%3d", int(curR)))
	gLabel := widget.NewLabel(fmt.Sprintf("%3d", int(curG)))
	bLabel := widget.NewLabel(fmt.Sprintf("%3d", int(curB)))

	rSlider := widget.NewSlider(0, 255)
	rSlider.SetValue(curR)
	gSlider := widget.NewSlider(0, 255)
	gSlider.SetValue(curG)
	bSlider := widget.NewSlider(0, 255)
	bSlider.SetValue(curB)

	updating := false

	updateFromRGB := func() {
		if updating {
			return
		}
		updating = true
		currentColor = color.NRGBA{
			R: uint8(math.Round(rSlider.Value)),
			G: uint8(math.Round(gSlider.Value)),
			B: uint8(math.Round(bSlider.Value)),
			A: 255,
		}
		previewRect.FillColor = currentColor
		previewRect.Refresh()

		hexStr := config.ColorToHex(currentColor)
		hexLabel.SetText(hexStr)
		hexEntry.SetText(hexStr)

		rLabel.SetText(fmt.Sprintf("%3d", int(rSlider.Value)))
		gLabel.SetText(fmt.Sprintf("%3d", int(gSlider.Value)))
		bLabel.SetText(fmt.Sprintf("%3d", int(bSlider.Value)))
		updating = false
	}

	rSlider.OnChanged = func(v float64) { updateFromRGB() }
	gSlider.OnChanged = func(v float64) { updateFromRGB() }
	bSlider.OnChanged = func(v float64) { updateFromRGB() }

	hexEntry.OnChanged = func(s string) {
		if updating {
			return
		}
		clean := strings.TrimSpace(s)
		if !strings.HasPrefix(clean, "#") {
			clean = "#" + clean
		}
		if len(clean) == 7 {
			fallback := color.NRGBA{A: 0}
			c := config.HexToColor(clean, fallback)
			if c.A != 0 {
				updating = true
				rVal := float64(c.R)
				gVal := float64(c.G)
				bVal := float64(c.B)

				currentColor = c
				previewRect.FillColor = currentColor
				previewRect.Refresh()

				rSlider.SetValue(rVal)
				gSlider.SetValue(gVal)
				bSlider.SetValue(bVal)

				rLabel.SetText(fmt.Sprintf("%3d", int(rVal)))
				gLabel.SetText(fmt.Sprintf("%3d", int(gVal)))
				bLabel.SetText(fmt.Sprintf("%3d", int(bVal)))
				hexLabel.SetText(clean)
				updating = false
			}
		}
	}

	// Curated Palette of Preset Swatches (18 colors: darks, lights, and vibrant accents)
	presetHexes := []string{
		"#1e1e2e", "#181825", "#11111b", "#1a1b26", "#0f172a", "#1e293b",
		"#ffffff", "#f8fafc", "#f1f5f9", "#e2e8f0", "#94a3b8", "#475569",
		"#ef4444", "#f97316", "#eab308", "#22c55e", "#06b6d4", "#3b82f6",
		"#8b5cf6", "#ec4899", "#d946ef", "#0ea5e9", "#14b8a6", "#84cc16",
	}

	paletteGrid := container.NewGridWithColumns(6)
	for _, hexCode := range presetHexes {
		h := hexCode
		c := config.HexToColor(h, color.NRGBA{R: 255, G: 255, B: 255, A: 255})
		targetColor := c
		targetHex := h

		sw := newClickableSwatch(targetColor, func() {
			updating = true
			rVal := float64(targetColor.R)
			gVal := float64(targetColor.G)
			bVal := float64(targetColor.B)

			currentColor = targetColor
			previewRect.FillColor = currentColor
			previewRect.Refresh()

			rSlider.SetValue(rVal)
			gSlider.SetValue(gVal)
			bSlider.SetValue(bVal)

			rLabel.SetText(fmt.Sprintf("%3d", int(rVal)))
			gLabel.SetText(fmt.Sprintf("%3d", int(gVal)))
			bLabel.SetText(fmt.Sprintf("%3d", int(bVal)))
			hexLabel.SetText(targetHex)
			hexEntry.SetText(targetHex)
			updating = false
		})
		paletteGrid.Add(sw)
	}

	// Sliders Container
	rRow := container.NewBorder(nil, nil, widget.NewLabelWithStyle("Kırmızı (R):", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}), rLabel, rSlider)
	gRow := container.NewBorder(nil, nil, widget.NewLabelWithStyle("Yeşil (G):", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}), gLabel, gSlider)
	bRow := container.NewBorder(nil, nil, widget.NewLabelWithStyle("Mavi (B):", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}), bLabel, bSlider)

	slidersBox := container.NewVBox(rRow, gRow, bRow)

	hexRow := container.NewBorder(nil, nil, widget.NewLabelWithStyle("HEX Kodu:", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}), nil, hexEntry)

	dialogContent := container.NewVBox(
		container.NewCenter(container.NewVBox(previewStack, hexLabel)),
		widget.NewSeparator(),
		widget.NewLabelWithStyle("🎨 Hızlı Renk Paleti:", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		paletteGrid,
		widget.NewSeparator(),
		widget.NewLabelWithStyle("🎛️ Hassas RGB Renk Ayarı:", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		slidersBox,
		hexRow,
	)

	d := fyneDialog.NewCustomConfirm(
		title,
		"Rengi Seç",
		"İptal",
		dialogContent,
		func(ok bool) {
			if ok && onChosen != nil {
				onChosen(currentColor)
			}
		},
		parent,
	)
	d.Resize(fyne.NewSize(420, 500))
	d.Show()
}

// clickableSwatch is a clean, frameless-safe tappable color rectangle for palettes
type clickableSwatch struct {
	widget.BaseWidget
	rect     *canvas.Rectangle
	onTapped func()
}

func newClickableSwatch(c color.Color, onTapped func()) *clickableSwatch {
	w := &clickableSwatch{
		rect:     canvas.NewRectangle(c),
		onTapped: onTapped,
	}
	w.rect.SetMinSize(fyne.NewSize(36, 28))
	w.rect.CornerRadius = 5.0
	w.rect.StrokeColor = color.NRGBA{R: 200, G: 200, B: 200, A: 200}
	w.rect.StrokeWidth = 1.0
	w.ExtendBaseWidget(w)
	return w
}

func (s *clickableSwatch) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(s.rect)
}

func (s *clickableSwatch) Tapped(_ *fyne.PointEvent) {
	if s.onTapped != nil {
		s.onTapped()
	}
}

func (s *clickableSwatch) MinSize() fyne.Size {
	return fyne.NewSize(36, 28)
}

