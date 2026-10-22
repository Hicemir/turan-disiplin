package ui

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"turandisiplin/model"
	"turandisiplin/pdf"
)

// BuildFormsTab creates the Investigation, Counseling & Disciplinary Records UI
func BuildFormsTab(state *AppState) fyne.CanvasObject {
	var currentStudents []model.Student
	var currentIncidents []model.Incident
	var selectedStudent *model.Student
	var selectedIncident *model.Incident

	stSelect := NewTurkishSelect([]string{}, nil)
	incSelect := NewTurkishSelect([]string{}, nil)

	docTypeSelect := NewTurkishSelect([]string{
		"1. İfade Tutanağı",
		"2. Savunma Tutanağı",
		"3. Bilgi Toplama Formu - Olay",
		"4. Bilgi Toplama Formu - Öğretmen",
		"5. Bilgi Toplama Formu - Sınıf Rehber Öğretmeni",
		"6. Bilgi Toplama Formu - Veli",
		"7. Bilgi Toplama Formu - Öğrenci",
		"8. Ek Süre Talep Formu (7 Gün)",
		"9. Şikayet Formu",
		"10. Sözlü Uyarı Tutanağı",
		"11. Öğrenci Sözleşmesi",
		"12. Veli Görüşme Tutanağı",
		"13. Veli Görüşme Formu (Boş)",
		"14. Rehberlik ve Psikolojik Danışma Görüş Formu",
	}, nil)
	docTypeSelect.SetSelected("1. İfade Tutanağı")

	// Dynamic inputs container depending on form selected
	extraBox := container.NewVBox()

	// Reusable inputs
	statementEntry := NewSmartMultiLineEntry()
	statementEntry.SetPlaceHolder("Yazılı ifade / savunma metni (boş bırakılırsa elle yazmak için çizgili alan oluşturulur)...")

	extraDaysEntry := NewSingleLineEntry()
	extraDaysEntry.SetText("7")
	extraReasonEntry := NewSingleLineEntry()
	extraReasonEntry.SetText("Delillerin ve tanık ifadelerinin toplanmasının tamamlanamaması")

	sozluKonuEntry := NewSingleLineEntry()
	sozluKonuEntry.SetText("Kılık kıyafet ve ders düzeni kurallarına uyulmaması")
	sozluDetayEntry := NewSmartMultiLineEntry()
	sozluDetayEntry.SetPlaceHolder("Öğrenciye yapılan sözlü uyarının detayları...")
	sozluOgretmenEntry := NewSingleLineEntry()
	sozluOgretmenEntry.SetText(state.SchoolInfo.PrincipalName)

	veliGorusenEntry := NewSingleLineEntry()
	veliGorusenEntry.SetText(state.SchoolInfo.BoardChair)
	veliKonuEntry := NewSingleLineEntry()
	veliKonuEntry.SetText("Öğrencinin devamsızlık ve ders içi davranışları")
	veliNotlarEntry := NewSmartMultiLineEntry()
	veliKararEntry := NewSmartMultiLineEntry()

	rehber1Entry := NewSmartMultiLineEntry()
	rehber1Entry.SetPlaceHolder("1- Öğrencinin Ailevi Durumu...")
	rehber2Entry := NewSmartMultiLineEntry()
	rehber2Entry.SetPlaceHolder("2- Öğrencinin Psikolojik ve Sosyal Özellikleri...")
	rehber3Entry := NewSmartMultiLineEntry()
	rehber3Entry.SetPlaceHolder("3- Olayın Nedenleri Hakkındaki Görüşler...")
	rehber4Entry := NewSmartMultiLineEntry()
	rehber4Entry.SetPlaceHolder("4- Uygulanabilecek Yaptırım veya Rehberlik Önerileri...")

	sec1Box := container.NewVBox(
		widget.NewLabelWithStyle("İfade / Savunma Metni:", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		statementEntry,
	)
	sec8Form := widget.NewForm(
		widget.NewFormItem("Talep Edilen Ek Süre (Gün)", extraDaysEntry),
		widget.NewFormItem("Ek Süre Talep Nedeni", extraReasonEntry),
	)
	sec10Form := widget.NewForm(
		widget.NewFormItem("Uyarı Konusu", sozluKonuEntry),
		widget.NewFormItem("Uyaran Öğretmen / Yönetici", sozluOgretmenEntry),
		widget.NewFormItem("Açıklama / Görüşme Özeti", sozluDetayEntry),
	)
	sec12Form := widget.NewForm(
		widget.NewFormItem("Görüşmeyi Yapan", veliGorusenEntry),
		widget.NewFormItem("Görüşme Konusu", veliKonuEntry),
		widget.NewFormItem("Görüşme Notları", veliNotlarEntry),
		widget.NewFormItem("Alınan Kararlar", veliKararEntry),
	)
	sec14Form := widget.NewForm(
		widget.NewFormItem("1. Ailevi Durumu", rehber1Entry),
		widget.NewFormItem("2. Psikolojik ve Sosyal Özellikleri", rehber2Entry),
		widget.NewFormItem("3. Olayın Nedenleri", rehber3Entry),
		widget.NewFormItem("4. Rehberlik Önerileri", rehber4Entry),
	)

	extraBox.Add(sec1Box)
	extraBox.Add(sec8Form)
	extraBox.Add(sec10Form)
	extraBox.Add(sec12Form)
	extraBox.Add(sec14Form)

	updateExtraInputs := func(formName string) {
		sec1Box.Hide()
		sec8Form.Hide()
		sec10Form.Hide()
		sec12Form.Hide()
		sec14Form.Hide()
		switch {
		case strings.HasPrefix(formName, "1.") || strings.HasPrefix(formName, "2."):
			sec1Box.Show()
		case strings.HasPrefix(formName, "8."):
			sec8Form.Show()
		case strings.HasPrefix(formName, "10."):
			sec10Form.Show()
		case strings.HasPrefix(formName, "12."):
			sec12Form.Show()
		case strings.HasPrefix(formName, "14."):
			sec14Form.Show()
		}
	}

	docTypeSelect.OnChanged = func(s string) {
		updateExtraInputs(s)
	}

	const noIncidentOption = "(Olay yok - öğrenciyi ara)"
	const maxSearchResults = 50

	// Student search box: only used when no related incident is selected.
	studentSearch := NewSingleLineEntry()
	studentSearch.SetPlaceHolder("Ad soyad veya okul no yazarak arayın...")
	studentSearchItem := widget.NewFormItem("Öğrenci Ara", studentSearch)

	setStudentOptions := func(list []model.Student, placeHolder string) {
		currentStudents = list
		selectedStudent = nil
		var opts []string
		for _, s := range list {
			opts = append(opts, fmt.Sprintf("%s - %s (%s)", s.SchoolNo, s.FullName, s.ClassSection))
		}
		stSelect.Options = opts
		stSelect.PlaceHolder = placeHolder
		stSelect.ClearSelected()
		stSelect.Refresh()
		if len(opts) == 1 || (len(opts) > 0 && selectedIncident != nil) {
			stSelect.SetSelectedIndex(0)
		}
	}

	// Shows students of the selected incident, or search results when none is chosen.
	refreshStudents := func() {
		if selectedIncident != nil {
			studentSearchItem.HintText = ""
			studentSearch.Hide()
			incStudents, err := state.Store.GetIncidentStudents(state.AcademicYear, selectedIncident.IncidentCode)
			if err != nil {
				setStudentOptions(nil, "(Seçiniz)")
				return
			}
			all, _ := state.Store.GetStudents(state.AcademicYear, "", "")
			byNo := make(map[string]model.Student, len(all))
			for _, s := range all {
				byNo[s.SchoolNo] = s
			}
			var list []model.Student
			seen := map[string]bool{}
			for _, is := range incStudents {
				if seen[is.SchoolNo] {
					continue
				}
				seen[is.SchoolNo] = true
				if st, ok := byNo[is.SchoolNo]; ok {
					list = append(list, st)
				} else {
					list = append(list, model.Student{
						AcademicYear: is.AcademicYear, SchoolNo: is.SchoolNo, FullName: is.FullName,
						ClassSection: is.ClassSection, Gender: is.Gender, TCNo: is.TCNo,
						BirthPlace: is.BirthPlace, BirthDate: is.BirthDate, FatherName: is.FatherName,
						MotherName: is.MotherName, GuardianName: is.GuardianName,
						GuardianPhone: is.GuardianPhone, GuardianAddr: is.GuardianAddr,
						BoardingStatus: is.BoardingStatus,
					})
				}
			}
			ph := "(Seçiniz)"
			if len(list) == 0 {
				ph = "(Bu olaya kayıtlı öğrenci yok)"
			}
			setStudentOptions(list, ph)
			return
		}

		studentSearch.Show()
		query := foldTR(strings.TrimSpace(studentSearch.Text))
		if query == "" {
			setStudentOptions(nil, "(Önce yukarıdan öğrenci arayın)")
			return
		}
		all, err := state.Store.GetStudents(state.AcademicYear, "", "")
		if err != nil {
			setStudentOptions(nil, "(Seçiniz)")
			return
		}
		var list []model.Student
		for _, s := range all {
			hay := foldTR(s.SchoolNo + " " + s.FullName + " " + s.ClassSection)
			if strings.Contains(hay, query) {
				list = append(list, s)
				if len(list) >= maxSearchResults {
					break
				}
			}
		}
		ph := "(Seçiniz)"
		if len(list) == 0 {
			ph = "(Sonuç bulunamadı)"
		}
		setStudentOptions(list, ph)
	}

	studentSearch.OnChanged = func(string) { refreshStudents() }

	reloadDropdowns := func() {
		incList, err := state.Store.GetIncidents(state.AcademicYear)
		if err == nil {
			currentIncidents = incList
			opts := []string{noIncidentOption}
			for _, inc := range incList {
				opts = append(opts, fmt.Sprintf("#%s - %s", inc.IncidentCode, inc.Title))
			}
			incSelect.Options = opts
			incSelect.Refresh()
		}
		selectedIncident = nil
		incSelect.SetSelected(noIncidentOption)
		refreshStudents()
	}

	stSelect.OnChanged = func(selected string) {
		for i := range currentStudents {
			s := currentStudents[i]
			if strings.HasPrefix(selected, s.SchoolNo+" - ") {
				selectedStudent = &s
				break
			}
		}
	}

	incSelect.OnChanged = func(selected string) {
		selectedIncident = nil
		for i := range currentIncidents {
			inc := currentIncidents[i]
			prefix := fmt.Sprintf("#%s - ", inc.IncidentCode)
			if strings.HasPrefix(selected, prefix) {
				selectedIncident = &inc
				break
			}
		}
		refreshStudents()
	}

	// Big Action Button: "PDF Oluştur ve Aç"
	pdfBtn := widget.NewButtonWithIcon("PDF OLUŞTUR VE AÇ", theme.DocumentIcon(), func() {
		if selectedStudent == nil {
			showError(state.Window, fmt.Errorf("Lütfen listeden bir öğrenci seçiniz"))
			return
		}

		docType := docTypeSelect.Selected
		stNo := selectedStudent.SchoolNo
		stName := strings.ReplaceAll(selectedStudent.FullName, " ", "_")

		switch {
		case strings.HasPrefix(docType, "1."): // İfade Tutanağı
			outPath, _ := pdf.GenerateUniquePDFPath("", fmt.Sprintf("Ifade_Tutanagi_%s_%s", stNo, stName))
			err := pdf.GenerateIfadeSavunmaTutanagi(state.SchoolInfo, selectedStudent, selectedIncident, false, statementEntry.Text, outPath)
			if err != nil {
				showError(state.Window, err)
				return
			}
			openPDF(state, outPath)

		case strings.HasPrefix(docType, "2."): // Savunma Tutanağı
			outPath, _ := pdf.GenerateUniquePDFPath("", fmt.Sprintf("Savunma_Tutanagi_%s_%s", stNo, stName))
			err := pdf.GenerateIfadeSavunmaTutanagi(state.SchoolInfo, selectedStudent, selectedIncident, true, statementEntry.Text, outPath)
			if err != nil {
				showError(state.Window, err)
				return
			}
			openPDF(state, outPath)

		case strings.HasPrefix(docType, "3."): // Bilgi Toplama - Olay
			outPath, _ := pdf.GenerateUniquePDFPath("", fmt.Sprintf("Bilgi_Toplama_Olay_%s", stNo))
			err := pdf.GenerateBilgiToplamaFormu("Olay", state.SchoolInfo, selectedStudent, selectedIncident, outPath)
			if err != nil {
				showError(state.Window, err)
				return
			}
			openPDF(state, outPath)

		case strings.HasPrefix(docType, "4."): // Bilgi Toplama - Öğretmen
			outPath, _ := pdf.GenerateUniquePDFPath("", fmt.Sprintf("Bilgi_Toplama_Ogretmen_%s", stNo))
			err := pdf.GenerateBilgiToplamaFormu("Öğretmen", state.SchoolInfo, selectedStudent, selectedIncident, outPath)
			if err != nil {
				showError(state.Window, err)
				return
			}
			openPDF(state, outPath)

		case strings.HasPrefix(docType, "5."): // Bilgi Toplama - Sınıf Rehber Öğretmeni
			outPath, _ := pdf.GenerateUniquePDFPath("", fmt.Sprintf("Bilgi_Toplama_SinifRehber_%s", stNo))
			err := pdf.GenerateBilgiToplamaFormu("Sınıf Rehber Öğretmeni", state.SchoolInfo, selectedStudent, selectedIncident, outPath)
			if err != nil {
				showError(state.Window, err)
				return
			}
			openPDF(state, outPath)

		case strings.HasPrefix(docType, "6."): // Bilgi Toplama - Veli
			outPath, _ := pdf.GenerateUniquePDFPath("", fmt.Sprintf("Bilgi_Toplama_Veli_%s", stNo))
			err := pdf.GenerateBilgiToplamaFormu("Veli", state.SchoolInfo, selectedStudent, selectedIncident, outPath)
			if err != nil {
				showError(state.Window, err)
				return
			}
			openPDF(state, outPath)

		case strings.HasPrefix(docType, "7."): // Bilgi Toplama - Öğrenci
			outPath, _ := pdf.GenerateUniquePDFPath("", fmt.Sprintf("Bilgi_Toplama_Ogrenci_%s", stNo))
			err := pdf.GenerateBilgiToplamaFormu("Öğrenci", state.SchoolInfo, selectedStudent, selectedIncident, outPath)
			if err != nil {
				showError(state.Window, err)
				return
			}
			openPDF(state, outPath)

		case strings.HasPrefix(docType, "8."): // Ek Süre Talep
			days, _ := strconv.Atoi(extraDaysEntry.Text)
			if days <= 0 {
				days = 7
			}
			outPath, _ := pdf.GenerateUniquePDFPath("", fmt.Sprintf("Ek_Sure_Talep_%s", stNo))
			err := pdf.GenerateEkSureTalepFormu(state.SchoolInfo, selectedStudent, selectedIncident, days, extraReasonEntry.Text, outPath)
			if err != nil {
				showError(state.Window, err)
				return
			}
			openPDF(state, outPath)

		case strings.HasPrefix(docType, "9."): // Şikayet Formu
			outPath, _ := pdf.GenerateUniquePDFPath("", fmt.Sprintf("Sikayet_Formu_%s", stNo))
			reportedStr := fmt.Sprintf("%s (%s - %s)", selectedStudent.FullName, selectedStudent.ClassSection, selectedStudent.SchoolNo)
			complainantStr := ""
			incDateStr := time.Now().Format("02.01.2006")
			incTopicStr := ""
			incDescStr := ""
			if selectedIncident != nil {
				complainantStr = selectedIncident.Complainant
				incDateStr = fmt.Sprintf("%s - %s", selectedIncident.Date, selectedIncident.Time)
				incTopicStr = selectedIncident.Title
				incDescStr = selectedIncident.Summary
			}
			err := pdf.GenerateSikayetFormu(state.SchoolInfo, complainantStr, reportedStr, "—", "—", incDateStr, incTopicStr, incDescStr, "", outPath)
			if err != nil {
				showError(state.Window, err)
				return
			}
			openPDF(state, outPath)

		case strings.HasPrefix(docType, "10."): // Sözlü Uyarı
			outPath, _ := pdf.GenerateUniquePDFPath("", fmt.Sprintf("Sozlu_Uyari_%s", stNo))
			err := pdf.GenerateSozluUyariTutanagi(state.SchoolInfo, selectedStudent, time.Now().Format("02.01.2006"), sozluKonuEntry.Text, sozluDetayEntry.Text, sozluOgretmenEntry.Text, outPath)
			if err != nil {
				showError(state.Window, err)
				return
			}
			openPDF(state, outPath)

		case strings.HasPrefix(docType, "11."): // Öğrenci Sözleşmesi
			outPath, _ := pdf.GenerateUniquePDFPath("", fmt.Sprintf("Ogrenci_Sozlesmesi_%s", stNo))
			err := pdf.GenerateOgrenciSozlesmesi(state.SchoolInfo, selectedStudent, time.Now().Format("02.01.2006"), outPath)
			if err != nil {
				showError(state.Window, err)
				return
			}
			openPDF(state, outPath)

		case strings.HasPrefix(docType, "12."): // Veli Görüşme Tutanağı
			outPath, _ := pdf.GenerateUniquePDFPath("", fmt.Sprintf("Veli_Gorusme_Tutanagi_%s", stNo))
			gName := model.EffectiveGuardianName(selectedStudent.FullName, selectedStudent.FatherName, selectedStudent.GuardianName)
			err := pdf.GenerateVeliGorusmeTutanagi(state.SchoolInfo, selectedStudent, gName, time.Now().Format("02.01.2006"), veliKonuEntry.Text, veliNotlarEntry.Text, veliKararEntry.Text, veliGorusenEntry.Text, false, outPath)
			if err != nil {
				showError(state.Window, err)
				return
			}
			openPDF(state, outPath)

		case strings.HasPrefix(docType, "13."): // Veli Görüşme Formu (Boş)
			outPath, _ := pdf.GenerateUniquePDFPath("", fmt.Sprintf("Veli_Gorusme_Formu_%s", stNo))
			err := pdf.GenerateVeliGorusmeTutanagi(state.SchoolInfo, selectedStudent, "", "", "", "", "", "", true, outPath)
			if err != nil {
				showError(state.Window, err)
				return
			}
			openPDF(state, outPath)

		case strings.HasPrefix(docType, "14."): // Rehberlik Görüş Formu
			outPath, _ := pdf.GenerateUniquePDFPath("", fmt.Sprintf("Rehberlik_Gorus_Formu_%s", stNo))
			notes := [4]string{rehber1Entry.Text, rehber2Entry.Text, rehber3Entry.Text, rehber4Entry.Text}
			err := pdf.GenerateRehberlikGorusFormu(state.SchoolInfo, selectedStudent, notes, outPath)
			if err != nil {
				showError(state.Window, err)
				return
			}
			openPDF(state, outPath)
		}
	})
	pdfBtn.Importance = widget.SuccessImportance

	mainForm := widget.NewForm(
		widget.NewFormItem("Evrak Türü *", docTypeSelect),
		widget.NewFormItem("İlişkili Olay (Varsa)", incSelect),
		studentSearchItem,
		widget.NewFormItem("Öğrenci Seçimi *", stSelect),
	)

	contentCard := widget.NewCard("Resmi Tutanak ve Form Seçimi", "Oluşturmak istediğiniz resmi evrakı ve öğrenciyi seçiniz", container.NewVBox(
		mainForm,
		widget.NewSeparator(),
		extraBox,
		widget.NewSeparator(),
		pdfBtn,
	))

	reloadDropdowns()
	updateExtraInputs(docTypeSelect.Selected)

	pageScroll := container.NewVScroll(contentCard)
	statementEntry.ParentScroll = pageScroll
	sozluDetayEntry.ParentScroll = pageScroll
	veliNotlarEntry.ParentScroll = pageScroll
	veliKararEntry.ParentScroll = pageScroll
	rehber1Entry.ParentScroll = pageScroll
	rehber2Entry.ParentScroll = pageScroll
	rehber3Entry.ParentScroll = pageScroll
	rehber4Entry.ParentScroll = pageScroll

	return pageScroll
}

// foldTR lowercases text and removes Turkish diacritics so searches work
// regardless of case or whether the user typed "ı/i", "ş/s", etc.
func foldTR(s string) string {
	var sb strings.Builder
	for _, r := range s {
		switch r {
		case 'ı', 'İ', 'I', 'i':
			sb.WriteRune('i')
		case 'ğ', 'Ğ':
			sb.WriteRune('g')
		case 'ü', 'Ü':
			sb.WriteRune('u')
		case 'ş', 'Ş':
			sb.WriteRune('s')
		case 'ö', 'Ö':
			sb.WriteRune('o')
		case 'ç', 'Ç':
			sb.WriteRune('c')
		default:
			sb.WriteString(strings.ToLower(string(r)))
		}
	}
	return sb.String()
}
