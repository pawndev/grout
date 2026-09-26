package ui

import (
	"testing"

	"grout/cache"

	gaba "github.com/BrandonKowalski/gabagool/v2/pkg/gabagool"
)

func categoryByKey(t *testing.T, key string) filterCategory {
	t.Helper()
	for _, category := range filterCategories {
		if category.key == key {
			return category
		}
	}
	t.Fatalf("no filter category %q", key)
	return filterCategory{}
}

// Every category has to reach its own field of the filter. They used to be
// wired by position across four separate switches and a fixed-size array, so
// adding one in the wrong order sent a filter to the wrong field.
func TestFilterCategories_GetAndSetAgree(t *testing.T) {
	for _, category := range filterCategories {
		t.Run(category.key, func(t *testing.T) {
			var filter cache.GameFilter
			category.set(&filter, []string{"chosen"})

			got := category.get(filter)
			if len(got) != 1 || got[0] != "chosen" {
				t.Errorf("get after set = %v, want the value that was set", got)
			}

			category.set(&filter, nil)
			if len(category.get(filter)) != 0 {
				t.Errorf("clearing left %v behind", category.get(filter))
			}
		})
	}
}

// Two categories writing the same field would make one of them silently
// overwrite the other whenever both are chosen.
func TestFilterCategories_WriteDistinctFields(t *testing.T) {
	for _, category := range filterCategories {
		var filter cache.GameFilter
		category.set(&filter, []string{"chosen"})

		for _, other := range filterCategories {
			if other.key == category.key {
				continue
			}
			if len(other.get(filter)) != 0 {
				t.Errorf("setting %q also filled %q", category.key, other.key)
			}
		}
	}
}

func TestFilterCategories_KeysAreUnique(t *testing.T) {
	seen := make(map[string]bool)
	for _, category := range filterCategories {
		if category.key == "" {
			t.Errorf("category %q has no key", category.labelDefault)
		}
		if seen[category.key] {
			t.Errorf("two categories share the key %q", category.key)
		}
		seen[category.key] = true
	}
}

// Only the platform row is answered by the collection rather than by a
// metadata table, and it needs all three table names or none.
func TestFilterCategories_TableNames(t *testing.T) {
	for _, category := range filterCategories {
		if category.isPlatform() {
			if category.junctionTable != "" || category.fkCol != "" {
				t.Errorf("%q has no lookup table but names other tables", category.key)
			}
			continue
		}
		if category.junctionTable == "" || category.fkCol == "" {
			t.Errorf("%q names a lookup table but not the rest", category.key)
		}
	}
}

func row(key, value string) gaba.ItemWithOptions {
	return gaba.ItemWithOptions{
		Item:           gaba.MenuItem{Metadata: key},
		Options:        []gaba.Option{{Value: ""}, {Value: value}},
		SelectedOption: 1,
	}
}

func TestSelectedFilter(t *testing.T) {
	filter := selectedFilter([]gaba.ItemWithOptions{
		row("genre", "Platform"),
		row("region", "USA"),
		row("platform", "snes"),
	})

	if len(filter.Genres) != 1 || filter.Genres[0] != "Platform" {
		t.Errorf("Genres = %v, want the chosen genre", filter.Genres)
	}
	if len(filter.Regions) != 1 || filter.Regions[0] != "USA" {
		t.Errorf("Regions = %v, want the chosen region", filter.Regions)
	}
	if len(filter.PlatformSlugs) != 1 || filter.PlatformSlugs[0] != "snes" {
		t.Errorf("PlatformSlugs = %v, want the chosen platform", filter.PlatformSlugs)
	}
	if len(filter.Tags) != 0 {
		t.Errorf("Tags = %v, want nothing for a category that was not shown", filter.Tags)
	}
}

