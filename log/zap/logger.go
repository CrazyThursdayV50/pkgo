package zap

import (
	"os"
	"time"

	"github.com/CrazyThursdayV50/pkgo/log"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Field = zap.Field

var _ log.DescLogger[Field] = (*zap.Logger)(nil)

var loggerLevelMap = map[string]zapcore.Level{
	"debug":  zapcore.DebugLevel,
	"info":   zapcore.InfoLevel,
	"warn":   zapcore.WarnLevel,
	"error":  zapcore.ErrorLevel,
	"dpanic": zapcore.DPanicLevel,
	"panic":  zapcore.PanicLevel,
	"fatal":  zapcore.FatalLevel,
}

func getLoggerLevel(loggerLevel string) zapcore.Level {
	level, exist := loggerLevelMap[loggerLevel]
	if !exist {
		return zapcore.DebugLevel
	}

	return level
}

func New(cfg *Config) *zap.Logger {
	logLevel := getLoggerLevel(cfg.Level)
	logWriter := zapcore.AddSync(os.Stdout)

	var encoderCfg zapcore.EncoderConfig
	var encoder zapcore.Encoder

	switch {
	case cfg.Development && cfg.Console:
		encoderCfg = zap.NewDevelopmentEncoderConfig()
		encoderCfg.EncodeLevel = zapcore.CapitalColorLevelEncoder
		encoder = zapcore.NewConsoleEncoder(encoderCfg)

	case cfg.Development && !cfg.Console:
		encoderCfg = zap.NewDevelopmentEncoderConfig()
		encoder = zapcore.NewJSONEncoder(encoderCfg)

	case !cfg.Development && cfg.Console:
		encoderCfg = zap.NewProductionEncoderConfig()
		encoderCfg.EncodeLevel = zapcore.LowercaseColorLevelEncoder
		encoderCfg.EncodeTime = zapcore.TimeEncoderOfLayout(time.StampMilli)
		encoder = zapcore.NewConsoleEncoder(encoderCfg)

	default:
		encoderCfg = zap.NewProductionEncoderConfig()
		encoderCfg.EncodeTime = zapcore.TimeEncoderOfLayout(time.StampMilli)
		encoder = zapcore.NewJSONEncoder(encoderCfg)
	}

	core := zapcore.NewCore(encoder, logWriter, zap.NewAtomicLevelAt(logLevel))
	logger := zap.New(core, zap.AddCaller(), zap.AddCallerSkip(cfg.CallerSkip))
	return logger
}
