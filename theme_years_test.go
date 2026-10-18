package main

import (
	"image/color"
	"testing"

	"fyne.io/fyne/v2/theme"
	"turandisiplin/config"
	"turandisiplin/store"
	appTheme "turandisiplin/theme"
)

func TestAcademicYearSortingAndDeduplication(t *testing.T) {
	st, err := store.GetStore()
	if err != nil {
		t.Fatalf("failed to get store: %v", err)
	}

	years := st.GetAcademicYears()
	if len(years) < 4 {
		t.Fatalf("expected at least 4 academic years, got %d", len(years))
	}

	// Verify descending chronological order
	for i := 0; i < len(years)-1; i++ {
		y1 := store.ExtractStartYear(years[i])
		y2 := store.ExtractStartYear(years[i+1])
		if y1 < y2 {
			t.Errorf("years not sorted descending: %s before %s", years[i], years[i+1])
		}
		if years[i] == years[i+1] {
			t.Errorf("duplicate year detected: %s", years[i])
		}
	}

	// Verify normalization
	if store.NormalizeAcademicYear("2027-2028") != "2027 - 2028" {
		t.Errorf("normalization failed for 2027-2028: got %s", store.NormalizeAcademicYear("2027-2028"))
	}
	if store.NormalizeAcademicYear(" 2029 / 2030 ") != "2029 - 2030" {
		t.Errorf("normalization failed for 2029 / 2030: got %s", store.NormalizeAcademicYear(" 2029 / 2030 "))
	}

	// Test GetNextNearestAcademicYear
	nextYear := store.GetNextNearestAcademicYear([]string{"2028 - 2029", "2027 - 2028"})
	if nextYear != "2029 - 2030" {
		t.Errorf("expected 2029 - 2030, got %s", nextYear)
	}

	// Test DeleteAcademicYear
	if err := st.DeleteAcademicYear("9998 - 9999"); err != nil {
		t.Errorf("failed to delete non-existent year: %v", err)
	}
}

func TestThemesPaletteAndCustomTheme(t *testing.T) {
	customBg := color.NRGBA{R: 20, G: 20, B: 30, A: 255}
	customAccent := color.NRGBA{R: 255, G: 165, B: 0, A: 255}

	allThemes := appTheme.GetAllThemes(customBg, customAccent)
	if len(allThemes) != 3 { // dark_cyan + light + custom
		t.Fatalf("expected 3 themes including custom, got %d", len(allThemes))
	}

	expectedThemes := []string{"dark_cyan", "light"}
	for _, id := range expectedThemes {
		th := appTheme.NewTurkishTheme(id, 15.0)
		fg := th.Color(theme.ColorNameForeground, theme.VariantDark)
		bg := th.Color(theme.ColorNameBackground, theme.VariantDark)
		if fg == nil || bg == nil {
			t.Errorf("theme %s returned nil color", id)
		}
	}

	// Test Custom Theme
	thCustom := appTheme.NewTurkishTheme("custom", 15.0, customBg, customAccent)
	customBgRes := thCustom.Color(theme.ColorNameBackground, theme.VariantDark)
	customAccentRes := thCustom.Color(theme.ColorNamePrimary, theme.VariantDark)
	if customBgRes == nil || customAccentRes == nil {
		t.Errorf("custom theme returned nil colors")
	}

	// Test ColorToHex and HexToColor
	hex := config.ColorToHex(color.NRGBA{R: 255, G: 0, B: 128, A: 255})
	if hex != "#ff0080" {
		t.Errorf("expected #ff0080, got %s", hex)
	}
	cBack := config.HexToColor(hex, color.NRGBA{})
	if cBack.R != 255 || cBack.G != 0 || cBack.B != 128 {
		t.Errorf("hex roundtrip failed: %+v", cBack)
	}
}
