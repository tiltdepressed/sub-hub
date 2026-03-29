package logging

import (
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"sub-hub/internal/config"
)

func New(cfg config.LogConfig, env, version, commit string) (*zap.Logger, error) {
	level := zapcore.InfoLevel
	_ = level.Set(strings.ToLower(cfg.Level))

	encoderCfg := zap.NewProductionEncoderConfig()
	encoderCfg.TimeKey = "ts"
	encoderCfg.EncodeTime = zapcore.ISO8601TimeEncoder

	loggerCfg := zap.NewProductionConfig()
	loggerCfg.Level = zap.NewAtomicLevelAt(level)
	loggerCfg.EncoderConfig = encoderCfg
	loggerCfg.OutputPaths = []string{"stdout"}
	loggerCfg.ErrorOutputPaths = []string{"stderr"}
	logger, err := loggerCfg.Build(
		zap.Fields(
			zap.String("env", env),
			zap.String("version", version),
			zap.String("commit", commit),
		),
	)
	if err != nil {
		return nil, err
	}
	return logger, nil
}
