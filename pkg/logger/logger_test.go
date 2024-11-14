package logger

import "testing"

func TestXxx(t *testing.T) {
	logger := L().With("traceID", "avcd")
	logger.Warnf("hahah %s", "Sdfdsf")
}
