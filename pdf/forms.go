package pdf

import (
	"fmt"
	"strings"
	"time"

	"github.com/signintech/gopdf"

	"turandisiplin/model"
)

// GenerateIfadeSavunmaTutanagi generates either İfade Tutanağı or Savunma Tutanağı.
func GenerateIfadeSavunmaTutanagi(school *model.SchoolInfo, student *model.Student, incident *model.Incident, isSavunma bool, statementText string, outputPath string) error {
	doc, err := NewDoc(false)
	if err != nil {
		return err
	}

	doc.PDF.AddPage()

	// Title
	title := "İ F A D E   T U T A N A Ğ I"
	subTitle := "İFADESİ ALINANIN"
	sigTitle := "İFADEYİ VERENİN ADI SOYADI - İMZASI"
	if isSavunma {
		title = "S A V U N M A   T U T A N A Ğ I"
		subTitle = "SAVUNMASI ALINANIN"
		sigTitle = "SAVUNMAYI VERENİN ADI SOYADI - İMZASI"
	}

	doc.DrawTitle(title, 14, 10)

	// Subtitle
	doc.SetBoldFont(9)
	doc.PDF.SetX(doc.MarginLeft)
	_ = doc.PDF.CellWithOption(&gopdf.Rect{W: doc.UsableWidth, H: 14}, subTitle, gopdf.CellOption{Align: gopdf.Left})
	doc.PDF.Br(16)

	// Student Info Table (2-column layout)
	lblW := 140.0
	valW := doc.UsableWidth - lblW

	infoRows := [][2]string{
		{"Adı Soyadı", ": " + student.FullName},
		{"Sınıfı ve Okul No", fmt.Sprintf(": %s  -  %s", student.ClassSection, student.SchoolNo)},
		{"Baba - Anne Adı", fmt.Sprintf(": %s - %s", student.FatherName, student.MotherName)},
		{"Doğum Yeri ve Tarihi", fmt.Sprintf(": %s - %s", student.BirthPlace, student.BirthDate)},
		{choose(isSavunma, "Savunma Tarihi", "İfade Tarihi"), ": " + time.Now().Format("02/01/2006")},
	}

	for _, row := range infoRows {
		doc.SetBoldFont(8.5)
		doc.PDF.SetX(doc.MarginLeft)
		_ = doc.PDF.CellWithOption(&gopdf.Rect{W: lblW, H: 13}, row[0], gopdf.CellOption{Align: gopdf.Left})

		doc.SetRegularFont(8.5)
		doc.PDF.SetX(doc.MarginLeft + lblW)
		_ = doc.PDF.CellWithOption(&gopdf.Rect{W: valW, H: 13}, row[1], gopdf.CellOption{Align: gopdf.Left})
		doc.PDF.Br(14)
	}

	doc.PDF.Br(6)

	// Incident prompt
	doc.PDF.SetLineWidth(0.5)
	yLineTop := doc.PDF.GetY()
	doc.PDF.Line(doc.MarginLeft, yLineTop, doc.MarginLeft+doc.UsableWidth, yLineTop)
	doc.PDF.Br(4)

	doc.SetBoldFont(8.5)
	doc.PDF.SetX(doc.MarginLeft)
	_ = doc.PDF.CellWithOption(&gopdf.Rect{W: doc.UsableWidth, H: 12}, "OLAY:", gopdf.CellOption{Align: gopdf.Left})
	doc.PDF.Br(12)

	doc.SetRegularFont(8.0)
	doc.PDF.SetX(doc.MarginLeft)
	eventSummary := ""
	if incident != nil {
		eventSummary = incident.Summary
		if eventSummary == "" {
			eventSummary = incident.Title
		}
	}
	_ = doc.PDF.MultiCell(&gopdf.Rect{W: doc.UsableWidth, H: 12}, eventSummary)
	doc.PDF.Br(4)

	yLineBot := doc.PDF.GetY()
	doc.PDF.Line(doc.MarginLeft, yLineBot, doc.MarginLeft+doc.UsableWidth, yLineBot)
	doc.PDF.Br(10)

	// If statement text is provided, render it, otherwise render ruled lines
	if strings.TrimSpace(statementText) != "" {
		doc.SetRegularFont(8.5)
		doc.PDF.SetX(doc.MarginLeft)
		_ = doc.PDF.MultiCell(&gopdf.Rect{W: doc.UsableWidth, H: 14}, statementText)
		doc.PDF.Br(15)
	} else {
		// Calculate remaining space for ruled lines (up to ~22 lines so signatures fit comfortably)
		doc.DrawRuledLines(20, 15)
	}

	// Student Signature
	doc.PDF.Br(15)
	doc.SetBoldFont(8.0)
	doc.PDF.SetX(doc.MarginLeft)
	_ = doc.PDF.CellWithOption(&gopdf.Rect{W: doc.UsableWidth, H: 12}, sigTitle, gopdf.CellOption{Align: gopdf.Right})
	doc.PDF.Br(12)
	doc.SetRegularFont(7.5)
	doc.PDF.SetX(doc.MarginLeft)
	_ = doc.PDF.CellWithOption(&gopdf.Rect{W: doc.UsableWidth, H: 12}, "İmza: ........................................", gopdf.CellOption{Align: gopdf.Right})
	doc.PDF.Br(15)

	// Board Members Signature (5 columns)
	boardLabel := "Ödül ve Disiplin"
	if strings.ToUpper(school.SchoolType) == "ORTAOKUL" {
		boardLabel = "Öğrenci Davranışlarını Değerlendirme"
	}
	doc.DrawSignatures5(school, boardLabel)

	return doc.PDF.WritePdf(outputPath)
}

