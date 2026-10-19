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

	"turandisiplin/config"
	"turandisiplin/pdf"
	"turandisiplin/store"
)

// BuildHeader creates the persistent senior-friendly top bar
func BuildHeader(state *AppState) fyne.CanvasObject {
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

	title := widget.NewLabelWithStyle("🏛️ TURANDİSİPLİN", fyne.TextAlignLeading, fyne.TextStyle{Bold: true})

	yearSelect := NewTurkishSelect(yearOptions, func(selected string) {
		if selected != "" && selected != state.AcademicYear {
			state.SetAcademicYear(selected)
		}
	})
	yearSelect.SetSelected(state.AcademicYear)
	yearSelectContainer := container.NewGridWrap(fyne.NewSize(130, 36), yearSelect)

	// Add Custom Year Button (Compact + icon)
	addYearBtn := widget.NewButtonWithIcon("", theme.ContentAddIcon(), func() {
		entry := NewSingleLineEntry()
		nextYear := store.GetNextNearestAcademicYear(yearOptions)
		entry.SetText(nextYear)
		dialog.ShowCustomConfirm("Yeni Öğretim Yılı Ekle", "Ekle", "İptal", entry, func(ok bool) {
			if ok && strings.TrimSpace(entry.Text) != "" {
				norm := store.NormalizeAcademicYear(entry.Text)
				state.SetAcademicYear(norm)
			}
		}, state.Window)
	})

	// Delete Active Year Button (Compact trash icon)
	deleteYearBtn := widget.NewButtonWithIcon("", theme.DeleteIcon(), func() {
		if len(yearOptions) <= 1 {
			ShowTurkishInfo("Uyarı", "Sistemde en az bir öğretim yılı kalmalıdır. Tek kalan yılı silemezsiniz.", state.Window)
			return
		}
		ShowTurkishDangerConfirm(
			"Öğretim Yılını Sil",
			fmt.Sprintf("DİKKAT:\n'%s' öğretim yılını silmek istediğinize emin misiniz?\n\nBu yıla ait TÜM öğrenci, olay ve disiplin kararları kalıcı olarak silinecektir!\nBu işlem geri alınamaz.", state.AcademicYear),
			func(ok bool) {
				if !ok {
					return
				}
				toDelete := state.AcademicYear
				_ = state.Store.DeleteAcademicYear(toDelete)
				state.Config.DeletedYears = append(state.Config.DeletedYears, toDelete)
				_ = config.SaveConfig(state.Config)

				remaining := state.Store.GetAcademicYears()
				if len(remaining) > 0 {
					state.SetAcademicYear(remaining[0])
				}
				ShowTurkishInfo("Silindi", fmt.Sprintf("'%s' öğretim yılı başarıyla silindi.", toDelete), state.Window)
			},
			state.Window,
		)
	})

	var liseBtn, ortaokulBtn *widget.Button
	updateSchoolTypeButtons := func() {
		if state.SchoolInfo.SchoolType == "ORTAOKUL" {
			liseBtn.Importance = widget.LowImportance
			ortaokulBtn.Importance = widget.HighImportance
		} else {
			liseBtn.Importance = widget.HighImportance
			ortaokulBtn.Importance = widget.LowImportance
		}
		liseBtn.Refresh()
		ortaokulBtn.Refresh()
	}

	liseBtn = widget.NewButton("LİSE", func() {
		if state.SchoolInfo.SchoolType != "LİSE" {
			state.SetSchoolType("LİSE")
			updateSchoolTypeButtons()
		}
	})

	ortaokulBtn = widget.NewButton("ORTAOKUL", func() {
		if state.SchoolInfo.SchoolType != "ORTAOKUL" {
			state.SetSchoolType("ORTAOKUL")
			updateSchoolTypeButtons()
		}
	})

	updateSchoolTypeButtons()
	schoolTypeBox := container.NewHBox(liseBtn, ortaokulBtn)

	// Open Output Folder Button
	openFolderBtn := widget.NewButtonWithIcon("Belgeler", theme.FolderOpenIcon(), func() {
		home, err := os.UserHomeDir()
		outDir := "."
		if err == nil {
			outDir = filepath.Join(home, "TuranDisiplin_belgeler")
		}
		_ = os.MkdirAll(outDir, 0755)
		_ = pdf.OpenPDF(outDir)
	})

	// Backup Database Button
	backupBtn := widget.NewButtonWithIcon("Yedek", theme.DownloadIcon(), func() {
		BackupDatabase(state)
	})

	// Regulations (Yönetmelik) Button in Top Bar
	yonetmelikBtn := widget.NewButtonWithIcon("Yönetmelik", theme.DocumentIcon(), func() {
		state.NavigateTo("yonetmelik")
	})

	// Settings Button
	settingsBtn := widget.NewButtonWithIcon("Ayarlar", theme.SettingsIcon(), func() {
		state.NavigateTo("settings")
	})

	leftSide := container.NewHBox(
		title,
		widget.NewSeparator(),
		widget.NewLabelWithStyle("Yıl:", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		yearSelectContainer,
		addYearBtn,
		deleteYearBtn,
		widget.NewSeparator(),
		schoolTypeBox,
	)

	rightSide := container.NewHBox(
		openFolderBtn,
		backupBtn,
		widget.NewSeparator(),
		yonetmelikBtn,
		settingsBtn,
	)

	headerRow := container.NewBorder(nil, nil, leftSide, rightSide, nil)

	return container.NewVBox(
		headerRow,
		widget.NewSeparator(),
	)
}

func showInfo(w fyne.Window, title, msg string) {
	ShowTurkishInfo(title, msg, w)
}

func showError(w fyne.Window, err error) {
	if err != nil {
		dialog.ShowError(err, w)
	}
}

func showConfirm(w fyne.Window, title, msg string, onConfirm func()) {
	ShowTurkishConfirm(title, msg, func(ok bool) {
		if ok {
			onConfirm()
		}
	}, w)
}

func formatPoints(p int) string {
	return fmt.Sprintf("%d", p)
}
