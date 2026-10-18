package excel

import (
	"bytes"
	"encoding/csv"
	"encoding/xml"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/extrame/xls"
	"github.com/xuri/excelize/v2"
	"golang.org/x/net/html"

	"turandisiplin/model"
)

// ReadSheetData reads all sheets from any Excel, HTML, XML or CSV export into map[sheetName][][]string.
// Guaranteed never to panic or crash.
func ReadSheetData(filePath string) (res map[string][][]string, sheets []string, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("dosya açılırken beklenmeyen hata: %v", r)
		}
	}()

	headerBytes := make([]byte, 4096)
	f, openErr := os.Open(filePath)
	if openErr == nil {
		_, _ = f.Read(headerBytes)
		_ = f.Close()
	}
	contentStr := strings.TrimSpace(string(headerBytes))
	lowerContent := strings.ToLower(contentStr)

	// 1. OpenXML (.xlsx / .xlsm) magic bytes: "PK\x03\x04"
	if len(headerBytes) >= 4 && headerBytes[0] == 0x50 && headerBytes[1] == 0x4B && headerBytes[2] == 0x03 && headerBytes[3] == 0x04 {
		if data, sList, e := readXLSX(filePath); e == nil && len(sList) > 0 {
			return data, sList, nil
		}
	}

	// 2. HTML Table (frequently exported by e-Okul / ActiveReports with .xls extension)
	if strings.Contains(lowerContent, "<table") || strings.Contains(lowerContent, "<html") ||
		strings.Contains(lowerContent, "<!doctype") || strings.Contains(lowerContent, "<tr") {
		if data, sList, e := readHTMLTable(filePath); e == nil && len(sList) > 0 {
			return data, sList, nil
		}
	}

	// 3. XML Spreadsheet 2003 (SpreadsheetML, e-Okul "Excel (Sadece Veri)" export with .xls extension)
	if strings.Contains(lowerContent, "<?xml") || strings.Contains(lowerContent, "<workbook") || strings.Contains(lowerContent, "xmlns:ss") {
		if data, sList, e := readXMLSpreadsheet(filePath); e == nil && len(sList) > 0 {
			return data, sList, nil
		}
	}

	// 4. Binary OLE2 BIFF8 (.xls) magic bytes: D0 CF 11 E0 A1 B1 1A E1
	if len(headerBytes) >= 8 && headerBytes[0] == 0xD0 && headerBytes[1] == 0xCF && headerBytes[2] == 0x11 && headerBytes[3] == 0xE0 {
		if data, sList, e := readXLS(filePath); e == nil && len(sList) > 0 {
			return data, sList, nil
		}
	}

	// 5. Fallback cascade: try each parser with full panic isolation
	ext := strings.ToLower(filepath.Ext(filePath))
	if ext == ".xlsx" || ext == ".xlsm" {
		if data, sList, e := readXLSX(filePath); e == nil && len(sList) > 0 {
			return data, sList, nil
		}
	}
	if data, sList, e := readHTMLTable(filePath); e == nil && len(sList) > 0 {
		return data, sList, nil
	}
	if data, sList, e := readXMLSpreadsheet(filePath); e == nil && len(sList) > 0 {
		return data, sList, nil
	}
	if data, sList, e := readXLS(filePath); e == nil && len(sList) > 0 {
		return data, sList, nil
	}
	return readDelimited(filePath)
}

func readXLSX(filePath string) (res map[string][][]string, sheets []string, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("xlsx okuma hatası: %v", r)
		}
	}()

	f, err := excelize.OpenFile(filePath)
	if err != nil {
		return nil, nil, err
	}
	defer f.Close()

	sheets = f.GetSheetList()
	if len(sheets) == 0 {
		return nil, nil, fmt.Errorf("çalışma sayfasında sayfa bulunamadı")
	}

	res = make(map[string][][]string)
	for _, sheet := range sheets {
		rows, err := f.GetRows(sheet)
		if err != nil {
			continue
		}
		res[sheet] = rows
	}
	return res, sheets, nil
}

