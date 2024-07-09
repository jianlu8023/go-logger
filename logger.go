package go_logger

import (
	rotateloggers "github.com/lestrrat-go/file-rotatelogs"
	"go.uber.org/zap"
	"gopkg.in/natefinch/lumberjack.v2"
)

type Logger struct {
	ZapLogger        *zap.Logger
	LumberjackLogger *lumberjack.Logger
	RotateLogger     *rotateloggers.RotateLogs
}

func NewLogger(config *Config) *zap.Logger {
	return nil
}

type Option func()

func NewSugaredLogger(config *Config, options ...Option) *zap.SugaredLogger {
	return NewLogger(config).Sugar()
}
