package main

import (
	"go.uber.org/zap"
	"mediachop/config"
	"mediachop/mediaServer"
	"mediachop/service/mediaStore"
	"os"
)

func initDepends() {
	// init log without logger config
	err := config.LoadConfig("conf/mediachop.yaml")
	if err != nil {
		os.Exit(0)
	}
}

func main() {
	initDepends()
	defer zap.L().Sync()
	mediaStore.Init()
	mediaServer.Start(config.MediaServer)
}
