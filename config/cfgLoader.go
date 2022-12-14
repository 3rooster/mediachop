package config

import (
	"go.uber.org/zap"
	"gopkg.in/yaml.v3"
	"io/ioutil"
)

var MediaServer *MediaServerConfig
var Env *env

func init() {
	initRootDir()
}

// LoadConfig load config from file
func LoadConfig(path string) error {
	// log to console
	initLogger(nil)

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
	return processLogCfg(config.Logger)
}

func processLogCfg(config *zap.Config) error {

	for i := 0; i < len(config.OutputPaths); i++ {
		config.OutputPaths[i] = GetAbsPath(config.OutputPaths[i])
	}
	for i := 0; i < len(config.ErrorOutputPaths); i++ {
		config.ErrorOutputPaths[i] = GetAbsPath(config.ErrorOutputPaths[i])
	}
	if config.Encoding == "" {
		config.Encoding = "json"
		config.EncoderConfig = zap.NewProductionEncoderConfig()
	}
	return initLogger(config)
}
