package excel

import (
	"os"
	"path/filepath"
	"testing"
)

func TestReadSheetData_HTMLTable(t *testing.T) {
	htmlContent := `<!DOCTYPE html>
<html>
<head><meta charset="utf-8"></head>
<body>
<table>
  <tr>
    <th>Öğrenci No</th>
    <th>Adı Soyadı</th>
    <th>Sınıfı</th>
    <th>TC Kimlik</th>
  </tr>
  <tr>
    <td>101</td>
    <td>Ahmet Yılmaz</td>
    <td>9-A</td>
    <td>11111111111</td>
  </tr>
  <tr>
    <td>102</td>
    <td>Ayşe Demir</td>
    <td>9-B</td>
    <td>22222222222</td>
  </tr>
</table>
</body>
</html>`

	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "eokul_test.xls")
	if err := os.WriteFile(filePath, []byte(htmlContent), 0644); err != nil {
		t.Fatalf("Temp file error: %v", err)
	}

	data, sheets, err := ReadSheetData(filePath)
	if err != nil {
		t.Fatalf("ReadSheetData failed on HTML table: %v", err)
	}
	if len(sheets) == 0 {
		t.Fatalf("Expected at least 1 sheet, got 0")
	}
	rows := data[sheets[0]]
	for rIdx, r := range rows {
		t.Logf("Raw Row %d: %q", rIdx, r)
	}

	students, err := ParseStudentRows(rows, "2024 - 2025")
	if err != nil {
		t.Fatalf("ParseStudentRows failed: %v", err)
	}
	for i, st := range students {
		t.Logf("Student %d: SchoolNo=%q FullName=%q Class=%q TC=%q", i, st.SchoolNo, st.FullName, st.ClassSection, st.TCNo)
	}
	if len(students) != 2 {
		t.Fatalf("Expected 2 students, got %d", len(students))
	}
	if students[0].SchoolNo != "101" || students[0].FullName != "Ahmet Yılmaz" {
		t.Errorf("Unexpected student 0: %+v", students[0])
	}
	if students[1].SchoolNo != "102" || students[1].FullName != "Ayşe Demir" {
		t.Errorf("Unexpected student 1: %+v", students[1])
	}
}

func TestReadSheetData_XMLSpreadsheet(t *testing.T) {
	xmlContent := `<?xml version="1.0"?>
<?mso-application progid="Excel.Sheet"?>
<Workbook xmlns="urn:schemas-microsoft-com:office:spreadsheet"
 xmlns:o="urn:schemas-microsoft-com:office:office"
 xmlns:x="urn:schemas-microsoft-com:office:excel"
 xmlns:ss="urn:schemas-microsoft-com:office:spreadsheet"
 xmlns:html="http://www.w3.org/TR/REC-html40">
 <Worksheet ss:Name="Öğrenci Listesi">
  <Table>
   <Row>
    <Cell><Data ss:Type="String">Okul No</Data></Cell>
    <Cell><Data ss:Type="String">Adı Soyadı</Data></Cell>
    <Cell><Data ss:Type="String">Sınıf</Data></Cell>
   </Row>
   <Row>
    <Cell><Data ss:Type="String">205</Data></Cell>
    <Cell><Data ss:Type="String">Mehmet Kaya</Data></Cell>
    <Cell><Data ss:Type="String">11-C</Data></Cell>
   </Row>
  </Table>
 </Worksheet>
</Workbook>`

	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "eokul_data_only.xls")
	if err := os.WriteFile(filePath, []byte(xmlContent), 0644); err != nil {
		t.Fatalf("Temp file error: %v", err)
	}

	data, sheets, err := ReadSheetData(filePath)
	if err != nil {
		t.Fatalf("ReadSheetData failed on XML SpreadsheetML: %v", err)
	}
	if len(sheets) == 0 {
		t.Fatalf("Expected at least 1 sheet, got 0")
	}
	rows := data[sheets[0]]
	if len(rows) != 2 {
		t.Fatalf("Expected 2 rows, got %d", len(rows))
	}

	students, err := ParseStudentRows(rows, "2024 - 2025")
	if err != nil {
		t.Fatalf("ParseStudentRows failed: %v", err)
	}
	if len(students) != 1 {
		t.Fatalf("Expected 1 student, got %d", len(students))
	}
	if students[0].SchoolNo != "205" || students[0].FullName != "Mehmet Kaya" {
		t.Errorf("Unexpected student: %+v", students[0])
	}
}

