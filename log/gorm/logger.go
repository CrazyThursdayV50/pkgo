package gormlogger

import (
	"time"

	"github.com/CrazyThursdayV50/pkgo/log"
	"gorm.io/gorm/logger"
)

func New(l log.Logger, cfg *logger.Config) logger.Interface {
	logCfg := logger.Config{
		SlowThreshold:             cfg.SlowThreshold * time.Millisecond,
		Colorful:                  cfg.Colorful,
		IgnoreRecordNotFoundError: cfg.IgnoreRecordNotFoundError,
		ParameterizedQueries:      cfg.ParameterizedQueries,
		LogLevel:                  cfg.LogLevel,
	}

	return logger.New(l, logCfg)
}
