package pdf

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/signintech/gopdf"

	"okulevrak/assets"
	"okulevrak/model"
)

const (
	FontFamily = "LiberationSans"

	PageWidthPortrait   = 595.28
	PageHeightPortrait  = 841.89
	PageWidthLandscape  = 841.89
	PageHeightLandscape = 595.28

	DefaultMarginLeft   = 30.0
	DefaultMarginRight  = 20.0
	DefaultMarginTop    = 25.0
	DefaultMarginBottom = 25.0
)

// DocWrapper wraps a gopdf.GoPdf instance with helper methods.
type DocWrapper struct {
	PDF         *gopdf.GoPdf
	Landscape   bool
	MarginLeft  float64
	MarginRight float64
	MarginTop   float64
	MarginBottom float64
	UsableWidth float64
}

// NewDoc creates and initializes a new PDF document with embedded Turkish fonts.
func NewDoc(landscape bool) (*DocWrapper, error) {
	pdf := &gopdf.GoPdf{}
	config := gopdf.Config{}
	if landscape {
		config.PageSize = *gopdf.PageSizeA4Landscape
	} else {
		config.PageSize = *gopdf.PageSizeA4
	}
	pdf.Start(config)

	// Add embedded LiberationSans fonts
	err := pdf.AddTTFFontDataWithOption(FontFamily, assets.FontRegular, gopdf.TtfOption{Style: gopdf.Regular})
	if err != nil {
		return nil, fmt.Errorf("regular font yüklenemedi: %w", err)
	}
	err = pdf.AddTTFFontDataWithOption(FontFamily, assets.FontBold, gopdf.TtfOption{Style: gopdf.Bold})
	if err != nil {
		return nil, fmt.Errorf("bold font yüklenemedi: %w", err)
	}
	err = pdf.AddTTFFontDataWithOption(FontFamily, assets.FontItalic, gopdf.TtfOption{Style: gopdf.Italic})
	if err != nil {
		return nil, fmt.Errorf("italic font yüklenemedi: %w", err)
	}

	ml := DefaultMarginLeft
	mr := DefaultMarginRight
	if landscape {
		ml = 20.0
		mr = 20.0
	}
	pw := PageWidthPortrait
	if landscape {
		pw = PageWidthLandscape
	}

	doc := &DocWrapper{
		PDF:          pdf,
		Landscape:    landscape,
		MarginLeft:   ml,
		MarginRight:  mr,
		MarginTop:    DefaultMarginTop,
		MarginBottom: DefaultMarginBottom,
		UsableWidth:  pw - ml - mr,
	}

	pdf.SetMargins(ml, doc.MarginTop, mr, doc.MarginBottom)
	return doc, nil
}

// SetRegularFont sets the font to regular at the given size.
func (d *DocWrapper) SetRegularFont(size float64) {
	_ = d.PDF.SetFont(FontFamily, "", size)
}

// SetBoldFont sets the font to bold at the given size.
func (d *DocWrapper) SetBoldFont(size float64) {
	_ = d.PDF.SetFont(FontFamily, "B", size)
}

// SetItalicFont sets the font to italic at the given size.
func (d *DocWrapper) SetItalicFont(size float64) {
	_ = d.PDF.SetFont(FontFamily, "I", size)
}

// DrawAntet draws the official school header (TC, MEB, Valilik/Kaymakamlık, Okul Adı)
func (d *DocWrapper) DrawAntet(school *model.SchoolInfo) {
	d.SetBoldFont(9)
	lines := []string{
		"T.C.",
	}

	prov := strings.TrimSpace(school.Province)
	dist := strings.TrimSpace(school.District)
	if dist != "" && prov != "" {
		lines = append(lines, fmt.Sprintf("%s KAYMAKAMLIĞI / İLÇE MİLLİ EĞİTİM MÜDÜRLÜĞÜ", strings.ToUpper(dist)))
	} else if prov != "" {
		lines = append(lines, fmt.Sprintf("%s VALİLİĞİ / İL MİLLİ EĞİTİM MÜDÜRLÜĞÜ", strings.ToUpper(prov)))
	}

	if school.Letterhead != "" {
		// Use custom letterhead if present
		l := strings.Split(school.Letterhead, "\n")
		for _, s := range l {
			s = strings.TrimSpace(s)
			if s != "" {
				lines = append(lines, s)
			}
		}
	} else if school.SchoolName != "" {
		lines = append(lines, strings.ToUpper(school.SchoolName))
	}

	for _, line := range lines {
		d.PDF.SetX(d.MarginLeft)
		_ = d.PDF.CellWithOption(&gopdf.Rect{W: d.UsableWidth, H: 14}, line, gopdf.CellOption{
			Align: gopdf.Center,
		})
		d.PDF.Br(14)
	}
	d.PDF.Br(4)
}

