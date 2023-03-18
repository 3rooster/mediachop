package main

import (
	"go.uber.org/zap"
	"mediachop/config"
	"mediachop/mediaServer"
	"mediachop/service/cache"
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
	cache.InitDefault(config.Cache.CommonCache)
	mediaStore.Init(config.Cache.Stream, config.Cache.MediaFile)
	mediaServer.Start(config.MediaServer)
}
