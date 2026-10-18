package ui

import (
	"fmt"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"okulevrak/model"
	"okulevrak/pdf"
)

// BuildStatsTab creates the Disciplinary Statistics & Penalized Students Roster UI
func BuildStatsTab(state *AppState) fyne.CanvasObject {
	var currentList []model.IncidentStudent
	var semesterStats []model.RewardPenaltyCount

	classSelect := NewTurkishSelect(append([]string{"Tüm Sınıflar"}, state.Store.GetClasses(state.AcademicYear)...), nil)
	classSelect.SetSelected("Tüm Sınıflar")

	penaltySelect := NewTurkishSelect([]string{"Tüm Cezalar", "Kınama", "Okuldan Kısa Süreli Uzaklaştırma", "Okul Değiştirme", "Örgün Eğitim Dışına Çıkarma", "Uyarma"}, nil)
	penaltySelect.SetSelected("Tüm Cezalar")

	totalLabel := widget.NewLabelWithStyle("Kayıt: 0", fyne.TextAlignLeading, fyne.TextStyle{Bold: true})

	// Table 1: Penalized Students List
	penalizedList := widget.NewList(
		func() int {
			return len(currentList)
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
				widget.NewLabel("Karar Tarih/No"),
			)
		},
		func(i widget.ListItemID, o fyne.CanvasObject) {
			if i >= len(currentList) {
				return
			}
			st := currentList[i]
			box := o.(*fyne.Container)
			box.Objects[0].(*widget.Label).SetText(st.SchoolNo)
			box.Objects[2].(*widget.Label).SetText(st.FullName)
			box.Objects[4].(*widget.Label).SetText(st.ClassSection)
			box.Objects[6].(*widget.Label).SetText(st.ImposedPenalty)
			box.Objects[8].(*widget.Label).SetText(fmt.Sprintf("%s (No:%s)", st.VerdictDate, st.VerdictNo))
		},
	)

	// Table 2: Semester Breakdown List
	statsList := widget.NewList(
		func() int {
			return len(semesterStats)
		},
		func() fyne.CanvasObject {
			return container.NewHBox(
				widget.NewLabelWithStyle("Ceza / Ödül Türü", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
				widget.NewSeparator(),
				widget.NewLabel("1. Dönem: 0"),
				widget.NewSeparator(),
				widget.NewLabel("2. Dönem: 0"),
				widget.NewSeparator(),
				widget.NewLabelWithStyle("Toplam: 0", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
			)
		},
		func(i widget.ListItemID, o fyne.CanvasObject) {
			if i >= len(semesterStats) {
				return
			}
			s := semesterStats[i]
			box := o.(*fyne.Container)
			box.Objects[0].(*widget.Label).SetText(s.Name)
			box.Objects[2].(*widget.Label).SetText(fmt.Sprintf("1. Dönem: %d", s.Sem1Count))
			box.Objects[4].(*widget.Label).SetText(fmt.Sprintf("2. Dönem: %d", s.Sem2Count))
			box.Objects[6].(*widget.Label).SetText(fmt.Sprintf("Toplam: %d", s.TotalCount))
		},
	)

	reload := func() {
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
			currentList = list
			totalLabel.SetText(fmt.Sprintf("Toplam Ceza Alan Öğrenci: %d", len(currentList)))
			penalizedList.Refresh()
		}

		sList, err := state.Store.CalculateSemesterPenaltyCounts(state.AcademicYear)
		if err == nil {
			semesterStats = sList
			statsList.Refresh()
		}
	}

	classSelect.OnChanged = func(s string) { reload() }
	penaltySelect.OnChanged = func(s string) { reload() }


	// PDF 1: Ceza Alan Öğrenciler Raporu
	pdfRosterBtn := widget.NewButtonWithIcon("PDF CEZA ALANLAR LİSTESİ", theme.DocumentIcon(), func() {
		if len(currentList) == 0 {
			showError(state.Window, fmt.Errorf("Yazdırılacak ceza kaydı bulunamadı"))
			return
		}

		var stats []model.StatItem
		for _, st := range currentList {
			stats = append(stats, model.StatItem{
				Name:  fmt.Sprintf("%s (%s - %s) - %s (K.No: %s)", st.FullName, st.ClassSection, st.SchoolNo, st.ImposedPenalty, st.VerdictNo),
				Count: 1,
			})
		}

		outPath, err := pdf.GenerateUniquePDFPath("", "Ceza_Alan_Ogrenciler_Raporu")
		if err != nil {
			showError(state.Window, err)
			return
		}

		title := fmt.Sprintf("%s\n%s ÖĞRETİM YILI CEZA ALAN ÖĞRENCİLER LİSTESİ", strings.ToUpper(state.SchoolInfo.SchoolName), state.AcademicYear)
		err = pdf.GenerateCezaIstatistikRaporu(state.SchoolInfo, title, stats, outPath)
		if err != nil {
			showError(state.Window, fmt.Errorf("PDF oluşturulamadı: %w", err))
			return
		}

		_ = pdf.OpenPDF(outPath)
	})
	pdfRosterBtn.Importance = widget.SuccessImportance

	// PDF 2: Ceza İstatistik Raporu
	pdfStatsBtn := widget.NewButtonWithIcon("PDF DÖNEMLİK İSTATİSTİK RAPORU", theme.DocumentIcon(), func() {
		if len(semesterStats) == 0 {
			showError(state.Window, fmt.Errorf("İstatistik verisi bulunamadı"))
			return
		}

		var stats []model.StatItem
		for _, s := range semesterStats {
			stats = append(stats, model.StatItem{
				Name:  fmt.Sprintf("%s (1.Dönem: %d, 2.Dönem: %d)", s.Name, s.Sem1Count, s.Sem2Count),
				Count: s.TotalCount,
			})
		}

		outPath, err := pdf.GenerateUniquePDFPath("", "Ceza_Istatistik_Raporu")
		if err != nil {
			showError(state.Window, err)
			return
		}

		title := fmt.Sprintf("%s\n%s EĞİTİM ÖĞRETİM YILI DÖNEMLERE GÖRE DİSİPLİN CEZASI İSTATİSTİKLERİ", strings.ToUpper(state.SchoolInfo.SchoolName), state.AcademicYear)
		err = pdf.GenerateCezaIstatistikRaporu(state.SchoolInfo, title, stats, outPath)
		if err != nil {
			showError(state.Window, fmt.Errorf("PDF oluşturulamadı: %w", err))
			return
		}

		_ = pdf.OpenPDF(outPath)
	})
	pdfStatsBtn.Importance = widget.SuccessImportance

	// Filter Row
	filterRow := container.NewHBox(
		widget.NewLabelWithStyle("Filtre:", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		classSelect,
		penaltySelect,
		totalLabel,
	)

	leftCard := container.NewBorder(
		container.NewVBox(
			widget.NewLabelWithStyle("📋 Ceza Alan Öğrenciler Listesi", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
			filterRow,
		),
		pdfRosterBtn,
		nil, nil,
		penalizedList,
	)

	rightCard := container.NewBorder(
		widget.NewLabelWithStyle("📊 Dönemlik Ceza Sayıları (1. ve 2. Dönem)", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		pdfStatsBtn,
		nil, nil,
		statsList,
	)

	split := container.NewHSplit(leftCard, rightCard)
	split.SetOffset(0.55)

	reload()
	return split
}
