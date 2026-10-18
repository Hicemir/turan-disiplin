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

// BuildNoticesTab creates the Penalty Notification Slips UI (Öğrenci & Veli Karar Tebligatı)
func BuildNoticesTab(state *AppState) fyne.CanvasObject {
	var penalizedStudents []model.IncidentStudent
	var selectedStudent *model.IncidentStudent

	// Filters
	searchEntry := widget.NewEntry()
	searchEntry.SetPlaceHolder("Öğrenci no veya ad soyad ile ara...")

	classes := append([]string{"Tüm Sınıflar"}, state.Store.GetClasses(state.AcademicYear)...)
	classSelect := NewTurkishSelect(classes, nil)
	classSelect.SetSelected("Tüm Sınıflar")

	penalties := []string{"Tüm Cezalar", "Sözlü Uyarma", "Yazılı Uyarma", "Kınama", "Okuldan Kısa Süreli Uzaklaştırma (1-5 Gün)", "Okul Değiştirme", "Örgün Eğitim Dışına Çıkarma"}
	penaltySelect := NewTurkishSelect(penalties, nil)
	penaltySelect.SetSelected("Tüm Cezalar")

	// Notice Controls
	sayiEntry := NewSingleLineEntry()
	sayiEntry.SetText(fmt.Sprintf("%s%d", state.SchoolInfo.OfficialCode, 100))

	tebligTarihiEntry := NewSingleLineEntry()
	tebligTarihiEntry.SetText(time.Now().Format("02.01.2006"))

	// Verdict Summary Labels
	summaryLabel := widget.NewLabelWithStyle("Bir öğrenci seçiniz.", fyne.TextAlignLeading, fyne.TextStyle{Italic: true})

	// Student List
	noticesList := widget.NewList(
		func() int {
			return len(penalizedStudents)
		},
		func() fyne.CanvasObject {
			return container.NewHBox(
				widget.NewLabelWithStyle("No", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
				widget.NewSeparator(),
				widget.NewLabel("Öğrenci Adı"),
				widget.NewSeparator(),
				widget.NewLabel("Sınıf"),
				widget.NewSeparator(),
				widget.NewLabelWithStyle("Verilen Ceza", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
				widget.NewSeparator(),
				widget.NewLabel("Karar No"),
			)
		},
		func(i widget.ListItemID, o fyne.CanvasObject) {
			if i >= len(penalizedStudents) {
				return
			}
			st := penalizedStudents[i]
			box := o.(*fyne.Container)
			box.Objects[0].(*widget.Label).SetText(st.SchoolNo)
			box.Objects[2].(*widget.Label).SetText(st.FullName)
			box.Objects[4].(*widget.Label).SetText(st.ClassSection)
			box.Objects[6].(*widget.Label).SetText(st.ImposedPenalty)
			box.Objects[8].(*widget.Label).SetText(fmt.Sprintf("K.No: %s", st.VerdictNo))
		},
	)

	// Action 1: Öğrenci Karar Tebligat Pusulası PDF
	pdfOgrenciTebligatBtn := widget.NewButtonWithIcon("PDF OLUŞTUR VE AÇ (ÖĞRENCİ TEBLİGAT PUSULASI)", theme.DocumentIcon(), func() {
		if selectedStudent == nil {
			showError(state.Window, fmt.Errorf("Lütfen listeden bir öğrenci seçiniz"))
			return
		}

		vNo, _ := strconv.Atoi(selectedStudent.VerdictNo)
		suspDays, _ := strconv.Atoi(selectedStudent.SuspensionDays)

		verdict := &model.Verdict{
			DecisionNo:      vNo,
			DecisionDate:    selectedStudent.VerdictDate,
			OffenseType:     selectedStudent.OffenseCategory,
			RegClause:       selectedStudent.ClauseCode,
			RegClauseDesc:   selectedStudent.ClauseDescription,
			PenaltyName:     selectedStudent.ImposedPenalty,
			SuspensionDays:  suspDays,
			SuspensionDates: selectedStudent.SuspensionDates,
			DeductedPoints:  selectedStudent.DeductedPoints,
			RemainingPoints: selectedStudent.RemainingPoints,
		}

		mStudent, _ := state.Store.GetStudentByNo(state.AcademicYear, selectedStudent.SchoolNo)
		if mStudent == nil {
			mStudent = &model.Student{
				SchoolNo:     selectedStudent.SchoolNo,
				FullName:     selectedStudent.FullName,
				ClassSection: selectedStudent.ClassSection,
			}
		}

		inc, _ := state.Store.GetIncidentByCode(state.AcademicYear, selectedStudent.IncidentCode)

		outPath, err := pdf.GenerateUniquePDFPath("", fmt.Sprintf("Ogrenci_Tebligat_%s_%s", selectedStudent.SchoolNo, strings.ReplaceAll(selectedStudent.FullName, " ", "_")))
		if err != nil {
			showError(state.Window, err)
			return
		}

		err = pdf.GenerateCezaTebligatiOgrenci(state.SchoolInfo, mStudent, inc, selectedStudent, verdict, outPath)
		if err != nil {
			showError(state.Window, fmt.Errorf("PDF oluşturulamadı: %w", err))
			return
		}

		_ = pdf.OpenPDF(outPath)
	})
	pdfOgrenciTebligatBtn.Importance = widget.SuccessImportance

	// Action 2: Veli Karar Tebligat Pusulası PDF
	pdfVeliTebligatBtn := widget.NewButtonWithIcon("PDF OLUŞTUR VE AÇ (VELİ TEBLİGAT PUSULASI)", theme.DocumentIcon(), func() {
		if selectedStudent == nil {
			showError(state.Window, fmt.Errorf("Lütfen listeden bir öğrenci seçiniz"))
			return
		}

		vNo, _ := strconv.Atoi(selectedStudent.VerdictNo)
		suspDays, _ := strconv.Atoi(selectedStudent.SuspensionDays)

		verdict := &model.Verdict{
			DecisionNo:      vNo,
			DecisionDate:    selectedStudent.VerdictDate,
			OffenseType:     selectedStudent.OffenseCategory,
			RegClause:       selectedStudent.ClauseCode,
			RegClauseDesc:   selectedStudent.ClauseDescription,
			PenaltyName:     selectedStudent.ImposedPenalty,
			SuspensionDays:  suspDays,
			SuspensionDates: selectedStudent.SuspensionDates,
		}

		mStudent, _ := state.Store.GetStudentByNo(state.AcademicYear, selectedStudent.SchoolNo)
		if mStudent == nil {
			mStudent = &model.Student{
				SchoolNo:     selectedStudent.SchoolNo,
				FullName:     selectedStudent.FullName,
				ClassSection: selectedStudent.ClassSection,
				GuardianName: selectedStudent.GuardianName,
				FatherName:   selectedStudent.FatherName,
			}
		}

		inc, _ := state.Store.GetIncidentByCode(state.AcademicYear, selectedStudent.IncidentCode)

		outPath, err := pdf.GenerateUniquePDFPath("", fmt.Sprintf("Veli_Tebligat_%s_%s", selectedStudent.SchoolNo, strings.ReplaceAll(selectedStudent.FullName, " ", "_")))
		if err != nil {
			showError(state.Window, err)
			return
		}

		err = pdf.GenerateCezaTebligatiVeli(state.SchoolInfo, mStudent, inc, selectedStudent, verdict, sayiEntry.Text, outPath)
		if err != nil {
			showError(state.Window, fmt.Errorf("PDF oluşturulamadı: %w", err))
			return
		}

		_ = pdf.OpenPDF(outPath)
	})
	pdfOgrenciTebligatBtn.Disable()
	pdfVeliTebligatBtn.Disable()

	// Empty state message when no verdicts exist yet in the academic year
	emptyNoticeBox := container.NewCenter(container.NewVBox(
		widget.NewLabelWithStyle("ℹ️ Bu öğretim yılında henüz kesinleşmiş bir disiplin kurulu kararı bulunmamaktadır.", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		widget.NewLabelWithStyle("Ceza Tebligatı: Disiplin Kurulu tarafından kesinleşen cezaların (Kınama, Uzaklaştırma vb.)\nöğrenciye ve veliye resmi olarak tebliğ edilmesi için tebligat pusulası oluşturur.\n\nÖğrenciye ceza tebligatı düzenlemek için önce 'DİSİPLİN KURULU KARARI (EK-1)' bölümünden ilgili olay için kurul kararı oluşturunuz.", fyne.TextAlignCenter, fyne.TextStyle{Italic: true}),
	))

	listStack := container.NewStack(noticesList, emptyNoticeBox)

	reloadPenalized := func() {
		cFilter := classSelect.Selected
		if cFilter == "Tüm Sınıflar" {
			cFilter = ""
		}
		pFilter := penaltySelect.Selected
		if pFilter == "Tüm Cezalar" {
			pFilter = ""
		}
		list, err := state.Store.GetPenalizedStudents(state.AcademicYear, cFilter, pFilter)
		if err == nil {
			// Apply search filter
			var filtered []model.IncidentStudent
			sTerm := strings.ToLower(searchEntry.Text)
			for _, st := range list {
				if sTerm == "" || strings.Contains(strings.ToLower(st.FullName), sTerm) || strings.Contains(st.SchoolNo, sTerm) {
					filtered = append(filtered, st)
				}
			}
			penalizedStudents = filtered
			noticesList.Refresh()

			if len(penalizedStudents) == 0 {
				emptyNoticeBox.Show()
				noticesList.Hide()
			} else {
				emptyNoticeBox.Hide()
				noticesList.Show()
			}
		}
	}

	noticesList.OnSelected = func(id widget.ListItemID) {
		if id >= 0 && id < len(penalizedStudents) {
			st := penalizedStudents[id]
			selectedStudent = &st
			summaryLabel.SetText(fmt.Sprintf("Seçilen: %s (%s - %s)\nKarar No: %s  |  Tarih: %s\nVerilen Ceza: %s\nYönetmelik Maddesi: %s (%s)\nİndirilen Puan: %d  |  Kalan Puan: %d",
				st.FullName, st.ClassSection, st.SchoolNo, st.VerdictNo, st.VerdictDate, st.ImposedPenalty, st.ClauseCode, st.ClauseDescription, st.DeductedPoints, st.RemainingPoints))
			pdfOgrenciTebligatBtn.Enable()
			pdfVeliTebligatBtn.Enable()
		}
	}

	searchEntry.OnChanged = func(s string) { reloadPenalized() }
	classSelect.OnChanged = func(s string) { reloadPenalized() }
	penaltySelect.OnChanged = func(s string) { reloadPenalized() }

	// Filter Row (Responsive with full-width search input)
	filterBar := container.NewBorder(
		nil, nil,
		widget.NewLabelWithStyle("Filtrele:", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		container.NewHBox(classSelect, penaltySelect),
		searchEntry,
	)

	// Settings Form
	tebligForm := widget.NewForm(
		widget.NewFormItem("Resmi Yazı Sayısı (Veli Tebligatı)", sayiEntry),
		widget.NewFormItem("Tebliğ Tarihi", tebligTarihiEntry),
	)

	summaryCard := widget.NewCard("Seçilen Öğrenci Karar Özeti", "Tebligatı yapılacak disiplin kararı", container.NewVBox(
		summaryLabel,
		tebligForm,
		container.NewHBox(pdfOgrenciTebligatBtn, pdfVeliTebligatBtn),
	))

	content := container.NewBorder(
		container.NewVBox(
			widget.NewLabelWithStyle("📬 Disiplin / Yaptırım Karar Tebligatları (Öğrenci & Veli)", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
			filterBar,
		),
		summaryCard,
		nil, nil,
		listStack,
	)

	reloadPenalized()
	return content
}
