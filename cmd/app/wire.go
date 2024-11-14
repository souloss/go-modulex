//go:build wireinject
// +build wireinject

//go:generate go run -mod=mod github.com/google/wire/cmd/wire
package main

import (
	"context"

	"github.com/google/wire"
	"github.com/souloss/go-clean-arch/bootstrap"
	"github.com/souloss/go-clean-arch/entities"
	"github.com/souloss/go-clean-arch/internal/query"
	"github.com/souloss/go-clean-arch/internal/repository"
	"github.com/souloss/go-clean-arch/internal/routes"
	"github.com/souloss/go-clean-arch/internal/usecase"
	"github.com/souloss/go-clean-arch/pkg/database"
	"github.com/souloss/go-clean-arch/pkg/modulex"
	databasemod "github.com/souloss/go-clean-arch/pkg/modulex/x/database"
	"github.com/souloss/go-clean-arch/pkg/modulex/x/logger"
	servermod "github.com/souloss/go-clean-arch/pkg/modulex/x/server"
	"github.com/souloss/go-clean-arch/pkg/server"
	"github.com/souloss/go-clean-arch/pkg/server/gin/middleware"
	"gorm.io/gen"
	"gorm.io/gorm"
)

func provideDOOptions() []gen.DOOption {
	return []gen.DOOption{}
}

func provideGORMDB() *gorm.DB {
	return database.GetDB(context.Background())
}

// provideRouteGroups 构建路由组
//
//	@param articleUsecase
//	@param authorUsecase
//	@return *server.RouteGroup
func provideRouteGroups(
	articleUsecase *usecase.ArticleUsecase,
	authorUsecase *usecase.AuthorUsecase,
) *server.RouteGroup {
	versionRoute := routes.GetVersionRouter()
	articleRG := routes.GetArticleRouteGroup(articleUsecase)
	authorRG := routes.GetAuthorRouteGroup(authorUsecase)
	sg := server.NewGroup(
		"/api/v1",
		server.WithRoutes(versionRoute),
		server.WithGroups(*articleRG, *authorRG),
		server.WithMiddleware(
			// otelgin.Middleware(version.ProjectName),
			middleware.NewTraceMiddleware(),
			middleware.RequestLogMiddleware(),
			middleware.ResponseLogMiddleware(),
		),
	)
	return sg
}

// GetRouteGroups 构建和获取路由组
//
//	@return *server.RouteGroup
func GetRouteGroups() *server.RouteGroup {
	wire.Build(
		// 提供ORM
		provideGORMDB,
		query.Use,
		provideDOOptions,
		repository.NewAuthorRepo,
		repository.NewArticleRepo,
		// 提供用例
		usecase.NewAuthorUsecase,
		usecase.NewArticleUsecase,
		// 提供 routeGroup
		provideRouteGroups,
		wire.Bind(new(entities.AuthorRepository), new(*repository.AuthorRepo)),
		wire.Bind(new(entities.ArticleRepository), new(*repository.ArticleRepo)),
	)
	return &server.RouteGroup{}
}

// initializeWEBModules 初始化WEB框架模块，在启动前注册路由组
//
//	@return modulex.Module
func initializeWEBModules() modulex.Module {
	serverModule := servermod.New()
	hookedServerModule := modulex.NewHookedBootableModule(serverModule)
	hookedServerModule.AddHook(modulex.PhaseBeforeStart, func(ctx context.Context) error {
		// 启动前
		rg := GetRouteGroups()
		if rg != nil {
			serverModule.Engine.RegisterGroups([]server.RouteGroup{*rg})
		}
		return nil
	})
	return hookedServerModule
}

// provideModules 提供应用模块
//
//	@return []modulex.Module
func provideModules() []modulex.Module {
	loggerMod := logger.New()
	databaseMod := databasemod.New()
	// telemetryMod := telemetry.New()

	return []modulex.Module{
		// telemetryMod,
		loggerMod,
		databaseMod,
		initializeWEBModules(),
	}
}

// initializeApp 初始化 Application 实例
//
//	@param cfgPath
//	@return *bootstrap.Application
//	@return error
func initializeApp(_ string) (*bootstrap.Application, error) {
	wire.Build(
		bootstrap.NewApp,
		provideModules,
	)
	return &bootstrap.Application{}, nil
}