func TestReadSheetData_CardBasedEokulReport(t *testing.T) {
	// e-Okul "Fotoğraflı-Kimlik Bilgili Öğrenci Listesi" HTML layout
	cardHTML := `<html>
<body>
<table>
  <tr><td colspan="4">T.C. MİLLÎ EĞİTİM BAKANLIĞI</td></tr>
  <tr><td colspan="4">Fotoğraflı-Kimlik Bilgili Öğrenci Listesi</td></tr>
  <tr>
    <td>T.C. Kimlik No :</td><td>12345678901</td>
    <td>Öğrenci No :</td><td>301</td>
  </tr>
  <tr>
    <td>Adı Soyadı :</td><td>FATMA ÇELİK</td>
    <td>Sınıfı / Şubesi :</td><td>10 / B</td>
  </tr>
  <tr>
    <td>Baba Adı :</td><td>ALİ</td>
    <td>Anne Adı :</td><td>HATİCE</td>
  </tr>
  <tr>
    <td>Cinsiyeti :</td><td>Kız</td>
    <td>Telefon :</td><td>05551112233</td>
  </tr>
  <tr>
    <td>T.C. Kimlik No :</td><td>98765432109</td>
    <td>Öğrenci No :</td><td>302</td>
  </tr>
  <tr>
    <td>Adı Soyadı :</td><td>BURAK ÖZTÜRK</td>
    <td>Sınıfı / Şubesi :</td><td>10 / B</td>
  </tr>
  <tr>
    <td>Baba Adı :</td><td>HASAN</td>
    <td>Anne Adı :</td><td>EMİNE</td>
  </tr>
</table>
</body>
</html>`

	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "fotografli_liste.xls")
	if err := os.WriteFile(filePath, []byte(cardHTML), 0644); err != nil {
		t.Fatalf("Temp file error: %v", err)
	}

	data, sheets, err := ReadSheetData(filePath)
	if err != nil {
		t.Fatalf("ReadSheetData failed on card-based layout: %v", err)
	}
	rows := data[sheets[0]]

	students, err := ParseStudentRows(rows, "2024 - 2025")
	if err != nil {
		t.Fatalf("ParseStudentRows failed on card layout: %v", err)
	}
	for i, st := range students {
		t.Logf("Card Student %d: SchoolNo=%q FullName=%q Class=%q TC=%q", i, st.SchoolNo, st.FullName, st.ClassSection, st.TCNo)
	}
	if len(students) != 2 {
		t.Fatalf("Expected 2 card-based students, got %d", len(students))
	}
	if students[0].SchoolNo != "301" || students[0].FullName != "FATMA ÇELİK" || students[0].TCNo != "12345678901" {
		t.Errorf("Unexpected student 0: %+v", students[0])
	}
	if students[0].ClassSection != "10-B" {
		t.Errorf("Expected class 10-B, got %s", students[0].ClassSection)
	}
	if students[1].SchoolNo != "302" || students[1].FullName != "BURAK ÖZTÜRK" || students[1].TCNo != "98765432109" {
		t.Errorf("Unexpected student 1: %+v", students[1])
	}
}

func TestReadSheetData_CorruptedFileCrashProof(t *testing.T) {
	// Verify that totally corrupted bytes or incomplete files never panic or crash
	corruptBytes := []byte{0xD0, 0xCF, 0x11, 0xE0, 0xA1, 0xB1, 0x1A, 0xE1, 0xFF, 0x00, 0x12, 0x34}
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "corrupted.xls")
	if err := os.WriteFile(filePath, corruptBytes, 0644); err != nil {
		t.Fatalf("Temp file error: %v", err)
	}

	// Should not panic, should return gracefully
	_, _, err := ReadSheetData(filePath)
	// Any result is fine as long as it does not panic/crash
	_ = err
}
