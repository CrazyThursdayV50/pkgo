package gormlogger

import (
	"time"

	"gorm.io/gorm/logger"
)

func New(l logger.Writer, cfg *logger.Config) logger.Interface {
	logCfg := logger.Config{
		SlowThreshold:             cfg.SlowThreshold * time.Millisecond,
		Colorful:                  cfg.Colorful,
		IgnoreRecordNotFoundError: cfg.IgnoreRecordNotFoundError,
		ParameterizedQueries:      cfg.ParameterizedQueries,
		LogLevel:                  cfg.LogLevel,
	}

	return logger.New(l, logCfg)
}
