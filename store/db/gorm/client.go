package gorm

import (
	"context"
	"log"
	"time"

	sql "github.com/go-sql-driver/mysql"
	"github.com/opentracing/opentracing-go"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type DB struct {
	*gorm.DB
}

func (db *DB) Db(ctx context.Context) *gorm.DB {
	return db.DB.WithContext(ctx)
}

func (db *DB) Tx(ctx context.Context, f func() bool) (*gorm.DB, func()) {
	tx := db.DB.Begin().WithContext(ctx)
	fn := func() {
		if f() {
			tx.Commit()
			return
		}

		tx.Rollback()
	}
	return tx, fn
}

func NewDB(logger logger.Interface, tracer opentracing.Tracer, cfg *Config) *DB {
	cfg.Gorm.Logger = logger
	dsnConf, _ := sql.ParseDSN(cfg.DSN)
	dialector := mysql.New(mysql.Config{
		DSN:           cfg.DSN,
		DSNConfig:     dsnConf,
		ServerVersion: cfg.ServerVersion,
	})

	db, err := gorm.Open(dialector, &cfg.Gorm)
	if err != nil {
		log.Fatalf("new gorm client failed: %v", err)
	}

	inner, err := db.DB()
	if err != nil {
		log.Fatalf("get db failed: %v", err)
	}

	inner.SetMaxIdleConns(cfg.MaxIdleConn)
	inner.SetMaxOpenConns(cfg.MaxOpenConn)
	inner.SetConnMaxLifetime(time.Duration(cfg.MaxLifeTime) * time.Minute)
	inner.SetConnMaxIdleTime(time.Duration(cfg.MaxIdleTime) * time.Minute)

	registerInterceptors(db, traceInterceptor(tracer))
	return &DB{db}
}

func DefaultFindInBatchesCallback(tx *gorm.DB, batch int) error { return tx.Error }
