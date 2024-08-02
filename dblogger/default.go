package dblogger

import (
	"go.uber.org/zap/zapcore"

	glog "github.com/jianlu8023/go-logger"
)

var (
	defaultConsoleConfig = zapcore.EncoderConfig{
		MessageKey:     "msg",
		LevelKey:       "",
		TimeKey:        "",
		NameKey:        "",
		CallerKey:      "",
		FunctionKey:    "",
		StacktraceKey:  "",
		SkipLineEnding: false,
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    glog.CustomColorCapitalLevelEncoder,
		EncodeTime:     glog.CustomTimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
		EncodeName:     zapcore.FullNameEncoder,
	}

	defaultDBLogger = glog.NewLogger(
		&glog.Config{
			LogLevel:    "DEBUG",
			DevelopMode: true,
			StackLevel:  "ERROR",
			Caller:      false,
			ModuleName:  "",
		},
		glog.WithConsoleFormat(),
		glog.WithConsoleConfig(defaultConsoleConfig),
	)
)
