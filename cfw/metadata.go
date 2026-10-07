package cfw

import (
	"log/slog"
	"os"

	"grout/cfw/muos"
	"grout/gamelist"
)

// esRestartFlag is watched by the EmulationStation frontend, which reloads its
// gamelists when it appears. The path is relative to the directory grout was
// launched from, which is where the frontend expects it.
const esRestartFlag = "./es_restart_request"

// scheduleESRestart asks the frontend to reload, so metadata just written shows
// up without the user restarting the device.
func scheduleESRestart() {
	file, err := os.Create(esRestartFlag)
	if err != nil {
		slog.Default().Debug("Unable to schedule ES restart", "error", err)
		return
	}
	file.Close()
}

// AddGroutToGamelist adds grout's launcher entry to the frontend's game list.
// Only the EmulationStation family has somewhere to put it.
func AddGroutToGamelist(c CFW) {
	f := Lookup(c)
	path, launcher := f.GroutGamelist(), f.GroutLauncherPath()
	if path == "" || launcher == "" {
		return
	}

	gamelist.AddGroutEntry(path, launcher)
	scheduleESRestart()
}

// FillGamesMetadata writes game metadata in the form the firmware reads.
func FillGamesMetadata(entries []gamelist.RomGameEntry) {
	logger := slog.Default()

	switch ActiveFirmware().Gamelist() {
	case GamelistEmulationStation:
		if err := gamelist.AddRomGamesToGamelist(entries, gamelist.GameListFileName); err != nil {
			logger.Warn("Failed to add games to ES gamelist.xml", "error", err)
		}
		scheduleESRestart()

	case GamelistMiyoo:
		if err := gamelist.AddRomGamesToGamelist(entries, gamelist.MiyooGameListFileName); err != nil {
			logger.Warn("Failed to add games to miyoogamelist.xml", "error", err)
		}

	case GamelistMuOSText:
		for _, entry := range entries {
			muos.AddGameDescription(entry)
		}

	case GamelistNone:
	}
}

// HasGamesMetadata reports whether the firmware reads game metadata grout can
// write.
func HasGamesMetadata() bool {
	return ActiveFirmware().Gamelist() != GamelistNone
}

// RefreshGamesMetadata rewrites the metadata of games already on the device.
//
// Unlike FillGamesMetadata it keeps what the frontend recorded about each game
// in a gamelist, such as play count or favourites, and only replaces what
// grout writes. muOS keeps grout's text in files of its own, which are simply
// rewritten.
func RefreshGamesMetadata(entries []gamelist.RomGameEntry) error {
	switch ActiveFirmware().Gamelist() {
	case GamelistEmulationStation:
		err := gamelist.RefreshRomGamesInGamelist(entries, gamelist.GameListFileName)
		scheduleESRestart()
		return err

	case GamelistMiyoo:
		return gamelist.RefreshRomGamesInGamelist(entries, gamelist.MiyooGameListFileName)

	case GamelistMuOSText:
		for _, entry := range entries {
			muos.AddGameDescription(entry)
		}
	}
	return nil
}
