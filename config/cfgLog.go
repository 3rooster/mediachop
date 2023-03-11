package config

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

// initLogger init logger
func initLogger(config *zap.Config, rotationConfig *lumberjack.Logger) error {
	var logger *zap.Logger
	var err error

	if config == nil {
		logger, err = zap.NewDevelopment(zap.AddCaller(), zap.AddCallerSkip(1), zap.AddStacktrace(zap.WarnLevel))
	} else if config != nil && rotationConfig != nil {
		w := zapcore.AddSync(rotationConfig)
		core := zapcore.NewCore(
			zapcore.NewJSONEncoder(config.EncoderConfig),
			w,
			config.Level,
		)
		logger = zap.New(core,
			zap.ErrorOutput(w),
			zap.AddCaller(),
			zap.AddCallerSkip(1),
			zap.AddStacktrace(zap.WarnLevel))
	}
	if err != nil {
		zap.S().Fatalf("init logger failed, err=%v ", err)
		return err
	}
	zap.ReplaceGlobals(logger)
	return nil
}
