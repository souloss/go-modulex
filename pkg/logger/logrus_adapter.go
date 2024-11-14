package logger

import (
	"context"

	"github.com/sirupsen/logrus"
)

type LogrusAdapter struct {
	name   string
	logger *logrus.Entry
}

func NewLogrusAdapter(logger *logrus.Logger) *LogrusAdapter {
	return &LogrusAdapter{logger: logger.WithContext(context.TODO())}
}

func (l *LogrusAdapter) Named(name string) FormatStrLogger {
	newName := l.name
	if newName != "" {
		newName = l.name + "." + name
	} else {
		newName = name
	}
	return &LogrusAdapter{
		name: newName,
		logger: l.logger.WithFields(logrus.Fields{
			"logger": newName,
		}),
	}
}

func (l *LogrusAdapter) With(key string, val interface{}) FormatStrLogger {
	return &LogrusAdapter{
		name:   l.name,
		logger: l.logger.WithField(key, val),
	}
}

func (l *LogrusAdapter) GetLevel() Level {
	switch l.logger.Level {
	case logrus.DebugLevel:
		return DebugLevel
	case logrus.ErrorLevel:
		return ErrorLevel
	case logrus.InfoLevel:
		return InfoLevel
	case logrus.WarnLevel:
		return WarnLevel
	default:
		return DebugLevel
	}
}

func (l *LogrusAdapter) Debug(args ...interface{}) {
	l.logger.Debug(args...)
}

func (l *LogrusAdapter) Debugf(format string, args ...interface{}) {
	l.logger.Debugf(format, args...)
}

func (l *LogrusAdapter) Info(args ...interface{}) {
	l.logger.Info(args...)
}

func (l *LogrusAdapter) Infof(format string, args ...interface{}) {
	l.logger.Infof(format, args...)
}

func (l *LogrusAdapter) Warn(args ...interface{}) {
	l.logger.Warn(args...)
}

func (l *LogrusAdapter) Warnf(format string, args ...interface{}) {
	l.logger.Warnf(format, args...)
}

func (l *LogrusAdapter) Error(args ...interface{}) {
	l.logger.Error(args...)
}

func (l *LogrusAdapter) Errorf(format string, args ...interface{}) {
	l.logger.Errorf(format, args...)
}

func (l *LogrusAdapter) Printf(format string, args ...interface{}) {
	l.logger.Printf(format, args...)
}
