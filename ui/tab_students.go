package ui

import (
	"fmt"
	"strconv"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"turandisiplin/excel"
	"turandisiplin/model"
)

// BuildStudentsTab creates the Student Roster and Excel Import UI
func BuildStudentsTab(state *AppState) fyne.CanvasObject {
	var currentStudents []model.Student
	var selectedStudent *model.Student

	// Count Label
	countLabel := widget.NewLabelWithStyle("Kayıt: 0 öğrenci", fyne.TextAlignLeading, fyne.TextStyle{Bold: true})

	studentRolesLabel := widget.NewLabel("Seçili öğrencinin olay / disiplin geçmişi ve rolleri burada görüntülenir.")
	studentRolesLabel.Wrapping = fyne.TextWrapWord

	// Search & Filter
	searchEntry := NewSingleLineEntry()
	searchEntry.SetPlaceHolder("Öğrenci no veya ad soyad ile ara...")

	classOptions := append([]string{"Tüm Sınıflar"}, state.Store.GetClasses(state.AcademicYear)...)
	classSelect := NewTurkishSelect(classOptions, nil)
	classSelect.SetSelected("Tüm Sınıflar")

	// Student Detail Form Inputs
	idEntry := NewSingleLineEntry()
	idEntry.Disable()

	schoolNoEntry := NewSingleLineEntry()
	fullNameEntry := NewSingleLineEntry()
	classEntry := NewSingleLineEntry()
	genderSelect := NewTurkishSelect([]string{"ERKEK", "KIZ"}, nil)
	tcEntry := NewSingleLineEntry()
	birthPlaceEntry := NewSingleLineEntry()
	birthDateEntry := NewSingleLineEntry()
	fatherEntry := NewSingleLineEntry()
	motherEntry := NewSingleLineEntry()
	guardianEntry := NewSingleLineEntry()
	phoneEntry := NewSingleLineEntry()
	addrEntry := NewSingleLineEntry()
	boardingSelect := NewTurkishSelect([]string{"GÜNDÜZLÜ", "PANSİYON / YATILI"}, nil)
	scoreEntry := NewSingleLineEntry()
	scoreEntry.SetText("100")

	clearForm := func() {
		selectedStudent = nil
		idEntry.SetText("")
		schoolNoEntry.SetText("")
		fullNameEntry.SetText("")
		classEntry.SetText("")
		genderSelect.ClearSelected()
		tcEntry.SetText("")
		birthPlaceEntry.SetText("")
		birthDateEntry.SetText("")
		fatherEntry.SetText("")
		motherEntry.SetText("")
		guardianEntry.SetText("")
		phoneEntry.SetText("")
		addrEntry.SetText("")
		boardingSelect.SetSelected("GÜNDÜZLÜ")
		scoreEntry.SetText("100")
		studentRolesLabel.SetText("Seçili öğrencinin olay / disiplin geçmişi ve rolleri burada görüntülenir.")
	}

	populateForm := func(st *model.Student) {
		selectedStudent = st
		idEntry.SetText(fmt.Sprintf("%d", st.ID))
		schoolNoEntry.SetText(st.SchoolNo)
		fullNameEntry.SetText(st.FullName)
		classEntry.SetText(st.ClassSection)
		genderSelect.SetSelected(st.Gender)
		tcEntry.SetText(st.TCNo)
		birthPlaceEntry.SetText(st.BirthPlace)
		birthDateEntry.SetText(st.BirthDate)
		fatherEntry.SetText(st.FatherName)
		motherEntry.SetText(st.MotherName)
		guardianEntry.SetText(st.GuardianName)
		phoneEntry.SetText(st.GuardianPhone)
		addrEntry.SetText(st.GuardianAddr)
		boardingSelect.SetSelected(st.BoardingStatus)
		scoreEntry.SetText(fmt.Sprintf("%d", st.BehaviorScore))

		// Review student's roles in incidents & disciplinary history
		records, _ := state.Store.GetIncidentStudentsBySchoolNo(state.AcademicYear, st.SchoolNo)
		if len(records) == 0 {
			studentRolesLabel.SetText("Bu öğrencinin bu öğretim yılında kayıtlı bir olay veya disiplin soruşturması bulunmamaktadır.")
		} else {
			var sb strings.Builder
			for _, r := range records {
				roleStr := r.Role
				if roleStr == "" {
					roleStr = "ROL BELİRTİLMEMİŞ"
				}
				outcome := "Soruşturma Aşamasında"
				if r.VerdictNo != "" && r.ImposedPenalty != "" {
					outcome = fmt.Sprintf("Karar No %s: %s", r.VerdictNo, r.ImposedPenalty)
				}
				sb.WriteString(fmt.Sprintf("• Olay: %s (Olay Kodu: %s)\n  Olaydaki Rolü: %s\n  Durum / Ceza: %s\n", r.IncidentTitle, r.IncidentCode, roleStr, outcome))
				if r.StatementText != "" {
					stmt := r.StatementText
					if len(stmt) > 90 {
						stmt = stmt[:87] + "..."
					}
					sb.WriteString(fmt.Sprintf("  Öğrenci İfadesi: \"%s\"\n", stmt))
				}
				sb.WriteString("\n")
			}
			studentRolesLabel.SetText(strings.TrimSpace(sb.String()))
		}
	}

	// Student List (using widget.List for smooth senior-friendly scrolling)
	studentList := widget.NewList(
		func() int {
			return len(currentStudents)
		},
		func() fyne.CanvasObject {
			return container.NewHBox(
				widget.NewLabelWithStyle("00000", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
				widget.NewSeparator(),
				widget.NewLabelWithStyle("9/A", fyne.TextAlignLeading, fyne.TextStyle{Italic: true}),
				widget.NewSeparator(),
				widget.NewLabel("Öğrenci Adı Soyadı"),
				widget.NewSeparator(),
				widget.NewLabel("Puan: 100"),
			)
		},
		func(i widget.ListItemID, o fyne.CanvasObject) {
			if i >= len(currentStudents) {
				return
			}
			st := currentStudents[i]
			box := o.(*fyne.Container)
			noLbl := box.Objects[0].(*widget.Label)
			classLbl := box.Objects[2].(*widget.Label)
			nameLbl := box.Objects[4].(*widget.Label)
			scoreLbl := box.Objects[6].(*widget.Label)

			noLbl.SetText(fmt.Sprintf("%-5s", st.SchoolNo))
			classLbl.SetText(st.ClassSection)
			nameLbl.SetText(st.FullName)
			scoreLbl.SetText(fmt.Sprintf("Puan: %d", st.BehaviorScore))
		},
	)

	reloadStudents := func() {
		filterClass := classSelect.Selected
		if filterClass == "Tüm Sınıflar" {
			filterClass = ""
		}
		list, err := state.Store.GetStudents(state.AcademicYear, searchEntry.Text, filterClass)
		if err == nil {
			currentStudents = list
			countLabel.SetText(fmt.Sprintf("Kayıt: %d öğrenci", len(currentStudents)))
			studentList.Refresh()
		}
	}

	studentList.OnSelected = func(id widget.ListItemID) {
		if id >= 0 && id < len(currentStudents) {
			st := currentStudents[id]
			populateForm(&st)
			state.SelectedStudentNo = st.SchoolNo
		}
	}

	searchEntry.OnChanged = func(s string) {
		reloadStudents()
	}

	classSelect.OnChanged = func(s string) {
		reloadStudents()
	}


	// Excel Import Dialog using native OS file dialog
	excelBtn := widget.NewButtonWithIcon("Excel'den Öğrenci Yükle (.xlsx / .xls)", theme.DownloadIcon(), func() {
		SelectExcelFile(state.Window, func(filePath string) {
			// Prompt with the 3 import modes
			modeRadio := widget.NewRadioGroup([]string{
				"1. Sadece Yeni Öğrencileri Ekle (Var olanlara dokunma)",
				"2. Tüm Öğrencileri Sil ve Yeniden Yükle",
				"3. Mevcut Öğrencilerin Sınıflarını Güncelle",
			}, nil)
			modeRadio.SetSelected("1. Sadece Yeni Öğrencileri Ekle (Var olanlara dokunma)")

			infoLabel := widget.NewLabel(fmt.Sprintf("Dosya: %s\nÖğretim Yılı: %s", filePath, state.AcademicYear))
			content := container.NewVBox(infoLabel, widget.NewSeparator(), modeRadio)

			dialog.ShowCustomConfirm("Excel İçe Aktarma Modu", "İçe Aktar", "İptal", content, func(ok bool) {
				if !ok {
					return
				}

				parsed, err := excel.ImportStudentsFromFile(filePath, state.AcademicYear)
				if err != nil {
					showError(state.Window, fmt.Errorf("Excel okunurken hata oluştu: %w", err))
					return
				}

				switch modeRadio.Selected {
				case "2. Tüm Öğrencileri Sil ve Yeniden Yükle":
					_ = state.Store.DeleteAllStudents(state.AcademicYear)
					for _, st := range parsed {
						_ = state.Store.SaveStudent(st)
					}
				case "3. Mevcut Öğrencilerin Sınıflarını Güncelle":
					for _, st := range parsed {
						existing, err := state.Store.GetStudentByNo(state.AcademicYear, st.SchoolNo)
						if err == nil && existing != nil {
							existing.ClassSection = st.ClassSection
							_ = state.Store.SaveStudent(existing)
						} else {
							_ = state.Store.SaveStudent(st)
						}
					}
				default: // Mode 1
					for _, st := range parsed {
						existing, _ := state.Store.GetStudentByNo(state.AcademicYear, st.SchoolNo)
						if existing == nil {
							_ = state.Store.SaveStudent(st)
						}
					}
				}

				showInfo(state.Window, "İçe Aktarma Tamamlandı", fmt.Sprintf("%d öğrenci başarıyla işlendi.", len(parsed)))
				state.NotifyChange()
			}, state.Window)
		})
	})
	excelBtn.Importance = widget.MediumImportance

	// Student Actions: Save, Delete, Clear
	saveStudentBtn := widget.NewButtonWithIcon("Öğrenciyi Kaydet", theme.DocumentSaveIcon(), func() {
		no := strings.TrimSpace(schoolNoEntry.Text)
		name := strings.TrimSpace(fullNameEntry.Text)
		if no == "" || name == "" {
			showError(state.Window, fmt.Errorf("Okul numarası ve Ad Soyad zorunludur"))
			return
		}

		score := 100
		if s, err := strconv.Atoi(scoreEntry.Text); err == nil {
			score = s
		}

		var id int64
		if selectedStudent != nil {
			id = selectedStudent.ID
		}

		st := &model.Student{
			ID:             id,
			AcademicYear:   state.AcademicYear,
			SchoolNo:       no,
			FullName:       name,
			ClassSection:   strings.TrimSpace(classEntry.Text),
			Gender:         genderSelect.Selected,
			TCNo:           strings.TrimSpace(tcEntry.Text),
			BirthPlace:     strings.TrimSpace(birthPlaceEntry.Text),
			BirthDate:      strings.TrimSpace(birthDateEntry.Text),
			FatherName:     strings.TrimSpace(fatherEntry.Text),
			MotherName:     strings.TrimSpace(motherEntry.Text),
			GuardianName:   strings.TrimSpace(guardianEntry.Text),
			GuardianPhone:  strings.TrimSpace(phoneEntry.Text),
			GuardianAddr:   strings.TrimSpace(addrEntry.Text),
			BoardingStatus: boardingSelect.Selected,
			BehaviorScore:  score,
		}

		err := state.Store.SaveStudent(st)
		if err != nil {
			showError(state.Window, err)
			return
		}

		showInfo(state.Window, "Başarılı", "Öğrenci kaydı kaydedildi.")
		clearForm()
		state.NotifyChange()
	})
	saveStudentBtn.Importance = widget.HighImportance

	deleteStudentBtn := widget.NewButtonWithIcon("Öğrenciyi Sil", theme.DeleteIcon(), func() {
		if selectedStudent == nil {
			showError(state.Window, fmt.Errorf("Silinecek öğrenci seçilmedi"))
			return
		}
		ShowTurkishDangerConfirm(
			"Öğrenci Sil",
			fmt.Sprintf("%s (%s) kaydını silmek istediğinize emin misiniz?", selectedStudent.FullName, selectedStudent.SchoolNo),
			func(confirmed bool) {
				if !confirmed {
					return
				}
				_ = state.Store.DeleteStudent(selectedStudent.ID)
				clearForm()
				state.NotifyChange()
			},
			state.Window,
		)
	})
	deleteStudentBtn.Importance = widget.DangerImportance

	newStudentBtn := widget.NewButtonWithIcon("Formu Temizle / Yeni", theme.ContentAddIcon(), func() {
		clearForm()
	})
	newStudentBtn.Importance = widget.MediumImportance

	// Layout Construction
	filterRow := container.NewBorder(
		nil, nil,
		countLabel,
		classSelect,
		searchEntry,
	)

	leftPanel := container.NewBorder(
		container.NewVBox(filterRow, excelBtn),
		nil, nil, nil,
		studentList,
	)

	formItems := widget.NewForm(
		widget.NewFormItem("Okul No *", schoolNoEntry),
		widget.NewFormItem("Adı Soyadı *", fullNameEntry),
		widget.NewFormItem("Sınıfı / Şubesi", classEntry),
		widget.NewFormItem("Cinsiyeti", genderSelect),
		widget.NewFormItem("T.C. Kimlik No", tcEntry),
		widget.NewFormItem("Doğum Yeri", birthPlaceEntry),
		widget.NewFormItem("Doğum Tarihi", birthDateEntry),
		widget.NewFormItem("Baba Adı", fatherEntry),
		widget.NewFormItem("Anne Adı", motherEntry),
		widget.NewFormItem("Veli Adı Soyadı", guardianEntry),
		widget.NewFormItem("Veli Telefonu", phoneEntry),
		widget.NewFormItem("Veli Adresi", addrEntry),
		widget.NewFormItem("Pansiyon / Gündüzlü", boardingSelect),
		widget.NewFormItem("Davranış Puanı", scoreEntry),
	)

	formCard := widget.NewCard("Öğrenci Detayları", "Listeden bir öğrenci seçin veya yeni öğrenci bilgilerini girip kaydedin.", formItems)
	roleHistoryCard := widget.NewCard("Disiplin Olayları ve Rol Geçmişi", "Öğrencinin bu öğretim yılındaki olay rolleri ve soruşturma geçmişi", studentRolesLabel)
	buttonsRow := container.NewHBox(saveStudentBtn, deleteStudentBtn, newStudentBtn)

	rightPanel := container.NewVScroll(container.NewVBox(formCard, roleHistoryCard, buttonsRow))

	split := container.NewHSplit(leftPanel, rightPanel)
	split.SetOffset(0.5)

	reloadStudents()
	return split
}
