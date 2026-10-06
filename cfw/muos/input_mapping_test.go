package muos

import (
	"encoding/json"
	"reflect"
	"testing"
)

// The line muOS shipped for Anbernic devices until August 2026.
const sdlMapBefore = "19000000010000000100000000010000,muOS-Keys,a:b3,b:b4,x:b6,y:b5,leftshoulder:b7,rightshoulder:b8,lefttrigger:b12,righttrigger:b13,guide:b11,start:b10,back:b9,dpup:h0.1,dpleft:h0.8,dpright:h0.2,dpdown:h0.4,volumedown:b1,volumeup:b2,platform:Linux,"

// The line it ships now, with every button moved.
const sdlMapNow = "060000006d750000070700006f730000,muOS-Keys,a:b2,b:b3,x:b4,y:b5,leftshoulder:b6,rightshoulder:b7,lefttrigger:b8,righttrigger:b9,guide:b12,start:b11,back:b10,dpup:h0.1,dpleft:h0.8,dpright:h0.2,dpdown:h0.4,volumedown:b0,volumeup:b1,platform:Linux,"

type buttonMaps struct {
	Buttons map[string]int `json:"joystick_button_map"`
	Hats    map[string]int `json:"joystick_hat_map"`
}

func decode(t *testing.T, data []byte) buttonMaps {
	t.Helper()
	var maps buttonMaps
	if err := json.Unmarshal(data, &maps); err != nil {
		t.Fatal(err)
	}
	return maps
}

// Built from the line muOS used to ship, the mapping is the one grout has
// always embedded, so devices that have not updated see no change.
func TestMappingFromSDLMap_MatchesTheEmbeddedOne(t *testing.T) {
	built, ok := mappingFromSDLMap(sdlMapBefore)
	if !ok {
		t.Fatal("the old muOS line did not parse")
	}
	embedded, err := embeddedInputMappings.ReadFile("input_mappings/anbernic.json")
	if err != nil {
		t.Fatal(err)
	}
	if got, want := decode(t, built), decode(t, embedded); !reflect.DeepEqual(got, want) {
		t.Errorf("built %+v\nwant %+v", got, want)
	}
}

// muOS renumbered the buttons. Following its line puts A back on A, where the
// embedded mapping had B on A and X on B.
func TestMappingFromSDLMap_FollowsTheNewNumbering(t *testing.T) {
	built, ok := mappingFromSDLMap(sdlMapNow)
	if !ok {
		t.Fatal("the new muOS line did not parse")
	}
	buttons := decode(t, built).Buttons
	want := map[string]int{
		"2": buttonA, "3": buttonB, "4": buttonX, "5": buttonY,
		"6": buttonL1, "7": buttonR1, "8": buttonL2, "9": buttonR2,
		"10": buttonSelect, "11": buttonStart, "12": buttonMenu,
	}
	if !reflect.DeepEqual(buttons, want) {
		t.Errorf("buttons = %v, want %v", buttons, want)
	}
}

func TestMappingFromSDLMap_RejectsWhatIsNotAMapping(t *testing.T) {
	for _, line := range []string{"", "not,a,mapping", "guid,name,platform:Linux"} {
		if _, ok := mappingFromSDLMap(line); ok {
			t.Errorf("%q was accepted", line)
		}
	}
}

// Anbernic devices follow muOS's own file when it has one, and the embedded
// mapping when it does not.
func TestAnbernicMapping_PrefersMuOSFile(t *testing.T) {
	root := t.TempDir()
	if got := decode(t, anbernicMapping(root)).Buttons["3"]; got != buttonA {
		t.Errorf("without muOS's file, 3 = %d, want the embedded A", got)
	}

	write(t, root, sdlMapFile, sdlMapNow)
	if got := decode(t, anbernicMapping(root)).Buttons["2"]; got != buttonA {
		t.Errorf("with muOS's file, 2 = %d, want A", got)
	}
}
