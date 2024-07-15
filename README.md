# go-logger 基于 zap 结合 lumberjack rotate-logs

## Quick Start

```go
import (
glog "github.com/jianlu8023/go-logger"
)

logger := glog.NewLogger(
    &glog.Config{
        LogLevel:    "debug",
        DevelopMode: true,
    },
    glog.WithRotateLog(&glog.RotateLogConfig{
        FileName:  "./logs/rotatelog-db-test.log",
        LocalTime: true,
    }),
    glog.WithLumberjack(&glog.LumberjackConfig{
        FileName:  "./logs/lumberjack-db-test.log",
        Localtime: true,
    }),
    glog.WithConsoleConfig(zapcore.EncoderConfig{
        MessageKey:     "msg",
        LevelKey:       "level",
        TimeKey:        "time",
        NameKey:        "logger",
        CallerKey:      "caller",
        FunctionKey:    "",
        StacktraceKey:  "",
        LineEnding:     zapcore.DefaultLineEnding,
        EncodeLevel:    glog.CustomCapitalStringLevelEncoder,
        EncodeTime:     glog.CustomTimeEncoder,
        EncodeDuration: zapcore.SecondsDurationEncoder,
        EncodeCaller:   zapcore.ShortCallerEncoder,
        EncodeName:     zapcore.FullNameEncoder,
    }),
	glog.WithFileConfig(zapcore.EncoderConfig{
        MessageKey:     "msg",
        LevelKey:       "level",
        TimeKey:        "time",
        NameKey:        "logger",
        CallerKey:      "caller",
        FunctionKey:    "func",
        StacktraceKey:  "stacktrace",
        LineEnding:     zapcore.DefaultLineEnding,
        EncodeLevel:    glog.CustomCapitalStringLevelEncoder,
        EncodeTime:     glog.CustomTimeEncoder,
        EncodeDuration: zapcore.SecondsDurationEncoder,
        EncodeCaller:   zapcore.ShortCallerEncoder,
        EncodeName:     zapcore.FullNameEncoder,
    }),

    glog.WithConsoleFormat(),
)

```