func readXLS(filePath string) (res map[string][][]string, sheets []string, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("xls ayrıştırma hatası: %v", r)
		}
	}()

	wb, err := xls.Open(filePath, "utf-8")
	if err != nil {
		return nil, nil, err
	}

	numSheets := wb.NumSheets()
	if numSheets == 0 {
		return nil, nil, fmt.Errorf("xls dosyasında sayfa bulunamadı")
	}

	res = make(map[string][][]string)
	for i := 0; i < numSheets; i++ {
		sheet := wb.GetSheet(i)
		if sheet == nil {
			continue
		}
		sheets = append(sheets, sheet.Name)
		var rows [][]string
		maxRow := int(sheet.MaxRow)
		for r := 0; r <= maxRow; r++ {
			rowObj := sheet.Row(r)
			if rowObj == nil {
				continue
			}
			lastCol := rowObj.LastCol()
			var cols []string
			hasContent := false
			for c := 0; c <= lastCol; c++ {
				val := strings.TrimSpace(rowObj.Col(c))
				if val != "" {
					hasContent = true
				}
				cols = append(cols, val)
			}
			if hasContent {
				rows = append(rows, cols)
			}
		}
		res[sheet.Name] = rows
	}
	return res, sheets, nil
}

// readHTMLTable parses HTML-formatted reports exported by e-Okul / ActiveReports
func readHTMLTable(filePath string) (map[string][][]string, []string, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, nil, err
	}

	doc, err := html.Parse(bytes.NewReader(data))
	if err != nil {
		return nil, nil, err
	}

	var rows [][]string
	var walk func(*html.Node)
	walk = func(n *html.Node) {
		if n.Type == html.ElementNode && strings.EqualFold(n.Data, "tr") {
			var cols []string
			for c := n.FirstChild; c != nil; c = c.NextSibling {
				if c.Type == html.ElementNode && (strings.EqualFold(c.Data, "td") || strings.EqualFold(c.Data, "th")) {
					text := extractHTMLText(c)
					text = strings.ReplaceAll(text, "\u00a0", " ") // non-breaking space
					text = strings.TrimSpace(text)
					cols = append(cols, text)
				}
			}
			if len(cols) > 0 {
				hasContent := false
				for _, v := range cols {
					if v != "" {
						hasContent = true
						break
					}
				}
				if hasContent {
					rows = append(rows, cols)
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			walk(c)
		}
	}
	walk(doc)

	if len(rows) == 0 {
		return nil, nil, fmt.Errorf("HTML tablosunda satır bulunamadı")
	}

	return map[string][][]string{"Sayfa1": rows}, []string{"Sayfa1"}, nil
}

func extractHTMLText(n *html.Node) string {
	if n.Type == html.TextNode {
		return n.Data
	}
	var sb strings.Builder
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		sb.WriteString(extractHTMLText(c))
	}
	return sb.String()
}

// XML Spreadsheet 2003 (SpreadsheetML) models
type xmlWorkbook struct {
	XMLName    xml.Name       `xml:"Workbook"`
	Worksheets []xmlWorksheet `xml:"Worksheet"`
}

type xmlWorksheet struct {
	Name  string   `xml:"Name,attr"`
	Table xmlTable `xml:"Table"`
}

type xmlTable struct {
	Rows []xmlRow `xml:"Row"`
}

type xmlRow struct {
	Cells []xmlCell `xml:"Cell"`
}

type xmlCell struct {
	Data string `xml:"Data"`
}

func readXMLSpreadsheet(filePath string) (map[string][][]string, []string, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, nil, err
	}

	var wb xmlWorkbook
	if err := xml.Unmarshal(data, &wb); err != nil {
		return nil, nil, err
	}

	if len(wb.Worksheets) == 0 {
		return nil, nil, fmt.Errorf("XML çalışma sayfasında sayfa bulunamadı")
	}

	res := make(map[string][][]string)
	var sheets []string

	for idx, ws := range wb.Worksheets {
		sheetName := ws.Name
		if sheetName == "" {
			sheetName = fmt.Sprintf("Sayfa%d", idx+1)
		}
		sheets = append(sheets, sheetName)

		var rows [][]string
		for _, r := range ws.Table.Rows {
			var cols []string
			hasContent := false
			for _, c := range r.Cells {
				val := strings.TrimSpace(c.Data)
				if val != "" {
					hasContent = true
				}
				cols = append(cols, val)
			}
			if hasContent {
				rows = append(rows, cols)
			}
		}
		res[sheetName] = rows
	}

	return res, sheets, nil
}

