package ui

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"turandisiplin/config"
	appTheme "turandisiplin/theme"
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

	// 2. Custom Colors and Themes
	currentCustomBg := config.HexToColor(state.Config.CustomBgColor, color.NRGBA{R: 30, G: 30, B: 46, A: 255})
	currentCustomAccent := config.HexToColor(state.Config.CustomAccentColor, color.NRGBA{R: 137, G: 180, B: 250, A: 255})

	allThemes := appTheme.GetAllThemes(currentCustomBg, currentCustomAccent)
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

	// Custom Color Picker Card
	bgHexLabel := widget.NewLabelWithStyle(config.ColorToHex(currentCustomBg), fyne.TextAlignLeading, fyne.TextStyle{Monospace: true})
	accentHexLabel := widget.NewLabelWithStyle(config.ColorToHex(currentCustomAccent), fyne.TextAlignLeading, fyne.TextStyle{Monospace: true})

	customBgPreview := canvas.NewRectangle(currentCustomBg)
	customAccentPreview := canvas.NewRectangle(currentCustomAccent)

	makeSwatchBox := func(rect *canvas.Rectangle) fyne.CanvasObject {
		rect.SetMinSize(fyne.NewSize(26, 26))
		rect.CornerRadius = 4.0
		rect.StrokeColor = color.NRGBA{R: 255, G: 255, B: 255, A: 230}
		rect.StrokeWidth = 2.0

		frame := canvas.NewRectangle(color.NRGBA{R: 35, G: 35, B: 35, A: 220})
		frame.SetMinSize(fyne.NewSize(30, 30))
		frame.CornerRadius = 6.0
		frame.StrokeColor = color.NRGBA{R: 120, G: 120, B: 120, A: 200}
		frame.StrokeWidth = 1.0

		return container.NewStack(frame, container.NewCenter(rect))
	}

	bgPickBtn := widget.NewButtonWithIcon("Arka Plan Rengi Seç", theme.ColorPaletteIcon(), func() {
		ShowAdvancedColorPicker("Arka Plan Rengi Seçin", currentCustomBg, func(c color.Color) {
			if c != nil {
				currentCustomBg = config.HexToColor(config.ColorToHex(c), color.NRGBA{R: 30, G: 30, B: 46, A: 255})
				customBgPreview.FillColor = currentCustomBg
				customBgPreview.Refresh()
				bgHexLabel.SetText(config.ColorToHex(currentCustomBg))
				selectedThemeID = "custom"
				allThemes = appTheme.GetAllThemes(currentCustomBg, currentCustomAccent)
				refreshThemeRows()
			}
		}, state.Window)
	})

	accentPickBtn := widget.NewButtonWithIcon("Vurgu / Buton Rengi Seç", theme.ColorPaletteIcon(), func() {
		ShowAdvancedColorPicker("Vurgu / Buton Rengi Seçin", currentCustomAccent, func(c color.Color) {
			if c != nil {
				currentCustomAccent = config.HexToColor(config.ColorToHex(c), color.NRGBA{R: 137, G: 180, B: 250, A: 255})
				customAccentPreview.FillColor = currentCustomAccent
				customAccentPreview.Refresh()
				accentHexLabel.SetText(config.ColorToHex(currentCustomAccent))
				selectedThemeID = "custom"
				allThemes = appTheme.GetAllThemes(currentCustomBg, currentCustomAccent)
				refreshThemeRows()
			}
		}, state.Window)
	})

	customColorBox := widget.NewCard("🎨 Özel Tema Renk Paleti", "Kendi renklerinizi seçerek özel temanızı oluşturun",
		container.NewVBox(
			container.NewHBox(makeSwatchBox(customBgPreview), bgHexLabel, bgPickBtn),
			container.NewHBox(makeSwatchBox(customAccentPreview), accentHexLabel, accentPickBtn),
		),
	)

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

	// 4. Startup Preferences
	rememberYearCheck := widget.NewCheck("Açılışta Öğretim Yılı Sorma (Kayıtlı Yıl ile Doğrudan Başla)", nil)
	rememberYearCheck.SetChecked(state.Config.RememberYear)

	// Save button
	applyBtn := widget.NewButtonWithIcon("Ayarları Uygula ve Kaydet", theme.DocumentSaveIcon(), func() {
		newUIMode := uiModeMap[uiRadio.Selected]
		newThemeMode := selectedThemeID
		newFontSize := sizeMap[sizeRadio.Selected]

		state.Config.UIMode = newUIMode
		state.Config.ThemeMode = newThemeMode
		state.Config.FontSize = newFontSize
		state.Config.RememberYear = rememberYearCheck.Checked
		state.Config.CustomBgColor = config.ColorToHex(currentCustomBg)
		state.Config.CustomAccentColor = config.ColorToHex(currentCustomAccent)
		_ = config.SaveConfig(state.Config)

		state.ApplyTheme(newThemeMode, newFontSize)

		ShowTurkishInfo("Başarılı", "Görünüm ve tema ayarlarınız başarıyla güncellendi.", state.Window)
	})
	applyBtn.Importance = widget.HighImportance

	// 5. Backup & Restore Card
	backupCard := widget.NewCard(
		"💾 Veritabanı Yedekleme ve Kurtarma",
		"Tüm okul, öğrenci, olay ve disiplin kayıtlarınızı güvenle yedekleyin veya eski yedeği geri yükleyin",
		container.NewHBox(
			widget.NewButtonWithIcon("Veritabanını Yedekle", theme.DownloadIcon(), func() {
				BackupDatabase(state)
			}),
			widget.NewButtonWithIcon("Yedekten Geri Yükle", theme.UploadIcon(), func() {
				RestoreDatabase(state)
			}),
		),
	)

	form := container.NewVBox(
		widget.NewLabelWithStyle("Ana Ekran Görünüm Düzeni", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		uiRadio,
		widget.NewSeparator(),
		widget.NewLabelWithStyle("Renk Teması (Arka Plan ve Vurgu Rengi)", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		themeContainer,
		customColorBox,
		widget.NewSeparator(),
		widget.NewLabelWithStyle("Yazı Boyutu", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		sizeRadio,
		widget.NewSeparator(),
		widget.NewLabelWithStyle("Açılış Tercihleri", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		rememberYearCheck,
		widget.NewSeparator(),
		backupCard,
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
