package config

import (
	"encoding/json"
	"fmt"
	"image/color"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"turandisiplin/assets"
	"turandisiplin/model"
)

var (
	currentYonetmelik *model.YonetmelikData
	yonetmelikMutex   sync.RWMutex
)

// GetConfigDir returns the directory where user configurations and SQLite database are stored
func GetConfigDir() (string, error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		configDir = "."
	}
	appDir := filepath.Join(configDir, "TuranDisiplin")
	if err := os.MkdirAll(appDir, 0755); err != nil {
		return ".", err
	}
	return appDir, nil
}

// GetUserYonetmelikPath returns the path to user-edited yonetmelik.json
func GetUserYonetmelikPath() string {
	dir, err := GetConfigDir()
	if err != nil {
		return "yonetmelik.json"
	}
	return filepath.Join(dir, "yonetmelik.json")
}

// LoadYonetmelik loads from user config directory if present, otherwise from embedded assets
func LoadYonetmelik() (*model.YonetmelikData, error) {
	yonetmelikMutex.Lock()
	defer yonetmelikMutex.Unlock()

	userPath := GetUserYonetmelikPath()
	var dataBytes []byte

	if info, err := os.Stat(userPath); err == nil && info.Size() > 0 {
		readBytes, err := os.ReadFile(userPath)
		if err == nil && len(readBytes) > 0 {
			dataBytes = readBytes
		}
	}

	if len(dataBytes) == 0 {
		dataBytes = assets.YonetmelikJSON
	}

	var parsed model.YonetmelikData
	if err := json.Unmarshal(dataBytes, &parsed); err != nil {
		// Fallback to embedded if user file was corrupted
		if errFallback := json.Unmarshal(assets.YonetmelikJSON, &parsed); errFallback != nil {
			return nil, fmt.Errorf("yonetmelik JSON parse hatası: %w", errFallback)
		}
	}

	currentYonetmelik = &parsed
	return currentYonetmelik, nil
}

// GetYonetmelik returns the active in-memory yonetmelik configuration
func GetYonetmelik() *model.YonetmelikData {
	yonetmelikMutex.RLock()
	if currentYonetmelik != nil {
		defer yonetmelikMutex.RUnlock()
		return currentYonetmelik
	}
	yonetmelikMutex.RUnlock()

	data, err := LoadYonetmelik()
	if err != nil {
		panic(err)
	}
	return data
}

// SaveYonetmelik saves the updated yonetmelik data to os.UserConfigDir()/TuranDisiplin/yonetmelik.json
func SaveYonetmelik(data *model.YonetmelikData) error {
	yonetmelikMutex.Lock()
	defer yonetmelikMutex.Unlock()

	bytes, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return fmt.Errorf("JSON serileştirme hatası: %w", err)
	}

	userPath := GetUserYonetmelikPath()
	if err := os.WriteFile(userPath, bytes, 0644); err != nil {
		return fmt.Errorf("dosyaya yazma hatası (%s): %w", userPath, err)
	}

	currentYonetmelik = data
	return nil
}

// ResetToDefault removes any custom user file and reloads the embedded default yonetmelik
func ResetToDefault() (*model.YonetmelikData, error) {
	yonetmelikMutex.Lock()
	defer yonetmelikMutex.Unlock()

	userPath := GetUserYonetmelikPath()
	_ = os.Remove(userPath)

	var parsed model.YonetmelikData
	if err := json.Unmarshal(assets.YonetmelikJSON, &parsed); err != nil {
		return nil, err
	}

	currentYonetmelik = &parsed
	return currentYonetmelik, nil
}

// ResetYonetmelik is an alias for ResetToDefault
func ResetYonetmelik() (*model.YonetmelikData, error) {
	return ResetToDefault()
}

// AppConfig stores user preferences like last selected academic year, theme settings and custom colors
type AppConfig struct {
	DefaultYear       string   `json:"default_year"`
	RememberYear      bool     `json:"remember_year"`        // If true, skip startup year dialog
	UIMode            string   `json:"ui_mode"`              // "modern" or "legacy"
	ThemeMode         string   `json:"theme_mode"`           // theme ID or "custom"
	ThemeUserChosen   bool     `json:"theme_user_chosen,omitempty"` // true once the user picks a theme; otherwise the OS theme is followed
	FontSize          float32  `json:"font_size"`            // 15, 17, 19
	DeletedYears      []string `json:"deleted_years,omitempty"`
	CustomBgColor     string   `json:"custom_bg_color,omitempty"`
	CustomAccentColor string   `json:"custom_accent_color,omitempty"`
}

// ColorToHex converts any color.Color to #rrggbb hex string
func ColorToHex(c color.Color) string {
	if c == nil {
		return "#1e1e2e"
	}
	r, g, b, _ := c.RGBA()
	return fmt.Sprintf("#%02x%02x%02x", uint8(r>>8), uint8(g>>8), uint8(b>>8))
}

// HexToColor parses a #rrggbb or rrggbb hex string to color.NRGBA
func HexToColor(s string, fallback color.NRGBA) color.NRGBA {
	s = strings.TrimPrefix(s, "#")
	if len(s) == 6 {
		var r, g, b uint8
		if _, err := fmt.Sscanf(s, "%02x%02x%02x", &r, &g, &b); err == nil {
			return color.NRGBA{R: r, G: g, B: b, A: 255}
		}
	}
	return fallback
}

func GetConfigFilePath() string {
	dir, err := GetConfigDir()
	if err != nil {
		return "config.json"
	}
	return filepath.Join(dir, "config.json")
}

// LoadConfig loads the application settings
func LoadConfig() (*AppConfig, error) {
	cfgPath := GetConfigFilePath()
	cfg := &AppConfig{
		DefaultYear: "2026 - 2027",
		UIMode:      "modern",
		ThemeMode:   "", // empty indicates unconfigured -> auto-detect system theme
		FontSize:    15.0,
	}
	data, err := os.ReadFile(cfgPath)
	if err == nil {
		_ = json.Unmarshal(data, cfg)
	}
	if cfg.UIMode == "" {
		cfg.UIMode = "modern"
	}
	if cfg.FontSize < 13 {
		cfg.FontSize = 15.0
	}
	return cfg, nil
}

// SaveConfig saves the application settings
func SaveConfig(cfg *AppConfig) error {
	cfgPath := GetConfigFilePath()
	b, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(cfgPath, b, 0644)
}

// GetClausesForSchoolType returns all clauses filtered by school type ("LİSE" or "ORTAOKUL")
func GetClausesForSchoolType(schoolType string) []model.Clause {
	y := GetYonetmelik()
	var res []model.Clause
	for _, c := range y.Clauses {
		if c.SchoolType == schoolType {
			res = append(res, c)
		}
	}
	return res
}