// DrawTitle draws a centered bold title.
func (d *DocWrapper) DrawTitle(title string, size float64, spacingAfter float64) {
	d.SetBoldFont(size)
	d.PDF.SetX(d.MarginLeft)
	_ = d.PDF.CellWithOption(&gopdf.Rect{W: d.UsableWidth, H: size * 1.5}, title, gopdf.CellOption{
		Align: gopdf.Center,
	})
	d.PDF.Br(size*1.5 + spacingAfter)
}

// DrawRect draws a border rectangle from top-left.
func (d *DocWrapper) DrawRect(x, y, w, h float64) {
	d.PDF.RectFromUpperLeftWithStyle(x, y, w, h, "D")
}

// DrawRuledLines draws n horizontal ruled lines for handwritten statements.
func (d *DocWrapper) DrawRuledLines(count int, lineHeight float64) {
	d.PDF.SetLineWidth(0.4)
	for i := 0; i < count; i++ {
		y := d.PDF.GetY() + lineHeight
		d.PDF.Line(d.MarginLeft, y, d.MarginLeft+d.UsableWidth, y)
		d.PDF.SetY(y)
	}
}

// DrawDottedLines draws n lines of dots across the usable width.
func (d *DocWrapper) DrawDottedLines(count int, lineHeight float64) {
	d.SetRegularFont(8)
	dots := strings.Repeat(". ", int(d.UsableWidth/6))
	for i := 0; i < count; i++ {
		d.PDF.SetX(d.MarginLeft)
		_ = d.PDF.CellWithOption(&gopdf.Rect{W: d.UsableWidth, H: lineHeight}, dots, gopdf.CellOption{Align: gopdf.Left})
		d.PDF.Br(lineHeight)
	}
}

// DrawSignatures5 draws the standard 5-column board member signatures table.
func (d *DocWrapper) DrawSignatures5(school *model.SchoolInfo, titleLabel string) {
	d.PDF.Br(15)
	colW := d.UsableWidth / 5.0
	startY := d.PDF.GetY()

	members := []struct {
		Role  string
		Name  string
		Title string
	}{
		{Role: "Kurul Başkanı", Name: school.BoardChair, Title: school.BoardChairTitle},
		{Role: "Üye", Name: school.Member1, Title: school.Member1Title},
		{Role: "Üye", Name: school.Member2, Title: school.Member2Title},
		{Role: "Üye", Name: school.Member3, Title: school.Member3Title},
		{Role: "Üye", Name: school.Member4, Title: school.Member4Title},
	}

	// Roles
	d.SetBoldFont(7.5)
	for i, m := range members {
		d.PDF.SetX(d.MarginLeft + float64(i)*colW)
		d.PDF.SetY(startY)
		_ = d.PDF.CellWithOption(&gopdf.Rect{W: colW, H: 12}, m.Role, gopdf.CellOption{Align: gopdf.Center})
	}

	// Names & Titles
	d.SetRegularFont(7.0)
	for i, m := range members {
		nameStr := m.Name
		if nameStr == "" {
			nameStr = "........................"
		}
		titleStr := m.Title

		d.PDF.SetX(d.MarginLeft + float64(i)*colW)
		d.PDF.SetY(startY + 32)
		_ = d.PDF.CellWithOption(&gopdf.Rect{W: colW, H: 10}, nameStr, gopdf.CellOption{Align: gopdf.Center})

		d.PDF.SetX(d.MarginLeft + float64(i)*colW)
		d.PDF.SetY(startY + 43)
		_ = d.PDF.CellWithOption(&gopdf.Rect{W: colW, H: 10}, titleStr, gopdf.CellOption{Align: gopdf.Center})
	}

	d.PDF.SetY(startY + 60)
}

