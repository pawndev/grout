package ui

import (
	"fmt"

	"grout/catalog"
	"grout/download"
	"grout/romm"
	"grout/settings"

	gaba "github.com/BrandonKowalski/gabagool/v2/pkg/gabagool"
)

type MetadataSyncInput struct {
	Config settings.Config
	Host   settings.Host
}

type MetadataSyncOutput struct{}

// MetadataSyncScreen rewrites the metadata of games already on the device from
// what the server knows now, for when a download was cut short or the server's
// metadata changed since. What the firmware recorded about each game, such as
// play count or favourites, is kept.
type MetadataSyncScreen struct{}

func NewMetadataSyncScreen() *MetadataSyncScreen {
	return &MetadataSyncScreen{}
}

func (s *MetadataSyncScreen) Execute(input MetadataSyncInput) MetadataSyncOutput {
	s.draw(input)
	return MetadataSyncOutput{}
}

func (s *MetadataSyncScreen) draw(input MetadataSyncInput) {
	platforms, err := catalog.MappedPlatforms(input.Host, input.Config.DirectoryMappings, input.Config.ApiTimeout.Duration())
	if err != nil {
		gaba.GetLogger().Error("Failed to fetch platforms", "error", err)
		s.tell(fmt.Sprintf("Failed to fetch platforms: %v", err))
		return
	}
	if len(platforms) == 0 {
		s.tell(localize("artwork_sync_no_platforms", "No platforms with directory mappings found."))
		return
	}

	found := s.scan(input, platforms)
	if len(found) == 0 {
		s.tell(localize("metadata_sync_no_games", "No downloaded games found."))
		return
	}

	chosen, ok := choosePlatforms(found, localize("button_update", "Update"))
	if !ok {
		return
	}

	s.refresh(input, chosen)
}

// scan finds the downloaded games of each platform, showing progress since it
// reads every platform's library and looks for each game on the card.
func (s *MetadataSyncScreen) scan(input MetadataSyncInput, platforms []romm.Platform) []platformRoms {
	logger := gaba.GetLogger()

	var found []platformRoms
	for i, platform := range platforms {
		// ProcessMessage runs the closure on the calling goroutine, so
		// appending from inside is safe
		gaba.ProcessMessage(
			fmt.Sprintf(localize("artwork_sync_scanning", "Scanning platform %d/%d: %s..."), i+1, len(platforms), platform.Name),
			gaba.ProcessMessageOptions{ShowThemeBackground: true},
			func() (any, error) {
				games, err := catalog.Games(catalog.GameSource{Platform: platform})
				if err != nil {
					logger.Error("Failed to read platform games", "platform", platform.Name, "error", err)
					return nil, nil
				}

				downloaded := make([]romm.Rom, 0, len(games))
				for _, game := range games {
					if catalog.IsDownloaded(input.Config, game) {
						downloaded = append(downloaded, game)
					}
				}
				if len(downloaded) > 0 {
					found = append(found, platformRoms{platform: platform, roms: downloaded})
				}
				return nil, nil
			},
		)
	}
	return found
}

func (s *MetadataSyncScreen) refresh(input MetadataSyncInput, chosen []platformRoms) {
	batches := make([]download.PlatformGames, 0, len(chosen))
	for _, entry := range chosen {
		batches = append(batches, download.PlatformGames{Platform: entry.platform, Games: entry.roms})
	}

	updated, err := gaba.ProcessMessage(
		localize("metadata_sync_updating", "Updating metadata..."),
		gaba.ProcessMessageOptions{ShowThemeBackground: true},
		func() (int, error) {
			return download.RefreshMetadata(input.Config, batches)
		},
	)
	if err != nil {
		gaba.GetLogger().Error("Metadata update failed", "error", err)
		s.tell(localize("metadata_sync_failed", "Could not update the metadata.\nCheck the logs for more info."))
		return
	}

	gaba.GetLogger().Info("Metadata update complete", "games", updated)
	s.tell(fmt.Sprintf(localize("metadata_sync_complete", "Updated metadata for %d games."), updated))
}

func (s *MetadataSyncScreen) tell(message string) {
	gaba.ConfirmationMessage(message, ContinueFooter(), gaba.MessageOptions{})
}