// GenerateBilgiToplamaFormu generates information collection forms for incident, teachers, counselor, parents, or student.
func GenerateBilgiToplamaFormu(formType string, school *model.SchoolInfo, student *model.Student, incident *model.Incident, outputPath string) error {
	doc, err := NewDoc(false)
	if err != nil {
		return err
	}

	doc.PDF.AddPage()

	boardName := "ÖDÜL VE DİSİPLİN"
	if strings.ToUpper(school.SchoolType) == "ORTAOKUL" {
		boardName = "ÖĞRENCİ DAVRANIŞLARINI DEĞERLENDİRME"
	}

	switch formType {
	case "Olay":
		doc.DrawTitle(fmt.Sprintf("%s KURULU BİLGİ TOPLAMA FORMU - Olay", boardName), 12, 10)
		doc.SetBoldFont(8.5)
		doc.PDF.SetX(doc.MarginLeft)
		_ = doc.PDF.CellWithOption(&gopdf.Rect{W: doc.UsableWidth, H: 13}, "Olayla İlgili Bilgilerim:", gopdf.CellOption{Align: gopdf.Left})
		doc.PDF.Br(14)
		doc.DrawRuledLines(14, 15)
		doc.PDF.Br(12)

		doc.SetBoldFont(8.5)
		doc.PDF.SetX(doc.MarginLeft)
		_ = doc.PDF.CellWithOption(&gopdf.Rect{W: doc.UsableWidth, H: 13}, "Arkadaşım Hakkındaki Görüşlerim:", gopdf.CellOption{Align: gopdf.Left})
		doc.PDF.Br(14)
		doc.DrawRuledLines(14, 15)

		doc.PDF.Br(16)
		doc.SetRegularFont(8.0)
		doc.PDF.SetX(doc.MarginLeft)
		_ = doc.PDF.CellWithOption(&gopdf.Rect{W: doc.UsableWidth, H: 12}, "Öğrencinin Adı Soyadı / Sınıfı - Numarası / İmzası", gopdf.CellOption{Align: gopdf.Right})
		doc.PDF.Br(12)
		_ = doc.PDF.CellWithOption(&gopdf.Rect{W: doc.UsableWidth, H: 12}, "Tarih: ......./......./20.......   İmza: ........................................", gopdf.CellOption{Align: gopdf.Right})

	case "Öğretmen":
		doc.DrawTitle(fmt.Sprintf("%s KURULU BİLGİ TOPLAMA FORMU - Öğretmen", boardName), 12, 10)
		renderStudentMiniHeader(doc, student)

		sections := []string{"1. Dersimdeki Tutumu:", "2. Ders Dışı Faaliyetleri:", "3. Arkadaşları İle İlişkisi:"}
		for _, sec := range sections {
			doc.SetBoldFont(8.5)
			doc.PDF.SetX(doc.MarginLeft)
			_ = doc.PDF.CellWithOption(&gopdf.Rect{W: doc.UsableWidth, H: 13}, sec, gopdf.CellOption{Align: gopdf.Left})
			doc.PDF.Br(13)
			doc.DrawRuledLines(9, 15)
			doc.PDF.Br(8)
		}

		doc.PDF.Br(10)
		doc.SetRegularFont(8.0)
		doc.PDF.SetX(doc.MarginLeft)
		_ = doc.PDF.CellWithOption(&gopdf.Rect{W: doc.UsableWidth, H: 12}, "Öğretmenin Adı Soyadı / Branşı / İmzası", gopdf.CellOption{Align: gopdf.Right})
		doc.PDF.Br(12)
		_ = doc.PDF.CellWithOption(&gopdf.Rect{W: doc.UsableWidth, H: 12}, "Tarih: ......./......./20.......   İmza: ........................................", gopdf.CellOption{Align: gopdf.Right})

	case "Sınıf Rehber Öğretmeni":
		doc.DrawTitle(fmt.Sprintf("%s KURULU BİLGİ TOPLAMA FORMU - Sınıf Rehber Öğretmeni", boardName), 12, 10)
		renderStudentMiniHeader(doc, student)

		sections := []string{"1. Kişisel Özellikleri:", "2. Sosyal Özellikleri:"}
		for _, sec := range sections {
			doc.SetBoldFont(8.5)
			doc.PDF.SetX(doc.MarginLeft)
			_ = doc.PDF.CellWithOption(&gopdf.Rect{W: doc.UsableWidth, H: 13}, sec, gopdf.CellOption{Align: gopdf.Left})
			doc.PDF.Br(13)
			doc.DrawRuledLines(14, 15)
			doc.PDF.Br(10)
		}

		doc.PDF.Br(12)
		doc.SetRegularFont(8.0)
		doc.PDF.SetX(doc.MarginLeft)
		_ = doc.PDF.CellWithOption(&gopdf.Rect{W: doc.UsableWidth, H: 12}, "Sınıf Rehber Öğretmeni Adı Soyadı / İmzası", gopdf.CellOption{Align: gopdf.Right})
		doc.PDF.Br(12)
		_ = doc.PDF.CellWithOption(&gopdf.Rect{W: doc.UsableWidth, H: 12}, "Tarih: ......./......./20.......   İmza: ........................................", gopdf.CellOption{Align: gopdf.Right})

	case "Veli":
		doc.DrawTitle(fmt.Sprintf("%s KURULU VELİ BİLGİ TOPLAMA FORMU", boardName), 12, 10)
		renderStudentMiniHeader(doc, student)

		doc.SetBoldFont(8.5)
		doc.PDF.SetX(doc.MarginLeft)
		_ = doc.PDF.CellWithOption(&gopdf.Rect{W: doc.UsableWidth, H: 13}, "Öğrenci Hakkındaki Görüşlerim:", gopdf.CellOption{Align: gopdf.Left})
		doc.PDF.Br(13)
		doc.DrawRuledLines(26, 15)

		doc.PDF.Br(16)
		doc.SetRegularFont(8.0)
		doc.PDF.SetX(doc.MarginLeft)
		_ = doc.PDF.CellWithOption(&gopdf.Rect{W: doc.UsableWidth, H: 12}, "Velinin Adı Soyadı / Yakınlığı / İmzası", gopdf.CellOption{Align: gopdf.Right})
		doc.PDF.Br(12)
		_ = doc.PDF.CellWithOption(&gopdf.Rect{W: doc.UsableWidth, H: 12}, "Tarih: ......./......./20.......   İmza: ........................................", gopdf.CellOption{Align: gopdf.Right})

	default: // Öğrenci
		doc.DrawTitle(fmt.Sprintf("%s KURULU BİLGİ TOPLAMA FORMU - Öğrenci", boardName), 12, 10)
		renderStudentMiniHeader(doc, student)

		sections := []string{"1. Dersteki Davranışları:", "2. Teneffüslerdeki Davranışları:", "3. Okul Dışındaki Davranışları:"}
		for _, sec := range sections {
			doc.SetBoldFont(8.5)
			doc.PDF.SetX(doc.MarginLeft)
			_ = doc.PDF.CellWithOption(&gopdf.Rect{W: doc.UsableWidth, H: 13}, sec, gopdf.CellOption{Align: gopdf.Left})
			doc.PDF.Br(13)
			doc.DrawRuledLines(9, 15)
			doc.PDF.Br(8)
		}

		doc.PDF.Br(10)
		doc.SetRegularFont(8.0)
		doc.PDF.SetX(doc.MarginLeft)
		_ = doc.PDF.CellWithOption(&gopdf.Rect{W: doc.UsableWidth, H: 12}, "Öğrencinin Adı Soyadı / İmzası", gopdf.CellOption{Align: gopdf.Right})
		doc.PDF.Br(12)
		_ = doc.PDF.CellWithOption(&gopdf.Rect{W: doc.UsableWidth, H: 12}, "Tarih: ......./......./20.......   İmza: ........................................", gopdf.CellOption{Align: gopdf.Right})
	}

	return doc.PDF.WritePdf(outputPath)
}

