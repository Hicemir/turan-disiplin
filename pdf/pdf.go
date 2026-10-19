package pdf

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"runtime/debug"
	"strings"
	"time"

	"github.com/signintech/gopdf"

	"turandisiplin/assets"
	"turandisiplin/model"
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

// SplitLines splits text by newlines and wraps each paragraph to fit within maxWidth using the currently active font.
func (d *DocWrapper) SplitLines(text string, maxWidth float64) []string {
	text = strings.ReplaceAll(text, "\r\n", "\n")
	text = strings.TrimSpace(text)
	if text == "" {
		return []string{"—"}
	}
	if maxWidth <= 10 {
		return []string{text}
	}
	paragraphs := strings.Split(text, "\n")
	var out []string
	for _, p := range paragraphs {
		p = strings.TrimSpace(p)
		if p == "" {
			continue
		}
		wrapped, err := d.PDF.SplitTextWithWordWrap(p, maxWidth)
		if err != nil || len(wrapped) == 0 {
			out = append(out, p)
		} else {
			out = append(out, wrapped...)
		}
	}
	if len(out) == 0 {
		return []string{"—"}
	}
	return out
}

// DrawWrappedText word-wraps text within width at horizontal offset x, advancing Y by lineH per line.
func (d *DocWrapper) DrawWrappedText(x, width, lineH float64, text string, align int) {
	if strings.TrimSpace(text) == "" {
		return
	}
	pageH := PageHeightPortrait
	if d.Landscape {
		pageH = PageHeightLandscape
	}
	lines := d.SplitLines(text, width)
	for _, line := range lines {
		if d.PDF.GetY()+lineH > pageH-d.MarginBottom-25 {
			d.PDF.AddPage()
			d.PDF.SetY(d.MarginTop)
		}
		d.PDF.SetX(x)
		_ = d.PDF.CellWithOption(&gopdf.Rect{W: width, H: lineH}, line, gopdf.CellOption{Align: align})
		d.PDF.Br(lineH)
	}
}

// DrawWrappedCellLines renders pre-wrapped lines inside a table cell starting at (x, startY).
func (d *DocWrapper) DrawWrappedCellLines(x, startY, width, lineH float64, lines []string, align int) {
	for i, line := range lines {
		d.PDF.SetX(x)
		d.PDF.SetY(startY + float64(i)*lineH)
		_ = d.PDF.CellWithOption(&gopdf.Rect{W: width, H: lineH}, line, gopdf.CellOption{Align: align})
	}
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
		wrapped := d.SplitLines(line, d.UsableWidth)
		for _, wLine := range wrapped {
			d.PDF.SetX(d.MarginLeft)
			_ = d.PDF.CellWithOption(&gopdf.Rect{W: d.UsableWidth, H: 13}, wLine, gopdf.CellOption{
				Align: gopdf.Center,
			})
			d.PDF.Br(13)
		}
	}
	d.PDF.Br(4)
}

// DrawTitle draws a centered bold title (supports multi-line \n and automatic word-wrapping).
func (d *DocWrapper) DrawTitle(title string, size float64, spacingAfter float64) {
	d.SetBoldFont(size)
	lineH := size * 1.35
	wrapped := d.SplitLines(title, d.UsableWidth)
	for _, line := range wrapped {
		d.PDF.SetX(d.MarginLeft)
		_ = d.PDF.CellWithOption(&gopdf.Rect{W: d.UsableWidth, H: lineH}, line, gopdf.CellOption{
			Align: gopdf.Center,
		})
		d.PDF.Br(lineH)
	}
	if spacingAfter > 0 {
		d.PDF.Br(spacingAfter)
	}
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
	maxEndY := startY + 55
	for i, m := range members {
		nameStr := strings.TrimSpace(m.Name)
		if nameStr == "" {
			nameStr = "........................"
		}
		titleStr := strings.TrimSpace(m.Title)

		colX := d.MarginLeft + float64(i)*colW + 2
		nameLines := d.SplitLines(nameStr, colW-4)
		d.DrawWrappedCellLines(colX, startY+30, colW-4, 9.5, nameLines, gopdf.Center)

		titleY := startY + 30 + float64(len(nameLines))*9.5 + 1
		if titleStr != "" {
			titleLines := d.SplitLines(titleStr, colW-4)
			d.DrawWrappedCellLines(colX, titleY, colW-4, 9.5, titleLines, gopdf.Center)
			endY := titleY + float64(len(titleLines))*9.5
			if endY > maxEndY {
				maxEndY = endY
			}
		}
	}

	d.PDF.SetY(maxEndY + 6)
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
	var parts []string
	if strings.TrimSpace(school.Address) != "" {
		parts = append(parts, "Adres: "+strings.TrimSpace(school.Address))
	}
	if strings.TrimSpace(school.Phone) != "" {
		parts = append(parts, "Tel: "+strings.TrimSpace(school.Phone))
	}
	if strings.TrimSpace(school.Fax) != "" {
		parts = append(parts, "Fax: "+strings.TrimSpace(school.Fax))
	}
	contact := strings.Join(parts, "   |   ")
	if contact != "" {
		lines := d.SplitLines(contact, d.UsableWidth)
		d.DrawWrappedCellLines(d.MarginLeft, footerY+3, d.UsableWidth, 9.0, lines, gopdf.Center)
	}
}

// GenerateUniquePDFPath returns an available file path avoiding collision.
func GenerateUniquePDFPath(baseDir string, prefix string) (string, error) {
	if baseDir == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			baseDir = "."
		} else {
			baseDir = filepath.Join(home, "TuranDisiplin_belgeler")
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
	// The PDF has just been generated; return its temporary buffers to the OS.
	go debug.FreeOSMemory()

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
