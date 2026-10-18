package ui

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"fyne.io/fyne/v2"
	fyneDialog "fyne.io/fyne/v2/dialog"

	"okulevrak/config"
	"okulevrak/picker"
	"okulevrak/store"
)

// BackupDatabase copies the active OkulEvrak.db SQLite database to a user-chosen location
func BackupDatabase(state *AppState) {
	go func() {
		cfgDir, err := config.GetConfigDir()
		if err != nil {
			showError(state.Window, fmt.Errorf("Uygulama klasörüne erişilemedi: %w", err))
			return
		}
		srcPath := filepath.Join(cfgDir, "OkulEvrak.db")

		defaultName := fmt.Sprintf("OkulEvrak_Yedek_%s.db", time.Now().Format("2006-01-02_15-04"))
		filters := []picker.FileFilter{
			{Name: "OkulEvrak Veritabanı (*.db)", Extensions: []string{"db"}},
		}

		destPath, err := picker.SaveFile("Veritabanı Yedeği Kaydet", defaultName, filters)

		if err != nil {
			if errors.Is(err, picker.ErrCancelled) {
				return
			}
			// If native tool unavailable, fallback to Fyne dialog
			fyneDialog.ShowFileSave(func(writer fyne.URIWriteCloser, err error) {
				if err != nil || writer == nil {
					return
				}
				p := writer.URI().Path()
				_ = writer.Close()
				if p != "" {
					executeBackup(state, srcPath, p)
				}
			}, state.Window)
			return
		}

		if destPath == "" {
			return
		}

		if filepath.Ext(destPath) != ".db" {
			destPath += ".db"
		}

		executeBackup(state, srcPath, destPath)
	}()
}

func executeBackup(state *AppState, srcPath, destPath string) {
	if err := copyFile(srcPath, destPath); err != nil {
		showError(state.Window, fmt.Errorf("Yedekleme dosyası kaydedilemedi: %w", err))
		return
	}

	fyneDialog.ShowInformation(
		"Yedekleme Başarılı",
		fmt.Sprintf("Veritabanı yedeğiniz başarıyla oluşturuldu:\n\n%s", destPath),
		state.Window,
	)
}

// RestoreDatabase restores OkulEvrak.db from a chosen .db backup file
func RestoreDatabase(state *AppState) {
	go func() {
		cfgDir, err := config.GetConfigDir()
		if err != nil {
			showError(state.Window, fmt.Errorf("Uygulama klasörüne erişilemedi: %w", err))
			return
		}
		targetPath := filepath.Join(cfgDir, "OkulEvrak.db")

		filters := []picker.FileFilter{
			{Name: "OkulEvrak Veritabanı (*.db)", Extensions: []string{"db"}},
		}

		backupPath, err := picker.OpenFile("Geri Yüklenecek Yedek Dosyasını Seçin (.db)", filters)

		if err != nil {
			if errors.Is(err, picker.ErrCancelled) {
				return
			}
			// If native tool unavailable, fallback to Fyne dialog
			fyneDialog.ShowFileOpen(func(reader fyne.URIReadCloser, err error) {
				if err != nil || reader == nil {
					return
				}
				p := reader.URI().Path()
				_ = reader.Close()
				if p != "" {
					executeRestoreConfirmation(state, targetPath, p)
				}
			}, state.Window)
			return
		}

		if backupPath == "" {
			return
		}

		executeRestoreConfirmation(state, targetPath, backupPath)
	}()
}

func executeRestoreConfirmation(state *AppState, targetPath, backupPath string) {
	// Prompt user with clear warning and confirmation
	fyneDialog.ShowConfirm(
		"Geri Yükleme Onayı",
		fmt.Sprintf(
			"DİKKAT: Seçtiğiniz '%s' yedek dosyası mevcut veritabanınızın üzerine yazılacaktır.\n\n"+
				"Bu işlem geri alınamaz! Mevcut veritabanınızı bu yedekle değiştirmek istediğinize emin misiniz?",
			filepath.Base(backupPath),
		),
		func(ok bool) {
			if !ok {
				return
			}

			if err := copyFile(backupPath, targetPath); err != nil {
				showError(state.Window, fmt.Errorf("Geri yükleme başarısız oldu: %w", err))
				return
			}

			// Re-read school info and refresh application
			st, err := store.GetStore()
			if err == nil && st != nil {
				if info, err := st.GetSchoolInfo(); err == nil && info != nil {
					state.SchoolInfo = info
				}
			}
			state.NotifyChange()

			fyneDialog.ShowInformation(
				"Geri Yükleme Başarılı",
				"Veritabanınız başarıyla geri yüklendi ve güncellendi.",
				state.Window,
			)
		},
		state.Window,
	)
}

func copyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	if _, err = io.Copy(out, in); err != nil {
		return err
	}
	return out.Sync()
}