func renderStudentMiniHeader(doc *DocWrapper, student *model.Student) {
	if student == nil {
		return
	}
	doc.SetBoldFont(8.5)
	doc.PDF.SetX(doc.MarginLeft)
	info := fmt.Sprintf("Öğrencinin Adı Soyadı: %s      Sınıfı / Okul No: %s - %s", student.FullName, student.ClassSection, student.SchoolNo)
	_ = doc.PDF.CellWithOption(&gopdf.Rect{W: doc.UsableWidth, H: 14}, info, gopdf.CellOption{Align: gopdf.Left})
	doc.PDF.Br(16)
}

// GenerateEkSureTalepFormu generates the 7-day investigation extension request.
func GenerateEkSureTalepFormu(school *model.SchoolInfo, student *model.Student, incident *model.Incident, extraDays int, reason string, outputPath string) error {
	doc, err := NewDoc(false)
	if err != nil {
		return err
	}

	doc.PDF.AddPage()

	boardName := "ÖDÜL VE DİSİPLİN"
	if strings.ToUpper(school.SchoolType) == "ORTAOKUL" {
		boardName = "ÖĞRENCİ DAVRANIŞLARINI DEĞERLENDİRME"
	}

	doc.DrawTitle(fmt.Sprintf("%s KURULU\nOLAY EK SÜRE TALEP FORMU", boardName), 12, 12)

	// 2-Column Bordered Table
	lblW := 140.0
	valW := doc.UsableWidth - lblW

	stName := ""
	if student != nil {
		stName = fmt.Sprintf("%s - %s - %s", student.ClassSection, student.SchoolNo, student.FullName)
	}
	incName := ""
	incSum := ""
	if incident != nil {
		incName = incident.Title
		incSum = incident.Summary
	}

	rows := [][2]string{
		{"Olay", incName},
		{"Olay Açıklaması", incSum},
		{"Öğrenci Adı Soyadı", stName},
		{"Ek Süre", fmt.Sprintf("%d GÜN", extraDays)},
		{"Ek Süre Nedeni", reason},
		{"Tarih", time.Now().Format("02/01/2006")},
	}

	for _, r := range rows {
		h := 22.0
		if len(r[1]) > 50 {
			h = 32.0
		}
		y := doc.PDF.GetY()

		doc.DrawRect(doc.MarginLeft, y, lblW, h)
		doc.DrawRect(doc.MarginLeft+lblW, y, valW, h)

		doc.SetBoldFont(8.0)
		doc.PDF.SetX(doc.MarginLeft + 4)
		doc.PDF.SetY(y + 4)
		_ = doc.PDF.CellWithOption(&gopdf.Rect{W: lblW - 8, H: 12}, r[0], gopdf.CellOption{Align: gopdf.Left})

		doc.SetRegularFont(8.0)
		doc.PDF.SetX(doc.MarginLeft + lblW + 4)
		doc.PDF.SetY(y + 4)
		_ = doc.PDF.MultiCell(&gopdf.Rect{W: valW - 8, H: 12}, r[1])

		doc.PDF.SetY(y + h)
	}

	doc.PDF.Br(25)

	// Petition block
	doc.SetBoldFont(10)
	doc.PDF.SetX(doc.MarginLeft)
	_ = doc.PDF.CellWithOption(&gopdf.Rect{W: doc.UsableWidth, H: 14}, strings.ToUpper(school.SchoolName)+" MÜDÜRLÜĞÜNE", gopdf.CellOption{Align: gopdf.Center})
	doc.PDF.Br(20)

	doc.SetRegularFont(8.5)
	doc.PDF.SetX(doc.MarginLeft)
	petitionText := "Yukarıda belirtilen olay ve öğrenci ile ilgili yürütülen disiplin / davranış değerlendirme incelemesinin sağlıklı bir şekilde tamamlanabilmesi için ilgili mevzuat gereğince belirtilen sürenin ek süre olarak kurulumuza verilmesini arz ederim."
	_ = doc.PDF.MultiCell(&gopdf.Rect{W: doc.UsableWidth, H: 14}, petitionText)
	doc.PDF.Br(25)

	// Board Chair Signature
	doc.SetBoldFont(8.0)
	doc.PDF.SetX(doc.MarginLeft)
	_ = doc.PDF.CellWithOption(&gopdf.Rect{W: doc.UsableWidth, H: 12}, school.BoardChair, gopdf.CellOption{Align: gopdf.Right})
	doc.PDF.Br(12)
	doc.SetRegularFont(7.5)
	doc.PDF.SetX(doc.MarginLeft)
	_ = doc.PDF.CellWithOption(&gopdf.Rect{W: doc.UsableWidth, H: 12}, boardName+" Kurulu Başkanı", gopdf.CellOption{Align: gopdf.Right})
	doc.PDF.Br(30)

	// Approval Block
	doc.DrawPrincipalApproval(school.PrincipalName)

	return doc.PDF.WritePdf(outputPath)
}

