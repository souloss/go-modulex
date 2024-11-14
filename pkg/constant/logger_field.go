package constant

// LoggerKey is a custom type for context keys to avoid collisions
type LoggerKey string

// Define Logger field keys
const (
	// Session and Trace related keys
	TraceIDLoggerKey LoggerKey = "trace_id"
)