func readDelimited(filePath string) (map[string][][]string, []string, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, nil, err
	}

	// Try delimiters: tab, semicolon, comma
	delims := []rune{'\t', ';', ','}
	for _, d := range delims {
		r := csv.NewReader(bytes.NewReader(data))
		r.Comma = d
		r.LazyQuotes = true
		r.FieldsPerRecord = -1
		rows, err := r.ReadAll()
		if err == nil && len(rows) > 0 {
			var cleanRows [][]string
			for _, row := range rows {
				hasContent := false
				var cleanCols []string
				for _, cell := range row {
					v := strings.TrimSpace(cell)
					if v != "" {
						hasContent = true
					}
					cleanCols = append(cleanCols, v)
				}
				if hasContent {
					cleanRows = append(cleanRows, cleanCols)
				}
			}
			if len(cleanRows) > 0 {
				return map[string][][]string{"Sayfa1": cleanRows}, []string{"Sayfa1"}, nil
			}
		}
	}
	return nil, nil, fmt.Errorf("dosya biçimi tanınamadı")
}

// SelectDefaultSheet chooses the most likely sheet name from e-Okul exports.
func SelectDefaultSheet(sheets []string) string {
	if len(sheets) == 0 {
		return ""
	}
	priorities := []string{"rapor", "sayfa1", "sheet1", "rapor1", "ogrenci", "öğrenci", "liste"}
	for _, p := range priorities {
		for _, s := range sheets {
			if strings.EqualFold(strings.TrimSpace(s), p) {
				return s
			}
		}
	}
	for _, p := range priorities {
		for _, s := range sheets {
			if strings.Contains(strings.ToLower(s), p) {
				return s
			}
		}
	}
	return sheets[0]
}

// ParseStudentRows parses student rows supporting tabular, card-based ("Fotoğraflı"), and heuristic formats.
func ParseStudentRows(rows [][]string, schoolYear string) (students []*model.Student, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("öğrenci verileri ayrıştırılırken hata: %v", r)
		}
	}()

	if len(rows) == 0 {
		return nil, fmt.Errorf("excel verisi boş")
	}

	seenNo := make(map[string]bool)
	currentClass := ""

	// Strategy 1: Tabular table with detected header row
	headerRowIdx := -1
	colMap := make(map[string]int)

	for i, row := range rows {
		if i > 120 {
			break
		}
		m := detectHeaderMapping(row)
		hasNo := (hasKey(m, "okulno") && m["okulno"] >= 0) || (hasKey(m, "tcno") && m["tcno"] >= 0)
		hasName := (hasKey(m, "adsoyad") && m["adsoyad"] >= 0) || (hasKey(m, "ad") && hasKey(m, "soyad"))
		if len(m) >= 2 && hasNo && hasName {
			headerRowIdx = i
			colMap = m
			break
		}
	}

	if headerRowIdx >= 0 {
		for r := headerRowIdx + 1; r < len(rows); r++ {
			row := rows[r]
			if len(row) == 0 {
				continue
			}
			banner := extractClassBanner(row)
			if banner != "" {
				currentClass = banner
				continue
			}
			st := studentFromMappedRow(row, colMap, currentClass, schoolYear)
			if st != nil && st.SchoolNo != "" && !seenNo[st.SchoolNo] {
				seenNo[st.SchoolNo] = true
				students = append(students, st)
			}
		}
		if len(students) > 0 {
			return students, nil
		}
	}

	// Strategy 2: Card / Block / Multi-row layout (e.g. e-Okul "Fotoğraflı-Kimlik Bilgili Öğrenci Listesi")
	cardStudents := parseCardBasedStudents(rows, schoolYear)
	if len(cardStudents) > 0 {
		return cardStudents, nil
	}

	// Strategy 3: Heuristic Row Parser (pure tabular data rows without headers)
	for r := 0; r < len(rows); r++ {
		row := rows[r]
		if len(row) == 0 {
			continue
		}
		banner := extractClassBanner(row)
		if banner != "" {
			currentClass = banner
			continue
		}
		st := studentFromHeuristicRow(row, currentClass, schoolYear)
		if st != nil && st.SchoolNo != "" && !seenNo[st.SchoolNo] {
			seenNo[st.SchoolNo] = true
			students = append(students, st)
		}
	}

	if len(students) == 0 {
		return nil, fmt.Errorf("geçerli öğrenci kaydı bulunamadı (okul numarası veya ad bilgileri tespit edilemedi)")
	}

	return students, nil
}

