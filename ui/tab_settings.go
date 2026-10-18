package ui

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"okulevrak/config"
	appTheme "okulevrak/theme"
)

// BuildSettingsView creates the Theme and Interface Preferences view
func BuildSettingsView(state *AppState) fyne.CanvasObject {
	// 1. UI Mode Selector (Modern vs Klasik)
	uiModeMap := map[string]string{
		"Modern Arayüz": "modern",
		"Klasik Arayüz": "legacy",
	}
	var currentUIModeLabel string
	if state.Config.UIMode == "legacy" {
		currentUIModeLabel = "Klasik Arayüz"
	} else {
		currentUIModeLabel = "Modern Arayüz"
	}

	uiRadio := widget.NewRadioGroup([]string{
		"Modern Arayüz",
		"Klasik Arayüz",
	}, nil)
	uiRadio.SetSelected(currentUIModeLabel)

	// 2. Color Theme Selector with 2-Square Color Previews
	allThemes := appTheme.GetAllThemes()
	selectedThemeID := state.Config.ThemeMode
	if selectedThemeID == "" {
		selectedThemeID = "light"
	}

	themeContainer := container.NewVBox()

	var refreshThemeRows func()
	refreshThemeRows = func() {
		themeContainer.Objects = nil
		for _, th := range allThemes {
			thCopy := th
			isSelected := thCopy.ID == selectedThemeID

			var radioIcon *widget.Icon
			if isSelected {
				radioIcon = widget.NewIcon(theme.RadioButtonCheckedIcon())
			} else {
				radioIcon = widget.NewIcon(theme.RadioButtonIcon())
			}

			// Square 1: Background Color Swatch
			bgSquare := canvas.NewRectangle(thCopy.BgColor)
			bgSquare.SetMinSize(fyne.NewSize(20, 20))
			bgSquare.StrokeColor = color.NRGBA{R: 160, G: 160, B: 160, A: 255}
			bgSquare.StrokeWidth = 1.0
			bgSquare.CornerRadius = 3.0

			// Square 2: Accent / Button Color Swatch
			accentSquare := canvas.NewRectangle(thCopy.AccentColor)
			accentSquare.SetMinSize(fyne.NewSize(20, 20))
			accentSquare.StrokeColor = color.NRGBA{R: 100, G: 100, B: 100, A: 200}
			accentSquare.StrokeWidth = 1.0
			accentSquare.CornerRadius = 3.0

			swatches := container.NewHBox(bgSquare, accentSquare)
			label := widget.NewLabelWithStyle(thCopy.Name, fyne.TextAlignLeading, fyne.TextStyle{Bold: isSelected})

			rowContent := container.NewHBox(
				radioIcon,
				swatches,
				label,
			)

			btn := widget.NewButton("", func() {
				selectedThemeID = thCopy.ID
				refreshThemeRows()
			})

			themeContainer.Add(container.NewStack(btn, container.NewPadded(rowContent)))
		}
		themeContainer.Refresh()
	}
	refreshThemeRows()

	// 3. Font Size Selector (15pt, 17pt, 19pt)
	sizeMap := map[string]float32{
		"15pt": 15.0,
		"17pt": 17.0,
		"19pt": 19.0,
	}
	var currentSizeLabel string
	if state.Config.FontSize >= 18.5 {
		currentSizeLabel = "19pt"
	} else if state.Config.FontSize >= 16.5 {
		currentSizeLabel = "17pt"
	} else {
		currentSizeLabel = "15pt"
	}

	sizeRadio := widget.NewRadioGroup([]string{
		"15pt",
		"17pt",
		"19pt",
	}, nil)
	sizeRadio.SetSelected(currentSizeLabel)

	// Save button (Single, prominent action button - NO duplicate back buttons)
	applyBtn := widget.NewButtonWithIcon("Ayarları Uygula ve Kaydet", theme.DocumentSaveIcon(), func() {
		newUIMode := uiModeMap[uiRadio.Selected]
		newThemeMode := selectedThemeID
		newFontSize := sizeMap[sizeRadio.Selected]

		state.Config.UIMode = newUIMode
		state.Config.ThemeMode = newThemeMode
		state.Config.FontSize = newFontSize
		_ = config.SaveConfig(state.Config)

		if state.App != nil {
			state.App.Settings().SetTheme(appTheme.NewTurkishTheme(newThemeMode, newFontSize))
		}

		showInfo(state.Window, "Başarılı", "Görünüm ve tema ayarlarınız başarıyla güncellendi.")
	})
	applyBtn.Importance = widget.HighImportance

	form := container.NewVBox(
		widget.NewLabelWithStyle("Ana Ekran Görünüm Düzeni", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		uiRadio,
		widget.NewSeparator(),
		widget.NewLabelWithStyle("Renk Teması (Arka Plan ve Vurgu Rengi)", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		themeContainer,
		widget.NewSeparator(),
		widget.NewLabelWithStyle("Yazı Boyutu", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		sizeRadio,
	)

	card := widget.NewCard(
		"⚙️ Görünüm ve Tema Ayarları",
		"Ekran tasarımını, tema rengini ve yazı boyutunu seçiniz",
		container.NewVBox(
			form,
			widget.NewSeparator(),
			container.NewHBox(applyBtn),
		),
	)

	return container.NewVScroll(container.NewPadded(card))
}
