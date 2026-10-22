package ui

import (
	"fmt"
	"strings"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"turandisiplin/model"
	"turandisiplin/pdf"
)

// BuildSummonsTab creates the Student & Member Summons Slips UI
func BuildSummonsTab(state *AppState) fyne.CanvasObject {
	var currentIncidents []model.Incident
	var currentStudents []model.Student
	var selectedIncident *model.Incident

	// --- 1. Student Summons Controls ---
	incSelect := NewTurkishSelect([]string{}, nil)
	callTypeRadio := widget.NewRadioGroup([]string{"Savunma Çağrı Pusulası", "İfade Çağrı Pusulası"}, nil)
	callTypeRadio.Horizontal = true
	callTypeRadio.SetSelected("Savunma Çağrı Pusulası")

	placeEntry := NewSingleLineEntry()
	placeEntry.SetText("Disiplin Kurulu Odası")
	dateEntry := NewSingleLineEntry()
	dateEntry.SetText(time.Now().Format("02.01.2006"))
	timeEntry := NewSingleLineEntry()
	timeEntry.SetText("10:00")

	// Virtualized multi-student selection checklist
	var incidentStudents []model.IncidentStudent
	var studentChecked []bool

	studentCheckList := NewFocusableList(
		func() int {
			return len(incidentStudents)
		},
		func() fyne.CanvasObject {
			return widget.NewCheck("00000 - Öğrenci Adı Soyadı (ROL)", nil)
		},
		func(i widget.ListItemID, o fyne.CanvasObject) {
			if i >= len(incidentStudents) {
				return
			}
			is := incidentStudents[i]
			chk := o.(*widget.Check)
			chk.OnChanged = nil
			chk.Text = fmt.Sprintf("%s - %s (%s)", is.SchoolNo, is.FullName, is.Role)
			if i < len(studentChecked) {
				chk.Checked = studentChecked[i]
			}
			chk.Refresh()
			idx := i
			chk.OnChanged = func(b bool) {
				if idx < len(studentChecked) {
					studentChecked[idx] = b
				}
			}
		},
		fyne.NewSize(400, 120),
	)

	reloadStudentChecks := func(incCode string) {
		incidentStudents = nil
		studentChecked = nil

		list, err := state.Store.GetIncidentStudents(state.AcademicYear, incCode)
		if err == nil {
			incidentStudents = list
			studentChecked = make([]bool, len(list))
			for i := range studentChecked {
				studentChecked[i] = true
			}
		}
		studentCheckList.Refresh()
	}

	reloadIncidents := func() {
		list, err := state.Store.GetIncidents(state.AcademicYear)
		if err == nil {
			currentIncidents = list
			var opts []string
			for _, inc := range list {
				opts = append(opts, fmt.Sprintf("#%s - %s", inc.IncidentCode, inc.Title))
			}
			incSelect.Options = opts
			incSelect.Refresh()
			if len(opts) > 0 {
				incSelect.SetSelectedIndex(0)
			}
		}
	}

	incSelect.OnChanged = func(selected string) {
		for _, inc := range currentIncidents {
			prefix := fmt.Sprintf("#%s - ", inc.IncidentCode)
			if strings.HasPrefix(selected, prefix) {
				selectedIncident = &inc
				reloadStudentChecks(inc.IncidentCode)
				break
			}
		}
	}

	getSelectedStudents := func() []*model.Student {
		var res []*model.Student
		if selectedIncident == nil {
			return res
		}
		for i, is := range incidentStudents {
			if i < len(studentChecked) && studentChecked[i] {
				st, _ := state.Store.GetStudentByNo(state.AcademicYear, is.SchoolNo)
				if st == nil {
					st = &model.Student{
						SchoolNo:     is.SchoolNo,
						FullName:     is.FullName,
						ClassSection: is.ClassSection,
					}
				}
				res = append(res, st)
			}
		}
		return res
	}

	pdfStudentSummonsBtn := widget.NewButtonWithIcon("PDF OLUŞTUR VE AÇ (ÖĞRENCİ ÇAĞRI PUSULASI)", theme.DocumentIcon(), func() {
		selectedList := getSelectedStudents()
		if len(selectedList) == 0 {
			showError(state.Window, fmt.Errorf("Lütfen en az bir öğrenci işaretleyiniz"))
			return
		}

		isSavunma := callTypeRadio.Selected == "Savunma Çağrı Pusulası"
		prefix := "Savunma_Cagri_Pusulasi"
		if !isSavunma {
			prefix = "Ifade_Cagri_Pusulasi"
		}

		outPath, err := pdf.GenerateUniquePDFPath("", prefix)
		if err != nil {
			showError(state.Window, err)
			return
		}

		err = pdf.GenerateCagriPusulasi(state.SchoolInfo, selectedList, selectedIncident, isSavunma, placeEntry.Text, dateEntry.Text, timeEntry.Text, outPath)
		if err != nil {
			showError(state.Window, fmt.Errorf("PDF oluşturulamadı: %w", err))
			return
		}

		openPDF(state, outPath)
	})
	pdfStudentSummonsBtn.Importance = widget.SuccessImportance

	pdfStudentListBtn := widget.NewButtonWithIcon("PDF ÇAĞRI LİSTESİ OLUŞTUR", theme.DocumentPrintIcon(), func() {
		selectedList := getSelectedStudents()
		if len(selectedList) == 0 {
			showError(state.Window, fmt.Errorf("Lütfen en az bir öğrenci işaretleyiniz"))
			return
		}

		outPath, err := pdf.GenerateUniquePDFPath("", "Cagri_Pusulasi_Ogrenci_Listesi")
		if err != nil {
			showError(state.Window, err)
			return
		}

		err = pdf.GenerateCagriOgrenciListesi(state.SchoolInfo, selectedList, outPath)
		if err != nil {
			showError(state.Window, fmt.Errorf("PDF oluşturulamadı: %w", err))
			return
		}

		openPDF(state, outPath)
	})

	// --- 2. Board Member Summons Controls ---
	sayiEntry := NewSingleLineEntry()
	sayiEntry.SetText(state.SchoolInfo.OfficialCode + "100")
	konuEntry := NewSingleLineEntry()
	konuEntry.SetText("Disiplin Kurulu Toplantı Çağrısı")
	uyeTarihEntry := NewSingleLineEntry()
	uyeTarihEntry.SetText(time.Now().Format("02.01.2006"))
	uyeSaatEntry := NewSingleLineEntry()
	uyeSaatEntry.SetText("11:00")
	uyeYerEntry := NewSingleLineEntry()
	uyeYerEntry.SetText("Toplantı Salonu")

	// Agenda Items
	agendaText := NewSmartMultiLineEntry()
	agendaText.SetText("1. Açılış ve yoklama\n2. Hakkında soruşturma açılan öğrencilerin durumlarının görüşülmesi\n3. Kararların alınması ve imzalanması\n4. Kapanış")

	pdfMemberSummonsBtn := widget.NewButtonWithIcon("PDF OLUŞTUR VE AÇ (ÜYE ÇAĞRI PUSULASI)", theme.DocumentIcon(), func() {
		lines := strings.Split(agendaText.Text, "\n")
		var agenda []string
		for _, l := range lines {
			l = strings.TrimSpace(l)
			if l != "" {
				agenda = append(agenda, l)
			}
		}

		outPath, err := pdf.GenerateUniquePDFPath("", "Uye_Toplanti_Cagri_Pusulasi")
		if err != nil {
			showError(state.Window, err)
			return
		}

		err = pdf.GenerateUyeCagriPusulasi(state.SchoolInfo, sayiEntry.Text, konuEntry.Text, uyeTarihEntry.Text, uyeSaatEntry.Text, uyeYerEntry.Text, agenda, outPath)
		if err != nil {
			showError(state.Window, fmt.Errorf("PDF oluşturulamadı: %w", err))
			return
		}

		openPDF(state, outPath)
	})
	pdfMemberSummonsBtn.Importance = widget.SuccessImportance

	// Student Section Layout
	stForm := widget.NewForm(
		widget.NewFormItem("İncelenen Olay", incSelect),
		widget.NewFormItem("Çağrı Pusulası Türü", callTypeRadio),
		widget.NewFormItem("Görüşme Yeri", placeEntry),
		widget.NewFormItem("Tarih", dateEntry),
		widget.NewFormItem("Saat", timeEntry),
	)

	cardStudent := widget.NewCard("Öğrenci Savunma ve İfade Çağrı Pusulası", "Çağrı yapılacak öğrencileri seçip resmi pusulayı oluşturun", container.NewVBox(
		stForm,
		widget.NewLabelWithStyle("Pusula Çıkarılacak Öğrenciler:", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		studentCheckList.Wrap(),
		container.NewHBox(pdfStudentSummonsBtn, pdfStudentListBtn),
	))

	// Member Section Layout
	memberForm := widget.NewForm(
		widget.NewFormItem("Resmi Yazı Sayısı", sayiEntry),
		widget.NewFormItem("Konu", konuEntry),
		widget.NewFormItem("Toplantı Tarihi", uyeTarihEntry),
		widget.NewFormItem("Toplantı Saati", uyeSaatEntry),
		widget.NewFormItem("Toplantı Yeri", uyeYerEntry),
	)

	// Agenda sits in the Border center so it absorbs all spare vertical space;
	// SetMinRowsVisible keeps a usable minimum (the entry scrolls internally).
	agendaText.SetMinRowsVisible(4)
	agendaBox := container.NewBorder(
		widget.NewLabelWithStyle("Gündem Maddeleri (Her satıra bir madde)", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		nil, nil, nil,
		agendaText,
	)

	cardMember := widget.NewCard("Disiplin / Davranış Kurulu Üyeleri Toplantı Çağrısı", "Gündemli resmi kurul toplantı çağrı pusulası", container.NewBorder(
		memberForm,
		pdfMemberSummonsBtn,
		nil, nil,
		agendaBox,
	))

	content := container.NewVScroll(container.NewBorder(
		container.NewVBox(cardStudent, widget.NewSeparator()),
		nil, nil, nil,
		cardMember,
	))

	agendaText.ParentScroll = content
	studentCheckList.ParentScroll = content

	reloadIncidents()
	_ = currentStudents
	return content
}
