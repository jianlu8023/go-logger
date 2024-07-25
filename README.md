# go-logger 基于 zap 结合 lumberjack rotate-logs

## Quick Start

```go
package main

import (
	glog "github.com/jianlu8023/go-logger"
)

func main() {
	logger := NewLogger(
		&Config{
			LogLevel:    "debug",
			DevelopMode: true,
			StackLevel:  "",
			ModuleName:  "[SDK]",
			Caller:      true,
		},
		WithRotateLog(&RotateLogConfig{
			FileName:  "./logs/rotatelog-logger.log",
			LocalTime: true,
		}),
		WithRotateLog(RotateLogDefaultConfig()),
		WithLumberjack(&LumberjackConfig{
			FileName: "./logs/lumberjack-logger.log",
		}),
		WithLumberjack(LumberjackDefaultConfig()),
		WithFileConfig(zapcore.EncoderConfig{
			MessageKey:     "msg",
			LevelKey:       "level",
			TimeKey:        "time",
			NameKey:        "logger",
			CallerKey:      "caller",
			StacktraceKey:  "stacktrace",
			FunctionKey:    "func",
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeLevel:    glog.CustomColorCapitalLevelEncoder,
			EncodeTime:     glog.CustomTimeEncoder,
			EncodeDuration: zapcore.SecondsDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
		}),
		WithConsoleConfig(zapcore.EncoderConfig{
			MessageKey:    "msg",
			LevelKey:      "level",
			TimeKey:       "time",
			NameKey:       "logger",
			CallerKey:     "caller",
			StacktraceKey: "stacktrace",
			FunctionKey:   "func",
			LineEnding:    zapcore.DefaultLineEnding,
			EncodeLevel:   CustomColorCapitalLevelEncoder,
			EncodeTime: func(date time.Time, encoder zapcore.PrimitiveArrayEncoder) {
				encoder.AppendString(date.Format("2006-01-02 15:04:05.000"))
			},
			EncodeDuration: zapcore.SecondsDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
		}),
		WithConsoleFormat(),
	)
}
```





