package db_logger

import (
	"fmt"
	"os"
	"time"

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

// Logger 是 db-logger 的核心结构体,同时适配 GORM 与 XORM。
//
// 并发说明:LogLevel、Colorful、showSql 等字段为非原子读写,
// Logger 设计为创建后不可变(SetLevel/ShowSQL 仅在初始化阶段调用)。
// 若需要在运行时并发修改级别,请通过 NewDBLogger 重建实例。
type Logger struct {
	zapLogger                 *zap.Logger `json:"-"`
	LogLevel                  LogLevel
	SlowThreshold             time.Duration
	Colorful                  bool
	IgnoreRecordNotFoundError bool
	ParameterizedQueries      bool
	showSql                   bool
}

// String 输出 Logger 的关键配置信息。
// zapLogger 字段不可序列化(含 sync.Mutex 等),通过 json:"-" 跳过。
func (l Logger) String() string {
	return fmt.Sprintf("Logger{LogLevel: %d, SlowThreshold: %v, Colorful: %v, IgnoreRecordNotFoundError: %v, ParameterizedQueries: %v, showSql: %v}",
		l.LogLevel, l.SlowThreshold, l.Colorful,
		l.IgnoreRecordNotFoundError, l.ParameterizedQueries, l.showSql)
}

// NewDBLogger 根据配置创建 Logger。
// 可通过 WithCustomLogger 传入自定义 *zap.Logger;若未传入或类型断言失败,
// 会降级到默认 logger(懒加载)并在 stderr 输出告警(不会 panic)。
func NewDBLogger(config Config, options ...Option) *Logger {
	zapLogger := getDefaultDBLogger()
	for _, opt := range options {
		if opt.Name() == customLoggerKey {
			custom, ok := opt.Value().(*zap.Logger)
			if !ok {
				fmt.Fprintf(os.Stderr, "[db-logger] WithCustomLogger value is not *zap.Logger, fallback to default logger\n")
				continue
			}
			zapLogger = custom
		}
	}
	return &Logger{
		zapLogger:                 zapLogger,
		LogLevel:                  config.LogLevel,
		SlowThreshold:             config.SlowThreshold,
		Colorful:                  config.Colorful,
		IgnoreRecordNotFoundError: config.IgnoreRecordNotFoundError,
		ParameterizedQueries:      config.ParameterizedQueries,
		showSql:                   config.ShowSql,
	}
}
