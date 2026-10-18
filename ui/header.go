package ui

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"okulevrak/pdf"
	"okulevrak/store"
)

// BuildHeader creates the persistent senior-friendly top bar
func BuildHeader(state *AppState) fyne.CanvasObject {
	title := widget.NewLabelWithStyle("🏛️ OKUL EVRAK", fyne.TextAlignLeading, fyne.TextStyle{Bold: true})

	// Academic Year Dropdown
	state.AcademicYear = store.NormalizeAcademicYear(state.AcademicYear)
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

	yearSelect := widget.NewSelect(yearOptions, func(selected string) {
		if selected != "" && selected != state.AcademicYear {
			state.SetAcademicYear(selected)
		}
	})
	yearSelect.SetSelected(state.AcademicYear)

	// Add Custom Year Button
	addYearBtn := widget.NewButtonWithIcon("Yeni Yıl", theme.ContentAddIcon(), func() {
		entry := widget.NewEntry()
		entry.SetPlaceHolder("Örn: 2029 - 2030")
		dialog.ShowCustomConfirm("Yeni Öğretim Yılı Ekle", "Ekle", "İptal", entry, func(ok bool) {
			if ok && strings.TrimSpace(entry.Text) != "" {
				norm := store.NormalizeAcademicYear(entry.Text)
				state.SetAcademicYear(norm)
			}
		}, state.Window)
	})

	// School Type Radio
	schoolTypeRadio := widget.NewRadioGroup([]string{"LİSE", "ORTAOKUL"}, func(selected string) {
		if selected != "" && selected != state.SchoolInfo.SchoolType {
			state.SetSchoolType(selected)
		}
	})
	schoolTypeRadio.Horizontal = true
	schoolTypeRadio.SetSelected(state.SchoolInfo.SchoolType)

	// Open Output Folder Button
	openFolderBtn := widget.NewButtonWithIcon("Çıktı Klasörü", theme.FolderOpenIcon(), func() {
		home, err := os.UserHomeDir()
		outDir := "."
		if err == nil {
			outDir = filepath.Join(home, "OkulEvrak_Ciktilar")
		}
		_ = os.MkdirAll(outDir, 0755)
		_ = pdf.OpenPDF(outDir)
	})

	// Settings Button
	settingsBtn := widget.NewButtonWithIcon("Ayarlar & Tema", theme.SettingsIcon(), func() {
		state.NavigateTo("settings")
	})

	headerBox := container.NewHBox(
		title,
		widget.NewSeparator(),
		widget.NewLabelWithStyle("Öğretim Yılı:", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		yearSelect,
		addYearBtn,
		widget.NewSeparator(),
		widget.NewLabelWithStyle("Kurum:", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		schoolTypeRadio,
		widget.NewSeparator(),
		openFolderBtn,
		widget.NewSeparator(),
		settingsBtn,
	)

	return container.NewVBox(
		headerBox,
		widget.NewSeparator(),
	)
}

func showInfo(w fyne.Window, title, msg string) {
	dialog.ShowInformation(title, msg, w)
}

func showError(w fyne.Window, err error) {
	if err != nil {
		dialog.ShowError(err, w)
	}
}

func showConfirm(w fyne.Window, title, msg string, onConfirm func()) {
	dialog.ShowConfirm(title, msg, func(ok bool) {
		if ok {
			onConfirm()
		}
	}, w)
}

func formatPoints(p int) string {
	return fmt.Sprintf("%d", p)
}
