package store

import (
	"database/sql"
	"fmt"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	_ "modernc.org/sqlite"
	"okulevrak/config"
	"okulevrak/model"
)

type Store struct {
	db   *sql.DB
	lock sync.RWMutex
}

var (
	globalStore *Store
	storeOnce   sync.Once
)

// GetStore returns the singleton SQLite store instance
func GetStore() (*Store, error) {
	var initErr error
	storeOnce.Do(func() {
		cfgDir, err := config.GetConfigDir()
		if err != nil {
			initErr = err
			return
		}
		dbPath := filepath.Join(cfgDir, "OkulEvrak.db")
		db, err := sql.Open("sqlite", dbPath)
		if err != nil {
			initErr = fmt.Errorf("veritabanı açılamadı (%s): %w", dbPath, err)
			return
		}

		s := &Store{db: db}
		if err := s.migrate(); err != nil {
			initErr = fmt.Errorf("tablo kurulumu başarısız: %w", err)
			return
		}
		globalStore = s
	})
	return globalStore, initErr
}

func (s *Store) migrate() error {
	queries := []string{
		`CREATE TABLE IF NOT EXISTS school_info (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			school_type TEXT DEFAULT 'LİSE',
			academic_year TEXT DEFAULT '2026 - 2027',
			institution_code TEXT DEFAULT '',
			province TEXT DEFAULT '',
			district TEXT DEFAULT '',
			school_name TEXT DEFAULT '',
			principal_name TEXT DEFAULT '',
			letterhead TEXT DEFAULT '',
			official_code TEXT DEFAULT '',
			address TEXT DEFAULT '',
			phone TEXT DEFAULT '',
			fax TEXT DEFAULT '',
			warning_counter INTEGER DEFAULT 0,
			board_chair TEXT DEFAULT '',
			board_chair_title TEXT DEFAULT 'Müdür Yardımcısı',
			member1 TEXT DEFAULT '',
			member1_title TEXT DEFAULT 'Öğretmen',
			member2 TEXT DEFAULT '',
			member2_title TEXT DEFAULT 'Öğretmen',
			member3 TEXT DEFAULT '',
			member3_title TEXT DEFAULT 'Onur Kurulu 2. Başkanı',
			member4 TEXT DEFAULT '',
			member4_title TEXT DEFAULT 'Öğrenci Velisi',
			sub_board_chair TEXT DEFAULT '',
			sub_board_chair_title TEXT DEFAULT 'Müdür Yardımcısı',
			sub_member1 TEXT DEFAULT '',
			sub_member1_title TEXT DEFAULT 'Öğretmen',
			sub_member2 TEXT DEFAULT '',
			sub_member2_title TEXT DEFAULT 'Öğretmen',
			sub_member3 TEXT DEFAULT '',
			sub_member3_title TEXT DEFAULT 'Öğretmen',
			sub_member4 TEXT DEFAULT '',
			sub_member4_title TEXT DEFAULT 'Öğrenci Velisi'
		);`,
		`CREATE TABLE IF NOT EXISTS students (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			academic_year TEXT NOT NULL,
			school_no TEXT NOT NULL,
			full_name TEXT NOT NULL,
			class_section TEXT DEFAULT '',
			gender TEXT DEFAULT '',
			tc_no TEXT DEFAULT '',
			birth_place TEXT DEFAULT '',
			birth_date TEXT DEFAULT '',
			father_name TEXT DEFAULT '',
			mother_name TEXT DEFAULT '',
			guardian_name TEXT DEFAULT '',
			guardian_phone TEXT DEFAULT '',
			guardian_addr TEXT DEFAULT '',
			boarding_status TEXT DEFAULT 'GÜNDÜZLÜ',
			behavior_score INTEGER DEFAULT 100,
			photo_path TEXT DEFAULT '',
			UNIQUE(academic_year, school_no)
		);`,
		`CREATE TABLE IF NOT EXISTS incidents (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			academic_year TEXT NOT NULL,
			incident_code TEXT NOT NULL,
			title TEXT NOT NULL,
			location TEXT DEFAULT '',
			date TEXT DEFAULT '',
			time TEXT DEFAULT '',
			complainant TEXT DEFAULT '',
			petition_ref TEXT DEFAULT '',
			summary TEXT DEFAULT '',
			UNIQUE(academic_year, incident_code)
		);`,
		`CREATE TABLE IF NOT EXISTS incident_students (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			academic_year TEXT NOT NULL,
			incident_code TEXT NOT NULL,
			incident_title TEXT DEFAULT '',
			role TEXT DEFAULT 'SUÇLANAN ÖĞRENCİ',
			school_no TEXT NOT NULL,
			full_name TEXT NOT NULL,
			class_section TEXT DEFAULT '',
			gender TEXT DEFAULT '',
			tc_no TEXT DEFAULT '',
			birth_place TEXT DEFAULT '',
			birth_date TEXT DEFAULT '',
			father_name TEXT DEFAULT '',
			mother_name TEXT DEFAULT '',
			guardian_name TEXT DEFAULT '',
			guardian_phone TEXT DEFAULT '',
			guardian_addr TEXT DEFAULT '',
			boarding_status TEXT DEFAULT 'GÜNDÜZLÜ',
			academic_standing TEXT DEFAULT 'ORTA',
			health_status TEXT DEFAULT 'NORMAL',
			economic_status TEXT DEFAULT 'NORMAL',
			living_with TEXT DEFAULT 'AİLESİNİN YANINDA KALIYOR',
			parents_living_status TEXT DEFAULT 'ANNE BABA SAĞ',
			parents_bio_status TEXT DEFAULT 'ANNE VE BABA ÖZ',
			school_commute TEXT DEFAULT 'AİLESİNİN YANINDA OKUYOR',
			environment_raised TEXT DEFAULT 'İLÇE MERKEZİ',
			prior_penalty_status TEXT DEFAULT 'CEZA ALMADI',
			graduated_school TEXT DEFAULT '',
			statement_text TEXT DEFAULT '',
			written_warning TEXT DEFAULT '',
			warning_date TEXT DEFAULT '',
			warning_behavior TEXT DEFAULT '',
			warning_number TEXT DEFAULT '',
			warning_penalty TEXT DEFAULT '',
			warning_clause TEXT DEFAULT '',
			mother_age TEXT DEFAULT '',
			mother_education TEXT DEFAULT '',
			mother_alive TEXT DEFAULT 'EVET',
			mother_bio TEXT DEFAULT 'ÖZ',
			mother_together TEXT DEFAULT 'EVET',
			father_age TEXT DEFAULT '',
			father_education TEXT DEFAULT '',
			father_alive TEXT DEFAULT 'EVET',
			father_bio TEXT DEFAULT 'ÖZ',
			father_together TEXT DEFAULT 'EVET',
			sibling_count TEXT DEFAULT '',
			offense_place TEXT DEFAULT '',
			offense_date_time TEXT DEFAULT '',
			offense_category TEXT DEFAULT '',
			offense_reason TEXT DEFAULT '',
			defense_summary TEXT DEFAULT 'SUÇUNU KABUL EDİYOR',
			witness_summary TEXT DEFAULT '',
			evidence TEXT DEFAULT '',
			mitigating_factors TEXT DEFAULT '',
			aggravating_factors TEXT DEFAULT '',
			board_opinion TEXT DEFAULT 'SUÇLU OLDUĞU YÖNÜNDE',
			clause_code TEXT DEFAULT '',
			clause_penalty TEXT DEFAULT '',
			clause_description TEXT DEFAULT '',
			imposed_penalty TEXT DEFAULT '',
			suspension_days TEXT DEFAULT '',
			suspension_dates TEXT DEFAULT '',
			deducted_points INTEGER DEFAULT 0,
			remaining_points INTEGER DEFAULT 100,
			verdict_text TEXT DEFAULT '',
			verdict_no TEXT DEFAULT '',
			verdict_date TEXT DEFAULT '',
			restoration_status TEXT DEFAULT '',
			restoration_date TEXT DEFAULT ''
		);`,
		`CREATE TABLE IF NOT EXISTS agenda_items (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			order_no INTEGER,
			item_text TEXT
		);`,
		`CREATE TABLE IF NOT EXISTS meeting_decisions (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			academic_year TEXT NOT NULL,
			meeting_name TEXT NOT NULL,
			meeting_date TEXT,
			meeting_time TEXT,
			meeting_place TEXT,
			meeting_chair TEXT,
			decision_no INTEGER,
			decision_text TEXT
		);`,
		`CREATE TABLE IF NOT EXISTS reward_penalty_counts (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			academic_year TEXT NOT NULL,
			meeting_name TEXT NOT NULL,
			meeting_date TEXT,
			meeting_time TEXT,
			meeting_place TEXT,
			meeting_chair TEXT,
			kind TEXT,
			name TEXT,
			sem1_count INTEGER DEFAULT 0,
			sem2_count INTEGER DEFAULT 0,
			total_count INTEGER DEFAULT 0
		);`,
	}

	for _, q := range queries {
		if _, err := s.db.Exec(q); err != nil {
			return err
		}
	}

	// Seed school_info row if empty
	var count int
	_ = s.db.QueryRow("SELECT COUNT(*) FROM school_info").Scan(&count)
	if count == 0 {
		_, err := s.db.Exec(`INSERT INTO school_info (
			school_type, academic_year, school_name, board_chair_title, member1_title, member2_title, member3_title, member4_title,
			sub_board_chair_title, sub_member1_title, sub_member2_title, sub_member3_title, sub_member4_title
		) VALUES (
			'LİSE', '2026 - 2027', 'Okul Adı', 'Müdür Yardımcısı', 'Öğretmen', 'Öğretmen', 'Onur Kurulu 2. Başkanı', 'Öğrenci Velisi',
			'Müdür Yardımcısı', 'Öğretmen', 'Öğretmen', 'Öğretmen', 'Öğrenci Velisi'
		)`)
		if err != nil {
			return err
		}
	}

	return nil
}

