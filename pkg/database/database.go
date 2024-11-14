package database

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/souloss/go-clean-arch/pkg/constant"
	"github.com/souloss/go-clean-arch/pkg/logger"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var (
	instance *gorm.DB
	once     sync.Once
)

// GetDB returns the global database instance
func GetDB(ctx context.Context) *gorm.DB {
	v := ctx.Value(constant.GormDBCtxKey)
	if v != nil {
		if tx, ok := v.(*gorm.DB); ok {
			return tx
		}
	}
	return instance.WithContext(ctx)
}

// InitDB initializes the global database instance
func InitDB(ctx context.Context, cfg *DataSource, log logger.FormatStrLogger) (err error) {
	once.Do(func() {
		instance, err = newGormDB(ctx, cfg, log)
	})
	return err
}

func newGormDB(ctx context.Context, cfg *DataSource, log logger.FormatStrLogger) (*gorm.DB, error) {
	var (
		db  *gorm.DB
		err error
	)

	dbLogger := NewGormLogger(log)

	dsn := cfg.Raw
	if dsn == "" {
		dsn = cfg.GenerateDSN()
	}
	if cfg.Scheme == "" {
		cfg = Parse(cfg.Raw)
	}

	gConfig := &gorm.Config{
		// 在执行任何 SQL 时都会创建一个 prepared statement 并将其缓存，以提高后续的效率
		PrepareStmt: true,
		// 更改创建时间使用的函数
		NowFunc: func() time.Time {
			return time.Now().Local()
		},
		Logger: dbLogger,
	}

	switch cfg.Scheme {
	case "mysql":
		db, err = gorm.Open(mysql.Open(dsn), gConfig)
	case "postgres":
		db, err = gorm.Open(postgres.New(postgres.Config{
			DSN:                  dsn,
			PreferSimpleProtocol: true, // disables implicit prepared statement usage
		}), gConfig)
	case "sqlite":
		db, err = gorm.Open(sqlite.Open(dsn), gConfig)
	default:
		return nil, fmt.Errorf("unknown database dialect: %s", cfg.Scheme)
	}

	if err != nil {
		logger.L().Error("new gorm db error: %v", err, ctx)
		return nil, err
	}

	// Connection Pool config
	sqlDB, err := db.DB()
	if err != nil {
		panic(err)
	}
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	return db, nil
}
