package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

// BuildMainUI constructs the dynamic application container
func BuildMainUI(state *AppState) fyne.CanvasObject {
	mainHolder := container.NewStack()

	renderCurrentView := func() {
		mainHolder.Objects = nil

		if state.CurrentView == "hub" || state.CurrentView == "" {
			mainHolder.Add(BuildHubView(state))
			mainHolder.Refresh()
			return
		}

		// Subview Navigation: "Geri Dön" in Legacy Mode, "Ana Menüye Dön" in Modern Mode
		var backBtn *widget.Button
		if state.Config.UIMode == "legacy" {
			backBtn = widget.NewButtonWithIcon("Geri Dön", theme.NavigateBackIcon(), func() {
				state.NavigateTo("hub")
			})
			backBtn.Importance = widget.HighImportance
		} else {
			backBtn = widget.NewButtonWithIcon("Ana Menüye Dön", theme.NavigateBackIcon(), func() {
				state.NavigateTo("hub")
			})
			backBtn.Importance = widget.HighImportance
		}

		yearBtn := widget.NewButtonWithIcon("Öğretim Yılı: "+state.AcademicYear, theme.HistoryIcon(), func() {
			ShowAcademicYearPicker(state)
		})

		var subContent fyne.CanvasObject
		var pageTitle string

		switch state.CurrentView {
		case "school":
			pageTitle = "🏛️ Genel Bilgiler (Kurum ve Kurul Üyeleri)"
			subContent = BuildSchoolInfoTab(state)
		case "students":
			pageTitle = "👥 Öğrenci İşleri"
			subContent = BuildStudentsTab(state)
		case "incidents":
			pageTitle = "🐞 Olay Girişi ve İnceleme Profilleri"
			subContent = BuildIncidentsTab(state)
		case "verdict":
			pageTitle = "⚖️ DİSİPLİN KURULU KARARI (EK-1 / EK-10)"
			subContent = BuildVerdictTab(state)
		case "summons":
			pageTitle = "📨 Çağrı Pusulaları (Öğrenci & Kurul)"
			subContent = BuildSummonsTab(state)
		case "forms", "written_warning":
			pageTitle = "📝 Resmi Tutanak ve İfade/Savunma Formları"
			subContent = BuildFormsTab(state)
		case "notices":
			pageTitle = "📬 Öğrenci ve Veli Ceza Tebligat Pusulaları"
			subContent = BuildNoticesTab(state)
		case "restoration":
			pageTitle = "🔄 Davranış Puanı İadesi ve Yaptırım Kaldırma"
			subContent = BuildRestorationTab(state)
		case "minutes":
			pageTitle = "☀️ Sene Başı / Sene Sonu Kurul Tutanakları"
			subContent = BuildMinutesTab(state)
		case "dizi":
			pageTitle = "📚 Soruşturma Dosyası Dizi Pusulası"
			subContent = BuildMinutesTab(state)
		case "stats":
			pageTitle = "📊 Ceza İstatistikleri ve Raporlar"
			subContent = BuildStatsTab(state)
		case "yonetmelik":
			pageTitle = "🗂️ Yönetmelik Maddeleri ve Şablon Düzenleyici"
			subContent = BuildYonetmelikTab(state)
		case "settings":
			pageTitle = "⚙️ Görünüm ve Tema Ayarları"
			subContent = BuildSettingsView(state)
		default:
			pageTitle = "Ana Menü"
			subContent = BuildHubView(state)
		}

		topBar := container.NewBorder(
			nil, nil,
			backBtn,
			yearBtn,
			widget.NewLabelWithStyle(pageTitle, fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		)

		viewContainer := container.NewBorder(
			container.NewVBox(topBar, widget.NewSeparator()),
			nil, nil, nil,
			subContent,
		)

		mainHolder.Add(viewContainer)
		mainHolder.Refresh()
	}

	state.OnViewRefresh = renderCurrentView

	renderCurrentView()
	return mainHolder
}
