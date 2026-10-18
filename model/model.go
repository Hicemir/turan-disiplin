package model

import "time"

// SchoolInfo stores the general school and board settings
type SchoolInfo struct {
	ID                     int64  `json:"id"`
	SchoolType             string `json:"school_type"` // "LİSE" or "ORTAOKUL"
	AcademicYear           string `json:"academic_year"` // e.g. "2026 - 2027"
	InstitutionCode        string `json:"institution_code"`
	Province               string `json:"province"`
	District               string `json:"district"`
	SchoolName             string `json:"school_name"`
	PrincipalName          string `json:"principal_name"`
	Letterhead             string `json:"letterhead"`
	OfficialCode           string `json:"official_code"` // e.g. Sayı
	Address                string `json:"address"`
	Phone                  string `json:"phone"`
	Fax                    string `json:"fax"`
	WarningCounter         int    `json:"warning_counter"`
	BoardChair             string `json:"board_chair"`
	BoardChairTitle        string `json:"board_chair_title"`
	Member1                string `json:"member1"`
	Member1Title           string `json:"member1_title"`
	Member2                string `json:"member2"`
	Member2Title           string `json:"member2_title"`
	Member3                string `json:"member3"`
	Member3Title           string `json:"member3_title"`
	Member4                string `json:"member4"`
	Member4Title           string `json:"member4_title"`
	SubBoardChair          string `json:"sub_board_chair"`
	SubBoardChairTitle     string `json:"sub_board_chair_title"`
	SubMember1             string `json:"sub_member1"`
	SubMember1Title        string `json:"sub_member1_title"`
	SubMember2             string `json:"sub_member2"`
	SubMember2Title        string `json:"sub_member2_title"`
	SubMember3             string `json:"sub_member3"`
	SubMember3Title        string `json:"sub_member3_title"`
	SubMember4             string `json:"sub_member4"`
	SubMember4Title        string `json:"sub_member4_title"`
}

// Student represents a student record in the master roster
type Student struct {
	ID             int64  `json:"id"`
	AcademicYear   string `json:"academic_year"`
	SchoolNo       string `json:"school_no"`
	FullName       string `json:"full_name"`
	ClassSection   string `json:"class_section"`
	Gender         string `json:"gender"`
	TCNo           string `json:"tc_no"`
	BirthPlace     string `json:"birth_place"`
	BirthDate      string `json:"birth_date"`
	FatherName     string `json:"father_name"`
	MotherName     string `json:"mother_name"`
	GuardianName   string `json:"guardian_name"`
	GuardianPhone  string `json:"guardian_phone"`
	GuardianAddr   string `json:"guardian_addr"`
	BoardingStatus string `json:"boarding_status"` // "GÜNDÜZLÜ", "PANSİYON / YATILI"
	BehaviorScore  int    `json:"behavior_score"`  // Default 100
	PhotoPath      string `json:"photo_path"`
}

// Incident represents a disciplinary incident
type Incident struct {
	ID             int64  `json:"id"`
	AcademicYear   string `json:"academic_year"`
	IncidentCode   string `json:"incident_code"` // Sequential numeric string "1", "2"...
	Title          string `json:"title"`
	Location       string `json:"location"`
	Date           string `json:"date"` // dd.MM.yyyy
	Time           string `json:"time"`
	Complainant    string `json:"complainant"`
	PetitionRef    string `json:"petition_ref"`
	Summary        string `json:"summary"`
	StudentCount   int    `json:"student_count,omitempty"`
}

