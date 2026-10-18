package ui

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	fyneDialog "fyne.io/fyne/v2/dialog"
	"github.com/sqweek/dialog"

	"okulevrak/config"
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
		destPath, err := dialog.File().
			Title("Veritabanı Yedeği Kaydet").
			SetStartFile(defaultName).
			Filter("OkulEvrak Veritabanı (*.db)", "db").
			Save()

		if err != nil {
			if errors.Is(err, dialog.ErrCancelled) {
				return
			}
			showError(state.Window, fmt.Errorf("Yedekleme konumu seçilemedi: %w", err))
			return
		}

		if destPath == "" {
			return
		}

		if filepath.Ext(destPath) != ".db" {
			destPath += ".db"
		}

		// Perform safe file copy
		if err := copyFile(srcPath, destPath); err != nil {
			showError(state.Window, fmt.Errorf("Yedekleme dosyası kaydedilemedi: %w", err))
			return
		}

		fyneDialog.ShowInformation(
			"Yedekleme Başarılı",
			fmt.Sprintf("Veritabanı yedeğiniz başarıyla oluşturuldu:\n\n%s", destPath),
			state.Window,
		)
	}()
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

		backupPath, err := dialog.File().
			Title("Geri Yüklenecek Yedek Dosyasını Seçin (.db)").
			Filter("OkulEvrak Veritabanı (*.db)", "db").
			Load()

		if err != nil {
			if errors.Is(err, dialog.ErrCancelled) {
				return
			}
			showError(state.Window, fmt.Errorf("Yedek dosyası seçilemedi: %w", err))
			return
		}

		if backupPath == "" {
			return
		}

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
	}()
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
