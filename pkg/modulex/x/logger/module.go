package logger

import (
	"context"

	"github.com/mitchellh/mapstructure"
	"github.com/souloss/go-clean-arch/pkg/logger"
	"github.com/souloss/go-clean-arch/pkg/modulex"
	"go.uber.org/zap"
)

type LoggerModule struct {
	modulex.BaseBootableModule
}

func New() *LoggerModule {
	return &LoggerModule{
		BaseBootableModule: *modulex.NewBaseBootableModule(),
	}
}

func (l *LoggerModule) Name() string {
	return "logger"
}

// Init 初始化全局 zapLogger 和 pkgLogger
//
//	@receiver l
//	@param ctx
//	@param cfg
//	@return error
func (l *LoggerModule) Init(ctx context.Context, cfg any) error {
	logCfg := &logger.LoggerConfig{}
	if err := mapstructure.Decode(cfg, logCfg); err != nil {
		return err
	}
	if err := logger.InitZapGlobalLogger(ctx, logCfg); err != nil {
		return err
	}
	logger.ReplaceGlobals(logger.NewZapAdapter(zap.L()))
	return nil
}
