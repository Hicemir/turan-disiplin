package model

import (
	"strconv"
	"strings"
	"time"
)

// NormalizeBirthDate converts Excel serial date numbers (e.g. "41050" or
// "41050.0") into dd.MM.yyyy. Any non-serial text is returned unchanged.
func NormalizeBirthDate(s string) string {
	s = strings.TrimSpace(s)
	if s == "" {
		return s
	}
	f, err := strconv.ParseFloat(strings.ReplaceAll(s, ",", "."), 64)
	if err != nil || f < 1 || f > 80000 {
		return s
	}
	// Excel 1900 date system epoch (accounts for the fictitious 1900-02-29).
	base := time.Date(1899, 12, 30, 0, 0, 0, 0, time.UTC)
	return base.AddDate(0, 0, int(f)).Format("02.01.2006")
}
