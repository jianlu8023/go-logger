package dblogger

import (
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

	defaultDBLogger = glog.NewLogger(
		glog.WithCaller(),
		glog.WithDevelopMode(),
		glog.WithDefaultLogLevel("DEBUG"),
		glog.WithStackLogLevel("ERROR"),
		glog.WithConsoleFormat(),
		glog.WithConsoleConfig(defaultConsoleConfig),
	)
)
