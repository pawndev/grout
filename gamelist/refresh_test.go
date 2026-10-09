package gamelist

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/beevik/etree"
)

func parsed(t *testing.T, xml string) *GameList {
	t.Helper()
	gl := New()
	if err := gl.Parse([]byte(xml)); err != nil {
		t.Fatalf("parse: %v", err)
	}
	return gl
}

func onlyGame(t *testing.T, gl *GameList) *etree.Element {
	t.Helper()
	games := gl.root().SelectElements(GameElement)
	if len(games) != 1 {
		t.Fatalf("expected 1 <game> element, got %d: %v", len(games), gameElements(t, gl))
	}
	return games[0]
}

func childText(game *etree.Element, element string) (string, bool) {
	child := game.FindElement(element)
	if child == nil {
		return "", false
	}
	return child.Text(), true
}

func TestRefreshRomGame_KeepsFrontendFields(t *testing.T) {
	gl := parsed(t, `<?xml version="1.0" encoding="UTF-8"?>
<gameList>
    <game>
        <path>./Sonic the Hedgehog.gba</path>
        <name>Sonic</name>
        <desc>Old description</desc>
        <playcount>12</playcount>
        <lastplayed>20260101T120000</lastplayed>
        <gametime>3600</gametime>
        <favorite>true</favorite>
        <hidden>false</hidden>
    </game>
</gameList>`)

	r := rom("Sonic the Hedgehog", "Sonic the Hedgehog.gba", "USA")
	r.Summary = "New description"
	r.Genres = []string{"Platform"}
	gl.RefreshRomGame(entry(r, "/roms/gba"))

	game := onlyGame(t, gl)
	for element, want := range map[string]string{
		NameElement:  "Sonic the Hedgehog (USA)",
		DescElement:  "New description",
		GenreElement: "Platform",
		"playcount":  "12",
		"lastplayed": "20260101T120000",
		"gametime":   "3600",
		"favorite":   "true",
		"hidden":     "false",
	} {
		if got, ok := childText(game, element); !ok || got != want {
			t.Errorf("<%s> = %q (present %v), want %q", element, got, ok, want)
		}
	}
}

// The existing <path> is how the firmware finds the rom, possibly relative.
// A refresh must not move it.
func TestRefreshRomGame_KeepsExistingPath(t *testing.T) {
	gl := parsed(t, `<gameList><game><path>./Sonic.gba</path><name>Sonic</name></game></gameList>`)

	gl.RefreshRomGame(entry(rom("Sonic", "Sonic.gba"), "/roms/gba"))

	if got, _ := childText(onlyGame(t, gl), PathElement); got != "./Sonic.gba" {
		t.Errorf("<path> = %q, want it left as %q", got, "./Sonic.gba")
	}
}

// A field the server has nothing for must not blank what a scraper filled in.
func TestRefreshRomGame_EmptyValuesDoNotOverwrite(t *testing.T) {
	gl := parsed(t, `<gameList><game><path>./Sonic.gba</path><desc>Scraped description</desc><image>./media/Sonic.png</image></game></gameList>`)

	gl.RefreshRomGame(entry(rom("Sonic", "Sonic.gba"), "/roms/gba"))

	game := onlyGame(t, gl)
	if got, _ := childText(game, DescElement); got != "Scraped description" {
		t.Errorf("<desc> = %q, want the scraped one kept", got)
	}
	if got, _ := childText(game, ImageElement); got != "./media/Sonic.png" {
		t.Errorf("<image> = %q, want it kept when no art is on the device", got)
	}
}

