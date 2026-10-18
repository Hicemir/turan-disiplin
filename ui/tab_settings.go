package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"okulevrak/config"
	appTheme "okulevrak/theme"
)

// BuildSettingsView creates the Theme and Interface Preferences view
func BuildSettingsView(state *AppState) fyne.CanvasObject {
	// UI Mode Selector
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

	// Theme Selector
	themeMap := map[string]string{
		"Açık Tema":   "light",
		"Sarı / Krem": "cream",
		"Koyu Tema":   "dark",
	}
	var currentThemeLabel string
	switch state.Config.ThemeMode {
	case "cream":
		currentThemeLabel = "Sarı / Krem"
	case "dark":
		currentThemeLabel = "Koyu Tema"
	default:
		currentThemeLabel = "Açık Tema"
	}

	themeRadio := widget.NewRadioGroup([]string{
		"Açık Tema",
		"Sarı / Krem",
		"Koyu Tema",
	}, nil)
	themeRadio.SetSelected(currentThemeLabel)

	// Font Size Selector
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

	// Save button
	applyBtn := widget.NewButtonWithIcon("Ayarları Uygula ve Kaydet", theme.DocumentSaveIcon(), func() {
		newUIMode := uiModeMap[uiRadio.Selected]
		newThemeMode := themeMap[themeRadio.Selected]
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

	var backBtn *widget.Button
	if state.Config.UIMode == "legacy" {
		backBtn = widget.NewButtonWithIcon("❌ KAPAT", theme.CancelIcon(), func() {
			state.NavigateTo("hub")
		})
		backBtn.Importance = widget.DangerImportance
	} else {
		backBtn = widget.NewButtonWithIcon("⬅️ Ana Menüye Dön", theme.NavigateBackIcon(), func() {
			state.NavigateTo("hub")
		})
	}

	form := widget.NewForm(
		widget.NewFormItem("Ana Ekran Görünüm Düzeni", uiRadio),
		widget.NewFormItem("Renk Teması", themeRadio),
		widget.NewFormItem("Yazı Boyutu", sizeRadio),
	)

	card := widget.NewCard("⚙️ Görünüm ve Tema Ayarları", "Ekran tasarımını, tema rengini ve yazı boyutunu seçiniz", container.NewVBox(
		form,
		widget.NewSeparator(),
		container.NewHBox(applyBtn, widget.NewSeparator(), backBtn),
	))

	header := container.NewBorder(
		nil, nil,
		backBtn,
		nil,
		widget.NewLabelWithStyle("⚙️ Ayarlar & Tema Tercihleri", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
	)

	return container.NewVScroll(container.NewVBox(
		header,
		widget.NewSeparator(),
		card,
	))
}
