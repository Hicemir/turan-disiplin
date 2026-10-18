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

func normalizeTurkish(s string) string {
	var sb strings.Builder
	for _, r := range s {
		switch r {
		case 'ı', 'I', 'İ', 'i':
			sb.WriteRune('I')
		case 'ğ', 'Ğ':
			sb.WriteRune('G')
		case 'ü', 'Ü':
			sb.WriteRune('U')
		case 'ş', 'Ş':
			sb.WriteRune('S')
		case 'ö', 'Ö':
			sb.WriteRune('O')
		case 'ç', 'Ç':
			sb.WriteRune('C')
		default:
			if r >= 'a' && r <= 'z' {
				sb.WriteRune(r - ('a' - 'A'))
			} else {
				sb.WriteRune(r)
			}
		}
	}
	return sb.String()
}

// GetPenaltyPointsDeduction calculates official MEB points deduction for a given disciplinary penalty.
func GetPenaltyPointsDeduction(penalty string) int {
	p := normalizeTurkish(penalty)
	if strings.Contains(p, "ORGUN") && strings.Contains(p, "DISI") {
		return 80
	} else if strings.Contains(p, "OKUL DEGISTIRME") {
		return 40
	} else if strings.Contains(p, "UZAKLASTIRMA") {
		return 20
	} else if strings.Contains(p, "KINAMA") {
		return 10
	}
	// Sözlü Uyarma, Yazılı Uyarma, Uyarma
	return 0
}

