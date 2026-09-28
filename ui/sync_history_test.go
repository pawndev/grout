package ui

import (
	"testing"
	"time"

	"grout/saves"
)

// Each day is split into a table per platform, in the order the platforms
// were last synced, and the platform moves from a column to the title.
func TestHistorySections_OneTablePerPlatform(t *testing.T) {
	day := time.Date(2026, 9, 27, 0, 0, 0, 0, time.Local)
	at := func(h, m int) time.Time { return day.Add(time.Duration(h)*time.Hour + time.Duration(m)*time.Minute) }

	sections := historySections([]saves.SyncDay{{
		Date: day,
		Events: []saves.SyncEvent{
			{Action: "upload", RomName: "Advance Wars", Platform: "Game Boy Advance", At: at(18, 30)},
			{Action: "download", RomName: "Chrono Trigger", Platform: "Super Nintendo", At: at(18, 29)},
			{Action: "upload", RomName: "Mother 3", Platform: "Game Boy Advance", At: at(18, 28)},
		},
	}}, day)

	if len(sections) != 2 {
		t.Fatalf("got %d sections, want one per platform", len(sections))
	}
	if want := "Sep 27 · Game Boy Advance"; sections[0].Title != want {
		t.Errorf("first title = %q, want %q", sections[0].Title, want)
	}
	if want := "Sep 27 · Super Nintendo"; sections[1].Title != want {
		t.Errorf("second title = %q, want %q", sections[1].Title, want)
	}

	gba := sections[0]
	if len(gba.TableHeaders) != 3 {
		t.Errorf("headers = %v, want icon, game and time only", gba.TableHeaders)
	}
	if len(gba.TableRows) != 2 || gba.TableRows[0].Cells[1] != "Advance Wars" || gba.TableRows[1].Cells[1] != "Mother 3" {
		t.Errorf("GBA rows = %+v, want both GBA games, newest first", gba.TableRows)
	}
	if gba.TableRows[0].Cells[2] != "18:30" {
		t.Errorf("time cell = %q, want 18:30", gba.TableRows[0].Cells[2])
	}
}

// The year is left off for this year's days and kept for older ones.
func TestHistoryDate(t *testing.T) {
	now := time.Date(2026, 9, 27, 12, 0, 0, 0, time.Local)
	if got := historyDate(time.Date(2026, 3, 4, 0, 0, 0, 0, time.Local), now); got != "Mar 4" {
		t.Errorf("this year = %q, want Mar 4", got)
	}
	if got := historyDate(time.Date(2025, 12, 31, 0, 0, 0, 0, time.Local), now); got != "Dec 31, 2025" {
		t.Errorf("last year = %q, want Dec 31, 2025", got)
	}
}
