package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"

	"okulevrak/config"
	"okulevrak/model"
	"okulevrak/store"
	appTheme "okulevrak/theme"
)

// AppState stores the active state across the application
type AppState struct {
	App          fyne.App
	Window       fyne.Window
	Store        *store.Store
	Config       *config.AppConfig
	Yonetmelik   *model.YonetmelikData
	SchoolInfo   *model.SchoolInfo
	AcademicYear string

	// Navigation State
	CurrentView string // "hub", "school", "students", "incidents", "verdict", "summons", "forms", "notices", "restoration", "minutes", "stats", "yonetmelik", "settings", "dizi"

	// Active Selections
	SelectedStudentNo    string
	SelectedIncidentCode string
	SelectedIncStudentID int64

	// View refresh callback (clean, non-recursive, deadlock-free)
	OnViewRefresh func()
}

var GlobalState *AppState

// InitState initializes the shared application state
func InitState(a fyne.App, w fyne.Window) (*AppState, error) {
	st, err := store.GetStore()
	if err != nil {
		return nil, err
	}

	cfg, err := config.LoadConfig()
	if err != nil {
		return nil, err
	}

	ym, err := config.LoadYonetmelik()
	if err != nil {
		return nil, err
	}

	info, err := st.GetSchoolInfo()
	if err != nil {
		return nil, err
	}

	currentYear := cfg.DefaultYear
	if currentYear == "" {
		if info.AcademicYear != "" {
			currentYear = info.AcademicYear
		} else {
			currentYear = "2026 - 2027"
		}
	}

	GlobalState = &AppState{
		App:          a,
		Window:       w,
		Store:        st,
		Config:       cfg,
		Yonetmelik:   ym,
		SchoolInfo:   info,
		AcademicYear: currentYear,
		CurrentView:  "hub",
	}

	// Apply configured theme on startup
	if a != nil {
		tMode := cfg.ThemeMode
		if tMode == "" {
			// First run: detect OS system theme variant
			if a.Settings().ThemeVariant() == theme.VariantDark {
				tMode = "dark"
			} else {
				tMode = "light"
			}
			cfg.ThemeMode = tMode
			_ = config.SaveConfig(cfg)
		}
		a.Settings().SetTheme(appTheme.NewTurkishTheme(tMode, cfg.FontSize))
	}

	return GlobalState, nil
}

// NavigateTo switches the visible functional module cleanly
func (s *AppState) NavigateTo(view string) {
	s.CurrentView = view
	if s.OnViewRefresh != nil {
		s.OnViewRefresh()
	}
}

// SetAcademicYear changes the active year and updates DB & Config
func (s *AppState) SetAcademicYear(year string) {
	if year == "" || year == s.AcademicYear {
		return
	}
	s.AcademicYear = year
	s.Config.DefaultYear = year
	_ = config.SaveConfig(s.Config)

	s.SchoolInfo.AcademicYear = year
	_ = s.Store.SaveSchoolInfo(s.SchoolInfo)

	if s.OnViewRefresh != nil {
		s.OnViewRefresh()
	}
}

// SetSchoolType changes between LİSE and ORTAOKUL
func (s *AppState) SetSchoolType(sType string) {
	if sType == "" || sType == s.SchoolInfo.SchoolType {
		return
	}
	s.SchoolInfo.SchoolType = sType
	_ = s.Store.SaveSchoolInfo(s.SchoolInfo)

	if s.OnViewRefresh != nil {
		s.OnViewRefresh()
	}
}

// ApplyTheme switches color mode ("light", "cream", "dark") and font size
func (s *AppState) ApplyTheme(mode string, fontSize float32) {
	s.Config.ThemeMode = mode
	s.Config.FontSize = fontSize
	_ = config.SaveConfig(s.Config)

	if s.App != nil {
		s.App.Settings().SetTheme(appTheme.NewTurkishTheme(mode, fontSize))
	}
}

// SetUIMode switches between "modern" and "legacy" while preserving the user's chosen theme
func (s *AppState) SetUIMode(mode string) {
	if s.Config.UIMode == mode {
		return
	}
	s.Config.UIMode = mode
	_ = config.SaveConfig(s.Config)

	if s.OnViewRefresh != nil {
		s.OnViewRefresh()
	}
}

// NotifyChange refreshes the current active view
func (s *AppState) NotifyChange() {
	if s.OnViewRefresh != nil {
		s.OnViewRefresh()
	}
}
