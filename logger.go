package go_logger

import (
	"fmt"
	"os"

	rotateloggers "github.com/lestrrat-go/file-rotatelogs"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"

	_ "github.com/jianlu8023/go-logger/internal/bootstrap"
)

type Logger struct {
	ZapLogger        *zap.Logger
	LumberjackLogger *lumberjack.Logger
	RotateLogger     *rotateloggers.RotateLogs
}

func NewLogger(config *Config, options ...Option) *zap.Logger {
	mode := config.Mode
	var cores []zapcore.Core
	var lv zapcore.Level

	switch config.LogLevel {
	case "info":
		lv = zapcore.InfoLevel
	case "debug":
		lv = zapcore.DebugLevel
	case "warn":
		lv = zapcore.WarnLevel
	case "error":
		lv = zapcore.ErrorLevel
	case "panic":
		lv = zapcore.PanicLevel
	case "fatal":
		lv = zapcore.FatalLevel
	default:
		lv = zapcore.InfoLevel
	}
	for _, m := range mode {
		switch m {
		case "stdout":
			encoder := zapcore.NewConsoleEncoder(consoleEncoderConfig)
			core := zapcore.NewCore(encoder,
				zapcore.NewMultiWriteSyncer(
					zapcore.AddSync(os.Stdout)),
				lv,
			)
			cores = append(cores, core)
		case "file":
			encoder := zapcore.NewConsoleEncoder(fileEncoderConfig)
			writeSyncer, cancel, err := zap.Open(NewLumberjackUrl(nil))
			defer func() {
				fmt.Println("file cancel")
				cancel()
			}()
			if err == nil {
				core := zapcore.NewCore(encoder, zapcore.NewMultiWriteSyncer(
					zapcore.AddSync(writeSyncer),
				),
					lv,
				)
				cores = append(cores, core)
			}
		case "date":
			encoder := zapcore.NewConsoleEncoder(fileEncoderConfig)
			writeSyncer, cancel, err := zap.Open(NewRotateLogURL(&RotateLogConfig{
				LocalTime: true,
			}))

			defer func() {
				fmt.Println("date cancel")
				cancel()
			}()
			if err == nil {
				core := zapcore.NewCore(encoder, zapcore.NewMultiWriteSyncer(
					zapcore.AddSync(writeSyncer),
				),
					lv,
				)
				cores = append(cores, core)
			}
		}
	}

	core := zapcore.NewTee(cores...)
	if config.DevelopMode {
		return zap.New(core, zap.AddCaller(), zap.Development(), zap.AddStacktrace(zapcore.ErrorLevel))
	} else {
		return zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))
	}
}

func NewSugaredLogger(config *Config, options ...Option) *zap.SugaredLogger {

	return NewLogger(config, options...).Sugar()
}
