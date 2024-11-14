package server

import (
	"context"

	"github.com/mitchellh/mapstructure"
	"github.com/souloss/go-clean-arch/pkg/modulex"
	"github.com/souloss/go-clean-arch/pkg/server"
	"github.com/souloss/go-clean-arch/pkg/server/gin"
)

type ServerConfig struct {
	server.ServerConfig `mapstructure:",squash"`
	Options             server.Options `mapstructure:"options"`
}

type ServerModule struct {
	modulex.BaseBootableModule
	Engine server.WebFramework
}

func New() *ServerModule {
	ginFrame := gin.NewGinFramework()
	ginFrame.RegisterConverters(&gin.JSONWithDataConvert{}, &gin.JSONWithReqAndDataConvert{}, &gin.JSONWithReqConvert{})
	return &ServerModule{
		Engine:             ginFrame,
		BaseBootableModule: *modulex.NewBaseBootableModule(),
	}
}

func (l *ServerModule) Name() string {
	return "server"
}

// Init 初始化全局 zapLogger 和 pkgLogger
//
//	@receiver s
//	@param ctx
//	@param cfg
//	@return error
func (s *ServerModule) Init(ctx context.Context, cfg any) error {
	serverCfg := &ServerConfig{}

	if err := mapstructure.Decode(cfg, serverCfg); err != nil {
		return err
	}
	// init engine and routes
	if err := s.Engine.Init(&serverCfg.ServerConfig, &serverCfg.Options); err != nil {
		return err
	}
	return nil
}

// Start starts the module
//
//	@receiver s
//	@param ctx
//	@return error
func (s *ServerModule) Start(ctx context.Context) error {
	if err := s.Engine.Start(ctx); err != nil {
		return err
	}
	return nil
}

// Stop stops the module
//
//	@receiver s
//	@param ctx
//	@return error
func (s *ServerModule) Stop(ctx context.Context) error {
	if err := s.Engine.Shutdown(ctx); err != nil {
		return err
	}
	return nil
}
