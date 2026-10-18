package excel

import (
	"fmt"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/extrame/xls"
	"github.com/xuri/excelize/v2"

	"okulevrak/model"
)

// ReadSheetData reads all sheets from an .xlsx or .xls file into a map of sheetName -> [][]string.
func ReadSheetData(filePath string) (map[string][][]string, []string, error) {
	ext := strings.ToLower(filepath.Ext(filePath))
	if ext == ".xlsx" || ext == ".xlsm" {
		return readXLSX(filePath)
	} else if ext == ".xls" {
		return readXLS(filePath)
	}
	// Try XLSX first, then fallback to XLS
	data, sheets, err := readXLSX(filePath)
	if err == nil && len(sheets) > 0 {
		return data, sheets, nil
	}
	return readXLS(filePath)
}

func readXLSX(filePath string) (map[string][][]string, []string, error) {
	f, err := excelize.OpenFile(filePath)
	if err != nil {
		return nil, nil, fmt.Errorf("excelize açma hatası: %w", err)
	}
	defer f.Close()

	sheets := f.GetSheetList()
	if len(sheets) == 0 {
		return nil, nil, fmt.Errorf("çalışma kitabında sayfa bulunamadı")
	}

	result := make(map[string][][]string)
	for _, sheet := range sheets {
		rows, err := f.GetRows(sheet)
		if err != nil {
			continue
		}
		result[sheet] = rows
	}

	return result, sheets, nil
}

func readXLS(filePath string) (map[string][][]string, []string, error) {
	wb, err := xls.Open(filePath, "utf-8")
	if err != nil {
		return nil, nil, fmt.Errorf("xls açma hatası: %w", err)
	}

	numSheets := wb.NumSheets()
	if numSheets == 0 {
		return nil, nil, fmt.Errorf("xls dosyasında sayfa bulunamadı")
	}

	result := make(map[string][][]string)
	var sheets []string

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
		result[sheet.Name] = rows
	}

	return result, sheets, nil
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

// ParseStudentRows intelligently parses rows into Student structs.
func ParseStudentRows(rows [][]string, schoolYear string) ([]*model.Student, error) {
	if len(rows) == 0 {
		return nil, fmt.Errorf("excel verisi boş")
	}

	var students []*model.Student
	seenNo := make(map[string]bool)
	currentClass := ""

	// Check if any row looks like a table header
	headerRowIdx := -1
	colMap := make(map[string]int)

	for i, row := range rows {
		if i > 25 {
			break // Header is usually within the first 25 rows
		}
		m := detectHeaderMapping(row)
		if len(m) >= 2 && (m["okulno"] >= 0 || m["adsoyad"] >= 0 || (m["ad"] >= 0 && m["soyad"] >= 0)) {
			headerRowIdx = i
			colMap = m
			break
		}
	}

	startRow := 0
	if headerRowIdx >= 0 {
		startRow = headerRowIdx + 1
	}

	for r := startRow; r < len(rows); r++ {
		row := rows[r]
		if len(row) == 0 {
			continue
		}

		// 1. Detect Class Header Banner (e.g. "9. Sınıf / A Şubesi", "10 / A Sınıfı")
		banner := extractClassBanner(row)
		if banner != "" {
			currentClass = banner
			continue
		}

		// 2. If we found a structured header mapping, use it
		if headerRowIdx >= 0 {
			st := studentFromMappedRow(row, colMap, currentClass, schoolYear)
			if st != nil && st.SchoolNo != "" && !seenNo[st.SchoolNo] {
				seenNo[st.SchoolNo] = true
				students = append(students, st)
			}
			continue
		}

		// 3. Fallback heuristic detection for unheadered/complex e-Okul layouts
		st := studentFromHeuristicRow(row, currentClass, schoolYear)
		if st != nil && st.SchoolNo != "" && !seenNo[st.SchoolNo] {
			seenNo[st.SchoolNo] = true
			students = append(students, st)
		}
	}

	if len(students) == 0 {
		return nil, fmt.Errorf("geçerli öğrenci kaydı bulunamadı (okul numarası veya ad sütunları tespit edilemedi)")
	}

	return students, nil
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
				strings.Contains(upper, "HAZIRLIK") || strings.Contains(upper, "SINIFI /")) {
			// Clean up banner
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
		c := strings.ToUpper(strings.TrimSpace(rawCell))
		c = strings.ReplaceAll(c, "İ", "I")
		c = strings.ReplaceAll(c, "I", "I")
		c = strings.ReplaceAll(c, ".", "")

		if strings.Contains(c, "OKUL NO") || strings.Contains(c, "OGRENCI NO") || c == "NO" || c == "NUMARASI" {
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
		g := strings.ToUpper(getCell(row, idx))
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
			upper := strings.ToUpper(val)
			if strings.Contains(upper, "SINIF") || strings.Contains(upper, "ŞUBE") {
				if classCol == -1 {
					classCol = idx
				}
			} else if nameCol == -1 && len(val) >= 3 && !strings.Contains(upper, "TOPLAM") && !strings.Contains(upper, "MÜDÜR") {
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
		// Check if next column is surname
		if nameCol+1 < len(row) && hasAlphabetic(row[nameCol+1]) && !strings.Contains(strings.ToUpper(row[nameCol+1]), "SINIF") {
			nextVal := strings.TrimSpace(row[nameCol+1])
			if nextVal != "" && len(nextVal) < 25 {
				// Likely surname
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

func cleanNumeric(s string) string {
	s = strings.TrimSpace(s)
	// If float format like 123.0
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

// ImportStudentsFromFile opens an excel file and extracts students for a given school year.
func ImportStudentsFromFile(filePath string, schoolYear string) ([]*model.Student, error) {
	sheetData, sheets, err := ReadSheetData(filePath)
	if err != nil {
		return nil, err
	}

	chosenSheet := SelectDefaultSheet(sheets)
	rows, ok := sheetData[chosenSheet]
	if !ok || len(rows) == 0 {
		// try any other sheet
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
