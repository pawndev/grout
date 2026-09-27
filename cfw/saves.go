package cfw

import (
	"grout/settings"
	"grout/textmatch"
	"path/filepath"
	"slices"
	"strings"
)

// SaveBasename returns the on-disk basename (the part before the save-file extension) an
// emulator uses for a ROM's save files, given whether the device keeps the ROM extension.
//
// keepRomExt=false is the RetroArch convention: the save is named after the ROM basename
// WITHOUT its extension (e.g. ROM "Game (USA).gba" -> save "Game (USA).srm"). keepRomExt=true
// is the minarch convention (NextUI/MinUI default): the save keeps the FULL ROM filename,
// extension included (e.g. ROM "Game (USA).sfc" -> save "Game (USA).sfc.sav"). Reading or
// writing a save under the wrong convention silently breaks sync (issue #245). NextUI
// exposes both as a setting, so the style is detected per-device rather than assumed by CFW.
func SaveBasename(keepRomExt bool, romFileName string) string {
	if keepRomExt {
		return romFileName
	}
	return strings.TrimSuffix(romFileName, filepath.Ext(romFileName))
}

// DefaultKeepsRomExt reports the CFW's default save-naming style, used only as
// a fallback when the actual on-device convention can't be detected from
// existing saves (issue #245).
func DefaultKeepsRomExt(c CFW) bool {
	return Lookup(c).KeepsRomExtInSaves()
}

// EmulatorFolderMap returns the emulator save folders the given CFW uses,
// keyed by RomM filesystem slug. Firmwares that keep saves in the rom folder
// return their platform table.
func EmulatorFolderMap(c CFW) map[string][]string {
	return Lookup(c).SaveDirectories()
}

// EmulatorFoldersForFSSlug returns the emulator folders for a given filesystem slug.
func EmulatorFoldersForFSSlug(fsSlug string) []string {
	saveDirectoriesMap := EmulatorFolderMap(GetCFW())
	if saveDirectoriesMap == nil {
		return nil
	}
	return saveDirectoriesMap[fsSlug]
}

// GetSaveDirectory returns the full save directory path for a given filesystem slug.
// Falls back to the first emulator folder if no match is found.
func GetSaveDirectory(fsSlug string) string {
	baseSavePath := BaseSavePath()
	if baseSavePath == "" {
		return ""
	}

	emulatorDirs := EmulatorFoldersForFSSlug(fsSlug)
	if len(emulatorDirs) == 0 {
		return ""
	}

	return filepath.Join(baseSavePath, emulatorDirs[0])
}

// GetSaveDirectoryForRomPath resolves the save directory associated with the ROM's
// actual emulator folder. Tag-based CFWs can expose multiple launchers for one platform
// (for example NextUI's Game Boy Advance (GBA) and Game Boy Advance (MGBA)); choosing the
// platform's first save directory would write a valid save where the selected launcher
// never looks for it.
//
// If the ROM folder cannot be correlated with a configured save directory, this falls
// back to the platform default used by GetSaveDirectory.
func GetSaveDirectoryForRomPath(fsSlug, romPath string) string {
	baseSavePath := BaseSavePath()
	if baseSavePath == "" {
		return ""
	}

	emulatorDirs := EmulatorFoldersForFSSlug(fsSlug)
	if len(emulatorDirs) == 0 {
		return ""
	}

	romDir := filepath.Base(filepath.Dir(romPath))
	romTag := textmatch.ParseTag(romDir)
	for _, emulatorDir := range emulatorDirs {
		if strings.EqualFold(emulatorDir, romDir) || (romTag != "" && strings.EqualFold(emulatorDir, romTag)) {
			return filepath.Join(baseSavePath, emulatorDir)
		}
	}

	return filepath.Join(baseSavePath, emulatorDirs[0])
}

// SaveFolders lists where a platform's saves may be, relative to the save
// root, the one new saves go in first.
func SaveFolders(config settings.Config, fsSlug string) []string {
	return ActiveFirmware().SaveFolders(config, fsSlug)
}

// PreferredSaveFolder is where new saves must go for the active firmware, or
// "" when it has no opinion. See Firmware.PreferredSaveFolder.
func PreferredSaveFolder(config settings.Config, fsSlug string) string {
	return ActiveFirmware().PreferredSaveFolder(config, fsSlug)
}

func (f *Firmware) SaveFolders(config settings.Config, fsSlug string) []string {
	folders := f.SaveDirectories()[config.ResolveFSSlug(fsSlug)]

	extra, preferred := f.extraSaveFoldersFor(config, fsSlug)
	if preferred {
		return appendNew(append([]string(nil), extra...), folders...)
	}
	return appendNew(append([]string(nil), folders...), extra...)
}

// PreferredSaveFolder is where new saves must go when the firmware says so,
// overriding the save mapping. muOS decides per rom folder whether Pickles runs
// it, and a save outside Pickles' folders is one Pickles never loads. A mapping
// to another Pickles core is kept.
func (f *Firmware) PreferredSaveFolder(config settings.Config, fsSlug string) string {
	extra, preferred := f.extraSaveFoldersFor(config, fsSlug)
	if !preferred {
		return ""
	}
	if mapped := config.SaveDirectoryMappings[fsSlug]; slices.Contains(extra, mapped) {
		return mapped
	}
	return extra[0]
}

func (f *Firmware) extraSaveFoldersFor(config settings.Config, fsSlug string) ([]string, bool) {
	if f == nil || f.extraSaveFolders == nil {
		return nil, false
	}
	romFolder := fsSlug
	if mapping, ok := config.DirectoryMappings[fsSlug]; ok && mapping.RelativePath != "" {
		romFolder = mapping.RelativePath
	}
	return f.extraSaveFolders(romFolder)
}

func appendNew(list []string, more ...string) []string {
	for _, item := range more {
		if !slices.Contains(list, item) {
			list = append(list, item)
		}
	}
	return list
}

// EmulatorLabel is what to call a save folder on screen for firmware c.
func EmulatorLabel(c CFW, dir string) string {
	return Lookup(c).EmulatorLabel(dir)
}
