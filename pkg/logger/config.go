package logger

type LoggerConfig struct {
	Level    string `mapstructure:"level"`
	FileName string `mapstructure:"filename"`
	Console  bool   `mapstructure:"console"`

	MaxSize      int  `mapstructure:"max_size"`
	MaxAge       int  `mapstructure:"max_age"`
	RotationTime int  `mapstructure:"rotation_time"`
	MaxBackups   int  `mapstructure:"max_backups"`
	Compress     bool `mapstructure:"compress"`
	LocalTime    bool `mapstructure:"local_time"`
}
