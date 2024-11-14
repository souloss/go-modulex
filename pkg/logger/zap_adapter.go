package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type ZapAdapter struct {
	logger *zap.Logger
}

func NewZapAdapter(logger *zap.Logger) *ZapAdapter {
	return &ZapAdapter{logger: logger}
}

// 确保实现了接口
var _ FormatStrLogger = (*ZapAdapter)(nil)

func (l *ZapAdapter) Named(name string) FormatStrLogger {
	return NewZapAdapter(l.logger.Named(name))
}

func (l *ZapAdapter) With(key string, val interface{}) FormatStrLogger {
	return &ZapAdapter{
		logger: l.logger.With(zap.Any(key, val)),
	}
}

func (l *ZapAdapter) GetLevel() Level {
	switch l.logger.Level() {
	case zapcore.DebugLevel:
		return DebugLevel
	case zapcore.ErrorLevel:
		return ErrorLevel
	case zapcore.InfoLevel:
		return InfoLevel
	case zapcore.WarnLevel:
		return WarnLevel
	default:
		return DebugLevel
	}
}

func (l *ZapAdapter) Debug(args ...interface{}) {
	l.logger.Sugar().Debug(args...)
}

func (l *ZapAdapter) Debugf(format string, args ...interface{}) {
	l.logger.Sugar().Debugf(format, args...)
}

func (l *ZapAdapter) Info(args ...interface{}) {
	l.logger.Sugar().Info(args...)
}

func (l *ZapAdapter) Infof(format string, args ...interface{}) {
	l.logger.Sugar().Infof(format, args...)
}

func (l *ZapAdapter) Warn(args ...interface{}) {
	l.logger.Sugar().Warn(args...)
}

func (l *ZapAdapter) Warnf(format string, args ...interface{}) {
	l.logger.Sugar().Warnf(format, args...)
}

func (l *ZapAdapter) Error(args ...interface{}) {
	l.logger.Sugar().Error(args...)
}

func (l *ZapAdapter) Errorf(format string, args ...interface{}) {
	l.logger.Sugar().Errorf(format, args...)
}

func (l *ZapAdapter) Printf(format string, args ...interface{}) {
	l.logger.Sugar().Infof(format, args...)
}
