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

	"okulevrak/config"
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

	yearSelect := NewTurkishSelect(yearOptions, func(selected string) {
		if selected != "" && selected != state.AcademicYear {
			state.SetAcademicYear(selected)
		}
	})
	yearSelect.SetSelected(state.AcademicYear)

	// Add Custom Year Button with auto-filled next year
	addYearBtn := widget.NewButtonWithIcon("Yeni Yıl", theme.ContentAddIcon(), func() {
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

	// Delete Active Year Button with confirmation
	deleteYearBtn := widget.NewButtonWithIcon("Yılı Sil", theme.DeleteIcon(), func() {
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

	// School Type Radio
	schoolTypeRadio := widget.NewRadioGroup([]string{"LİSE", "ORTAOKUL"}, func(selected string) {
		if selected != "" && selected != state.SchoolInfo.SchoolType {
			state.SetSchoolType(selected)
		}
	})
	schoolTypeRadio.Horizontal = true
	schoolTypeRadio.SetSelected(state.SchoolInfo.SchoolType)

	// Open Output Folder Button
	openFolderBtn := widget.NewButtonWithIcon("Belgeler Klasörü", theme.FolderOpenIcon(), func() {
		home, err := os.UserHomeDir()
		outDir := "."
		if err == nil {
			outDir = filepath.Join(home, "OkulEvrak_belgeler")
			oldDir := filepath.Join(home, "OkulEvrak_Ciktilar")
			if _, errOld := os.Stat(oldDir); errOld == nil {
				if _, errNew := os.Stat(outDir); os.IsNotExist(errNew) {
					_ = os.Rename(oldDir, outDir)
				}
			}
		}
		_ = os.MkdirAll(outDir, 0755)
		_ = pdf.OpenPDF(outDir)
	})

	// Backup Database Button
	backupBtn := widget.NewButtonWithIcon("Yedekleme", theme.DownloadIcon(), func() {
		BackupDatabase(state)
	})

	// Regulations (Yönetmelik) Button in Top Bar
	yonetmelikBtn := widget.NewButtonWithIcon("Yönetmelik & Şablon", theme.DocumentIcon(), func() {
		state.NavigateTo("yonetmelik")
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
		deleteYearBtn,
		widget.NewSeparator(),
		widget.NewLabelWithStyle("Kurum:", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		schoolTypeRadio,
		widget.NewSeparator(),
		openFolderBtn,
		backupBtn,
		widget.NewSeparator(),
		yonetmelikBtn,
		settingsBtn,
	)

	return container.NewVBox(
		headerBox,
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
