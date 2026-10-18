package pdf

import (
	"fmt"
	"strings"

	"github.com/signintech/gopdf"

	"turandisiplin/model"
)

// GenerateCagriPusulasi generates student summons slips (1 page per student).
func GenerateCagriPusulasi(school *model.SchoolInfo, students []*model.Student, incident *model.Incident, isSavunma bool, place, dateStr, timeStr string, outputPath string) error {
	doc, err := NewDoc(false)
	if err != nil {
		return err
	}

	boardName := "ÖDÜL VE DİSİPLİN"
	citation := "Milli Eğitim Bakanlığı Ortaöğretim Kurumları Yönetmeliğinin 190. maddesi gereğince"
	if strings.ToUpper(school.SchoolType) == "ORTAOKUL" {
		boardName = "ÖĞRENCİ DAVRANIŞLARINI DEĞERLENDİRME"
		citation = "Milli Eğitim Bakanlığı Okul Öncesi Eğitim ve İlköğretim Kurumları Yönetmeliğinin ilgili maddesi gereğince"
	}

	title := "İ F A D E     Ç A Ğ R I     P U S U L A S I"
	callType := "İfade"
	callLocationLabel := "İFADENİN ALINACAĞI YER :"
	if isSavunma {
		title = "S A V U N M A     Ç A Ğ R I     P U S U L A S I"
		callType = "Savunma"
		callLocationLabel = "SAVUNMANIN ALINACAĞI YER :"
	}

	for idx, student := range students {
		doc.PDF.AddPage()

		// Table 1: Antet & Header
		doc.DrawAntet(school)
		doc.SetBoldFont(9)
		doc.PDF.SetX(doc.MarginLeft)
		_ = doc.PDF.CellWithOption(&gopdf.Rect{W: doc.UsableWidth, H: 14}, boardName+" KURULU", gopdf.CellOption{Align: gopdf.Center})
		doc.PDF.Br(18)

		doc.SetBoldFont(14)
		doc.PDF.SetX(doc.MarginLeft)
		_ = doc.PDF.CellWithOption(&gopdf.Rect{W: doc.UsableWidth, H: 20}, title, gopdf.CellOption{Align: gopdf.Center})
		doc.PDF.Br(24)

		// Table 2: Place & Time (indented by 30)
		indent := 30.0
		innerW := doc.UsableWidth - (indent * 2)
		lblW := 150.0
		valW := innerW - lblW

		timeBoxY := doc.PDF.GetY()
		doc.PDF.SetLineWidth(1.2)
		doc.PDF.Line(doc.MarginLeft+indent, timeBoxY, doc.MarginLeft+indent+innerW, timeBoxY)
		doc.PDF.Br(6)

		boxRows := [][2]string{
			{"ÖĞRETİM YILI :", school.AcademicYear},
			{callLocationLabel, place},
			{"TARİH VE SAATİ :", fmt.Sprintf("%s  -  %s", dateStr, timeStr)},
		}

		for _, br := range boxRows {
			doc.SetBoldFont(8.0)
			doc.PDF.SetX(doc.MarginLeft + indent)
			_ = doc.PDF.CellWithOption(&gopdf.Rect{W: lblW, H: 14}, br[0], gopdf.CellOption{Align: gopdf.Left})

			doc.SetRegularFont(8.0)
			doc.PDF.SetX(doc.MarginLeft + indent + lblW)
			_ = doc.PDF.CellWithOption(&gopdf.Rect{W: valW, H: 14}, br[1], gopdf.CellOption{Align: gopdf.Left})
			doc.PDF.Br(14)
		}

		doc.PDF.Br(2)
		timeBoxEnd := doc.PDF.GetY()
		doc.PDF.Line(doc.MarginLeft+indent, timeBoxEnd, doc.MarginLeft+indent+innerW, timeBoxEnd)
		doc.PDF.Br(15)

		// Table 3: Student info
		doc.SetBoldFont(8.0)
		doc.PDF.SetX(doc.MarginLeft + indent)
		_ = doc.PDF.CellWithOption(&gopdf.Rect{W: innerW, H: 12}, "Öğrencinin;", gopdf.CellOption{Align: gopdf.Left})
		doc.PDF.Br(14)

		doc.PDF.SetX(doc.MarginLeft + indent)
		_ = doc.PDF.CellWithOption(&gopdf.Rect{W: lblW, H: 13}, "ADI SOYADI :", gopdf.CellOption{Align: gopdf.Left})
		doc.SetRegularFont(8.0)
		doc.PDF.SetX(doc.MarginLeft + indent + lblW)
		_ = doc.PDF.CellWithOption(&gopdf.Rect{W: valW, H: 13}, student.FullName, gopdf.CellOption{Align: gopdf.Left})
		doc.PDF.Br(14)

		doc.SetBoldFont(8.0)
		doc.PDF.SetX(doc.MarginLeft + indent)
		_ = doc.PDF.CellWithOption(&gopdf.Rect{W: lblW, H: 13}, "SINIFI - NUMARASI :", gopdf.CellOption{Align: gopdf.Left})
		doc.SetRegularFont(8.0)
		doc.PDF.SetX(doc.MarginLeft + indent + lblW)
		_ = doc.PDF.CellWithOption(&gopdf.Rect{W: valW, H: 13}, fmt.Sprintf("%s - %s", student.ClassSection, student.SchoolNo), gopdf.CellOption{Align: gopdf.Left})
		doc.PDF.Br(20)

		// Table 4: Summons paragraph
		doc.SetRegularFont(8.0)
		doc.PDF.SetX(doc.MarginLeft + 15)
		para := fmt.Sprintf("Hakkınızda açılmış bulunan inceleme / soruşturma nedeniyle %s, yukarıda belirtilen gün ve saatte %snızı vermek üzere hazır bulunmanız rica olunur.", citation, strings.ToLower(callType))
		_ = doc.PDF.MultiCell(&gopdf.Rect{W: doc.UsableWidth - 30, H: 14}, para)
		doc.PDF.Br(25)

		// Chair signature
		doc.SetBoldFont(8.0)
		doc.PDF.SetX(doc.MarginLeft)
		_ = doc.PDF.CellWithOption(&gopdf.Rect{W: doc.UsableWidth - 20, H: 12}, school.BoardChair, gopdf.CellOption{Align: gopdf.Right})
		doc.PDF.Br(12)
		doc.SetRegularFont(7.5)
		doc.PDF.SetX(doc.MarginLeft)
		_ = doc.PDF.CellWithOption(&gopdf.Rect{W: doc.UsableWidth - 20, H: 11}, boardName+" Kurulu Başkanı", gopdf.CellOption{Align: gopdf.Right})
		doc.PDF.Br(11)
		doc.PDF.SetX(doc.MarginLeft)
		_ = doc.PDF.CellWithOption(&gopdf.Rect{W: doc.UsableWidth - 20, H: 11}, school.BoardChairTitle, gopdf.CellOption{Align: gopdf.Right})
		doc.PDF.Br(30)

		// Table 5: "OKUDUM" acknowledgement block
		doc.SetBoldFont(8.5)
		doc.PDF.SetX(doc.MarginLeft + indent)
		_ = doc.PDF.CellWithOption(&gopdf.Rect{W: innerW, H: 13}, "OKUDUM", gopdf.CellOption{Align: gopdf.Left})
		doc.PDF.Br(15)

		ackRows := [][2]string{
			{"ADI SOYADI :", student.FullName},
			{"TARİH :", "......./......./20......."},
			{"İMZA :", "................................................"},
		}

		for _, ar := range ackRows {
			doc.SetBoldFont(7.5)
			doc.PDF.SetX(doc.MarginLeft + indent)
			_ = doc.PDF.CellWithOption(&gopdf.Rect{W: 100, H: 13}, ar[0], gopdf.CellOption{Align: gopdf.Left})
			doc.SetRegularFont(7.5)
			doc.PDF.SetX(doc.MarginLeft + indent + 100)
			_ = doc.PDF.CellWithOption(&gopdf.Rect{W: innerW - 100, H: 13}, ar[1], gopdf.CellOption{Align: gopdf.Left})
			doc.PDF.Br(13)
		}

		// Table 6: Bottom address
		doc.DrawFooterAddress(school)

		_ = idx
	}

	return doc.PDF.WritePdf(outputPath)
}

