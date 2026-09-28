package nextui

import (
	"os"
	"path/filepath"
	"testing"
)

func withSettings(t *testing.T, content string) {
	t.Helper()
	base := t.TempDir()
	t.Setenv("BASE_PATH", base)
	if content == "" {
		return
	}
	dir := filepath.Join(base, ".userdata", "shared")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "minuisettings.txt"), []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}

// NextUI's Save format setting decides whether a save keeps the rom's
// extension and whether it ends in .sav or .srm.
func TestSaveNaming(t *testing.T) {
	for _, tc := range []struct {
		name     string
		settings string
		keep     bool
		ext      string
	}{
		{"no settings file is NextUI's default", "", true, "sav"},
		{"MinUI", "fontSize=2\nsaveFormat=0\n", true, "sav"},
		{"RetroArch", "saveFormat=1\n", false, "srm"},
		{"Generic", "saveFormat=2\n", false, "sav"},
		{"RetroArch uncompressed", "saveFormat=3\n", false, "srm"},
		{"setting missing from the file", "fontSize=2\n", true, "sav"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			withSettings(t, tc.settings)
			keep, ext := SaveNaming()
			if keep != tc.keep || ext != tc.ext {
				t.Errorf("keep=%v ext=%q, want keep=%v ext=%q", keep, ext, tc.keep, tc.ext)
			}
		})
	}
}
