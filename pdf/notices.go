package pdf

import (
	"fmt"
	"strings"

	"github.com/signintech/gopdf"

	"turandisiplin/model"
)

// GenerateCezaTebligatiOgrenci generates student penalty notification slip.
func GenerateCezaTebligatiOgrenci(school *model.SchoolInfo, student *model.Student, incident *model.Incident, incStudent *model.IncidentStudent, verdict *model.Verdict, outputPath string) error {
	doc, err := NewDoc(false)
	if err != nil {
		return err
	}

	doc.PDF.AddPage()

	boardName := "ÖDÜL VE DİSİPLİN"
	noticeLaw := "Milli Eğitim Bakanlığı Ortaöğretim Kurumları Yönetmeliğinin 169. maddesi gereğince"
	if strings.ToUpper(school.SchoolType) == "ORTAOKUL" {
		boardName = "ÖĞRENCİ DAVRANIŞLARINI DEĞERLENDİRME"
		noticeLaw = "Milli Eğitim Bakanlığı Okul Öncesi Eğitim ve İlköğretim Kurumları Yönetmeliğinin ilgili maddesi gereğince"
	}

	doc.DrawAntet(school)
	doc.SetBoldFont(9)
	doc.PDF.SetX(doc.MarginLeft)
	_ = doc.PDF.CellWithOption(&gopdf.Rect{W: doc.UsableWidth, H: 14}, boardName+" KURULU", gopdf.CellOption{Align: gopdf.Center})
	doc.PDF.Br(18)

	doc.SetBoldFont(13)
	doc.PDF.SetX(doc.MarginLeft)
	_ = doc.PDF.CellWithOption(&gopdf.Rect{W: doc.UsableWidth, H: 20}, "Ö Ğ R E N C İ   K A R A R   T E B L İ G A T   P U S U L A S I", gopdf.CellOption{Align: gopdf.Center})
	doc.PDF.Br(24)

	// Indented table
	indent := 30.0
	innerW := doc.UsableWidth - (indent * 2)
	lblW := 150.0
	valW := innerW - lblW

	boxTopY := doc.PDF.GetY()
	doc.PDF.SetLineWidth(1.2)
	doc.PDF.Line(doc.MarginLeft+indent, boxTopY, doc.MarginLeft+indent+innerW, boxTopY)
	doc.PDF.Br(6)

	kararTarihNo := "—"
	verilenCeza := "—"
	cezaUygTarihi := "—"
	cezaNedeni := "—"
	yonMadde := "—"

	if verdict != nil {
		kararTarihNo = fmt.Sprintf("%s  -  %d", verdict.DecisionDate, verdict.DecisionNo)
		verilenCeza = verdict.PenaltyName
		if verdict.SuspensionDates != "" {
			cezaUygTarihi = verdict.SuspensionDates
		}
		cezaNedeni = strings.TrimSpace(verdict.OffenseType)
		yonMadde = strings.TrimSpace(fmt.Sprintf("%s (%s)", verdict.RegClause, verdict.RegClauseDesc))
	}
	if cezaNedeni == "" && incStudent != nil && strings.TrimSpace(incStudent.OffenseCategory) != "" {
		cezaNedeni = strings.TrimSpace(incStudent.OffenseCategory)
	}
	if cezaNedeni == "" && incident != nil && strings.TrimSpace(incident.Title) != "" {
		cezaNedeni = strings.TrimSpace(incident.Title)
	}
	if cezaNedeni == "" && incStudent != nil && strings.TrimSpace(incStudent.IncidentTitle) != "" {
		cezaNedeni = strings.TrimSpace(incStudent.IncidentTitle)
	}
	if cezaNedeni == "" && verdict != nil && strings.TrimSpace(verdict.RegClauseDesc) != "" {
		cezaNedeni = strings.TrimSpace(verdict.RegClauseDesc)
	}
	if cezaNedeni == "" {
		cezaNedeni = "—"
	}

	stName := ""
	stClassNo := ""
	if student != nil {
		stName = student.FullName
		stClassNo = fmt.Sprintf("%s  -  %s", student.ClassSection, student.SchoolNo)
	}

	rows := [][2]string{
		{"ÖĞRETİM YILI :", school.AcademicYear},
		{"KARAR TARİHİ VE NO :", kararTarihNo},
		{"ÖĞRENCİNİN ADI SOYADI :", stName},
		{"SINIFI VE OKUL NO :", stClassNo},
		{"VERİLEN CEZA / YAPTIRIM :", verilenCeza},
		{"CEZA UYGULAMA TARİHİ :", cezaUygTarihi},
		{"CEZA NEDENİ DAVRANIŞ :", cezaNedeni},
		{"YÖNETMELİK MADDESİ :", yonMadde},
	}

	for _, r := range rows {
		rowY := doc.PDF.GetY()
		doc.SetBoldFont(8.0)
		doc.PDF.SetX(doc.MarginLeft + indent)
		doc.PDF.SetY(rowY)
		_ = doc.PDF.CellWithOption(&gopdf.Rect{W: lblW, H: 13}, r[0], gopdf.CellOption{Align: gopdf.Left})

		doc.SetRegularFont(8.0)
		valLines := doc.SplitLines(r[1], valW)
		doc.DrawWrappedCellLines(doc.MarginLeft+indent+lblW, rowY, valW, 12.0, valLines, gopdf.Left)
		rowH := float64(len(valLines)) * 12.0
		if rowH < 14.0 {
			rowH = 14.0
		}
		doc.PDF.SetY(rowY + rowH)
	}

	doc.PDF.Br(4)
	boxBotY := doc.PDF.GetY()
	doc.PDF.Line(doc.MarginLeft+indent, boxBotY, doc.MarginLeft+indent+innerW, boxBotY)
	doc.PDF.Br(16)

	// Notice text
	doc.SetRegularFont(8.5)
	nText := fmt.Sprintf("%s, hakkınızda kurul tarafından verilen yukarıdaki ceza/yaptırım kararı tarafınıza tebliğ edilmiştir. Karara karşı süresi içinde itiraz hakkınız bulunmaktadır.", noticeLaw)
	doc.DrawWrappedText(doc.MarginLeft+15, doc.UsableWidth-30, 13.0, nText, gopdf.Left)
	doc.PDF.Br(18)

	// Principal signature
	doc.SetBoldFont(8.0)
	doc.PDF.SetX(doc.MarginLeft)
	_ = doc.PDF.CellWithOption(&gopdf.Rect{W: doc.UsableWidth - 20, H: 12}, school.PrincipalName, gopdf.CellOption{Align: gopdf.Right})
	doc.PDF.Br(12)
	doc.SetRegularFont(7.5)
	doc.PDF.SetX(doc.MarginLeft)
	_ = doc.PDF.CellWithOption(&gopdf.Rect{W: doc.UsableWidth - 20, H: 11}, "Okul Müdürü", gopdf.CellOption{Align: gopdf.Right})
	doc.PDF.Br(30)

	// Dual Tebliğ / Tebellüğ Block
	colW := doc.UsableWidth / 2.0
	ySig := doc.PDF.GetY()

	doc.SetBoldFont(8.0)
	doc.PDF.SetX(doc.MarginLeft)
	doc.PDF.SetY(ySig)
	_ = doc.PDF.CellWithOption(&gopdf.Rect{W: colW, H: 12}, "TEBLİĞ EDEN", gopdf.CellOption{Align: gopdf.Center})

	doc.PDF.SetX(doc.MarginLeft + colW)
	doc.PDF.SetY(ySig)
	_ = doc.PDF.CellWithOption(&gopdf.Rect{W: colW, H: 12}, "TEBELLÜĞ EDEN (ÖĞRENCİ)", gopdf.CellOption{Align: gopdf.Center})

	doc.SetRegularFont(7.5)
	doc.PDF.SetX(doc.MarginLeft)
	doc.PDF.SetY(ySig + 25)
	_ = doc.PDF.CellWithOption(&gopdf.Rect{W: colW, H: 11}, "Adı Soyadı: ....................................", gopdf.CellOption{Align: gopdf.Center})

	doc.PDF.SetX(doc.MarginLeft + colW)
	doc.PDF.SetY(ySig + 25)
	_ = doc.PDF.CellWithOption(&gopdf.Rect{W: colW, H: 11}, stName, gopdf.CellOption{Align: gopdf.Center})

	doc.PDF.SetX(doc.MarginLeft)
	doc.PDF.SetY(ySig + 38)
	_ = doc.PDF.CellWithOption(&gopdf.Rect{W: colW, H: 11}, "Tarih: ......./......./20.......", gopdf.CellOption{Align: gopdf.Center})

	doc.PDF.SetX(doc.MarginLeft + colW)
	doc.PDF.SetY(ySig + 38)
	_ = doc.PDF.CellWithOption(&gopdf.Rect{W: colW, H: 11}, "Tarih: ......./......./20.......", gopdf.CellOption{Align: gopdf.Center})

	doc.PDF.SetX(doc.MarginLeft)
	doc.PDF.SetY(ySig + 51)
	_ = doc.PDF.CellWithOption(&gopdf.Rect{W: colW, H: 11}, "İmza: ....................................", gopdf.CellOption{Align: gopdf.Center})

	doc.PDF.SetX(doc.MarginLeft + colW)
	doc.PDF.SetY(ySig + 51)
	_ = doc.PDF.CellWithOption(&gopdf.Rect{W: colW, H: 11}, "İmza: ....................................", gopdf.CellOption{Align: gopdf.Center})

	doc.DrawFooterAddress(school)
	return doc.PDF.WritePdf(outputPath)
}

