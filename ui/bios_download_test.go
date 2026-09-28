package ui

import (
	"testing"

	"grout/bios"
	"grout/romm"
)

// Two firmware files can share a name in different folders. Staged under the
// same name, one download overwrote the other and the wrong bytes were
// installed.
func TestBIOSStagingName_UniquePerFirmware(t *testing.T) {
	a := bios.Requirement{Firmware: romm.Firmware{ID: 9, FileName: "bios.bin", FilePath: "bios/psx"}}
	b := bios.Requirement{Firmware: romm.Firmware{ID: 10, FileName: "bios.bin", FilePath: "bios/ps2"}}

	if biosStagingName(a) == biosStagingName(b) {
		t.Errorf("both stage as %q", biosStagingName(a))
	}
	if got := biosStagingName(bios.Requirement{Firmware: romm.Firmware{ID: 3, FileName: "../../x.bin"}}); got != "bios_3_x.bin" {
		t.Errorf("staging name = %q, want bios_3_x.bin", got)
	}
}