// parseCardBasedStudents handles multi-row card layouts (such as Fotoğraflı-Kimlik Bilgili Öğrenci Listesi)
func parseCardBasedStudents(rows [][]string, schoolYear string) []*model.Student {
	var list []*model.Student
	seen := make(map[string]bool)
	currentClass := ""

	var current *model.Student

	saveCurrent := func() {
		if current != nil && current.SchoolNo != "" && current.FullName != "" && !seen[current.SchoolNo] {
			seen[current.SchoolNo] = true
			if current.ClassSection == "" {
				current.ClassSection = currentClass
			}
			if current.BehaviorScore == 0 {
				current.BehaviorScore = 100
			}
			if current.BoardingStatus == "" {
				current.BoardingStatus = "GÜNDÜZLÜ"
			}
			list = append(list, current)
		}
	}

	for _, row := range rows {
		// Detect class banner
		banner := extractClassBanner(row)
		if banner != "" {
			currentClass = banner
		}

		for cIdx, cell := range row {
			val := strings.TrimSpace(cell)
			if val == "" {
				continue
			}
			norm := normalizeTurkish(val)
			nextCell := ""
			if cIdx+1 < len(row) {
				nextCell = strings.TrimSpace(row[cIdx+1])
			}

			// Check for T.C. Kimlik No
			if (strings.Contains(norm, "TC") || strings.Contains(norm, "KIMLIK NO")) && !strings.Contains(norm, "BAKANLIG") {
				rawNum := extractValueAfterLabel(val)
				if rawNum == "" {
					rawNum = nextCell
				}
				num := cleanNumeric(rawNum)
				if len(num) == 11 {
					// If current student already has TC and SchoolNo, this starts a new student card
					if current != nil && current.TCNo != "" && current.SchoolNo != "" {
						saveCurrent()
						current = &model.Student{AcademicYear: schoolYear, ClassSection: currentClass}
					} else if current == nil {
						current = &model.Student{AcademicYear: schoolYear, ClassSection: currentClass}
					}
					current.TCNo = num
				}
			}

			// Check for Okul No / Öğrenci No
			if strings.Contains(norm, "OKUL NO") || strings.Contains(norm, "OGRENCI NO") {
				rawNo := extractValueAfterLabel(val)
				if rawNo == "" {
					rawNo = nextCell
				}
				no := cleanNumeric(rawNo)
				if len(no) >= 1 && len(no) <= 6 {
					if current != nil && current.SchoolNo != "" && current.SchoolNo != no {
						saveCurrent()
						current = &model.Student{AcademicYear: schoolYear, ClassSection: currentClass}
					} else if current == nil {
						current = &model.Student{AcademicYear: schoolYear, ClassSection: currentClass}
					}
					current.SchoolNo = no
				}
			}

			// Check for Adı Soyadı
			if strings.Contains(norm, "ADI SOYADI") || strings.Contains(norm, "AD SOYAD") {
				name := extractValueAfterLabel(val)
				if name == "" {
					name = nextCell
				}
				if name != "" && hasAlphabetic(name) {
					if current == nil {
						current = &model.Student{AcademicYear: schoolYear, ClassSection: currentClass}
					}
					current.FullName = strings.TrimSpace(name)
				}
			} else if (norm == "ADI:" || norm == "ADI") && nextCell != "" {
				if current == nil {
					current = &model.Student{AcademicYear: schoolYear, ClassSection: currentClass}
				}
				if current.FullName == "" {
					current.FullName = nextCell
				}
			} else if (norm == "SOYADI:" || norm == "SOYADI") && nextCell != "" {
				if current != nil {
					current.FullName = strings.TrimSpace(current.FullName + " " + nextCell)
				}
			}

			// Check for Sınıfı / Şubesi
			if strings.Contains(norm, "SINIFI") || strings.Contains(norm, "SUBESI") || strings.Contains(norm, "SINIF") {
				cls := extractValueAfterLabel(val)
				if cls == "" {
					cls = nextCell
				}
				if cls != "" && len(cls) < 30 {
					cleanCls := strings.TrimSpace(cls)
					cleanCls = strings.ReplaceAll(cleanCls, " / ", "-")
					cleanCls = strings.ReplaceAll(cleanCls, "/", "-")
					cleanCls = strings.ReplaceAll(cleanCls, " ", "")
					if current != nil {
						current.ClassSection = cleanCls
					}
				}
			}

			// Check for Cinsiyet
			if strings.Contains(norm, "CINSIYET") {
				gen := extractValueAfterLabel(val)
				if gen == "" {
					gen = nextCell
				}
				genUpper := normalizeTurkish(gen)
				if strings.HasPrefix(genUpper, "K") {
					if current != nil {
						current.Gender = "KIZ"
					}
				} else if strings.HasPrefix(genUpper, "E") {
					if current != nil {
						current.Gender = "ERKEK"
					}
				}
			}

			// Check for Baba Adı
			if strings.Contains(norm, "BABA ADI") || strings.Contains(norm, "BABA") {
				baba := extractValueAfterLabel(val)
				if baba == "" {
					baba = nextCell
				}
				if current != nil && baba != "" {
					current.FatherName = baba
				}
			}

			// Check for Anne Adı
			if strings.Contains(norm, "ANNE ADI") || strings.Contains(norm, "ANNE") {
				anne := extractValueAfterLabel(val)
				if anne == "" {
					anne = nextCell
				}
				if current != nil && anne != "" {
					current.MotherName = anne
				}
			}

			// Check for Doğum Tarihi
			if strings.Contains(norm, "DOGUM TARIHI") {
				dt := extractValueAfterLabel(val)
				if dt == "" {
					dt = nextCell
				}
				if current != nil && dt != "" {
					current.BirthDate = dt
				}
			}

			// Check for Doğum Yeri
			if strings.Contains(norm, "DOGUM YERI") {
				dy := extractValueAfterLabel(val)
				if dy == "" {
					dy = nextCell
				}
				if current != nil && dy != "" {
					current.BirthPlace = dy
				}
			}
		}
	}
	saveCurrent()

	return list
}

