package bootstrap

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
	_ "time/tzdata"

	"github.com/souloss/go-clean-arch/pkg/config"
	"github.com/souloss/go-clean-arch/pkg/logger"
	"github.com/souloss/go-clean-arch/pkg/modulex"
)

type Application struct {
	modulex.BaseApplication
	name      string
	AppConfig *AppConfig
}

func NewApp(cfgPath string, modules ...modulex.Module) *Application {
	var err error
	app := &Application{
		name:            "MyApplication",
		BaseApplication: *modulex.NewBaseApplication(),
	}
	appConfig := &AppConfig{}
	app.ConfigManager, err = config.New(appConfig, config.WithDefaults(defaults))
	if err != nil {
		logger.L().Error("App config or engine is not initialized properly")
		panic("new app failed")
	}
	app.AppConfig = app.ConfigManager.Values().(*AppConfig)

	for _, module := range modules {
		app.AddModule(module)
	}

	return app
}

func (a *Application) Name() string {
	return a.name
}

// Run run the application with graceful shutdown
func (app *Application) Run() {
	// Set timezone
	tz, err := time.LoadLocation(app.AppConfig.TimeZone)
	if err != nil {
		log.Fatal(err)
	}
	time.Local = tz

	ctx := context.Background()

	if err := app.Init(ctx); err != nil {
		logger.L().Errorf("Application initialized failed, err: %v", err)
		panic("app init failed")
	}

	if err := app.Start(ctx); err != nil {
		logger.L().Errorf("Application start failed, err: %v", err)
		panic("app start failed")
	}

	// Create a channel to listen for an interrupt or terminate signal from the OS.
	// Use a buffered channel because the signal package requires it.
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	// Block until we receive our signal.
	select {
	case <-shutdown:
		logger.L().Info("Starting graceful shutdown...")

		// Create a deadline to wait for.
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		// Doesn't block if no connections, but will otherwise wait
		// until the timeout deadline.
		err := app.Stop(ctx)
		if err != nil {
			// If we encounter an error during shutdown, log it and then try to close the server.
			logger.L().Errorf("Could not stop server gracefully: %v", err)
			if closeErr := app.Stop(ctx); closeErr != nil {
				// If the server couldn't close, log it as a fatal error.
				logger.L().Errorf("Could not close http server: %v", closeErr)
			}
		}
	}
}