// GenerateCezaTebligatiVeli generates parent penalty notification.
func GenerateCezaTebligatiVeli(school *model.SchoolInfo, student *model.Student, incident *model.Incident, incStudent *model.IncidentStudent, verdict *model.Verdict, sayiNo string, outputPath string) error {
	doc, err := NewDoc(false)
	if err != nil {
		return err
	}

	doc.PDF.AddPage()
	doc.DrawAntet(school)

	// Sayı & Tarih
	doc.SetBoldFont(8.0)
	doc.PDF.SetX(doc.MarginLeft)
	_ = doc.PDF.CellWithOption(&gopdf.Rect{W: 50, H: 12}, "SAYI :", gopdf.CellOption{Align: gopdf.Left})
	doc.SetRegularFont(8.0)
	doc.PDF.SetX(doc.MarginLeft + 50)
	_ = doc.PDF.CellWithOption(&gopdf.Rect{W: 200, H: 12}, sayiNo, gopdf.CellOption{Align: gopdf.Left})

	tarih := ""
	if verdict != nil {
		tarih = verdict.DecisionDate
	}
	doc.PDF.SetX(doc.MarginLeft + doc.UsableWidth - 100)
	_ = doc.PDF.CellWithOption(&gopdf.Rect{W: 100, H: 12}, tarih, gopdf.CellOption{Align: gopdf.Right})
	doc.PDF.Br(14)

	doc.SetBoldFont(8.0)
	doc.PDF.SetX(doc.MarginLeft)
	_ = doc.PDF.CellWithOption(&gopdf.Rect{W: 50, H: 12}, "KONU :", gopdf.CellOption{Align: gopdf.Left})
	doc.SetRegularFont(8.0)
	doc.PDF.SetX(doc.MarginLeft + 50)
	_ = doc.PDF.CellWithOption(&gopdf.Rect{W: 200, H: 12}, "Ceza Bildirimi", gopdf.CellOption{Align: gopdf.Left})
	doc.PDF.Br(24)

	// Salutation
	parentName := ""
	stName := ""
	stClass := ""
	if student != nil {
		parentName = model.EffectiveGuardianName(student.FullName, student.FatherName, student.GuardianName)
		stName = student.FullName
		stClass = fmt.Sprintf("%s sınıfı %s numaralı", student.ClassSection, student.SchoolNo)
	}
	if parentName == "" && incStudent != nil {
		parentName = model.EffectiveGuardianName(incStudent.FullName, incStudent.FatherName, incStudent.GuardianName)
	}

	doc.SetBoldFont(9.0)
	doc.PDF.SetX(doc.MarginLeft)
	_ = doc.PDF.CellWithOption(&gopdf.Rect{W: doc.UsableWidth, H: 14}, fmt.Sprintf("SAYIN : %s", parentName), gopdf.CellOption{Align: gopdf.Left})
	doc.PDF.Br(18)

	appealText := "Milli Eğitim Bakanlığı Ortaöğretim Kurumları Yönetmeliğinin 170. maddesi gereğince bu karara karşı tebliğ tarihinden itibaren 5 (beş) iş günü içinde okul müdürlüğüne itiraz hakkınız bulunmaktadır."
	if strings.ToUpper(school.SchoolType) == "ORTAOKUL" {
		appealText = "Milli Eğitim Bakanlığı Okul Öncesi Eğitim ve İlköğretim Kurumları Yönetmeliğinin 56/10 maddesi gereğince bu karara karşı tebliğ tarihinden itibaren 5 (beş) iş günü içinde itiraz hakkınız bulunmaktadır."
	}

	doc.SetRegularFont(8.5)
	letterBody := fmt.Sprintf("Okulumuz %s öğrencisi %s hakkında aşağıda belirtilen disiplin kararı alınmıştır. %s Bilgilerinize rica ederim.", stClass, stName, appealText)
	doc.DrawWrappedText(doc.MarginLeft, doc.UsableWidth, 13.0, letterBody, gopdf.Left)
	doc.PDF.Br(18)

	// Principal
	doc.SetBoldFont(8.0)
	doc.PDF.SetX(doc.MarginLeft)
	_ = doc.PDF.CellWithOption(&gopdf.Rect{W: doc.UsableWidth - 20, H: 12}, school.PrincipalName, gopdf.CellOption{Align: gopdf.Right})
	doc.PDF.Br(12)
	doc.SetRegularFont(7.5)
	doc.PDF.SetX(doc.MarginLeft)
	_ = doc.PDF.CellWithOption(&gopdf.Rect{W: doc.UsableWidth - 20, H: 11}, "Okul Müdürü", gopdf.CellOption{Align: gopdf.Right})
	doc.PDF.Br(22)

	// 4-Column Table: [Karar Tarih ve No, Verilen Ceza, Yönetmelik Maddesi, Ceza Nedeni]
	w1 := doc.UsableWidth * 0.18
	w2 := doc.UsableWidth * 0.22
	w3 := doc.UsableWidth * 0.32
	w4 := doc.UsableWidth * 0.28

	yT := doc.PDF.GetY()
	doc.DrawRect(doc.MarginLeft, yT, w1, 20)
	doc.DrawRect(doc.MarginLeft+w1, yT, w2, 20)
	doc.DrawRect(doc.MarginLeft+w1+w2, yT, w3, 20)
	doc.DrawRect(doc.MarginLeft+w1+w2+w3, yT, w4, 20)

	doc.SetBoldFont(7.5)
	doc.PDF.SetX(doc.MarginLeft + 2)
	doc.PDF.SetY(yT + 4)
	_ = doc.PDF.CellWithOption(&gopdf.Rect{W: w1 - 4, H: 12}, "KARAR TARİH VE NO", gopdf.CellOption{Align: gopdf.Center})

	doc.PDF.SetX(doc.MarginLeft + w1 + 2)
	doc.PDF.SetY(yT + 4)
	_ = doc.PDF.CellWithOption(&gopdf.Rect{W: w2 - 4, H: 12}, "VERİLEN CEZA", gopdf.CellOption{Align: gopdf.Center})

	doc.PDF.SetX(doc.MarginLeft + w1 + w2 + 2)
	doc.PDF.SetY(yT + 4)
	_ = doc.PDF.CellWithOption(&gopdf.Rect{W: w3 - 4, H: 12}, "YÖNETMELİK MADDESİ", gopdf.CellOption{Align: gopdf.Center})

	doc.PDF.SetX(doc.MarginLeft + w1 + w2 + w3 + 2)
	doc.PDF.SetY(yT + 4)
	_ = doc.PDF.CellWithOption(&gopdf.Rect{W: w4 - 4, H: 12}, "CEZA NEDENİ DAVRANIŞ", gopdf.CellOption{Align: gopdf.Center})

	// Data Row
	kararTarihNo := "—"
	verilenCeza := "—"
	yonMadde := "—"
	cezaNedeni := ""
	if verdict != nil {
		kararTarihNo = fmt.Sprintf("%s / %d", verdict.DecisionDate, verdict.DecisionNo)
		verilenCeza = verdict.PenaltyName
		yonMadde = strings.TrimSpace(fmt.Sprintf("%s (%s)", verdict.RegClause, verdict.RegClauseDesc))
		cezaNedeni = strings.TrimSpace(verdict.OffenseType)
	}
	if cezaNedeni == "" && incStudent != nil && strings.TrimSpace(incStudent.OffenseCategory) != "" {
		cezaNedeni = strings.TrimSpace(incStudent.OffenseCategory)
	}
	if cezaNedeni == "" && incident != nil && strings.TrimSpace(incident.Title) != "" {
		cezaNedeni = strings.TrimSpace(incident.Title)
	}
	if cezaNedeni == "" && incStudent != nil && strings.TrimSpace(incStudent.IncidentTitle) != "" {
		cezaNedeni = strings.TrimSpace(incStudent.IncidentTitle)
	}
	if cezaNedeni == "" && verdict != nil && strings.TrimSpace(verdict.RegClauseDesc) != "" {
		cezaNedeni = strings.TrimSpace(verdict.RegClauseDesc)
	}
	if cezaNedeni == "" {
		cezaNedeni = "—"
	}

	lineH := 10.5
	doc.SetRegularFont(7.0)
	lines1 := doc.SplitLines(kararTarihNo, w1-6)
	doc.SetBoldFont(7.0)
	lines2 := doc.SplitLines(verilenCeza, w2-6)
	doc.SetRegularFont(6.8)
	lines3 := doc.SplitLines(yonMadde, w3-6)
	lines4 := doc.SplitLines(cezaNedeni, w4-6)

	maxLines := len(lines1)
	if len(lines2) > maxLines {
		maxLines = len(lines2)
	}
	if len(lines3) > maxLines {
		maxLines = len(lines3)
	}
	if len(lines4) > maxLines {
		maxLines = len(lines4)
	}

	curY := yT + 20
	hData := float64(maxLines)*lineH + 8.0
	if hData < 28.0 {
		hData = 28.0
	}
	doc.DrawRect(doc.MarginLeft, curY, w1, hData)
	doc.DrawRect(doc.MarginLeft+w1, curY, w2, hData)
	doc.DrawRect(doc.MarginLeft+w1+w2, curY, w3, hData)
	doc.DrawRect(doc.MarginLeft+w1+w2+w3, curY, w4, hData)

	doc.SetRegularFont(7.0)
	doc.DrawWrappedCellLines(doc.MarginLeft+3, curY+4, w1-6, lineH, lines1, gopdf.Center)

	doc.SetBoldFont(7.0)
	doc.DrawWrappedCellLines(doc.MarginLeft+w1+3, curY+4, w2-6, lineH, lines2, gopdf.Center)

	doc.SetRegularFont(6.8)
	doc.DrawWrappedCellLines(doc.MarginLeft+w1+w2+3, curY+4, w3-6, lineH, lines3, gopdf.Left)
	doc.DrawWrappedCellLines(doc.MarginLeft+w1+w2+w3+3, curY+4, w4-6, lineH, lines4, gopdf.Left)

	// Dual signature block
	colW := doc.UsableWidth / 2.0
	yAck := curY + hData + 30

	doc.SetBoldFont(8.0)
	doc.PDF.SetX(doc.MarginLeft)
	doc.PDF.SetY(yAck)
	_ = doc.PDF.CellWithOption(&gopdf.Rect{W: colW, H: 12}, "TEBLİĞ EDEN", gopdf.CellOption{Align: gopdf.Center})

	doc.PDF.SetX(doc.MarginLeft + colW)
	doc.PDF.SetY(yAck)
	_ = doc.PDF.CellWithOption(&gopdf.Rect{W: colW, H: 12}, "TEBELLÜĞ EDEN (VELİ)", gopdf.CellOption{Align: gopdf.Center})

	doc.SetRegularFont(7.5)
	doc.PDF.SetX(doc.MarginLeft)
	doc.PDF.SetY(yAck + 25)
	_ = doc.PDF.CellWithOption(&gopdf.Rect{W: colW, H: 11}, "Tarih: ......./......./20.......", gopdf.CellOption{Align: gopdf.Center})

	doc.PDF.SetX(doc.MarginLeft + colW)
	doc.PDF.SetY(yAck + 25)
	_ = doc.PDF.CellWithOption(&gopdf.Rect{W: colW, H: 11}, "Tarih: ......./......./20.......", gopdf.CellOption{Align: gopdf.Center})

	doc.PDF.SetX(doc.MarginLeft)
	doc.PDF.SetY(yAck + 38)
	_ = doc.PDF.CellWithOption(&gopdf.Rect{W: colW, H: 11}, "İmza: ....................................", gopdf.CellOption{Align: gopdf.Center})

	doc.PDF.SetX(doc.MarginLeft + colW)
	doc.PDF.SetY(yAck + 38)
	_ = doc.PDF.CellWithOption(&gopdf.Rect{W: colW, H: 11}, "İmza: ....................................", gopdf.CellOption{Align: gopdf.Center})

	doc.DrawFooterAddress(school)
	return doc.PDF.WritePdf(outputPath)
}

