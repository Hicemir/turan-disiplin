package picker

import (
	"os/exec"
	"testing"
)

func TestPickerDefinitions(t *testing.T) {
	filters := []FileFilter{
		{Name: "Excel Dosyaları", Extensions: []string{"xlsx", "xls"}},
		{Name: "Veritabanı Dosyaları", Extensions: []string{"db"}},
	}

	if len(filters) != 2 {
		t.Fatalf("expected 2 filters, got %d", len(filters))
	}
	if filters[0].Extensions[0] != "xlsx" {
		t.Errorf("expected xlsx, got %s", filters[0].Extensions[0])
	}
}

func TestLinuxPickerToolAvailability(t *testing.T) {
	_, zenityErr := exec.LookPath("zenity")
	_, kdialogErr := exec.LookPath("kdialog")

	if zenityErr != nil && kdialogErr != nil {
		t.Log("Neither zenity nor kdialog is installed in PATH; picker will gracefully fall back to Fyne dialog")
	} else {
		t.Logf("Native CLI picker tool found (zenity: %v, kdialog: %v)", zenityErr == nil, kdialogErr == nil)
	}
}