// DrawPrincipalApproval draws the centered "UYGUNDUR / Okul Müdürü" block.
func (d *DocWrapper) DrawPrincipalApproval(principalName string) {
	d.PDF.Br(12)
	d.SetBoldFont(8)
	d.PDF.SetX(d.MarginLeft)
	dateStr := time.Now().Format("02/01/2006")
	_ = d.PDF.CellWithOption(&gopdf.Rect{W: d.UsableWidth, H: 12}, dateStr, gopdf.CellOption{Align: gopdf.Center})
	d.PDF.Br(12)
	_ = d.PDF.CellWithOption(&gopdf.Rect{W: d.UsableWidth, H: 12}, "UYGUNDUR", gopdf.CellOption{Align: gopdf.Center})
	d.PDF.Br(28)

	pName := strings.TrimSpace(principalName)
	if pName == "" {
		pName = "........................................"
	}
	_ = d.PDF.CellWithOption(&gopdf.Rect{W: d.UsableWidth, H: 12}, pName, gopdf.CellOption{Align: gopdf.Center})
	d.PDF.Br(12)
	d.SetRegularFont(8)
	_ = d.PDF.CellWithOption(&gopdf.Rect{W: d.UsableWidth, H: 12}, "Okul Müdürü", gopdf.CellOption{Align: gopdf.Center})
	d.PDF.Br(12)
}

// DrawFooterAddress draws bottom address, phone, and fax.
func (d *DocWrapper) DrawFooterAddress(school *model.SchoolInfo) {
	pageH := PageHeightPortrait
	if d.Landscape {
		pageH = PageHeightLandscape
	}
	footerY := pageH - d.MarginBottom - 25

	d.PDF.SetLineWidth(0.8)
	d.PDF.Line(d.MarginLeft, footerY, d.MarginLeft+d.UsableWidth, footerY)

	d.SetItalicFont(6.5)
	d.PDF.SetX(d.MarginLeft)
	d.PDF.SetY(footerY + 3)

	contact := fmt.Sprintf("Adres: %s   |   Tel: %s   |   Fax: %s", school.Address, school.Phone, school.Fax)
	_ = d.PDF.CellWithOption(&gopdf.Rect{W: d.UsableWidth, H: 10}, contact, gopdf.CellOption{Align: gopdf.Center})
}

// GenerateUniquePDFPath returns an available file path avoiding collision.
func GenerateUniquePDFPath(baseDir string, prefix string) (string, error) {
	if baseDir == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			baseDir = "."
		} else {
			baseDir = filepath.Join(home, "OkulEvrak_belgeler")
			// Migrate legacy folder if exists and new folder does not
			oldDir := filepath.Join(home, "OkulEvrak_Ciktilar")
			if _, errOld := os.Stat(oldDir); errOld == nil {
				if _, errNew := os.Stat(baseDir); os.IsNotExist(errNew) {
					_ = os.Rename(oldDir, baseDir)
				}
			}
		}
	}

	_ = os.MkdirAll(baseDir, 0755)

	safePrefix := strings.ReplaceAll(prefix, " ", "_")
	safePrefix = strings.ReplaceAll(safePrefix, "/", "-")
	safePrefix = strings.ReplaceAll(safePrefix, "\\", "-")

	for i := 1; i <= 9999; i++ {
		candidate := filepath.Join(baseDir, fmt.Sprintf("%s_%d.pdf", safePrefix, i))
		if _, err := os.Stat(candidate); os.IsNotExist(err) {
			return candidate, nil
		}
	}

	return filepath.Join(baseDir, fmt.Sprintf("%s_%d.pdf", safePrefix, time.Now().Unix())), nil
}

// OpenPDF opens the generated PDF with the default system application.
func OpenPDF(pdfPath string) error {
	if _, err := os.Stat(pdfPath); err != nil {
		return fmt.Errorf("dosya bulunamadı: %w", err)
	}

	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("rundll32", "url.dll,FileProtocolHandler", pdfPath)
	case "darwin":
		cmd = exec.Command("open", pdfPath)
	default:
		cmd = exec.Command("xdg-open", pdfPath)
	}

	err := cmd.Start()
	if err != nil {
		return fmt.Errorf("PDF açılamadı: %w", err)
	}
	return nil
}