// GenerateDavranisPuaniIadeCizelgesi generates the Landscape A4 behavior points refund chart.
func GenerateDavranisPuaniIadeCizelgesi(school *model.SchoolInfo, rows []model.BehaviorScoreRefundRow, outputPath string) error {
	doc, err := NewDoc(true) // Landscape
	if err != nil {
		return err
	}

	doc.PDF.AddPage()

	header := fmt.Sprintf("%s\n%s EĞİTİM ÖĞRETİM YILI\nMEB ORTAÖĞRETİM KURUMLARI YÖNETMELİĞİNİN 171/2 MADDESİ GEREĞİ SENE SONU DAVRANIŞ PUANI İADE ÇİZELGESİ", strings.ToUpper(school.SchoolName), school.AcademicYear)
	doc.DrawTitle(header, 10, 10)

	// 14 columns
	colWidths := []float64{25, 75, 40, 45, 70, 40, 45, 75, 55, 45, 45, 45, 50, 60}
	headers := []string{"SIRA", "ADI SOYADI", "OKUL NO", "SINIFI", "OLAY ADI", "KARAR NO", "KARAR TARİHİ", "VERİLEN CEZA", "YÖN. MADDE", "İND. PUAN", "KALAN", "İADE", "İADE TARİHİ", "DURUM"}

	curY := doc.PDF.GetY()
	doc.SetBoldFont(6.5)

	x := doc.MarginLeft
	for i, w := range colWidths {
		doc.DrawRect(x, curY, w, 18)
		doc.PDF.SetX(x)
		doc.PDF.SetY(curY + 3)
		_ = doc.PDF.CellWithOption(&gopdf.Rect{W: w, H: 12}, headers[i], gopdf.CellOption{Align: gopdf.Center})
		x += w
	}

	curY += 18
	doc.SetRegularFont(6.0)
	lineH := 8.5

	for i, r := range rows {
		vals := []string{
			fmt.Sprintf("%d", i+1),
			r.FullName,
			r.SchoolNo,
			r.ClassSection,
			r.IncidentName,
			fmt.Sprintf("%d", r.DecisionNo),
			r.DecisionDate,
			r.PenaltyName,
			r.RegClause,
			fmt.Sprintf("%d", r.DeductedPoints),
			fmt.Sprintf("%d", r.RemainingPoints),
			fmt.Sprintf("%d", r.RefundPoints),
			r.RefundDate,
			r.Status,
		}

		wrappedCols := make([][]string, len(colWidths))
		maxLines := 1
		for j, w := range colWidths {
			wrappedCols[j] = doc.SplitLines(vals[j], w-4)
			if len(wrappedCols[j]) > maxLines {
				maxLines = len(wrappedCols[j])
			}
		}

		rowH := float64(maxLines)*lineH + 5.0
		if rowH < 16.0 {
			rowH = 16.0
		}

		if curY+rowH > PageHeightLandscape-doc.MarginBottom-30 {
			doc.PDF.AddPage()
			curY = doc.MarginTop
			doc.SetRegularFont(6.0)
		}

		x = doc.MarginLeft
		for j, w := range colWidths {
			doc.DrawRect(x, curY, w, rowH)
			align := gopdf.Center
			if j == 1 || j == 4 || j == 7 {
				align = gopdf.Left
			}
			doc.DrawWrappedCellLines(x+2, curY+2.5, w-4, lineH, wrappedCols[j], align)
			x += w
		}
		curY += rowH
	}

	doc.PDF.SetY(curY + 10)
	doc.DrawSignatures5(school, "Ödül ve Disiplin")
	doc.DrawPrincipalApproval(school.PrincipalName)

	return doc.PDF.WritePdf(outputPath)
}

