package romm

const (
	EndpointLogin = "/api/login"

	EndpointPlatforms    = "/api/platforms"
	EndpointPlatformByID = "/api/platforms/%d"

	EndpointRoms         = "/api/roms"
	EndpointRomByID      = "/api/roms/%d"
	EndpointRomContent   = "/api/roms/%d/content/%s"
	EndpointRomsDownload = "/api/roms/download"

	EndpointRomFileByID    = "/api/romsfiles/%d"
	EndpointRomFileContent = "/api/romsfiles/%d/content/%s"

	EndpointCollections           = "/api/collections"
	EndpointCollectionByID        = "/api/collections/%d"
	EndpointSmartCollections      = "/api/collections/smart"
	EndpointSmartCollectionByID   = "/api/collections/smart/%d"
	EndpointVirtualCollections    = "/api/collections/virtual"
	EndpointVirtualCollectionByID = "/api/collections/virtual/%d"

	EndpointFirmware = "/api/firmware"

	EndpointSaves = "/api/saves"
)
