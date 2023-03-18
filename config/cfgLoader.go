package config

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"mediachop/service/cache"
)

var MediaServer *MediaServerConfig
var Env *env

var Cache = &cacheConfig{
	Cache: &cache.Config{
		ClearIntervalMs: 5000,
		DefaultTTLMs:    5000,
		Shards:          8,
	},
	Stream: &cache.Config{
		ClearIntervalMs: 60 * 1000,
		DefaultTTLMs:    300 * 1000,
		Shards:          1,
	},
	MediaFile: &cache.Config{
		ClearIntervalMs: 10 * 1000,
		DefaultTTLMs:    60 * 1000,
		Shards:          1,
	},
}

func init() {
	initRootDir()
}

// LoadConfig load config from file
func LoadConfig(path string) error {
	// log to console
	initLogger(nil, nil)

	// load config
	finalPath := GetAbsPath(path)
	zap.S().Infof("load config from %s", finalPath)
	content, err := ioutil.ReadFile(finalPath)
	if err != nil {
		zap.S().Infof("failed to load config, err=%v", err)
		return err
	}
	var config cfg
	err = yaml.Unmarshal([]byte(content), &config)

	if err != nil {
		zap.S().Infof("failed to parse config file, path=%s", path)
		return err
	}
	Env = config.Env
	MediaServer = config.Server
	if config.CacheCfg != nil && config.CacheCfg.Cache != nil {
		Cache.Cache = config.CacheCfg.Cache
	}
	if config.CacheCfg != nil && config.CacheCfg.Stream != nil {
		Cache.Stream = config.CacheCfg.Stream
	}
	if config.CacheCfg != nil && config.CacheCfg.MediaFile != nil {
		Cache.MediaFile = config.CacheCfg.MediaFile
	}
	return processLogCfg(config.Logger, config.RotationConfig)
}

func processLogCfg(config *zap.Config, rotationConfig *lumberjack.Logger) error {

	for i := 0; i < len(config.OutputPaths); i++ {
		config.OutputPaths[i] = GetAbsPath(config.OutputPaths[i])
	}
	for i := 0; i < len(config.ErrorOutputPaths); i++ {
		config.ErrorOutputPaths[i] = GetAbsPath(config.ErrorOutputPaths[i])
	}
	if config.Encoding == "" {
		config.Encoding = "json"
		config.EncoderConfig = zap.NewProductionEncoderConfig()
		config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	}
	return initLogger(config, rotationConfig)
}
