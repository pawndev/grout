package download

import (
	"grout/catalog"
	"grout/cfw"
	"grout/files"
	"grout/gamelist"
	"grout/library"
	"grout/romm"
	"grout/settings"
	"grout/textmatch"
)

// MetadataEntries rebuilds the metadata entries of the games already on the
// device, from what the server knows about them now.
// Artwork is pointed at only where its file is present, so a refresh never
// leaves an entry naming art that was not downloaded.
func MetadataEntries(config settings.Config, platform romm.Platform, games []romm.Rom) []gamelist.RomGameEntry {
	activeCFW := cfw.GetCFW()
	isESBased := activeCFW.IsBasedOnEmulationStation()

	entries := make([]gamelist.RomGameEntry, 0, len(games))
	for _, game := range games {
		path := catalog.LocalRomPath(config, game)
		if path == "" {
			continue
		}
		gamePlatform := platformFor(platform, game)

		regions := game.Regions
		if config.GamelistOmitsRegion {
			regions = nil
		}

		entries = append(entries, gamelist.RomGameEntry{
			Game:         game.ToGame(textmatch.PrepareRomName(game.Name, regions), path, presentArt(config, game, gamePlatform, activeCFW, isESBased)),
			Platform:     gamePlatform.ToPlatform(),
			RomDirectory: cfw.PlatformRomDirectory(config, gamePlatform.FSSlug),
		})
	}
	return entries
}

// PlatformGames is a platform and some of its games.
type PlatformGames struct {
	Platform romm.Platform
	Games    []romm.Rom
}

// RefreshMetadata rewrites the metadata of the games already on the device and
// returns how many games it covered. What the firmware recorded about each
// game is kept
func RefreshMetadata(config settings.Config, batches []PlatformGames) (int, error) {
	var entries []gamelist.RomGameEntry
	for _, batch := range batches {
		entries = append(entries, MetadataEntries(config, batch.Platform, batch.Games)...)
	}
	if len(entries) == 0 {
		return 0, nil
	}
	if err := cfw.RefreshGamesMetadata(entries); err != nil {
		return 0, err
	}
	return len(entries), nil
}

// presentArt records the art of a game whose file is on the device, wherever
// the firmware would look for it
func presentArt(config settings.Config, game romm.Rom, platform romm.Platform, activeCFW cfw.CFW, isESBased bool) library.ArtPaths {
	var paths library.ArtPaths
	for _, spec := range artSpecs {
		location := spec.location(config, game, platform, activeCFW)
		if location != "" && files.FileExists(location) {
			spec.record(&paths, location, isESBased)
		}
	}
	return paths
}
