package main

import (
	"testing"

	"fyne.io/fyne/v2/theme"
	"okulevrak/store"
	appTheme "okulevrak/theme"
)

func TestAcademicYearSortingAndDeduplication(t *testing.T) {
	st, err := store.GetStore()
	if err != nil {
		t.Fatalf("failed to get store: %v", err)
	}

	years := st.GetAcademicYears()
	if len(years) < 6 {
		t.Fatalf("expected at least 6 academic years, got %d", len(years))
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
}

func TestThemesPaletteAndDarkModes(t *testing.T) {
	allThemes := appTheme.GetAllThemes()
	if len(allThemes) != 11 {
		t.Fatalf("expected 11 themes, got %d", len(allThemes))
	}

	expectedDark := []string{"dark", "dark_emerald", "dark_purple", "dark_amber", "dark_ruby", "dark_cyan"}
	for _, id := range expectedDark {
		th := appTheme.NewTurkishTheme(id, 15.0)
		fg := th.Color(theme.ColorNameForeground, theme.VariantDark)
		bg := th.Color(theme.ColorNameBackground, theme.VariantDark)
		if fg == nil || bg == nil {
			t.Errorf("theme %s returned nil color", id)
		}
	}

	expectedLight := []string{"light", "cream", "green", "bordeaux", "navy"}
	for _, id := range expectedLight {
		th := appTheme.NewTurkishTheme(id, 15.0)
		fg := th.Color(theme.ColorNameForeground, theme.VariantLight)
		bg := th.Color(theme.ColorNameBackground, theme.VariantLight)
		if fg == nil || bg == nil {
			t.Errorf("theme %s returned nil color", id)
		}
	}
}