// GenerateYaptirimKaldirmaRaporu generates the Landscape A4 middle-school sanction lift report.
func GenerateYaptirimKaldirmaRaporu(school *model.SchoolInfo, rows []model.SanctionLiftRow, outputPath string) error {
	doc, err := NewDoc(true) // Landscape
	if err != nil {
		return err
	}

	doc.PDF.AddPage()

	header := fmt.Sprintf("%s\n%s EĞİTİM ÖĞRETİM YILI\nMEB OKUL ÖNCESİ EĞİTİM VE İLKÖĞRETİM KURUMLARI YÖNETMELİĞİNİN 62/3 MADDESİ GEREĞİ YAPTIRIM KARARI KALDIRMA RAPORU", strings.ToUpper(school.SchoolName), school.AcademicYear)
	doc.DrawTitle(header, 10, 10)

	// 11 columns
	colWidths := []float64{25, 80, 45, 50, 85, 45, 55, 95, 60, 55, 70}
	headers := []string{"SIRA", "ADI SOYADI", "OKUL NO", "SINIFI", "OLAY ADI", "KARAR NO", "KARAR TARİHİ", "UYGULANAN YAPTIRIM", "YÖN. MADDE", "KALDIRMA TARİHİ", "DURUMU"}

	curY := doc.PDF.GetY()
	doc.SetBoldFont(6.5)

	x := doc.MarginLeft
	for i, w := range colWidths {
		doc.DrawRect(x, curY, w, 18)
		doc.PDF.SetX(x)
		doc.PDF.SetY(curY + 3)
		_ = doc.PDF.CellWithOption(&gopdf.Rect{W: w, H: 12}, headers[i], gopdf.CellOption{Align: gopdf.Center})
		x += w
	}

	curY += 18
	doc.SetRegularFont(6.0)
	lineH := 8.5

	for i, r := range rows {
		vals := []string{
			fmt.Sprintf("%d", i+1),
			r.FullName,
			r.SchoolNo,
			r.ClassSection,
			r.IncidentName,
			fmt.Sprintf("%d", r.DecisionNo),
			r.DecisionDate,
			r.PenaltyName,
			r.RegClause,
			r.LiftDate,
			r.Status,
		}

		wrappedCols := make([][]string, len(colWidths))
		maxLines := 1
		for j, w := range colWidths {
			wrappedCols[j] = doc.SplitLines(vals[j], w-4)
			if len(wrappedCols[j]) > maxLines {
				maxLines = len(wrappedCols[j])
			}
		}

		rowH := float64(maxLines)*lineH + 5.0
		if rowH < 16.0 {
			rowH = 16.0
		}

		if curY+rowH > PageHeightLandscape-doc.MarginBottom-30 {
			doc.PDF.AddPage()
			curY = doc.MarginTop
			doc.SetRegularFont(6.0)
		}

		x = doc.MarginLeft
		for j, w := range colWidths {
			doc.DrawRect(x, curY, w, rowH)
			align := gopdf.Center
			if j == 1 || j == 4 || j == 7 {
				align = gopdf.Left
			}
			doc.DrawWrappedCellLines(x+2, curY+2.5, w-4, lineH, wrappedCols[j], align)
			x += w
		}
		curY += rowH
	}

	doc.PDF.SetY(curY + 10)
	doc.DrawSignatures5(school, "Öğrenci Davranışlarını Değerlendirme")
	doc.DrawPrincipalApproval(school.PrincipalName)

	return doc.PDF.WritePdf(outputPath)
}

