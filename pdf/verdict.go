package pdf

import (
	"fmt"
	"strings"

	"github.com/signintech/gopdf"

	"turandisiplin/model"
)

// GenerateKurulKarari generates the official EK-1 (Lise) or EK-10 (Ortaokul) Disciplinary Board Decision Document.
func GenerateKurulKarari(school *model.SchoolInfo, student *model.Student, incident *model.Incident, incStudent *model.IncidentStudent, verdict *model.Verdict, outputPath string) error {
	doc, err := NewDoc(false)
	if err != nil {
		return err
	}

	doc.PDF.AddPage()

	isOrta := strings.ToUpper(school.SchoolType) == "ORTAOKUL"
	boardName := "ÖDÜL VE DİSİPLİN"
	ekCode := "EK-1"
	if isOrta {
		boardName = "ÖĞRENCİ DAVRANIŞLARINI DEĞERLENDİRME"
		ekCode = "EK-10"
	}

	// 1. Top Right EK tag
	doc.SetBoldFont(9)
	doc.PDF.SetX(doc.MarginLeft)
	_ = doc.PDF.CellWithOption(&gopdf.Rect{W: doc.UsableWidth, H: 12}, ekCode, gopdf.CellOption{Align: gopdf.Right})
	doc.PDF.Br(14)

	// 2. Centered Title
	doc.DrawTitle(boardName+" KURULU KARARI", 12, 10)

	// 3. 3-Column Main Table: [150pt, 10pt, remaining pt]
	wCol1 := 145.0
	wCol2 := 10.0
	wCol3 := doc.UsableWidth - wCol1 - wCol2

	// Gather rows
	type tableRow struct {
		Label string
		Value string
		Bold  bool
	}

	var rows []tableRow

	kararNoStr := ""
	kararTarihStr := ""
	if verdict != nil {
		kararNoStr = fmt.Sprintf("%d", verdict.DecisionNo)
		kararTarihStr = verdict.DecisionDate
	}

	rows = append(rows, tableRow{Label: "Karar No", Value: kararNoStr})
	rows = append(rows, tableRow{Label: "Karar Tarihi", Value: kararTarihStr})

	if student != nil {
		rows = append(rows, tableRow{Label: "Öğrencinin Adı ve Soyadı", Value: student.FullName, Bold: true})
		rows = append(rows, tableRow{Label: "Sınıfı ve Numarası", Value: fmt.Sprintf("%s  -  %s", student.ClassSection, student.SchoolNo)})
		rows = append(rows, tableRow{Label: "Cinsiyeti", Value: student.Gender})
		rows = append(rows, tableRow{Label: "T.C. Kimlik No", Value: student.TCNo})
		rows = append(rows, tableRow{Label: "Doğum Yeri ve Tarihi", Value: fmt.Sprintf("%s / %s", student.BirthPlace, model.NormalizeBirthDate(student.BirthDate))})
		rows = append(rows, tableRow{Label: "Baba ve Anne Adı", Value: fmt.Sprintf("%s - %s", student.FatherName, student.MotherName)})
		rows = append(rows, tableRow{Label: "Velisinin Adı Soyadı", Value: student.GuardianName})
		rows = append(rows, tableRow{Label: "Velisinin Telefon ve Adresi", Value: fmt.Sprintf("%s  |  %s", student.GuardianPhone, student.GuardianAddr)})
		rows = append(rows, tableRow{Label: "Pansiyon / Gündüzlü Durumu", Value: student.BoardingStatus})
	}

	if incStudent != nil {
		rows = append(rows, tableRow{Label: "Başarı Durumu", Value: incStudent.AcademicStanding})
		rows = append(rows, tableRow{Label: "Sağlık Durumu", Value: incStudent.HealthStatus})
		rows = append(rows, tableRow{Label: "Ailenin Ekonomik Durumu", Value: incStudent.EconomicStatus})
		rows = append(rows, tableRow{Label: "Kiminle Kaldığı", Value: incStudent.LivingWith})
		rows = append(rows, tableRow{Label: "Anne - Baba Durumu", Value: fmt.Sprintf("%s - %s", incStudent.ParentsLivingStatus, incStudent.ParentsBioStatus)})
		rows = append(rows, tableRow{Label: "Daha Önce Okuduğu Okullar", Value: incStudent.GraduatedSchool})
		rows = append(rows, tableRow{Label: "Büyüdüğü Çevre", Value: incStudent.EnvironmentRaised})
		rows = append(rows, tableRow{Label: "Daha Önce Ceza / Yaptırım Durumu", Value: incStudent.PriorPenaltyStatus})
	}

	if incident != nil {
		rows = append(rows, tableRow{Label: "Olayın Adı", Value: incident.Title})
		rows = append(rows, tableRow{Label: "Olayın Tarihi ve Yeri", Value: fmt.Sprintf("%s - %s", incident.Date, incident.Location)})
		rows = append(rows, tableRow{Label: "Olayın Özeti", Value: incident.Summary})
	}

	if verdict != nil {
		rows = append(rows, tableRow{Label: "Suçun Çeşidi", Value: verdict.OffenseType})
		rows = append(rows, tableRow{Label: "Suçun Nedeni", Value: verdict.OffenseReason})
		rows = append(rows, tableRow{Label: "Tanık İfadeleri ve Deliller", Value: fmt.Sprintf("Tanıklar: %s | Deliller: %s", verdict.WitnessStatements, verdict.EvidenceDetails)})
		rows = append(rows, tableRow{Label: "Hafifletici / Ağırlaştırıcı Nedenler", Value: fmt.Sprintf("Hafifleten: %s | Ağırlaştıran: %s", verdict.MitigatingFactors, verdict.AggravatingFactors)})
		rows = append(rows, tableRow{Label: "Kurulun Kanaati", Value: verdict.BoardOpinion})
		rows = append(rows, tableRow{Label: "Uygulanan Yönetmelik Maddesi", Value: fmt.Sprintf("%s (%s)", verdict.RegClause, verdict.RegClauseDesc)})
		rows = append(rows, tableRow{Label: "Verilen Ceza / Yaptırım", Value: verdict.PenaltyName, Bold: true})
		if verdict.SuspensionDays > 0 {
			rows = append(rows, tableRow{Label: "Uzaklaştırma Süresi ve Tarihi", Value: fmt.Sprintf("%d Gün  (%s)", verdict.SuspensionDays, verdict.SuspensionDates)})
		}
		if !isOrta {
			rows = append(rows, tableRow{Label: "İndirilen / Kalan Davranış Puanı", Value: fmt.Sprintf("%d Puan İndirildi / Kalan Puan: %d", verdict.DeductedPoints, verdict.RemainingPoints)})
		}
		rows = append(rows, tableRow{Label: "Kurul Kararı", Value: verdict.DecisionText, Bold: true})
	}

	fontSize := 7.0
	lineH := 10.5
	doc.PDF.SetLineWidth(0.3)

	for _, r := range rows {
		val := strings.TrimSpace(r.Value)
		if val == "" {
			val = "—"
		}

		// Set font first so gopdf measures text glyph widths accurately
		if r.Bold {
			doc.SetBoldFont(fontSize)
		} else {
			doc.SetRegularFont(fontSize)
		}

		lines, err := doc.PDF.SplitTextWithWordWrap(val, wCol3-8)
		if err != nil || len(lines) == 0 {
			lines = []string{val}
		}

		cellH := float64(len(lines))*lineH + 5.0
		if cellH < 14.0 {
			cellH = 14.0
		}

		// Check if we need a new page
		curY := doc.PDF.GetY()
		if curY+cellH > PageHeightPortrait-doc.MarginBottom-60 {
			doc.PDF.AddPage()
			curY = doc.MarginTop
			doc.PDF.SetY(curY)
		}

		doc.DrawRect(doc.MarginLeft, curY, wCol1, cellH)
		doc.DrawRect(doc.MarginLeft+wCol1, curY, wCol2, cellH)
		doc.DrawRect(doc.MarginLeft+wCol1+wCol2, curY, wCol3, cellH)

		// Col 1: Label
		doc.SetBoldFont(fontSize)
		doc.PDF.SetX(doc.MarginLeft + 3)
		doc.PDF.SetY(curY + 2)
		_ = doc.PDF.CellWithOption(&gopdf.Rect{W: wCol1 - 6, H: cellH - 4}, r.Label, gopdf.CellOption{Align: gopdf.Left})

		// Col 2: ":"
		doc.PDF.SetX(doc.MarginLeft + wCol1)
		doc.PDF.SetY(curY + 2)
		_ = doc.PDF.CellWithOption(&gopdf.Rect{W: wCol2, H: cellH - 4}, ":", gopdf.CellOption{Align: gopdf.Center})

		// Col 3: Exact wrapped lines rendered safely within boundaries
		if r.Bold {
			doc.SetBoldFont(fontSize)
		} else {
			doc.SetRegularFont(fontSize)
		}
		for i, line := range lines {
			doc.PDF.SetX(doc.MarginLeft + wCol1 + wCol2 + 4)
			doc.PDF.SetY(curY + 2.5 + float64(i)*lineH)
			_ = doc.PDF.CellWithOption(&gopdf.Rect{W: wCol3 - 8, H: lineH}, line, gopdf.CellOption{Align: gopdf.Left})
		}

		doc.PDF.SetY(curY + cellH)
	}

	// Signatures: check space
	curY := doc.PDF.GetY()
	if curY > PageHeightPortrait-doc.MarginBottom-110 {
		doc.PDF.AddPage()
	}

	// Substitute Member replacement check:
	// If the reporting teacher or substitute member override was used
	effectiveChair := school.BoardChair
	effectiveChairTitle := school.BoardChairTitle
	effectiveM1 := school.Member1
	effectiveM1Title := school.Member1Title
	effectiveM2 := school.Member2
	effectiveM2Title := school.Member2Title
	effectiveM3 := school.Member3
	effectiveM3Title := school.Member3Title
	effectiveM4 := school.Member4
	effectiveM4Title := school.Member4Title

	if verdict != nil {
		if verdict.ChairIsSubstitute && school.SubBoardChair != "" {
			effectiveChair = school.SubBoardChair
			effectiveChairTitle = school.SubBoardChairTitle
		}
		if verdict.M1IsSubstitute && school.SubMember1 != "" {
			effectiveM1 = school.SubMember1
			effectiveM1Title = school.SubMember1Title
		}
		if verdict.M2IsSubstitute && school.SubMember2 != "" {
			effectiveM2 = school.SubMember2
			effectiveM2Title = school.SubMember2Title
		}
		if verdict.M3IsSubstitute && school.SubMember3 != "" {
			effectiveM3 = school.SubMember3
			effectiveM3Title = school.SubMember3Title
		}
		if verdict.M4IsSubstitute && school.SubMember4 != "" {
			effectiveM4 = school.SubMember4
			effectiveM4Title = school.SubMember4Title
		}
	}

	tempSchool := *school
	tempSchool.BoardChair = effectiveChair
	tempSchool.BoardChairTitle = effectiveChairTitle
	tempSchool.Member1 = effectiveM1
	tempSchool.Member1Title = effectiveM1Title
	tempSchool.Member2 = effectiveM2
	tempSchool.Member2Title = effectiveM2Title
	tempSchool.Member3 = effectiveM3
	tempSchool.Member3Title = effectiveM3Title
	tempSchool.Member4 = effectiveM4
	tempSchool.Member4Title = effectiveM4Title

	doc.DrawSignatures5(&tempSchool, boardName)
	doc.DrawPrincipalApproval(school.PrincipalName)

	return doc.PDF.WritePdf(outputPath)
}
