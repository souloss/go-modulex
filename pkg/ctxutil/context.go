package ctxutil

import (
	"context"

	"github.com/souloss/go-clean-arch/pkg/constant"
	"github.com/souloss/go-clean-arch/pkg/logger"
)

// WithLoggerContext 获取带日志实例的上下文
//
//	@param ctx
//	@param logger
//	@return context.Context
func WithLoggerContext(ctx context.Context, logger logger.FormatStrLogger) context.Context {
	return context.WithValue(ctx, constant.LoggerCtxKey, logger)
}

// GetLoggerFromContext 从上下文中获取日志实例
//
//	@param ctx
//	@return logger.FormatStrLogger
func GetLoggerFromContext(ctx context.Context) logger.FormatStrLogger {
	log, ok := ctx.Value(constant.LoggerCtxKey).(logger.FormatStrLogger)
	if !ok {
		return logger.L()
	}
	return log
}

// WithTraceIDContext 获取带TraceID的上下文
//
//	@param ctx
//	@param id
//	@return context.Context
func WithTraceIDContext(ctx context.Context, traceID string) context.Context {
	return context.WithValue(ctx, constant.Traceidctxkey, traceID)
}

// GetTraceIDFromContext 从上下文中获取TraceID
//
//	@param ctx
//	@return string
//	@return bool
func GetTraceIDFromContext(ctx context.Context) (string, bool) {
	traceID, ok := ctx.Value(constant.Traceidctxkey).(string)
	return traceID, ok
}
