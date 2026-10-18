package ui

import (
	"errors"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	fyneDialog "fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/sqweek/dialog"

	"okulevrak/store"
)

// SelectExcelFile opens the operating system's native file picker (Windows Explorer or Linux Zenity/KDialog)
func SelectExcelFile(w fyne.Window, onSelected func(filePath string)) {
	go func() {
		filePath, err := dialog.File().
			Title("Excel Öğrenci Listesi Seçin (.xlsx / .xls)").
			Filter("Excel Dosyaları (*.xlsx, *.xls)", "xlsx", "xls").
			Load()

		if err != nil {
			if errors.Is(err, dialog.ErrCancelled) {
				return
			}
			// If native dialog fails (e.g. headless or missing GUI tool), fall back to Fyne's dialog
			fyneDialog.ShowFileOpen(func(reader fyne.URIReadCloser, err error) {
				if err != nil || reader == nil {
					return
				}
				p := reader.URI().Path()
				_ = reader.Close()
				onSelected(p)
			}, w)
			return
		}

		if filePath != "" {
			onSelected(filePath)
		}
	}()
}

// ShowAcademicYearPicker opens a clean modal allowing switching or adding academic years
func ShowAcademicYearPicker(state *AppState) {
	years := state.Store.GetAcademicYears()
	var selectedYear = state.AcademicYear

	yearSelect := widget.NewSelect(years, func(s string) {
		selectedYear = s
	})
	yearSelect.SetSelected(state.AcademicYear)

	newYearEntry := widget.NewEntry()
	newYearEntry.SetPlaceHolder("Yeni Yıl (Örn: 2029 - 2030)")

	content := container.NewVBox(
		widget.NewLabelWithStyle("Aktif Öğretim Yılını Seçin:", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		yearSelect,
		widget.NewSeparator(),
		widget.NewLabel("Veya yeni bir öğretim yılı yazıp ekleyin:"),
		newYearEntry,
	)

	fyneDialog.ShowCustomConfirm("Öğretim Yılı Seçimi", "Uygula", "İptal", content, func(ok bool) {
		if !ok {
			return
		}
		chosen := selectedYear
		if strings.TrimSpace(newYearEntry.Text) != "" {
			chosen = store.NormalizeAcademicYear(newYearEntry.Text)
		}
		if chosen != "" && chosen != state.AcademicYear {
			state.SetAcademicYear(chosen)
		}
	}, state.Window)
}

// ShowHelpDialog displays comprehensive, senior-friendly guidance on e-Okul imports, incidents, and procedures
func ShowHelpDialog(state *AppState) {
	titleLabel := widget.NewLabelWithStyle("📖 OkulEvrak Kullanım Rehberi ve Yardım", fyne.TextAlignCenter, fyne.TextStyle{Bold: true})

	makeSection := func(heading string, body string) fyne.CanvasObject {
		h := widget.NewLabelWithStyle(heading, fyne.TextAlignLeading, fyne.TextStyle{Bold: true})
		b := widget.NewLabel(body)
		b.Wrapping = fyne.TextWrapWord
		return container.NewVBox(h, b, widget.NewSeparator())
	}

	sec1 := makeSection(
		"1. e-Okul ve e-Mesem Excel Dosyasını İçe Aktarma",
		"• e-Okul Yönetim Bilgi Sistemi > Kurum İşlemleri > Öğrenci İşlemleri > Raporlar bölümünden 'Sınıf Listesi' veya 'Genel Öğrenci Listesi' Excel dosyasını (.xlsx veya .xls) bilgisayarınıza indirin.\n"+
			"• OkulEvrak'ta 'Öğrenciler' ekranını açın ve 'Excel'den Öğrenci Yükle' butonuna tıklayın.\n"+
			"• Açılan standart sistem dosya yöneticisinden indirdiğiniz Excel dosyasını seçin.\n"+
			"• Program; Öğrenci Numarası, T.C. Kimlik No, Ad Soyad, Sınıf/Şube ve Cinsiyet sütunlarını otomatik algılar.\n"+
			"• Karşınıza gelen 3 aktarım modundan birini seçebilirsiniz:\n"+
			"   1) Sadece yeni öğrencileri ekle (Mevcut kayıtlar korunur)\n"+
			"   2) Tüm öğrencileri sil ve sıfırdan yükle\n"+
			"   3) Mevcut öğrencilerin sınıf ve şubelerini güncelle",
	)

	sec2 := makeSection(
		"2. Disiplin Olayı ve Soruşturma Girişi",
		"• 'Olay Girişi' ekranından olay konusunu, gerçekleştiği yeri ve tarih/saatini kaydedin.\n"+
			"• 'Olaya Öğrenci Ekleme' bölümünden olaya dahil olan öğrencileri arayarak Fail, Mağdur veya Tanık rolleriyle olaya bağlayın.\n"+
			"• Sosyo-ekonomik ve psikolojik inceleme profillerini doldurabilirsiniz.",
	)

	sec3 := makeSection(
		"3. Çağrı Pusulaları ve İfade/Savunma Formları",
		"• Kurul toplanmadan önce öğrenciye, veliye ve kurul üyelerine gönderilecek resmi davetiyeleri 'Çağrı Pusulaları' ekranından tek tıkla resmi PDF olarak çıkartabilirsiniz.\n"+
			"• Öğrenci savunma ve tanık ifade tutanaklarını 'Tutanak ve Formlar' menüsünden yazdırabilirsiniz.",
	)

	sec4 := makeSection(
		"4. Disiplin Kurulu Kararı (EK-1 / EK-10)",
		"• 'KARAR' ekranında olayı ve öğrenciyi seçin.\n"+
			"• MEB Yönetmeliği maddelerinden uygun cezayı seçin. Sistem resmi gerekçeli karar metnini otomatik oluşturur.\n"+
			"• 'Kararı Kaydet' dedikten sonra 'PDF OLUŞTUR VE AÇ' butonuna basarak resmi MEB formatındaki EK-1 veya EK-10 karar tutanağını yazdırılabilir PDF olarak alın.",
	)

	sec5 := makeSection(
		"5. Ceza Tebligatı, Dizi Pusulası ve İade Çizelgeleri",
		"• Kararın ardından öğrenci ve veli tebligat pusulalarını 'Ceza Tebligatı' ekranından yazdırın.\n"+
			"• Evrakların üst makama veya arşive teslimi için 'Dizi Pusulası' belgesini tek tıkla oluşturun.\n"+
			"• Dönem veya sene sonunda öğrencilerin davranış puanlarının iadesini ve silinmesini 'Davranış Puanı İade' ekranından yönetin.",
	)

	sec6 := makeSection(
		"6. Güvenlik, Yedekleme ve Bağımsız Çalışma",
		"• Program tüm verileri bilgisayarınızdaki yerel veritabanında (OkulEvrak.db) saklar.\n"+
			"• İnternet bağlantısı veya dış sunucu gerektirmez; verileriniz kurumunuzda güvenle kalır.\n"+
			"• Oluşturulan tüm resmi evraklar uygulama klasöründeki 'Çıktılar' dizininde arşivlenir.",
	)

	helpContent := container.NewVBox(
		titleLabel,
		widget.NewSeparator(),
		sec1,
		sec2,
		sec3,
		sec4,
		sec5,
		sec6,
	)

	scroll := container.NewVScroll(helpContent)
	scroll.SetMinSize(fyne.NewSize(650, 480))

	fyneDialog.ShowCustom("Kullanım Rehberi & Yardım", "Kapat", scroll, state.Window)
}
