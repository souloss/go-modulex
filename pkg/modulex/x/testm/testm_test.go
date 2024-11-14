package testm

import (
	"context"
	"fmt"
	"testing"

	"github.com/souloss/go-clean-arch/pkg/modulex"
)

type TestModule struct {
	*modulex.BaseBootableModule
}

func New() *TestModule {
	return &TestModule{
		BaseBootableModule: modulex.NewBaseBootableModule(),
	}
}

func (t *TestModule) init(context.Context, any) error {
	fmt.Println("testmodule init")
	return nil
}

func TestXXXX(t *testing.T) {
	tm := New()
	tm.Init(context.Background(), "")
	// tm.init(context.Background(), "")
}
