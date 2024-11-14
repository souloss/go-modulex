package logger

import (
	"context"
	"sync"
)

type loggerKey string

const logKey loggerKey = "LOGGER"

// Level 定义日志级别
type Level int

const (
	DebugLevel Level = iota
	InfoLevel
	WarnLevel
	ErrorLevel

	_minLevel = DebugLevel
	_maxLevel = ErrorLevel

	InvalidLevel = _maxLevel + 1
)

// FormatStrLogger 定义日志接口
type FormatStrLogger interface {
	// 命名空间
	Named(string) FormatStrLogger
	// 日志上下文
	With(string, interface{}) FormatStrLogger
	GetLevel() Level
	Debug(args ...interface{})
	Debugf(format string, args ...interface{})
	Info(args ...interface{})
	Infof(format string, args ...interface{})
	Warn(args ...interface{})
	Warnf(format string, args ...interface{})
	Error(args ...interface{})
	Errorf(format string, args ...interface{})
	Printf(string, ...interface{})
}

// JSONLogger 定义JSON格式日志接口
type JSONLogger interface {
	Named(string) JSONLogger
	GetLevel() Level
	Debug(msg string, fields map[string]interface{})
	Info(msg string, fields map[string]interface{})
	Warn(msg string, fields map[string]interface{})
	Error(msg string, fields map[string]interface{})
}

var (
	_globalMu sync.RWMutex
	_globalL  FormatStrLogger = NewSimpleLogger("", InfoLevel)
)

func L() FormatStrLogger {
	_globalMu.RLock()
	l := _globalL
	_globalMu.RUnlock()
	return l
}

func ReplaceGlobals(logger FormatStrLogger) func() {
	_globalMu.Lock()
	prev := _globalL
	_globalL = logger
	_globalMu.Unlock()
	return func() { ReplaceGlobals(prev) }
}

// WithLoggerContext 获取带日志实例的上下文
//
//	@param ctx
//	@param logger
//	@return context.Context
func WithLoggerContext(ctx context.Context, logger FormatStrLogger) context.Context {
	return context.WithValue(ctx, logKey, logger)
}

// FromContext 从上下文中获取日志实例
//
//	@param ctx
//	@return FormatStrLogger
func FromContext(ctx context.Context) FormatStrLogger {
	logger, ok := ctx.Value(logKey).(FormatStrLogger)
	if !ok {
		return L()
	}
	return logger
}

// getLevelStr 获取日志级别字符串
//
//	@param l
//	@return string
func getLevelStr(l Level) string {
	switch l {
	case DebugLevel:
		return "DEBUG"
	case InfoLevel:
		return "INFO"
	case WarnLevel:
		return "WARN"
	case ErrorLevel:
		return "ERROR"
	default:
		return "-"
	}
}
