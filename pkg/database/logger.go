package database

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/souloss/go-clean-arch/pkg/logger"

	gormlogger "gorm.io/gorm/logger"
	"gorm.io/gorm/utils"
)

const ctxLoggerKey = "Logger"

type GormLogger struct {
	gormlogger.Config
	Logger logger.FormatStrLogger
}

func NewGormLogger(log logger.FormatStrLogger) gormlogger.Interface {
	return &GormLogger{
		Config: gormlogger.Config{
			LogLevel:                  gormlogger.Warn,
			SlowThreshold:             100 * time.Millisecond,
			Colorful:                  false,
			IgnoreRecordNotFoundError: false,
			ParameterizedQueries:      false,
		},
		Logger: log,
	}
}

func (l *GormLogger) LogMode(level gormlogger.LogLevel) gormlogger.Interface {
	newlogger := *l
	newlogger.LogLevel = level
	return &newlogger
}

// Info print info
func (l GormLogger) Info(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= gormlogger.Info {
		l.logger(ctx).Infof(msg, data...)
	}
}

// Warn print warn messages
func (l GormLogger) Warn(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= gormlogger.Warn {
		l.logger(ctx).Warnf(msg, data...)
	}
}

// Error print error messages
func (l GormLogger) Error(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= gormlogger.Error {
		l.logger(ctx).Errorf(msg, data...)
	}
}

func (l GormLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	if l.LogLevel <= gormlogger.Silent {
		return
	}

	elapsed := time.Since(begin)
	elapsedStr := fmt.Sprintf("%.3fms", float64(elapsed.Nanoseconds())/1e6)
	logger := l.logger(ctx)
	switch {
	case err != nil && l.LogLevel >= gormlogger.Error && (!errors.Is(err, gormlogger.ErrRecordNotFound) || !l.IgnoreRecordNotFoundError):
		sql, rows := fc()
		if rows == -1 {
			logger.Errorf("trace %s err: %s, elapsed: %d, rows: %d, sql: %s", utils.FileWithLineNum(), err, elapsedStr, rows, sql)
		} else {
			logger.Errorf("trace %s err: %s, elapsed: %d, rows: %s, sql: %s", utils.FileWithLineNum(), err, elapsedStr, "-", sql)
		}
	case elapsed > l.SlowThreshold && l.SlowThreshold != 0 && l.LogLevel >= gormlogger.Warn:
		sql, rows := fc()
		slowLog := fmt.Sprintf("SLOW SQL >= %v", l.SlowThreshold)
		if rows == -1 {
			logger.Errorf("trace %s slow: %s, elapsed: %d, rows: %d, sql: %s", utils.FileWithLineNum(), slowLog, elapsedStr, rows, sql)
		} else {
			logger.Errorf("trace %s slow: %s, elapsed: %d, rows: %s, sql: %s", utils.FileWithLineNum(), slowLog, elapsedStr, "-", sql)
		}
	case l.LogLevel == gormlogger.Info:
		sql, rows := fc()
		if rows == -1 {
			logger.Errorf("trace %s elapsed: %d, rows: %d, sql: %s", utils.FileWithLineNum(), elapsedStr, rows, sql)
		} else {
			logger.Errorf("trace %s elapsed: %d, rows: %s, sql: %s", utils.FileWithLineNum(), elapsedStr, "-", sql)
		}
	}
}

func (l GormLogger) logger(ctx context.Context) logger.FormatStrLogger {
	log := l.Logger
	if ctx != nil {
		if c, ok := ctx.(*gin.Context); ok {
			ctx = c.Request.Context()
		}
		l := ctx.Value(ctxLoggerKey)
		ctxLogger, ok := l.(logger.FormatStrLogger)
		if ok {
			log = ctxLogger
		}
	}
	return log
}