// GenerateUyeCagriPusulasi generates summons for board members with agenda and signatures.
func GenerateUyeCagriPusulasi(school *model.SchoolInfo, sayi, konu, dateStr, timeStr, place string, agendaItems []string, outputPath string) error {
	doc, err := NewDoc(false)
	if err != nil {
		return err
	}

	doc.PDF.AddPage()

	boardName := "ÖDÜL VE DİSİPLİN"
	if strings.ToUpper(school.SchoolType) == "ORTAOKUL" {
		boardName = "ÖĞRENCİ DAVRANIŞLARINI DEĞERLENDİRME"
	}

	doc.DrawAntet(school)

	// Sayı & Konu
	doc.SetBoldFont(8.0)
	doc.PDF.SetX(doc.MarginLeft)
	_ = doc.PDF.CellWithOption(&gopdf.Rect{W: 50, H: 12}, "Sayı :", gopdf.CellOption{Align: gopdf.Left})
	doc.SetRegularFont(8.0)
	doc.PDF.SetX(doc.MarginLeft + 50)
	_ = doc.PDF.CellWithOption(&gopdf.Rect{W: 250, H: 12}, sayi, gopdf.CellOption{Align: gopdf.Left})

	doc.PDF.SetX(doc.MarginLeft + doc.UsableWidth - 100)
	_ = doc.PDF.CellWithOption(&gopdf.Rect{W: 100, H: 12}, dateStr, gopdf.CellOption{Align: gopdf.Right})
	doc.PDF.Br(13)

	doc.SetBoldFont(8.0)
	doc.PDF.SetX(doc.MarginLeft)
	_ = doc.PDF.CellWithOption(&gopdf.Rect{W: 50, H: 12}, "Konu :", gopdf.CellOption{Align: gopdf.Left})
	doc.SetRegularFont(8.0)
	doc.PDF.SetX(doc.MarginLeft + 50)
	_ = doc.PDF.CellWithOption(&gopdf.Rect{W: 300, H: 12}, konu, gopdf.CellOption{Align: gopdf.Left})
	doc.PDF.Br(20)

	// Header to members
	doc.SetBoldFont(9.0)
	doc.PDF.SetX(doc.MarginLeft)
	_ = doc.PDF.CellWithOption(&gopdf.Rect{W: doc.UsableWidth, H: 14}, boardName+" KURULU ÜYELERİNE", gopdf.CellOption{Align: gopdf.Center})
	doc.PDF.Br(18)

	// Body invitation
	doc.SetRegularFont(8.5)
	doc.PDF.SetX(doc.MarginLeft)
	body := fmt.Sprintf("Kurulumuz aşağıda belirtilen gündem maddelerini görüşmek üzere %s günü saat %s'de %s toplanacaktır. Belirtilen gün ve saatte toplantıda hazır bulunmanızı rica ederim.", dateStr, timeStr, place)
	_ = doc.PDF.MultiCell(&gopdf.Rect{W: doc.UsableWidth, H: 14}, body)
	doc.PDF.Br(22)

	// Dual signatures: Left Chair, Right Principal
	colW := doc.UsableWidth / 2.0
	ySig := doc.PDF.GetY()

	doc.SetBoldFont(8.0)
	doc.PDF.SetX(doc.MarginLeft)
	doc.PDF.SetY(ySig)
	_ = doc.PDF.CellWithOption(&gopdf.Rect{W: colW, H: 12}, school.BoardChair, gopdf.CellOption{Align: gopdf.Center})

	doc.PDF.SetX(doc.MarginLeft + colW)
	doc.PDF.SetY(ySig)
	_ = doc.PDF.CellWithOption(&gopdf.Rect{W: colW, H: 12}, school.PrincipalName, gopdf.CellOption{Align: gopdf.Center})

	doc.SetRegularFont(7.5)
	doc.PDF.SetX(doc.MarginLeft)
	doc.PDF.SetY(ySig + 13)
	_ = doc.PDF.CellWithOption(&gopdf.Rect{W: colW, H: 11}, boardName+" Kurulu Başkanı", gopdf.CellOption{Align: gopdf.Center})

	doc.PDF.SetX(doc.MarginLeft + colW)
	doc.PDF.SetY(ySig + 13)
	_ = doc.PDF.CellWithOption(&gopdf.Rect{W: colW, H: 11}, "Okul Müdürü", gopdf.CellOption{Align: gopdf.Center})

	doc.PDF.SetY(ySig + 35)

	// Agenda Items Table
	doc.SetBoldFont(9.0)
	doc.PDF.SetX(doc.MarginLeft + 20)
	_ = doc.PDF.CellWithOption(&gopdf.Rect{W: doc.UsableWidth - 40, H: 13}, "G Ü N D E M", gopdf.CellOption{Align: gopdf.Left})
	doc.PDF.Br(15)

	doc.SetRegularFont(8.0)
	for i, item := range agendaItems {
		item = strings.TrimSpace(item)
		if item == "" {
			continue
		}
		doc.PDF.SetX(doc.MarginLeft + 20)
		txt := fmt.Sprintf("%d) %s", i+1, item)
		_ = doc.PDF.MultiCell(&gopdf.Rect{W: doc.UsableWidth - 40, H: 13}, txt)
		doc.PDF.Br(4)
	}

	doc.PDF.Br(18)

	// Members signature table
	indentM := 60.0
	tableW := doc.UsableWidth - (indentM * 2)
	wName := tableW * 0.65
	wSign := tableW * 0.35

	yTable := doc.PDF.GetY()
	doc.DrawRect(doc.MarginLeft+indentM, yTable, wName, 18)
	doc.DrawRect(doc.MarginLeft+indentM+wName, yTable, wSign, 18)

	doc.SetBoldFont(8.0)
	doc.PDF.SetX(doc.MarginLeft + indentM + 5)
	doc.PDF.SetY(yTable + 3)
	_ = doc.PDF.CellWithOption(&gopdf.Rect{W: wName - 10, H: 12}, "ÜYELER", gopdf.CellOption{Align: gopdf.Left})

	doc.PDF.SetX(doc.MarginLeft + indentM + wName + 5)
	doc.PDF.SetY(yTable + 3)
	_ = doc.PDF.CellWithOption(&gopdf.Rect{W: wSign - 10, H: 12}, "İMZA", gopdf.CellOption{Align: gopdf.Center})

	members := []string{school.Member1, school.Member2, school.Member3, school.Member4}
	currentY := yTable + 18

	doc.SetRegularFont(7.5)
	for _, m := range members {
		mName := m
		if mName == "" {
			mName = "........................................"
		}
		doc.DrawRect(doc.MarginLeft+indentM, currentY, wName, 20)
		doc.DrawRect(doc.MarginLeft+indentM+wName, currentY, wSign, 20)

		doc.PDF.SetX(doc.MarginLeft + indentM + 5)
		doc.PDF.SetY(currentY + 4)
		_ = doc.PDF.CellWithOption(&gopdf.Rect{W: wName - 10, H: 12}, mName, gopdf.CellOption{Align: gopdf.Left})

		doc.PDF.SetX(doc.MarginLeft + indentM + wName + 5)
		doc.PDF.SetY(currentY + 4)
		_ = doc.PDF.CellWithOption(&gopdf.Rect{W: wSign - 10, H: 12}, "........................", gopdf.CellOption{Align: gopdf.Center})

		currentY += 20
	}

	doc.DrawFooterAddress(school)
	return doc.PDF.WritePdf(outputPath)
}