// SchoolInfo methods

func (s *Store) GetSchoolInfo() (*model.SchoolInfo, error) {
	s.lock.RLock()
	defer s.lock.RUnlock()

	var info model.SchoolInfo
	row := s.db.QueryRow(`SELECT 
		id, school_type, academic_year, institution_code, province, district, school_name, principal_name,
		letterhead, official_code, address, phone, fax, warning_counter,
		board_chair, board_chair_title, member1, member1_title, member2, member2_title,
		member3, member3_title, member4, member4_title, sub_board_chair, sub_board_chair_title,
		sub_member1, sub_member1_title, sub_member2, sub_member2_title, sub_member3, sub_member3_title,
		sub_member4, sub_member4_title
		FROM school_info LIMIT 1`)

	err := row.Scan(
		&info.ID, &info.SchoolType, &info.AcademicYear, &info.InstitutionCode, &info.Province, &info.District,
		&info.SchoolName, &info.PrincipalName, &info.Letterhead, &info.OfficialCode, &info.Address,
		&info.Phone, &info.Fax, &info.WarningCounter, &info.BoardChair, &info.BoardChairTitle,
		&info.Member1, &info.Member1Title, &info.Member2, &info.Member2Title, &info.Member3,
		&info.Member3Title, &info.Member4, &info.Member4Title, &info.SubBoardChair, &info.SubBoardChairTitle,
		&info.SubMember1, &info.SubMember1Title, &info.SubMember2, &info.SubMember2Title,
		&info.SubMember3, &info.SubMember3Title, &info.SubMember4, &info.SubMember4Title,
	)
	if err != nil {
		return nil, err
	}
	return &info, nil
}

func (s *Store) SaveSchoolInfo(info *model.SchoolInfo) error {
	s.lock.Lock()
	defer s.lock.Unlock()

	_, err := s.db.Exec(`UPDATE school_info SET
		school_type = ?, academic_year = ?, institution_code = ?, province = ?, district = ?,
		school_name = ?, principal_name = ?, letterhead = ?, official_code = ?, address = ?,
		phone = ?, fax = ?, warning_counter = ?, board_chair = ?, board_chair_title = ?,
		member1 = ?, member1_title = ?, member2 = ?, member2_title = ?, member3 = ?,
		member3_title = ?, member4 = ?, member4_title = ?, sub_board_chair = ?, sub_board_chair_title = ?,
		sub_member1 = ?, sub_member1_title = ?, sub_member2 = ?, sub_member2_title = ?,
		sub_member3 = ?, sub_member3_title = ?, sub_member4 = ?, sub_member4_title = ?
		WHERE id = ?`,
		info.SchoolType, info.AcademicYear, info.InstitutionCode, info.Province, info.District,
		info.SchoolName, info.PrincipalName, info.Letterhead, info.OfficialCode, info.Address,
		info.Phone, info.Fax, info.WarningCounter, info.BoardChair, info.BoardChairTitle,
		info.Member1, info.Member1Title, info.Member2, info.Member2Title, info.Member3,
		info.Member3Title, info.Member4, info.Member4Title, info.SubBoardChair, info.SubBoardChairTitle,
		info.SubMember1, info.SubMember1Title, info.SubMember2, info.SubMember2Title,
		info.SubMember3, info.SubMember3Title, info.SubMember4, info.SubMember4Title,
		info.ID,
	)
	return err
}

