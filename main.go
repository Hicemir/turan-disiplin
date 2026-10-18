package main

import (
	"log"
	"runtime/debug"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"

	"turandisiplin/config"
	"turandisiplin/ui"
)

func init() {
	// Proactive memory debloating:
	// Set target GC percentage to 20 to collect garbage frequently,
	// and set memory limit to 80MB so Go releases memory aggressively to OS.
	debug.SetGCPercent(20)
	debug.SetMemoryLimit(80 * 1024 * 1024)
}

func main() {
	a := app.NewWithID("com.turandisiplin.app")

	w := a.NewWindow("TuranDisiplin - MEB Öğrenci Disiplin ve Davranış Değerlendirme Sistemi")
	w.Resize(fyne.NewSize(1260, 840))
	w.CenterOnScreen()

	// Initialize application state (automatically applies saved theme & UI mode)
	state, err := ui.InitState(a, w)
	if err != nil {
		log.Fatalf("Uygulama durumu başlatılamadı: %v", err)
	}

	// Build the main UI (starts at the Hub view)
	mainUI := ui.BuildMainUI(state)
	w.SetContent(mainUI)

	// Startup Academic Year Selector Modal
	showYearSelectionModal(w, state)

	w.ShowAndRun()
}

func showYearSelectionModal(w fyne.Window, state *ui.AppState) {
	if state.Config.RememberYear && state.AcademicYear != "" {
		return
	}

	yearOptions := state.Store.GetAcademicYears()
	foundCurrent := false
	for _, y := range yearOptions {
		if y == state.AcademicYear {
			foundCurrent = true
			break
		}
	}
	if !foundCurrent {
		yearOptions = append([]string{state.AcademicYear}, yearOptions...)
	}

	yearSelect := ui.NewTurkishSelect(yearOptions, nil)
	yearSelect.SetSelected(state.AcademicYear)

	newYearEntry := ui.NewSingleLineEntry()
	newYearEntry.SetPlaceHolder("veya yeni öğretim yılı yazınız (örn: 2027 - 2028)")

	schoolTypeRadio := widget.NewRadioGroup([]string{"LİSE", "ORTAOKUL"}, nil)
	schoolTypeRadio.Horizontal = true
	schoolTypeRadio.SetSelected(state.SchoolInfo.SchoolType)

	rememberCheck := widget.NewCheck("Bu seçimi hatırla ve açılışta tekrar sorma", nil)
	rememberCheck.SetChecked(state.Config.RememberYear)

	form := widget.NewForm(
		widget.NewFormItem("Öğretim Yılı Seçin", yearSelect),
		widget.NewFormItem("Yeni Öğretim Yılı", newYearEntry),
		widget.NewFormItem("Okul Kademesi", schoolTypeRadio),
		widget.NewFormItem("", rememberCheck),
	)

	d := dialog.NewCustomConfirm(
		"📅 Öğretim Yılı ve Kurum Seçimi",
		"Programa Başla",
		"Varsayılan ile Devam Et",
		container.NewVBox(
			widget.NewLabelWithStyle("Hoş Geldiniz! Lütfen işlem yapmak istediğiniz öğretim yılını seçiniz:", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
			form,
		),
		func(ok bool) {
			chosenYear := yearSelect.Selected
			if newYearEntry.Text != "" {
				chosenYear = newYearEntry.Text
			}
			if chosenYear != "" {
				state.SetAcademicYear(chosenYear)
			}
			state.SetSchoolType(schoolTypeRadio.Selected)

			state.Config.RememberYear = rememberCheck.Checked
			if rememberCheck.Checked {
				state.Config.DefaultYear = chosenYear
			}
			_ = config.SaveConfig(state.Config)
		},
		w,
	)
	d.Resize(fyne.NewSize(540, 310))
	d.Show()
}
