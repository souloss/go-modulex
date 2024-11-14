package logger

import (
	"io"
	"time"

	"github.com/natefinch/lumberjack/v3"
)

type RotateConfig struct {
	// 共用配置
	Filename string        // 完整文件名
	MaxAge   time.Duration // 保留旧日志文件的最大天数

	// 按时间轮转配置
	RotationTime time.Duration // 日志文件轮转时间

	// 按大小轮转配置
	MaxSize    int64 // 日志文件最大大小（MB）
	MaxBackups int   // 保留日志文件的最大数量
	Compress   bool  // 是否对日志文件进行压缩归档
	LocalTime  bool  // 是否使用本地时间，默认 UTC 时间
}

func NewRotateWriter(c RotateConfig) (io.Writer, error) {
	return lumberjack.NewRoller(c.Filename, c.MaxSize, &lumberjack.Options{
		MaxAge:     c.MaxAge,
		MaxBackups: c.MaxBackups,
		LocalTime:  c.LocalTime,
		Compress:   c.Compress,
	})
}
