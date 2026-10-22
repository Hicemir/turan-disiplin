package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"

	"turandisiplin/assets"
	"turandisiplin/config"
	"turandisiplin/ui"
)

func main() {
	ensureLinuxDesktopIntegration()

	a := app.NewWithID("com.turandisiplin.app")
	a.SetIcon(assets.AppIcon)

	w := a.NewWindow("TuranDisiplin " + config.AppVersion + " - MEB Öğrenci Disiplin ve Davranış Değerlendirme Sistemi")
	w.SetIcon(assets.AppIcon)
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

func ensureLinuxDesktopIntegration() {
	if runtime.GOOS != "linux" {
		return
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return
	}

	execPath, err := os.Executable()
	if err != nil {
		execPath = "TuranDisiplin-Linux"
	}

	// 1. Install hicolor and pixmaps icons (Shield emblem)
	iconDirs := []string{
		filepath.Join(home, ".local/share/icons/hicolor/512x512/apps"),
		filepath.Join(home, ".local/share/icons/hicolor/256x256/apps"),
		filepath.Join(home, ".local/share/icons/hicolor/48x48/apps"),
		filepath.Join(home, ".local/share/pixmaps"),
		filepath.Join(home, ".local/share/icons"),
	}
	iconNames := []string{"com.turandisiplin.app.png", "TuranDisiplin-Linux.png", "turandisiplin.png"}

	for _, dir := range iconDirs {
		_ = os.MkdirAll(dir, 0755)
		for _, name := range iconNames {
			p := filepath.Join(dir, name)
			if len(assets.LogoPNG) > 0 {
				_ = os.WriteFile(p, assets.LogoPNG, 0644)
			}
		}
	}

	// 2. Install desktop files for Wayland and X11 matching
	appDir := filepath.Join(home, ".local/share/applications")
	_ = os.MkdirAll(appDir, 0755)

	desktopContent := fmt.Sprintf(`[Desktop Entry]
Version=1.0
Type=Application
Name=TuranDisiplin
GenericName=MEB Öğrenci Disiplin ve Davranış Değerlendirme Sistemi
Comment=MEB Öğrenci Disiplin ve Davranış Değerlendirme Sistemi
Exec=%s
Icon=com.turandisiplin.app
Terminal=false
Categories=Office;Education;
StartupWMClass=com.turandisiplin.app
`, execPath)

	desktopFiles := []string{
		filepath.Join(appDir, "com.turandisiplin.app.desktop"),
		filepath.Join(appDir, "TuranDisiplin-Linux.desktop"),
		filepath.Join(appDir, "turandisiplin.desktop"),
	}

	for _, df := range desktopFiles {
		_ = os.WriteFile(df, []byte(desktopContent), 0644)
	}

	// 3. Show the app icon on the executable file itself in the file manager
	// (Nautilus/Nemo/Caja/Thunar honour the gio "custom-icon" metadata). A raw ELF
	// binary cannot embed an icon, so this is what replaces the grey generic icon.
	if execPath != "" && len(assets.LogoPNG) > 0 {
		iconPath := filepath.Join(home, ".local/share/icons/com.turandisiplin.app.png")
		if _, statErr := os.Stat(iconPath); statErr == nil {
			_ = exec.Command("gio", "set", execPath, "metadata::custom-icon", "file://"+iconPath).Run()
		}
	}

	// Notify desktop database and icon cache
	_ = exec.Command("update-desktop-database", appDir).Run()
	_ = exec.Command("gtk-update-icon-cache", filepath.Join(home, ".local/share/icons/hicolor")).Run()
}