// GenerateCagriOgrenciListesi prints the table of summoned students.
func GenerateCagriOgrenciListesi(school *model.SchoolInfo, students []*model.Student, outputPath string) error {
	doc, err := NewDoc(false)
	if err != nil {
		return err
	}

	doc.PDF.AddPage()
	doc.DrawAntet(school)
	doc.DrawTitle("ÇAĞRI PUSULASI ÖĞRENCİ LİSTESİ", 12, 15)

	wSira := 40.0
	wName := 250.0
	wClass := doc.UsableWidth - wSira - wName

	// Header row
	y := doc.PDF.GetY()
	doc.DrawRect(doc.MarginLeft, y, wSira, 20)
	doc.DrawRect(doc.MarginLeft+wSira, y, wName, 20)
	doc.DrawRect(doc.MarginLeft+wSira+wName, y, wClass, 20)

	doc.SetBoldFont(8.0)
	doc.PDF.SetX(doc.MarginLeft)
	doc.PDF.SetY(y + 4)
	_ = doc.PDF.CellWithOption(&gopdf.Rect{W: wSira, H: 12}, "SIRA", gopdf.CellOption{Align: gopdf.Center})

	doc.PDF.SetX(doc.MarginLeft + wSira + 5)
	doc.PDF.SetY(y + 4)
	_ = doc.PDF.CellWithOption(&gopdf.Rect{W: wName - 10, H: 12}, "ADI SOYADI", gopdf.CellOption{Align: gopdf.Left})

	doc.PDF.SetX(doc.MarginLeft + wSira + wName + 5)
	doc.PDF.SetY(y + 4)
	_ = doc.PDF.CellWithOption(&gopdf.Rect{W: wClass - 10, H: 12}, "SINIFI - NUMARASI", gopdf.CellOption{Align: gopdf.Left})

	curY := y + 20
	doc.SetRegularFont(8.0)

	for i, st := range students {
		doc.DrawRect(doc.MarginLeft, curY, wSira, 18)
		doc.DrawRect(doc.MarginLeft+wSira, curY, wName, 18)
		doc.DrawRect(doc.MarginLeft+wSira+wName, curY, wClass, 18)

		doc.PDF.SetX(doc.MarginLeft)
		doc.PDF.SetY(curY + 3)
		_ = doc.PDF.CellWithOption(&gopdf.Rect{W: wSira, H: 12}, fmt.Sprintf("%d", i+1), gopdf.CellOption{Align: gopdf.Center})

		doc.PDF.SetX(doc.MarginLeft + wSira + 5)
		doc.PDF.SetY(curY + 3)
		_ = doc.PDF.CellWithOption(&gopdf.Rect{W: wName - 10, H: 12}, st.FullName, gopdf.CellOption{Align: gopdf.Left})

		doc.PDF.SetX(doc.MarginLeft + wSira + wName + 5)
		doc.PDF.SetY(curY + 3)
		_ = doc.PDF.CellWithOption(&gopdf.Rect{W: wClass - 10, H: 12}, fmt.Sprintf("%s - %s", st.ClassSection, st.SchoolNo), gopdf.CellOption{Align: gopdf.Left})

		curY += 18
	}

	doc.PDF.SetY(curY + 25)
	doc.SetBoldFont(8.0)
	doc.PDF.SetX(doc.MarginLeft)
	_ = doc.PDF.CellWithOption(&gopdf.Rect{W: doc.UsableWidth, H: 12}, school.BoardChair, gopdf.CellOption{Align: gopdf.Right})
	doc.PDF.Br(12)
	doc.SetRegularFont(7.5)
	doc.PDF.SetX(doc.MarginLeft)
	_ = doc.PDF.CellWithOption(&gopdf.Rect{W: doc.UsableWidth, H: 12}, "Kurul Başkanı", gopdf.CellOption{Align: gopdf.Right})

	return doc.PDF.WritePdf(outputPath)
}