// GenerateDiziPusulasi generates the case file index (dizi pusulası).
func GenerateDiziPusulasi(school *model.SchoolInfo, items []model.DiziItem, outputPath string) error {
	doc, err := NewDoc(false)
	if err != nil {
		return err
	}

	doc.PDF.AddPage()

	boardName := "ÖDÜL VE DİSİPLİN"
	if strings.ToUpper(school.SchoolType) == "ORTAOKUL" {
		boardName = "ÖĞRENCİ DAVRANIŞLARINI DEĞERLENDİRME"
	}

	doc.DrawTitle(fmt.Sprintf("%s MÜDÜRLÜĞÜ\n%s KURULU DİZİ PUSULASI", strings.ToUpper(school.SchoolName), boardName), 12, 20)

	indent := 40.0
	tableW := doc.UsableWidth - (indent * 2)
	wSira := 40.0
	wCount := 60.0
	wSubject := tableW - wSira - wCount

	curY := doc.PDF.GetY()
	doc.DrawRect(doc.MarginLeft+indent, curY, wSira, 20)
	doc.DrawRect(doc.MarginLeft+indent+wSira, curY, wSubject, 20)
	doc.DrawRect(doc.MarginLeft+indent+wSira+wSubject, curY, wCount, 20)

	doc.SetBoldFont(8.0)
	doc.PDF.SetX(doc.MarginLeft + indent)
	doc.PDF.SetY(curY + 4)
	_ = doc.PDF.CellWithOption(&gopdf.Rect{W: wSira, H: 12}, "SIRA", gopdf.CellOption{Align: gopdf.Center})

	doc.PDF.SetX(doc.MarginLeft + indent + wSira + 5)
	doc.PDF.SetY(curY + 4)
	_ = doc.PDF.CellWithOption(&gopdf.Rect{W: wSubject - 10, H: 12}, "KONUSU", gopdf.CellOption{Align: gopdf.Left})

	doc.PDF.SetX(doc.MarginLeft + indent + wSira + wSubject)
	doc.PDF.SetY(curY + 4)
	_ = doc.PDF.CellWithOption(&gopdf.Rect{W: wCount, H: 12}, "ADET / SAYI", gopdf.CellOption{Align: gopdf.Center})

	curY += 20
	doc.SetRegularFont(8.0)
	totalCount := 0
	lineH := 11.5

	for i, item := range items {
		subLines := doc.SplitLines(item.Subject, wSubject-10)
		rowH := float64(len(subLines))*lineH + 6.0
		if rowH < 18.0 {
			rowH = 18.0
		}

		doc.DrawRect(doc.MarginLeft+indent, curY, wSira, rowH)
		doc.DrawRect(doc.MarginLeft+indent+wSira, curY, wSubject, rowH)
		doc.DrawRect(doc.MarginLeft+indent+wSira+wSubject, curY, wCount, rowH)

		doc.PDF.SetX(doc.MarginLeft + indent)
		doc.PDF.SetY(curY + 3)
		_ = doc.PDF.CellWithOption(&gopdf.Rect{W: wSira, H: 12}, fmt.Sprintf("%d", i+1), gopdf.CellOption{Align: gopdf.Center})

		doc.DrawWrappedCellLines(doc.MarginLeft+indent+wSira+5, curY+3, wSubject-10, lineH, subLines, gopdf.Left)

		doc.PDF.SetX(doc.MarginLeft + indent + wSira + wSubject)
		doc.PDF.SetY(curY + 3)
		_ = doc.PDF.CellWithOption(&gopdf.Rect{W: wCount, H: 12}, fmt.Sprintf("%d", item.Count), gopdf.CellOption{Align: gopdf.Center})

		totalCount += item.Count
		curY += rowH
	}

	// Total Row
	doc.DrawRect(doc.MarginLeft+indent, curY, wSira+wSubject, 20)
	doc.DrawRect(doc.MarginLeft+indent+wSira+wSubject, curY, wCount, 20)

	doc.SetBoldFont(8.0)
	doc.PDF.SetX(doc.MarginLeft + indent)
	doc.PDF.SetY(curY + 4)
	_ = doc.PDF.CellWithOption(&gopdf.Rect{W: wSira + wSubject - 10, H: 12}, "TOPLAM", gopdf.CellOption{Align: gopdf.Right})

	doc.PDF.SetX(doc.MarginLeft + indent + wSira + wSubject)
	doc.PDF.SetY(curY + 4)
	_ = doc.PDF.CellWithOption(&gopdf.Rect{W: wCount, H: 12}, fmt.Sprintf("%d", totalCount), gopdf.CellOption{Align: gopdf.Center})

	// Handover signatures
	yHand := curY + 50
	colW := doc.UsableWidth / 2.0

	doc.SetBoldFont(8.0)
	doc.PDF.SetX(doc.MarginLeft)
	doc.PDF.SetY(yHand)
	_ = doc.PDF.CellWithOption(&gopdf.Rect{W: colW, H: 12}, "Teslim Eden", gopdf.CellOption{Align: gopdf.Center})

	doc.PDF.SetX(doc.MarginLeft + colW)
	doc.PDF.SetY(yHand)
	_ = doc.PDF.CellWithOption(&gopdf.Rect{W: colW, H: 12}, "Teslim Alan", gopdf.CellOption{Align: gopdf.Center})

	doc.SetRegularFont(7.5)
	doc.PDF.SetX(doc.MarginLeft)
	doc.PDF.SetY(yHand + 30)
	_ = doc.PDF.CellWithOption(&gopdf.Rect{W: colW, H: 12}, "İmza: ........................................", gopdf.CellOption{Align: gopdf.Center})

	doc.PDF.SetX(doc.MarginLeft + colW)
	doc.PDF.SetY(yHand + 30)
	_ = doc.PDF.CellWithOption(&gopdf.Rect{W: colW, H: 12}, "İmza: ........................................", gopdf.CellOption{Align: gopdf.Center})

	return doc.PDF.WritePdf(outputPath)
}

