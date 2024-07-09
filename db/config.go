package db

import (
	"encoding/json"
	"path/filepath"
	"time"

	"go.uber.org/zap"
	gormlogger "gorm.io/gorm/logger"
)

type Config struct {
	Logger                    *zap.Logger         `json:"logger,omitempty"`
	LogLevel                  gormlogger.LogLevel `json:"logLevel,omitempty"`
	SlowThreshold             time.Duration       `json:"slowThreshold,omitempty"`
	Colorful                  bool                `json:"colorful,omitempty"`
	IgnoreRecordNotFoundError bool                `json:"ignoreRecordNotFoundError,omitempty"`
	ParameterizedQueries      bool                `json:"parameterizedQueries,omitempty"`
}

func (c *Config) String() string {
	bytes, _ := json.Marshal(c)
	return string(bytes)
}

const ctxLoggerKey = "zapLogger"

var (
	gormPackage = filepath.Join("gorm.io", "gorm")
)