// The leading All entry means the category is not filtering, which is not the
// same as filtering by an empty value.
func TestSelectedFilter_AllMeansNoFilter(t *testing.T) {
	item := row("genre", "Platform")
	item.SelectedOption = 0

	if got := selectedFilter([]gaba.ItemWithOptions{item}); len(got.Genres) != 0 {
		t.Errorf("Genres = %v, want none when All is selected", got.Genres)
	}
}

// A row the screen did not build has no key, and must not be mistaken for one
// that does.
func TestSelectedFilter_IgnoresUnknownRows(t *testing.T) {
	item := row("not_a_category", "value")

	got := selectedFilter([]gaba.ItemWithOptions{item})
	for _, category := range filterCategories {
		if values := category.get(got); len(values) != 0 {
			t.Errorf("an unrecognised row filled %q with %v", category.key, values)
		}
	}
}

// Choosing one filter narrows the others. A row whose value survives that
// must stay on it rather than jumping back to All.
func TestKeepSelection_KeepsASurvivingValue(t *testing.T) {
	item := row("genre", "Platform")

	keepSelection(&item, []gaba.Option{{Value: ""}, {Value: "Adventure"}, {Value: "Platform"}})

	if got := selectedValue(item); got != "Platform" {
		t.Errorf("selection = %q, want it left on the value the user chose", got)
	}
}

// A value that no longer matches anything has to fall back to All, or the row
// would point past the end of its own options.
func TestKeepSelection_FallsBackWhenTheValueIsGone(t *testing.T) {
	item := row("genre", "Platform")

	keepSelection(&item, []gaba.Option{{Value: ""}, {Value: "Adventure"}})

	if item.SelectedOption != 0 {
		t.Errorf("selection = %d, want All", item.SelectedOption)
	}
	if got := selectedValue(item); got != "" {
		t.Errorf("value = %q, want nothing filtered", got)
	}
}

func TestChosenValue(t *testing.T) {
	if got := chosenValue([]string{"one"}); got != "one" {
		t.Errorf("chosenValue = %q, want the single value", got)
	}
	if got := chosenValue(nil); got != "" {
		t.Errorf("chosenValue = %q, want nothing", got)
	}
	// The screen offers one value at a time, so a filter carrying several came
	// from somewhere else and none of them is "the" selection.
	if got := chosenValue([]string{"one", "two"}); got != "" {
		t.Errorf("chosenValue = %q, want nothing for a multi-value filter", got)
	}
}

// The downloaded row opens on the setting, or on the override once there is
// one, and only a change from the setting comes back as an override.
func TestDownloadedRow(t *testing.T) {
	hide, show := true, false

	for _, tc := range []struct {
		name          string
		hideByDefault bool
		current       *bool
		pick          bool // true picks Show
		wantOpen      bool // true opens on Show
		wantOverride  *bool
	}{
		{"untouched, setting shows", false, nil, true, true, nil},
		{"untouched, setting hides", true, nil, false, false, nil},
		{"hide while the setting shows", false, nil, false, true, &hide},
		{"show while the setting hides", true, nil, true, false, &show},
		{"reopens on the override", false, &hide, false, false, &hide},
		{"back to the setting drops the override", false, &hide, true, false, nil},
	} {
		t.Run(tc.name, func(t *testing.T) {
			item := downloadedItem(cache.GameFilter{HideDownloaded: tc.current}, tc.hideByDefault)

			if opened := item.Options[item.SelectedOption].Value.(bool); opened != tc.wantOpen {
				t.Errorf("opened on show=%v, want %v", opened, tc.wantOpen)
			}

			item.SelectedOption = optionIndex(item.Options, tc.pick)
			got := downloadedOverride([]gaba.ItemWithOptions{item}, tc.hideByDefault)

			switch {
			case tc.wantOverride == nil && got != nil:
				t.Errorf("override = %v, want none", *got)
			case tc.wantOverride != nil && (got == nil || *got != *tc.wantOverride):
				t.Errorf("override = %v, want hide=%v", got, *tc.wantOverride)
			}
		})
	}
}