func extractValueAfterLabel(s string) string {
	parts := strings.SplitN(s, ":", 2)
	if len(parts) == 2 {
		return strings.TrimSpace(parts[1])
	}
	return ""
}

// extractClassBanner checks if row is a section divider banner.
func extractClassBanner(row []string) string {
	for _, cell := range row {
		c := strings.TrimSpace(cell)
		if c == "" {
			continue
		}
		upper := strings.ToUpper(c)
		if (strings.Contains(upper, "SINIF") || strings.Contains(upper, "ŞUBE") || strings.Contains(upper, "SUBE")) &&
			(strings.Contains(upper, "9.") || strings.Contains(upper, "10.") || strings.Contains(upper, "11.") || strings.Contains(upper, "12.") ||
				strings.Contains(upper, "5.") || strings.Contains(upper, "6.") || strings.Contains(upper, "7.") || strings.Contains(upper, "8.") ||
				strings.Contains(upper, "HAZIRLIK") || strings.Contains(upper, "SINIFI /") || strings.Contains(upper, "/")) {
			return cleanBanner(c)
		}
	}
	return ""
}

func cleanBanner(raw string) string {
	parts := strings.Split(raw, ":")
	if len(parts) > 1 {
		return strings.TrimSpace(parts[len(parts)-1])
	}
	return strings.TrimSpace(raw)
}