// NormalizeAcademicYear standardizes a year string like "2027-2028", "2027 - 2028", "2027/2028" -> "2027 - 2028"
func NormalizeAcademicYear(raw string) string {
	raw = strings.TrimSpace(raw)
	raw = strings.ReplaceAll(raw, "/", "-")
	parts := strings.Split(raw, "-")
	if len(parts) == 2 {
		y1 := strings.TrimSpace(parts[0])
		y2 := strings.TrimSpace(parts[1])
		if len(y1) == 4 && len(y2) == 4 {
			return y1 + " - " + y2
		}
	}
	return raw
}

// ExtractStartYear parses the starting 4-digit year as an integer for chronological sorting
func ExtractStartYear(yearStr string) int {
	yearStr = strings.TrimSpace(yearStr)
	var num int
	for _, ch := range yearStr {
		if ch >= '0' && ch <= '9' {
			num = num*10 + int(ch-'0')
			if num >= 1000 && num <= 9999 {
				return num
			}
		} else if num > 0 {
			break
		}
	}
	return num
}

// GetAcademicYears returns distinct academic years sorted descending (newest first) with deduplication
func (s *Store) GetAcademicYears() []string {
	s.lock.RLock()
	defer s.lock.RUnlock()

	cfg, _ := config.LoadConfig()
	deletedMap := make(map[string]bool)
	if cfg != nil {
		for _, dy := range cfg.DeletedYears {
			deletedMap[NormalizeAcademicYear(dy)] = true
		}
	}

	yearMap := make(map[string]bool)
	defaults := []string{
		"2028 - 2029",
		"2027 - 2028",
		"2026 - 2027",
		"2025 - 2026",
		"2024 - 2025",
		"2023 - 2024",
	}
	for _, y := range defaults {
		norm := NormalizeAcademicYear(y)
		if !deletedMap[norm] {
			yearMap[norm] = true
		}
	}

	rows, err := s.db.Query("SELECT DISTINCT academic_year FROM students UNION SELECT DISTINCT academic_year FROM incidents UNION SELECT DISTINCT academic_year FROM school_info")
	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var y string
			if err := rows.Scan(&y); err == nil {
				norm := NormalizeAcademicYear(y)
				if norm != "" && !deletedMap[norm] {
					yearMap[norm] = true
				}
			}
		}
	}

	var res []string
	for y := range yearMap {
		res = append(res, y)
	}

	// Sort cleanly descending (newest year first: 2028 - 2029, 2027 - 2028, 2026 - 2027...)
	sort.Slice(res, func(i, j int) bool {
		yi := ExtractStartYear(res[i])
		yj := ExtractStartYear(res[j])
		if yi != yj {
			return yi > yj
		}
		return res[i] > res[j]
	})

	return res
}

// GetNextNearestAcademicYear calculates the next chronological academic year, e.g. "2029 - 2030"
func GetNextNearestAcademicYear(years []string) string {
	maxYear := 0
	for _, y := range years {
		sy := ExtractStartYear(y)
		if sy > maxYear {
			maxYear = sy
		}
	}
	if maxYear == 0 {
		maxYear = time.Now().Year()
	}
	return fmt.Sprintf("%d - %d", maxYear+1, maxYear+2)
}

// DeleteAcademicYear removes all records associated with a specific academic year
func (s *Store) DeleteAcademicYear(academicYear string) error {
	s.lock.Lock()
	defer s.lock.Unlock()

	norm := NormalizeAcademicYear(academicYear)
	queries := []string{
		"DELETE FROM students WHERE academic_year = ?",
		"DELETE FROM incidents WHERE academic_year = ?",
		"DELETE FROM incident_students WHERE academic_year = ?",
		"DELETE FROM meeting_decisions WHERE academic_year = ?",
		"DELETE FROM reward_penalty_counts WHERE academic_year = ?",
	}
	for _, q := range queries {
		if _, err := s.db.Exec(q, norm); err != nil {
			return err
		}
	}
	return nil
}

// Student methods

func (s *Store) GetStudents(academicYear, search, classFilter string) ([]model.Student, error) {
	s.lock.RLock()
	defer s.lock.RUnlock()

	query := "SELECT id, academic_year, school_no, full_name, class_section, gender, tc_no, birth_place, birth_date, father_name, mother_name, guardian_name, guardian_phone, guardian_addr, boarding_status, behavior_score, photo_path FROM students WHERE academic_year = ?"
	args := []interface{}{academicYear}

	if classFilter != "" && classFilter != "TÜM SINIFLAR" {
		query += " AND class_section = ?"
		args = append(args, classFilter)
	}
	if strings.TrimSpace(search) != "" {
		sTrim := "%" + strings.TrimSpace(search) + "%"
		query += " AND (school_no LIKE ? OR full_name LIKE ? OR tc_no LIKE ?)"
		args = append(args, sTrim, sTrim, sTrim)
	}
	query += " ORDER BY class_section ASC, CAST(school_no AS INTEGER) ASC, full_name ASC"

	rows, err := s.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []model.Student
	for rows.Next() {
		var st model.Student
		if err := rows.Scan(
			&st.ID, &st.AcademicYear, &st.SchoolNo, &st.FullName, &st.ClassSection, &st.Gender,
			&st.TCNo, &st.BirthPlace, &st.BirthDate, &st.FatherName, &st.MotherName,
			&st.GuardianName, &st.GuardianPhone, &st.GuardianAddr, &st.BoardingStatus,
			&st.BehaviorScore, &st.PhotoPath,
		); err != nil {
			return nil, err
		}
		list = append(list, st)
	}
	return list, nil
}