// GenerateSeneBasiSonuKarar generates the official beginning/end of semester meeting minutes.
func GenerateSeneBasiSonuKarar(school *model.SchoolInfo, meetingType string, dateStr, timeStr, place, chair string, decisions []string, outputPath string) error {
	doc, err := NewDoc(false)
	if err != nil {
		return err
	}

	doc.PDF.AddPage()

	boardName := "ÖDÜL VE DİSİPLİN"
	statute := "Milli Eğitim Bakanlığı Ortaöğretim Kurumları Yönetmeliğinin 189. Maddesi gereğince"
	if strings.ToUpper(school.SchoolType) == "ORTAOKUL" {
		boardName = "ÖĞRENCİ DAVRANIŞLARINI DEĞERLENDİRME"
		statute = "Milli Eğitim Bakanlığı Okul Öncesi Eğitim ve İlköğretim Kurumları Yönetmeliğinin 58. Maddesi gereğince"
	}

	// Table 1: Title
	doc.DrawTitle(fmt.Sprintf("%s KURULU %s TOPLANTI TUTANAĞI", boardName, strings.ToUpper(meetingType)), 13, 15)

	// Table 2: Details
	lblW := 140.0
	valW := doc.UsableWidth - lblW

	rows := [][2]string{
		{"OKUL ADI", ": " + school.SchoolName},
		{"ÖĞRETİM YILI", ": " + school.AcademicYear},
		{"TOPLANTI TARİHİ", ": " + dateStr},
		{"TOPLANTI SAATİ", ": " + timeStr},
		{"TOPLANTI YERİ", ": " + place},
		{"TOPLANTI BAŞKANI", ": " + chair},
	}

	for _, r := range rows {
		rowY := doc.PDF.GetY()
		doc.SetBoldFont(8.0)
		doc.PDF.SetX(doc.MarginLeft)
		doc.PDF.SetY(rowY)
		_ = doc.PDF.CellWithOption(&gopdf.Rect{W: lblW, H: 13}, r[0], gopdf.CellOption{Align: gopdf.Left})

		doc.SetRegularFont(8.0)
		valLines := doc.SplitLines(r[1], valW)
		doc.DrawWrappedCellLines(doc.MarginLeft+lblW, rowY, valW, 13.0, valLines, gopdf.Left)
		doc.PDF.SetY(rowY + float64(len(valLines))*13.0 + 1.0)
	}

	doc.PDF.Br(10)

	// Statutory intro
	doc.SetRegularFont(8.5)
	intro := fmt.Sprintf("%s kurulumuz yukarıda belirtilen yer ve zamanda toplanarak gündemdeki maddeleri görüşmüş ve aşağıdaki kararları almıştır:", statute)
	doc.DrawWrappedText(doc.MarginLeft, doc.UsableWidth, 13.5, intro, gopdf.Left)
	doc.PDF.Br(10)

	// Numbered Decisions
	for i, d := range decisions {
		d = strings.TrimSpace(d)
		if d == "" {
			continue
		}

		if doc.PDF.GetY() > PageHeightPortrait-doc.MarginBottom-120 {
			doc.PDF.AddPage()
		}

		itemY := doc.PDF.GetY()
		doc.SetBoldFont(8.0)
		doc.PDF.SetX(doc.MarginLeft + 10)
		doc.PDF.SetY(itemY)
		_ = doc.PDF.CellWithOption(&gopdf.Rect{W: 30, H: 13}, fmt.Sprintf("%d.)", i+1), gopdf.CellOption{Align: gopdf.Right})

		doc.SetRegularFont(8.0)
		doc.PDF.SetY(itemY)
		doc.DrawWrappedText(doc.MarginLeft+45, doc.UsableWidth-55, 13.0, d, gopdf.Left)
		doc.PDF.Br(4)
	}

	doc.PDF.Br(10)

	// Closing clause
	doc.SetBoldFont(9.0)
	closing := "OY BİRLİĞİ İLE KARAR VERİLDİ."
	if strings.Contains(strings.ToUpper(meetingType), "SONU") {
		closing = "İŞ BU KURUL TUTANAĞININ SENE SONU ÖĞRETMENLER KURULUNDA OKUNMASI TALEBİYLE OKUL MÜDÜRLÜĞÜNÜN ONAYINA SUNULMASINA OY BİRLİĞİ İLE KARAR VERİLDİ."
	}
	doc.DrawWrappedText(doc.MarginLeft, doc.UsableWidth, 14.0, closing, gopdf.Left)
	doc.PDF.Br(12)

	doc.DrawSignatures5(school, boardName)
	doc.DrawPrincipalApproval(school.PrincipalName)

	return doc.PDF.WritePdf(outputPath)
}

