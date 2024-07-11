package db

import (
	"encoding/json"
	"path/filepath"
	"time"

	"go.uber.org/zap"
	gormlogger "gorm.io/gorm/logger"
)

type LogLevel int

const (
	// Silent silent log level
	Silent LogLevel = iota + 1
	// Error error log level
	Error
	// Warn warn log level
	Warn
	// Info info log level
	Info
)

var (
	dbMap = map[LogLevel]gormlogger.LogLevel{
		Silent: gormlogger.Silent,
		Error:  gormlogger.Error,
		Warn:   gormlogger.Warn,
		Info:   gormlogger.Info,
	}
)

type Config struct {
	Logger                    *zap.Logger   `json:"logger,omitempty"`
	LogLevel                  LogLevel      `json:"logLevel,omitempty"`
	SlowThreshold             time.Duration `json:"slowThreshold,omitempty"`
	Colorful                  bool          `json:"colorful,omitempty"`
	IgnoreRecordNotFoundError bool          `json:"ignoreRecordNotFoundError,omitempty"`
	ParameterizedQueries      bool          `json:"parameterizedQueries,omitempty"`
}

func (c *Config) String() string {
	bytes, _ := json.Marshal(c)
	return string(bytes)
}

const ctxLoggerKey = "zapLogger"

var (
	gormPackage    = filepath.Join("gorm.io", "gorm")
	zapgormPackage = filepath.Join("moul.io", "zapgorm2")
)