// GenerateYaziliUyariTutanagi generates the written warning document.
func GenerateYaziliUyariTutanagi(school *model.SchoolInfo, student *model.Student, sayi, dateStr, clause, explanation, vicePrincipal string, outputPath string) error {
	doc, err := NewDoc(false)
	if err != nil {
		return err
	}

	doc.PDF.AddPage()
	doc.DrawAntet(school)

	// Sayı & Tarih
	doc.SetBoldFont(8.0)
	doc.PDF.SetX(doc.MarginLeft)
	_ = doc.PDF.CellWithOption(&gopdf.Rect{W: 50, H: 12}, "Sayı :", gopdf.CellOption{Align: gopdf.Left})
	doc.SetRegularFont(8.0)
	doc.PDF.SetX(doc.MarginLeft + 50)
	_ = doc.PDF.CellWithOption(&gopdf.Rect{W: 200, H: 12}, sayi, gopdf.CellOption{Align: gopdf.Left})

	doc.PDF.SetX(doc.MarginLeft + doc.UsableWidth - 100)
	_ = doc.PDF.CellWithOption(&gopdf.Rect{W: 100, H: 12}, dateStr, gopdf.CellOption{Align: gopdf.Right})
	doc.PDF.Br(14)

	doc.SetBoldFont(8.0)
	doc.PDF.SetX(doc.MarginLeft)
	_ = doc.PDF.CellWithOption(&gopdf.Rect{W: 50, H: 12}, "Konu :", gopdf.CellOption{Align: gopdf.Left})
	doc.SetRegularFont(8.0)
	doc.PDF.SetX(doc.MarginLeft + 50)
	_ = doc.PDF.CellWithOption(&gopdf.Rect{W: 200, H: 12}, "Yazılı Uyarı", gopdf.CellOption{Align: gopdf.Left})
	doc.PDF.Br(20)

	doc.DrawTitle("ÖĞRENCİ YAZILI UYARI TUTANAĞI", 12, 15)

	lblW := 140.0
	valW := doc.UsableWidth - lblW

	rows := [][2]string{
		{"Adı Soyadı", student.FullName},
		{"Sınıfı - Okul No", fmt.Sprintf("%s - %s", student.ClassSection, student.SchoolNo)},
		{"Uyarı Tarihi", dateStr},
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

	doc.PDF.Br(18)
	doc.SetRegularFont(8.5)
	doc.PDF.SetX(doc.MarginLeft)
	warnText := fmt.Sprintf("Milli Eğitim Bakanlığı Ortaöğretim Kurumları Yönetmeliğinin 163/2 maddesi ve %s maddesi (%s) kapsamında gösterdiğiniz olumsuz davranış nedeniyle işbu yazılı uyarı tarafınıza tebliğ edilmiştir. Aynı veya benzer olumsuz davranışların tekrarlanması durumunda disiplin işlemi uygulanacağı hususunda bilginizi ve gereğini rica ederim.", clause, explanation)
	_ = doc.PDF.MultiCell(&gopdf.Rect{W: doc.UsableWidth, H: 14}, warnText)
	doc.PDF.Br(30)

	// Principal
	doc.SetBoldFont(8.0)
	doc.PDF.SetX(doc.MarginLeft)
	_ = doc.PDF.CellWithOption(&gopdf.Rect{W: doc.UsableWidth, H: 12}, school.PrincipalName, gopdf.CellOption{Align: gopdf.Right})
	doc.PDF.Br(12)
	doc.SetRegularFont(7.5)
	doc.PDF.SetX(doc.MarginLeft)
	_ = doc.PDF.CellWithOption(&gopdf.Rect{W: doc.UsableWidth, H: 12}, "Okul Müdürü", gopdf.CellOption{Align: gopdf.Right})
	doc.PDF.Br(35)

	// Acknowledgement
	doc.SetRegularFont(8.0)
	doc.PDF.SetX(doc.MarginLeft)
	ack := "Yukarıdaki yazılı uyarıyı okudum ve bir nüshasını elden teslim aldım."
	_ = doc.PDF.CellWithOption(&gopdf.Rect{W: doc.UsableWidth, H: 12}, ack, gopdf.CellOption{Align: gopdf.Left})
	doc.PDF.Br(20)

	colW := doc.UsableWidth / 2.0
	yAck := doc.PDF.GetY()

	doc.SetBoldFont(7.5)
	doc.PDF.SetX(doc.MarginLeft)
	doc.PDF.SetY(yAck)
	_ = doc.PDF.CellWithOption(&gopdf.Rect{W: colW, H: 12}, "ÖĞRENCİ", gopdf.CellOption{Align: gopdf.Center})

	doc.PDF.SetX(doc.MarginLeft + colW)
	doc.PDF.SetY(yAck)
	_ = doc.PDF.CellWithOption(&gopdf.Rect{W: colW, H: 12}, "TEBLİĞ EDEN MÜDÜR YARDIMCISI", gopdf.CellOption{Align: gopdf.Center})

	doc.SetRegularFont(7.5)
	doc.PDF.SetX(doc.MarginLeft)
	doc.PDF.SetY(yAck + 20)
	_ = doc.PDF.CellWithOption(&gopdf.Rect{W: colW, H: 11}, student.FullName, gopdf.CellOption{Align: gopdf.Center})

	vpName := vicePrincipal
	if vpName == "" {
		vpName = "........................................"
	}
	doc.PDF.SetX(doc.MarginLeft + colW)
	doc.PDF.SetY(yAck + 20)
	_ = doc.PDF.CellWithOption(&gopdf.Rect{W: colW, H: 11}, vpName, gopdf.CellOption{Align: gopdf.Center})

	doc.PDF.SetX(doc.MarginLeft)
	doc.PDF.SetY(yAck + 32)
	_ = doc.PDF.CellWithOption(&gopdf.Rect{W: colW, H: 11}, "İmza: ............................", gopdf.CellOption{Align: gopdf.Center})

	doc.PDF.SetX(doc.MarginLeft + colW)
	doc.PDF.SetY(yAck + 32)
	_ = doc.PDF.CellWithOption(&gopdf.Rect{W: colW, H: 11}, "İmza: ............................", gopdf.CellOption{Align: gopdf.Center})

	return doc.PDF.WritePdf(outputPath)
}
