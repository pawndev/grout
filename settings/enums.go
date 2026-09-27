package settings

type ReleaseChannel string

const (
	ReleaseChannelMatchRomM ReleaseChannel = "match_romm"
	ReleaseChannelStable    ReleaseChannel = "stable"
	ReleaseChannelBeta      ReleaseChannel = "beta"
)

type DownloadedGamesMode string

const (
	DownloadedGamesModeDoNothing DownloadedGamesMode = "do_nothing"
	DownloadedGamesModeMark      DownloadedGamesMode = "mark"
	DownloadedGamesModeFilter    DownloadedGamesMode = "filter"
)

type CollectionView string

const (
	CollectionViewPlatform CollectionView = "platform"
	CollectionViewUnified  CollectionView = "unified"
)

type LogLevel string

const (
	LogLevelDebug LogLevel = "DEBUG"
	LogLevelInfo  LogLevel = "INFO"
	LogLevelError LogLevel = "ERROR"
)

// RomStorage pins muOS roms to one card. Empty means automatic. The values
// are muOS's own names for each storage.
type RomStorage string

const (
	RomStorageAuto RomStorage = ""
	RomStorageSD1  RomStorage = "rom"
	RomStorageSD2  RomStorage = "sdcard"
	RomStorageUSB  RomStorage = "usb"
)

// RomStorageEnvVar carries the choice to the firmware layer, which reads the
// environment rather than the config, as it does for BASE_PATH.
const RomStorageEnvVar = "GROUT_MUOS_ROM_STORAGE"