// GenerateSikayetFormu generates the complaint form.
func GenerateSikayetFormu(school *model.SchoolInfo, complainant, reported, witnesses, evidences, incidentDate, topic, description, guidanceNotes string, outputPath string) error {
	doc, err := NewDoc(false)
	if err != nil {
		return err
	}

	doc.PDF.AddPage()
	doc.DrawAntet(school)
	doc.DrawTitle("Ş İ K A Y E T   F O R M U", 13, 10)

	lblW := 150.0
	valW := doc.UsableWidth - lblW

	rows := [][2]string{
		{"Şikayet Eden (Adı-Soyadı-Sınıfı-No)", complainant},
		{"Şikayet Edilen (Adı-Soyadı-Sınıfı-No)", reported},
		{"Tanıklar", witnesses},
		{"Deliller", evidences},
		{"Olay Tarihi ve Saati", incidentDate},
		{"Olayın Konusu", topic},
	}

	for _, r := range rows {
		h := 20.0
		y := doc.PDF.GetY()
		doc.DrawRect(doc.MarginLeft, y, lblW, h)
		doc.DrawRect(doc.MarginLeft+lblW, y, valW, h)

		doc.SetBoldFont(7.5)
		doc.PDF.SetX(doc.MarginLeft + 4)
		doc.PDF.SetY(y + 4)
		_ = doc.PDF.CellWithOption(&gopdf.Rect{W: lblW - 8, H: 11}, r[0], gopdf.CellOption{Align: gopdf.Left})

		doc.SetRegularFont(7.5)
		doc.PDF.SetX(doc.MarginLeft + lblW + 4)
		doc.PDF.SetY(y + 4)
		_ = doc.PDF.CellWithOption(&gopdf.Rect{W: valW - 8, H: 11}, r[1], gopdf.CellOption{Align: gopdf.Left})

		doc.PDF.SetY(y + h)
	}

	doc.PDF.Br(10)
	doc.SetBoldFont(8.0)
	doc.PDF.SetX(doc.MarginLeft)
	_ = doc.PDF.CellWithOption(&gopdf.Rect{W: doc.UsableWidth, H: 12}, "Olayın Açıklaması:", gopdf.CellOption{Align: gopdf.Left})
	doc.PDF.Br(14)

	if strings.TrimSpace(description) != "" {
		doc.SetRegularFont(8.0)
		doc.PDF.SetX(doc.MarginLeft)
		_ = doc.PDF.MultiCell(&gopdf.Rect{W: doc.UsableWidth, H: 13}, description)
		doc.PDF.Br(10)
	} else {
		doc.DrawDottedLines(5, 14)
	}

	doc.PDF.Br(10)
	doc.SetBoldFont(8.0)
	doc.PDF.SetX(doc.MarginLeft)
	_ = doc.PDF.CellWithOption(&gopdf.Rect{W: doc.UsableWidth, H: 12}, "Rehberlik ve Psikolojik Danışma Servisi Görüş ve Önerisi:", gopdf.CellOption{Align: gopdf.Left})
	doc.PDF.Br(14)

	if strings.TrimSpace(guidanceNotes) != "" {
		doc.SetRegularFont(8.0)
		doc.PDF.SetX(doc.MarginLeft)
		_ = doc.PDF.MultiCell(&gopdf.Rect{W: doc.UsableWidth, H: 13}, guidanceNotes)
		doc.PDF.Br(10)
	} else {
		doc.DrawDottedLines(6, 14)
	}

	doc.PDF.Br(15)
	doc.SetRegularFont(7.5)
	doc.PDF.SetX(doc.MarginLeft)
	_ = doc.PDF.CellWithOption(&gopdf.Rect{W: doc.UsableWidth, H: 11}, "Psikolojik Danışman / Rehber Öğretmen: ........................................     İmza: .........................", gopdf.CellOption{Align: gopdf.Right})

	doc.DrawPrincipalApproval(school.PrincipalName)

	return doc.PDF.WritePdf(outputPath)
}

