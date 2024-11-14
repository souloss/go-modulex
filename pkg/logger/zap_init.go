package logger

import (
	"context"
	"io"
	"os"
	"sync"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var doOnce sync.Once

func initLogger(logLevel string, replaceGlobal bool, writers ...io.Writer) (*zap.Logger, error) {
	level := zapcore.InfoLevel
	if err := level.UnmarshalText([]byte(logLevel)); err != nil {
		return nil, err
	}
	var zapWriters []zapcore.WriteSyncer
	for _, writer := range writers {
		zapWriters = append(zapWriters, zapcore.AddSync(writer))
	}
	encoder := zapcore.EncoderConfig{
		TimeKey:      "ts",
		LevelKey:     "level",
		MessageKey:   "msg",
		NameKey:      "name",
		CallerKey:    "caller",
		EncodeLevel:  zapcore.LowercaseColorLevelEncoder,
		EncodeTime:   zapcore.RFC3339TimeEncoder,
		EncodeName:   zapcore.FullNameEncoder,
		EncodeCaller: zapcore.ShortCallerEncoder,
	}
	core := zapcore.NewCore(
		zapcore.NewConsoleEncoder(encoder),
		// zapcore.NewJSONEncoder(encoder),
		zapcore.NewMultiWriteSyncer(zapWriters...),
		level,
	)
	l := zap.New(core).WithOptions(
		zap.AddCaller(),
		// 使用了适配器，所以需要跳过一层
		zap.AddCallerSkip(1),
	)
	if replaceGlobal {
		zap.ReplaceGlobals(l)
	}
	return l, nil
}

func getDefaultWrites(logConfig *LoggerConfig) ([]io.Writer, error) {
	writes := []io.Writer{}
	rotateWrite, err := NewRotateWriter(RotateConfig{
		Filename:     logConfig.FileName,
		MaxAge:       time.Duration(logConfig.MaxAge),
		RotationTime: time.Duration(logConfig.RotationTime),
		MaxSize:      int64(logConfig.MaxSize),
		MaxBackups:   logConfig.MaxBackups,
		Compress:     logConfig.Compress,
		LocalTime:    logConfig.LocalTime,
	})
	if err != nil {
		// 使用全局 Logger 打印警告日志
		L().Warnf("Rotate write create failed, err: %v", err)
		// 创建轮转日志写入器失败，创建普通文件作为写入器
		if logConfig.FileName != "" {
			file, err := os.OpenFile(logConfig.FileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
			if err != nil {
				L().Warnf("file write create failed, err: %v", err)
				return nil, err
			}
			writes = append(writes, file)
		}
	} else {
		writes = append(writes, rotateWrite)
	}
	if logConfig.Console {
		writes = append(writes, os.Stdout)
	}
	return writes, nil
}

func InitZapGlobalLogger(ctx context.Context, logConfig *LoggerConfig) error {
	var err error
	writes, err := getDefaultWrites(logConfig)
	if err != nil {
		L().Error("get logger writes error: %v", err, ctx)
		return err
	}
	doOnce.Do(func() {
		_, err = initLogger(logConfig.Level, true, writes...)
		if err != nil {
			L().Error("new logger error: %v", err, ctx)
		}
	})
	return err
}

func MustNewZapLogger(ctx context.Context, logConfig *LoggerConfig) *zap.Logger {
	writes, err := getDefaultWrites(logConfig)
	if err != nil {
		L().Error("get logger writes error: %v", err, ctx)
		panic(err)
	}
	l, err := initLogger(logConfig.Level, false, writes...)
	if err != nil {
		L().Error("new logger error: %v", err, ctx)
		panic(err)
	}
	return l
}