// IncidentStudent links a student to an incident and stores investigation / verdict details
type IncidentStudent struct {
	ID                  int64  `json:"id"`
	AcademicYear        string `json:"academic_year"`
	IncidentCode        string `json:"incident_code"`
	IncidentTitle       string `json:"incident_title"`
	Role                string `json:"role"` // "SUÇLANAN ÖĞRENCİ", "MAĞDUR / ŞİKAYETÇİ ÖĞRENCİ", "BİLGİSİNE BAŞVURULACAK ÖĞRENCİ"
	SchoolNo            string `json:"school_no"`
	FullName            string `json:"full_name"`
	ClassSection        string `json:"class_section"`
	Gender              string `json:"gender"`
	TCNo                string `json:"tc_no"`
	BirthPlace          string `json:"birth_place"`
	BirthDate           string `json:"birth_date"`
	FatherName          string `json:"father_name"`
	MotherName          string `json:"mother_name"`
	GuardianName        string `json:"guardian_name"`
	GuardianPhone       string `json:"guardian_phone"`
	GuardianAddr        string `json:"guardian_addr"`
	BoardingStatus      string `json:"boarding_status"`
	AcademicStanding    string `json:"academic_standing"`
	HealthStatus        string `json:"health_status"`
	EconomicStatus      string `json:"economic_status"`
	LivingWith          string `json:"living_with"`
	ParentsLivingStatus string `json:"parents_living_status"`
	ParentsBioStatus    string `json:"parents_bio_status"`
	SchoolCommute       string `json:"school_commute"`
	EnvironmentRaised   string `json:"environment_raised"`
	PriorPenaltyStatus  string `json:"prior_penalty_status"`
	GraduatedSchool     string `json:"graduated_school"`
	StatementText       string `json:"statement_text"`
	WrittenWarning      string `json:"written_warning"` // "ALDI" or ""
	WarningDate         string `json:"warning_date"`
	WarningBehavior     string `json:"warning_behavior"`
	WarningNumber       string `json:"warning_number"`
	WarningPenalty      string `json:"warning_penalty"`
	WarningClause       string `json:"warning_clause"`
	// Ortaokul parent fields
	MotherAge       string `json:"mother_age"`
	MotherEducation string `json:"mother_education"`
	MotherAlive     string `json:"mother_alive"`
	MotherBio       string `json:"mother_bio"`
	MotherTogether  string `json:"mother_together"`
	FatherAge       string `json:"father_age"`
	FatherEducation string `json:"father_education"`
	FatherAlive     string `json:"father_alive"`
	FatherBio       string `json:"father_bio"`
	FatherTogether  string `json:"father_together"`
	SiblingCount    string `json:"sibling_count"`
	// Verdict Fields
	OffensePlace      string `json:"offense_place"`
	OffenseDateTime   string `json:"offense_date_time"`
	OffenseCategory   string `json:"offense_category"`
	OffenseReason     string `json:"offense_reason"`
	DefenseSummary    string `json:"defense_summary"`
	WitnessSummary    string `json:"witness_summary"`
	Evidence          string `json:"evidence"`
	MitigatingFactors string `json:"mitigating_factors"`
	AggravatingFactors string `json:"aggravating_factors"`
	BoardOpinion      string `json:"board_opinion"` // "SUÇLU OLDUĞU YÖNÜNDE", etc.
	ClauseCode        string `json:"clause_code"`
	ClausePenalty     string `json:"clause_penalty"`
	ClauseDescription string `json:"clause_description"`
	ImposedPenalty    string `json:"imposed_penalty"`
	SuspensionDays    string `json:"suspension_days"`
	SuspensionDates   string `json:"suspension_dates"`
	DeductedPoints    int    `json:"deducted_points"`
	RemainingPoints   int    `json:"remaining_points"`
	VerdictText       string `json:"verdict_text"`
	VerdictNo         string `json:"verdict_no"`
	VerdictDate       string `json:"verdict_date"`
	RestorationStatus string `json:"restoration_status"` // "İADE EDİLDİ", "YAPTIRIM KALDIRILDI", etc.
	RestorationDate   string `json:"restoration_date"`
}

// AgendaItem represents a meeting agenda point
type AgendaItem struct {
	ID        int64  `json:"id"`
	OrderNo   int    `json:"order_no"`
	ItemText  string `json:"item_text"`
}

// MeetingDecision represents a decision recorded in beginning/end of year meeting
type MeetingDecision struct {
	ID            int64  `json:"id"`
	AcademicYear  string `json:"academic_year"`
	MeetingName   string `json:"meeting_name"` // "SENE BAŞI", "SENE SONU"
	MeetingDate   string `json:"meeting_date"`
	MeetingTime   string `json:"meeting_time"`
	MeetingPlace  string `json:"meeting_place"`
	MeetingChair  string `json:"meeting_chair"`
	DecisionNo    int    `json:"decision_no"`
	DecisionText  string `json:"decision_text"`
}

// RewardPenaltyCount stores 1st and 2nd semester counts for rewards and penalties
type RewardPenaltyCount struct {
	ID           int64  `json:"id"`
	AcademicYear string `json:"academic_year"`
	MeetingName  string `json:"meeting_name"`
	MeetingDate  string `json:"meeting_date"`
	MeetingTime  string `json:"meeting_time"`
	MeetingPlace string `json:"meeting_place"`
	MeetingChair string `json:"meeting_chair"`
	Kind         string `json:"kind"` // "ödül" or "ceza"
	Name         string `json:"name"`
	Sem1Count    int    `json:"sem1_count"`
	Sem2Count    int    `json:"sem2_count"`
	TotalCount   int    `json:"total_count"`
}

// Clause represents a regulation offense article
type Clause struct {
	Penalty     string `json:"ceza"`
	Article     string `json:"madde"`
	Description string `json:"açıklama"`
	SchoolType  string `json:"okultürü"`
}