// GenerateSozluUyariTutanagi generates the verbal warning record.
func GenerateSozluUyariTutanagi(school *model.SchoolInfo, student *model.Student, dateStr, topic, details, warningTeacher string, outputPath string) error {
	doc, err := NewDoc(false)
	if err != nil {
		return err
	}

	doc.PDF.AddPage()
	doc.DrawAntet(school)
	doc.DrawTitle("S Ö Z L Ü   U Y A R I   T U T A N A Ğ I", 13, 15)

	lblW := 140.0
	valW := doc.UsableWidth - lblW

	stName := ""
	stNo := ""
	if student != nil {
		stName = student.FullName
		stNo = student.ClassSection + " - " + student.SchoolNo
	}

	rows := [][2]string{
		{"Adı Soyadı", stName},
		{"Sınıfı / Numarası", stNo},
		{"Tarih", dateStr},
		{"Konu", topic},
	}

	for _, r := range rows {
		h := 20.0
		y := doc.PDF.GetY()
		doc.DrawRect(doc.MarginLeft, y, lblW, h)
		doc.DrawRect(doc.MarginLeft+lblW, y, valW, h)

		doc.SetBoldFont(8.0)
		doc.PDF.SetX(doc.MarginLeft + 4)
		doc.PDF.SetY(y + 4)
		_ = doc.PDF.CellWithOption(&gopdf.Rect{W: lblW - 8, H: 12}, r[0], gopdf.CellOption{Align: gopdf.Left})

		doc.SetRegularFont(8.0)
		doc.PDF.SetX(doc.MarginLeft + lblW + 4)
		doc.PDF.SetY(y + 4)
		_ = doc.PDF.CellWithOption(&gopdf.Rect{W: valW - 8, H: 12}, r[1], gopdf.CellOption{Align: gopdf.Left})

		doc.PDF.SetY(y + h)
	}

	doc.PDF.Br(15)
	doc.SetBoldFont(8.5)
	doc.PDF.SetX(doc.MarginLeft)
	_ = doc.PDF.CellWithOption(&gopdf.Rect{W: doc.UsableWidth, H: 13}, "Uyarı Açıklaması ve Görüşme Özeti:", gopdf.CellOption{Align: gopdf.Left})
	doc.PDF.Br(15)

	if strings.TrimSpace(details) != "" {
		doc.SetRegularFont(8.5)
		doc.PDF.SetX(doc.MarginLeft)
		_ = doc.PDF.MultiCell(&gopdf.Rect{W: doc.UsableWidth, H: 14}, details)
		doc.PDF.Br(15)
	} else {
		doc.DrawRuledLines(10, 16)
	}

	doc.PDF.Br(25)

	// Dual signature block
	colW := doc.UsableWidth / 2.0
	ySig := doc.PDF.GetY()

	doc.SetBoldFont(8.0)
	doc.PDF.SetX(doc.MarginLeft)
	doc.PDF.SetY(ySig)
	_ = doc.PDF.CellWithOption(&gopdf.Rect{W: colW, H: 12}, "UYARILAN ÖĞRENCİ", gopdf.CellOption{Align: gopdf.Center})

	doc.PDF.SetX(doc.MarginLeft + colW)
	doc.PDF.SetY(ySig)
	_ = doc.PDF.CellWithOption(&gopdf.Rect{W: colW, H: 12}, "UYARAN ÖĞRETMEN / YÖNETİCİ", gopdf.CellOption{Align: gopdf.Center})

	doc.SetRegularFont(7.5)
	doc.PDF.SetX(doc.MarginLeft)
	doc.PDF.SetY(ySig + 25)
	_ = doc.PDF.CellWithOption(&gopdf.Rect{W: colW, H: 11}, stName, gopdf.CellOption{Align: gopdf.Center})

	tName := warningTeacher
	if tName == "" {
		tName = "........................................"
	}
	doc.PDF.SetX(doc.MarginLeft + colW)
	doc.PDF.SetY(ySig + 25)
	_ = doc.PDF.CellWithOption(&gopdf.Rect{W: colW, H: 11}, tName, gopdf.CellOption{Align: gopdf.Center})

	doc.PDF.SetX(doc.MarginLeft)
	doc.PDF.SetY(ySig + 38)
	_ = doc.PDF.CellWithOption(&gopdf.Rect{W: colW, H: 11}, "İmza: ..............................", gopdf.CellOption{Align: gopdf.Center})

	doc.PDF.SetX(doc.MarginLeft + colW)
	doc.PDF.SetY(ySig + 38)
	_ = doc.PDF.CellWithOption(&gopdf.Rect{W: colW, H: 11}, "İmza: ..............................", gopdf.CellOption{Align: gopdf.Center})

	doc.PDF.SetY(ySig + 60)
	doc.DrawFooterAddress(school)

	return doc.PDF.WritePdf(outputPath)
}

