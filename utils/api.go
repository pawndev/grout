package utils

import (
	"grout/romm"
)

func GetRommClient(host romm.Host) *romm.Client {
	return romm.NewClient(host.URL(), romm.WithBasicAuth(host.Username, host.Password))
}