func (s *Store) GetStudentByNo(academicYear, schoolNo string) (*model.Student, error) {
	s.lock.RLock()
	defer s.lock.RUnlock()

	var st model.Student
	err := s.db.QueryRow(`SELECT id, academic_year, school_no, full_name, class_section, gender, tc_no, birth_place, birth_date, father_name, mother_name, guardian_name, guardian_phone, guardian_addr, boarding_status, behavior_score, photo_path 
		FROM students WHERE academic_year = ? AND school_no = ?`, academicYear, schoolNo).Scan(
		&st.ID, &st.AcademicYear, &st.SchoolNo, &st.FullName, &st.ClassSection, &st.Gender,
		&st.TCNo, &st.BirthPlace, &st.BirthDate, &st.FatherName, &st.MotherName,
		&st.GuardianName, &st.GuardianPhone, &st.GuardianAddr, &st.BoardingStatus,
		&st.BehaviorScore, &st.PhotoPath,
	)
	if err != nil {
		return nil, err
	}
	return &st, nil
}

func (s *Store) SaveStudent(st *model.Student) error {
	s.lock.Lock()
	defer s.lock.Unlock()

	if st.BehaviorScore <= 0 {
		st.BehaviorScore = 100
	}

	_, err := s.db.Exec(`INSERT INTO students (
		academic_year, school_no, full_name, class_section, gender, tc_no,
		birth_place, birth_date, father_name, mother_name, guardian_name,
		guardian_phone, guardian_addr, boarding_status, behavior_score, photo_path
	) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	ON CONFLICT(academic_year, school_no) DO UPDATE SET
		full_name = excluded.full_name,
		class_section = excluded.class_section,
		gender = excluded.gender,
		tc_no = excluded.tc_no,
		birth_place = excluded.birth_place,
		birth_date = excluded.birth_date,
		father_name = excluded.father_name,
		mother_name = excluded.mother_name,
		guardian_name = excluded.guardian_name,
		guardian_phone = excluded.guardian_phone,
		guardian_addr = excluded.guardian_addr,
		boarding_status = excluded.boarding_status,
		behavior_score = excluded.behavior_score,
		photo_path = excluded.photo_path`,
		st.AcademicYear, st.SchoolNo, st.FullName, st.ClassSection, st.Gender, st.TCNo,
		st.BirthPlace, st.BirthDate, st.FatherName, st.MotherName, st.GuardianName,
		st.GuardianPhone, st.GuardianAddr, st.BoardingStatus, st.BehaviorScore, st.PhotoPath,
	)
	return err
}

func (s *Store) DeleteStudent(id int64) error {
	s.lock.Lock()
	defer s.lock.Unlock()
	_, err := s.db.Exec("DELETE FROM students WHERE id = ?", id)
	return err
}

func (s *Store) DeleteAllStudents(academicYear string) error {
	s.lock.Lock()
	defer s.lock.Unlock()
	_, err := s.db.Exec("DELETE FROM students WHERE academic_year = ?", academicYear)
	return err
}

func (s *Store) GetClasses(academicYear string) []string {
	s.lock.RLock()
	defer s.lock.RUnlock()

	rows, err := s.db.Query("SELECT DISTINCT class_section FROM students WHERE academic_year = ? AND class_section != '' ORDER BY class_section ASC", academicYear)
	if err != nil {
		return nil
	}
	defer rows.Close()

	var res []string
	for rows.Next() {
		var c string
		if err := rows.Scan(&c); err == nil {
			res = append(res, c)
		}
	}
	return res
}

// Incident methods

func (s *Store) GetIncidents(academicYear string) ([]model.Incident, error) {
	s.lock.RLock()
	defer s.lock.RUnlock()

	rows, err := s.db.Query(`SELECT i.id, i.academic_year, i.incident_code, i.title, i.location, i.date, i.time, i.complainant, i.petition_ref, i.summary,
		COUNT(s.id) as student_count
		FROM incidents i
		LEFT JOIN incident_students s ON i.academic_year = s.academic_year AND i.incident_code = s.incident_code
		WHERE i.academic_year = ?
		GROUP BY i.id
		ORDER BY CAST(i.incident_code AS INTEGER) DESC`, academicYear)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []model.Incident
	for rows.Next() {
		var inc model.Incident
		if err := rows.Scan(
			&inc.ID, &inc.AcademicYear, &inc.IncidentCode, &inc.Title, &inc.Location,
			&inc.Date, &inc.Time, &inc.Complainant, &inc.PetitionRef, &inc.Summary,
			&inc.StudentCount,
		); err != nil {
			return nil, err
		}
		list = append(list, inc)
	}
	return list, nil
}

func (s *Store) GetIncidentByCode(academicYear, incidentCode string) (*model.Incident, error) {
	s.lock.RLock()
	defer s.lock.RUnlock()

	var inc model.Incident
	err := s.db.QueryRow(`SELECT id, academic_year, incident_code, title, location, date, time, complainant, petition_ref, summary
		FROM incidents WHERE academic_year = ? AND incident_code = ?`, academicYear, incidentCode).Scan(
		&inc.ID, &inc.AcademicYear, &inc.IncidentCode, &inc.Title, &inc.Location,
		&inc.Date, &inc.Time, &inc.Complainant, &inc.PetitionRef, &inc.Summary,
	)
	if err != nil {
		return nil, err
	}
	return &inc, nil
}

func (s *Store) GetNextIncidentCode(academicYear string) string {
	s.lock.RLock()
	defer s.lock.RUnlock()

	var maxCode sql.NullInt64
	_ = s.db.QueryRow("SELECT MAX(CAST(incident_code AS INTEGER)) FROM incidents WHERE academic_year = ?", academicYear).Scan(&maxCode)
	next := 1
	if maxCode.Valid {
		next = int(maxCode.Int64) + 1
	}
	return strconv.Itoa(next)
}

func (s *Store) SaveIncident(inc *model.Incident) error {
	s.lock.Lock()
	defer s.lock.Unlock()

	if strings.TrimSpace(inc.IncidentCode) == "" {
		inc.IncidentCode = s.GetNextIncidentCode(inc.AcademicYear)
	}

	_, err := s.db.Exec(`INSERT INTO incidents (
		academic_year, incident_code, title, location, date, time, complainant, petition_ref, summary
	) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
	ON CONFLICT(academic_year, incident_code) DO UPDATE SET
		title = excluded.title,
		location = excluded.location,
		date = excluded.date,
		time = excluded.time,
		complainant = excluded.complainant,
		petition_ref = excluded.petition_ref,
		summary = excluded.summary`,
		inc.AcademicYear, inc.IncidentCode, inc.Title, inc.Location, inc.Date, inc.Time,
		inc.Complainant, inc.PetitionRef, inc.Summary,
	)
	return err
}

