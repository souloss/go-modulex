package modulex

import (
	"context"
)

// Module 简单模块
type Module interface {
	// Name returns the name of the module
	Name() string
	// Init initializes the module with given context
	Init(context.Context, any) error
}

// BootableModule 可启动的模块
type BootableModule interface {
	// Module
	Module

	// Start starts the module
	Start(context.Context) error
	// Stop stops the module
	Stop(context.Context) error
}

type BaseBootableModule struct {
}

func NewBaseBootableModule() *BaseBootableModule {
	return &BaseBootableModule{}
}

func (b *BaseBootableModule) Name() string {
	return "BaseBootableModule"
}

func (b *BaseBootableModule) Init(ctx context.Context, c any) error {
	return nil
}

func (b *BaseBootableModule) init(context.Context, any) error {
	return nil
}

func (b *BaseBootableModule) Start(context.Context) error {
	return nil
}
func (b *BaseBootableModule) Stop(context.Context) error {
	return nil
}
