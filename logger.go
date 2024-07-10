package go_logger

import (
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

const (
	warn   = "warn"
	info   = "info"
	debug  = "debug"
	_error = "error"
	fatal  = "fatal"
	_panic = "panic"
)

func NewLogger(config *Config, options ...Option) *zap.Logger {

	var (
		cores []zapcore.Core
		// consoleConfig 默认是consoleEncoderConfig
		consoleConfig = consoleEncoderConfig
		// fileConfig 默认是fileEncoderConfig
		fileConfig = fileEncoderConfig
	)

	lv := logLevel(config.LogLevel)

	// 判断options 中实有option的name是consoleEncoderConfigKey
	if ok, opt := containsOptions(options, consoleEncoderConfigKey); ok {
		consoleConfig = opt.Value().(zapcore.EncoderConfig)

	}

	if ok, opt := containsOptions(options, fileEncoderConfigKey); ok {
		fileConfig = opt.Value().(zapcore.EncoderConfig)
	}

	{
		// default console 输出
		encoder := zapcore.NewConsoleEncoder(consoleConfig)
		core := zapcore.NewCore(encoder,
			zapcore.NewMultiWriteSyncer(
				zapcore.AddSync(os.Stdout)),
			lv,
		)
		cores = append(cores, core)
	}
	for _, option := range options {
		switch option.Name() {
		case lumberjackKey:
			lumberjackConfig := option.Value().(*LumberjackConfig)
			encoder := zapcore.NewConsoleEncoder(fileConfig)
			writeSyncer, cancel, err := zap.Open(NewLumberjackUrl(lumberjackConfig))
			defer func() {
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
		case rotatelogKey:
			logConfig := option.Value().(*RotateLogConfig)
			encoder := zapcore.NewConsoleEncoder(fileConfig)
			writeSyncer, cancel, err := zap.Open(NewRotateLogURL(logConfig))
			defer func() {
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

func logLevel(level string) zapcore.Level {
	switch level {
	case info:
		return zapcore.InfoLevel
	case debug:
		return zapcore.DebugLevel
	case warn:
		return zapcore.WarnLevel
	case _error:
		return zapcore.ErrorLevel
	case _panic:
		return zapcore.PanicLevel
	case fatal:
		return zapcore.FatalLevel
	default:
		return zapcore.InfoLevel
	}
}