func (s *Store) DeleteIncident(academicYear, incidentCode string) error {
	s.lock.Lock()
	defer s.lock.Unlock()

	// Cascade delete
	_, _ = s.db.Exec("DELETE FROM incident_students WHERE academic_year = ? AND incident_code = ?", academicYear, incidentCode)
	_, err := s.db.Exec("DELETE FROM incidents WHERE academic_year = ? AND incident_code = ?", academicYear, incidentCode)
	return err
}

// Incident Students & Verdicts

func (s *Store) GetIncidentStudents(academicYear, incidentCode string) ([]model.IncidentStudent, error) {
	s.lock.RLock()
	defer s.lock.RUnlock()

	rows, err := s.db.Query(`SELECT id, academic_year, incident_code, incident_title, role, school_no, full_name, class_section,
		gender, tc_no, birth_place, birth_date, father_name, mother_name, guardian_name, guardian_phone, guardian_addr,
		boarding_status, academic_standing, health_status, economic_status, living_with, parents_living_status, parents_bio_status,
		school_commute, environment_raised, prior_penalty_status, graduated_school, statement_text, written_warning, warning_date,
		warning_behavior, warning_number, warning_penalty, warning_clause, mother_age, mother_education, mother_alive, mother_bio,
		mother_together, father_age, father_education, father_alive, father_bio, father_together, sibling_count, offense_place,
		offense_date_time, offense_category, offense_reason, defense_summary, witness_summary, evidence, mitigating_factors,
		aggravating_factors, board_opinion, clause_code, clause_penalty, clause_description, imposed_penalty, suspension_days,
		suspension_dates, deducted_points, remaining_points, verdict_text, verdict_no, verdict_date, restoration_status, restoration_date
		FROM incident_students
		WHERE academic_year = ? AND incident_code = ?
		ORDER BY role ASC, full_name ASC`, academicYear, incidentCode)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []model.IncidentStudent
	for rows.Next() {
		var is model.IncidentStudent
		if err := rows.Scan(
			&is.ID, &is.AcademicYear, &is.IncidentCode, &is.IncidentTitle, &is.Role, &is.SchoolNo, &is.FullName, &is.ClassSection,
			&is.Gender, &is.TCNo, &is.BirthPlace, &is.BirthDate, &is.FatherName, &is.MotherName, &is.GuardianName, &is.GuardianPhone, &is.GuardianAddr,
			&is.BoardingStatus, &is.AcademicStanding, &is.HealthStatus, &is.EconomicStatus, &is.LivingWith, &is.ParentsLivingStatus, &is.ParentsBioStatus,
			&is.SchoolCommute, &is.EnvironmentRaised, &is.PriorPenaltyStatus, &is.GraduatedSchool, &is.StatementText, &is.WrittenWarning, &is.WarningDate,
			&is.WarningBehavior, &is.WarningNumber, &is.WarningPenalty, &is.WarningClause, &is.MotherAge, &is.MotherEducation, &is.MotherAlive, &is.MotherBio,
			&is.MotherTogether, &is.FatherAge, &is.FatherEducation, &is.FatherAlive, &is.FatherBio, &is.FatherTogether, &is.SiblingCount, &is.OffensePlace,
			&is.OffenseDateTime, &is.OffenseCategory, &is.OffenseReason, &is.DefenseSummary, &is.WitnessSummary, &is.Evidence, &is.MitigatingFactors,
			&is.AggravatingFactors, &is.BoardOpinion, &is.ClauseCode, &is.ClausePenalty, &is.ClauseDescription, &is.ImposedPenalty, &is.SuspensionDays,
			&is.SuspensionDates, &is.DeductedPoints, &is.RemainingPoints, &is.VerdictText, &is.VerdictNo, &is.VerdictDate, &is.RestorationStatus, &is.RestorationDate,
		); err != nil {
			return nil, err
		}
		list = append(list, is)
	}
	return list, nil
}

func (s *Store) GetIncidentStudentByID(id int64) (*model.IncidentStudent, error) {
	s.lock.RLock()
	defer s.lock.RUnlock()

	var is model.IncidentStudent
	err := s.db.QueryRow(`SELECT id, academic_year, incident_code, incident_title, role, school_no, full_name, class_section,
		gender, tc_no, birth_place, birth_date, father_name, mother_name, guardian_name, guardian_phone, guardian_addr,
		boarding_status, academic_standing, health_status, economic_status, living_with, parents_living_status, parents_bio_status,
		school_commute, environment_raised, prior_penalty_status, graduated_school, statement_text, written_warning, warning_date,
		warning_behavior, warning_number, warning_penalty, warning_clause, mother_age, mother_education, mother_alive, mother_bio,
		mother_together, father_age, father_education, father_alive, father_bio, father_together, sibling_count, offense_place,
		offense_date_time, offense_category, offense_reason, defense_summary, witness_summary, evidence, mitigating_factors,
		aggravating_factors, board_opinion, clause_code, clause_penalty, clause_description, imposed_penalty, suspension_days,
		suspension_dates, deducted_points, remaining_points, verdict_text, verdict_no, verdict_date, restoration_status, restoration_date
		FROM incident_students WHERE id = ?`, id).Scan(
		&is.ID, &is.AcademicYear, &is.IncidentCode, &is.IncidentTitle, &is.Role, &is.SchoolNo, &is.FullName, &is.ClassSection,
		&is.Gender, &is.TCNo, &is.BirthPlace, &is.BirthDate, &is.FatherName, &is.MotherName, &is.GuardianName, &is.GuardianPhone, &is.GuardianAddr,
		&is.BoardingStatus, &is.AcademicStanding, &is.HealthStatus, &is.EconomicStatus, &is.LivingWith, &is.ParentsLivingStatus, &is.ParentsBioStatus,
		&is.SchoolCommute, &is.EnvironmentRaised, &is.PriorPenaltyStatus, &is.GraduatedSchool, &is.StatementText, &is.WrittenWarning, &is.WarningDate,
		&is.WarningBehavior, &is.WarningNumber, &is.WarningPenalty, &is.WarningClause, &is.MotherAge, &is.MotherEducation, &is.MotherAlive, &is.MotherBio,
		&is.MotherTogether, &is.FatherAge, &is.FatherEducation, &is.FatherAlive, &is.FatherBio, &is.FatherTogether, &is.SiblingCount, &is.OffensePlace,
		&is.OffenseDateTime, &is.OffenseCategory, &is.OffenseReason, &is.DefenseSummary, &is.WitnessSummary, &is.Evidence, &is.MitigatingFactors,
		&is.AggravatingFactors, &is.BoardOpinion, &is.ClauseCode, &is.ClausePenalty, &is.ClauseDescription, &is.ImposedPenalty, &is.SuspensionDays,
		&is.SuspensionDates, &is.DeductedPoints, &is.RemainingPoints, &is.VerdictText, &is.VerdictNo, &is.VerdictDate, &is.RestorationStatus, &is.RestorationDate,
	)
	if err != nil {
		return nil, err
	}
	return &is, nil
}

