package logger

import (
	"os"

	"github.com/DTreshy/nubl9-recrutation/internal/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func New(log config.LogConfig) (*zap.Logger, error) {
	fileConfig := zap.NewProductionConfig()
	fileConfig.EncoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05")
	consoleConfig := fileConfig
	consoleConfig.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	fileConfig.EncoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	fileEncoder := zapcore.NewJSONEncoder(fileConfig.EncoderConfig)
	consoleEncoder := zapcore.NewConsoleEncoder(consoleConfig.EncoderConfig)
	logFile, err := os.OpenFile(log.File, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	writer := zapcore.AddSync(logFile)
	level := zapcore.DebugLevel
	invalidLevel := false

	switch log.Level {
	case "info":
		level = zapcore.InfoLevel
	case "warn":
		level = zapcore.WarnLevel
	case "debug":
		level = zapcore.DebugLevel
	case "error":
		level = zapcore.ErrorLevel
	default:
		invalidLevel = true
	}

	core := zapcore.NewTee(
		zapcore.NewCore(fileEncoder, writer, level),
		zapcore.NewCore(consoleEncoder, zapcore.AddSync(os.Stdout), level),
	)
	logger := zap.New(core, zap.AddStacktrace(zapcore.ErrorLevel))

	if invalidLevel {
		logger.Sugar().Warnf("Provided invalid log level: %s", log.Level)
	}

	return logger, err
}
