// Code generated by Wire. DO NOT EDIT.

//go:generate go run -mod=mod github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package main

import (
	"context"
	"github.com/souloss/go-clean-arch/bootstrap"
	"github.com/souloss/go-clean-arch/internal/query"
	"github.com/souloss/go-clean-arch/internal/repository"
	"github.com/souloss/go-clean-arch/internal/routes"
	"github.com/souloss/go-clean-arch/internal/usecase"
	"github.com/souloss/go-clean-arch/pkg/database"
	"github.com/souloss/go-clean-arch/pkg/modulex"
	database2 "github.com/souloss/go-clean-arch/pkg/modulex/x/database"
	"github.com/souloss/go-clean-arch/pkg/modulex/x/logger"
	server2 "github.com/souloss/go-clean-arch/pkg/modulex/x/server"
	"github.com/souloss/go-clean-arch/pkg/server"
	"github.com/souloss/go-clean-arch/pkg/server/gin/middleware"
	"gorm.io/gen"
	"gorm.io/gorm"
)

// Injectors from wire.go:

// GetRouteGroups 构建和获取路由组
//
//	@return *server.RouteGroup
func GetRouteGroups() *server.RouteGroup {
	db := provideGORMDB()
	v := provideDOOptions()
	queryQuery := query.Use(db, v...)
	articleRepo := repository.NewArticleRepo(queryQuery)
	articleUsecase := usecase.NewArticleUsecase(articleRepo)
	authorRepo := repository.NewAuthorRepo(queryQuery)
	authorUsecase := usecase.NewAuthorUsecase(authorRepo)
	routeGroup := provideRouteGroups(articleUsecase, authorUsecase)
	return routeGroup
}

// initializeApp 初始化 Application 实例
//
//	@param cfgPath
//	@return *bootstrap.Application
//	@return error
func initializeApp(string2 string) (*bootstrap.Application, error) {
	v := provideModules()
	application := bootstrap.NewApp(string2, v...)
	return application, nil
}

// wire.go:

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
		"/api/v1", server.WithRoutes(versionRoute), server.WithGroups(*articleRG, *authorRG), server.WithMiddleware(middleware.NewTraceMiddleware(), middleware.RequestLogMiddleware(), middleware.ResponseLogMiddleware()),
	)
	return sg
}

// initializeWEBModules 初始化WEB框架模块，在启动前注册路由组
//
//	@return modulex.Module
func initializeWEBModules() modulex.Module {
	serverModule := server2.New()
	hookedServerModule := modulex.NewHookedBootableModule(serverModule)
	hookedServerModule.AddHook(modulex.PhaseBeforeStart, func(ctx context.Context) error {

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
	databaseMod := database2.New()

	return []modulex.Module{

		loggerMod,
		databaseMod,
		initializeWEBModules(),
	}
}
