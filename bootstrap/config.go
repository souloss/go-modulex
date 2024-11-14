package bootstrap

import (
	"github.com/souloss/go-clean-arch/pkg/database"
	"github.com/souloss/go-clean-arch/pkg/logger"
	"github.com/souloss/go-clean-arch/pkg/modulex/x/server"
)

// AppConfig is a struct that holds the configuration for the application.
type AppConfig struct {

	// Server 服务器配置
	Server server.ServerConfig `mapstructure:"server"`

	// 数据源配置
	Database database.DataSource `mapstructure:"database"`

	// 日志配置
	Logger logger.LoggerConfig `mapstructure:"logger"`

	// 时区配置
	TimeZone string `mapstructure:"time_zone"`
}

var defaults = map[string]interface{}{
	"server.host":    "localhost",
	"database.dsn":   "postgres://postgres:postgres@localhost:5432/myapp?sslmode=disable&TimeZone=Asia/Shanghai",
	"server.port":    8080,
	"logger.level":   "debug",
	"logger.console": true,
	"time_zone":      "Local",
}
