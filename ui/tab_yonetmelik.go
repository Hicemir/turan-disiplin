package ui

import (
	"fmt"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"okulevrak/config"
	"okulevrak/model"
)

// BuildYonetmelikTab creates the In-App Regulation & Template Editor UI
func BuildYonetmelikTab(state *AppState) fyne.CanvasObject {
	var currentClauses []model.Clause
	var selectedIndex int = -1

	// Search filter
	searchEntry := widget.NewEntry()
	searchEntry.SetPlaceHolder("Madde no veya suç açıklaması ile ara...")

	schoolTypeFilter := widget.NewRadioGroup([]string{"TÜMÜ", "LİSE", "ORTAOKUL"}, nil)
	schoolTypeFilter.Horizontal = true
	schoolTypeFilter.SetSelected("TÜMÜ")

	countLabel := widget.NewLabelWithStyle("Kayıt: 0 madde", fyne.TextAlignLeading, fyne.TextStyle{Bold: true})

	// Form Inputs
	schoolTypeSelect := widget.NewSelect([]string{"LİSE", "ORTAOKUL"}, nil)
	schoolTypeSelect.SetSelected(state.SchoolInfo.SchoolType)

	articleEntry := widget.NewEntry()
	articleEntry.SetPlaceHolder("Örn: Madde 164/1-a")

	penaltySelect := widget.NewSelect([]string{
		"Kınama",
		"Okuldan Kısa Süreli Uzaklaştırma",
		"Okul Değiştirme",
		"Örgün Eğitim Dışına Çıkarma",
		"Uyarma",
		"Okul Değiştirme Yaptırımı",
	}, nil)

	descEntry := widget.NewMultiLineEntry()
	descEntry.SetPlaceHolder("Yönetmelikte belirtilen eylem ve suç tanımı...")

	clearForm := func() {
		selectedIndex = -1
		schoolTypeSelect.SetSelected(state.SchoolInfo.SchoolType)
		articleEntry.SetText("")
		penaltySelect.ClearSelected()
		descEntry.SetText("")
	}

	populateForm := func(c model.Clause, idx int) {
		selectedIndex = idx
		schoolTypeSelect.SetSelected(c.SchoolType)
		articleEntry.SetText(c.Article)
		penaltySelect.SetSelected(c.Penalty)
		descEntry.SetText(c.Description)
	}

	// Clauses List
	clauseList := widget.NewList(
		func() int {
			return len(currentClauses)
		},
		func() fyne.CanvasObject {
			return container.NewHBox(
				widget.NewLabelWithStyle("[LİSE]", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
				widget.NewSeparator(),
				widget.NewLabelWithStyle("Madde", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
				widget.NewSeparator(),
				widget.NewLabelWithStyle("Ceza", fyne.TextAlignLeading, fyne.TextStyle{Italic: true}),
				widget.NewSeparator(),
				widget.NewLabel("Açıklama özeti"),
			)
		},
		func(i widget.ListItemID, o fyne.CanvasObject) {
			if i >= len(currentClauses) {
				return
			}
			c := currentClauses[i]
			box := o.(*fyne.Container)
			typeLbl := box.Objects[0].(*widget.Label)
			artLbl := box.Objects[2].(*widget.Label)
			penLbl := box.Objects[4].(*widget.Label)
			descLbl := box.Objects[6].(*widget.Label)

			sType := c.SchoolType
			if sType == "" {
				sType = "GENEL"
			}
			typeLbl.SetText(fmt.Sprintf("[%s]", sType))
			artLbl.SetText(c.Article)
			penLbl.SetText(c.Penalty)

			shortDesc := c.Description
			if len(shortDesc) > 60 {
				shortDesc = shortDesc[:57] + "..."
			}
			descLbl.SetText(shortDesc)
		},
	)

	reloadList := func() {
		filterType := schoolTypeFilter.Selected
		searchTerm := strings.ToLower(searchEntry.Text)

		var filtered []model.Clause
		for _, c := range state.Yonetmelik.Clauses {
			if filterType != "TÜMÜ" && !strings.EqualFold(c.SchoolType, filterType) {
				continue
			}
			if searchTerm != "" && !strings.Contains(strings.ToLower(c.Article), searchTerm) && !strings.Contains(strings.ToLower(c.Description), searchTerm) {
				continue
			}
			filtered = append(filtered, c)
		}
		currentClauses = filtered
		countLabel.SetText(fmt.Sprintf("Kayıt: %d madde", len(currentClauses)))
		clauseList.Refresh()
	}

	clauseList.OnSelected = func(id widget.ListItemID) {
		if id >= 0 && id < len(currentClauses) {
			c := currentClauses[id]
			// find original index in state.Yonetmelik.Clauses
			origIdx := -1
			for i, orig := range state.Yonetmelik.Clauses {
				if orig.Article == c.Article && orig.Description == c.Description && orig.SchoolType == c.SchoolType {
					origIdx = i
					break
				}
			}
			populateForm(c, origIdx)
		}
	}

	searchEntry.OnChanged = func(s string) { reloadList() }
	schoolTypeFilter.OnChanged = func(s string) { reloadList() }

	// Save Clause Button (Updates existing or Appends new, saves to user config file)
	saveClauseBtn := widget.NewButtonWithIcon("Maddeyi Kaydet", theme.DocumentSaveIcon(), func() {
		art := strings.TrimSpace(articleEntry.Text)
		desc := strings.TrimSpace(descEntry.Text)
		pen := penaltySelect.Selected
		sType := schoolTypeSelect.Selected

		if art == "" || desc == "" || pen == "" {
			showError(state.Window, fmt.Errorf("Lütfen Madde No, Ceza Türü ve Açıklama alanlarını doldurunuz"))
			return
		}

		newClause := model.Clause{
			Penalty:     pen,
			Article:     art,
			Description: desc,
			SchoolType:  sType,
		}

		if selectedIndex >= 0 && selectedIndex < len(state.Yonetmelik.Clauses) {
			// Update existing
			state.Yonetmelik.Clauses[selectedIndex] = newClause
		} else {
			// Add new
			state.Yonetmelik.Clauses = append(state.Yonetmelik.Clauses, newClause)
		}

		err := config.SaveYonetmelik(state.Yonetmelik)
		if err != nil {
			showError(state.Window, fmt.Errorf("Yönetmelik kaydedilemedi: %w", err))
			return
		}

		showInfo(state.Window, "Başarılı", "Yönetmelik maddesi kaydedildi ve ayarlarınıza işlendi.")
		clearForm()
		reloadList()
		state.NotifyChange()
	})
	saveClauseBtn.Importance = widget.HighImportance

	// Delete Clause Button
	deleteClauseBtn := widget.NewButtonWithIcon("Seçilen Maddeyi Sil", theme.DeleteIcon(), func() {
		if selectedIndex < 0 || selectedIndex >= len(state.Yonetmelik.Clauses) {
			showError(state.Window, fmt.Errorf("Silinecek bir madde seçilmedi"))
			return
		}

		showConfirm(state.Window, "Maddeyi Sil", "Bu yönetmelik maddesini silmek istediğinize emin misiniz?", func() {
			state.Yonetmelik.Clauses = append(state.Yonetmelik.Clauses[:selectedIndex], state.Yonetmelik.Clauses[selectedIndex+1:]...)
			_ = config.SaveYonetmelik(state.Yonetmelik)
			clearForm()
			reloadList()
			state.NotifyChange()
		})
	})
	deleteClauseBtn.Importance = widget.DangerImportance

	newClauseBtn := widget.NewButtonWithIcon("Yeni Madde", theme.ContentAddIcon(), func() {
		clearForm()
	})

	// Reset to Factory Defaults Button
	resetDefaultsBtn := widget.NewButtonWithIcon("Varsayılana Dön (Fabrika Ayarları)", theme.ViewRefreshIcon(), func() {
		showConfirm(state.Window, "Yönetmeliği Sıfırla", "Tüm özel eklediğiniz/düzenlediğiniz yönetmelik maddeleri silinecek ve orijinal sisteme gömülü yönetmelik yüklenecektir. Devam edilsin mi?", func() {
			ym, err := config.ResetYonetmelik()
			if err != nil {
				showError(state.Window, fmt.Errorf("Sıfırlama başarısız: %w", err))
				return
			}
			state.Yonetmelik = ym
			showInfo(state.Window, "Sıfırlandı", "Orijinal yönetmelik başarıyla geri yüklendi.")
			clearForm()
			reloadList()
			state.NotifyChange()
		})
	})

	// Layout Structure
	filterRow := container.NewHBox(
		countLabel,
		widget.NewSeparator(),
		schoolTypeFilter,
		widget.NewSeparator(),
		resetDefaultsBtn,
	)

	leftPanel := container.NewBorder(
		container.NewVBox(searchEntry, filterRow),
		nil, nil, nil,
		clauseList,
	)

	formItems := widget.NewForm(
		widget.NewFormItem("Okul Kademesi *", schoolTypeSelect),
		widget.NewFormItem("Madde Numarası *", articleEntry),
		widget.NewFormItem("Verilecek Ceza / Yaptırım *", penaltySelect),
		widget.NewFormItem("Suç / Eylem Açıklaması *", descEntry),
	)

	cardEdit := widget.NewCard("Madde Düzenleme / Yeni Madde Ekleme", "Listeden bir madde seçerek metnini düzenleyebilir veya yeni ceza maddeleri ekleyebilirsiniz.", formItems)
	buttonsRow := container.NewHBox(saveClauseBtn, deleteClauseBtn, newClauseBtn)

	rightPanel := container.NewVScroll(container.NewVBox(
		cardEdit,
		buttonsRow,
		widget.NewSeparator(),
		widget.NewLabelWithStyle("💡 Bilgi: Yaptığınız tüm değişiklikler bilgisayarınızdaki OkulEvrak profil klasörüne kaydedilir ve bir sonraki açılışta otomatik yüklenir. 'Varsayılana Dön' ile istediğiniz an fabrika ayarlarına dönebilirsiniz.", fyne.TextAlignLeading, fyne.TextStyle{Italic: true}),
	))

	split := container.NewHSplit(leftPanel, rightPanel)
	split.SetOffset(0.55)

	reloadList()
	return split
}
