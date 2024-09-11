package gorm

import (
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Config struct {
	Schema        string
	DSN           string
	MaxIdleConn   int
	MaxOpenConn   int
	MaxLifeTime   int64
	MaxIdleTime   int64
	ServerVersion string
	Gorm          gorm.Config
	Logger        logger.Config
}
