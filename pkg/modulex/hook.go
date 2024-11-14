package modulex

import (
	"context"
	"sync"
)

// Hook 定义生命周期钩子函数
type Hook func(context.Context) error

// LifecyclePhase defines the phases of the application lifecycle.
type LifecyclePhase string

// Constants for lifecycle phases using iota for simplicity and ease of maintenance.
const (
	PhaseBeforeInit  LifecyclePhase = "before_init"
	PhaseAfterInit   LifecyclePhase = "after_init"
	PhaseBeforeStart LifecyclePhase = "before_start"
	PhaseAfterStart  LifecyclePhase = "after_start"
	PhaseBeforeStop  LifecyclePhase = "before_stop"
	PhaseAfterStop   LifecyclePhase = "after_stop"
)

type Hookable interface {
	AddHook(phase LifecyclePhase, hook Hook) error
	RunHooks(phase LifecyclePhase, ctx context.Context) error
}

// BootableModule 带 hook 的可启动的模块
type HookedBootableModule struct {
	BootableModule
	mu    sync.Mutex
	hooks map[LifecyclePhase][]Hook
}

func NewHookedBootableModule(m BootableModule) *HookedBootableModule {
	return &HookedBootableModule{
		BootableModule: m,
		mu:             sync.Mutex{},
		hooks:          make(map[LifecyclePhase][]Hook),
	}
}

func (h *HookedBootableModule) Init(ctx context.Context, c any) error {
	if err := h.RunHooks(PhaseBeforeInit, ctx); err != nil {
		return err
	}
	if err := h.BootableModule.Init(ctx, c); err != nil {
		return err
	}
	return h.RunHooks(PhaseAfterInit, ctx)
}

func (h *HookedBootableModule) Start(ctx context.Context) error {
	if err := h.RunHooks(PhaseBeforeStart, ctx); err != nil {
		return err
	}
	if err := h.BootableModule.Start(ctx); err != nil {
		return err
	}
	return h.RunHooks(PhaseAfterStart, ctx)
}

func (h *HookedBootableModule) Stop(ctx context.Context) error {
	if err := h.RunHooks(PhaseBeforeStop, ctx); err != nil {
		return err
	}
	if err := h.BootableModule.Stop(ctx); err != nil {
		return err
	}
	return h.RunHooks(PhaseAfterStop, ctx)
}

func (h *HookedBootableModule) AddHook(phase LifecyclePhase, hook Hook) error {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.hooks[phase] = append(h.hooks[phase], hook)
	return nil
}

func (h *HookedBootableModule) RunHooks(phase LifecyclePhase, ctx context.Context) error {
	for _, hook := range h.hooks[phase] {
		if err := hook(ctx); err != nil {
			return err
		}
	}
	return nil
}
