package config

import (
	"fmt"
	"sync"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"

	"github.com/souloss/go-clean-arch/pkg/logger"
)

const moduleName = "configManager"

// ConfigManager 配置管理器
type ConfigManager struct {
	v         *viper.Viper
	mu        sync.RWMutex
	opts      Options
	values    interface{}
	validator *validator.Validate
	logger    logger.FormatStrLogger
}

// New 创建新的配置管理器
func New(values interface{}, opts ...Option) (*ConfigManager, error) {
	// 默认选项
	options := Options{
		ReloadInterval: time.Second * 5,
		Logger:         logger.L().Named(moduleName),
	}

	for _, opt := range opts {
		opt(&options)
	}

	c := &ConfigManager{
		v:         viper.New(),
		opts:      options,
		values:    values,
		logger:    options.Logger,
		validator: validator.New(),
	}

	if err := c.init(); err != nil {
		return nil, err
	}

	return c, nil
}

// init 初始化配置
func (c *ConfigManager) init() error {
	// 设置默认值
	if c.opts.Defaults != nil {
		for k, v := range c.opts.Defaults {
			c.v.SetDefault(k, v)
		}
	}

	// 设置环境变量
	if c.opts.EnvPrefix != "" {
		c.v.SetEnvPrefix(c.opts.EnvPrefix)
	}
	if c.opts.AutomaticEnv {
		c.v.AutomaticEnv()
	}

	fileExist := true
	if c.opts.ConfigFile != "" {
		c.v.SetConfigFile(c.opts.ConfigFile)
		c.v.SetConfigType(c.opts.ConfigType)

		// 读取配置
		if err := c.v.ReadInConfig(); err != nil {
			c.logger.Warnf("Failed read in config: %v, use default config", err)
			fileExist = false
		}
	}

	// 解析到结构体
	if err := c.v.Unmarshal(c.values); err != nil {
		return fmt.Errorf("Failed unmarshal config: %w", err)
	}

	// 启用配置文件监听
	if c.opts.EnableWatcher && fileExist {
		c.v.WatchConfig()
		c.v.OnConfigChange(func(e fsnotify.Event) {
			c.mu.Lock()
			defer c.mu.Unlock()

			if err := c.v.Unmarshal(c.values); err != nil {
				c.logger.Warnf("Failed reload config: %v\n", err)
				return
			}

			if c.opts.OnConfigChange != nil {
				c.opts.OnConfigChange()
			}
		})
	}

	return c.validator.Struct(c.values)
}

// Get 获取配置值
func (c *ConfigManager) Get(key string) interface{} {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.v.Get(key)
}

// GetString 获取字符串配置
func (c *ConfigManager) GetString(key string) string {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.v.GetString(key)
}

// GetInt 获取整数配置
func (c *ConfigManager) GetInt(key string) int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.v.GetInt(key)
}

// GetBool 获取布尔配置
func (c *ConfigManager) GetBool(key string) bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.v.GetBool(key)
}

// GetDuration 获取时间间隔配置
func (c *ConfigManager) GetDuration(key string) time.Duration {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.v.GetDuration(key)
}

// Set 设置配置值
func (c *ConfigManager) Set(key string, value interface{}) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.v.Set(key, value)
}

// Values 获取完整的配置结构体
func (c *ConfigManager) Values() interface{} {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.values
}

// WriteConfig 写入配置文件
func (c *ConfigManager) WriteConfig() error {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.v.WriteConfig()
}
