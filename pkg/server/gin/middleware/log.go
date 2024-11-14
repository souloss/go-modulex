package middleware

import (
	"bytes"
	"io"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/souloss/go-clean-arch/pkg/constant"
	"github.com/souloss/go-clean-arch/pkg/ctxutil"
)

func RequestLogMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var bodyBytes []byte
		if ctx.Request.Body != nil {
			bodyBytes, _ = ctx.GetRawData()
			ctx.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
		}

		traceID, _ := ctxutil.GetTraceIDFromContext(ctx.Request.Context())

		// 为日志实例和Context添加RequestID字段，并存入上下文
		sessionLogger := ctxutil.GetLoggerFromContext(ctx.Request.Context())
		sessionLogger.Infof("Trace: %s, Request Method: %s, Header: %v, URL: %s, Params: %v",
			traceID,
			ctx.Request.Method,
			ctx.Request.Header,
			ctx.Request.URL.String(),
			string(bodyBytes))

		sessionLogger = sessionLogger.With(string(constant.TraceIDLoggerKey), traceID)
		ctx.Request = ctx.Request.WithContext(ctxutil.WithLoggerContext(ctx.Request.Context(), sessionLogger))

		ctx.Next()
	}
}
func ResponseLogMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		blw := &bodyLogWriter{body: bytes.NewBufferString(""), ResponseWriter: ctx.Writer}
		ctx.Writer = blw
		startTime := time.Now()
		ctx.Next()
		duration := time.Since(startTime).String()
		ctxutil.GetLoggerFromContext(ctx.Request.Context()).Infof("Response body: %s, time: %v", blw.body.String(), duration)
	}
}

type bodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w bodyLogWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}
