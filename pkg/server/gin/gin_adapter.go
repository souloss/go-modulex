package gin

import (
	"context"
	"fmt"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/souloss/go-clean-arch/pkg/logger"
	"github.com/souloss/go-clean-arch/pkg/server"
)

// GinFramework 实现基于 Gin 的框架
type GinFramework struct {
	*server.BaseFramework
	engine       *gin.Engine
	handlerAdapt *ginHandlerAdapt
}

func NewGinFramework() *GinFramework {
	return &GinFramework{
		handlerAdapt: &ginHandlerAdapt{},
	}
}

var _ server.WebFramework = (*GinFramework)(nil)

func (g *GinFramework) RegisterConverters(c ...CtrlToGinHandlerConverter) {
	for _, item := range c {
		g.handlerAdapt.RegisterConverter(item)
	}
}

func (g *GinFramework) Init(cfg *server.ServerConfig, opts *server.Options) error {
	g.BaseFramework = &server.BaseFramework{}
	if err := g.BaseFramework.Init(cfg, opts); err != nil {
		return err
	}

	// 设置 Gin 模式
	gin.SetMode(cfg.Mode)

	// 创建 Gin 引擎
	g.engine = gin.New()

	// 设置受信任代理
	if len(opts.TrustedProxies) > 0 {
		g.engine.SetTrustedProxies(opts.TrustedProxies)
	}

	// 配置 CORS
	if len(opts.AllowedOrigins) > 0 {
		g.engine.Use(cors.New(cors.Config{
			AllowOrigins: opts.AllowedOrigins,
			AllowMethods: []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
			AllowHeaders: []string{"Origin", "Content-Type", "Accept"},
		}))
	}
	// 默认文件上传限制
	g.engine.MaxMultipartMemory = 32 << 20 // 32 MB
	// 设置 HTTP 服务器处理器
	g.GetServer().Handler = g.engine

	return nil
}

func (g *GinFramework) RegisterRoutes(routes []server.Route) error {
	for _, route := range routes {
		handler, ok := route.Handler.(gin.HandlerFunc)
		if !ok {
			convertSucc := false
			for _, warpper := range g.handlerAdapt.converter {
				if warpper.Check(route.Handler) {
					handler = warpper.Convert(route.Handler)
					convertSucc = true
				}
			}
			if !convertSucc {
				logger.L().Errorf("handler %v register failed", handler)
				return fmt.Errorf("handler register failed")
			}
		}

		// 转换中间件
		middlewares := make([]gin.HandlerFunc, 0, len(route.Middlewares))
		for _, m := range route.Middlewares {
			if mw, ok := m.(gin.HandlerFunc); ok {
				middlewares = append(middlewares, mw)
			}
		}

		handlers := append(middlewares, handler)
		g.engine.Handle(route.Method, route.Path, handlers...)
	}
	return nil
}

// GinFramework 的路由组注册实现
func (g *GinFramework) RegisterGroups(groups []server.RouteGroup) error {
	for _, group := range groups {
		if err := g.registerGroup(g.engine, group); err != nil {
			return err
		}
	}
	return nil
}

// registerGroup 递归注册路由组
func (g *GinFramework) registerGroup(router gin.IRouter, group server.RouteGroup) error {
	// 创建当前层级的路由组
	ginGroup := router.Group(group.Prefix)

	// 添加组级中间件
	for _, m := range group.Middlewares {
		if mw, ok := m.(gin.HandlerFunc); ok {
			ginGroup.Use(mw)
		}
	}

	// 注册当前组的路由
	for _, route := range group.Routes {
		handler, ok := route.Handler.(gin.HandlerFunc)
		if !ok {
			convertSucc := false
			for _, warpper := range g.handlerAdapt.converter {
				if warpper.Check(route.Handler) {
					handler = warpper.Convert(route.Handler)
					convertSucc = true
					break
				}
			}
			if !convertSucc {
				logger.L().Errorf("handler %v register failed", handler)
				return fmt.Errorf("handler register failed")
			}
		}

		// 转换路由级中间件
		middlewares := make([]gin.HandlerFunc, 0, len(route.Middlewares))
		for _, m := range route.Middlewares {
			if mw, ok := m.(gin.HandlerFunc); ok {
				middlewares = append(middlewares, mw)
			}
		}

		ginGroup.Handle(route.Method, route.Path, append(middlewares, handler)...)
	}

	// 递归注册子组
	for _, subGroup := range group.Groups {
		if err := g.registerGroup(ginGroup, subGroup); err != nil {
			return err
		}
	}

	return nil
}

func (g *GinFramework) Use(middleware ...interface{}) {
	for _, m := range middleware {
		if mw, ok := m.(gin.HandlerFunc); ok {
			g.engine.Use(mw)
		}
	}
}

func (g *GinFramework) Start(ctx context.Context) error {
	return g.GetServer().ListenAndServe()
}

func (g *GinFramework) Shutdown(ctx context.Context) error {
	return g.GetServer().Shutdown(ctx)
}

func (g *GinFramework) Native() interface{} {
	return g.engine
}
