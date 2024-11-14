package telemetry

import (
	"context"

	"github.com/souloss/go-clean-arch/pkg/database"
	"github.com/souloss/go-clean-arch/pkg/logger"
	"github.com/souloss/go-clean-arch/pkg/modulex"
	"github.com/souloss/go-clean-arch/pkg/telemetry"
	"github.com/uptrace/opentelemetry-go-extra/otelgorm"
	"github.com/uptrace/opentelemetry-go-extra/otelzap"
	"go.uber.org/zap"
)

type TelemetryModule struct {
	modulex.BaseBootableModule
}

func New() *TelemetryModule {
	return &TelemetryModule{
		BaseBootableModule: *modulex.NewBaseBootableModule(),
	}
}

func (l *TelemetryModule) Name() string {
	return "telemetry"
}

// Init 初始化全局 Telemetry
//
//	@receiver l
//	@param ctx
//	@param cfg
//	@return error
func (l *TelemetryModule) Init(ctx context.Context, cfg any) error {
	// if err := mapstructure.Decode(cfg, databaseCfg); err != nil {
	// 	return err
	// }
	if err := telemetry.InitOpenTelemetry(); err != nil {
		return err
	}
	// 应用到数据库
	if err := database.GetDB(ctx).Use(otelgorm.NewPlugin()); err != nil {
		return err
	}
	// 应用到Zap日志
	logger.ReplaceGlobals(logger.NewZapAdapter(otelzap.New(zap.L()).Logger))
	return nil
}
