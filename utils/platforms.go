package utils

import (
	"fmt"

	"grout/romm"
)

func GetMappedPlatforms(host romm.Host, mappings map[string]DirectoryMapping) []romm.Platform {
	c := GetRommClient(host)

	rommPlatforms, err := c.GetPlatforms()
	if err != nil {
		LogStandardFatal(fmt.Sprintf("Failed to get platforms from RomM: %s", err), nil)
	}

	var platforms []romm.Platform

	for _, platform := range rommPlatforms {
		_, exists := mappings[platform.Slug]
		if exists {
			platforms = append(platforms, romm.Platform{
				Name: platform.Name,
				ID:   platform.ID,
				Slug: platform.Slug,
			})
		}
	}

	return platforms
}
