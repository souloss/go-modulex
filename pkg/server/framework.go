package server

import (
	"context"
	"fmt"
	"net/http"
	"path"
	"strings"
	"time"
)

// ServerConfig 定义基础服务器配置
type ServerConfig struct {
	Name           string        `json:"name" yaml:"name"` // 服务名称
	Host           string        `json:"host" yaml:"host"` // 服务地址
	Port           int           `json:"port" yaml:"port"` // 服务端口
	Mode           string        `json:"mode" yaml:"mode"` // 运行模式(debug/release/test)
	ReadTimeout    time.Duration `json:"read_timeout" yaml:"read_timeout"`
	WriteTimeout   time.Duration `json:"write_timeout" yaml:"write_timeout"`
	MaxHeaderBytes int           `json:"max_header_bytes" yaml:"max_header_bytes"`
}

// Options 定义框架配置选项
type Options struct {
	EnableMetrics  bool     `json:"enable_metrics" yaml:"enable_metrics"`   // 是否启用指标收集
	EnableTrace    bool     `json:"enable_trace" yaml:"enable_trace"`       // 是否启用链路追踪
	AllowedOrigins []string `json:"allowed_origins" yaml:"allowed_origins"` // CORS 配置
	TrustedProxies []string `json:"trusted_proxies" yaml:"trusted_proxies"` // 受信任的代理
}

// Route 定义路由
type Route struct {
	Method      string        // HTTP 方法
	Path        string        // 路径
	Handler     interface{}   // 处理函数
	Middlewares []interface{} // 中间件列表

	// meta
	AbsolutePath string
	// swagger
	Swagger  *SwaggerInfo  // Swagger 文档信息
	Request  *RequestInfo  // 请求信息
	Response *ResponseInfo // 响应信息
}

// RouteGroup 定义路由组
type RouteGroup struct {
	Prefix      string        // 路由组前缀
	Middlewares []interface{} // 路由组中间件
	Routes      []Route       // 路由列表
	Groups      []RouteGroup  // 子路由组

	// swagger
	Tags        []string // 该组的 Swagger 标签
	Description string   // 该组的描述信息
}

// WebFramework 定义框架接口
type WebFramework interface {
	// 初始化框架
	Init(cfg *ServerConfig, opts *Options) error

	// 注册路由
	RegisterRoutes(routes []Route) error

	// 注册路由组
	RegisterGroups(groups []RouteGroup) error

	// 添加中间件
	Use(middleware ...interface{})

	// 启动服务
	Start(ctx context.Context) error

	// 优雅关闭
	Shutdown(ctx context.Context) error

	// 获取原生框架实例
	Native() interface{}
}

// GroupOption 定义路由组配置函数
type GroupOption func(*RouteGroup)

// WithMiddleware 添加中间件的选项函数
func WithMiddleware(middlewares ...interface{}) GroupOption {
	return func(g *RouteGroup) {
		g.Middlewares = append(g.Middlewares, middlewares...)
	}
}

// WithRoutes 添加路由的选项函数
func WithRoutes(routes ...Route) GroupOption {
	return func(g *RouteGroup) {
		g.Routes = append(g.Routes, routes...)
	}
}

// WithGroups 添加子组的选项函数
func WithGroups(groups ...RouteGroup) GroupOption {
	return func(g *RouteGroup) {
		g.Groups = append(g.Groups, groups...)
	}
}

// NewGroup 创建新的路由组
func NewGroup(prefix string, opts ...GroupOption) *RouteGroup {
	group := &RouteGroup{
		Prefix:      prefix,
		Middlewares: make([]interface{}, 0),
		Routes:      make([]Route, 0),
		Groups:      make([]RouteGroup, 0),
	}

	for _, opt := range opts {
		opt(group)
	}

	return group
}

// NewRoute 创建路由
//
//	@param method
//	@param path
//	@param handler
//	@param middlewares
//	@return *Route
func NewRoute(method, path string, handler interface{}, middlewares ...interface{}) *Route {
	route := &Route{
		Method:      method,
		Path:        path,
		Handler:     handler,
		Middlewares: middlewares,
	}
	return route
}

// BaseFramework
type BaseFramework struct {
	config  *ServerConfig
	options *Options
	server  *http.Server
}

// Init
//
//	@receiver b
//	@param cfg
//	@param opts
//	@return error
func (b *BaseFramework) Init(cfg *ServerConfig, opts *Options) error {
	b.config = cfg
	b.options = opts
	b.server = &http.Server{
		Addr:           fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		ReadTimeout:    cfg.ReadTimeout,
		WriteTimeout:   cfg.WriteTimeout,
		MaxHeaderBytes: cfg.MaxHeaderBytes,
	}
	return nil
}

// GetServer
//
//	@receiver b
//	@return *http.Server
func (b *BaseFramework) GetServer() *http.Server {
	return b.server
}

func (b *BaseFramework) ParseRoutes(groups ...*RouteGroup) {
	for i := range groups {
		parseGroup(groups[i], "")
	}
}

// parseGroup 递归解析单个路由组
func parseGroup(group *RouteGroup, parentPrefix string) {
	// 计算当前组的完整前缀
	currentPrefix := path.Join(parentPrefix, group.Prefix)

	// 规范化前缀
	currentPrefix = normalizePath(currentPrefix)

	// 处理当前组的所有路由
	for i := range group.Routes {
		route := &group.Routes[i]
		absolutePath := path.Join(currentPrefix, route.Path)
		route.AbsolutePath = normalizePath(absolutePath)
	}

	// 递归处理所有子组
	for i := range group.Groups {
		parseGroup(&group.Groups[i], currentPrefix)
	}
}

// normalizePath 规范化路径格式
func normalizePath(p string) string {
	// 确保路径以 / 开头
	if !strings.HasPrefix(p, "/") {
		p = "/" + p
	}

	// 处理连续的斜杠
	for strings.Contains(p, "//") {
		p = strings.ReplaceAll(p, "//", "/")
	}

	// 处理路径末尾的斜杠（保留单个 / 的情况）
	if p != "/" && strings.HasSuffix(p, "/") {
		p = p[:len(p)-1]
	}

	return p
}
