package db_logger

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"path/filepath"
	"regexp"
	"time"
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

// String 输出 Config 的 JSON 表示,marshal 失败时回退到 fmt 格式。
func (c *Config) String() string {
	bytes, err := json.Marshal(c)
	if err != nil {
		return fmt.Sprintf("Config{LogLevel: %d, SlowThreshold: %v, Colorful: %v}", c.LogLevel, c.SlowThreshold, c.Colorful)
	}
	return string(bytes)
}

// ctxKey 是 context.Context 的 key 类型,使用自定义类型避免与其它包的 string key 冲突。
type ctxKey int

const (
	ctxLoggerKey ctxKey = iota
	SessionIDKey
)

var (
	gormPackage    = filepath.Join("gorm.io", "gorm")
	zapgormPackage = filepath.Join("moul.io", "zapgorm2")
)

type LogContext struct {
	Ctx         context.Context
	SQL         string
	Args        []interface{}
	Result      sql.Result
	ExecuteTime time.Duration
	Err         error
}

var (
	sqlRegexp                       = regexp.MustCompile(`\?`)
	numericPlaceHolderRegexp        = regexp.MustCompile(`\$\d+`)
	numericPlaceHolderReplaceRegexp = regexp.MustCompile(`\$(\d+)`)
)
