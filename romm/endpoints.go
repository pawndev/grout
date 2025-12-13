package romm

const (
	endpointLogin = "/api/login"

	endpointPlatforms    = "/api/platforms"
	endpointPlatformByID = "/api/platforms/%d"

	endpointRoms         = "/api/roms"
	endpointRomByID      = "/api/roms/%d"
	endpointRomContent   = "/api/roms/%d/content/%s"
	endpointRomsDownload = "/api/roms/download"

	endpointRomFileByID    = "/api/romsfiles/%d"
	endpointRomFileContent = "/api/romsfiles/%d/content/%s"

	endpointCollections           = "/api/collections"
	endpointCollectionByID        = "/api/collections/%d"
	endpointSmartCollections      = "/api/collections/smart"
	endpointSmartCollectionByID   = "/api/collections/smart/%d"
	endpointVirtualCollections    = "/api/collections/virtual"
	endpointVirtualCollectionByID = "/api/collections/virtual/%d"

	endpointFirmware = "/api/firmware"

	endpointSaves = "/api/saves"
)
