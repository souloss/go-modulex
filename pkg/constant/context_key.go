package constant

// contextKey is a custom type for context keys to avoid collisions
type contextKey string

// Define context keys
const (
	// Session and Trace related keys
	Traceidctxkey contextKey = "trace_id"

	// Framework instances
	LoggerCtxKey contextKey = "logger"
	GormDBCtxKey contextKey = "gorm_db"

	// Business context
	LanguageCtxKey contextKey = "language"
)
