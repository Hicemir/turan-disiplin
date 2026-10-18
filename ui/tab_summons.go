package ui

import (
	"fmt"
	"strings"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"okulevrak/model"
	"okulevrak/pdf"
)

// BuildSummonsTab creates the Student & Member Summons Slips UI
func BuildSummonsTab(state *AppState) fyne.CanvasObject {
	var currentIncidents []model.Incident
	var currentStudents []model.Student
	var selectedIncident *model.Incident

	// --- 1. Student Summons Controls ---
	incSelect := widget.NewSelect([]string{}, nil)
	callTypeRadio := widget.NewRadioGroup([]string{"Savunma Çağrı Pusulası", "İfade Çağrı Pusulası"}, nil)
	callTypeRadio.Horizontal = true
	callTypeRadio.SetSelected("Savunma Çağrı Pusulası")

	placeEntry := widget.NewEntry()
	placeEntry.SetText("Disiplin Kurulu Odası")
	dateEntry := widget.NewEntry()
	dateEntry.SetText(time.Now().Format("02.01.2006"))
	timeEntry := widget.NewEntry()
	timeEntry.SetText("10:00")

	// Multi-student selection checklist
	var studentChecks []*widget.Check
	studentCheckContainer := container.NewVBox()

	reloadStudentChecks := func(incCode string) {
		studentCheckContainer.Objects = nil
		studentChecks = nil

		list, err := state.Store.GetIncidentStudents(state.AcademicYear, incCode)
		if err == nil {
			for _, is := range list {
				currIS := is
				chk := widget.NewCheck(fmt.Sprintf("%s - %s (%s)", currIS.SchoolNo, currIS.FullName, currIS.Role), nil)
				chk.SetChecked(true)
				studentChecks = append(studentChecks, chk)
				studentCheckContainer.Add(chk)
			}
		}
		studentCheckContainer.Refresh()
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
		incStudents, _ := state.Store.GetIncidentStudents(state.AcademicYear, selectedIncident.IncidentCode)
		for i, chk := range studentChecks {
			if chk.Checked && i < len(incStudents) {
				is := incStudents[i]
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

		_ = pdf.OpenPDF(outPath)
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

		_ = pdf.OpenPDF(outPath)
	})

	// --- 2. Board Member Summons Controls ---
	sayiEntry := widget.NewEntry()
	sayiEntry.SetText(state.SchoolInfo.OfficialCode + "100")
	konuEntry := widget.NewEntry()
	konuEntry.SetText("Disiplin Kurulu Toplantı Çağrısı")
	uyeTarihEntry := widget.NewEntry()
	uyeTarihEntry.SetText(time.Now().Format("02.01.2006"))
	uyeSaatEntry := widget.NewEntry()
	uyeSaatEntry.SetText("11:00")
	uyeYerEntry := widget.NewEntry()
	uyeYerEntry.SetText("Toplantı Salonu")

	// Agenda Items
	agendaText := widget.NewMultiLineEntry()
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

		_ = pdf.OpenPDF(outPath)
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
		studentCheckContainer,
		container.NewHBox(pdfStudentSummonsBtn, pdfStudentListBtn),
	))

	// Member Section Layout
	memberForm := widget.NewForm(
		widget.NewFormItem("Resmi Yazı Sayısı", sayiEntry),
		widget.NewFormItem("Konu", konuEntry),
		widget.NewFormItem("Toplantı Tarihi", uyeTarihEntry),
		widget.NewFormItem("Toplantı Saati", uyeSaatEntry),
		widget.NewFormItem("Toplantı Yeri", uyeYerEntry),
		widget.NewFormItem("Gündem Maddeleri (Her satıra bir madde)", agendaText),
	)

	cardMember := widget.NewCard("Disiplin / Davranış Kurulu Üyeleri Toplantı Çağrısı", "Gündemli resmi kurul toplantı çağrı pusulası", container.NewVBox(
		memberForm,
		pdfMemberSummonsBtn,
	))

	content := container.NewVScroll(container.NewVBox(
		cardStudent,
		widget.NewSeparator(),
		cardMember,
	))

	reloadIncidents()
	_ = currentStudents
	return content
}