// GenerateCezaIstatistikRaporu generates a penalty statistics summary report.
func GenerateCezaIstatistikRaporu(school *model.SchoolInfo, title string, stats []model.StatItem, outputPath string) error {
	doc, err := NewDoc(false)
	if err != nil {
		return err
	}

	doc.PDF.AddPage()
	doc.DrawAntet(school)
	doc.DrawTitle(title, 12, 15)

	wSira := 40.0
	wCount := 80.0
	wName := doc.UsableWidth - wSira - wCount

	curY := doc.PDF.GetY()
	doc.DrawRect(doc.MarginLeft, curY, wSira, 20)
	doc.DrawRect(doc.MarginLeft+wSira, curY, wName, 20)
	doc.DrawRect(doc.MarginLeft+wSira+wName, curY, wCount, 20)

	doc.SetBoldFont(8.0)
	doc.PDF.SetX(doc.MarginLeft)
	doc.PDF.SetY(curY + 4)
	_ = doc.PDF.CellWithOption(&gopdf.Rect{W: wSira, H: 12}, "SIRA", gopdf.CellOption{Align: gopdf.Center})

	doc.PDF.SetX(doc.MarginLeft + wSira + 5)
	doc.PDF.SetY(curY + 4)
	_ = doc.PDF.CellWithOption(&gopdf.Rect{W: wName - 10, H: 12}, "BAŞLIK / KAPSAM", gopdf.CellOption{Align: gopdf.Left})

	doc.PDF.SetX(doc.MarginLeft + wSira + wName)
	doc.PDF.SetY(curY + 4)
	_ = doc.PDF.CellWithOption(&gopdf.Rect{W: wCount, H: 12}, "SAYI", gopdf.CellOption{Align: gopdf.Center})

	curY += 20
	doc.SetRegularFont(8.0)
	total := 0
	lineH := 11.5

	for i, s := range stats {
		nameLines := doc.SplitLines(s.Name, wName-10)
		rowH := float64(len(nameLines))*lineH + 6.0
		if rowH < 18.0 {
			rowH = 18.0
		}

		if curY+rowH > PageHeightPortrait-doc.MarginBottom-60 {
			doc.PDF.AddPage()
			curY = doc.MarginTop
			doc.SetRegularFont(8.0)
		}

		doc.DrawRect(doc.MarginLeft, curY, wSira, rowH)
		doc.DrawRect(doc.MarginLeft+wSira, curY, wName, rowH)
		doc.DrawRect(doc.MarginLeft+wSira+wName, curY, wCount, rowH)

		doc.PDF.SetX(doc.MarginLeft)
		doc.PDF.SetY(curY + 3)
		_ = doc.PDF.CellWithOption(&gopdf.Rect{W: wSira, H: 12}, fmt.Sprintf("%d", i+1), gopdf.CellOption{Align: gopdf.Center})

		doc.DrawWrappedCellLines(doc.MarginLeft+wSira+5, curY+3, wName-10, lineH, nameLines, gopdf.Left)

		doc.PDF.SetX(doc.MarginLeft + wSira + wName)
		doc.PDF.SetY(curY + 3)
		_ = doc.PDF.CellWithOption(&gopdf.Rect{W: wCount, H: 12}, fmt.Sprintf("%d", s.Count), gopdf.CellOption{Align: gopdf.Center})

		total += s.Count
		curY += rowH
	}

	// Total Row
	doc.DrawRect(doc.MarginLeft, curY, wSira+wName, 20)
	doc.DrawRect(doc.MarginLeft+wSira+wName, curY, wCount, 20)

	doc.SetBoldFont(8.0)
	doc.PDF.SetX(doc.MarginLeft)
	doc.PDF.SetY(curY + 4)
	_ = doc.PDF.CellWithOption(&gopdf.Rect{W: wSira + wName - 10, H: 12}, "TOPLAM", gopdf.CellOption{Align: gopdf.Right})

	doc.PDF.SetX(doc.MarginLeft + wSira + wName)
	doc.PDF.SetY(curY + 4)
	_ = doc.PDF.CellWithOption(&gopdf.Rect{W: wCount, H: 12}, fmt.Sprintf("%d", total), gopdf.CellOption{Align: gopdf.Center})

	doc.PDF.SetY(curY + 35)
	boardName := "Ödül ve Disiplin"
	if strings.ToUpper(school.SchoolType) == "ORTAOKUL" {
		boardName = "Öğrenci Davranışlarını Değerlendirme"
	}
	doc.DrawSignatures5(school, boardName)
	doc.DrawPrincipalApproval(school.PrincipalName)

	return doc.PDF.WritePdf(outputPath)
}
