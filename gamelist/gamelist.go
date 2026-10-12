package gamelist

import (
	"path/filepath"
	"strings"

	"github.com/beevik/etree"
)

const (
	NameElement        = "name"
	DescElement        = "desc"
	ImageElement       = "image"
	PlayersElement     = "players"
	GenreElement       = "genre"
	PathElement        = "path"
	GameListElement    = "gameList"
	GameElement        = "game"
	ReleaseDateElement = "releasedate"
	DeveloperElement   = "developer"
	PublisherElement   = "publisher"
	RatingElement      = "rating"
	MD5Element         = "md5"
	VideoElement       = "video"
	MarqueeElement     = "marquee"
	BezelElement       = "bezel"
	ManualElement      = "manual"
	FanartElement      = "fanart"
	BoxbackElement     = "boxback"
	ThumbnailElement   = "thumbnail"
	LangElement        = "lang"
	RegionElement      = "region"
	CheevosHashElement = "cheevosHash"
	CheevosIDElement   = "cheevosId"
	ScraperIDElement   = "scraperId"
)

type FileName string

const (
	GameListFileName      FileName = "gamelist.xml"
	MiyooGameListFileName FileName = "miyoogamelist.xml"
)

type GameList struct {
	document *etree.Document
	// byFileName indexes the entries by rom file name, so refreshing a whole
	// library is not one scan of the file per game. Built on first lookup.
	byFileName map[string]*etree.Element
}

func New() *GameList {
	return &GameList{
		document: emptyGameList(),
	}
}

func emptyGameList() *etree.Document {
	document := etree.NewDocument()
	document.CreateProcInst("xml", `version="1.0" encoding="UTF-8"`)
	document.CreateElement(GameListElement)
	return document
}

func (gl *GameList) Parse(b []byte) error {
	document := etree.NewDocument()
	if err := document.ReadFromBytes(b); err != nil {
		return err
	}

	gl.document = document
	gl.byFileName = nil
	return nil
}

// root returns the <gameList> element, or nil. A file written by another tool
// may have any root, so callers must tolerate its absence.
func (gl *GameList) root() *etree.Element {
	if gl.document == nil {
		return nil
	}
	return gl.document.SelectElement(GameListElement)
}

type gameMatcher func(*etree.Element) bool

func (gl *GameList) findGame(match gameMatcher) *etree.Element {
	root := gl.root()
	if root == nil {
		return nil
	}
	for _, game := range root.SelectElements(GameElement) {
		if match(game) {
			return game
		}
	}
	return nil
}

func byName(name string) gameMatcher {
	return func(game *etree.Element) bool {
		nameElement := game.SelectElement(NameElement)
		return nameElement != nil && nameElement.Text() == name
	}
}

// fileNameOf is the rom's file name from <path>, the stable identity of a rom
// entry. <name> is a display string and changes with regions, locale and
// settings.
//
// The base name is kept so grout adopts entries written by EmulationStation,
// which stores "./Game.gba" where grout stores an absolute path.
func fileNameOf(game *etree.Element) string {
	pathElement := game.SelectElement(PathElement)
	if pathElement == nil {
		return ""
	}
	return filepath.Base(strings.TrimSpace(pathElement.Text()))
}

func validFileName(fileName string) bool {
	return fileName != "" && fileName != "."
}

// fileNameIndex returns the entries keyed by rom file name, building it from
// the document on first use.
func (gl *GameList) fileNameIndex() map[string]*etree.Element {
	if gl.byFileName != nil {
		return gl.byFileName
	}
	gl.byFileName = make(map[string]*etree.Element)
	if root := gl.root(); root != nil {
		for _, game := range root.SelectElements(GameElement) {
			gl.index(game)
		}
	}
	return gl.byFileName
}

// index records game under its rom file name. Should several entries share
// it, the first one in the file is kept, as a scan of the file would find.
func (gl *GameList) index(game *etree.Element) {
	if gl.byFileName == nil {
		// Not built yet: it will be from the document, game included.
		return
	}
	fileName := fileNameOf(game)
	if !validFileName(fileName) {
		return
	}
	if _, taken := gl.byFileName[fileName]; !taken {
		gl.byFileName[fileName] = game
	}
}

// gameByAnyFileName returns the entry whose rom file name is the first of
// fileNames found.
//
// A game's file name on the server is not always the one on the device.
// A multi-disc game is launched through its playlist, and a rom with several
// versions is stored under the one that was picked.
func (gl *GameList) gameByAnyFileName(fileNames ...string) *etree.Element {
	for _, fileName := range fileNames {
		if !validFileName(fileName) {
			continue
		}
		if game := gl.fileNameIndex()[fileName]; game != nil {
			return game
		}
	}
	return nil
}

func (gl *GameList) Contains(element, value string) bool {
	return gl.findGame(func(game *etree.Element) bool {
		e := game.FindElement(element)
		return e != nil && e.Text() == value
	}) != nil
}

func (gl *GameList) GetGameElementByName(name string) *etree.Element {
	return gl.findGame(byName(name))
}

func (gl *GameList) GetGameElementByFileName(fileName string) *etree.Element {
	return gl.gameByAnyFileName(fileName)
}

func (gl *GameList) GameContainsElements(name string, elements []string) bool {
	e := gl.GetGameElementByName(name)
	if e == nil {
		return false
	}
	for _, element := range elements {
		if e.FindElement(element) == nil {
			return false
		}
	}
	return true
}

func (gl *GameList) Save(filepath string) error {
	gl.document.Indent(4)
	if err := gl.document.WriteToFile(filepath); err != nil {
		return err
	}
	return nil
}

func (gl *GameList) AddGameEntry(info map[string]string) *etree.Element {
	root := gl.root()
	if root == nil {
		return nil
	}
	newGame := root.CreateElement(GameElement)

	for key, value := range info {
		newGame.CreateElement(key).SetText(value)
	}
	gl.index(newGame)
	return newGame
}

// upsert updates game with info, creating the entry when game is nil. Returns
// the element so callers can set attributes, which requires it to exist.
func (gl *GameList) upsert(game *etree.Element, info map[string]string) *etree.Element {
	if game == nil {
		return gl.AddGameEntry(info)
	}

	for key, value := range info {
		gl.setChild(game, key, value)
	}
	return game
}

// setChild sets the text of game's child element, creating it if absent.
// A new <path> moves the entry to its new file name in the index.
func (gl *GameList) setChild(game *etree.Element, element, value string) {
	if element == PathElement && gl.byFileName != nil {
		if old := fileNameOf(game); gl.byFileName[old] == game {
			delete(gl.byFileName, old)
		}
		defer gl.index(game)
	}

	if child := game.SelectElement(element); child != nil {
		child.SetText(value)
	} else {
		game.CreateElement(element).SetText(value)
	}
}

// AddOrUpdateEntry upserts an entry keyed by <name>, for non-rom entries such
// as the Grout launcher shortcut. Rom entries must use AddOrUpdateRomEntry.
func (gl *GameList) AddOrUpdateEntry(name string, info map[string]string) *etree.Element {
	return gl.upsert(gl.findGame(byName(name)), info)
}

func (gl *GameList) AddOrUpdateRomEntry(fileName string, info map[string]string) *etree.Element {
	if fileName == "" {
		return gl.AddGameEntry(info)
	}
	return gl.upsert(gl.gameByAnyFileName(fileName), info)
}
