package ui

import "testing"

// The QR, the instructions and the footer share the screen, so the QR takes
// half its height at most. At a fixed 320 the text ran into the footer on a
// 480 line Miyoo.
func TestPairingQRSize(t *testing.T) {
	for _, tc := range []struct {
		screenHeight int32
		want         int
	}{
		{480, 240},
		{560, 280},
		{720, 320},
		{1080, 320},
	} {
		if got := pairingQRSize(tc.screenHeight); got != tc.want {
			t.Errorf("screen %d: QR %d, want %d", tc.screenHeight, got, tc.want)
		}
	}
}
