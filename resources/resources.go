package resources

import (
	"embed"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/BrandonKowalski/gabagool/v2/pkg/gabagool/i18n"
)

//go:embed locales/*.toml splash.png input_mappings/*.json
var embeddedFiles embed.FS

// MuOSDevice represents the detected device type on muOS
type MuOSDevice string

const (
	MuOSDeviceAnbernic       MuOSDevice = "anbernic"
	MuOSDeviceTrimuiBrick    MuOSDevice = "trimui_brick"
	MuOSDeviceTrimuiSmartPro MuOSDevice = "trimui_smart_pro"
)

const trimuiMainUIPath = "/usr/trimui/bin/MainUI"

// DetectMuOSDevice detects the device type when running on muOS.
// Logic:
//   - If /usr/trimui/bin/MainUI doesn't exist → Anbernic
//   - If it exists and `strings MainUI | grep ^Trimui` returns "Trimui Brick" → Brick
//   - Otherwise → Smart Pro
func DetectMuOSDevice() MuOSDevice {
	if _, err := os.Stat(trimuiMainUIPath); os.IsNotExist(err) {
		return MuOSDeviceAnbernic
	}

	// File exists, check if it's a Trimui Brick
	cmd := exec.Command("sh", "-c", fmt.Sprintf("strings %s | grep ^Trimui", trimuiMainUIPath))
	output, err := cmd.Output()
	if err != nil {
		// If strings/grep fails, default to Smart Pro
		return MuOSDeviceTrimuiSmartPro
	}

	trimmedOutput := strings.TrimSpace(string(output))
	if trimmedOutput == "Trimui Brick" {
		return MuOSDeviceTrimuiBrick
	}

	return MuOSDeviceTrimuiSmartPro
}

// GetMuOSInputMappingBytes returns the embedded input mapping JSON for the detected muOS device
func GetMuOSInputMappingBytes() ([]byte, error) {
	device := DetectMuOSDevice()
	return GetInputMappingBytesForDevice(device)
}

// GetInputMappingBytesForDevice returns the embedded input mapping JSON for a specific device
func GetInputMappingBytesForDevice(device MuOSDevice) ([]byte, error) {
	var filename string
	switch device {
	case MuOSDeviceAnbernic:
		filename = "input_mappings/anbernic.json"
	case MuOSDeviceTrimuiBrick:
		filename = "input_mappings/trimui_brick.json"
	case MuOSDeviceTrimuiSmartPro:
		filename = "input_mappings/trimui_smart_pro.json"
	default:
		filename = "input_mappings/anbernic.json"
	}

	data, err := embeddedFiles.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read embedded input mapping %s: %w", filename, err)
	}
	return data, nil
}

// LocaleFile represents a locale file name and path
type LocaleFile struct {
	Name string
	Path string
}

var localeFiles = []LocaleFile{
	{Name: "active.en.toml", Path: "locales/active.en.toml"},
	{Name: "active.es.toml", Path: "locales/active.es.toml"},
	{Name: "active.fr.toml", Path: "locales/active.fr.toml"},
}

// GetLocaleMessageFiles returns locale files as MessageFile structs for i18n initialization
func GetLocaleMessageFiles() ([]i18n.MessageFile, error) {
	var messageFiles []i18n.MessageFile

	for _, localeFile := range localeFiles {
		content, err := embeddedFiles.ReadFile(localeFile.Path)
		if err != nil {
			return nil, fmt.Errorf("failed to read embedded locale file %s: %w", localeFile.Path, err)
		}

		messageFiles = append(messageFiles, i18n.MessageFile{
			Name:    localeFile.Name,
			Content: content,
		})
	}

	return messageFiles, nil
}

// GetSplashImageBytes returns the splash screen image data as a byte slice
// The splash image is 1024x720 PNG (aspect ratio 1.42:1)
func GetSplashImageBytes() ([]byte, error) {
	data, err := embeddedFiles.ReadFile("splash.png")
	if err != nil {
		return nil, fmt.Errorf("failed to read embedded splash image: %w", err)
	}
	return data, nil
}