// A multi-disc game is written under its playlist, whose name is not the
// game's file name on the server. Matching on the server name alone would add
// a second entry.
func TestRefreshRomGame_MatchesOnLocalFileName(t *testing.T) {
	gl := parsed(t, `<gameList><game><path>/roms/psx/Final Fantasy VII.m3u</path><name>FF7</name><playcount>3</playcount></game></gameList>`)

	r := rom("Final Fantasy VII", "Final Fantasy VII")
	r.Summary = "Cloud and friends"
	e := entry(r, "/roms/psx")
	e.Game.Path = "/roms/psx/Final Fantasy VII.m3u"
	gl.RefreshRomGame(e)

	game := onlyGame(t, gl)
	if got, _ := childText(game, DescElement); got != "Cloud and friends" {
		t.Errorf("<desc> = %q, want it updated", got)
	}
	if got, _ := childText(game, "playcount"); got != "3" {
		t.Errorf("<playcount> = %q, want it kept", got)
	}
}

// A game on the device with no entry, as after a download cut short before the
// metadata was written, gets a complete one.
func TestRefreshRomGame_CreatesMissingEntry(t *testing.T) {
	gl := New()
	r := rom("Sonic", "Sonic.gba")
	r.ScreenScraperID = 42
	gl.RefreshRomGame(entry(r, "/roms/gba"))

	game := onlyGame(t, gl)
	if got, _ := childText(game, PathElement); got != "/roms/gba/Sonic.gba" {
		t.Errorf("<path> = %q, want the rom's location", got)
	}
	if attr := game.SelectAttr("id"); attr == nil || attr.Value != "42" {
		t.Errorf("id attribute = %v, want 42", attr)
	}
}

// Other games in the file are not the refresh's business.
func TestRefreshRomGamesInGamelist_LeavesOtherEntriesAlone(t *testing.T) {
	romDir := t.TempDir()
	path := filepath.Join(romDir, string(GameListFileName))
	if err := os.WriteFile(path, []byte(`<gameList>
<game><path>./Other.gba</path><name>Other</name><playcount>5</playcount></game>
<game><path>./Sonic.gba</path><name>Sonic</name></game>
</gameList>`), 0644); err != nil {
		t.Fatal(err)
	}

	r := rom("Sonic the Hedgehog", "Sonic.gba")
	if err := RefreshRomGamesInGamelist([]RomGameEntry{entry(r, romDir)}, GameListFileName); err != nil {
		t.Fatalf("refresh: %v", err)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	gl := parsed(t, string(data))
	if got := gameElements(t, gl); len(got) != 2 || got[0] != "Other" || got[1] != "Sonic the Hedgehog" {
		t.Fatalf("entries = %v, want [Other, Sonic the Hedgehog]", got)
	}
	if !gl.Contains("playcount", "5") {
		t.Error("expected the other game's play count to survive")
	}
}

// An unknown player count must not pass for a single-player game and replace
// the count already recorded.
func TestRefreshRomGame_UnknownPlayersKeepsExisting(t *testing.T) {
	gl := parsed(t, `<gameList><game><path>./Sonic.gba</path><players>1-4</players></game></gameList>`)

	gl.RefreshRomGame(entry(rom("Sonic", "Sonic.gba"), "/roms/gba"))

	if got, _ := childText(onlyGame(t, gl), PlayersElement); got != "1-4" {
		t.Errorf("<players> = %q, want %q kept", got, "1-4")
	}
}

// A known count still replaces what is there.
func TestRefreshRomGame_KnownPlayersReplaces(t *testing.T) {
	gl := parsed(t, `<gameList><game><path>./Sonic.gba</path><players>1-4</players></game></gameList>`)

	r := rom("Sonic", "Sonic.gba")
	r.MaxPlayers = 1
	gl.RefreshRomGame(entry(r, "/roms/gba"))

	if got, _ := childText(onlyGame(t, gl), PlayersElement); got != "1" {
		t.Errorf("<players> = %q, want %q", got, "1")
	}
}

// A fresh download keeps writing a single player when the count is unknown.
func TestAddRomGame_UnknownPlayersDefaultsToOne(t *testing.T) {
	gl := New()
	gl.AddRomGame(entry(rom("Sonic", "Sonic.gba"), "/roms/gba"))

	if got, _ := childText(onlyGame(t, gl), PlayersElement); got != "1" {
		t.Errorf("<players> = %q, want %q", got, "1")
	}
}
