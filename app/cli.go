package main

import (
	"grout/romm"
	"log"
	"os"
)

func main() {
	rc := romm.NewClient("http://192.168.1.20:1550",
		romm.WithBasicAuth(os.Getenv("DEV_ROMM_USERNAME"),
			os.Getenv("DEV_ROMM_PASSWORD")))

	saves, _ := rc.GetSaves(romm.SaveQuery{})

	log.Println(saves)
}