// BuildVerdictTab creates the Disciplinary Board Verdict UI
func BuildVerdictTab(state *AppState) fyne.CanvasObject {
	var currentIncidents []model.Incident
	var currentIncStudents []model.IncidentStudent
	var selectedIncStudent *model.IncidentStudent
	var selectedIncident *model.Incident

	// Selectors
	incSelect := NewTurkishSelect([]string{}, nil)
	stSelect := NewTurkishSelect([]string{}, nil)

	// Verdict Form Fields
	verdictNoEntry := NewSingleLineEntry()
	verdictDateEntry := NewSingleLineEntry()
	verdictDateEntry.SetText(time.Now().Format("02.01.2006"))

	offensePlaceEntry := NewSingleLineEntry()
	offenseDateEntry := NewSingleLineEntry()
	offenseCategorySelect := NewTurkishSelect([]string{
		"DİSİPLİNİ BOZUCU DAVRANIŞLAR",
		"ŞİDDET VE SALDIRGANLIK",
		"DEVAMSIZLIK VE OKUL DÜZENİ",
		"BİLİŞİM VE ELEKTRONİK CİHAZ KULLANIMI",
		"ZARAR VERME / TAHRİBAT",
		"DİĞER",
	}, nil)
	offenseReasonEntry := NewSingleLineEntry()
	offenseReasonEntry.SetText("Düşüncesizlik ve kurallara uymama")

	witnessEntry := NewSmartMultiLineEntry()
	evidenceEntry := NewSmartMultiLineEntry()
	mitigatingSelect := NewTurkishSelect([]string{
		"Daha önce disiplin cezası almamış olması ve pişmanlığı",
		"Olaydaki samimi itirafı ve davranışı düzeltme gayreti",
		"Tahrik altında kalmış olması",
		"Herhangi bir hafifletici neden bulunmamaktadır",
	}, nil)
	mitigatingSelect.SetSelected("Daha önce disiplin cezası almamış olması ve pişmanlığı")

	aggravatingSelect := NewTurkishSelect([]string{
		"Herhangi bir ağırlaştırıcı neden bulunmamaktadır",
		"Benzer fiilin daha önce tekrarlanmış olması",
		"Toplu olarak veya organize şekilde işlenmiş olması",
		"Okul düzenini ağır biçimde aksatmış olması",
	}, nil)
	aggravatingSelect.SetSelected("Herhangi bir ağırlaştırıcı neden bulunmamaktadır")

	opinionSelect := NewTurkishSelect([]string{
		"SUÇLU OLDUĞU YÖNÜNDE",
		"SUÇSUZ OLDUĞU YÖNÜNDE",
		"CEZA TAYİNİNE YER OLMADIĞI YÖNÜNDE",
	}, nil)
	opinionSelect.SetSelected("SUÇLU OLDUĞU YÖNÜNDE")

	// Regulation Clause picker
	clauses := state.Yonetmelik.Clauses
	var clauseOptions []string
	clauseMap := make(map[string]model.Clause)

	schoolTypeUpper := strings.ToUpper(state.SchoolInfo.SchoolType)
	for _, c := range clauses {
		if c.SchoolType == "" || strings.EqualFold(c.SchoolType, schoolTypeUpper) {
			display := fmt.Sprintf("[%s] %s - %s", c.Penalty, c.Article, c.Description)
			clauseOptions = append(clauseOptions, display)
			clauseMap[display] = c
		}
	}

	clauseSelect := NewTurkishSelect(clauseOptions, nil)

	clauseCodeEntry := NewSingleLineEntry()
	clauseDescEntry := NewSmartMultiLineEntry()


	penaltySelect := NewTurkishSelect([]string{
		"Sözlü Uyarma",
		"Yazılı Uyarma",
		"Kınama",
		"Okuldan Kısa Süreli Uzaklaştırma (1-5 Gün)",
		"Okul Değiştirme",
		"Örgün Eğitim Dışına Çıkarma",
	}, nil)

	suspensionDaysEntry := NewSingleLineEntry()
	suspensionDaysEntry.SetPlaceHolder("1-5 gün")
	suspensionDatesEntry := NewSingleLineEntry()
	suspensionDatesEntry.SetPlaceHolder("Örn: 10.11.2026 - 12.11.2026")

	deductedPointsEntry := NewSingleLineEntry()
	deductedPointsEntry.SetText("10")
	remainingPointsEntry := NewSingleLineEntry()
	remainingPointsEntry.SetText("90")

	// Auto-calculate remaining score on deducted change
	deductedPointsEntry.OnChanged = func(s string) {
		if d, err := strconv.Atoi(s); err == nil {
			rem := 100 - d
			if rem < 0 {
				rem = 0
			}
			remainingPointsEntry.SetText(fmt.Sprintf("%d", rem))
		}
	}

	verdictTextEntry := NewSmartMultiLineEntry()

	// Substitute Member Checkboxes
	subChairCheck := widget.NewCheck("Başkan yerine Yedek Başkan imzalayacak (Taraf / Mazeretli)", nil)
	subM1Check := widget.NewCheck("Üye 1 yerine Yedek Üye imzalayacak", nil)
	subM2Check := widget.NewCheck("Üye 2 yerine Yedek Üye imzalayacak", nil)
	subM3Check := widget.NewCheck("Üye 3 yerine Yedek Üye imzalayacak", nil)
	subM4Check := widget.NewCheck("Üye 4 yerine Yedek Üye imzalayacak", nil)

	generateVerdictText := func() {
		if selectedIncStudent == nil || selectedIncident == nil {
			return
		}

		st := selectedIncStudent
		inc := selectedIncident

		suspText := ""
		if suspensionDaysEntry.Text != "" {
			suspText = fmt.Sprintf(" ve %s gün süreyle okuldan uzaklaştırılmasına (%s)", suspensionDaysEntry.Text, suspensionDatesEntry.Text)
		}

		pointText := ""
		if strings.ToUpper(state.SchoolInfo.SchoolType) == "LİSE" {
			pointText = fmt.Sprintf(", davranış puanından %s puan indirilerek kalan puanının %s olarak belirlenmesine", deductedPointsEntry.Text, remainingPointsEntry.Text)
		}

		text := fmt.Sprintf("Yukarıda açık kimliği yazılı okulumuz %s sınıfı %s numaralı öğrencisi %s'nin, %s tarihinde meydana gelen '%s' olayındaki olumsuz eylemi nedeniyle; savunması, tanık beyanları ve toplanan deliller değerlendirilerek Milli Eğitim Bakanlığı Yönetmeliğinin %s maddesi gereğince '%s' cezası ile cezalandırılmasına%s%s; işbu karara tebliğ tarihinden itibaren 5 iş günü içinde itiraz edilebileceğine oy birliği ile karar verildi.",
			st.ClassSection, st.SchoolNo, st.FullName, inc.Date, inc.Title, clauseCodeEntry.Text, penaltySelect.Selected, suspText, pointText)

		verdictTextEntry.SetText(text)
	}

	penaltySelect.OnChanged = func(selected string) {
		ded := GetPenaltyPointsDeduction(selected)
		deductedPointsEntry.SetText(fmt.Sprintf("%d", ded))
		generateVerdictText()
	}

	autoFillBtn := widget.NewButtonWithIcon("Metni Otomatik Oluştur", theme.DocumentCreateIcon(), func() {
		generateVerdictText()
	})

	clauseSelect.OnChanged = func(selected string) {
		if c, ok := clauseMap[selected]; ok {
			clauseCodeEntry.SetText(c.Article)
			clauseDescEntry.SetText(c.Description)
			penaltySelect.SetSelected(c.Penalty)
			ded := GetPenaltyPointsDeduction(c.Penalty)
			deductedPointsEntry.SetText(fmt.Sprintf("%d", ded))
			generateVerdictText()
		}
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

	reloadStudentsForIncident := func(incCode string) {
		list, err := state.Store.GetIncidentStudents(state.AcademicYear, incCode)
		if err == nil {
			currentIncStudents = list
			var opts []string
			for _, is := range list {
				opts = append(opts, fmt.Sprintf("%s - %s (%s)", is.SchoolNo, is.FullName, is.Role))
			}
			stSelect.Options = opts
			stSelect.Refresh()
			if len(opts) > 0 {
				stSelect.SetSelectedIndex(0)
			}
		}
	}

	incSelect.OnChanged = func(selected string) {
		if selected == "" {
			return
		}
		for _, inc := range currentIncidents {
			prefix := fmt.Sprintf("#%s - ", inc.IncidentCode)
			if strings.HasPrefix(selected, prefix) {
				selectedIncident = &inc
				offensePlaceEntry.SetText(inc.Location)
				offenseDateEntry.SetText(fmt.Sprintf("%s - %s", inc.Date, inc.Time))
				reloadStudentsForIncident(inc.IncidentCode)
				break
			}
		}
	}

	populateVerdict := func(is *model.IncidentStudent) {
		selectedIncStudent = is
		if is.VerdictNo != "" {
			verdictNoEntry.SetText(is.VerdictNo)
		} else {
			verdictNoEntry.SetText(state.Store.GetNextVerdictNo(state.AcademicYear))
		}

		if is.VerdictDate != "" {
			verdictDateEntry.SetText(is.VerdictDate)
		}

		if is.OffensePlace != "" {
			offensePlaceEntry.SetText(is.OffensePlace)
		}
		if is.OffenseCategory != "" {
			offenseCategorySelect.SetSelected(is.OffenseCategory)
		}
		if is.OffenseReason != "" {
			offenseReasonEntry.SetText(is.OffenseReason)
		}
		witnessEntry.SetText(is.WitnessSummary)
		evidenceEntry.SetText(is.Evidence)
		if is.MitigatingFactors != "" {
			mitigatingSelect.SetSelected(is.MitigatingFactors)
		}
		if is.AggravatingFactors != "" {
			aggravatingSelect.SetSelected(is.AggravatingFactors)
		}
		if is.BoardOpinion != "" {
			opinionSelect.SetSelected(is.BoardOpinion)
		}
		if is.ClauseCode != "" {
			clauseCodeEntry.SetText(is.ClauseCode)
		}
		if is.ClauseDescription != "" {
			clauseDescEntry.SetText(is.ClauseDescription)
		}
		if is.ImposedPenalty != "" {
			penaltySelect.SetSelected(is.ImposedPenalty)
		}
		if is.SuspensionDays != "" {
			suspensionDaysEntry.SetText(is.SuspensionDays)
		}
		if is.SuspensionDates != "" {
			suspensionDatesEntry.SetText(is.SuspensionDates)
		}
		if is.DeductedPoints > 0 {
			deductedPointsEntry.SetText(fmt.Sprintf("%d", is.DeductedPoints))
			remainingPointsEntry.SetText(fmt.Sprintf("%d", is.RemainingPoints))
		}
		if is.VerdictText != "" {
			verdictTextEntry.SetText(is.VerdictText)
		} else {
			generateVerdictText()
		}

		// Conflict of interest check
		if selectedIncident != nil && strings.Contains(strings.ToUpper(selectedIncident.Complainant), strings.ToUpper(state.SchoolInfo.BoardChair)) && state.SchoolInfo.BoardChair != "" {
			subChairCheck.SetChecked(true)
		}
	}

	stSelect.OnChanged = func(selected string) {
		if selected == "" {
			return
		}
		for _, is := range currentIncStudents {
			if strings.HasPrefix(selected, is.SchoolNo+" - ") {
				populateVerdict(&is)
				break
			}
		}
	}


	// Save Verdict Button
	saveVerdictBtn := widget.NewButtonWithIcon("Kararı Kaydet", theme.DocumentSaveIcon(), func() {
		if selectedIncStudent == nil {
			showError(state.Window, fmt.Errorf("Lütfen önce bir öğrenci seçiniz"))
			return
		}

		vNo := strings.TrimSpace(verdictNoEntry.Text)
		if vNo == "" {
			vNo = state.Store.GetNextVerdictNo(state.AcademicYear)
			verdictNoEntry.SetText(vNo)
		}

		ded, _ := strconv.Atoi(deductedPointsEntry.Text)
		rem, _ := strconv.Atoi(remainingPointsEntry.Text)

		selectedIncStudent.VerdictNo = vNo
		selectedIncStudent.VerdictDate = verdictDateEntry.Text
		selectedIncStudent.OffensePlace = offensePlaceEntry.Text
		selectedIncStudent.OffenseDateTime = offenseDateEntry.Text
		selectedIncStudent.OffenseCategory = offenseCategorySelect.Selected
		selectedIncStudent.OffenseReason = offenseReasonEntry.Text
		selectedIncStudent.WitnessSummary = witnessEntry.Text
		selectedIncStudent.Evidence = evidenceEntry.Text
		selectedIncStudent.MitigatingFactors = mitigatingSelect.Selected
		selectedIncStudent.AggravatingFactors = aggravatingSelect.Selected
		selectedIncStudent.BoardOpinion = opinionSelect.Selected
		selectedIncStudent.ClauseCode = clauseCodeEntry.Text
		selectedIncStudent.ClauseDescription = clauseDescEntry.Text
		selectedIncStudent.ImposedPenalty = penaltySelect.Selected
		selectedIncStudent.SuspensionDays = suspensionDaysEntry.Text
		selectedIncStudent.SuspensionDates = suspensionDatesEntry.Text
		selectedIncStudent.DeductedPoints = ded
		selectedIncStudent.RemainingPoints = rem
		selectedIncStudent.VerdictText = verdictTextEntry.Text

		err := state.Store.SaveIncidentStudent(selectedIncStudent)
		if err != nil {
			showError(state.Window, err)
			return
		}

		// Also update master student behavior score
		st, err := state.Store.GetStudentByNo(state.AcademicYear, selectedIncStudent.SchoolNo)
		if err == nil && st != nil {
			st.BehaviorScore = rem
			_ = state.Store.SaveStudent(st)
		}

		showInfo(state.Window, "Başarılı", fmt.Sprintf("Karar No %s başarıyla kaydedildi.", vNo))
		state.NotifyChange()
	})
	saveVerdictBtn.Importance = widget.HighImportance

	// Prominent "PDF Oluştur ve Aç" Button
	pdfBtn := widget.NewButtonWithIcon("PDF OLUŞTUR VE AÇ (EK-1 / EK-10)", theme.DocumentIcon(), func() {
		if selectedIncStudent == nil || selectedIncident == nil {
			showError(state.Window, fmt.Errorf("Lütfen önce olay ve öğrenci seçiniz"))
			return
		}

		if selectedIncStudent.VerdictNo == "" {
			selectedIncStudent.VerdictNo = verdictNoEntry.Text
		}

		vNo, _ := strconv.Atoi(verdictNoEntry.Text)
		suspDays, _ := strconv.Atoi(suspensionDaysEntry.Text)
		ded, _ := strconv.Atoi(deductedPointsEntry.Text)
		rem, _ := strconv.Atoi(remainingPointsEntry.Text)

		verdict := &model.Verdict{
			DecisionNo:         vNo,
			DecisionDate:       verdictDateEntry.Text,
			OffensePlace:       offensePlaceEntry.Text,
			OffenseDateTime:    offenseDateEntry.Text,
			OffenseType:        offenseCategorySelect.Selected,
			OffenseReason:      offenseReasonEntry.Text,
			WitnessStatements:  witnessEntry.Text,
			EvidenceDetails:    evidenceEntry.Text,
			MitigatingFactors:  mitigatingSelect.Selected,
			AggravatingFactors: aggravatingSelect.Selected,
			BoardOpinion:       opinionSelect.Selected,
			RegClause:          clauseCodeEntry.Text,
			RegClauseDesc:      clauseDescEntry.Text,
			PenaltyName:        penaltySelect.Selected,
			SuspensionDays:     suspDays,
			SuspensionDates:    suspensionDatesEntry.Text,
			DeductedPoints:     ded,
			RemainingPoints:    rem,
			DecisionText:       verdictTextEntry.Text,
			ChairIsSubstitute:  subChairCheck.Checked,
			M1IsSubstitute:     subM1Check.Checked,
			M2IsSubstitute:     subM2Check.Checked,
			M3IsSubstitute:     subM3Check.Checked,
			M4IsSubstitute:     subM4Check.Checked,
		}

		// Master student info
		mStudent, _ := state.Store.GetStudentByNo(state.AcademicYear, selectedIncStudent.SchoolNo)
		if mStudent == nil {
			mStudent = &model.Student{
				SchoolNo:       selectedIncStudent.SchoolNo,
				FullName:       selectedIncStudent.FullName,
				ClassSection:   selectedIncStudent.ClassSection,
				Gender:         selectedIncStudent.Gender,
				TCNo:           selectedIncStudent.TCNo,
				BirthPlace:     selectedIncStudent.BirthPlace,
				BirthDate:      selectedIncStudent.BirthDate,
				FatherName:     selectedIncStudent.FatherName,
				MotherName:     selectedIncStudent.MotherName,
				GuardianName:   selectedIncStudent.GuardianName,
				GuardianPhone:  selectedIncStudent.GuardianPhone,
				GuardianAddr:   selectedIncStudent.GuardianAddr,
				BoardingStatus: selectedIncStudent.BoardingStatus,
			}
		}

		ekPrefix := "Karar_Tutanağı_EK-1"
		if strings.ToUpper(state.SchoolInfo.SchoolType) == "ORTAOKUL" {
			ekPrefix = "Karar_Tutanağı_EK-10"
		}
		prefix := fmt.Sprintf("%s_%s_%s", ekPrefix, selectedIncStudent.SchoolNo, strings.ReplaceAll(selectedIncStudent.FullName, " ", "_"))

		outPath, err := pdf.GenerateUniquePDFPath("", prefix)
		if err != nil {
			showError(state.Window, err)
			return
		}

		err = pdf.GenerateKurulKarari(state.SchoolInfo, mStudent, selectedIncident, selectedIncStudent, verdict, outPath)
		if err != nil {
			showError(state.Window, fmt.Errorf("PDF oluşturulamadı: %w", err))
			return
		}

		err = pdf.OpenPDF(outPath)
		if err != nil {
			showInfo(state.Window, "PDF Oluşturuldu", fmt.Sprintf("Dosya kaydedildi ancak otomatik açılamadı:\n%s", outPath))
		}
	})
	pdfBtn.Importance = widget.SuccessImportance

	// Assembly of UI Containers
	selectorsCard := widget.NewCard("Olay ve Öğrenci Seçimi", "Karar alınacak dosya ve öğrenciyi seçin", container.NewVBox(
		widget.NewLabelWithStyle("İncelenen Olay:", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		incSelect,
		widget.NewLabelWithStyle("Olaydaki Öğrenci:", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		stSelect,
	))

	decisionDetailsForm := widget.NewForm(
		widget.NewFormItem("Karar No *", verdictNoEntry),
		widget.NewFormItem("Karar Tarihi", verdictDateEntry),
		widget.NewFormItem("Suç Yeri", offensePlaceEntry),
		widget.NewFormItem("Suç Tarihi ve Saati", offenseDateEntry),
		widget.NewFormItem("Suç Çeşidi / Kategori", offenseCategorySelect),
		widget.NewFormItem("Suçun Nedeni", offenseReasonEntry),
		widget.NewFormItem("Tanıkların İfadeleri", witnessEntry),
		widget.NewFormItem("Deliller", evidenceEntry),
		widget.NewFormItem("Hafifletici Nedenler", mitigatingSelect),
		widget.NewFormItem("Ağırlaştırıcı Nedenler", aggravatingSelect),
		widget.NewFormItem("Kurulun Kanaati", opinionSelect),
	)

	clauseForm := widget.NewForm(
		widget.NewFormItem("Yönetmelik Maddesi Seçin", clauseSelect),
		widget.NewFormItem("Madde No", clauseCodeEntry),
		widget.NewFormItem("Madde Açıklaması", clauseDescEntry),
		widget.NewFormItem("Verilen Ceza / Yaptırım", penaltySelect),
		widget.NewFormItem("Uzaklaştırma Gün Sayısı", suspensionDaysEntry),
		widget.NewFormItem("Uzaklaştırma Tarihleri", suspensionDatesEntry),
		widget.NewFormItem("İndirilen Davranış Puanı", deductedPointsEntry),
		widget.NewFormItem("Kalan Davranış Puanı", remainingPointsEntry),
	)

	substituteBox := container.NewVBox(
		widget.NewLabelWithStyle("⚖️ Asil / Yedek Üye Görevlendirme", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		subChairCheck,
		subM1Check,
		subM2Check,
		subM3Check,
		subM4Check,
	)

	verdictTextCard := widget.NewCard("Gerekçeli Kurul Kararı Metni", "Otomatik oluşturabilir veya metni doğrudan düzenleyebilirsiniz", container.NewVBox(
		autoFillBtn,
		verdictTextEntry,
	))

	actionButtons := container.NewHBox(
		saveVerdictBtn,
		widget.NewSeparator(),
		pdfBtn,
	)

	scrollContent := container.NewVScroll(container.NewVBox(
		selectorsCard,
		widget.NewCard("Olay ve İnceleme Bulguları", "Disiplin Kurulu Karar Tutanağı (EK-1 / EK-10)", decisionDetailsForm),
		widget.NewCard("Yönetmelik Maddesi ve Ceza Belirleme", "Mevzuat hükümleri ve puan indirim cetveli", clauseForm),
		substituteBox,
		verdictTextCard,
		actionButtons,
	))

	witnessEntry.ParentScroll = scrollContent
	evidenceEntry.ParentScroll = scrollContent
	clauseDescEntry.ParentScroll = scrollContent
	verdictTextEntry.ParentScroll = scrollContent

	reloadIncidents()
	return scrollContent
}
