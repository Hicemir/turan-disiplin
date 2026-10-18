package ui

import (
	"fmt"
	"image/color"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"okulevrak/pdf"
)

// BuildHubView renders either the Modern Dashboard Hub or the 1:1 Legacy replica
func BuildHubView(state *AppState) fyne.CanvasObject {
	if state.Config.UIMode == "legacy" {
		return buildLegacyHub(state)
	}
	return buildModernHub(state)
}

// -------------------------------------------------------------
// 1. MODERN DASHBOARD HUB (Senior-friendly, big clean cards)
// -------------------------------------------------------------

func buildModernHub(state *AppState) fyne.CanvasObject {
	header := BuildHeader(state)

	createNavCard := func(title, desc string, icon fyne.Resource, isPrimary bool, targetView string) fyne.CanvasObject {
		btn := widget.NewButtonWithIcon("Aç", theme.NavigateNextIcon(), func() {
			state.NavigateTo(targetView)
		})
		if isPrimary {
			btn.Importance = widget.HighImportance
		}

		titleLabel := widget.NewLabelWithStyle(title, fyne.TextAlignLeading, fyne.TextStyle{Bold: true})
		descLabel := widget.NewLabel(desc)
		descLabel.Wrapping = fyne.TextWrapWord

		cardContent := container.NewVBox(
			container.NewHBox(widget.NewIcon(icon), titleLabel),
			descLabel,
			btn,
		)
		return widget.NewCard("", "", cardContent)
	}

	// Column 1: Girişler & Tanımlar
	col1 := container.NewVBox(
		widget.NewLabelWithStyle("📝 GİRİŞLER & YÖNETİM", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		createNavCard("🏛️ Kurum ve Kurul Bilgileri", "Okul müdürü, asil ve yedek kurul üyeleri, resmi yazı kodları ve antet ayarları.", theme.HomeIcon(), false, "school"),
		createNavCard("👥 Öğrenci İşlemleri & Excel", "Öğrenci listesi, e-Okul / e-Mesem'den 3 modlu Excel aktarımı, öğrenci arama ve düzenleme.", theme.AccountIcon(), false, "students"),
		createNavCard("🐞 Olay Girişi & Soruşturma", "İncelenen disiplin olayları, öğrenci rolleri ve sosyo-ekonomik inceleme profilleri.", theme.WarningIcon(), false, "incidents"),
		createNavCard("💬 Öğrenci Yazılı Uyarı", "MEB Yönetmeliği 163/2 uyarınca resmi öğrenci yazılı uyarı tutanağı oluşturma.", theme.MailComposeIcon(), false, "written_warning"),
	)

	// Column 2: Kararlar & Resmi Evraklar
	col2 := container.NewVBox(
		widget.NewLabelWithStyle("📜 KARARLAR & EVRAKLAR", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		createNavCard("⚖️ DİSİPLİN KURULU KARARI (EK-1 / EK-10)", "21/17 satırlı resmi karar tutanağı, puan hesaplama, otomatik gerekçeli karar ve PDF.", theme.DocumentIcon(), true, "verdict"),
		createNavCard("📨 Çağrı Pusulaları", "Öğrenci savunma ve ifade çağrı pusulaları, tebellüğ bloğu ve kurul üye toplantı çağrısı.", theme.MailSendIcon(), false, "summons"),
		createNavCard("📝 Tutanak ve Formlar", "İfade, savunma, bilgi toplama (5 tür), ek süre, şikayet ve rehberlik formları.", theme.FileTextIcon(), false, "forms"),
		createNavCard("📬 Ceza Tebligatları", "Öğrenci ve veli karar tebligat pusulaları (5 günlük yasal itiraz hakkı ile).", theme.FolderIcon(), false, "notices"),
		createNavCard("☀️ Sene Başı / Sonu Kurul Tutanağı", "Dönem başı ve sonu resmi disiplin kurulu toplantı tutanakları.", theme.ContentPasteIcon(), false, "minutes"),
	)

	// Column 3: Raporlar & İnceleme
	col3 := container.NewVBox(
		widget.NewLabelWithStyle("📊 RAPORLAR & İSTATİSTİKLER", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		createNavCard("🔄 Puan İadesi / Yaptırım Kaldırma", "Lise için Madde 171/2 Sene Sonu Puan İade Çizelgesi, Ortaokul için Yaptırım Kaldırma Raporu.", theme.ViewRefreshIcon(), false, "restoration"),
		createNavCard("📊 Ceza İstatistikleri Raporu", "1. ve 2. Dönem ceza dağılımları ve okul disiplin analiz raporları.", theme.GridIcon(), false, "stats"),
		createNavCard("📚 Dizi Pusulası (Dosya İndeksi)", "Soruşturma dosyasında yer alan evrakların resmi dizi pusulası.", theme.FolderNewIcon(), false, "dizi"),
		createNavCard("🗂️ Yönetmelik & Şablon Düzenle", "Yeni suç ve ceza maddeleri ekleme, mevcut yönetmelik maddelerini düzenleme.", theme.SettingsIcon(), false, "yonetmelik"),
		createNavCard("🎨 Ayarlar & Tema Tercihleri", "Modern/Klasik arayüz seçimi, renk teması ve yazı boyutu.", theme.ColorPaletteIcon(), false, "settings"),
	)

	statsCard := buildLiveStatsCard(state, false)
	columnsGrid := container.NewGridWithColumns(3, col1, col2, col3)
	mainBody := container.NewBorder(nil, nil, nil, statsCard, container.NewVScroll(columnsGrid))

	return container.NewBorder(
		header,
		nil, nil, nil,
		mainBody,
	)
}

// -------------------------------------------------------------
// 2. 1:1 LEGACY DASHBOARD (Authentic WinForms Replica)
// -------------------------------------------------------------

func buildLegacyHub(state *AppState) fyne.CanvasObject {
	// 1. Left Brand Column: "OkulEvrak" logo & utilities
	brandLogo := canvas.NewText("OkulEvrak", color.NRGBA{R: 2, G: 119, B: 189, A: 255})
	brandLogo.TextSize = 28
	brandLogo.TextStyle = fyne.TextStyle{Bold: true, Italic: true}
	brandLogo.Alignment = fyne.TextAlignCenter

	yearText := canvas.NewText(state.AcademicYear, color.NRGBA{R: 183, G: 149, B: 11, A: 255})
	yearText.TextSize = 22
	yearText.TextStyle = fyne.TextStyle{Bold: true}
	yearText.Alignment = fyne.TextAlignCenter

	leftBtn := func(label string, icon fyne.Resource, action func()) fyne.CanvasObject {
		btn := widget.NewButtonWithIcon(label, icon, action)
		btn.Alignment = widget.ButtonAlignLeading
		return btn
	}

	backupBtn := leftBtn("Yedekleme", theme.DownloadIcon(), func() {
		showInfo(state.Window, "Yedekleme", "Veritabanı OkulEvrak.db dosyanız güvenle korunmaktadır.")
	})
	restoreBtn := leftBtn("Geri Yükleme", theme.UploadIcon(), func() {
		showInfo(state.Window, "Geri Yükleme", "OkulEvrak.db dosyanız doğrudan uygulama klasöründe yönetilmektedir.")
	})
	modernSwitchBtn := leftBtn("Modern Arayüze Geç", theme.ViewRefreshIcon(), func() {
		state.SetUIMode("modern")
	})
	settingsBtn := leftBtn("Ayarlar & Tema", theme.SettingsIcon(), func() {
		state.NavigateTo("settings")
	})
	helpBtn := leftBtn("Yardım", theme.HelpIcon(), func() {
		ShowHelpDialog(state)
	})

	leftCol := container.NewVBox(
		brandLogo,
		yearText,
		widget.NewSeparator(),
		backupBtn,
		restoreBtn,
		modernSwitchBtn,
		settingsBtn,
		helpBtn,
	)

	// 2. Center 3 Columns with Vertical Tile Buttons (Icon on top, centered text below)
	createLegacyTile := func(title string, icon fyne.Resource, isPrimary bool, targetView string) fyne.CanvasObject {
		iconWidget := widget.NewIcon(icon)
		var label *widget.Label
		if isPrimary {
			label = widget.NewLabelWithStyle(title, fyne.TextAlignCenter, fyne.TextStyle{Bold: true})
		} else {
			label = widget.NewLabelWithStyle(title, fyne.TextAlignCenter, fyne.TextStyle{})
		}
		label.Wrapping = fyne.TextWrapWord

		btn := widget.NewButton("", func() {
			state.NavigateTo(targetView)
		})
		if isPrimary {
			btn.Importance = widget.HighImportance
		}

		content := container.NewVBox(
			container.NewCenter(iconWidget),
			label,
		)
		return container.NewStack(btn, container.NewPadded(content))
	}

	createLegacyGroupBox := func(title string, tiles ...fyne.CanvasObject) fyne.CanvasObject {
		header := widget.NewLabelWithStyle(title, fyne.TextAlignCenter, fyne.TextStyle{Bold: true})
		contentItems := []fyne.CanvasObject{
			header,
			widget.NewSeparator(),
		}
		contentItems = append(contentItems, tiles...)
		return widget.NewCard("", "", container.NewVBox(contentItems...))
	}

	// GİRİŞLER GroupBox
	boxGirisler := createLegacyGroupBox("GİRİŞLER",
		createLegacyTile("Genel Bilgiler", theme.HomeIcon(), false, "school"),
		createLegacyTile("Öğrenciler", theme.AccountIcon(), false, "students"),
		createLegacyTile("Olay Girişi", theme.WarningIcon(), false, "incidents"),
		createLegacyTile("Olaya Öğrenci Ekleme", theme.ContentAddIcon(), false, "incidents"),
		createLegacyTile("YAZILI UYARI", theme.MailComposeIcon(), false, "written_warning"),
	)

	// KARAR GroupBox
	boxKarar := createLegacyGroupBox("KARAR",
		createLegacyTile("Çağrı Pusulaları", theme.MailSendIcon(), false, "summons"),
		createLegacyTile("Tutanak ve Formlar", theme.FileTextIcon(), false, "forms"),
		createLegacyTile("K A R A R", theme.ConfirmIcon(), true, "verdict"),
		createLegacyTile("Ceza Tebligatı", theme.FolderIcon(), false, "notices"),
		createLegacyTile("Sene Başı / Sene Sonu\nKurul Tutanakları", theme.ContentPasteIcon(), false, "minutes"),
	)

	// DİĞER GroupBox
	boxDiger := createLegacyGroupBox("DİĞER",
		createLegacyTile("Davranış Puanı İade Girişi", theme.ViewRefreshIcon(), false, "restoration"),
		createLegacyTile("Davranış Puanı İade Çizelgesi", theme.DocumentPrintIcon(), false, "restoration"),
		createLegacyTile("İstatistik (Ceza / Olay)", theme.GridIcon(), false, "stats"),
		createLegacyTile("Ceza Alan Öğrenciler", theme.ListIcon(), false, "stats"),
		createLegacyTile("Yönetmelik Maddeleri", theme.SettingsIcon(), false, "yonetmelik"),
		createLegacyTile("Dizi Pusulası", theme.FolderNewIcon(), false, "dizi"),
	)

	centerGrid := container.NewGridWithColumns(3,
		container.NewPadded(boxGirisler),
		container.NewPadded(boxKarar),
		container.NewPadded(boxDiger),
	)

	centerContainer := container.NewVScroll(centerGrid)

	// 3. Right Column: CEZA ALAN ÖĞRENCİ SAYILARI (Guaranteed fixed minimum width)
	rightCol := buildLiveStatsCard(state, true)

	// Top Bar in Legacy Mode: "OKUL EVRAK" Title & Quick Actions
	topBar := container.NewBorder(
		nil, nil,
		widget.NewLabelWithStyle("OKUL EVRAK", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		container.NewHBox(
			widget.NewButtonWithIcon("Öğretim Yılı: "+state.AcademicYear, theme.HistoryIcon(), func() {
				ShowAcademicYearPicker(state)
			}),
			widget.NewSeparator(),
			widget.NewButtonWithIcon("Çıktılar", theme.FolderOpenIcon(), func() {
				_ = pdf.OpenPDF(".")
			}),
			widget.NewButtonWithIcon("Modern Görünüm", theme.ViewRefreshIcon(), func() {
				state.SetUIMode("modern")
			}),
		),
		nil,
	)

	return container.NewBorder(
		container.NewVBox(topBar, widget.NewSeparator()),
		nil,
		container.NewPadded(leftCol),
		rightCol,
		centerContainer,
	)
}

type fixedStatsLayout struct {
	minWidth float32
}

func (l *fixedStatsLayout) Layout(objects []fyne.CanvasObject, size fyne.Size) {
	for _, o := range objects {
		o.Resize(size)
		o.Move(fyne.NewPos(0, 0))
	}
}

func (l *fixedStatsLayout) MinSize(objects []fyne.CanvasObject) fyne.Size {
	var maxH float32
	var maxW float32 = l.minWidth
	for _, o := range objects {
		m := o.MinSize()
		if m.Height > maxH {
			maxH = m.Height
		}
		if m.Width > maxW {
			maxW = m.Width
		}
	}
	return fyne.NewSize(maxW, maxH)
}

// buildLiveStatsCard builds the real-time penalty counter widget with guaranteed width
func buildLiveStatsCard(state *AppState, isLegacy bool) fyne.CanvasObject {
	penalized, _ := state.Store.GetPenalizedStudents(state.AcademicYear, "", "")

	sozlu := 0
	yazili := 0
	toplum := 0
	kinama := 0
	uzaklastirma := 0
	okulDeg := 0
	orgunDis := 0

	for _, p := range penalized {
		u := strings.ToUpper(p.ImposedPenalty)
		if strings.Contains(u, "SÖZLÜ") {
			sozlu++
		} else if strings.Contains(u, "YAZILI") {
			yazili++
		} else if strings.Contains(u, "TOPLUM") {
			toplum++
		} else if strings.Contains(u, "KINAMA") {
			kinama++
		} else if strings.Contains(u, "UZAKLAŞTIRMA") {
			uzaklastirma++
		} else if strings.Contains(u, "OKUL DEĞİŞTİRME") {
			okulDeg++
		} else if strings.Contains(u, "ÖRGÜN EĞİTİM") {
			orgunDis++
		}
	}

	headerText := widget.NewLabelWithStyle("CEZA ALAN ÖĞRENCİ SAYILARI", fyne.TextAlignCenter, fyne.TextStyle{Bold: true})

	badge := canvas.NewText("CEZA", color.NRGBA{R: 2, G: 119, B: 189, A: 255})
	badge.TextSize = 13
	badge.TextStyle = fyne.TextStyle{Bold: true}
	badge.Alignment = fyne.TextAlignCenter

	statRow := func(name string, count int) fyne.CanvasObject {
		lbl := widget.NewLabelWithStyle(name, fyne.TextAlignLeading, fyne.TextStyle{})
		num := widget.NewLabelWithStyle(fmt.Sprintf("%d", count), fyne.TextAlignTrailing, fyne.TextStyle{Bold: true})
		return container.NewBorder(nil, nil, lbl, num, nil)
	}

	items := []fyne.CanvasObject{
		headerText,
		container.NewCenter(badge),
		widget.NewSeparator(),
		statRow("SÖZLÜ OLARAK UYARMA", sozlu),
		statRow("YAZILI OLARAK UYARMA", yazili),
		statRow("TOPLUM HİZMETİ", toplum),
		statRow("KINAMA", kinama),
		statRow("OKULDAN KISA SÜRELİ UZAKLAŞTIRMA", uzaklastirma),
		statRow("OKUL DEĞİŞTİRME", okulDeg),
		statRow("ÖRGÜN EĞİTİM DIŞINA ÇIKARMA", orgunDis),
		widget.NewSeparator(),
		statRow("TOPLAM CEZA", len(penalized)),
	}

	card := container.NewPadded(widget.NewCard("", "", container.NewVBox(items...)))
	return container.New(&fixedStatsLayout{minWidth: 260}, card)
}
