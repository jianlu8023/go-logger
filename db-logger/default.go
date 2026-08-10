package db_logger

import (
	"sync"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	glog "github.com/jianlu8023/go-logger/v2"
)

var (
	defaultConsoleConfig = zapcore.EncoderConfig{
		MessageKey:       "msg",
		LevelKey:         "",
		TimeKey:          "",
		NameKey:          "",
		CallerKey:        "",
		FunctionKey:      "",
		StacktraceKey:    "",
		SkipLineEnding:   false,
		LineEnding:       zapcore.DefaultLineEnding,
		EncodeLevel:      glog.CustomColorCapitalLevelEncoder,
		EncodeTime:       glog.CustomTimeEncoder,
		EncodeDuration:   zapcore.SecondsDurationEncoder,
		EncodeCaller:     zapcore.ShortCallerEncoder,
		EncodeName:       zapcore.FullNameEncoder,
		ConsoleSeparator: "  ",
	}

	defaultDBLogger     *zap.Logger
	defaultDBLoggerOnce sync.Once
)

// getDefaultDBLogger 懒加载默认 zap.Logger 实例。
// 使用 sync.Once 确保仅在首次调用时初始化,避免 import 时的全局副作用,
// 同时保证多次调用返回同一实例。
func getDefaultDBLogger() *zap.Logger {
	defaultDBLoggerOnce.Do(func() {
		defaultDBLogger = glog.NewLogger(
			glog.WithCaller(),
			glog.WithDevelopMode(),
			glog.WithDefaultLogLevel("DEBUG"),
			glog.WithStackLogLevel("ERROR"),
			glog.WithConsoleFormat(),
			glog.WithConsoleConfig(defaultConsoleConfig),
		)
	})
	return defaultDBLogger
}
