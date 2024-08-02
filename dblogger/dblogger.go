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
	zapLogger                 *zap.Logger
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

func NewDBLogger(config Config, options ...Option) *Logger {
	for _, opt := range options {
		if opt.Name() == customLoggerKey {
			return &Logger{
				zapLogger:                 opt.Value().(*zap.Logger),
				LogLevel:                  config.LogLevel,
				SlowThreshold:             config.SlowThreshold,
				Colorful:                  config.Colorful,
				IgnoreRecordNotFoundError: config.IgnoreRecordNotFoundError,
				ParameterizedQueries:      config.ParameterizedQueries,
				showSql:                   config.ShowSql,
			}
		}
	}
	return &Logger{
		zapLogger:                 defaultDBLogger,
		LogLevel:                  config.LogLevel,
		SlowThreshold:             config.SlowThreshold,
		Colorful:                  config.Colorful,
		IgnoreRecordNotFoundError: config.IgnoreRecordNotFoundError,
		ParameterizedQueries:      config.ParameterizedQueries,
		showSql:                   config.ShowSql,
	}
}
