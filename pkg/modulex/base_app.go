package modulex

import (
	"context"
	"fmt"
	"sync"

	"github.com/souloss/go-clean-arch/pkg/config"
	"github.com/souloss/go-clean-arch/pkg/logger"
	"github.com/souloss/go-clean-arch/pkg/version"
)

// BaseApplication 提供 Application 接口的基础实现
type BaseApplication struct {
	mu      sync.Mutex
	modules map[string]Module

	name          string
	hooks         map[LifecyclePhase][]Hook
	ConfigManager *config.ConfigManager
}

func NewBaseApplication() *BaseApplication {
	nopCfg, _ := config.New(&struct{}{})
	return &BaseApplication{
		name:          "ModuleX Application",
		modules:       make(map[string]Module),
		hooks:         make(map[LifecyclePhase][]Hook),
		ConfigManager: nopCfg,
	}
}

func (app *BaseApplication) Banner() string {
	banner := `
    __  ___              __            __         _  __
   /  |/  /  ____   ____/ /  __  __   / /  ___   | |/ /
  / /|_/ /  / __ \ / __  /  / / / /  / /  / _ \  |   / 
 / /  / /  / /_/ // /_/ /  / /_/ /  / /  /  __/ /   |  
/_/  /_/   \____/ \__,_/   \__,_/  /_/   \___/ /_/|_|  
                                                       
 :: ModuleX ::               (%s)
	`
	return fmt.Sprintf(banner, version.Version)
}

func (app *BaseApplication) Name() string {
	return app.name
}

func (app *BaseApplication) AddModule(module Module) error {
	app.mu.Lock()
	defer app.mu.Unlock()

	if _, exists := app.modules[module.Name()]; exists {
		return fmt.Errorf("module %s already exists", module.Name())
	}

	app.modules[module.Name()] = module
	return nil
}

func (app *BaseApplication) GetModule(name string) (Module, bool) {
	app.mu.Lock()
	defer app.mu.Unlock()
	module, exists := app.modules[name]
	return module, exists
}

func (app *BaseApplication) AddHook(phase LifecyclePhase, hook Hook) error {
	app.mu.Lock()
	defer app.mu.Unlock()

	app.hooks[phase] = append(app.hooks[phase], hook)
	return nil
}

func (app *BaseApplication) RunHooks(phase LifecyclePhase, ctx context.Context) error {
	for _, hook := range app.hooks[phase] {
		if err := hook(ctx); err != nil {
			return err
		}
	}
	return nil
}

func (app *BaseApplication) Init(ctx context.Context) error {
	// 打印 banner
	logger.L().Info(app.Banner())
	if err := app.RunHooks(PhaseBeforeInit, ctx); err != nil {
		return err
	}

	for _, module := range app.modules {
		if err := module.Init(ctx, app.ConfigManager.Get(module.Name())); err != nil {
			return err
		}
	}

	return app.RunHooks(PhaseAfterInit, ctx)
}

func (app *BaseApplication) Start(ctx context.Context) error {
	if err := app.RunHooks(PhaseBeforeStart, ctx); err != nil {
		return err
	}

	for _, module := range app.modules {
		if bootable, ok := module.(BootableModule); ok {
			if err := bootable.Start(ctx); err != nil {
				return err
			}
		}
	}

	return app.RunHooks(PhaseAfterStart, ctx)
}

func (app *BaseApplication) Stop(ctx context.Context) error {
	if err := app.RunHooks(PhaseBeforeStop, ctx); err != nil {
		return err
	}

	for _, module := range app.modules {
		if bootable, ok := module.(BootableModule); ok {
			if err := bootable.Stop(ctx); err != nil {
				return err
			}
		}
	}

	return app.RunHooks(PhaseAfterStop, ctx)
}