func (s *Store) SaveIncidentStudent(is *model.IncidentStudent) error {
	s.lock.Lock()
	defer s.lock.Unlock()

	if is.ID > 0 {
		_, err := s.db.Exec(`UPDATE incident_students SET
			incident_title = ?, role = ?, school_no = ?, full_name = ?, class_section = ?, gender = ?,
			tc_no = ?, birth_place = ?, birth_date = ?, father_name = ?, mother_name = ?, guardian_name = ?,
			guardian_phone = ?, guardian_addr = ?, boarding_status = ?, academic_standing = ?, health_status = ?,
			economic_status = ?, living_with = ?, parents_living_status = ?, parents_bio_status = ?,
			school_commute = ?, environment_raised = ?, prior_penalty_status = ?, graduated_school = ?,
			statement_text = ?, written_warning = ?, warning_date = ?, warning_behavior = ?, warning_number = ?,
			warning_penalty = ?, warning_clause = ?, mother_age = ?, mother_education = ?, mother_alive = ?,
			mother_bio = ?, mother_together = ?, father_age = ?, father_education = ?, father_alive = ?,
			father_bio = ?, father_together = ?, sibling_count = ?, offense_place = ?, offense_date_time = ?,
			offense_category = ?, offense_reason = ?, defense_summary = ?, witness_summary = ?, evidence = ?,
			mitigating_factors = ?, aggravating_factors = ?, board_opinion = ?, clause_code = ?,
			clause_penalty = ?, clause_description = ?, imposed_penalty = ?, suspension_days = ?,
			suspension_dates = ?, deducted_points = ?, remaining_points = ?, verdict_text = ?,
			verdict_no = ?, verdict_date = ?, restoration_status = ?, restoration_date = ?
			WHERE id = ?`,
			is.IncidentTitle, is.Role, is.SchoolNo, is.FullName, is.ClassSection, is.Gender,
			is.TCNo, is.BirthPlace, is.BirthDate, is.FatherName, is.MotherName, is.GuardianName,
			is.GuardianPhone, is.GuardianAddr, is.BoardingStatus, is.AcademicStanding, is.HealthStatus,
			is.EconomicStatus, is.LivingWith, is.ParentsLivingStatus, is.ParentsBioStatus,
			is.SchoolCommute, is.EnvironmentRaised, is.PriorPenaltyStatus, is.GraduatedSchool,
			is.StatementText, is.WrittenWarning, is.WarningDate, is.WarningBehavior, is.WarningNumber,
			is.WarningPenalty, is.WarningClause, is.MotherAge, is.MotherEducation, is.MotherAlive,
			is.MotherBio, is.MotherTogether, is.FatherAge, is.FatherEducation, is.FatherAlive,
			is.FatherBio, is.FatherTogether, is.SiblingCount, is.OffensePlace, is.OffenseDateTime,
			is.OffenseCategory, is.OffenseReason, is.DefenseSummary, is.WitnessSummary, is.Evidence,
			is.MitigatingFactors, is.AggravatingFactors, is.BoardOpinion, is.ClauseCode,
			is.ClausePenalty, is.ClauseDescription, is.ImposedPenalty, is.SuspensionDays,
			is.SuspensionDates, is.DeductedPoints, is.RemainingPoints, is.VerdictText,
			is.VerdictNo, is.VerdictDate, is.RestorationStatus, is.RestorationDate,
			is.ID,
		)
		return err
	}

	res, err := s.db.Exec(`INSERT INTO incident_students (
		academic_year, incident_code, incident_title, role, school_no, full_name, class_section, gender,
		tc_no, birth_place, birth_date, father_name, mother_name, guardian_name, guardian_phone, guardian_addr,
		boarding_status, academic_standing, health_status, economic_status, living_with, parents_living_status, parents_bio_status,
		school_commute, environment_raised, prior_penalty_status, graduated_school, statement_text, written_warning, warning_date,
		warning_behavior, warning_number, warning_penalty, warning_clause, mother_age, mother_education, mother_alive, mother_bio,
		mother_together, father_age, father_education, father_alive, father_bio, father_together, sibling_count, offense_place,
		offense_date_time, offense_category, offense_reason, defense_summary, witness_summary, evidence, mitigating_factors,
		aggravating_factors, board_opinion, clause_code, clause_penalty, clause_description, imposed_penalty, suspension_days,
		suspension_dates, deducted_points, remaining_points, verdict_text, verdict_no, verdict_date, restoration_status, restoration_date
	) VALUES (
		?, ?, ?, ?, ?, ?, ?, ?,
		?, ?, ?, ?, ?, ?, ?, ?,
		?, ?, ?, ?, ?, ?, ?,
		?, ?, ?, ?, ?, ?, ?,
		?, ?, ?, ?, ?, ?, ?, ?,
		?, ?, ?, ?, ?, ?, ?, ?,
		?, ?, ?, ?, ?, ?, ?,
		?, ?, ?, ?, ?, ?, ?,
		?, ?, ?, ?, ?, ?, ?, ?
	)`,
		is.AcademicYear, is.IncidentCode, is.IncidentTitle, is.Role, is.SchoolNo, is.FullName, is.ClassSection, is.Gender,
		is.TCNo, is.BirthPlace, is.BirthDate, is.FatherName, is.MotherName, is.GuardianName, is.GuardianPhone, is.GuardianAddr,
		is.BoardingStatus, is.AcademicStanding, is.HealthStatus, is.EconomicStatus, is.LivingWith, is.ParentsLivingStatus, is.ParentsBioStatus,
		is.SchoolCommute, is.EnvironmentRaised, is.PriorPenaltyStatus, is.GraduatedSchool, is.StatementText, is.WrittenWarning, is.WarningDate,
		is.WarningBehavior, is.WarningNumber, is.WarningPenalty, is.WarningClause, is.MotherAge, is.MotherEducation, is.MotherAlive, is.MotherBio,
		is.MotherTogether, is.FatherAge, is.FatherEducation, is.FatherAlive, is.FatherBio, is.FatherTogether, is.SiblingCount, is.OffensePlace,
		is.OffenseDateTime, is.OffenseCategory, is.OffenseReason, is.DefenseSummary, is.WitnessSummary, is.Evidence, is.MitigatingFactors,
		is.AggravatingFactors, is.BoardOpinion, is.ClauseCode, is.ClausePenalty, is.ClauseDescription, is.ImposedPenalty, is.SuspensionDays,
		is.SuspensionDates, is.DeductedPoints, is.RemainingPoints, is.VerdictText, is.VerdictNo, is.VerdictDate, is.RestorationStatus, is.RestorationDate,
	)
	if err != nil {
		return err
	}
	is.ID, _ = res.LastInsertId()
	return nil
}

