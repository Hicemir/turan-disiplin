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

	// 2. Color Themes
	allThemes := appTheme.GetAllThemes()
	selectedThemeID := state.Config.ThemeMode
	if selectedThemeID == "" {
		selectedThemeID = "light"
	}

	themeContainer := container.NewVBox()

	type themeRowRefs struct {
		id           string
		radioIcon    *widget.Icon
		accentSquare *canvas.Rectangle
		bgSquare     *canvas.Rectangle
		label        *widget.Label
	}
	var themeRows []*themeRowRefs

	refreshThemeRows := func() {
		for i, r := range themeRows {
			isSelected := r.id == selectedThemeID
			if isSelected {
				r.radioIcon.SetResource(theme.RadioButtonCheckedIcon())
			} else {
				r.radioIcon.SetResource(theme.RadioButtonIcon())
			}
			r.label.TextStyle = fyne.TextStyle{Bold: isSelected}
			r.label.Refresh()
			if i < len(allThemes) {
				r.bgSquare.FillColor = allThemes[i].BgColor
				r.bgSquare.Refresh()
				r.accentSquare.FillColor = allThemes[i].AccentColor
				r.accentSquare.Refresh()
			}
		}
	}

	for _, th := range allThemes {
		thCopy := th
		isSelected := thCopy.ID == selectedThemeID

		var radioIcon *widget.Icon
		if isSelected {
			radioIcon = widget.NewIcon(theme.RadioButtonCheckedIcon())
		} else {
			radioIcon = widget.NewIcon(theme.RadioButtonIcon())
		}

		bgSquare := canvas.NewRectangle(thCopy.BgColor)
		bgSquare.SetMinSize(fyne.NewSize(20, 20))
		bgSquare.StrokeColor = color.NRGBA{R: 160, G: 160, B: 160, A: 255}
		bgSquare.StrokeWidth = 1.0
		bgSquare.CornerRadius = 3.0

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

		themeRows = append(themeRows, &themeRowRefs{
			id:           thCopy.ID,
			radioIcon:    radioIcon,
			bgSquare:     bgSquare,
			accentSquare: accentSquare,
			label:        label,
		})
		themeContainer.Add(container.NewStack(btn, container.NewPadded(rowContent)))
	}

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
		widget.NewSeparator(),
		widget.NewLabelWithStyle("Yazı Boyutu", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		sizeRadio,
		widget.NewSeparator(),
		widget.NewLabelWithStyle("Açılış Tercihleri", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		rememberYearCheck,
		widget.NewSeparator(),
		backupCard,
	)

	applyBtnBottom := widget.NewButtonWithIcon("Ayarları Uygula ve Kaydet", theme.DocumentSaveIcon(), func() {
		applyBtn.OnTapped()
	})
	applyBtnBottom.Importance = widget.HighImportance

	card := widget.NewCard(
		"⚙️ Tercihler",
		"Ekran tasarımını, tema rengini ve yazı boyutunu seçiniz",
		container.NewVBox(
			form,
			widget.NewSeparator(),
			container.NewHBox(applyBtnBottom),
		),
	)

	topSaveBar := container.NewBorder(
		nil, nil,
		widget.NewLabelWithStyle("⚙️ Görünüm ve Tema Ayarları", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		applyBtn,
	)

	scrollContent := container.NewVScroll(container.NewPadded(card))
	return container.NewBorder(
		container.NewPadded(topSaveBar),
		nil, nil, nil,
		scrollContent,
	)
}
