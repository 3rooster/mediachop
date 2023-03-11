package config

import (
	"go.uber.org/zap"
	"gopkg.in/natefinch/lumberjack.v2"
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
	CacheCfg       *cacheCfg          `yaml:"cache"`
}

type cacheCfg struct {
	CacheTTLSec      int  `yaml:"cache_ttl_sec"`
	ClearIntervalSec int  `yaml:"clear_interval_sec"`
	Verbose          bool `yaml:"verbose"`
}
