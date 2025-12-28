package main

import (
	"fmt"
	"grout/utils"
	"os"
	"time"

	gaba "github.com/BrandonKowalski/gabagool/v2/pkg/gabagool"
	_ "github.com/UncleJunVIP/certifiable"
)

func main() {
	defer cleanup()

	appStart := time.Now()
	result := setup()
	config := result.Config
	platforms := result.Platforms

	logger := gaba.GetLogger()
	logger.Debug("Starting Grout")

	cfw := utils.GetCFW()
	quitOnBack := len(config.Hosts) == 1
	showCollections := utils.ShowCollections(config, config.Hosts[0])

	fsmStart := time.Now()
	fsm := buildFSM(config, cfw, platforms, quitOnBack, showCollections, appStart)
	logger.Debug("FSM built", "seconds", fmt.Sprintf("%.2f", time.Since(fsmStart).Seconds()))

	logger.Info("Starting FSM.Run()")
	if err := fsm.Run(); err != nil {
		logger.Error("FSM error", "error", err)
	}
}

func cleanup() {
	if err := os.RemoveAll(".tmp"); err != nil {
		gaba.GetLogger().Error("Failed to clean .tmp directory", "error", err)
	}
	gaba.Close()
}
