package ui

import (
	"time"

	"grout/saves"

	gaba "github.com/BrandonKowalski/gabagool/v2/pkg/gabagool"
	buttons "github.com/BrandonKowalski/gabagool/v2/pkg/gabagool/constants"
)

type SyncHistoryInput struct {
	DeviceID string
}

type SyncHistoryOutput struct {
	Action SyncHistoryAction
}

type SyncHistoryScreen struct{}

func NewSyncHistoryScreen() *SyncHistoryScreen {
	return &SyncHistoryScreen{}
}

func (s *SyncHistoryScreen) Draw(input SyncHistoryInput) (SyncHistoryOutput, error) {
	output := SyncHistoryOutput{Action: SyncHistoryActionBack}

	days := saves.SyncHistory(input.DeviceID)
	if len(days) == 0 {
		gaba.ConfirmationMessage(
			localize("sync_history_empty", "No sync history found."),
			ContinueFooter(),
			gaba.MessageOptions{},
		)
		return output, nil
	}

	gaba.GetLogger().Debug("Showing sync history", "days", len(days))

	options := gaba.DefaultInfoScreenOptions()
	options.Sections = historySections(days, time.Now())
	options.ShowThemeBackground = false
	options.ShowScrollbar = true
	options.ConfirmButton = buttons.VirtualButtonUnassigned

	gaba.DetailScreen(
		localize("sync_history_title", "Sync History"),
		options,
		[]gaba.FooterHelpItem{FooterBack()},
	)

	return output, nil
}

// The icons are Material Design codepoints in the toolkit's font. The plain
// cloud heads the column that says which way each save went.
const (
	cloudOutline         = "\U000F0163"
	cloudDownloadOutline = "\U000F0B7D"
	cloudUploadOutline   = "\U000F0B7E"
)

// actionIcon shows which way a save went. An action grout does not recognise
// is written out rather than dropped, so a new one shows up as itself instead
// of as a blank cell.
func actionIcon(action string) string {
	switch action {
	case "upload":
		return cloudUploadOutline
	case "download":
		return cloudDownloadOutline
	default:
		return action
	}
}

// historySections is a table per day and platform. The platform names are
// the widest text in the history, so moving them from a column into the title
// leaves the games and times room to fit on one line.
func historySections(days []saves.SyncDay, now time.Time) []gaba.Section {
	headers := []string{
		cloudOutline,
		localize("sync_history_col_game", "Game"),
		localize("sync_history_col_time", "Time"),
	}

	var sections []gaba.Section
	for _, day := range days {
		// Events are newest first, so platforms come out in the order they
		// were last synced.
		var platforms []string
		rows := map[string][]gaba.TableRow{}
		for _, event := range day.Events {
			if _, seen := rows[event.Platform]; !seen {
				platforms = append(platforms, event.Platform)
			}
			rows[event.Platform] = append(rows[event.Platform], gaba.TableRow{Cells: []string{
				actionIcon(event.Action),
				event.RomName,
				event.At.Format("15:04"),
			}})
		}

		date := historyDate(day.Date, now)
		for _, platform := range platforms {
			sections = append(sections, gaba.NewTableSection(
				date+" · "+platform, headers, rows[platform], gaba.TableGridRowDividers))
		}
	}

	return sections
}

// historyDate keeps section titles short enough to fit beside a platform
// name: the year only when it is not this one.
func historyDate(date, now time.Time) string {
	if date.Year() == now.Year() {
		return date.Format("Jan 2")
	}
	return date.Format("Jan 2, 2006")
}
