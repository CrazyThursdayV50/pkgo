package gorm

import (
	"context"
	"time"

	"github.com/CrazyThursdayV50/pkgo/log"
	gormlogger "github.com/CrazyThursdayV50/pkgo/log/gorm"
	"github.com/CrazyThursdayV50/pkgo/trace"
	sql "github.com/go-sql-driver/mysql"
	"gorm.io/driver/mysql"
	g "gorm.io/gorm"
)

type DB struct {
	*g.DB
}

func (db *DB) Db(ctx context.Context) *g.DB {
	return db.DB.WithContext(ctx)
}

func (db *DB) Tx(ctx context.Context, f func() bool) (*g.DB, func()) {
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

func NewDB(logger log.Logger, tracer trace.Tracer, cfg *Config) *DB {
	gl := gormlogger.New(logger, &cfg.Logger)
	cfg.Gorm.Logger = gl
	dsnConf, _ := sql.ParseDSN(cfg.DSN)
	dialector := mysql.New(mysql.Config{
		DSN:           cfg.DSN,
		DSNConfig:     dsnConf,
		ServerVersion: cfg.ServerVersion,
	})

	db, err := g.Open(dialector, &cfg.Gorm)
	if err != nil {
		logger.Fatalf("new gorm client failed: %v", err)
	}

	inner, err := db.DB()
	if err != nil {
		logger.Fatalf("get db failed: %v", err)
	}

	inner.SetMaxIdleConns(cfg.MaxIdleConn)
	inner.SetMaxOpenConns(cfg.MaxOpenConn)
	inner.SetConnMaxLifetime(time.Duration(cfg.MaxLifeTime) * time.Minute)
	inner.SetConnMaxIdleTime(time.Duration(cfg.MaxIdleTime) * time.Minute)

	registerInterceptors(db, traceInterceptor(tracer))
	return &DB{db}
}

func DefaultFindInBatchesCallback(tx *g.DB, batch int) error { return tx.Error }
