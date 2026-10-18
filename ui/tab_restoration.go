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

// BuildRestorationTab creates the Behavior Point Restoration & Sanction Lift UI
func BuildRestorationTab(state *AppState) fyne.CanvasObject {
	var penalizedList []model.IncidentStudent
	var selectedStudent *model.IncidentStudent

	isOrta := strings.ToUpper(state.SchoolInfo.SchoolType) == "ORTAOKUL"

	refundDateEntry := widget.NewEntry()
	refundDateEntry.SetText(time.Now().Format("02.01.2006"))

	refundPointsEntry := widget.NewEntry()
	refundPointsEntry.SetText("10")

	statusLabel := widget.NewLabelWithStyle("Bir öğrenci seçiniz.", fyne.TextAlignLeading, fyne.TextStyle{Italic: true})

	// Student Table List
	stList := widget.NewList(
		func() int {
			return len(penalizedList)
		},
		func() fyne.CanvasObject {
			return container.NewHBox(
				widget.NewLabelWithStyle("No", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
				widget.NewSeparator(),
				widget.NewLabel("Öğrenci Adı"),
				widget.NewSeparator(),
				widget.NewLabel("Sınıf"),
				widget.NewSeparator(),
				widget.NewLabel("Ceza / Yaptırım"),
				widget.NewSeparator(),
				widget.NewLabelWithStyle("Durum", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
			)
		},
		func(i widget.ListItemID, o fyne.CanvasObject) {
			if i >= len(penalizedList) {
				return
			}
			st := penalizedList[i]
			box := o.(*fyne.Container)
			box.Objects[0].(*widget.Label).SetText(st.SchoolNo)
			box.Objects[2].(*widget.Label).SetText(st.FullName)
			box.Objects[4].(*widget.Label).SetText(st.ClassSection)
			box.Objects[6].(*widget.Label).SetText(st.ImposedPenalty)

			stat := st.RestorationStatus
			if stat == "" {
				stat = "BEKLEMEDE"
			}
			box.Objects[8].(*widget.Label).SetText(stat)
		},
	)

	reloadList := func() {
		list, err := state.Store.GetPenalizedStudents(state.AcademicYear, "", "")
		if err == nil {
			penalizedList = list
			stList.Refresh()
		}
	}

	stList.OnSelected = func(id widget.ListItemID) {
		if id >= 0 && id < len(penalizedList) {
			st := penalizedList[id]
			selectedStudent = &st
			refundPointsEntry.SetText(fmt.Sprintf("%d", st.DeductedPoints))

			stat := st.RestorationStatus
			if stat == "" {
				stat = "BEKLEMEDE (İade Edilmedi / Kaldırılmadı)"
			}
			statusLabel.SetText(fmt.Sprintf("Öğrenci: %s (%s - %s)\nVerilen Ceza: %s | İndirilen Puan: %d | Kalan Puan: %d\nMevcut Durum: %s\nİade/Kaldırma Tarihi: %s",
				st.FullName, st.ClassSection, st.SchoolNo, st.ImposedPenalty, st.DeductedPoints, st.RemainingPoints, stat, st.RestorationDate))
		}
	}


	// Execute Point Refund / Sanction Lift
	applyBtn := widget.NewButtonWithIcon("İade / Kaldırma Kararını Uygula", theme.ConfirmIcon(), func() {
		if selectedStudent == nil {
			showError(state.Window, fmt.Errorf("Lütfen önce listeden bir öğrenci seçiniz"))
			return
		}

		if isOrta {
			selectedStudent.RestorationStatus = "YAPTIRIM KALDIRILDI"
			selectedStudent.RestorationDate = refundDateEntry.Text
		} else {
			pts, _ := strconv.Atoi(refundPointsEntry.Text)
			selectedStudent.RestorationStatus = "İADE EDİLDİ"
			selectedStudent.RestorationDate = refundDateEntry.Text

			// Restore points in master student record
			st, err := state.Store.GetStudentByNo(state.AcademicYear, selectedStudent.SchoolNo)
			if err == nil && st != nil {
				st.BehaviorScore = st.BehaviorScore + pts
				if st.BehaviorScore > 100 {
					st.BehaviorScore = 100
				}
				_ = state.Store.SaveStudent(st)
			}
		}

		_ = state.Store.SaveIncidentStudent(selectedStudent)
		showInfo(state.Window, "İşlem Tamamlandı", "Öğrencinin puan iade / yaptırım kaldırma durumu güncellendi.")
		reloadList()
		state.NotifyChange()
	})
	applyBtn.Importance = widget.HighImportance

	// Prominent PDF Generation Button
	pdfBtn := widget.NewButtonWithIcon("PDF OLUŞTUR VE AÇ (SENE SONU ÇİZELGESİ / RAPORU)", theme.DocumentIcon(), func() {
		if len(penalizedList) == 0 {
			showError(state.Window, fmt.Errorf("Bu öğretim yılına ait ceza/yaptırım kaydı bulunmamaktadır"))
			return
		}

		if isOrta {
			// Generate Middle School Sanction Lift Report (Landscape A4)
			var rows []model.SanctionLiftRow
			for _, st := range penalizedList {
				vNo, _ := strconv.Atoi(st.VerdictNo)
				stat := st.RestorationStatus
				if stat == "" {
					stat = "YÜRÜRLÜKTE"
				}
				rows = append(rows, model.SanctionLiftRow{
					FullName:     st.FullName,
					SchoolNo:     st.SchoolNo,
					ClassSection: st.ClassSection,
					IncidentName: st.IncidentTitle,
					DecisionNo:   vNo,
					DecisionDate: st.VerdictDate,
					PenaltyName:  st.ImposedPenalty,
					RegClause:    st.ClauseCode,
					LiftDate:     st.RestorationDate,
					Status:       stat,
				})
			}

			outPath, err := pdf.GenerateUniquePDFPath("", "Yaptirim_Karari_Kaldirma_Raporu")
			if err != nil {
				showError(state.Window, err)
				return
			}

			err = pdf.GenerateYaptirimKaldirmaRaporu(state.SchoolInfo, rows, outPath)
			if err != nil {
				showError(state.Window, fmt.Errorf("PDF oluşturulamadı: %w", err))
				return
			}
			_ = pdf.OpenPDF(outPath)

		} else {
			// Generate High School Behavior Score Refund Chart (Landscape A4)
			var rows []model.BehaviorScoreRefundRow
			for _, st := range penalizedList {
				vNo, _ := strconv.Atoi(st.VerdictNo)
				stat := st.RestorationStatus
				if stat == "" {
					stat = "İADE EDİLMEDİ"
				}
				refPts := 0
				if stat == "İADE EDİLDİ" {
					refPts = st.DeductedPoints
				}
				rows = append(rows, model.BehaviorScoreRefundRow{
					FullName:        st.FullName,
					SchoolNo:        st.SchoolNo,
					ClassSection:    st.ClassSection,
					IncidentName:    st.IncidentTitle,
					DecisionNo:      vNo,
					DecisionDate:    st.VerdictDate,
					PenaltyName:     st.ImposedPenalty,
					RegClause:       st.ClauseCode,
					DeductedPoints:  st.DeductedPoints,
					RemainingPoints: st.RemainingPoints,
					RefundPoints:    refPts,
					RefundDate:      st.RestorationDate,
					Status:          stat,
				})
			}

			outPath, err := pdf.GenerateUniquePDFPath("", "Davranis_Puani_Iade_Cizelgesi")
			if err != nil {
				showError(state.Window, err)
				return
			}

			err = pdf.GenerateDavranisPuaniIadeCizelgesi(state.SchoolInfo, rows, outPath)
			if err != nil {
				showError(state.Window, fmt.Errorf("PDF oluşturulamadı: %w", err))
				return
			}
			_ = pdf.OpenPDF(outPath)
		}
	})
	pdfBtn.Importance = widget.SuccessImportance

	// Layout
	headerTitle := "MEB Ortaöğretim Kurumları Yönetmeliği Madde 171/2 - Sene Sonu Davranış Puanı İade Çizelgesi"
	if isOrta {
		headerTitle = "MEB Okul Öncesi Eğitim ve İlköğretim Kurumları Yönetmeliği Madde 62/3 - Yaptırım Kararı Kaldırma Raporu"
	}

	formItems := widget.NewForm(
		widget.NewFormItem("İade / Kaldırma Karar Tarihi", refundDateEntry),
	)
	if !isOrta {
		formItems.Append("İade Edilecek Davranış Puanı", refundPointsEntry)
	}

	actionsCard := widget.NewCard("İade / Yaptırım Kaldırma İşlemi", "Öğrencinin olumlu davranış geliştirmesi nedeniyle kurul kararıyla puan iadesi", container.NewVBox(
		statusLabel,
		formItems,
		container.NewHBox(applyBtn, widget.NewSeparator(), pdfBtn),
	))

	content := container.NewBorder(
		widget.NewLabelWithStyle(headerTitle, fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		actionsCard,
		nil, nil,
		stList,
	)

	reloadList()
	return content
}
