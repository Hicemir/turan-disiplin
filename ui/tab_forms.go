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

	"okulevrak/model"
	"okulevrak/pdf"
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

	updateExtraInputs := func(formName string) {
		extraBox.Objects = nil
		switch {
		case strings.HasPrefix(formName, "1.") || strings.HasPrefix(formName, "2."):
			extraBox.Add(widget.NewLabelWithStyle("İfade / Savunma Metni:", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}))
			extraBox.Add(statementEntry)
		case strings.HasPrefix(formName, "8."):
			f := widget.NewForm(
				widget.NewFormItem("Talep Edilen Ek Süre (Gün)", extraDaysEntry),
				widget.NewFormItem("Ek Süre Talep Nedeni", extraReasonEntry),
			)
			extraBox.Add(f)
		case strings.HasPrefix(formName, "10."):
			f := widget.NewForm(
				widget.NewFormItem("Uyarı Konusu", sozluKonuEntry),
				widget.NewFormItem("Uyaran Öğretmen / Yönetici", sozluOgretmenEntry),
				widget.NewFormItem("Açıklama / Görüşme Özeti", sozluDetayEntry),
			)
			extraBox.Add(f)
		case strings.HasPrefix(formName, "12."):
			f := widget.NewForm(
				widget.NewFormItem("Görüşmeyi Yapan", veliGorusenEntry),
				widget.NewFormItem("Görüşme Konusu", veliKonuEntry),
				widget.NewFormItem("Görüşme Notları", veliNotlarEntry),
				widget.NewFormItem("Alınan Kararlar", veliKararEntry),
			)
			extraBox.Add(f)
		case strings.HasPrefix(formName, "14."):
			f := widget.NewForm(
				widget.NewFormItem("1. Ailevi Durumu", rehber1Entry),
				widget.NewFormItem("2. Psikolojik ve Sosyal Özellikleri", rehber2Entry),
				widget.NewFormItem("3. Olayın Nedenleri", rehber3Entry),
				widget.NewFormItem("4. Rehberlik Önerileri", rehber4Entry),
			)
			extraBox.Add(f)
		}
		extraBox.Refresh()
	}

	docTypeSelect.OnChanged = func(s string) {
		updateExtraInputs(s)
	}

	reloadDropdowns := func() {
		stList, err := state.Store.GetStudents(state.AcademicYear, "", "")
		if err == nil {
			currentStudents = stList
			var opts []string
			for _, s := range stList {
				opts = append(opts, fmt.Sprintf("%s - %s (%s)", s.SchoolNo, s.FullName, s.ClassSection))
			}
			stSelect.Options = opts
			stSelect.Refresh()
			if len(opts) > 0 {
				stSelect.SetSelectedIndex(0)
			}
		}

		incList, err := state.Store.GetIncidents(state.AcademicYear)
		if err == nil {
			currentIncidents = incList
			var opts []string
			for _, inc := range incList {
				opts = append(opts, fmt.Sprintf("#%s - %s", inc.IncidentCode, inc.Title))
			}
			incSelect.Options = opts
			incSelect.Refresh()
			if len(opts) > 0 {
				incSelect.SetSelectedIndex(0)
			}
		}
	}

	stSelect.OnChanged = func(selected string) {
		for _, s := range currentStudents {
			if strings.HasPrefix(selected, s.SchoolNo+" - ") {
				selectedStudent = &s
				break
			}
		}
	}

	incSelect.OnChanged = func(selected string) {
		for _, inc := range currentIncidents {
			prefix := fmt.Sprintf("#%s - ", inc.IncidentCode)
			if strings.HasPrefix(selected, prefix) {
				selectedIncident = &inc
				break
			}
		}
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
			_ = pdf.OpenPDF(outPath)

		case strings.HasPrefix(docType, "2."): // Savunma Tutanağı
			outPath, _ := pdf.GenerateUniquePDFPath("", fmt.Sprintf("Savunma_Tutanagi_%s_%s", stNo, stName))
			err := pdf.GenerateIfadeSavunmaTutanagi(state.SchoolInfo, selectedStudent, selectedIncident, true, statementEntry.Text, outPath)
			if err != nil {
				showError(state.Window, err)
				return
			}
			_ = pdf.OpenPDF(outPath)

		case strings.HasPrefix(docType, "3."): // Bilgi Toplama - Olay
			outPath, _ := pdf.GenerateUniquePDFPath("", fmt.Sprintf("Bilgi_Toplama_Olay_%s", stNo))
			err := pdf.GenerateBilgiToplamaFormu("Olay", state.SchoolInfo, selectedStudent, selectedIncident, outPath)
			if err != nil {
				showError(state.Window, err)
				return
			}
			_ = pdf.OpenPDF(outPath)

		case strings.HasPrefix(docType, "4."): // Bilgi Toplama - Öğretmen
			outPath, _ := pdf.GenerateUniquePDFPath("", fmt.Sprintf("Bilgi_Toplama_Ogretmen_%s", stNo))
			err := pdf.GenerateBilgiToplamaFormu("Öğretmen", state.SchoolInfo, selectedStudent, selectedIncident, outPath)
			if err != nil {
				showError(state.Window, err)
				return
			}
			_ = pdf.OpenPDF(outPath)

		case strings.HasPrefix(docType, "5."): // Bilgi Toplama - Sınıf Rehber Öğretmeni
			outPath, _ := pdf.GenerateUniquePDFPath("", fmt.Sprintf("Bilgi_Toplama_SinifRehber_%s", stNo))
			err := pdf.GenerateBilgiToplamaFormu("Sınıf Rehber Öğretmeni", state.SchoolInfo, selectedStudent, selectedIncident, outPath)
			if err != nil {
				showError(state.Window, err)
				return
			}
			_ = pdf.OpenPDF(outPath)

		case strings.HasPrefix(docType, "6."): // Bilgi Toplama - Veli
			outPath, _ := pdf.GenerateUniquePDFPath("", fmt.Sprintf("Bilgi_Toplama_Veli_%s", stNo))
			err := pdf.GenerateBilgiToplamaFormu("Veli", state.SchoolInfo, selectedStudent, selectedIncident, outPath)
			if err != nil {
				showError(state.Window, err)
				return
			}
			_ = pdf.OpenPDF(outPath)

		case strings.HasPrefix(docType, "7."): // Bilgi Toplama - Öğrenci
			outPath, _ := pdf.GenerateUniquePDFPath("", fmt.Sprintf("Bilgi_Toplama_Ogrenci_%s", stNo))
			err := pdf.GenerateBilgiToplamaFormu("Öğrenci", state.SchoolInfo, selectedStudent, selectedIncident, outPath)
			if err != nil {
				showError(state.Window, err)
				return
			}
			_ = pdf.OpenPDF(outPath)

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
			_ = pdf.OpenPDF(outPath)

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
			_ = pdf.OpenPDF(outPath)

		case strings.HasPrefix(docType, "10."): // Sözlü Uyarı
			outPath, _ := pdf.GenerateUniquePDFPath("", fmt.Sprintf("Sozlu_Uyari_%s", stNo))
			err := pdf.GenerateSozluUyariTutanagi(state.SchoolInfo, selectedStudent, time.Now().Format("02.01.2006"), sozluKonuEntry.Text, sozluDetayEntry.Text, sozluOgretmenEntry.Text, outPath)
			if err != nil {
				showError(state.Window, err)
				return
			}
			_ = pdf.OpenPDF(outPath)

		case strings.HasPrefix(docType, "11."): // Öğrenci Sözleşmesi
			outPath, _ := pdf.GenerateUniquePDFPath("", fmt.Sprintf("Ogrenci_Sozlesmesi_%s", stNo))
			err := pdf.GenerateOgrenciSozlesmesi(state.SchoolInfo, selectedStudent, time.Now().Format("02.01.2006"), outPath)
			if err != nil {
				showError(state.Window, err)
				return
			}
			_ = pdf.OpenPDF(outPath)

		case strings.HasPrefix(docType, "12."): // Veli Görüşme Tutanağı
			outPath, _ := pdf.GenerateUniquePDFPath("", fmt.Sprintf("Veli_Gorusme_Tutanagi_%s", stNo))
			err := pdf.GenerateVeliGorusmeTutanagi(state.SchoolInfo, selectedStudent, selectedStudent.GuardianName, time.Now().Format("02.01.2006"), veliKonuEntry.Text, veliNotlarEntry.Text, veliKararEntry.Text, veliGorusenEntry.Text, false, outPath)
			if err != nil {
				showError(state.Window, err)
				return
			}
			_ = pdf.OpenPDF(outPath)

		case strings.HasPrefix(docType, "13."): // Veli Görüşme Formu (Boş)
			outPath, _ := pdf.GenerateUniquePDFPath("", fmt.Sprintf("Veli_Gorusme_Formu_%s", stNo))
			err := pdf.GenerateVeliGorusmeTutanagi(state.SchoolInfo, selectedStudent, "", "", "", "", "", "", true, outPath)
			if err != nil {
				showError(state.Window, err)
				return
			}
			_ = pdf.OpenPDF(outPath)

		case strings.HasPrefix(docType, "14."): // Rehberlik Görüş Formu
			outPath, _ := pdf.GenerateUniquePDFPath("", fmt.Sprintf("Rehberlik_Gorus_Formu_%s", stNo))
			notes := [4]string{rehber1Entry.Text, rehber2Entry.Text, rehber3Entry.Text, rehber4Entry.Text}
			err := pdf.GenerateRehberlikGorusFormu(state.SchoolInfo, selectedStudent, notes, outPath)
			if err != nil {
				showError(state.Window, err)
				return
			}
			_ = pdf.OpenPDF(outPath)
		}
	})
	pdfBtn.Importance = widget.SuccessImportance

	mainForm := widget.NewForm(
		widget.NewFormItem("Evrak Türü *", docTypeSelect),
		widget.NewFormItem("Öğrenci Seçimi *", stSelect),
		widget.NewFormItem("İlişkili Olay (Varsa)", incSelect),
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

	return container.NewVScroll(contentCard)
}
