package echo

import (
	"context"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/souloss/go-clean-arch/pkg/server"
)

type EchoFramework struct {
	*server.BaseFramework
	engine *echo.Echo
}

func NewEchoFramework() *EchoFramework {
	return &EchoFramework{}
}

func (e *EchoFramework) Init(cfg *server.ServerConfig, opts *server.Options) error {
	if err := e.BaseFramework.Init(cfg, opts); err != nil {
		return err
	}

	e.engine = echo.New()

	// 设置运行模式
	e.engine.Debug = cfg.Mode == "debug"

	// 配置 CORS
	if len(opts.AllowedOrigins) > 0 {
		e.engine.Use(middleware.CORSWithConfig(middleware.CORSConfig{
			AllowOrigins: opts.AllowedOrigins,
			AllowMethods: []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete, http.MethodOptions},
		}))
	}

	// 设置 HTTP 服务器处理器
	e.GetServer().Handler = e.engine

	return nil
}

func (e *EchoFramework) RegisterRoutes(routes []server.Route) error {
	for _, route := range routes {
		handler, ok := route.Handler.(echo.HandlerFunc)
		if !ok {
			return fmt.Errorf("invalid handler type for echo")
		}

		// 转换中间件
		middlewares := make([]echo.MiddlewareFunc, 0, len(route.Middlewares))
		for _, m := range route.Middlewares {
			if mw, ok := m.(echo.MiddlewareFunc); ok {
				middlewares = append(middlewares, mw)
			}
		}

		e.engine.Add(route.Method, route.Path, handler, middlewares...)
	}
	return nil
}

// EchoFramework 的路由组注册实现
func (e *EchoFramework) RegisterGroups(groups []server.RouteGroup) error {
	for _, group := range groups {
		if err := e.registerGroup(e.engine.Group(group.Prefix), group); err != nil {
			return err
		}
	}
	return nil
}

// registerGroup 递归注册路由组
func (e *EchoFramework) registerGroup(echoGroup *echo.Group, group server.RouteGroup) error {
	// 添加组级中间件
	for _, m := range group.Middlewares {
		if mw, ok := m.(echo.MiddlewareFunc); ok {
			echoGroup.Use(mw)
		}
	}

	// 注册当前组的路由
	for _, route := range group.Routes {
		handler, ok := route.Handler.(echo.HandlerFunc)
		if !ok {
			return fmt.Errorf("invalid handler type for echo group")
		}

		// 转换路由级中间件
		middlewares := make([]echo.MiddlewareFunc, 0, len(route.Middlewares))
		for _, m := range route.Middlewares {
			if mw, ok := m.(echo.MiddlewareFunc); ok {
				middlewares = append(middlewares, mw)
			}
		}

		echoGroup.Add(route.Method, route.Path, handler, middlewares...)
	}

	// 递归注册子组
	for _, subGroup := range group.Groups {
		subEchoGroup := echoGroup.Group(subGroup.Prefix)
		if err := e.registerGroup(subEchoGroup, subGroup); err != nil {
			return err
		}
	}

	return nil
}

func (e *EchoFramework) Use(middleware ...interface{}) {
	for _, m := range middleware {
		if mw, ok := m.(echo.MiddlewareFunc); ok {
			e.engine.Use(mw)
		}
	}
}

func (e *EchoFramework) Start(ctx context.Context) error {
	return e.GetServer().ListenAndServe()
}

func (e *EchoFramework) Shutdown(ctx context.Context) error {
	return e.GetServer().Shutdown(ctx)
}

func (e *EchoFramework) Native() interface{} {
	return e.engine
}
