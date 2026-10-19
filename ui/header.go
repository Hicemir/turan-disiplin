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

	newYearControls := func() (fyne.CanvasObject, fyne.CanvasObject, fyne.CanvasObject) {
		yearSelect := NewTurkishSelect(yearOptions, func(selected string) {
			if selected != "" && selected != state.AcademicYear {
				state.SetAcademicYear(selected)
			}
		})
		yearSelect.SetSelected(state.AcademicYear)
		yearSelectContainer := container.NewGridWrap(fyne.NewSize(140, 36), yearSelect)

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

		return yearSelectContainer, addYearBtn, deleteYearBtn
	}

	// buildBar creates the header row; compact=true shows icon-only buttons
	// and the short "O.OKUL" label.
	buildBar := func(compact bool) fyne.CanvasObject {
		var liseBtn, ortaokulBtn *widget.Button
		ortaokulLabel := "ORTAOKUL"
		if compact {
			ortaokulLabel = "O.OKUL"
		}

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

		ortaokulBtn = widget.NewButton(ortaokulLabel, func() {
			if state.SchoolInfo.SchoolType != "ORTAOKUL" {
				state.SetSchoolType("ORTAOKUL")
				updateSchoolTypeButtons()
			}
		})

		updateSchoolTypeButtons()
		schoolTypeBox := container.NewHBox(liseBtn, ortaokulBtn)

		// iconBtn drops its text when compact (icon-only)
		iconBtn := func(label string, icon fyne.Resource, fn func()) *widget.Button {
			if compact {
				label = ""
			}
			return widget.NewButtonWithIcon(label, icon, fn)
		}

		// Open Output Folder Button
		openFolderBtn := iconBtn("Belgeler Klasörü", theme.FolderOpenIcon(), func() {
			home, err := os.UserHomeDir()
			outDir := "."
			if err == nil {
				outDir = filepath.Join(home, "TuranDisiplin_belgeler")
			}
			_ = os.MkdirAll(outDir, 0755)
			_ = pdf.OpenPDF(outDir)
		})

		// Backup Database Button
		backupBtn := iconBtn("Yedekleme", theme.DownloadIcon(), func() {
			BackupDatabase(state)
		})

		// Regulations (Yönetmelik) Button in Top Bar
		yonetmelikBtn := iconBtn("Yönetmelik & Şablon", theme.DocumentIcon(), func() {
			state.NavigateTo("yonetmelik")
		})

		// Settings Button
		settingsBtn := iconBtn("Ayarlar & Tema", theme.SettingsIcon(), func() {
			state.NavigateTo("settings")
		})

		barTitle := widget.NewLabelWithStyle("🏛️ TURANDİSİPLİN", fyne.TextAlignLeading, fyne.TextStyle{Bold: true})
		yearCtl, addCtl, delCtl := newYearControls()
		return container.NewHBox(
			barTitle,
			widget.NewSeparator(),
			widget.NewLabelWithStyle("Öğretim Yılı:", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
			yearCtl,
			addCtl,
			delCtl,
			widget.NewSeparator(),
			widget.NewLabelWithStyle("Kurum:", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
			schoolTypeBox,
			widget.NewSeparator(),
			openFolderBtn,
			backupBtn,
			widget.NewSeparator(),
			yonetmelikBtn,
			settingsBtn,
		)
	}

	fullBar := buildBar(false)
	compactBar := container.NewHScroll(buildBar(true))

	return container.NewVBox(
		container.New(&responsiveHeaderLayout{}, fullBar, compactBar),
		widget.NewSeparator(),
	)
}

// responsiveHeaderLayout shows objects[0] (full, with text) when it fits the
// available width and otherwise objects[1] (compact icon-only, horizontally
// scrollable as a last resort).
type responsiveHeaderLayout struct{}

func (l *responsiveHeaderLayout) Layout(objects []fyne.CanvasObject, size fyne.Size) {
	if len(objects) < 2 {
		return
	}
	full, compact := objects[0], objects[1]
	if full.MinSize().Width <= size.Width {
		compact.Hide()
		full.Show()
		full.Move(fyne.NewPos(0, 0))
		full.Resize(size)
		return
	}
	full.Hide()
	compact.Show()
	compact.Move(fyne.NewPos(0, 0))
	compact.Resize(size)
}

func (l *responsiveHeaderLayout) MinSize(objects []fyne.CanvasObject) fyne.Size {
	if len(objects) < 2 {
		return fyne.NewSize(0, 0)
	}
	// Height of the full bar, but allow the width to shrink so the window can
	// be resized small (the compact variant scrolls when needed).
	return fyne.NewSize(200, objects[0].MinSize().Height)
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
