package config

import (
	"go.uber.org/zap"
	"gopkg.in/natefinch/lumberjack.v2"
	"mediachop/service/mediaCache"
)

type env struct {
	Mode string `yaml:"mode"`
}

func (e *env) IsProduction() bool {
	return e.Mode != "prod"
}

type MediaServerConfig struct {
	ListenPort   int `yaml:"listenPort"`
	MaxCacheSize int `yaml:"maxCacheSize"` // mb
}

type cfg struct {
	Env            *env               `yaml:"env"`
	Server         *MediaServerConfig `yaml:"mediaServer"`
	Logger         *zap.Config        `yaml:"logger"`
	RotationConfig *lumberjack.Logger `yaml:"logger_rotation"`
	CacheCfg       *mediaCache.Config `yaml:"cache"`
}