func (s *Store) DeleteIncidentStudent(id int64) error {
	s.lock.Lock()
	defer s.lock.Unlock()
	_, err := s.db.Exec("DELETE FROM incident_students WHERE id = ?", id)
	return err
}

func (s *Store) GetNextVerdictNo(academicYear string) string {
	s.lock.RLock()
	defer s.lock.RUnlock()

	var maxNo sql.NullInt64
	_ = s.db.QueryRow("SELECT MAX(CAST(verdict_no AS INTEGER)) FROM incident_students WHERE academic_year = ?", academicYear).Scan(&maxNo)
	next := 1
	if maxNo.Valid {
		next = int(maxNo.Int64) + 1
	}
	return strconv.Itoa(next)
}

func (s *Store) GetPriorPenalties(studentName string) string {
	s.lock.RLock()
	defer s.lock.RUnlock()

	rows, err := s.db.Query("SELECT academic_year, imposed_penalty FROM incident_students WHERE full_name = ? AND verdict_no != '' AND imposed_penalty != ''", studentName)
	if err != nil {
		return ""
	}
	defer rows.Close()

	var sb strings.Builder
	for rows.Next() {
		var year, penalty string
		if err := rows.Scan(&year, &penalty); err == nil {
			sb.WriteString(fmt.Sprintf("%s Öğretim Yılı %s Cezası ", year, penalty))
		}
	}
	return sb.String()
}

// GetIncidentStudentsBySchoolNo returns all incident records for a specific student
func (s *Store) GetIncidentStudentsBySchoolNo(academicYear, schoolNo string) ([]model.IncidentStudent, error) {
	s.lock.RLock()
	defer s.lock.RUnlock()

	rows, err := s.db.Query(`SELECT id, academic_year, incident_code, incident_title, role, school_no, full_name,
		class_section, statement_text, imposed_penalty, verdict_no, verdict_date
		FROM incident_students
		WHERE academic_year = ? AND school_no = ?
		ORDER BY id DESC`, academicYear, schoolNo)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []model.IncidentStudent
	for rows.Next() {
		var is model.IncidentStudent
		if err := rows.Scan(&is.ID, &is.AcademicYear, &is.IncidentCode, &is.IncidentTitle, &is.Role,
			&is.SchoolNo, &is.FullName, &is.ClassSection, &is.StatementText, &is.ImposedPenalty,
			&is.VerdictNo, &is.VerdictDate); err == nil {
			list = append(list, is)
		}
	}
	return list, nil
}

// GetPenalizedStudents returns all students who received a verdict with imposed penalty
func (s *Store) GetPenalizedStudents(academicYear, classFilter, penaltyFilter string) ([]model.IncidentStudent, error) {
	s.lock.RLock()
	defer s.lock.RUnlock()

	query := `SELECT id, academic_year, incident_code, incident_title, role, school_no, full_name, class_section,
		gender, tc_no, birth_place, birth_date, father_name, mother_name, guardian_name, guardian_phone, guardian_addr,
		boarding_status, academic_standing, health_status, economic_status, living_with, parents_living_status, parents_bio_status,
		school_commute, environment_raised, prior_penalty_status, graduated_school, statement_text, written_warning, warning_date,
		warning_behavior, warning_number, warning_penalty, warning_clause, mother_age, mother_education, mother_alive, mother_bio,
		mother_together, father_age, father_education, father_alive, father_bio, father_together, sibling_count, offense_place,
		offense_date_time, offense_category, offense_reason, defense_summary, witness_summary, evidence, mitigating_factors,
		aggravating_factors, board_opinion, clause_code, clause_penalty, clause_description, imposed_penalty, suspension_days,
		suspension_dates, deducted_points, remaining_points, verdict_text, verdict_no, verdict_date, restoration_status, restoration_date
		FROM incident_students
		WHERE academic_year = ? AND verdict_no != '' AND imposed_penalty != ''`
	args := []interface{}{academicYear}

	if classFilter != "" && classFilter != "TÜM SINIFLAR" {
		query += " AND class_section = ?"
		args = append(args, classFilter)
	}
	if penaltyFilter != "" && penaltyFilter != "TÜM CEZALAR" {
		query += " AND imposed_penalty = ?"
		args = append(args, penaltyFilter)
	}
	query += " ORDER BY CAST(verdict_no AS INTEGER) ASC"

	rows, err := s.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []model.IncidentStudent
	for rows.Next() {
		var is model.IncidentStudent
		if err := rows.Scan(
			&is.ID, &is.AcademicYear, &is.IncidentCode, &is.IncidentTitle, &is.Role, &is.SchoolNo, &is.FullName, &is.ClassSection,
			&is.Gender, &is.TCNo, &is.BirthPlace, &is.BirthDate, &is.FatherName, &is.MotherName, &is.GuardianName, &is.GuardianPhone, &is.GuardianAddr,
			&is.BoardingStatus, &is.AcademicStanding, &is.HealthStatus, &is.EconomicStatus, &is.LivingWith, &is.ParentsLivingStatus, &is.ParentsBioStatus,
			&is.SchoolCommute, &is.EnvironmentRaised, &is.PriorPenaltyStatus, &is.GraduatedSchool, &is.StatementText, &is.WrittenWarning, &is.WarningDate,
			&is.WarningBehavior, &is.WarningNumber, &is.WarningPenalty, &is.WarningClause, &is.MotherAge, &is.MotherEducation, &is.MotherAlive, &is.MotherBio,
			&is.MotherTogether, &is.FatherAge, &is.FatherEducation, &is.FatherAlive, &is.FatherBio, &is.FatherTogether, &is.SiblingCount, &is.OffensePlace,
			&is.OffenseDateTime, &is.OffenseCategory, &is.OffenseReason, &is.DefenseSummary, &is.WitnessSummary, &is.Evidence, &is.MitigatingFactors,
			&is.AggravatingFactors, &is.BoardOpinion, &is.ClauseCode, &is.ClausePenalty, &is.ClauseDescription, &is.ImposedPenalty, &is.SuspensionDays,
			&is.SuspensionDates, &is.DeductedPoints, &is.RemainingPoints, &is.VerdictText, &is.VerdictNo, &is.VerdictDate, &is.RestorationStatus, &is.RestorationDate,
		); err != nil {
			return nil, err
		}
		list = append(list, is)
	}
	return list, nil
}

