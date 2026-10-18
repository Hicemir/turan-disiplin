package main

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"fyne.io/fyne/v2/test"

	"turandisiplin/model"
	"turandisiplin/pdf"
	"turandisiplin/ui"
)

func TestFullAppEndToEnd(t *testing.T) {
	// 1. Initialize Headless Fyne App & UI State
	testApp := test.NewApp()
	defer testApp.Quit()

	testWin := testApp.NewWindow("TuranDisiplin Test")
	state, err := ui.InitState(testApp, testWin)
	if err != nil {
		t.Fatalf("InitState failed: %v", err)
	}

	// Build Main UI Container
	mainUI := ui.BuildMainUI(state)
	if mainUI == nil {
		t.Fatalf("BuildMainUI returned nil")
	}

	t.Log("===> Step 1: UI Mode & Theme Rapid Switching (Deadlock & Hang Verification)")
	// Cycle between legacy and modern multiple times to guarantee zero deadlock or freeze
	for i := 0; i < 8; i++ {
		if i%2 == 0 {
			state.SetUIMode("legacy")
			if state.Config.UIMode != "legacy" {
				t.Fatalf("Failed to switch to legacy mode at iteration %d", i)
			}
		} else {
			state.SetUIMode("modern")
			if state.Config.UIMode != "modern" {
				t.Fatalf("Failed to switch to modern mode at iteration %d", i)
			}
		}
	}

	// Test Theme Switching across all modes
	themes := []string{"cream", "dark", "light", "cream", "light"}
	for _, th := range themes {
		state.ApplyTheme(th, 15.0)
		if state.Config.ThemeMode != th {
			t.Fatalf("Expected theme %s, got %s", th, state.Config.ThemeMode)
		}
	}

	t.Log("===> Step 2: View Navigation Test Across All 13 Tabs")
	views := []string{
		"school", "students", "incidents", "verdict", "summons",
		"forms", "notices", "restoration", "minutes", "dizi", "stats",
		"yonetmelik", "settings", "hub",
	}
	for _, v := range views {
		state.NavigateTo(v)
		if state.CurrentView != v {
			t.Fatalf("Failed navigating to view: %s", v)
		}
	}

	t.Log("===> Step 3: School Info Management")
	testSchool := &model.SchoolInfo{
		SchoolType:      "LİSE",
		AcademicYear:    state.AcademicYear,
		SchoolName:      "TEST ANADOLU LİSESİ",
		InstitutionCode: "999888",
		Province:        "ANKARA",
		District:        "ÇANKAYA",
		PrincipalName:   "Mehmet Müdür",
		Letterhead:      "T.C. MİLLÎ EĞİTİM BAKANLIĞI",
		OfficialCode:    "B.08.4.MEM.4.06.00.00",
		BoardChair:      "Ahmet Başkan",
		BoardChairTitle: "Müdür Başyardımcısı",
		Member1:         "Ayşe Öğretmen",
		Member1Title:    "Edebiyat Öğretmeni",
		Member2:         "Fatma Öğretmen",
		Member2Title:    "Matematik Öğretmeni",
	}
	if err := state.Store.SaveSchoolInfo(testSchool); err != nil {
		t.Fatalf("SaveSchoolInfo failed: %v", err)
	}

	t.Log("===> Step 4: Adding Sample Students")
	st1 := &model.Student{
		AcademicYear:   state.AcademicYear,
		SchoolNo:       "9901",
		FullName:       "Deneme Öğrenci Bir",
		ClassSection:   "10-A",
		Gender:         "E",
		FatherName:     "Kemal",
		MotherName:     "Zeynep",
		BehaviorScore:  100,
	}
	st2 := &model.Student{
		AcademicYear:   state.AcademicYear,
		SchoolNo:       "9902",
		FullName:       "Deneme Öğrenci İki",
		ClassSection:   "10-B",
		Gender:         "K",
		FatherName:     "Ali",
		MotherName:     "Elif",
		BehaviorScore:  100,
	}
	if err := state.Store.SaveStudent(st1); err != nil {
		t.Fatalf("SaveStudent st1 failed: %v", err)
	}
	if err := state.Store.SaveStudent(st2); err != nil {
		t.Fatalf("SaveStudent st2 failed: %v", err)
	}

	retrievedSt1, err := state.Store.GetStudentByNo(state.AcademicYear, "9901")
	if err != nil || retrievedSt1 == nil {
		t.Fatalf("Failed to retrieve student st1")
	}

	t.Log("===> Step 5: Creating Sample Incident & Assigning Students")
	incCode := fmt.Sprintf("TST-%d", time.Now().UnixNano()%100000)
	incident := &model.Incident{
		AcademicYear: state.AcademicYear,
		IncidentCode: incCode,
		Title:        "Deneme Sınavında Kopya Çekme ve Görevliye İtiraz",
		Date:         "18.10.2026",
		Time:         "10:30",
		Location:     "10-A Sınıfı",
		Summary:      "Öğrencinin sınav esnasında kopya hazırladığı ve görevli öğretmene itiraz ettiği tespit edilmiştir.",
	}
	if err := state.Store.SaveIncident(incident); err != nil {
		t.Fatalf("SaveIncident failed: %v", err)
	}

	// Assign students to incident
	incStudent1 := &model.IncidentStudent{
		AcademicYear:    state.AcademicYear,
		IncidentCode:    incCode,
		IncidentTitle:   incident.Title,
		Role:            "SUÇLANAN ÖĞRENCİ",
		SchoolNo:        "9901",
		FullName:        "Deneme Öğrenci Bir",
		ClassSection:    "10-A",
		Gender:          "E",
		RemainingPoints: 100,
	}
	if err := state.Store.SaveIncidentStudent(incStudent1); err != nil {
		t.Fatalf("SaveIncidentStudent 1 failed: %v", err)
	}

	incStudent2 := &model.IncidentStudent{
		AcademicYear:    state.AcademicYear,
		IncidentCode:    incCode,
		IncidentTitle:   incident.Title,
		Role:            "BİLGİSİNE BAŞVURULACAK ÖĞRENCİ",
		SchoolNo:        "9902",
		FullName:        "Deneme Öğrenci İki",
		ClassSection:    "10-B",
		Gender:          "K",
		RemainingPoints: 100,
	}
	if err := state.Store.SaveIncidentStudent(incStudent2); err != nil {
		t.Fatalf("SaveIncidentStudent 2 failed: %v", err)
	}

	incStudents, err := state.Store.GetIncidentStudents(state.AcademicYear, incCode)
	if err != nil || len(incStudents) != 2 {
		t.Fatalf("Expected 2 incident students, got %d (err: %v)", len(incStudents), err)
	}

	t.Log("===> Step 6: Creating Verdict (Disiplin Kurulu Kararı)")
	// Update incStudent1 with verdict details
	incStudents[0].VerdictNo = "2026/999"
	incStudents[0].VerdictDate = "18.10.2026"
	incStudents[0].ClauseCode = "164/1-a"
	incStudents[0].ClauseDescription = "Sınavlarda kopya yapmak veya yaptırmak"
	incStudents[0].ImposedPenalty = "Kınama"
	incStudents[0].DeductedPoints = 10
	incStudents[0].RemainingPoints = 90
	incStudents[0].VerdictText = "Öğrencinin eylemine uyan MEB Yönetmeliği 164/1-a maddesi uyarınca KINAMA cezası ile cezalandırılmasına ve 10 davranış puanının indirilmesine oy birliğiyle karar verildi."
	incStudents[0].BoardOpinion = "SUÇLU OLDUĞU YÖNÜNDE"
	incStudents[0].OffensePlace = "10-A Sınıfı"
	incStudents[0].OffenseDateTime = "18.10.2026 10:30"
	incStudents[0].OffenseCategory = "Sınavda Kopya Hazırlama"
	incStudents[0].OffenseReason = "Derse hazırlıksız girme"
	incStudents[0].WitnessSummary = "Gözetmen öğretmen öğrencinin sırasının altında yazılı notlar bulduğunu belirtmiştir."
	incStudents[0].Evidence = "Kopya kağıtları ve sınav tutanağı"
	incStudents[0].MitigatingFactors = "Öğrencinin ilk disiplin olayı olması ve pişmanlık duyması"
	incStudents[0].AggravatingFactors = "Yok"

	if err := state.Store.SaveIncidentStudent(&incStudents[0]); err != nil {
		t.Fatalf("SaveIncidentStudent with verdict failed: %v", err)
	}

	// Update student's master score
	retrievedSt1.BehaviorScore = 90
	if err := state.Store.SaveStudent(retrievedSt1); err != nil {
		t.Fatalf("SaveStudent updated score failed: %v", err)
	}

	verdictObj := &model.Verdict{
		DecisionNo:         999,
		DecisionDate:       "18.10.2026",
		OffensePlace:       "10-A Sınıfı",
		OffenseDateTime:    "18.10.2026 10:30",
		OffenseType:        "Sınavda Kopya Hazırlama",
		OffenseReason:      "Derse hazırlıksız girme",
		WitnessStatements:  "Gözetmen öğretmen notları teslim etmiştir.",
		EvidenceDetails:    "Kopya pusulası",
		MitigatingFactors:  "Pişmanlık",
		AggravatingFactors: "Yok",
		BoardOpinion:       "Suçlu olduğu yönünde",
		RegClause:          "164/1-a",
		RegClauseDesc:      "Sınavlarda kopya yapmak veya yaptırmak",
		PenaltyName:        "Kınama",
		DeductedPoints:     10,
		RemainingPoints:    90,
		DecisionText:       incStudents[0].VerdictText,
	}

	t.Log("===> Step 7: Generating Official PDFs for All Modules")
	tmpDir, err := os.MkdirTemp("", "turandisiplin_test_*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	verifyPDF := func(name, path string) {
		info, err := os.Stat(path)
		if err != nil {
			t.Fatalf("[%s] PDF stat failed: %v", name, err)
		}
		if info.Size() < 1000 {
			t.Fatalf("[%s] PDF file too small: %d bytes", name, info.Size())
		}
		header := make([]byte, 5)
		f, err := os.Open(path)
		if err != nil {
			t.Fatalf("[%s] PDF open failed: %v", name, err)
		}
		defer f.Close()
		_, _ = f.Read(header)
		if !bytes.Equal(header, []byte("%PDF-")) {
			t.Fatalf("[%s] PDF invalid magic header: %s", name, string(header))
		}
		t.Logf("  ✓ [%s] Successfully generated (%d bytes): %s", name, info.Size(), filepath.Base(path))
	}

	// 7.1 EK-1 / EK-10 Karar Tutanağı PDF
	pKarar := filepath.Join(tmpDir, "01_EK1_Karar.pdf")
	if err := pdf.GenerateKurulKarari(testSchool, retrievedSt1, incident, &incStudents[0], verdictObj, pKarar); err != nil {
		t.Fatalf("GenerateKurulKarari failed: %v", err)
	}
	verifyPDF("EK-1 Karar Tutanağı", pKarar)

	// 7.2 Çağrı Pusulaları PDF (İfade & Savunma)
	pCagri := filepath.Join(tmpDir, "02_Cagri_Pusulasi.pdf")
	if err := pdf.GenerateCagriPusulasi(testSchool, []*model.Student{retrievedSt1}, incident, true, "Disiplin Odası", "19.10.2026", "11:00", pCagri); err != nil {
		t.Fatalf("GenerateCagriPusulasi failed: %v", err)
	}
	verifyPDF("Çağrı Pusulası", pCagri)

	// 7.3 Savunma Tutanağı PDF
	pSavunma := filepath.Join(tmpDir, "03_Savunma_Tutanagi.pdf")
	if err := pdf.GenerateIfadeSavunmaTutanagi(testSchool, retrievedSt1, incident, true, "Öğrenci savunmasında pişman olduğunu belirtmiştir.", pSavunma); err != nil {
		t.Fatalf("GenerateIfadeSavunmaTutanagi failed: %v", err)
	}
	verifyPDF("Savunma Tutanağı", pSavunma)

	// 7.4 Öğrenci Karar Tebligat Pusulası PDF
	pTebligat := filepath.Join(tmpDir, "04_Ogrenci_Tebligat.pdf")
	if err := pdf.GenerateCezaTebligatiOgrenci(testSchool, retrievedSt1, incident, &incStudents[0], verdictObj, pTebligat); err != nil {
		t.Fatalf("GenerateCezaTebligatiOgrenci failed: %v", err)
	}
	verifyPDF("Öğrenci Karar Tebligat Pusulası", pTebligat)

	// 7.5 Veli Karar Tebligat Pusulası PDF
	pTebligatVeli := filepath.Join(tmpDir, "05_Veli_Tebligat.pdf")
	if err := pdf.GenerateCezaTebligatiVeli(testSchool, retrievedSt1, incident, &incStudents[0], verdictObj, "2026/999-1", pTebligatVeli); err != nil {
		t.Fatalf("GenerateCezaTebligatiVeli failed: %v", err)
	}
	verifyPDF("Veli Karar Tebligat Pusulası", pTebligatVeli)

	// 7.6 Sene Başı / Sene Sonu Toplantı Tutanağı PDF
	pToplanti := filepath.Join(tmpDir, "06_Toplanti_Tutanagi.pdf")
	decisions := []string{
		"1. Okul düzeni ve disiplin kuralları öğrencilere duyurulacak.",
		"2. Ödül ve disiplin kurulu ayda en az bir kez toplanacak.",
	}
	if err := pdf.GenerateSeneBasiSonuKarar(testSchool, "SENE BAŞI", "18.10.2026", "09:30", "Öğretmenler Odası", "Müdür Başyardımcısı", decisions, pToplanti); err != nil {
		t.Fatalf("GenerateSeneBasiSonuKarar failed: %v", err)
	}
	verifyPDF("Toplantı Tutanağı", pToplanti)

	// 7.7 Dizi Pusulası PDF
	pDizi := filepath.Join(tmpDir, "07_Dizi_Pusulasi.pdf")
	diziItems := []model.DiziItem{
		{Subject: "Olay İnceleme Raporu", Count: 1},
		{Subject: "İfade ve Savunma Tutanakları", Count: 2},
		{Subject: "Disiplin Kurulu Kararı", Count: 1},
	}
	if err := pdf.GenerateDiziPusulasi(testSchool, diziItems, pDizi); err != nil {
		t.Fatalf("GenerateDiziPusulasi failed: %v", err)
	}
	verifyPDF("Dizi Pusulası", pDizi)

	// 7.8 Davranış Puanı İade Çizelgesi PDF
	pIade := filepath.Join(tmpDir, "08_Puan_Iade_Cizelgesi.pdf")
	refundRows := []model.BehaviorScoreRefundRow{
		{
			FullName:        retrievedSt1.FullName,
			SchoolNo:        retrievedSt1.SchoolNo,
			ClassSection:    retrievedSt1.ClassSection,
			IncidentName:    incident.Title,
			DecisionNo:      999,
			DecisionDate:    "18.10.2026",
			PenaltyName:     "Kınama",
			RegClause:       "164/1-a",
			DeductedPoints:  10,
			RemainingPoints: 90,
			RefundPoints:    10,
			RefundDate:      "18.10.2026",
			Status:          "İADE EDİLDİ",
		},
	}
	if err := pdf.GenerateDavranisPuaniIadeCizelgesi(testSchool, refundRows, pIade); err != nil {
		t.Fatalf("GenerateDavranisPuaniIadeCizelgesi failed: %v", err)
	}
	verifyPDF("Davranış Puanı İade Çizelgesi", pIade)

	// 7.9 Ceza İstatistik Raporu PDF
	pStats := filepath.Join(tmpDir, "09_Ceza_Raporu.pdf")
	statItems := []model.StatItem{
		{Name: "Kınama", Count: 1},
		{Name: "Yazılı Uyarma", Count: 0},
	}
	if err := pdf.GenerateCezaIstatistikRaporu(testSchool, "1. DÖNEM CEZA RAPORU", statItems, pStats); err != nil {
		t.Fatalf("GenerateCezaIstatistikRaporu failed: %v", err)
	}
	verifyPDF("Ceza İstatistik Raporu", pStats)

	// 7.10 MEB Madde 163/2 Yazılı Uyarı Tutanağı PDF
	pUyari := filepath.Join(tmpDir, "10_Yazili_Uyari.pdf")
	if err := pdf.GenerateYaziliUyariTutanagi(testSchool, retrievedSt1, "2026/01", "18.10.2026", "163/2", "Öğrenciye ders düzenini bozmaması konusunda yazılı uyarı yapılmıştır.", "Müdür Yardımcısı", pUyari); err != nil {
		t.Fatalf("GenerateYaziliUyariTutanagi failed: %v", err)
	}
	verifyPDF("Yazılı Uyarı Tutanağı", pUyari)

	t.Log("===> Step 8: Point Restoration (Davranış Puanı İadesi)")
	incStudents[0].RestorationStatus = "İADE EDİLDİ"
	incStudents[0].RestorationDate = "18.10.2026"
	incStudents[0].RemainingPoints = 100
	if err := state.Store.SaveIncidentStudent(&incStudents[0]); err != nil {
		t.Fatalf("SaveIncidentStudent restoration failed: %v", err)
	}

	retrievedSt1.BehaviorScore = 100
	if err := state.Store.SaveStudent(retrievedSt1); err != nil {
		t.Fatalf("SaveStudent restoration failed: %v", err)
	}

	restoredSt1, _ := state.Store.GetStudentByNo(state.AcademicYear, "9901")
	if restoredSt1.BehaviorScore != 100 {
		t.Fatalf("Expected behavior score to be restored to 100, got %d", restoredSt1.BehaviorScore)
	}

	t.Log("===> Step 9: Clean Up All Sample Test Data (Return Fresh Clean Database)")
	// Delete incident (cascades to incident_students)
	if err := state.Store.DeleteIncident(state.AcademicYear, incCode); err != nil {
		t.Fatalf("DeleteIncident failed: %v", err)
	}
	// Delete sample students
	if err := state.Store.DeleteStudent(retrievedSt1.ID); err != nil {
		t.Fatalf("DeleteStudent st1 failed: %v", err)
	}
	st2Retrieved, _ := state.Store.GetStudentByNo(state.AcademicYear, "9902")
	if st2Retrieved != nil {
		_ = state.Store.DeleteStudent(st2Retrieved.ID)
	}

	// Verify database is completely clean
	afterStudents, _ := state.Store.GetStudents(state.AcademicYear, "", "")
	for _, s := range afterStudents {
		if s.SchoolNo == "9901" || s.SchoolNo == "9902" {
			t.Fatalf("Test student %s was not cleaned up!", s.SchoolNo)
		}
	}
	afterIncidents, _ := state.Store.GetIncidents(state.AcademicYear)
	for _, inc := range afterIncidents {
		if inc.IncidentCode == incCode {
			t.Fatalf("Test incident %s was not cleaned up!", inc.IncidentCode)
		}
	}

	t.Log("===> All tests completed successfully! Database is verified fresh, clean, and bugless.")
}

func TestPenaltyPointDeductions(t *testing.T) {
	testCases := []struct {
		penalty  string
		expected int
	}{
		{"Kınama", 10},
		{"KINAMA", 10},
		{"Okuldan Kısa Süreli Uzaklaştırma (1-5 Gün)", 20},
		{"Kısa Süreli Uzaklaştırma", 20},
		{"Okul Değiştirme", 40},
		{"Örgün Eğitim Dışına Çıkarma", 80},
		{"Sözlü Uyarma", 0},
		{"Yazılı Uyarma", 0},
		{"Uyarma", 0},
		{"", 0},
	}

	for _, tc := range testCases {
		got := ui.GetPenaltyPointsDeduction(tc.penalty)
		if got != tc.expected {
			t.Errorf("GetPenaltyPointsDeduction(%q) = %d, expected %d", tc.penalty, got, tc.expected)
		}
	}
}