// GenerateOgrenciSozlesmesi generates student contract.
func GenerateOgrenciSozlesmesi(school *model.SchoolInfo, student *model.Student, dateStr, outputPath string) error {
	doc, err := NewDoc(false)
	if err != nil {
		return err
	}

	doc.PDF.AddPage()
	doc.DrawAntet(school)
	doc.DrawTitle("Ö Ğ R E N C İ   S Ö Z L E Ş M E S İ", 14, 15)

	renderStudentMiniHeader(doc, student)

	pledge := "Okul kurallarına, öğretmenlerime ve arkadaşlarıma saygılı davranacağıma, derslerimi ve ödevlerimi aksatmadan takip edeceğime, okul araç ve gereçlerini koruyacağıma, huzur ve disiplini bozacak hiçbir tutum ve davranış içerisinde bulunmayacağıma; aksi takdirde ilgili mevzuat hükümlerinin uygulanacağını bildiğimi beyan ve taahhüt ederim."

	doc.SetRegularFont(9.0)
	doc.PDF.SetX(doc.MarginLeft)
	_ = doc.PDF.MultiCell(&gopdf.Rect{W: doc.UsableWidth, H: 16}, pledge)
	doc.PDF.Br(35)

	doc.SetBoldFont(8.5)
	doc.PDF.SetX(doc.MarginLeft)
	_ = doc.PDF.CellWithOption(&gopdf.Rect{W: doc.UsableWidth, H: 13}, "ÖĞRENCİ", gopdf.CellOption{Align: gopdf.Right})
	doc.PDF.Br(14)
	doc.SetRegularFont(8.0)
	doc.PDF.SetX(doc.MarginLeft)
	stName := ""
	if student != nil {
		stName = student.FullName
	}
	_ = doc.PDF.CellWithOption(&gopdf.Rect{W: doc.UsableWidth, H: 12}, stName, gopdf.CellOption{Align: gopdf.Right})
	doc.PDF.Br(12)
	_ = doc.PDF.CellWithOption(&gopdf.Rect{W: doc.UsableWidth, H: 12}, "Tarih: "+dateStr, gopdf.CellOption{Align: gopdf.Right})
	doc.PDF.Br(12)
	_ = doc.PDF.CellWithOption(&gopdf.Rect{W: doc.UsableWidth, H: 12}, "İmza: ........................................", gopdf.CellOption{Align: gopdf.Right})

	doc.DrawFooterAddress(school)
	return doc.PDF.WritePdf(outputPath)
}