// CalculateSemesterPenaltyCounts computes 1st and 2nd semester counts based on the January 30 boundary
func (s *Store) CalculateSemesterPenaltyCounts(academicYear string) ([]model.RewardPenaltyCount, error) {
	s.lock.RLock()
	defer s.lock.RUnlock()

	parts := strings.Split(academicYear, "-")
	y2 := time.Now().Year()
	if len(parts) >= 2 {
		if val, err := strconv.Atoi(strings.TrimSpace(parts[1])); err == nil {
			y2 = val
		}
	}
	cutoff := time.Date(y2, 1, 30, 23, 59, 59, 0, time.Local)

	penalties := []string{
		"SÖZLÜ OLARAK UYARMA",
		"YAZILI OLARAK UYARMA",
		"TOPLUM HİZMETİ",
		"KINAMA",
		"OKULDAN KISA SÜRELİ UZAKLAŞTIRMA",
		"OKUL DEĞİŞTİRME",
		"ÖRGÜN EĞİTİM DIŞINA ÇIKARMA",
		"UYARMA",
	}

	var results []model.RewardPenaltyCount
	for _, pen := range penalties {
		rows, err := s.db.Query("SELECT verdict_date, warning_date, written_warning FROM incident_students WHERE academic_year = ? AND (imposed_penalty = ? OR (warning_penalty = ? AND written_warning = 'ALDI'))", academicYear, pen, pen)
		if err != nil {
			continue
		}
		sem1, sem2 := 0, 0
		for rows.Next() {
			var vDate, wDate, wAlert string
			_ = rows.Scan(&vDate, &wDate, &wAlert)
			targetDateStr := vDate
			if pen == "YAZILI OLARAK UYARMA" && wDate != "" {
				targetDateStr = wDate
			}
			t, errParse := time.Parse("02.01.2006", targetDateStr)
			if errParse != nil {
				t, errParse = time.Parse("2006-01-02", targetDateStr)
			}
			if errParse == nil {
				if t.Before(cutoff) || t.Equal(cutoff) {
					sem1++
				} else {
					sem2++
				}
			} else {
				sem1++
			}
		}
		rows.Close()
		results = append(results, model.RewardPenaltyCount{
			AcademicYear: academicYear,
			Kind:         "ceza",
			Name:         pen,
			Sem1Count:    sem1,
			Sem2Count:    sem2,
			TotalCount:   sem1 + sem2,
		})
	}
	return results, nil
}

// Agenda items methods

func (s *Store) GetAgendaItems() ([]model.AgendaItem, error) {
	s.lock.RLock()
	defer s.lock.RUnlock()

	rows, err := s.db.Query("SELECT id, order_no, item_text FROM agenda_items ORDER BY order_no ASC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []model.AgendaItem
	for rows.Next() {
		var it model.AgendaItem
		if err := rows.Scan(&it.ID, &it.OrderNo, &it.ItemText); err == nil {
			list = append(list, it)
		}
	}
	return list, nil
}

func (s *Store) SaveAgendaItems(items []model.AgendaItem) error {
	s.lock.Lock()
	defer s.lock.Unlock()

	_, _ = s.db.Exec("DELETE FROM agenda_items")
	for i, it := range items {
		_, _ = s.db.Exec("INSERT INTO agenda_items (order_no, item_text) VALUES (?, ?)", i+1, it.ItemText)
	}
	return nil
}

// Meeting Decisions methods

func (s *Store) GetMeetingDecisions(academicYear, meetingName string) ([]model.MeetingDecision, error) {
	s.lock.RLock()
	defer s.lock.RUnlock()

	rows, err := s.db.Query("SELECT id, academic_year, meeting_name, meeting_date, meeting_time, meeting_place, meeting_chair, decision_no, decision_text FROM meeting_decisions WHERE academic_year = ? AND meeting_name = ? ORDER BY decision_no ASC", academicYear, meetingName)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []model.MeetingDecision
	for rows.Next() {
		var d model.MeetingDecision
		if err := rows.Scan(&d.ID, &d.AcademicYear, &d.MeetingName, &d.MeetingDate, &d.MeetingTime, &d.MeetingPlace, &d.MeetingChair, &d.DecisionNo, &d.DecisionText); err == nil {
			list = append(list, d)
		}
	}
	return list, nil
}

func (s *Store) SaveMeetingDecisions(academicYear, meetingName, date, time, place, chair string, decisions []string) error {
	s.lock.Lock()
	defer s.lock.Unlock()

	_, _ = s.db.Exec("DELETE FROM meeting_decisions WHERE academic_year = ? AND meeting_name = ?", academicYear, meetingName)
	for i, text := range decisions {
		if strings.TrimSpace(text) == "" {
			continue
		}
		_, _ = s.db.Exec(`INSERT INTO meeting_decisions (
			academic_year, meeting_name, meeting_date, meeting_time, meeting_place, meeting_chair, decision_no, decision_text
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
			academicYear, meetingName, date, time, place, chair, i+1, text,
		)
	}
	return nil
}
