package config

import "go.uber.org/zap"

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
	Env      *env               `yaml:"env"`
	Server   *MediaServerConfig `yaml:"mediaServer"`
	Logger   *zap.Config        `yaml:"logger"`
	CacheCfg *cacheCfg          `yaml:"mediaCache"`
}

type cacheCfg struct {
	CacheTTLSec      int64 `yaml:"cache_ttl_sec"`
	ClearIntervalSec int   `yaml:"clear_interval_sec"`
	Verbose          bool  `yaml:"verbose"`
}
