package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"turandisiplin/model"
)

// BuildSchoolInfoTab creates the School and Board settings UI
func BuildSchoolInfoTab(state *AppState) fyne.CanvasObject {
	info := state.SchoolInfo

	// School General Details
	schoolNameEntry := NewSingleLineEntry()
	schoolNameEntry.SetText(info.SchoolName)

	instCodeEntry := NewSingleLineEntry()
	instCodeEntry.SetText(info.InstitutionCode)

	provEntry := NewSingleLineEntry()
	provEntry.SetText(info.Province)

	distEntry := NewSingleLineEntry()
	distEntry.SetText(info.District)

	principalEntry := NewSingleLineEntry()
	principalEntry.SetText(info.PrincipalName)

	letterheadEntry := NewSmartMultiLineEntry()
	letterheadEntry.SetText(info.Letterhead)
	letterheadEntry.SetPlaceHolder("Boş bırakılırsa il/ilçe ve okul adından otomatik oluşturulur.")

	officialCodeEntry := NewSingleLineEntry()
	officialCodeEntry.SetText(info.OfficialCode)
	officialCodeEntry.SetPlaceHolder("Örn: E-12345678-020-")

	phoneEntry := NewSingleLineEntry()
	phoneEntry.SetText(info.Phone)

	faxEntry := NewSingleLineEntry()
	faxEntry.SetText(info.Fax)

	addressEntry := NewSingleLineEntry()
	addressEntry.SetText(info.Address)

	// Primary Board Members
	chairName := NewSingleLineEntry()
	chairName.SetText(info.BoardChair)
	chairTitle := NewSingleLineEntry()
	chairTitle.SetText(info.BoardChairTitle)

	m1Name := NewSingleLineEntry()
	m1Name.SetText(info.Member1)
	m1Title := NewSingleLineEntry()
	m1Title.SetText(info.Member1Title)

	m2Name := NewSingleLineEntry()
	m2Name.SetText(info.Member2)
	m2Title := NewSingleLineEntry()
	m2Title.SetText(info.Member2Title)

	m3Name := NewSingleLineEntry()
	m3Name.SetText(info.Member3)
	m3Title := NewSingleLineEntry()
	m3Title.SetText(info.Member3Title)

	m4Name := NewSingleLineEntry()
	m4Name.SetText(info.Member4)
	m4Title := NewSingleLineEntry()
	m4Title.SetText(info.Member4Title)

	// Substitute Board Members
	subChairName := NewSingleLineEntry()
	subChairName.SetText(info.SubBoardChair)
	subChairTitle := NewSingleLineEntry()
	subChairTitle.SetText(info.SubBoardChairTitle)

	subM1Name := NewSingleLineEntry()
	subM1Name.SetText(info.SubMember1)
	subM1Title := NewSingleLineEntry()
	subM1Title.SetText(info.SubMember1Title)

	subM2Name := NewSingleLineEntry()
	subM2Name.SetText(info.SubMember2)
	subM2Title := NewSingleLineEntry()
	subM2Title.SetText(info.SubMember2Title)

	subM3Name := NewSingleLineEntry()
	subM3Name.SetText(info.SubMember3)
	subM3Title := NewSingleLineEntry()
	subM3Title.SetText(info.SubMember3Title)

	subM4Name := NewSingleLineEntry()
	subM4Name.SetText(info.SubMember4)
	subM4Title := NewSingleLineEntry()
	subM4Title.SetText(info.SubMember4Title)


	saveBtn := widget.NewButtonWithIcon("Bilgileri Kaydet", theme.DocumentSaveIcon(), func() {
		updated := &model.SchoolInfo{
			ID:                 info.ID,
			SchoolType:         state.SchoolInfo.SchoolType,
			AcademicYear:       state.AcademicYear,
			InstitutionCode:    instCodeEntry.Text,
			Province:           provEntry.Text,
			District:           distEntry.Text,
			SchoolName:         schoolNameEntry.Text,
			PrincipalName:      principalEntry.Text,
			Letterhead:         letterheadEntry.Text,
			OfficialCode:       officialCodeEntry.Text,
			Address:            addressEntry.Text,
			Phone:              phoneEntry.Text,
			Fax:                faxEntry.Text,
			BoardChair:         chairName.Text,
			BoardChairTitle:    chairTitle.Text,
			Member1:            m1Name.Text,
			Member1Title:       m1Title.Text,
			Member2:            m2Name.Text,
			Member2Title:       m2Title.Text,
			Member3:            m3Name.Text,
			Member3Title:       m3Title.Text,
			Member4:            m4Name.Text,
			Member4Title:       m4Title.Text,
			SubBoardChair:      subChairName.Text,
			SubBoardChairTitle: subChairTitle.Text,
			SubMember1:         subM1Name.Text,
			SubMember1Title:    subM1Title.Text,
			SubMember2:         subM2Name.Text,
			SubMember2Title:    subM2Title.Text,
			SubMember3:         subM3Name.Text,
			SubMember3Title:    subM3Title.Text,
			SubMember4:         subM4Name.Text,
			SubMember4Title:    subM4Title.Text,
		}

		err := state.Store.SaveSchoolInfo(updated)
		if err != nil {
			showError(state.Window, err)
			return
		}
		state.SchoolInfo = updated
		showInfo(state.Window, "Başarılı", "Okul ve kurul bilgileri başarıyla kaydedildi.")
		state.NotifyChange()
	})
	saveBtn.Importance = widget.HighImportance

	// Forms
	schoolForm := widget.NewForm(
		widget.NewFormItem("Okul Adı", schoolNameEntry),
		widget.NewFormItem("Kurum Kodu", instCodeEntry),
		widget.NewFormItem("İl", provEntry),
		widget.NewFormItem("İlçe", distEntry),
		widget.NewFormItem("Okul Müdürü Adı", principalEntry),
		widget.NewFormItem("Resmi Yazı Kodu", officialCodeEntry),
		widget.NewFormItem("Telefon", phoneEntry),
		widget.NewFormItem("Faks", faxEntry),
		widget.NewFormItem("Adres", addressEntry),
		widget.NewFormItem("Özel Antet (İsteğe Bağlı)", letterheadEntry),
	)

	primaryForm := widget.NewForm(
		widget.NewFormItem("Kurul Başkanı Adı", chairName),
		widget.NewFormItem("Başkan Ünvanı", chairTitle),
		widget.NewFormItem("Asil Üye 1 Adı", m1Name),
		widget.NewFormItem("Üye 1 Ünvanı", m1Title),
		widget.NewFormItem("Asil Üye 2 Adı", m2Name),
		widget.NewFormItem("Üye 2 Ünvanı", m2Title),
		widget.NewFormItem("Asil Üye 3 Adı", m3Name),
		widget.NewFormItem("Üye 3 Ünvanı", m3Title),
		widget.NewFormItem("Asil Üye 4 Adı", m4Name),
		widget.NewFormItem("Üye 4 Ünvanı", m4Title),
	)

	subForm := widget.NewForm(
		widget.NewFormItem("Yedek Kurul Başkanı Adı", subChairName),
		widget.NewFormItem("Yedek Başkan Ünvanı", subChairTitle),
		widget.NewFormItem("Yedek Üye 1 Adı", subM1Name),
		widget.NewFormItem("Yedek Üye 1 Ünvanı", subM1Title),
		widget.NewFormItem("Yedek Üye 2 Adı", subM2Name),
		widget.NewFormItem("Yedek Üye 2 Ünvanı", subM2Title),
		widget.NewFormItem("Yedek Üye 3 Adı", subM3Name),
		widget.NewFormItem("Yedek Üye 3 Ünvanı", subM3Title),
		widget.NewFormItem("Yedek Üye 4 Adı", subM4Name),
		widget.NewFormItem("Yedek Üye 4 Ünvanı", subM4Title),
	)

	cardSchool := widget.NewCard("Kurum Bilgileri", "Resmi evraklarda ve antetlerde kullanılacak kurum bilgileri", schoolForm)
	cardPrimary := widget.NewCard("Asil Kurul Üyeleri", "Karar tutanakları ve toplantı çağrılarında imzası bulunan üyeler", primaryForm)
	cardSub := widget.NewCard("Yedek Kurul Üyeleri", "Asil üyenin olaya taraf olması veya mazereti durumunda yerine görev alacak üyeler", subForm)

	topBar := container.NewHBox(
		widget.NewLabelWithStyle("🏛️ Okul ve Disiplin / Davranış Değerlendirme Kurulu Bilgileri", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		saveBtn,
	)

	content := container.NewVBox(
		topBar,
		cardSchool,
		cardPrimary,
		cardSub,
		saveBtn,
	)

	pageScroll := container.NewVScroll(content)
	letterheadEntry.ParentScroll = pageScroll
	return pageScroll
}
