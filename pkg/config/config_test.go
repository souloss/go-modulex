package config_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/souloss/go-clean-arch/pkg/config"
	"github.com/stretchr/testify/assert"
)

// AppConfig is a struct that holds the configuration for the application.
type AppConfig struct {
	Server struct {
		Host string `mapstructure:"host"`
		Port int    `mapstructure:"port"`
	} `mapstructure:"server"`

	// Begain with the database type, followed by the connection string. For example:
	// - type:sqlite3, url:memory -> sqlite3://:memory: creates an in-memory database.
	// - type:sqlite3, url:path/to/database.db -> sqlite3:///path/to/database.db creates a database file at the specified path.
	// - type:sqlite3, url:localhost:5432/dbname, username:user, password=password -> postgres://user:password@localhost:5432/dbname connects to a PostgreSQL database.
	Database struct {
		Type     string `mapstructure:"type"`
		URL      string `mapstructure:"url"`
		UserName string `mapstructure:"username"`
		Password string `mapstructure:"password"`
	} `mapstructure:"database"`

	Logger struct {
		Level      string `mapstructure:"level"`
		Path       string `mapstructure:"path"`
		FileName   string `mapstructure:"filename"`
		MaxSize    int    `mapstructure:"max_size"`
		MaxAge     int    `mapstructure:"max_age"`
		MaxBackups int    `mapstructure:"max_backups"`
		Compress   bool   `mapstructure:"compress"`
		Console    bool   `mapstructure:"console"`
	} `mapstructure:"logger"`
}

func TestReadConfig(t *testing.T) {
	a := assert.New(t)
	// 创建配置结构体实例
	cfg := &AppConfig{}

	// 设置默认值
	defaults := map[string]interface{}{
		"server.host":  "localhost",
		"server.port":  8080,
		"logger.level": "info",
	}

	// 创建配置管理器
	c, err := config.New(
		cfg,
		config.WithConfigFile("config.yaml"),
		config.WithConfigType("yaml"),
		config.WithDefaults(defaults),
		config.WithEnvPrefix("APP"),
		config.WithAutomaticEnv(),
		config.WithWatcher(func() {
			fmt.Println("config updated!")
		}),
		config.WithReloadInterval(time.Second*5),
	)
	if err != nil {
		panic(err)
	}

	// 配置管理器获取配置
	a.Equal(c.GetString("server.host"), "localhost")
	a.Equal(c.GetInt("server.port"), 8080)
	a.Equal(c.GetString("logger.level"), "info")

	// 结构体获取配置
	a.Equal(cfg.Server.Host, "localhost")
	a.Equal(cfg.Server.Port, 8080)
	a.Equal(cfg.Logger.Level, "info")

	// 配置管理器获取结构体
	appCfg := c.Values().(*AppConfig)
	a.Equal(cfg, appCfg)
}

func TestNopConfig(t *testing.T) {
	a := assert.New(t)
	cfg, err := config.New(&struct{}{})
	a.Equal(nil, err)
	a.NotNil(cfg)
}