// GenerateVeliGorusmeTutanagi generates parent-teacher interview record or blank form.
func GenerateVeliGorusmeTutanagi(school *model.SchoolInfo, student *model.Student, parentName, dateStr, topic, notes, decision, interviewer string, isBlank bool, outputPath string) error {
	doc, err := NewDoc(false)
	if err != nil {
		return err
	}

	doc.PDF.AddPage()
	doc.DrawAntet(school)
	formTitle := "V E L İ   G Ö R Ü Ş M E   T U T A N A Ğ I"
	if isBlank {
		formTitle = "V E L İ   G Ö R Ü Ş M E   F O R M U"
	}
	doc.DrawTitle(formTitle, 13, 12)

	lblW := 140.0
	valW := doc.UsableWidth - lblW

	stName := ""
	stClass := ""
	if student != nil {
		stName = student.FullName
		stClass = fmt.Sprintf("%s - %s", student.ClassSection, student.SchoolNo)
	}

	rows := [][2]string{
		{"Öğrencinin Adı Soyadı", stName},
		{"Sınıfı / Okul No", stClass},
		{"Görüşülen Veli Adı Soyadı", parentName},
		{"Görüşme Tarihi", dateStr},
		{"Görüşme Konusu", topic},
	}

	for _, r := range rows {
		h := 20.0
		y := doc.PDF.GetY()
		doc.DrawRect(doc.MarginLeft, y, lblW, h)
		doc.DrawRect(doc.MarginLeft+lblW, y, valW, h)

		doc.SetBoldFont(8.0)
		doc.PDF.SetX(doc.MarginLeft + 4)
		doc.PDF.SetY(y + 4)
		_ = doc.PDF.CellWithOption(&gopdf.Rect{W: lblW - 8, H: 12}, r[0], gopdf.CellOption{Align: gopdf.Left})

		doc.SetRegularFont(8.0)
		doc.PDF.SetX(doc.MarginLeft + lblW + 4)
		doc.PDF.SetY(y + 4)
		_ = doc.PDF.CellWithOption(&gopdf.Rect{W: valW - 8, H: 12}, r[1], gopdf.CellOption{Align: gopdf.Left})

		doc.PDF.SetY(y + h)
	}

	doc.PDF.Br(12)
	doc.SetBoldFont(8.5)
	doc.PDF.SetX(doc.MarginLeft)
	_ = doc.PDF.CellWithOption(&gopdf.Rect{W: doc.UsableWidth, H: 13}, "Görüşme Notları:", gopdf.CellOption{Align: gopdf.Left})
	doc.PDF.Br(14)

	if strings.TrimSpace(notes) != "" {
		doc.SetRegularFont(8.0)
		doc.PDF.SetX(doc.MarginLeft)
		_ = doc.PDF.MultiCell(&gopdf.Rect{W: doc.UsableWidth, H: 13}, notes)
		doc.PDF.Br(12)
	} else {
		doc.DrawDottedLines(7, 14)
	}

	doc.PDF.Br(10)
	doc.SetBoldFont(8.5)
	doc.PDF.SetX(doc.MarginLeft)
	_ = doc.PDF.CellWithOption(&gopdf.Rect{W: doc.UsableWidth, H: 13}, "Alınan Kararlar / Sonuç:", gopdf.CellOption{Align: gopdf.Left})
	doc.PDF.Br(14)

	if strings.TrimSpace(decision) != "" {
		doc.SetRegularFont(8.0)
		doc.PDF.SetX(doc.MarginLeft)
		_ = doc.PDF.MultiCell(&gopdf.Rect{W: doc.UsableWidth, H: 13}, decision)
		doc.PDF.Br(12)
	} else {
		doc.DrawDottedLines(5, 14)
	}

	doc.PDF.Br(25)

	// Dual signature block
	colW := doc.UsableWidth / 2.0
	ySig := doc.PDF.GetY()

	doc.SetBoldFont(8.0)
	doc.PDF.SetX(doc.MarginLeft)
	doc.PDF.SetY(ySig)
	_ = doc.PDF.CellWithOption(&gopdf.Rect{W: colW, H: 12}, "VELİ", gopdf.CellOption{Align: gopdf.Center})

	doc.PDF.SetX(doc.MarginLeft + colW)
	doc.PDF.SetY(ySig)
	_ = doc.PDF.CellWithOption(&gopdf.Rect{W: colW, H: 12}, "GÖRÜŞMEYİ YAPAN", gopdf.CellOption{Align: gopdf.Center})

	doc.SetRegularFont(7.5)
	doc.PDF.SetX(doc.MarginLeft)
	doc.PDF.SetY(ySig + 25)
	_ = doc.PDF.CellWithOption(&gopdf.Rect{W: colW, H: 11}, parentName, gopdf.CellOption{Align: gopdf.Center})

	iName := interviewer
	if iName == "" {
		iName = "........................................"
	}
	doc.PDF.SetX(doc.MarginLeft + colW)
	doc.PDF.SetY(ySig + 25)
	_ = doc.PDF.CellWithOption(&gopdf.Rect{W: colW, H: 11}, iName, gopdf.CellOption{Align: gopdf.Center})

	doc.PDF.SetX(doc.MarginLeft)
	doc.PDF.SetY(ySig + 38)
	_ = doc.PDF.CellWithOption(&gopdf.Rect{W: colW, H: 11}, "İmza: ..............................", gopdf.CellOption{Align: gopdf.Center})

	doc.PDF.SetX(doc.MarginLeft + colW)
	doc.PDF.SetY(ySig + 38)
	_ = doc.PDF.CellWithOption(&gopdf.Rect{W: colW, H: 11}, "İmza: ..............................", gopdf.CellOption{Align: gopdf.Center})

	doc.DrawFooterAddress(school)
	return doc.PDF.WritePdf(outputPath)
}

