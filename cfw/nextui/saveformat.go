package nextui

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"
)

// NextUI's Save format setting, as its settings file stores it.
const (
	saveFormatMinUI                 = "0" // Game.gba.sav, the default
	saveFormatRetroArch             = "1" // Game.srm, compressed
	saveFormatGeneric               = "2" // Game.sav
	saveFormatRetroArchUncompressed = "3" // Game.srm
)

// SaveNaming is how NextUI's emulator names a save: whether it keeps the rom's
// extension, and the save's own extension. It follows the Save format setting,
// and NextUI's default when there is none.
func SaveNaming() (keepRomExt bool, ext string) {
	switch saveFormat() {
	case saveFormatRetroArch, saveFormatRetroArchUncompressed:
		return false, "srm"
	case saveFormatGeneric:
		return false, "sav"
	default:
		return true, "sav"
	}
}

func saveFormat() string {
	file, err := os.Open(filepath.Join(GetBasePath(), ".userdata", "shared", "minuisettings.txt"))
	if err != nil {
		return saveFormatMinUI
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		if value, ok := strings.CutPrefix(strings.TrimSpace(scanner.Text()), "saveFormat="); ok {
			return value
		}
	}
	return saveFormatMinUI
}