// YonetmelikData is the root structure of yonetmelik.json
type YonetmelikData struct {
	Metadata               map[string]interface{}    `json:"metadata"`
	SchoolTypes            map[string]SchoolTypeMeta `json:"school_types"`
	Rewards                []string                  `json:"rewards"`
	StudentProfileOptions  map[string]interface{}    `json:"student_profile_options"`
	VerdictTemplates       map[string]map[string]string `json:"verdict_templates"`
	DocumentTemplates      map[string]interface{}    `json:"document_templates"`
	Clauses                []Clause                  `json:"cezamaddeleri"`
}

type SchoolTypeMeta struct {
	Code                   string      `json:"code"`
	BoardName              string      `json:"board_name"`
	BoardFullTitle         string      `json:"board_full_title"`
	RegulationName         string      `json:"regulation_name"`
	OffenseArticleNumber   string      `json:"offense_article_number"`
	BoardEstablishment     interface{} `json:"board_establishment_article"`
	OffenseCategories      []string    `json:"offense_categories"`
	BoardOpinions          []string    `json:"board_opinions"`
	Penalties              []PenaltyMeta `json:"penalties"`
	InitialBehaviorPoints  int         `json:"initial_behavior_points"`
}

type PenaltyMeta struct {
	Name                   string   `json:"name"`
	BehaviorPointDeduction int      `json:"behavior_point_deduction"`
	RequiresSuspensionDays bool     `json:"requires_suspension_days"`
	SuspensionDayOptions   []string `json:"suspension_day_options,omitempty"`
	IsRegulationCategory   bool     `json:"is_regulation_clause_category"`
}

// CurrentDateFormatted returns dd.MM.yyyy
func CurrentDateFormatted() string {
	return time.Now().Format("02.01.2006")
}

// CurrentTimeFormatted returns HH:mm
func CurrentTimeFormatted() string {
	return time.Now().Format("15:04")
}

type DiziItem struct {
	Subject string `json:"subject"`
	Count   int    `json:"count"`
}

type StatItem struct {
	Name  string `json:"name"`
	Count int    `json:"count"`
}

type BehaviorScoreRefundRow struct {
	FullName        string `json:"full_name"`
	SchoolNo        string `json:"school_no"`
	ClassSection    string `json:"class_section"`
	IncidentName    string `json:"incident_name"`
	DecisionNo      int    `json:"decision_no"`
	DecisionDate    string `json:"decision_date"`
	PenaltyName     string `json:"penalty_name"`
	RegClause       string `json:"reg_clause"`
	DeductedPoints  int    `json:"deducted_points"`
	RemainingPoints int    `json:"remaining_points"`
	RefundPoints    int    `json:"refund_points"`
	RefundDate      string `json:"refund_date"`
	Status          string `json:"status"`
}

type SanctionLiftRow struct {
	FullName     string `json:"full_name"`
	SchoolNo     string `json:"school_no"`
	ClassSection string `json:"class_section"`
	IncidentName string `json:"incident_name"`
	DecisionNo   int    `json:"decision_no"`
	DecisionDate string `json:"decision_date"`
	PenaltyName  string `json:"penalty_name"`
	RegClause    string `json:"reg_clause"`
	LiftDate     string `json:"lift_date"`
	Status       string `json:"status"`
}

type Verdict struct {
	DecisionNo         int    `json:"decision_no"`
	DecisionDate       string `json:"decision_date"`
	OffensePlace       string `json:"offense_place"`
	OffenseDateTime    string `json:"offense_date_time"`
	OffenseType        string `json:"offense_type"`
	OffenseReason      string `json:"offense_reason"`
	WitnessStatements  string `json:"witness_statements"`
	EvidenceDetails    string `json:"evidence_details"`
	MitigatingFactors  string `json:"mitigating_factors"`
	AggravatingFactors string `json:"aggravating_factors"`
	BoardOpinion       string `json:"board_opinion"`
	RegClause          string `json:"reg_clause"`
	RegClauseDesc      string `json:"reg_clause_desc"`
	PenaltyName        string `json:"penalty_name"`
	SuspensionDays     int    `json:"suspension_days"`
	SuspensionDates    string `json:"suspension_dates"`
	DeductedPoints     int    `json:"deducted_points"`
	RemainingPoints    int    `json:"remaining_points"`
	DecisionText       string `json:"decision_text"`
	ChairIsSubstitute  bool   `json:"chair_is_substitute"`
	M1IsSubstitute     bool   `json:"m1_is_substitute"`
	M2IsSubstitute     bool   `json:"m2_is_substitute"`
	M3IsSubstitute     bool   `json:"m3_is_substitute"`
	M4IsSubstitute     bool   `json:"m4_is_substitute"`
}
