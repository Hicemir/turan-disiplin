package pdf

import (
	"fmt"
	"strings"

	"github.com/signintech/gopdf"

	"turandisiplin/model"
)

// toUpperTR converts a string to uppercase using Turkish locale rules ('i' -> 'İ', 'ı' -> 'I').
func toUpperTR(s string) string {
	var sb strings.Builder
	for _, r := range s {
		switch r {
		case 'i':
			sb.WriteRune('İ')
		case 'ı':
			sb.WriteRune('I')
		case 'ğ':
			sb.WriteRune('Ğ')
		case 'ü':
			sb.WriteRune('Ü')
		case 'ş':
			sb.WriteRune('Ş')
		case 'ö':
			sb.WriteRune('Ö')
		case 'ç':
			sb.WriteRune('Ç')
		default:
			sb.WriteString(strings.ToUpper(string(r)))
		}
	}
	return sb.String()
}

func cleanJoin(a, b, sep string) string {
	a = strings.TrimSpace(a)
	b = strings.TrimSpace(b)
	if a != "" && b != "" {
		return a + sep + b
	}
	if a != "" {
		return a
	}
	return b
}

// GenerateKurulKarari generates the official EK-1 (Lise) or EK-10 (Ortaokul) Disciplinary Board Decision Document.
func GenerateKurulKarari(school *model.SchoolInfo, student *model.Student, incident *model.Incident, incStudent *model.IncidentStudent, verdict *model.Verdict, outputPath string) error {
	doc, err := NewDoc(false)
	if err != nil {
		return err
	}
	doc.MarginTop = 18.0
	doc.MarginBottom = 18.0
	doc.PDF.SetMargins(doc.MarginLeft, doc.MarginTop, doc.MarginRight, doc.MarginBottom)

	doc.PDF.AddPage()

	isOrta := strings.ToUpper(strings.TrimSpace(school.SchoolType)) == "ORTAOKUL"
	titleText := "OKUL ÖĞRENCİ ÖDÜL VE DİSİPLİN KURULU KARARI"
	ekCode := "EK-1"
	if isOrta {
		titleText = "ÖĞRENCİ DAVRANIŞLARINI DEĞERLENDİRME KURULU KARARI"
		ekCode = "EK-10"
	}

	// 1. Top Right EK tag (8pt)
	doc.PDF.SetY(doc.MarginTop)
	doc.SetRegularFont(8)
	doc.PDF.SetX(doc.MarginLeft)
	_ = doc.PDF.CellWithOption(&gopdf.Rect{W: doc.UsableWidth, H: 10}, ekCode, gopdf.CellOption{Align: gopdf.Right})
	doc.PDF.Br(10)

	// 2. Centered Title (12pt Bold)
	doc.SetBoldFont(11.5)
	doc.PDF.SetX(doc.MarginLeft)
	_ = doc.PDF.CellWithOption(&gopdf.Rect{W: doc.UsableWidth, H: 14}, titleText, gopdf.CellOption{Align: gopdf.Center})
	doc.PDF.Br(16)

	// 3. Column proportions matching legacy Frm07Karar.cs: [12, 1, 20]
	wCol1 := doc.UsableWidth * 12.0 / 33.0
	wCol2 := doc.UsableWidth * 1.0 / 33.0
	wCol3 := doc.UsableWidth - wCol1 - wCol2

	type tableRow struct {
		Label         string
		Value         string
		IsSection     bool
		IsParentTable bool
		BoldValue     bool
		MinHeight     float64
	}

	var rows []tableRow

	kararNoStr := ""
	kararTarihStr := ""
	if verdict != nil {
		if verdict.DecisionNo > 0 {
			kararNoStr = fmt.Sprintf("%d", verdict.DecisionNo)
		}
		kararTarihStr = strings.TrimSpace(verdict.DecisionDate)
	}
	if kararNoStr == "" && incStudent != nil {
		kararNoStr = strings.TrimSpace(incStudent.VerdictNo)
	}
	if kararTarihStr == "" && incStudent != nil {
		kararTarihStr = strings.TrimSpace(incStudent.VerdictDate)
	}

	stFullName := ""
	stTC := ""
	stBirthDate := ""
	stBirthPlace := ""
	stClass := ""
	stNo := ""
	stGender := ""
	stBoarding := ""
	stFather := ""
	stMother := ""
	stGuardianName := ""
	stGuardianPhone := ""
	stGuardianAddr := ""

	if student != nil {
		stFullName = student.FullName
		stTC = student.TCNo
		stBirthDate = model.NormalizeBirthDate(student.BirthDate)
		stBirthPlace = student.BirthPlace
		stClass = student.ClassSection
		stNo = student.SchoolNo
		stGender = student.Gender
		stBoarding = student.BoardingStatus
		stFather = student.FatherName
		stMother = student.MotherName
		stGuardianName = student.GuardianName
		stGuardianPhone = student.GuardianPhone
		stGuardianAddr = student.GuardianAddr
	}
	if incStudent != nil {
		if stFullName == "" {
			stFullName = incStudent.FullName
		}
		if stTC == "" {
			stTC = incStudent.TCNo
		}
		if stBirthDate == "" {
			stBirthDate = model.NormalizeBirthDate(incStudent.BirthDate)
		}
		if stBirthPlace == "" {
			stBirthPlace = incStudent.BirthPlace
		}
		if stClass == "" {
			stClass = incStudent.ClassSection
		}
		if stNo == "" {
			stNo = incStudent.SchoolNo
		}
		if stGender == "" {
			stGender = incStudent.Gender
		}
		if stBoarding == "" {
			stBoarding = incStudent.BoardingStatus
		}
		if stFather == "" {
			stFather = incStudent.FatherName
		}
		if stMother == "" {
			stMother = incStudent.MotherName
		}
		if stGuardianName == "" {
			stGuardianName = incStudent.GuardianName
		}
		if stGuardianPhone == "" {
			stGuardianPhone = incStudent.GuardianPhone
		}
		if stGuardianAddr == "" {
			stGuardianAddr = incStudent.GuardianAddr
		}
	}

	academicStanding := ""
	healthStatus := ""
	economicStatus := ""
	livingWith := ""
	parentsLiving := ""
	parentsBio := ""
	schoolCommute := ""
	envRaised := ""
	priorPenalty := ""
	defenseText := ""
	siblingCount := ""

	if incStudent != nil {
		academicStanding = incStudent.AcademicStanding
		healthStatus = incStudent.HealthStatus
		economicStatus = incStudent.EconomicStatus
		livingWith = incStudent.LivingWith
		parentsLiving = incStudent.ParentsLivingStatus
		parentsBio = incStudent.ParentsBioStatus
		schoolCommute = incStudent.SchoolCommute
		envRaised = incStudent.EnvironmentRaised
		priorPenalty = incStudent.PriorPenaltyStatus
		siblingCount = incStudent.SiblingCount
		defenseText = strings.TrimSpace(incStudent.DefenseSummary)
		if defenseText == "" {
			defenseText = strings.TrimSpace(incStudent.StatementText)
		}
	}
	if strings.TrimSpace(schoolCommute) == "" && strings.Contains(toUpperTR(stBoarding), "GÜNDÜZLÜ") {
		schoolCommute = "AİLESİNİN YANINDA OKUYOR"
	}

	// Format address + phone (+ guardian name fallback if address empty)
	addrVal := strings.TrimSpace(stGuardianAddr)
	phoneVal := strings.TrimSpace(stGuardianPhone)
	if phoneVal != "" {
		if addrVal != "" {
			addrVal += "\nTel No:" + phoneVal
		} else {
			addrVal = "Tel No:" + phoneVal
		}
	}

	offensePlaceTime := ""
	offenseType := ""
	offenseReason := ""
	witnessText := ""
	evidenceText := ""
	mitigatingText := ""
	boardOpinion := ""
	penaltyAndClause := ""
	decisionText := ""

	if verdict != nil {
		place := strings.TrimSpace(verdict.OffensePlace)
		dt := strings.TrimSpace(verdict.OffenseDateTime)
		if place == "" && incident != nil {
			place = strings.TrimSpace(incident.Location)
		}
		if dt == "" && incident != nil {
			dt = cleanJoin(incident.Date, incident.Time, " - ")
		}
		offensePlaceTime = cleanJoin(place, dt, " / ")

		offenseType = toUpperTR(strings.TrimSpace(verdict.OffenseType))
		if offenseType == "" && incident != nil {
			offenseType = toUpperTR(strings.TrimSpace(incident.Title))
		}
		offenseReason = strings.TrimSpace(verdict.OffenseReason)
		witnessText = strings.TrimSpace(verdict.WitnessStatements)
		evidenceText = strings.TrimSpace(verdict.EvidenceDetails)

		mit := strings.TrimSpace(verdict.MitigatingFactors)
		agg := strings.TrimSpace(verdict.AggravatingFactors)
		if strings.EqualFold(agg, "Herhangi bir ağırlaştırıcı neden bulunmamaktadır") || strings.EqualFold(agg, "Yok") {
			agg = ""
		}
		if strings.EqualFold(mit, "Herhangi bir hafifletici neden bulunmamaktadır") && agg != "" {
			mit = ""
		}
		mitigatingText = toUpperTR(cleanJoin(mit, agg, " "))

		boardOpinion = toUpperTR(strings.TrimSpace(verdict.BoardOpinion))

		penLine := toUpperTR(strings.TrimSpace(verdict.PenaltyName))
		if verdict.SuspensionDays > 0 && !strings.HasPrefix(penLine, fmt.Sprintf("%d", verdict.SuspensionDays)) {
			penLine = fmt.Sprintf("%d GÜN %s", verdict.SuspensionDays, penLine)
		}
		clauseLine := strings.TrimSpace(verdict.RegClause)
		if strings.TrimSpace(verdict.RegClauseDesc) != "" {
			clauseLine = cleanJoin(clauseLine, "("+strings.TrimSpace(verdict.RegClauseDesc)+")", " ")
		}
		penaltyAndClause = cleanJoin(penLine, clauseLine, "\n")

		decisionText = toUpperTR(strings.TrimSpace(verdict.DecisionText))
	}

	if !isOrta {
		// Authentic EK-1 (Lise) rows matching official scan
		rows = []tableRow{
			{Label: "Karar No", Value: kararNoStr},
			{Label: "Karar Tarihi", Value: kararTarihStr},
			{Label: "       ÖĞRENCİNİN;", IsSection: true},
			{Label: "TC KİMLİK NO", Value: stTC},
			{Label: "Adı Soyadı", Value: toUpperTR(stFullName)},
			{Label: "Doğum Tarihi ve yeri", Value: cleanJoin(stBirthDate, stBirthPlace, " - ")},
			{Label: "Sınıfı, alanı/dalı ve okul numarası", Value: cleanJoin(stClass, stNo, " - ")},
			{Label: "Paralı veya parasız yatılı ya da gündüzlü olduğu", Value: toUpperTR(stBoarding)},
			{Label: "Başarı durumu", Value: toUpperTR(academicStanding)},
			{Label: "Sağlık durumu", Value: toUpperTR(healthStatus)},
			{Label: "Ailesinin ekonomik durumu", Value: toUpperTR(economicStatus)},
			{Label: "Ailesi ile birlikte oturup oturmadığı", Value: toUpperTR(livingWith)},
			{Label: "Anne-babasının sağ olup olmadığı", Value: toUpperTR(parentsLiving)},
			{Label: "Anne-babasının öz olup olmadığı", Value: toUpperTR(parentsBio)},
			{Label: "Ailesinin yanında okuyup okumadığı", Value: toUpperTR(schoolCommute)},
			{Label: "Büyüyüp yetiştiği çevre", Value: toUpperTR(envRaised)},
			{Label: "Ailesinin oturduğu yer / ikamet adresi", Value: addrVal, MinHeight: 22},
			{Label: "Şimdiye kadar aldığı cezalar ve genel durumu", Value: toUpperTR(priorPenalty)},
			{Label: "Cezayı gerektiren davranışının yapıldığı yer ve tarih", Value: offensePlaceTime},
			{Label: "Cezayı gerektiren davranışının çeşidi", Value: offenseType},
			{Label: "Cezayı gerektiren davranışının nedeni", Value: offenseReason},
			{Label: "    Olayla ilgili olarak;", IsSection: true},
			{Label: "a) Cezalandırılan öğrencinin ifadesinin özeti", Value: defenseText, MinHeight: 18},
			{Label: "b) Tanıkların ifadesinin özeti", Value: witnessText, MinHeight: 18},
			{Label: "c) Varsa cezayı gerektiren davranışının tespitine yarayan diğer deliller", Value: evidenceText, MinHeight: 18},
			{Label: "Cezayı hafifleten veya şiddetlendiren nedenler", Value: mitigatingText, MinHeight: 18},
			{Label: "Okul öğrenci ödül ve disiplin kurulunun kanaati", Value: boardOpinion, MinHeight: 18},
			{Label: "Verilen cezanın çeşidi ve dayandığı yönetmelik maddesi", Value: penaltyAndClause, MinHeight: 25},
			{Label: "OKUL ÖĞRENCİ ÖDÜL VE DİSİPLİN KURULU KARARI", Value: decisionText, BoldValue: true, MinHeight: 45},
		}
	} else {
		// Authentic EK-10 (Ortaokul) rows matching Frm07KararOrta.cs
		rows = []tableRow{
			{Label: "Karar No", Value: kararNoStr},
			{Label: "Karar Tarihi", Value: kararTarihStr},
			{Label: "       ÖĞRENCİNİN;", IsSection: true},
			{Label: "Adı Soyadı", Value: stFullName},
			{Label: "Cinsiyeti", Value: stGender},
			{Label: "Okulu", Value: school.SchoolName},
			{Label: "Okul Numarası", Value: stNo},
			{Label: "Sınıfı", Value: stClass},
			{Label: "Doğum Tarihi", Value: stBirthDate},
			{Label: "Yarı Yıl", Value: school.AcademicYear},
			{Label: "Başarı durumu", Value: academicStanding},
			{Label: "Sağlık durumu", Value: healthStatus},
			{IsParentTable: true},
			{Label: "Ailesinin ekonomik durumu", Value: economicStatus},
			{Label: "Kardeş sayısı ve yaşları", Value: siblingCount},
			{Label: "Ailesinin oturduğu yer ve açık adresi", Value: addrVal, MinHeight: 20},
			{Label: "Şimdiye kadar aldığı yaptırımlar ve genel durumu", Value: priorPenalty},
			{Label: "Yaptırımı gerektiren davranışının yapıldığı yer ve tarihi", Value: offensePlaceTime},
			{Label: "Yaptırımı gerektiren davranışının çeşidi", Value: offenseType},
			{Label: "Yaptırımı gerektiren davranışının nedeni", Value: offenseReason},
			{Label: "     Olayla ilgili olarak;", IsSection: true},
			{Label: "     a) Yaptırım uygulanacak öğrencinin ifadesinin özeti", Value: defenseText, MinHeight: 20},
			{Label: "     b) Tanıkların ifadesinin özeti", Value: witnessText, MinHeight: 20},
			{Label: "     c) Varsa yaptırım gerektiren davranışının tespitine yarayan diğer deliller", Value: evidenceText, MinHeight: 20},
			{Label: "Yaptırımı hafifleten veya şiddetlendiren nedenler", Value: mitigatingText, MinHeight: 20},
			{Label: "Verilen yaptırımın çeşidi ve dayandığı yönetmelik maddesi", Value: penaltyAndClause, MinHeight: 26},
			{Label: "Öğrenci Davranışlarını Değerlendirme Kurulunun kanaati", Value: boardOpinion, MinHeight: 18},
			{Label: "ÖĞRENCİ DAVRANIŞLARINI DEĞERLENDİRME KURULU KARARI", Value: decisionText, BoldValue: true, MinHeight: 45},
		}
	}

	fontSize := 7.5
	lineH := 9.5
	defaultMinH := 14.5
	if isOrta {
		fontSize = 7.0
		lineH = 8.8
		defaultMinH = 13.0
	}

	tableStartY := doc.PDF.GetY()
	leftX := doc.MarginLeft
	rightX := doc.MarginLeft + doc.UsableWidth

	// Draw thick top border of the table
	doc.PDF.SetLineWidth(1.8)
	doc.PDF.Line(leftX, tableStartY, rightX, tableStartY)

	for idx, r := range rows {
		curY := doc.PDF.GetY()

		if r.IsParentTable {
			// Ortaokul ANNE / BABA side-by-side subtable
			subRowH := 11.5
			totalSubH := subRowH * 7.0
			if curY+totalSubH > PageHeightPortrait-doc.MarginBottom-20 {
				doc.PDF.SetLineWidth(1.8)
				doc.PDF.Line(leftX, curY, rightX, curY)
				doc.PDF.AddPage()
				curY = doc.MarginTop
				doc.PDF.SetY(curY)
				doc.PDF.SetLineWidth(1.8)
				doc.PDF.Line(leftX, curY, rightX, curY)
			}

			midX := leftX + doc.UsableWidth/2.0
			wSubLabel := (doc.UsableWidth / 2.0) * 5.0 / 16.0
			wSubColon := (doc.UsableWidth / 2.0) * 1.0 / 16.0
			wSubVal := (doc.UsableWidth / 2.0) - wSubLabel - wSubColon

			mAge, mEdu, mAlive, mBio, mTog := "", "", "", "", ""
			fAge, fEdu, fAlive, fBio, fTog := "", "", "", "", ""
			if incStudent != nil {
				mAge = incStudent.MotherAge
				mEdu = incStudent.MotherEducation
				mAlive = incStudent.MotherAlive
				mBio = incStudent.MotherBio
				mTog = incStudent.MotherTogether
				fAge = incStudent.FatherAge
				fEdu = incStudent.FatherEducation
				fAlive = incStudent.FatherAlive
				fBio = incStudent.FatherBio
				fTog = incStudent.FatherTogether
			}

			parentRows := []struct {
				Label string
				MVal  string
				FVal  string
			}{
				{"Adı", stMother, stFather},
				{"Yaşı", mAge, fAge},
				{"Eğitimi", mEdu, fEdu},
				{"Hayatta mı", mAlive, fAlive},
				{"Öz mü", mBio, fBio},
				{"Birlikte mi", mTog, fTog},
			}

			// Header row: ANNE | BABA
			doc.SetBoldFont(fontSize)
			doc.PDF.SetX(leftX)
			doc.PDF.SetY(curY + 1.5)
			_ = doc.PDF.CellWithOption(&gopdf.Rect{W: doc.UsableWidth / 2.0, H: subRowH}, "ANNE", gopdf.CellOption{Align: gopdf.Center})
			doc.PDF.SetX(midX)
			doc.PDF.SetY(curY + 1.5)
			_ = doc.PDF.CellWithOption(&gopdf.Rect{W: doc.UsableWidth / 2.0, H: subRowH}, "BABA", gopdf.CellOption{Align: gopdf.Center})

			doc.PDF.SetLineWidth(0.3)
			doc.PDF.Line(leftX, curY+subRowH, rightX, curY+subRowH)

			for pi, pr := range parentRows {
				py := curY + float64(pi+1)*subRowH
				doc.SetBoldFont(fontSize)
				doc.PDF.SetX(leftX + 4)
				doc.PDF.SetY(py + 1.5)
				_ = doc.PDF.CellWithOption(&gopdf.Rect{W: wSubLabel - 4, H: subRowH}, pr.Label, gopdf.CellOption{Align: gopdf.Left})
				doc.PDF.SetX(leftX + wSubLabel)
				doc.PDF.SetY(py + 1.5)
				_ = doc.PDF.CellWithOption(&gopdf.Rect{W: wSubColon, H: subRowH}, ":", gopdf.CellOption{Align: gopdf.Center})

				doc.SetRegularFont(fontSize)
				doc.PDF.SetX(leftX + wSubLabel + wSubColon + 2)
				doc.PDF.SetY(py + 1.5)
				_ = doc.PDF.CellWithOption(&gopdf.Rect{W: wSubVal - 4, H: subRowH}, pr.MVal, gopdf.CellOption{Align: gopdf.Left})

				doc.SetBoldFont(fontSize)
				doc.PDF.SetX(midX + 4)
				doc.PDF.SetY(py + 1.5)
				_ = doc.PDF.CellWithOption(&gopdf.Rect{W: wSubLabel - 4, H: subRowH}, pr.Label, gopdf.CellOption{Align: gopdf.Left})
				doc.PDF.SetX(midX + wSubLabel)
				doc.PDF.SetY(py + 1.5)
				_ = doc.PDF.CellWithOption(&gopdf.Rect{W: wSubColon, H: subRowH}, ":", gopdf.CellOption{Align: gopdf.Center})

				doc.SetRegularFont(fontSize)
				doc.PDF.SetX(midX + wSubLabel + wSubColon + 2)
				doc.PDF.SetY(py + 1.5)
				_ = doc.PDF.CellWithOption(&gopdf.Rect{W: wSubVal - 4, H: subRowH}, pr.FVal, gopdf.CellOption{Align: gopdf.Left})

				doc.PDF.SetLineWidth(0.3)
				doc.PDF.Line(leftX, py+subRowH, rightX, py+subRowH)
			}

			// Center divider between ANNE and BABA
			doc.PDF.SetLineWidth(0.5)
			doc.PDF.Line(midX, curY, midX, curY+totalSubH)

			// Outer left/right thick borders for this block
			doc.PDF.SetLineWidth(1.8)
			doc.PDF.Line(leftX, curY, leftX, curY+totalSubH)
			doc.PDF.Line(rightX, curY, rightX, curY+totalSubH)

			doc.PDF.SetY(curY + totalSubH)
			continue
		}

		if r.IsSection {
			cellH := defaultMinH
			if curY+cellH > PageHeightPortrait-doc.MarginBottom-20 {
				doc.PDF.SetLineWidth(1.8)
				doc.PDF.Line(leftX, curY, rightX, curY)
				doc.PDF.AddPage()
				curY = doc.MarginTop
				doc.PDF.SetY(curY)
				doc.PDF.SetLineWidth(1.8)
				doc.PDF.Line(leftX, curY, rightX, curY)
			}

			doc.SetBoldFont(fontSize)
			doc.PDF.SetX(leftX + 3)
			doc.PDF.SetY(curY + 2.5)
			_ = doc.PDF.CellWithOption(&gopdf.Rect{W: doc.UsableWidth - 6, H: lineH}, r.Label, gopdf.CellOption{Align: gopdf.Left})

			// Outer left & right thick borders
			doc.PDF.SetLineWidth(1.8)
			doc.PDF.Line(leftX, curY, leftX, curY+cellH)
			doc.PDF.Line(rightX, curY, rightX, curY+cellH)

			// Thin bottom row divider
			doc.PDF.SetLineWidth(0.35)
			doc.PDF.Line(leftX, curY+cellH, rightX, curY+cellH)

			doc.PDF.SetY(curY + cellH)
			continue
		}

		val := strings.TrimSpace(r.Value)

		// Measure label lines
		doc.SetBoldFont(fontSize)
		lblLines := doc.SplitLines(r.Label, wCol1-6)

		// Measure value lines (keep blank instead of em-dash if empty, matching authentic form)
		var valLines []string
		if val != "" {
			if r.BoldValue {
				doc.SetBoldFont(fontSize)
			} else {
				doc.SetRegularFont(fontSize)
			}
			valLines = doc.SplitLines(val, wCol3-8)
		}

		maxLines := len(lblLines)
		if len(valLines) > maxLines {
			maxLines = len(valLines)
		}
		if maxLines < 1 {
			maxLines = 1
		}

		cellH := float64(maxLines)*lineH + 4.5
		minH := defaultMinH
		if r.MinHeight > minH {
			minH = r.MinHeight
		}
		if cellH < minH {
			cellH = minH
		}

		// Page break check if table overflows page
		if curY+cellH > PageHeightPortrait-doc.MarginBottom-20 {
			doc.PDF.SetLineWidth(1.8)
			doc.PDF.Line(leftX, curY, rightX, curY)
			doc.PDF.AddPage()
			curY = doc.MarginTop
			doc.PDF.SetY(curY)
			doc.PDF.SetLineWidth(1.8)
			doc.PDF.Line(leftX, curY, rightX, curY)
		}

		// Col 1: Label (Bold)
		doc.SetBoldFont(fontSize)
		doc.DrawWrappedCellLines(leftX+3, curY+2.5, wCol1-6, lineH, lblLines, gopdf.Left)

		// Col 2: ":" (Bold, centered)
		doc.PDF.SetX(leftX + wCol1)
		doc.PDF.SetY(curY + 2.5)
		_ = doc.PDF.CellWithOption(&gopdf.Rect{W: wCol2, H: lineH}, ":", gopdf.CellOption{Align: gopdf.Center})

		// Col 3: Value
		if len(valLines) > 0 {
			if r.BoldValue {
				doc.SetBoldFont(fontSize)
			} else {
				doc.SetRegularFont(fontSize)
			}
			doc.DrawWrappedCellLines(leftX+wCol1+wCol2+3, curY+2.5, wCol3-6, lineH, valLines, gopdf.Left)
		}

		// Outer left & right thick borders
		doc.PDF.SetLineWidth(1.8)
		doc.PDF.Line(leftX, curY, leftX, curY+cellH)
		doc.PDF.Line(rightX, curY, rightX, curY+cellH)

		// Bottom border: thick on final row, thin on internal rows
		if idx == len(rows)-1 {
			doc.PDF.SetLineWidth(1.8)
		} else {
			doc.PDF.SetLineWidth(0.35)
		}
		doc.PDF.Line(leftX, curY+cellH, rightX, curY+cellH)

		doc.PDF.SetY(curY + cellH)
	}

	// Substitute Member replacement check
	effectiveChair := school.BoardChair
	effectiveChairTitle := school.BoardChairTitle
	effectiveM1 := school.Member1
	effectiveM2 := school.Member2
	effectiveM3 := school.Member3
	effectiveM4 := school.Member4

	if verdict != nil {
		if verdict.ChairIsSubstitute && school.SubBoardChair != "" {
			effectiveChair = school.SubBoardChair
			effectiveChairTitle = school.SubBoardChairTitle
		}
		if verdict.M1IsSubstitute && school.SubMember1 != "" {
			effectiveM1 = school.SubMember1
		}
		if verdict.M2IsSubstitute && school.SubMember2 != "" {
			effectiveM2 = school.SubMember2
		}
		if verdict.M3IsSubstitute && school.SubMember3 != "" {
			effectiveM3 = school.SubMember3
		}
		if verdict.M4IsSubstitute && school.SubMember4 != "" {
			effectiveM4 = school.SubMember4
		}
	}

	// Check if signatures fit on current page; otherwise add new page
	curY := doc.PDF.GetY()
	if curY+70 > PageHeightPortrait-doc.MarginBottom {
		doc.PDF.AddPage()
		curY = doc.MarginTop + 10
		doc.PDF.SetY(curY)
	} else {
		curY += 24.0 // Clear white space for physical signatures/stamps
		doc.PDF.SetY(curY)
	}

	// 4. Authentic 5-Column Signatures Block matching official scan
	colW := doc.UsableWidth / 5.0
	chairRoleLine := "Ödül ve Disiplin Kurulu"
	if isOrta {
		chairRoleLine = "Öğrenci Davranışlarını\nDeğerlendirme Kurulu"
	}
	if strings.TrimSpace(effectiveChairTitle) == "" {
		effectiveChairTitle = "Müdür Yardımcısı"
	}

	sigCols := []struct {
		Name  string
		Lines []string
	}{
		{
			Name:  strings.TrimSpace(effectiveChair),
			Lines: append(strings.Split(chairRoleLine, "\n"), "Başkanı", strings.TrimSpace(effectiveChairTitle)),
		},
		{
			Name:  strings.TrimSpace(effectiveM1),
			Lines: []string{"Üye"},
		},
		{
			Name:  strings.TrimSpace(effectiveM2),
			Lines: []string{"Üye"},
		},
		{
			Name:  strings.TrimSpace(effectiveM3),
			Lines: []string{"Üye"},
		},
		{
			Name:  strings.TrimSpace(effectiveM4),
			Lines: []string{"Üye"},
		},
	}

	sigStartY := doc.PDF.GetY()

	for i, sc := range sigCols {
		colX := leftX + float64(i)*colW + 2
		yPos := sigStartY

		if sc.Name != "" {
			doc.SetBoldFont(7.5)
			nameLines := doc.SplitLines(sc.Name, colW-4)
			doc.DrawWrappedCellLines(colX, yPos, colW-4, 9.0, nameLines, gopdf.Center)
			yPos += float64(len(nameLines)) * 9.0
		} else {
			yPos += 9.0
		}

		doc.SetRegularFont(7.0)
		for _, rl := range sc.Lines {
			rl = strings.TrimSpace(rl)
			if rl == "" {
				continue
			}
			rLines := doc.SplitLines(rl, colW-4)
			doc.DrawWrappedCellLines(colX, yPos, colW-4, 8.5, rLines, gopdf.Center)
			yPos += float64(len(rLines)) * 8.5
		}
	}

	return doc.PDF.WritePdf(outputPath)
}
