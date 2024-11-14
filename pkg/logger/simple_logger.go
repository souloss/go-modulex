package logger

import (
	"fmt"
	"log"
	"path/filepath"
	"runtime"
	"time"
)

// SimpleLogger 是一个简单的日志记录器实现
type SimpleLogger struct {
	name   string
	level  Level
	values map[string]interface{}
}

// 确保实现了接口
var _ FormatStrLogger = (*SimpleLogger)(nil)

// NewSimpleLogger 创建一个新的 SimpleLogger 实例
func NewSimpleLogger(name string, level Level) *SimpleLogger {
	log.SetFlags(0)
	return &SimpleLogger{
		name:   name,
		level:  level,
		values: make(map[string]interface{}),
	}
}

// Named 设置日志记录器的名字并返回自身
func (l *SimpleLogger) Named(name string) FormatStrLogger {
	newName := l.name
	if newName != "" {
		newName = l.name + "." + name
	} else {
		newName = name
	}
	return &SimpleLogger{
		name:  newName,
		level: l.level,
	}
}

func (l *SimpleLogger) With(key string, val interface{}) FormatStrLogger {
	values := l.values
	values[key] = val
	return &SimpleLogger{
		name:   l.name,
		level:  l.level,
		values: values,
	}
}

// GetLevel 返回当前日志记录器的级别
func (l *SimpleLogger) GetLevel() Level {
	return l.level
}

// Debug 记录一条调试信息
func (l *SimpleLogger) Debug(args ...interface{}) {
	if l.level <= DebugLevel {
		l.log(DebugLevel, fmt.Sprint(args...))
	}
}

// Debugf 记录一条格式化的调试信息
func (l *SimpleLogger) Debugf(format string, args ...interface{}) {
	if l.level <= DebugLevel {
		l.log(DebugLevel, fmt.Sprintf(format, args...))
	}
}

// Info 记录一条信息
func (l *SimpleLogger) Info(args ...interface{}) {
	if l.level <= InfoLevel {
		l.log(InfoLevel, fmt.Sprint(args...))
	}
}

// Infof 记录一条格式化的信息
func (l *SimpleLogger) Infof(format string, args ...interface{}) {
	if l.level <= InfoLevel {
		l.log(InfoLevel, fmt.Sprintf(format, args...))
	}
}

// Warn 记录一条警告信息
func (l *SimpleLogger) Warn(args ...interface{}) {
	if l.level <= WarnLevel {
		l.log(WarnLevel, fmt.Sprint(args...))
	}
}

// Warnf 记录一条格式化的警告信息
func (l *SimpleLogger) Warnf(format string, args ...interface{}) {
	if l.level <= WarnLevel {
		l.log(WarnLevel, fmt.Sprintf(format, args...))
	}
}

// Error 记录一条错误信息
func (l *SimpleLogger) Error(args ...interface{}) {
	if l.level <= ErrorLevel {
		l.log(ErrorLevel, fmt.Sprint(args...))
	}
}

// Errorf 记录一条格式化的错误信息
func (l *SimpleLogger) Errorf(format string, args ...interface{}) {
	if l.level <= ErrorLevel {
		l.log(ErrorLevel, fmt.Sprintf(format, args...))
	}
}

// Printf 打印一条无级别日志
func (l *SimpleLogger) Printf(format string, args ...interface{}) {
	l.log(_maxLevel, fmt.Sprintf(format, args...))
}

// log
func (l *SimpleLogger) log(level Level, message string) {
	levelStr := getLevelStr(level)
	now := time.Now().Format("2006-01-02 15:04:05")
	_, file, line, ok := runtime.Caller(2) // 跳过 log 和 public 函数
	if !ok {
		file = "???"
		line = 0
	} else {
		file = filepath.Base(file)
	}
	valuesMes := ""
	for key, val := range l.values {
		valuesMes += fmt.Sprintf("%s: %v", key, val)
	}
	logMessage := fmt.Sprintf("%s [%s] %s:%d %s %s\n", now, levelStr, file, line, message, valuesMes)

	log.Println(logMessage)
}
