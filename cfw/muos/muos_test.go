package muos

import (
	"os"
	"path/filepath"
	"testing"

	"grout/settings"
)

// device lays out a muOS root with the given storage mounted, each holding a
// ROMS folder, and the storage config muOS keeps for them.
func device(t *testing.T, withRoms ...string) string {
	t.Helper()
	root := t.TempDir()

	for _, storage := range []string{"rom", "sdcard", "usb"} {
		mount := map[string]string{"rom": "/mnt/mmc", "sdcard": "/mnt/sdcard", "usb": "/mnt/usb"}[storage]
		config := filepath.Join(root, storageConfigDir, storage)
		if err := os.MkdirAll(config, 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(config, "mount"), []byte(mount+"\n"), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	for _, mount := range withRoms {
		if err := os.MkdirAll(filepath.Join(root, mount, "ROMS"), 0o755); err != nil {
			t.Fatal(err)
		}
	}
	return root
}

// Without the union mount, roms are looked for on USB, then SD2, then SD1,
// the order the union took writes in.
func TestRomDirectory_FollowsTheUnionOrder(t *testing.T) {
	for _, tc := range []struct {
		name     string
		withRoms []string
		want     string
	}{
		{"only SD1", []string{"/mnt/mmc"}, "/mnt/mmc/ROMS"},
		{"SD2 before SD1", []string{"/mnt/mmc", "/mnt/sdcard"}, "/mnt/sdcard/ROMS"},
		{"USB before both", []string{"/mnt/mmc", "/mnt/sdcard", "/mnt/usb"}, "/mnt/usb/ROMS"},
		{"nothing anywhere falls back to SD1", nil, "/mnt/mmc/ROMS"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			root := device(t, tc.withRoms...)
			if got := romDirectory(root, ""); got != filepath.Join(root, tc.want) {
				t.Errorf("rom directory = %s, want %s", got, tc.want)
			}
		})
	}
}

// An older muOS still merges every card at /mnt/union, and that wins.
func TestRomDirectory_UsesTheUnionWhenItExists(t *testing.T) {
	root := device(t, "/mnt/mmc", "/mnt/sdcard")
	if err := os.MkdirAll(filepath.Join(root, RomsFolderUnion), 0o755); err != nil {
		t.Fatal(err)
	}

	if got := romDirectory(root, settings.RomStorageSD2); got != filepath.Join(root, RomsFolderUnion) {
		t.Errorf("rom directory = %s, want the union", got)
	}
}

// The setting picks a card, but one without a ROMS folder cannot hold the
// library, so the automatic order takes over.
func TestRomDirectory_HonoursTheChosenStorage(t *testing.T) {
	root := device(t, "/mnt/mmc", "/mnt/sdcard")

	if got := romDirectory(root, settings.RomStorageSD1); got != filepath.Join(root, "/mnt/mmc/ROMS") {
		t.Errorf("SD1 chosen: %s", got)
	}
	if got := romDirectory(root, settings.RomStorageUSB); got != filepath.Join(root, "/mnt/sdcard/ROMS") {
		t.Errorf("USB chosen but has no roms: %s, want the automatic SD2", got)
	}
}

// muOS reads each mount point from its own config, so a device that mounts a
// card elsewhere is followed there.
func TestRomDirectory_ReadsMountPointsFromMuOS(t *testing.T) {
	root := device(t)
	config := filepath.Join(root, storageConfigDir, "sdcard", "mount")
	if err := os.WriteFile(config, []byte("/media/card2\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(filepath.Join(root, "/media/card2/ROMS"), 0o755); err != nil {
		t.Fatal(err)
	}

	if got := romDirectory(root, ""); got != filepath.Join(root, "/media/card2/ROMS") {
		t.Errorf("rom directory = %s, want the configured SD2 mount", got)
	}
}
