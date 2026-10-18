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

// BuildMinutesTab creates the Beginning/End of Year Meeting Minutes & Case File Index (Dizi Pusulası) UI
func BuildMinutesTab(state *AppState) fyne.CanvasObject {
	// --- SECTION 1: Toplantı Tutanağı (Sene Başı / Sene Sonu) ---
	meetingTypeRadio := widget.NewRadioGroup([]string{"SENE BAŞI", "SENE SONU"}, nil)
	meetingTypeRadio.Horizontal = true
	meetingTypeRadio.SetSelected("SENE BAŞI")

	mDateEntry := widget.NewEntry()
	mDateEntry.SetText(time.Now().Format("02.01.2006"))
	mTimeEntry := widget.NewEntry()
	mTimeEntry.SetText("14:00")
	mPlaceEntry := widget.NewEntry()
	mPlaceEntry.SetText("Öğretmenler Odası")
	mChairEntry := widget.NewEntry()
	mChairEntry.SetText(state.SchoolInfo.BoardChair)

	// Decision items entry (multi-line)
	decisionsText := widget.NewMultiLineEntry()

	loadDefaultDecisions := func(mType string) {
		if mType == "SENE BAŞI" {
			decisionsText.SetText(
				"1. Milli Eğitim Bakanlığı Ortaöğretim Kurumları Yönetmeliğinin ilgili maddeleri doğrultusunda kurulun görev ve sorumlulukları gözden geçirilmiştir.\n" +
					"2. Okulda düzen ve disiplini sağlamak amacıyla rehberlik ve psikolojik danışma servisi ile iş birliği yapılarak önleyici rehberlik çalışmalarına ağırlık verilmesi kararlaştırılmıştır.\n" +
					"3. Öğrenci devamsızlıklarının düzenli olarak takip edilmesi ve devamsızlığı alışkanlık haline getiren öğrencilerin velileri ile görüşülmesi kararlaştırılmıştır.\n" +
					"4. Bilişim araçlarının okul içerisinde izinsiz ve kurallara aykırı kullanımına karşı gerekli tedbirlerin alınması uygun görülmüştür.\n" +
					"5. Olumlu davranış gösteren öğrencilerin Onur Belgesi ile ödüllendirilmesi için sınıf rehber öğretmenleri ile koordinasyon sağlanması kararlaştırılmıştır.",
			)
		} else {
			decisionsText.SetText(
				"1. Eğitim ve öğretim yılı boyunca disiplin kuruluna intikal eden olaylar ve alınan kararlar değerlendirilmiştir.\n" +
					"2. Alınan disiplin cezalarının dönemler itibariyle dağılımı incelenmiş, sene sonu itibariyle davranışlarında olumlu yönde gelişme gösteren öğrencilerin durumları görüşülmüştür.\n" +
					"3. Yönetmeliğin 171/2 maddesi uyarınca olumlu davranış sergileyen ve pişmanlık duyan öğrencilerin davranış puanlarının iadesine karar verilmiştir.\n" +
					"4. Gelecek eğitim ve öğretim yılında öğrenci devamsızlığı ve akran zorbalığı konularında alınacak tedbirler belirlenmiştir.\n" +
					"5. İş bu kurul tutanağının sene sonu öğretmenler kurulunda okunması talebiyle okul müdürlüğünün onayına sunulmasına karar verilmiştir.",
			)
		}
	}

	loadDefaultDecisions("SENE BAŞI")

	meetingTypeRadio.OnChanged = func(s string) {
		loadDefaultDecisions(s)
	}

	// Append Statistics button
	appendStatsBtn := widget.NewButtonWithIcon("Ceza İstatistiklerini Metne Ekle", theme.ContentAddIcon(), func() {
		stats, err := state.Store.CalculateSemesterPenaltyCounts(state.AcademicYear)
		if err != nil || len(stats) == 0 {
			showError(state.Window, fmt.Errorf("Bu öğretim yılına ait istatistik verisi bulunamadı"))
			return
		}

		var sb strings.Builder
		sb.WriteString("\n--- DÖNEMLERE GÖRE CEZA DAĞILIMI ---\n")
		for _, s := range stats {
			sb.WriteString(fmt.Sprintf("- %s: 1. Dönem: %d, 2. Dönem: %d, Toplam: %d\n", s.Name, s.Sem1Count, s.Sem2Count, s.TotalCount))
		}
		decisionsText.SetText(decisionsText.Text + sb.String())
	})

	pdfMeetingBtn := widget.NewButtonWithIcon("📄 PDF OLUŞTUR VE AÇ (TOPLANTI TUTANAĞI)", theme.DocumentIcon(), func() {
		rawLines := strings.Split(decisionsText.Text, "\n")
		var decs []string
		for _, l := range rawLines {
			l = strings.TrimSpace(l)
			if l != "" {
				decs = append(decs, l)
			}
		}

		mType := meetingTypeRadio.Selected
		prefix := fmt.Sprintf("%s_Toplanti_Tutanagi", strings.ReplaceAll(mType, " ", "_"))

		outPath, err := pdf.GenerateUniquePDFPath("", prefix)
		if err != nil {
			showError(state.Window, err)
			return
		}

		err = pdf.GenerateSeneBasiSonuKarar(state.SchoolInfo, mType, mDateEntry.Text, mTimeEntry.Text, mPlaceEntry.Text, mChairEntry.Text, decs, outPath)
		if err != nil {
			showError(state.Window, fmt.Errorf("PDF oluşturulamadı: %w", err))
			return
		}

		_ = pdf.OpenPDF(outPath)
	})
	pdfMeetingBtn.Importance = widget.SuccessImportance

	// --- SECTION 2: Dizi Pusulası (Case File Document Index) ---
	diziEntries := widget.NewMultiLineEntry()
	diziEntries.SetText(
		"1. Olay Bildirim / Şikayet Dilekçesi | 1\n" +
			"2. İfade Tutanağı | 1\n" +
			"3. Savunma Çağrı Pusulası | 1\n" +
			"4. Savunma Tutanağı | 1\n" +
			"5. Bilgi Toplama Formları ve Tanık İfadeleri | 3\n" +
			"6. Rehberlik ve Psikolojik Danışma Raporu | 1\n" +
			"7. Disiplin Kurulu Karar Tutanağı (EK-1) | 1\n" +
			"8. Ceza Tebligat Belgeleri (Öğrenci ve Veli) | 2",
	)

	pdfDiziBtn := widget.NewButtonWithIcon("📄 PDF OLUŞTUR VE AÇ (DİZİ PUSULASI)", theme.DocumentIcon(), func() {
		lines := strings.Split(diziEntries.Text, "\n")
		var items []model.DiziItem
		for _, l := range lines {
			l = strings.TrimSpace(l)
			if l == "" {
				continue
			}
			parts := strings.Split(l, "|")
			subj := strings.TrimSpace(parts[0])
			cnt := 1
			if len(parts) > 1 {
				c, err := strconv.Atoi(strings.TrimSpace(parts[1]))
				if err == nil && c > 0 {
					cnt = c
				}
			}
			items = append(items, model.DiziItem{Subject: subj, Count: cnt})
		}

		outPath, err := pdf.GenerateUniquePDFPath("", "Dizi_Pusulasi")
		if err != nil {
			showError(state.Window, err)
			return
		}

		err = pdf.GenerateDiziPusulasi(state.SchoolInfo, items, outPath)
		if err != nil {
			showError(state.Window, fmt.Errorf("PDF oluşturulamadı: %w", err))
			return
		}

		_ = pdf.OpenPDF(outPath)
	})
	pdfDiziBtn.Importance = widget.SuccessImportance

	// Forms Assembly
	meetingForm := widget.NewForm(
		widget.NewFormItem("Toplantı Türü", meetingTypeRadio),
		widget.NewFormItem("Tarih", mDateEntry),
		widget.NewFormItem("Saat", mTimeEntry),
		widget.NewFormItem("Toplantı Yeri", mPlaceEntry),
		widget.NewFormItem("Toplantı Başkanı", mChairEntry),
		widget.NewFormItem("Alınan Kararlar", decisionsText),
	)

	cardMeeting := widget.NewCard("Sene Başı / Sene Sonu Kurul Toplantı Tutanağı", "Mevzuat gereği dönem başı ve sonunda düzenlenen resmi kurul toplantı tutanağı", container.NewVBox(
		meetingForm,
		container.NewHBox(appendStatsBtn, widget.NewSeparator(), pdfMeetingBtn),
	))

	cardDizi := widget.NewCard("Dizi Pusulası (Dosya Evrak İndeksi)", "Disiplin soruşturma dosyasındaki resmi evrakların listesi (Konu | Sayı formatında)", container.NewVBox(
		diziEntries,
		pdfDiziBtn,
	))

	content := container.NewVScroll(container.NewVBox(
		cardMeeting,
		widget.NewSeparator(),
		cardDizi,
	))

	return content
}
