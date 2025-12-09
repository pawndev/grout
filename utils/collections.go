package utils

import (
	"grout/romm"
)

func ShowCollections(host romm.Host) bool {
	rc := GetRommClient(host)
	col, err := rc.GetCollections()

	return err == nil && len(col) > 0
}
