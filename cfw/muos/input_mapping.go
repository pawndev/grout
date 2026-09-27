package muos

import (
	"embed"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

//go:embed input_mappings/*.json
var embeddedInputMappings embed.FS

// Device represents the detected device type on muOS
type Device string

const (
	DeviceAnbernic       Device = "anbernic"
	DeviceTrimui         Device = "trimui"
	DeviceTrimuiSmartPro Device = "trimui-smart-pro"
)

// DetectDevice detects the device type when running on muOS by checking input devices.
// Returns DeviceTrimui if "TRIMUI" is found in /proc/bus/input/devices, otherwise DeviceAnbernic.
func DetectDevice() Device {
	logger := slog.Default()
	logger.Info("Detecting muOS device type...")

	cmd := exec.Command("sh", "-c", "cat /proc/bus/input/devices | grep TRIMUI")
	output, err := cmd.Output()

	if err != nil || len(output) == 0 {
		return DeviceAnbernic
	}

	if len(output) >= 1 {
		for _, line := range strings.Split(string(output), "\n") {
			if strings.Contains(line, "TRIMUI Smart Pro") {
				return DeviceTrimuiSmartPro
			}
		}
	}

	return DeviceTrimui
}

// GetInputMappingBytes returns the embedded input mapping JSON for the detected muOS device
func GetInputMappingBytes() ([]byte, error) {
	device := DetectDevice()
	return GetInputMappingBytesForDevice(device)
}

// GetInputMappingBytesForDevice returns the embedded input mapping JSON for a specific device
func GetInputMappingBytesForDevice(device Device) ([]byte, error) {
	var filename string
	switch device {
	case DeviceAnbernic:
		if _, err := os.Stat(filepath.Join("overrides", "cfw", "muos", "input_mappings", "anbernic.json")); err != nil {
			return anbernicMapping("/"), nil
		}
		filename = "input_mappings/anbernic.json"
	case DeviceTrimui:
		filename = "input_mappings/trimui.json"
	case DeviceTrimuiSmartPro:
		filename = "input_mappings/trimui-smart-pro.json"
	default:
		filename = "input_mappings/anbernic.json"
	}

	overridePath := filepath.Join("overrides", "cfw", "muos", filename)
	data, err := os.ReadFile(overridePath)
	if err != nil {
		data, err = embeddedInputMappings.ReadFile(filename)
		if err != nil {
			return nil, fmt.Errorf("failed to read embedded input mapping %s: %w", filename, err)
		}
	}

	return data, nil
}

// sdlMapFile is where muOS keeps the SDL mapping for the device's own pad.
const sdlMapFile = "/opt/muos/device/config/board/sdl_map"

// Virtual button values, as the toolkit numbers them. The embedded mappings
// use the same numbers.
const (
	buttonUp     = 1
	buttonDown   = 2
	buttonLeft   = 3
	buttonRight  = 4
	buttonA      = 5
	buttonB      = 6
	buttonX      = 7
	buttonY      = 8
	buttonL1     = 9
	buttonL2     = 10
	buttonR1     = 11
	buttonR2     = 12
	buttonStart  = 13
	buttonSelect = 14
	buttonMenu   = 15
)

// sdlButtons maps the names in muOS's SDL line to virtual buttons. On
// Anbernic devices muOS names the face buttons as printed on the pad.
var sdlButtons = map[string]int{
	"a": buttonA, "b": buttonB, "x": buttonX, "y": buttonY,
	"leftshoulder": buttonL1, "rightshoulder": buttonR1,
	"lefttrigger": buttonL2, "righttrigger": buttonR2,
	"start": buttonStart, "back": buttonSelect, "guide": buttonMenu,
}

var sdlHats = map[string]int{"dpup": buttonUp, "dpdown": buttonDown, "dpleft": buttonLeft, "dpright": buttonRight}

// anbernicMapping follows muOS's own SDL line for the pad. muOS renumbered its
// buttons in August 2026, so a fixed mapping put B on A. Without the file, as
// on older muOS, the embedded mapping is used.
func anbernicMapping(root string) []byte {
	if line, err := os.ReadFile(filepath.Join(root, sdlMapFile)); err == nil {
		if mapping, ok := mappingFromSDLMap(strings.TrimSpace(string(line))); ok {
			return mapping
		}
	}
	data, _ := embeddedInputMappings.ReadFile("input_mappings/anbernic.json")
	return data
}

// mappingFromSDLMap turns an SDL controller line, "guid,name,a:b2,...", into
// grout's joystick mapping.
func mappingFromSDLMap(line string) ([]byte, bool) {
	buttons := map[string]int{}
	hats := map[string]int{}

	fields := strings.Split(line, ",")
	if len(fields) < 3 {
		return nil, false
	}
	for _, field := range fields[2:] {
		name, value, ok := strings.Cut(field, ":")
		if !ok {
			continue
		}
		if button, known := sdlButtons[name]; known && strings.HasPrefix(value, "b") {
			buttons[strings.TrimPrefix(value, "b")] = button
		}
		// Hats are "h0.<mask>"; the mask is what the toolkit keys on.
		if button, known := sdlHats[name]; known && strings.HasPrefix(value, "h0.") {
			hats[strings.TrimPrefix(value, "h0.")] = button
		}
	}
	if len(buttons) == 0 {
		return nil, false
	}

	mapping, err := json.Marshal(map[string]any{
		"keyboard_map":          map[string]int{},
		"controller_button_map": map[string]int{},
		"controller_hat_map":    map[string]int{},
		"joystick_axis_map":     map[string]int{},
		"joystick_button_map":   buttons,
		"joystick_hat_map":      hats,
	})
	return mapping, err == nil
}
