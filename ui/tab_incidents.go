package ui

import (
	"fmt"
	"strings"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	fyneDialog "fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"turandisiplin/model"
)

// BuildIncidentsTab creates a clean 2-panel Incident Management UI
func BuildIncidentsTab(state *AppState) fyne.CanvasObject {
	var currentIncidents []model.Incident
	var selectedIncident *model.Incident

	var incidentStudents []model.IncidentStudent
	var selectedIncStudent *model.IncidentStudent

	// Incident Form Inputs (Single-line entries to prevent trapping mouse scroll)
	codeEntry := NewSingleLineEntry()
	titleEntry := NewSingleLineEntry()
	dateEntry := NewSingleLineEntry()
	dateEntry.SetText(time.Now().Format("02.01.2006"))
	timeEntry := NewSingleLineEntry()
	timeEntry.SetText(time.Now().Format("15:04"))
	locationEntry := NewSingleLineEntry()
	locationEntry.SetPlaceHolder("Örn: Okul Bahçesi, 10/B Sınıfı")
	complainantEntry := NewSingleLineEntry()
	complainantEntry.SetPlaceHolder("Örn: Nöbetçi Öğretmen Ahmet Yılmaz")
	petitionEntry := NewSingleLineEntry()
	summaryEntry := NewSmartMultiLineEntry()
	summaryEntry.SetPlaceHolder("Olayın detaylı açıklaması...")

	var incStudentsList *FocusableList
	var incidentList *widget.List
	var editStudentProfileBtn *widget.Button
	var removeStudentFromIncBtn *widget.Button
	var saveStudentStatementBtn *widget.Button
	var stStatementLabel *widget.Label
	var stInlineRoleSelect *widget.Select
	var stInlineStatementEntry *SmartMultiLineEntry

	clearSelectedStudent := func() {
		selectedIncStudent = nil
		state.SelectedIncStudentID = 0
		if stStatementLabel != nil {
			stStatementLabel.SetText("İfade girmek veya düzenlemek için yukarıdaki listeden bir öğrenci seçiniz.")
		}
		if stInlineStatementEntry != nil {
			stInlineStatementEntry.SetText("")
		}
		if saveStudentStatementBtn != nil {
			saveStudentStatementBtn.Disable()
		}
		if editStudentProfileBtn != nil {
			editStudentProfileBtn.Disable()
		}
		if removeStudentFromIncBtn != nil {
			removeStudentFromIncBtn.Disable()
		}
	}

	clearIncidentForm := func() {
		selectedIncident = nil
		state.SelectedIncidentCode = ""
		codeEntry.SetText(state.Store.GetNextIncidentCode(state.AcademicYear))
		titleEntry.SetText("")
		dateEntry.SetText(time.Now().Format("02.01.2006"))
		timeEntry.SetText(time.Now().Format("15:04"))
		locationEntry.SetText("")
		complainantEntry.SetText("")
		petitionEntry.SetText("")
		summaryEntry.SetText("")

		incidentStudents = nil
		clearSelectedStudent()
		if incStudentsList != nil {
			incStudentsList.Refresh()
		}
	}

	reloadIncStudents := func() {
		if selectedIncident == nil {
			incidentStudents = nil
			clearSelectedStudent()
			if incStudentsList != nil {
				incStudentsList.Refresh()
			}
			return
		}
		list, err := state.Store.GetIncidentStudents(state.AcademicYear, selectedIncident.IncidentCode)
		if err == nil {
			incidentStudents = list
			if incStudentsList != nil {
				incStudentsList.Refresh()
			}
		}
	}

	selectIncident := func(inc *model.Incident) {
		if inc == nil {
			return
		}
		isDifferent := (selectedIncident == nil || selectedIncident.IncidentCode != inc.IncidentCode)
		selectedIncident = inc
		state.SelectedIncidentCode = inc.IncidentCode

		codeEntry.SetText(inc.IncidentCode)
		titleEntry.SetText(inc.Title)
		dateEntry.SetText(inc.Date)
		timeEntry.SetText(inc.Time)
		locationEntry.SetText(inc.Location)
		complainantEntry.SetText(inc.Complainant)
		petitionEntry.SetText(inc.PetitionRef)
		summaryEntry.SetText(inc.Summary)

		if isDifferent {
			clearSelectedStudent()
			reloadIncStudents()
		}
	}

	// Attached Student List
	incStudentsList = NewFocusableList(
		func() int {
			return len(incidentStudents)
		},
		func() fyne.CanvasObject {
			return container.NewHBox(
				widget.NewLabelWithStyle("No", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
				widget.NewSeparator(),
				widget.NewLabel("Öğrenci Adı Soyadı"),
				widget.NewSeparator(),
				widget.NewLabel("Sınıf"),
				widget.NewSeparator(),
				widget.NewLabelWithStyle("Rol", fyne.TextAlignLeading, fyne.TextStyle{Italic: true}),
			)
		},
		func(i widget.ListItemID, o fyne.CanvasObject) {
			if i >= len(incidentStudents) {
				return
			}
			is := incidentStudents[i]
			box := o.(*fyne.Container)
			noLbl := box.Objects[0].(*widget.Label)
			nameLbl := box.Objects[2].(*widget.Label)
			classLbl := box.Objects[4].(*widget.Label)
			roleLbl := box.Objects[6].(*widget.Label)

			noLbl.SetText(fmt.Sprintf("%-5s", is.SchoolNo))
			nameLbl.SetText(is.FullName)
			classLbl.SetText(is.ClassSection)
			roleLbl.SetText(is.Role)
		},
		fyne.NewSize(500, 180),
	)

	// Incident List
	incidentList = widget.NewList(
		func() int {
			return len(currentIncidents)
		},
		func() fyne.CanvasObject {
			return container.NewHBox(
				widget.NewLabelWithStyle("KOD", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
				widget.NewSeparator(),
				widget.NewLabel("Olay Adı / Başlığı"),
				widget.NewSeparator(),
				widget.NewLabel("Tarih"),
			)
		},
		func(i widget.ListItemID, o fyne.CanvasObject) {
			if i >= len(currentIncidents) {
				return
			}
			inc := currentIncidents[i]
			box := o.(*fyne.Container)
			codeLbl := box.Objects[0].(*widget.Label)
			titleLbl := box.Objects[2].(*widget.Label)
			dateLbl := box.Objects[4].(*widget.Label)

			codeLbl.SetText(fmt.Sprintf("#%-3s", inc.IncidentCode))
			titleLbl.SetText(inc.Title)
			dateLbl.SetText(inc.Date)
		},
	)

	reloadIncidents := func() {
		list, err := state.Store.GetIncidents(state.AcademicYear)
		if err == nil {
			currentIncidents = list
			incidentList.Refresh()

			// Maintain active selection
			if selectedIncident != nil {
				for idx, inc := range currentIncidents {
					if inc.IncidentCode == selectedIncident.IncidentCode {
						incidentList.Select(idx)
						break
					}
				}
			}
		}
	}

	incidentList.OnSelected = func(id widget.ListItemID) {
		if id >= 0 && id < len(currentIncidents) {
			inc := currentIncidents[id]
			selectIncident(&inc)
		}
	}

	// Buttons
	saveIncidentBtn := widget.NewButtonWithIcon("Olayı Kaydet", theme.DocumentSaveIcon(), func() {
		code := strings.TrimSpace(codeEntry.Text)
		title := strings.TrimSpace(titleEntry.Text)
		if code == "" || title == "" {
			showError(state.Window, fmt.Errorf("Olay kodu ve Olay adı zorunludur"))
			return
		}

		var id int64
		if selectedIncident != nil {
			id = selectedIncident.ID
		}

		inc := &model.Incident{
			ID:           id,
			AcademicYear: state.AcademicYear,
			IncidentCode: code,
			Title:        title,
			Date:         dateEntry.Text,
			Time:         timeEntry.Text,
			Location:     locationEntry.Text,
			Complainant:  complainantEntry.Text,
			PetitionRef:  petitionEntry.Text,
			Summary:      summaryEntry.Text,
		}

		err := state.Store.SaveIncident(inc)
		if err != nil {
			showError(state.Window, err)
			return
		}

		ShowTurkishInfo("Başarılı", "Olay kaydı başarıyla kaydedildi.", state.Window)
		selectedIncident = inc
		state.SelectedIncidentCode = code
		reloadIncidents()
	})
	saveIncidentBtn.Importance = widget.HighImportance

	deleteIncidentBtn := widget.NewButtonWithIcon("Olayı Sil", theme.DeleteIcon(), func() {
		if selectedIncident == nil {
			showError(state.Window, fmt.Errorf("Silinecek olay seçilmedi"))
			return
		}
		ShowTurkishDangerConfirm(
			"Olay Sil",
			fmt.Sprintf("'%s' olayını ve bu olaya bağlı tüm öğrenci/karar kayıtlarını silmek istediğinize emin misiniz?", selectedIncident.Title),
			func(confirmed bool) {
				if !confirmed {
					return
				}
				_ = state.Store.DeleteIncident(state.AcademicYear, selectedIncident.IncidentCode)
				clearIncidentForm()
				reloadIncidents()
			},
			state.Window,
		)
	})
	deleteIncidentBtn.Importance = widget.DangerImportance

	newIncidentBtn := widget.NewButtonWithIcon("Yeni Olay", theme.ContentAddIcon(), func() {
		clearIncidentForm()
	})

	// Add Student to Incident Modal
	addStudentToIncBtn := widget.NewButtonWithIcon("Olaya Öğrenci Ekle", theme.ContentAddIcon(), func() {
		if selectedIncident == nil {
			showError(state.Window, fmt.Errorf("Lütfen önce bir olay seçin veya kaydedin"))
			return
		}

		// Searchable student picker
		stSearch := NewSingleLineEntry()
		stSearch.SetPlaceHolder("Öğrenci ara (No veya İsim)...")
		var foundList []model.Student

		pickList := widget.NewList(
			func() int { return len(foundList) },
			func() fyne.CanvasObject {
				return container.NewHBox(
					widget.NewLabelWithStyle("00000", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
					widget.NewSeparator(),
					widget.NewLabel("Ad Soyad"),
					widget.NewSeparator(),
					widget.NewLabel("Sınıf"),
				)
			},
			func(i widget.ListItemID, o fyne.CanvasObject) {
				if i >= len(foundList) {
					return
				}
				s := foundList[i]
				box := o.(*fyne.Container)
				box.Objects[0].(*widget.Label).SetText(s.SchoolNo)
				box.Objects[2].(*widget.Label).SetText(s.FullName)
				box.Objects[4].(*widget.Label).SetText(s.ClassSection)
			},
		)

		updatePickList := func() {
			res, _ := state.Store.GetStudents(state.AcademicYear, stSearch.Text, "")
			foundList = res
			pickList.Refresh()
		}

		stSearch.OnChanged = func(s string) {
			updatePickList()
		}

		var chosen *model.Student
		pickList.OnSelected = func(id widget.ListItemID) {
			if id >= 0 && id < len(foundList) {
				s := foundList[id]
				chosen = &s
			}
		}

		pickerContent := container.NewBorder(stSearch, nil, nil, nil, pickList)
		pickerContent.Resize(fyne.NewSize(500, 350))

		d := fyneDialog.NewCustomConfirm("Olaya Öğrenci Ekle", "Ekle", "İptal", pickerContent, func(ok bool) {
			if !ok || chosen == nil {
				return
			}

			// Prior penalty check
			prior := state.Store.GetPriorPenalties(chosen.FullName)
			if prior == "" {
				prior = "DAHA ÖNCE CEZA ALMAMIŞTIR"
			}

			is := &model.IncidentStudent{
				AcademicYear:        state.AcademicYear,
				IncidentCode:        selectedIncident.IncidentCode,
				IncidentTitle:       selectedIncident.Title,
				Role:                "SUÇLANAN ÖĞRENCİ",
				SchoolNo:            chosen.SchoolNo,
				FullName:            chosen.FullName,
				ClassSection:        chosen.ClassSection,
				Gender:              chosen.Gender,
				TCNo:                chosen.TCNo,
				BirthPlace:          chosen.BirthPlace,
				BirthDate:           chosen.BirthDate,
				FatherName:          chosen.FatherName,
				MotherName:          chosen.MotherName,
				GuardianName:        model.EffectiveGuardianName(chosen.FullName, chosen.FatherName, chosen.GuardianName),
				GuardianPhone:       chosen.GuardianPhone,
				GuardianAddr:        chosen.GuardianAddr,
				BoardingStatus:      chosen.BoardingStatus,
				AcademicStanding:    "ORTA",
				HealthStatus:        "SAĞLAM",
				EconomicStatus:      "ORTA",
				LivingWith:          "AİLESİYLE BİRLİKTE",
				ParentsLivingStatus: "HER İKİSİ DE SAĞ",
				ParentsBioStatus:    "ÖZ VE BİRLİKTE YAŞIYORLAR",
				EnvironmentRaised:   "İL / İLÇE MERKEZİ",
				PriorPenaltyStatus:  prior,
			}

			_ = state.Store.SaveIncidentStudent(is)
			reloadIncStudents()

			// Ensure incident stays actively selected in incidentList
			if selectedIncident != nil {
				for idx, inc := range currentIncidents {
					if inc.IncidentCode == selectedIncident.IncidentCode {
						incidentList.Select(idx)
						break
					}
				}
			}

			// Automatically select the newly added student in incStudentsList
			for idx, st := range incidentStudents {
				if st.SchoolNo == chosen.SchoolNo {
					incStudentsList.Select(idx)
					break
				}
			}

			ShowTurkishInfo("Öğrenci Eklendi", fmt.Sprintf("%s (%s) başarıyla olaya eklendi.", chosen.FullName, chosen.SchoolNo), state.Window)
		}, state.Window)

		d.Resize(fyne.NewSize(520, 420))
		updatePickList()
		d.Show()
	})
	addStudentToIncBtn.Importance = widget.MediumImportance

	// Edit Student Profile & Statement Modal
	showStudentProfileModal := func(is *model.IncidentStudent) {
		if is == nil {
			return
		}

		stRoleSelect := NewTurkishSelect([]string{
			"SUÇLANAN ÖĞRENCİ",
			"MAĞDUR / ŞİKAYETÇİ ÖĞRENCİ",
			"BİLGİSİNE BAŞVURULACAK ÖĞRENCİ (TANIK)",
		}, nil)
		stRoleSelect.SetSelected(is.Role)

		stAcademicSelect := NewTurkishSelect([]string{"ÇOK İYİ", "İYİ", "ORTA", "ZAYIF"}, nil)
		stAcademicSelect.SetSelected(is.AcademicStanding)

		stHealthSelect := NewTurkishSelect([]string{"SAĞLAM", "KRONİK RAHATSIZLIĞI VAR", "ENGELLİ", "DİĞER"}, nil)
		stHealthSelect.SetSelected(is.HealthStatus)

		stEconomicSelect := NewTurkishSelect([]string{"İYİ", "ORTA", "DAR GELİRLİ", "ÇOK DÜŞÜK"}, nil)
		stEconomicSelect.SetSelected(is.EconomicStatus)

		stLivingSelect := NewTurkishSelect([]string{"AİLESİYLE BİRLİKTE", "YALNIZ ANNESİYLE", "YALNIZ BABASIYLA", "AKRABALARIYLA", "YURTTA / PANSİYONDA", "YALNIZ"}, nil)
		stLivingSelect.SetSelected(is.LivingWith)

		stParentsLivingSelect := NewTurkishSelect([]string{"HER İKİSİ DE SAĞ", "BABA VEFAT ETMİŞ", "ANNE VEFAT ETMİŞ", "HER İKİSİ DE VEFAT ETMİŞ"}, nil)
		stParentsLivingSelect.SetSelected(is.ParentsLivingStatus)

		stParentsBioSelect := NewTurkishSelect([]string{"ÖZ VE BİRLİKTE YAŞIYORLAR", "BOŞANMIŞLAR / AYRI YAŞIYORLAR", "ÜVEY ANNE", "ÜVEY BABA"}, nil)
		stParentsBioSelect.SetSelected(is.ParentsBioStatus)

		stEnvSelect := NewTurkishSelect([]string{"BÜYÜKŞEHİR", "İL / İLÇE MERKEZİ", "KASABA / BELDE", "KÖY"}, nil)
		stEnvSelect.SetSelected(is.EnvironmentRaised)

		stPriorPenaltySelect := NewTurkishSelect([]string{"DAHA ÖNCE CEZA ALMAMIŞTIR", "UYARMA CEZASI ALMIŞTIR", "KINAMA CEZASI ALMIŞTIR", "OKULDAN KISA SÜRELİ UZAKLAŞTIRMA ALMIŞTIR"}, nil)
		stPriorPenaltySelect.SetSelected(is.PriorPenaltyStatus)

		stStatementEntry := NewSmartMultiLineEntry()
		stStatementEntry.SetText(is.StatementText)
		stStatementEntry.SetPlaceHolder("Öğrencinin olayla ilgili sözlü / yazılı ifadesi ve savunması...")

		form := widget.NewForm(
			widget.NewFormItem("Olaydaki Rolü *", stRoleSelect),
			widget.NewFormItem("Başarı Durumu", stAcademicSelect),
			widget.NewFormItem("Sağlık Durumu", stHealthSelect),
			widget.NewFormItem("Ekonomik Durum", stEconomicSelect),
			widget.NewFormItem("Kiminle Kaldığı", stLivingSelect),
			widget.NewFormItem("Anne-Baba Sağlık Durumu", stParentsLivingSelect),
			widget.NewFormItem("Anne-Baba Birliktelik", stParentsBioSelect),
			widget.NewFormItem("Büyüdüğü Çevre", stEnvSelect),
			widget.NewFormItem("Daha Önce Ceza Durumu", stPriorPenaltySelect),
			widget.NewFormItem("Öğrenci İfadesi / Savunması", stStatementEntry),
		)

		scrollContent := container.NewVScroll(form)
		stStatementEntry.ParentScroll = scrollContent
		scrollContent.SetMinSize(fyne.NewSize(580, 420))

		d := fyneDialog.NewCustomConfirm(
			fmt.Sprintf("Öğrenci Rol & İnceleme Profili: %s (%s)", is.FullName, is.SchoolNo),
			"Profili Kaydet",
			"İptal",
			scrollContent,
			func(ok bool) {
				if !ok {
					return
				}
				is.Role = stRoleSelect.Selected
				is.AcademicStanding = stAcademicSelect.Selected
				is.HealthStatus = stHealthSelect.Selected
				is.EconomicStatus = stEconomicSelect.Selected
				is.LivingWith = stLivingSelect.Selected
				is.ParentsLivingStatus = stParentsLivingSelect.Selected
				is.ParentsBioStatus = stParentsBioSelect.Selected
				is.EnvironmentRaised = stEnvSelect.Selected
				is.PriorPenaltyStatus = stPriorPenaltySelect.Selected
				is.StatementText = stStatementEntry.Text

				err := state.Store.SaveIncidentStudent(is)
				if err != nil {
					showError(state.Window, err)
					return
				}
				reloadIncStudents()
				ShowTurkishInfo("Başarılı", "Öğrenci rol ve profil bilgileri kaydedildi.", state.Window)
			},
			state.Window,
		)
		d.Resize(fyne.NewSize(620, 520))
		d.Show()
	}

	editStudentProfileBtn = widget.NewButtonWithIcon("Öğrenci Rol & İnceleme Profili", theme.DocumentCreateIcon(), func() {
		if selectedIncStudent != nil {
			showStudentProfileModal(selectedIncStudent)
		}
	})
	editStudentProfileBtn.Disable()

	removeStudentFromIncBtn = widget.NewButtonWithIcon("Öğrenciyi Olaydan Çıkar", theme.DeleteIcon(), func() {
		if selectedIncStudent == nil {
			showError(state.Window, fmt.Errorf("Çıkarılacak öğrenci seçilmedi"))
			return
		}
		ShowTurkishDangerConfirm(
			"Öğrenciyi Çıkar",
			fmt.Sprintf("%s adlı öğrenciyi bu olaydan çıkarmak istiyor musunuz?", selectedIncStudent.FullName),
			func(ok bool) {
				if !ok {
					return
				}
				_ = state.Store.DeleteIncidentStudent(selectedIncStudent.ID)
				selectedIncStudent = nil
				editStudentProfileBtn.Disable()
				removeStudentFromIncBtn.Disable()
				reloadIncStudents()
				if selectedIncident != nil {
					for idx, inc := range currentIncidents {
						if inc.IncidentCode == selectedIncident.IncidentCode {
							incidentList.Select(idx)
							break
						}
					}
				}
			},
			state.Window,
		)
	})
	removeStudentFromIncBtn.Importance = widget.DangerImportance
	removeStudentFromIncBtn.Disable()

	// Inline Student Statement & Role Editor (Under incident and student list)
	stStatementLabel = widget.NewLabelWithStyle("İfade girmek veya düzenlemek için yukarıdaki listeden bir öğrenci seçiniz.", fyne.TextAlignLeading, fyne.TextStyle{Italic: true})

	stInlineRoleSelect = NewTurkishSelect([]string{
		"SUÇLANAN ÖĞRENCİ",
		"MAĞDUR / ŞİKAYETÇİ ÖĞRENCİ",
		"BİLGİSİNE BAŞVURULACAK ÖĞRENCİ (TANIK)",
	}, nil)

	stInlineStatementEntry = NewSmartMultiLineEntry()
	stInlineStatementEntry.SetPlaceHolder("Öğrencinin olayla ilgili sözlü / yazılı ifadesi ve savunması...")

	saveStudentStatementBtn = widget.NewButtonWithIcon("Öğrenci İfadesini & Rolünü Kaydet", theme.DocumentSaveIcon(), func() {
		if selectedIncStudent == nil {
			showError(state.Window, fmt.Errorf("Lütfen önce listeden bir öğrenci seçiniz"))
			return
		}
		selectedIncStudent.Role = stInlineRoleSelect.Selected
		selectedIncStudent.StatementText = stInlineStatementEntry.Text

		err := state.Store.SaveIncidentStudent(selectedIncStudent)
		if err != nil {
			showError(state.Window, err)
			return
		}
		reloadIncStudents()
		ShowTurkishInfo("Kaydedildi", fmt.Sprintf("%s adlı öğrencinin ifadesi ve rolü kaydedildi.", selectedIncStudent.FullName), state.Window)
	})
	saveStudentStatementBtn.Importance = widget.HighImportance
	saveStudentStatementBtn.Disable()

	studentStatementForm := widget.NewForm(
		widget.NewFormItem("Olaydaki Rolü *", stInlineRoleSelect),
		widget.NewFormItem("Öğrenci İfadesi / Savunması *", stInlineStatementEntry),
	)

	statementCard := widget.NewCard("Öğrenci İfadesi & Olaydaki Rolü", "Olaya bağlı seçili öğrencinin ifadesi ve olaydaki rolü", container.NewVBox(
		stStatementLabel,
		studentStatementForm,
		container.NewHBox(saveStudentStatementBtn, editStudentProfileBtn),
	))

	incStudentsList.OnSelected = func(id widget.ListItemID) {
		incStudentsList.isActive = true
		if id >= 0 && id < len(incidentStudents) {
			is := incidentStudents[id]
			selectedIncStudent = &is
			stStatementLabel.SetText(fmt.Sprintf("Seçili Öğrenci: %s (%s - No: %s)", is.FullName, is.ClassSection, is.SchoolNo))
			stInlineRoleSelect.SetSelected(is.Role)
			stInlineStatementEntry.SetText(is.StatementText)
			saveStudentStatementBtn.Enable()
			editStudentProfileBtn.Enable()
			removeStudentFromIncBtn.Enable()
		}
	}

	// Forms Assembly
	incidentForm := widget.NewForm(
		widget.NewFormItem("Olay Kodu *", codeEntry),
		widget.NewFormItem("Olay Adı / Başlığı *", titleEntry),
		widget.NewFormItem("Tarih", dateEntry),
		widget.NewFormItem("Saat", timeEntry),
		widget.NewFormItem("Olay Yeri", locationEntry),
		widget.NewFormItem("Şikayetçi / İhbar Eden", complainantEntry),
		widget.NewFormItem("Dilekçe / Evrak Kayıt No", petitionEntry),
		widget.NewFormItem("Olayın Özeti", summaryEntry),
	)

	incidentCard := widget.NewCard("Olay Bilgileri", "İncelenen disiplin olayının genel tanımı", incidentForm)
	incidentBtns := container.NewHBox(saveIncidentBtn, deleteIncidentBtn, newIncidentBtn)

	studentsHeader := container.NewHBox(
		widget.NewLabelWithStyle("👥 Olaya Bağlı Öğrenciler", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		addStudentToIncBtn,
		removeStudentFromIncBtn,
	)

	studentsCard := widget.NewCard("Olaya Bağlı Öğrenci Listesi", "Soruşturmaya dahil edilen öğrenciler ve rolleri (Liste içinde kaydırmak için listeye tıklayın)", container.NewVBox(
		studentsHeader,
		incStudentsList.Wrap(),
	))

	// Layout Structure: 2-Panel Layout (Left: Incidents List, Right: Incident Form & Students)
	leftIncList := container.NewBorder(
		widget.NewLabelWithStyle("⚖️ Kayıtlı Olaylar", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		newIncidentBtn, nil, nil,
		incidentList,
	)

	rightPanelContent := container.NewVBox(
		incidentCard,
		incidentBtns,
		widget.NewSeparator(),
		studentsCard,
		widget.NewSeparator(),
		statementCard,
	)

	rightPanelScroll := container.NewVScroll(rightPanelContent)
	summaryEntry.ParentScroll = rightPanelScroll
	stInlineStatementEntry.ParentScroll = rightPanelScroll
	incStudentsList.ParentScroll = rightPanelScroll

	splitMain := container.NewHSplit(leftIncList, rightPanelScroll)
	splitMain.SetOffset(0.28)

	reloadIncidents()

	// Automatically select the active or first incident if available
	if state.SelectedIncidentCode != "" {
		for idx, inc := range currentIncidents {
			if inc.IncidentCode == state.SelectedIncidentCode {
				incidentList.Select(idx)
				break
			}
		}
	} else if len(currentIncidents) > 0 {
		incidentList.Select(0)
	} else {
		clearIncidentForm()
	}

	return splitMain
}
