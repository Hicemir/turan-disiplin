package ui

import (
	"fmt"
	"strings"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"okulevrak/model"
)

// BuildIncidentsTab creates the Incident Management & Student Assignment UI
func BuildIncidentsTab(state *AppState) fyne.CanvasObject {
	var currentIncidents []model.Incident
	var selectedIncident *model.Incident

	var incidentStudents []model.IncidentStudent
	var selectedIncStudent *model.IncidentStudent

	// Incident Form Inputs
	codeEntry := widget.NewEntry()
	titleEntry := widget.NewEntry()
	dateEntry := widget.NewEntry()
	dateEntry.SetText(time.Now().Format("02.01.2006"))
	timeEntry := widget.NewEntry()
	timeEntry.SetText(time.Now().Format("15:04"))
	locationEntry := widget.NewEntry()
	locationEntry.SetPlaceHolder("Örn: Okul Bahçesi, 10/B Sınıfı")
	complainantEntry := widget.NewEntry()
	complainantEntry.SetPlaceHolder("Örn: Nöbetçi Öğretmen Ahmet Yılmaz")
	petitionEntry := widget.NewEntry()
	summaryEntry := widget.NewMultiLineEntry()
	summaryEntry.SetPlaceHolder("Olayın detaylı açıklaması...")

	// Attached Student Inputs
	stRoleSelect := widget.NewSelect([]string{
		"SUÇLANAN ÖĞRENCİ",
		"MAĞDUR / ŞİKAYETÇİ ÖĞRENCİ",
		"BİLGİSİNE BAŞVURULACAK ÖĞRENCİ",
	}, nil)
	stRoleSelect.SetSelected("SUÇLANAN ÖĞRENCİ")

	stAcademicSelect := widget.NewSelect([]string{"ÇOK İYİ", "İYİ", "ORTA", "ZAYIF"}, nil)
	stHealthSelect := widget.NewSelect([]string{"SAĞLAM", "KRONİK RAHATSIZLIĞI VAR", "ENGELLİ", "DİĞER"}, nil)
	stEconomicSelect := widget.NewSelect([]string{"İYİ", "ORTA", "DAR GELİRLİ", "ÇOK DÜŞÜK"}, nil)
	stLivingSelect := widget.NewSelect([]string{"AİLESİYLE BİRLİKTE", "YALNIZ ANNESİYLE", "YALNIZ BABASIYLA", "AKRABALARIYLA", "YURTTA / PANSİYONDA", "YALNIZ"}, nil)
	stParentsLivingSelect := widget.NewSelect([]string{"HER İKİSİ DE SAĞ", "BABA VEFAT ETMİŞ", "ANNE VEFAT ETMİŞ", "HER İKİSİ DE VEFAT ETMİŞ"}, nil)
	stParentsBioSelect := widget.NewSelect([]string{"ÖZ VE BİRLİKTE YAŞIYORLAR", "BOŞANMIŞLAR / AYRI YAŞIYORLAR", "ÜVEY ANNE", "ÜVEY BABA"}, nil)
	stEnvSelect := widget.NewSelect([]string{"BÜYÜKŞEHİR", "İL / İLÇE MERKEZİ", "KASABA / BELDE", "KÖY"}, nil)
	stPriorPenaltySelect := widget.NewSelect([]string{"DAHA ÖNCE CEZA ALMAMIŞTIR", "UYARMA CEZASI ALMIŞTIR", "KINAMA CEZASI ALMIŞTIR", "OKULDAN KISA SÜRELİ UZAKLAŞTIRMA ALMIŞTIR"}, nil)
	stStatementEntry := widget.NewMultiLineEntry()
	stStatementEntry.SetPlaceHolder("Öğrencinin olayla ilgili sözlü / yazılı ifadesi...")

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
		selectedIncStudent = nil
		state.SelectedIncStudentID = 0
	}

	clearStudentForm := func() {
		selectedIncStudent = nil
		state.SelectedIncStudentID = 0
		stRoleSelect.SetSelected("SUÇLANAN ÖĞRENCİ")
		stAcademicSelect.ClearSelected()
		stHealthSelect.ClearSelected()
		stEconomicSelect.ClearSelected()
		stLivingSelect.ClearSelected()
		stParentsLivingSelect.ClearSelected()
		stParentsBioSelect.ClearSelected()
		stEnvSelect.ClearSelected()
		stPriorPenaltySelect.ClearSelected()
		stStatementEntry.SetText("")
	}

	populateStudentForm := func(is *model.IncidentStudent) {
		selectedIncStudent = is
		state.SelectedIncStudentID = is.ID
		stRoleSelect.SetSelected(is.Role)
		stAcademicSelect.SetSelected(is.AcademicStanding)
		stHealthSelect.SetSelected(is.HealthStatus)
		stEconomicSelect.SetSelected(is.EconomicStatus)
		stLivingSelect.SetSelected(is.LivingWith)
		stParentsLivingSelect.SetSelected(is.ParentsLivingStatus)
		stParentsBioSelect.SetSelected(is.ParentsBioStatus)
		stEnvSelect.SetSelected(is.EnvironmentRaised)
		stPriorPenaltySelect.SetSelected(is.PriorPenaltyStatus)
		stStatementEntry.SetText(is.StatementText)
	}

	// Attached Student List
	incStudentsList := widget.NewList(
		func() int {
			return len(incidentStudents)
		},
		func() fyne.CanvasObject {
			return container.NewHBox(
				widget.NewLabelWithStyle("No", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
				widget.NewSeparator(),
				widget.NewLabel("Öğrenci Adı"),
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
			roleLbl := box.Objects[4].(*widget.Label)

			noLbl.SetText(fmt.Sprintf("%-5s", is.SchoolNo))
			nameLbl.SetText(is.FullName)
			roleLbl.SetText(is.Role)
		},
	)

	reloadIncStudents := func() {
		if selectedIncident == nil {
			incidentStudents = nil
			incStudentsList.Refresh()
			return
		}
		list, err := state.Store.GetIncidentStudents(state.AcademicYear, selectedIncident.IncidentCode)
		if err == nil {
			incidentStudents = list
			incStudentsList.Refresh()
		}
	}

	incStudentsList.OnSelected = func(id widget.ListItemID) {
		if id >= 0 && id < len(incidentStudents) {
			is := incidentStudents[id]
			populateStudentForm(&is)
		}
	}

	// Incident List
	incidentList := widget.NewList(
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
		}
	}

	incidentList.OnSelected = func(id widget.ListItemID) {
		if id >= 0 && id < len(currentIncidents) {
			inc := currentIncidents[id]
			selectedIncident = &inc
			state.SelectedIncidentCode = inc.IncidentCode

			codeEntry.SetText(inc.IncidentCode)
			titleEntry.SetText(inc.Title)
			dateEntry.SetText(inc.Date)
			timeEntry.SetText(inc.Time)
			locationEntry.SetText(inc.Location)
			complainantEntry.SetText(inc.Complainant)
			petitionEntry.SetText(inc.PetitionRef)
			summaryEntry.SetText(inc.Summary)

			clearStudentForm()
			reloadIncStudents()
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

		showInfo(state.Window, "Başarılı", "Olay kaydı başarıyla kaydedildi.")
		selectedIncident = inc
		state.SelectedIncidentCode = code
		reloadIncidents()
		state.NotifyChange()
	})
	saveIncidentBtn.Importance = widget.HighImportance

	deleteIncidentBtn := widget.NewButtonWithIcon("Olayı Sil", theme.DeleteIcon(), func() {
		if selectedIncident == nil {
			showError(state.Window, fmt.Errorf("Silinecek olay seçilmedi"))
			return
		}
		showConfirm(state.Window, "Olay Sil", fmt.Sprintf("'%s' olayını ve bu olaya bağlı tüm öğrenci/karar kayıtlarını silmek istediğinize emin misiniz?", selectedIncident.Title), func() {
			_ = state.Store.DeleteIncident(state.AcademicYear, selectedIncident.IncidentCode)
			clearIncidentForm()
			reloadIncidents()
			state.NotifyChange()
		})
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
		stSearch := widget.NewEntry()
		stSearch.SetPlaceHolder("Öğrenci ara...")
		var foundList []model.Student

		pickList := widget.NewList(
			func() int { return len(foundList) },
			func() fyne.CanvasObject {
				return container.NewHBox(
					widget.NewLabel("00000"),
					widget.NewLabel("Ad Soyad"),
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
				box.Objects[1].(*widget.Label).SetText(s.FullName)
				box.Objects[2].(*widget.Label).SetText(s.ClassSection)
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

		d := dialog.NewCustomConfirm("Olaya Öğrenci Ekle", "Ekle", "İptal", pickerContent, func(ok bool) {
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
				GuardianName:        chosen.GuardianName,
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
			showInfo(state.Window, "Öğrenci Eklendi", fmt.Sprintf("%s (%s) olaya eklendi.", chosen.FullName, chosen.SchoolNo))
			state.NotifyChange()
		}, state.Window)

		d.Resize(fyne.NewSize(500, 400))
		updatePickList()
		d.Show()
	})
	addStudentToIncBtn.Importance = widget.MediumImportance

	saveStudentProfileBtn := widget.NewButtonWithIcon("Öğrenci Profilini Kaydet", theme.DocumentSaveIcon(), func() {
		if selectedIncStudent == nil {
			showError(state.Window, fmt.Errorf("Kaydedilecek öğrenci seçilmedi"))
			return
		}

		selectedIncStudent.Role = stRoleSelect.Selected
		selectedIncStudent.AcademicStanding = stAcademicSelect.Selected
		selectedIncStudent.HealthStatus = stHealthSelect.Selected
		selectedIncStudent.EconomicStatus = stEconomicSelect.Selected
		selectedIncStudent.LivingWith = stLivingSelect.Selected
		selectedIncStudent.ParentsLivingStatus = stParentsLivingSelect.Selected
		selectedIncStudent.ParentsBioStatus = stParentsBioSelect.Selected
		selectedIncStudent.EnvironmentRaised = stEnvSelect.Selected
		selectedIncStudent.PriorPenaltyStatus = stPriorPenaltySelect.Selected
		selectedIncStudent.StatementText = stStatementEntry.Text

		err := state.Store.SaveIncidentStudent(selectedIncStudent)
		if err != nil {
			showError(state.Window, err)
			return
		}

		showInfo(state.Window, "Başarılı", "Öğrenci rol ve sosyo-ekonomik inceleme bilgileri kaydedildi.")
		reloadIncStudents()
		state.NotifyChange()
	})
	saveStudentProfileBtn.Importance = widget.HighImportance

	removeStudentFromIncBtn := widget.NewButtonWithIcon("Öğrenciyi Olaydan Çıkar", theme.DeleteIcon(), func() {
		if selectedIncStudent == nil {
			showError(state.Window, fmt.Errorf("Çıkarılacak öğrenci seçilmedi"))
			return
		}
		showConfirm(state.Window, "Öğrenciyi Çıkar", fmt.Sprintf("%s adlı öğrenciyi bu olaydan çıkarmak istiyor musunuz?", selectedIncStudent.FullName), func() {
			_ = state.Store.DeleteIncidentStudent(selectedIncStudent.ID)
			clearStudentForm()
			reloadIncStudents()
			state.NotifyChange()
		})
	})
	removeStudentFromIncBtn.Importance = widget.DangerImportance

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

	studentProfileForm := widget.NewForm(
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

	studentCard := widget.NewCard("Öğrenci Rol & İnceleme Profili", "Karar tutanağı ve evraklara girecek detaylı profil", studentProfileForm)
	studentBtns := container.NewHBox(saveStudentProfileBtn, removeStudentFromIncBtn)

	// Layout Structure
	leftIncList := container.NewBorder(
		widget.NewLabelWithStyle("⚖️ Kayıtlı Olaylar", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		newIncidentBtn, nil, nil,
		incidentList,
	)

	middleIncDetails := container.NewVScroll(container.NewVBox(
		incidentCard,
		incidentBtns,
		widget.NewSeparator(),
		container.NewHBox(
			widget.NewLabelWithStyle("👥 Olaya Bağlı Öğrenciler", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
			addStudentToIncBtn,
		),
		incStudentsList,
	))

	rightStudentDetails := container.NewVScroll(container.NewVBox(
		studentCard,
		studentBtns,
	))

	splitRight := container.NewHSplit(middleIncDetails, rightStudentDetails)
	splitRight.SetOffset(0.5)

	splitMain := container.NewHSplit(leftIncList, splitRight)
	splitMain.SetOffset(0.25)

	reloadIncidents()
	clearIncidentForm()

	return splitMain
}
