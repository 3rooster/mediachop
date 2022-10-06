package config

import (
	"go.uber.org/zap"
)

// initLogger init logger
func initLogger(config *zap.Config) error {
	var logger *zap.Logger
	var err error
	if config == nil {
		logger, err = zap.NewDevelopment(zap.AddCaller(), zap.AddCallerSkip(1), zap.AddStacktrace(zap.WarnLevel))
	} else {
		logger, err = config.Build(zap.AddCaller(), zap.AddCallerSkip(1), zap.AddStacktrace(zap.WarnLevel))
	}
	if err != nil {
		zap.S().Fatalf("init logger failed, err=%v ", err)
		return err
	}
	zap.ReplaceGlobals(logger)
	return nil
}
