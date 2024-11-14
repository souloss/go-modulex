package config

import (
	"time"

	"github.com/souloss/go-clean-arch/pkg/logger"
)

// Options 定义配置选项
type Options struct {
	// 配置文件路径
	ConfigFile string
	// 配置文件类型(如 "yaml", "json")
	ConfigType string
	// 是否监听配置文件变化
	EnableWatcher bool
	// 配置变化回调函数
	OnConfigChange func()
	// 自动重载配置的时间间隔
	ReloadInterval time.Duration
	// 环境变量前缀
	EnvPrefix string
	// 是否自动加载环境变量
	AutomaticEnv bool
	// 默认配置映射
	Defaults map[string]interface{}
	// 日志接口
	Logger logger.FormatStrLogger
}

// Option 定义配置选项的函数类型
type Option func(*Options)

// WithConfigFile 设置配置文件路径
func WithConfigFile(file string) Option {
	return func(o *Options) {
		o.ConfigFile = file
	}
}

// WithConfigType 设置配置文件类型
func WithConfigType(configType string) Option {
	return func(o *Options) {
		o.ConfigType = configType
	}
}

// WithWatcher 启用配置文件监听
func WithWatcher(onChange func()) Option {
	return func(o *Options) {
		o.EnableWatcher = true
		o.OnConfigChange = onChange
	}
}

// WithReloadInterval 设置自动重载间隔
func WithReloadInterval(d time.Duration) Option {
	return func(o *Options) {
		o.ReloadInterval = d
	}
}

// WithEnvPrefix 设置环境变量前缀
func WithEnvPrefix(prefix string) Option {
	return func(o *Options) {
		o.EnvPrefix = prefix
	}
}

// WithAutomaticEnv 启用自动加载环境变量
func WithAutomaticEnv() Option {
	return func(o *Options) {
		o.AutomaticEnv = true
	}
}

// WithDefaults 设置默认值
func WithDefaults(defaults map[string]interface{}) Option {
	return func(o *Options) {
		o.Defaults = defaults
	}
}

// WithLogger 设置日志器
func WithLogger(logger logger.FormatStrLogger) Option {
	return func(o *Options) {
		o.Logger = logger
	}
}