func detectHeaderMapping(row []string) map[string]int {
	m := make(map[string]int)
	for idx, rawCell := range row {
		c := normalizeTurkish(rawCell)
		c = strings.ReplaceAll(c, ".", "")
		c = strings.ReplaceAll(c, ":", "")
		c = strings.ReplaceAll(c, "-", "")
		c = strings.ReplaceAll(c, "_", "")
		c = strings.TrimSpace(c)

		if strings.Contains(c, "OKUL NO") || strings.Contains(c, "OGRENCI NO") || c == "NO" || c == "NUMARASI" || c == "O NO" || strings.Contains(c, "NUMARA") {
			m["okulno"] = idx
		} else if strings.Contains(c, "TC") || strings.Contains(c, "KIMLIK") {
			m["tcno"] = idx
		} else if strings.Contains(c, "ADI SOYADI") || strings.Contains(c, "AD SOYAD") {
			m["adsoyad"] = idx
		} else if c == "ADI" || c == "AD" {
			m["ad"] = idx
		} else if c == "SOYADI" || c == "SOYAD" {
			m["soyad"] = idx
		} else if strings.Contains(c, "SINIF") || strings.Contains(c, "SUBE") {
			m["sinif"] = idx
		} else if strings.Contains(c, "BABA") {
			m["baba"] = idx
		} else if strings.Contains(c, "ANNE") {
			m["anne"] = idx
		} else if strings.Contains(c, "CINSIYET") {
			m["cinsiyet"] = idx
		} else if strings.Contains(c, "DOGUM TARIHI") {
			m["dogumtarihi"] = idx
		} else if strings.Contains(c, "DOGUM YERI") {
			m["dogumyeri"] = idx
		}
	}
	return m
}

func hasKey(m map[string]int, k string) bool {
	_, ok := m[k]
	return ok
}

func normalizeTurkish(s string) string {
	var sb strings.Builder
	for _, r := range s {
		switch r {
		case 'ı', 'I', 'İ', 'i':
			sb.WriteRune('I')
		case 'ğ', 'Ğ':
			sb.WriteRune('G')
		case 'ü', 'Ü':
			sb.WriteRune('U')
		case 'ş', 'Ş':
			sb.WriteRune('S')
		case 'ö', 'Ö':
			sb.WriteRune('O')
		case 'ç', 'Ç':
			sb.WriteRune('C')
		default:
			if r >= 'a' && r <= 'z' {
				sb.WriteRune(r - ('a' - 'A'))
			} else {
				sb.WriteRune(r)
			}
		}
	}
	return sb.String()
}

func getCell(row []string, idx int) string {
	if idx >= 0 && idx < len(row) {
		return strings.TrimSpace(row[idx])
	}
	return ""
}

func studentFromMappedRow(row []string, m map[string]int, defaultClass string, schoolYear string) *model.Student {
	okulNo := ""
	if idx, ok := m["okulno"]; ok {
		okulNo = cleanNumeric(getCell(row, idx))
	}
	if okulNo == "" {
		return nil
	}

	fullName := ""
	if idx, ok := m["adsoyad"]; ok {
		fullName = getCell(row, idx)
	} else {
		ad := ""
		soyad := ""
		if idx, ok := m["ad"]; ok {
			ad = getCell(row, idx)
		}
		if idx, ok := m["soyad"]; ok {
			soyad = getCell(row, idx)
		}
		fullName = strings.TrimSpace(ad + " " + soyad)
	}

	if fullName == "" {
		return nil
	}

	classStr := defaultClass
	if idx, ok := m["sinif"]; ok {
		val := getCell(row, idx)
		if val != "" {
			classStr = val
		}
	}

	gender := ""
	if idx, ok := m["cinsiyet"]; ok {
		g := normalizeTurkish(getCell(row, idx))
		if strings.HasPrefix(g, "K") {
			gender = "KIZ"
		} else if strings.HasPrefix(g, "E") {
			gender = "ERKEK"
		}
	}

	return &model.Student{
		AcademicYear:   schoolYear,
		SchoolNo:       okulNo,
		TCNo:           getCell(row, m["tcno"]),
		ClassSection:   classStr,
		FullName:       fullName,
		FatherName:     getCell(row, m["baba"]),
		MotherName:     getCell(row, m["anne"]),
		Gender:         gender,
		BirthDate:      getCell(row, m["dogumtarihi"]),
		BirthPlace:     getCell(row, m["dogumyeri"]),
		BoardingStatus: "GÜNDÜZLÜ",
		BehaviorScore:  100,
	}
}

