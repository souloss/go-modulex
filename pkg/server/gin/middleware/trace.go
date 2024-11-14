package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/rs/xid"
	"github.com/souloss/go-clean-arch/pkg/ctxutil"
)

type Option func(*xRquestTraceConfig)
type HeaderStrKey string

type (
	Generator func() string
)

// Config defines the xRquestTraceConfig for RequestID middleware
type xRquestTraceConfig struct {
	// Generator defines a function to generate an ID.
	// Optional. Default: func() string {
	//   return xid.New().String()
	// }
	generator Generator
	headerKey HeaderStrKey
}

var headerXRequestIDKey string

// WithGenerator set generator function
func WithGenerator(g Generator) Option {
	return func(cfg *xRquestTraceConfig) {
		cfg.generator = g
	}
}

// WithCustomHeaderStrKey set custom header key for request id
func WithCustomHeaderStrKey(s HeaderStrKey) Option {
	return func(cfg *xRquestTraceConfig) {
		cfg.headerKey = s
	}
}

// 链路追踪中间件
func NewTraceMiddleware(opts ...Option) gin.HandlerFunc {
	cfg := &xRquestTraceConfig{
		generator: func() string {
			return xid.New().String()
		},
		headerKey: "X-Request-ID",
	}
	for _, opt := range opts {
		opt(cfg)
	}
	headerXRequestIDKey = string(cfg.headerKey)

	return func(c *gin.Context) {
		// Get id from request
		rid := c.GetHeader(headerXRequestIDKey)
		if rid == "" {
			rid = cfg.generator()
			c.Request.Header.Add(headerXRequestIDKey, rid)
		}
		// 将 traceID 存入 ctx
		c.Request = c.Request.WithContext(ctxutil.WithTraceIDContext(c.Request.Context(), rid))
		// Set the id to ensure that the requestid is in the response
		c.Header(headerXRequestIDKey, rid)
		c.Next()
	}
}

// Get returns the request identifier
func GetTraceID(c *gin.Context) string {
	return c.GetHeader(headerXRequestIDKey)
}
