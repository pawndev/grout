package gamelist

import (
	"fmt"
	"grout/files"
	"log/slog"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"grout/library"

	"github.com/beevik/etree"
)

type GameListEntry struct {
	GL   *GameList
	Path string
}

// RomGameEntry is one game to write into a gamelist. It carries library.Game
// rather than the wire type so this package performs no text transformation:
// the display name arrives rendered, and identity is the file name.
type RomGameEntry struct {
	Game         library.Game
	Platform     library.Platform
	RomDirectory string
}

func (e RomGameEntry) FileName() string {
	if e.Game.FileName != "" {
		return e.Game.FileName
	}
	if e.Game.Path == "" {
		return ""
	}
	return filepath.Base(e.Game.Path)
}

func (gl *GameList) AddRomGame(entry RomGameEntry) {
	element := gl.AddOrUpdateRomEntry(entry.FileName(), romMetadata(entry.Game))
	setScraperID(element, entry.Game)
}

// RefreshRomGame rewrites the metadata grout manages on a rom's entry, creating
// the entry when there is none.
//
// Only the elements grout writes are touched: what the firmware records on its
// own, such as play count, play time or favourites, is left as it is.
// Values the server has nothing for are skipped rather than blanking what
// a scraper may have filled in.
func (gl *GameList) RefreshRomGame(entry RomGameEntry) {
	metadata := romMetadata(entry.Game)
	for element, value := range metadata {
		if value == "" {
			delete(metadata, element)
		}
	}

	game := gl.findGame(byAnyFileName(entry.FileName(), filepath.Base(entry.Game.Path)))
	if game == nil {
		setScraperID(gl.AddGameEntry(metadata), entry.Game)
		return
	}

	delete(metadata, PathElement)
	for element, value := range metadata {
		setChild(game, element, value)
	}
	setScraperID(game, entry.Game)
}

// romMetadata is every element grout writes for a game, keyed by element name.
func romMetadata(game library.Game) map[string]string {
	gameMetadata := map[string]string{
		NameElement: game.DisplayName,
		DescElement: game.Summary,
		MD5Element:  game.MD5,
	}

	if game.HasRating() {
		gameMetadata[RatingElement] = fmt.Sprintf("%.1f", game.Rating)
	}

	if game.HasReleaseDate() {
		gameMetadata[ReleaseDateElement] = game.ReleaseDate.Format("20060102T150405")
	}

	for element, path := range map[string]string{
		ImageElement:     game.Art.Cover,
		ThumbnailElement: game.Art.Thumbnail,
		MarqueeElement:   game.Art.Marquee,
		VideoElement:     game.Art.Video,
		BezelElement:     game.Art.Bezel,
		ManualElement:    game.Art.Manual,
		BoxbackElement:   game.Art.BoxBack,
		FanartElement:    game.Art.Fanart,
		PathElement:      game.Path,
	} {
		if path != "" {
			gameMetadata[element] = path
		}
	}

	if game.MaxPlayers > 1 {
		gameMetadata[PlayersElement] = fmt.Sprintf("1-%d", game.MaxPlayers)
	} else {
		gameMetadata[PlayersElement] = "1"
	}

	for element, values := range map[string][]string{
		RegionElement:    game.Regions,
		LangElement:      game.Languages,
		GenreElement:     game.Genres,
		DeveloperElement: game.Developers,
	} {
		if len(values) > 0 {
			gameMetadata[element] = strings.Join(values, ", ")
		}
	}

	if game.ScreenScraperID > 0 {
		gameMetadata[ScraperIDElement] = strconv.Itoa(game.ScreenScraperID)
	}

	if game.RetroAchievementsID > 0 {
		gameMetadata[CheevosIDElement] = strconv.Itoa(game.RetroAchievementsID)
	}

	if game.RetroAchievementsHash != "" {
		gameMetadata[CheevosHashElement] = game.RetroAchievementsHash
	}

	return gameMetadata
}

// setScraperID mirrors the ScreenScraper id onto the entry's id attribute,
// which needs the element to exist and so follows the upsert.
func setScraperID(element *etree.Element, game library.Game) {
	if element != nil && game.ScreenScraperID > 0 {
		element.CreateAttr("id", strconv.Itoa(game.ScreenScraperID))
	}
}

// AddRomGamesToGamelist adds or updates the entries of freshly downloaded games.
func AddRomGamesToGamelist(entries []RomGameEntry, gamelistFilename FileName) error {
	return applyToGamelists(entries, gamelistFilename, (*GameList).AddRomGame)
}

// RefreshRomGamesInGamelist brings the entries of games already on the device
// in line with the server, keeping what the frontend recorded about them. See
// RefreshRomGame.
func RefreshRomGamesInGamelist(entries []RomGameEntry, gamelistFilename FileName) error {
	return applyToGamelists(entries, gamelistFilename, (*GameList).RefreshRomGame)
}

// applyToGamelists loads the gamelist of each platform the entries belong to
// once, applies to apply to every entry, then saves each file.
func applyToGamelists(entries []RomGameEntry, gamelistFilename FileName, apply func(*GameList, RomGameEntry)) error {
	gamelists := make(map[string]GameListEntry)
	for _, game := range entries {
		glEntry, exists := gamelists[game.Platform.FSSlug]
		if !exists {
			gl := New()
			gamelistPath := fmt.Sprintf("%s/%s", game.RomDirectory, gamelistFilename)
			if files.FileExists(gamelistPath) {
				data, err := os.ReadFile(gamelistPath)
				if err != nil {
					// Saving over a file that could not be read would wipe
					// every entry in it.
					slog.Default().Error("Error reading gamelist file, skipping platform", "error", err, "path", gamelistPath)
					continue
				}
				if len(data) > 0 {
					if err := gl.Parse(data); err != nil {
						slog.Default().Error("gamelist not found or can't be parsed, skipping platform", "path", gamelistPath, "error", err)
						continue
					}
				}
			}
			glEntry = GameListEntry{Path: gamelistPath, GL: gl}
			gamelists[game.Platform.FSSlug] = glEntry
		}

		apply(glEntry.GL, game)
	}

	for _, glEntry := range gamelists {
		if err := glEntry.GL.Save(glEntry.Path); err != nil {
			slog.Default().Error("Unable to save gamelist file", "error", err, "path", glEntry.Path)
			return err
		}
		slog.Default().Debug("Successfully saved gamelist file", "path", glEntry.Path)
	}

	return nil
}