func studentFromHeuristicRow(row []string, defaultClass string, schoolYear string) *model.Student {
	var noCol = -1
	var nameCol = -1
	var classCol = -1
	var tcCol = -1

	for idx, cell := range row {
		val := strings.TrimSpace(cell)
		if val == "" {
			continue
		}
		num := cleanNumeric(val)
		if num != "" {
			if len(num) == 11 && tcCol == -1 {
				tcCol = idx
			} else if len(num) >= 1 && len(num) <= 6 && noCol == -1 {
				noCol = idx
			}
		} else if hasAlphabetic(val) {
			norm := normalizeTurkish(val)
			if isClassCandidate(val) {
				if classCol == -1 {
					classCol = idx
				}
			} else if strings.Contains(norm, "SINIF") || strings.Contains(norm, "SUBE") {
				if classCol == -1 {
					classCol = idx
				}
			} else if nameCol == -1 && len(val) >= 3 && !strings.Contains(norm, "TOPLAM") && !strings.Contains(norm, "MUDUR") {
				nameCol = idx
			}
		}
	}

	if noCol == -1 {
		return nil
	}

	okulNo := cleanNumeric(row[noCol])
	fullName := ""
	if nameCol != -1 {
		fullName = strings.TrimSpace(row[nameCol])
		if nameCol+1 < len(row) {
			nextVal := strings.TrimSpace(row[nameCol+1])
			normNext := normalizeTurkish(nextVal)
			if isPureLetters(nextVal) && len(nextVal) < 25 &&
				!strings.Contains(normNext, "SINIF") && !strings.Contains(normNext, "SUBE") {
				fullName = fullName + " " + nextVal
			}
		}
	}

	if fullName == "" {
		return nil
	}

	classStr := defaultClass
	if classCol != -1 {
		classStr = strings.TrimSpace(row[classCol])
	}

	tc := ""
	if tcCol != -1 {
		tc = cleanNumeric(row[tcCol])
	}

	return &model.Student{
		AcademicYear:   schoolYear,
		SchoolNo:       okulNo,
		TCNo:           tc,
		ClassSection:   classStr,
		FullName:       fullName,
		BoardingStatus: "GÜNDÜZLÜ",
		BehaviorScore:  100,
	}
}

func isPureLetters(s string) bool {
	hasLetter := false
	for _, r := range s {
		if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') ||
			r == 'ç' || r == 'Ç' || r == 'ğ' || r == 'Ğ' || r == 'ı' || r == 'I' ||
			r == 'İ' || r == 'i' || r == 'ö' || r == 'Ö' || r == 'ş' || r == 'Ş' ||
			r == 'ü' || r == 'Ü' {
			hasLetter = true
		} else if r != ' ' {
			return false
		}
	}
	return hasLetter
}

func isClassCandidate(s string) bool {
	s = strings.TrimSpace(s)
	if len(s) < 2 || len(s) > 10 {
		return false
	}
	hasDigit := false
	hasLetter := false
	for _, r := range s {
		if r >= '0' && r <= '9' {
			hasDigit = true
		} else if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') {
			hasLetter = true
		}
	}
	return hasDigit && hasLetter
}

func cleanNumeric(s string) string {
	s = strings.TrimSpace(s)
	if f, err := strconv.ParseFloat(s, 64); err == nil && f >= 0 {
		return strconv.FormatInt(int64(f), 10)
	}
	var sb strings.Builder
	for _, r := range s {
		if r >= '0' && r <= '9' {
			sb.WriteRune(r)
		}
	}
	return sb.String()
}

func hasAlphabetic(s string) bool {
	for _, r := range s {
		if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') ||
			r == 'ç' || r == 'Ç' || r == 'ğ' || r == 'Ğ' ||
			r == 'ı' || r == 'İ' || r == 'ö' || r == 'Ö' ||
			r == 'ş' || r == 'Ş' || r == 'ü' || r == 'Ü' {
			return true
		}
	}
	return false
}

// ImportStudentsFromFile opens an excel/html/xml/csv file and extracts students for a given school year.
func ImportStudentsFromFile(filePath string, schoolYear string) ([]*model.Student, error) {
	sheetData, sheets, err := ReadSheetData(filePath)
	if err != nil {
		return nil, err
	}

	chosenSheet := SelectDefaultSheet(sheets)
	rows, ok := sheetData[chosenSheet]
	if !ok || len(rows) == 0 {
		for _, s := range sheets {
			if len(sheetData[s]) > 0 {
				rows = sheetData[s]
				break
			}
		}
	}

	if len(rows) == 0 {
		return nil, fmt.Errorf("sayfada veri bulunamadı")
	}

	return ParseStudentRows(rows, schoolYear)
}
