package database

import (
	"context"

	"github.com/mitchellh/mapstructure"
	"github.com/souloss/go-clean-arch/pkg/database"
	"github.com/souloss/go-clean-arch/pkg/logger"
	"github.com/souloss/go-clean-arch/pkg/modulex"
)

type DataBaseModule struct {
	modulex.BaseBootableModule
}

func New() *DataBaseModule {
	return &DataBaseModule{
		BaseBootableModule: *modulex.NewBaseBootableModule(),
	}
}

func (l *DataBaseModule) Name() string {
	return "database"
}

// Init 初始化全局 zapDataBase 和 pkgDataBase
//
//	@receiver l
//	@param ctx
//	@param cfg
//	@return error
func (l *DataBaseModule) Init(ctx context.Context, cfg any) error {
	databaseCfg := &database.DataSource{}
	if err := mapstructure.Decode(cfg, databaseCfg); err != nil {
		return err
	}
	if err := database.InitDB(ctx, databaseCfg, logger.L().Named("database")); err != nil {
		return err
	}
	return nil
}
