package modulex

import (
	"context"
	"testing"
)

func TestApplication(t *testing.T) {
	app := NewBaseApplication()
	app.Init(context.Background())
	app.Start(context.Background())
}
