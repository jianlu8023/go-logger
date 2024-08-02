package dblogger

import (
	"path/filepath"
	"time"

	"github.com/bytedance/sonic"
)

type Config struct {
	// Logger                    *zap.Logger   `json:"logger,omitempty"`
	LogLevel                  LogLevel      `json:"logLevel,omitempty"`
	SlowThreshold             time.Duration `json:"slowThreshold,omitempty"`
	Colorful                  bool          `json:"colorful,omitempty"`
	IgnoreRecordNotFoundError bool          `json:"ignoreRecordNotFoundError,omitempty"`
	ParameterizedQueries      bool          `json:"parameterizedQueries,omitempty"`
	ShowSql                   bool          `json:"showSql,omitempty"`
}

func (c *Config) String() string {
	bytes, _ := sonic.Marshal(c)
	return string(bytes)
}

const (
	ctxLoggerKey = "zapLogger"
	SessionIDKey = "__xorm_session_id"
)

var (
	gormPackage    = filepath.Join("gorm.io", "gorm")
	zapgormPackage = filepath.Join("moul.io", "zapgorm2")
)