// GenerateRehberlikGorusFormu generates school counselor's report.
func GenerateRehberlikGorusFormu(school *model.SchoolInfo, student *model.Student, counselorNotes [4]string, outputPath string) error {
	doc, err := NewDoc(false)
	if err != nil {
		return err
	}

	doc.PDF.AddPage()

	boardName := "Ödül ve Disiplin"
	if strings.ToUpper(school.SchoolType) == "ORTAOKUL" {
		boardName = "Öğrenci Davranışlarını Değerlendirme"
	}

	doc.DrawTitle(fmt.Sprintf("%s MÜDÜRLÜĞÜNE\n(%s Kurulu Başkanlığına)", strings.ToUpper(school.SchoolName), boardName), 11, 12)

	stText := ""
	if student != nil {
		stText = fmt.Sprintf("Okulumuz %s sınıfı %s numaralı öğrencisi %s hakkında istenen rehberlik ve psikolojik danışmanlık inceleme görüşleri aşağıda belirtilmiştir.", student.ClassSection, student.SchoolNo, student.FullName)
	}

	doc.SetRegularFont(8.5)
	doc.PDF.SetX(doc.MarginLeft)
	_ = doc.PDF.MultiCell(&gopdf.Rect{W: doc.UsableWidth, H: 13}, stText)
	doc.PDF.Br(14)

	sections := []string{
		"1- Öğrencinin Ailevi Durumu:",
		"2- Öğrencinin Psikolojik ve Sosyal Özellikleri:",
		"3- Olayın Nedenleri Hakkındaki Görüşler:",
		"4- Öğrenciye Uygulanabilecek Yaptırım / Ceza veya Rehberlik Önlemleri Hakkındaki Öneriler:",
	}

	for i, sec := range sections {
		doc.SetBoldFont(8.0)
		doc.PDF.SetX(doc.MarginLeft)
		_ = doc.PDF.CellWithOption(&gopdf.Rect{W: doc.UsableWidth, H: 12}, sec, gopdf.CellOption{Align: gopdf.Left})
		doc.PDF.Br(12)

		note := counselorNotes[i]
		if strings.TrimSpace(note) != "" {
			doc.SetRegularFont(8.0)
			doc.PDF.SetX(doc.MarginLeft)
			_ = doc.PDF.MultiCell(&gopdf.Rect{W: doc.UsableWidth, H: 12}, note)
			doc.PDF.Br(8)
		} else {
			doc.DrawDottedLines(5, 13)
			doc.PDF.Br(4)
		}
	}

	doc.PDF.Br(20)
	doc.SetBoldFont(8.0)
	doc.PDF.SetX(doc.MarginLeft)
	_ = doc.PDF.CellWithOption(&gopdf.Rect{W: doc.UsableWidth, H: 12}, "Psikolojik Danışman / Rehber Öğretmen", gopdf.CellOption{Align: gopdf.Right})
	doc.PDF.Br(12)
	doc.SetRegularFont(7.5)
	doc.PDF.SetX(doc.MarginLeft)
	_ = doc.PDF.CellWithOption(&gopdf.Rect{W: doc.UsableWidth, H: 12}, "İmza: ........................................", gopdf.CellOption{Align: gopdf.Right})

	return doc.PDF.WritePdf(outputPath)
}

func choose(cond bool, a, b string) string {
	if cond {
		return a
	}
	return b
}
