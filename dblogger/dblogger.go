package dblogger

import (
	"time"

	"github.com/bytedance/sonic"
	"go.uber.org/zap"
)

type LogLevel int

const (
	OFF = iota + 1
	DEBUG
	INFO
	WARN
	ERROR
)

type Logger struct {
	ZapLogger                 *zap.Logger
	LogLevel                  LogLevel
	SlowThreshold             time.Duration
	Colorful                  bool
	IgnoreRecordNotFoundError bool
	ParameterizedQueries      bool
	showSql                   bool
}

func (l *Logger) String() string {
	bytes, _ := sonic.Marshal(l)
	return string(bytes)
}

func NewDBLogger(config Config) *Logger {
	return &Logger{
		ZapLogger:                 config.Logger,
		LogLevel:                  config.LogLevel,
		SlowThreshold:             config.SlowThreshold,
		Colorful:                  config.Colorful,
		IgnoreRecordNotFoundError: config.IgnoreRecordNotFoundError,
		ParameterizedQueries:      config.ParameterizedQueries,
		showSql:                   config.ShowSql,
	}
}
