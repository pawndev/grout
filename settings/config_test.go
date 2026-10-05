package settings

import "testing"

func TestSlotPreference_DefaultsToAutosave(t *testing.T) {
	c := Config{}
	if got := c.GetSlotPreference(1); got != "autosave" {
		t.Errorf("default slot = %q, want autosave", got)
	}
}

// Picking "autosave" must persist as an EXPLICIT preference (not be discarded), so it
// can override a sticky recorded slot. Otherwise a user can never switch a ROM back to
// autosave once another slot has been recorded (issue #250).
func TestSetSlotPreference_AutosavePersistsAsExplicit(t *testing.T) {
	c := Config{}
	c.SetSlotPreference(1, "quicksave")
	c.SetSlotPreference(1, "autosave") // user explicitly chooses autosave

	slot, ok := c.SlotPreferenceExplicit(1)
	if !ok || slot != "autosave" {
		t.Errorf("explicit autosave should persist: got (%q, %v), want (\"autosave\", true)", slot, ok)
	}
	if got := c.GetSlotPreference(1); got != "autosave" {
		t.Errorf("GetSlotPreference = %q, want autosave", got)
	}
}

// Issue #285: Non-regression tests
// RomM lowercases the folder names in PLATFORMS_BINDING but platforms and
// ROMs keep the folder's case in their fs_slug, so lookups must ignore case.
func TestResolveFSSlug_IgnoresBindingKeyCase(t *testing.T) {
	c := Config{PlatformsBinding: map[string]string{"game boy advance": "gba"}}
	if got := c.ResolveFSSlug("Game Boy Advance"); got != "gba" {
		t.Errorf("ResolveFSSlug = %q, want gba", got)
	}
	if got := c.ResolveFSSlug("Unbound"); got != "Unbound" {
		t.Errorf("unbound slug = %q, want it unchanged", got)
	}
}

func TestGetDirectoryMapping_IgnoresSlugCase(t *testing.T) {
	c := Config{DirectoryMappings: map[string]DirectoryMapping{
		"Game Boy Advance": {RomMSlug: "Game Boy Advance", RelativePath: "Nintendo Game Boy Advance"},
	}}
	got, ok := c.GetDirectoryMapping("game boy advance")
	if !ok || got != "Nintendo Game Boy Advance" {
		t.Errorf("GetDirectoryMapping = (%q, %v), want (\"Nintendo Game Boy Advance\", true)", got, ok)
	}
}
