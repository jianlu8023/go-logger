package go_logger

import (
	"os"
	"strings"

	zaplogfmt "github.com/sykesm/zap-logfmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

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

func consoleCore(encoder zapcore.Encoder, lv zap.AtomicLevel) zapcore.Core {
	return zapcore.NewCore(
		encoder,
		// zapcore.NewMultiWriteSyncer(
		zapcore.AddSync(os.Stdout),
		// ),
		lv,
	)
}

func lumberjackCore(conf *LumberjackConfig, encoder zapcore.Encoder, lv zap.AtomicLevel) zapcore.Core {
	writeSyncer, cancel, err := zap.Open(NewLumberjackUrl(conf))
	defer func() {
		cancel()
	}()
	if err == nil {
		core := zapcore.NewCore(
			encoder,
			// zapcore.NewMultiWriteSyncer(
			zapcore.AddSync(writeSyncer),
			// ),
			lv,
		)
		return core
	}
	return nil
}

func rotatelogCore(conf *RotateLogConfig, encoder zapcore.Encoder, lv zap.AtomicLevel) zapcore.Core {
	writeSyncer, cancel, err := zap.Open(NewRotateLogURL(conf))
	defer func() {
		cancel()
	}()
	if err == nil {
		core := zapcore.NewCore(
			encoder,
			// zapcore.NewMultiWriteSyncer(
			zapcore.AddSync(writeSyncer),
			// ),
			lv,
		)
		return core
	}
	return nil
}

func consoleLogger(config *Config, options ...Option) *zap.Logger {
	var (
		cores []zapcore.Core
		// consoleConfig 默认是consoleEncoderConfig
		consoleConfig = consoleEncoderConfig
		// fileConfig 默认是fileEncoderConfig
		fileConfig = fileEncoderConfig
		encoder    zapcore.Encoder
	)

	lv := logLevel(strings.ToLower(config.LogLevel))
	alv := zap.NewAtomicLevel()
	alv.SetLevel(lv)

	// 判断options 中实有option的name是consoleEncoderConfigKey
	if ok, opt := containsOptions(options, consoleEncoderConfigKey); ok {
		consoleConfig = opt.Value().(zapcore.EncoderConfig)
	}

	// 判断options 中实有option的name是fileEncoderConfigKey
	if ok, opt := containsOptions(options, fileEncoderConfigKey); ok {
		fileConfig = opt.Value().(zapcore.EncoderConfig)
	}

	{
		// default console 输出
		encoder = zapcore.NewConsoleEncoder(consoleConfig)
		cores = append(cores, consoleCore(encoder, alv))
	}
	for _, option := range options {
		switch option.Name() {
		case lumberjackKey:
			lumberjackConfig := option.Value().(*LumberjackConfig)
			encoder = zapcore.NewConsoleEncoder(fileConfig)
			core := lumberjackCore(lumberjackConfig, encoder, alv)
			if core != nil {
				cores = append(cores, core)
			}
		case rotatelogKey:
			logConfig := option.Value().(*RotateLogConfig)
			encoder = zapcore.NewConsoleEncoder(fileConfig)
			core := rotatelogCore(logConfig, encoder, alv)
			if core != nil {
				cores = append(cores, core)
			}
		}
	}
	return genLogger(cores, config)
}

func jsonLogger(config *Config, options ...Option) *zap.Logger {
	var (
		cores []zapcore.Core
		// consoleConfig 默认是consoleEncoderConfig
		consoleConfig = consoleEncoderConfig
		// fileConfig 默认是fileEncoderConfig
		fileConfig = fileEncoderConfig
		encoder    zapcore.Encoder
	)

	lv := logLevel(strings.ToLower(config.LogLevel))
	alv := zap.NewAtomicLevel()
	alv.SetLevel(lv)

	// 判断options 中实有option的name是consoleEncoderConfigKey
	if ok, opt := containsOptions(options, consoleEncoderConfigKey); ok {
		consoleConfig = opt.Value().(zapcore.EncoderConfig)
	}

	// 判断options 中实有option的name是fileEncoderConfigKey
	if ok, opt := containsOptions(options, fileEncoderConfigKey); ok {
		fileConfig = opt.Value().(zapcore.EncoderConfig)
	}

	{
		encoder = zapcore.NewJSONEncoder(consoleConfig)
		cores = append(cores, consoleCore(encoder, alv))
	}
	for _, option := range options {
		switch option.Name() {
		case lumberjackKey:
			lumberjackConfig := option.Value().(*LumberjackConfig)
			encoder = zapcore.NewJSONEncoder(fileConfig)
			core := lumberjackCore(lumberjackConfig, encoder, alv)
			if core != nil {
				cores = append(cores, core)
			}
		case rotatelogKey:
			logConfig := option.Value().(*RotateLogConfig)
			encoder = zapcore.NewJSONEncoder(fileConfig)
			core := rotatelogCore(logConfig, encoder, alv)
			if core != nil {
				cores = append(cores, core)
			}
		}
	}
	return genLogger(cores, config)
}

func zaplogfmtLogger(config *Config, options ...Option) *zap.Logger {
	var (
		cores []zapcore.Core
		// consoleConfig 默认是consoleEncoderConfig
		consoleConfig = consoleEncoderConfig
		// fileConfig 默认是fileEncoderConfig
		fileConfig = fileEncoderConfig
		encoder    zapcore.Encoder
	)

	lv := logLevel(strings.ToLower(config.LogLevel))
	alv := zap.NewAtomicLevel()
	alv.SetLevel(lv)

	// 判断options 中实有option的name是consoleEncoderConfigKey
	if ok, opt := containsOptions(options, consoleEncoderConfigKey); ok {
		consoleConfig = opt.Value().(zapcore.EncoderConfig)
	}

	// 判断options 中实有option的name是fileEncoderConfigKey
	if ok, opt := containsOptions(options, fileEncoderConfigKey); ok {
		fileConfig = opt.Value().(zapcore.EncoderConfig)
	}

	{
		encoder = zaplogfmt.NewEncoder(consoleConfig)
		cores = append(cores, consoleCore(encoder, alv))
	}

	for _, option := range options {
		switch option.Name() {
		case lumberjackKey:
			lumberjackConfig := option.Value().(*LumberjackConfig)
			encoder = zaplogfmt.NewEncoder(fileConfig)
			core := lumberjackCore(lumberjackConfig, encoder, alv)
			if core != nil {
				cores = append(cores, core)
			}
		case rotatelogKey:
			logConfig := option.Value().(*RotateLogConfig)
			encoder = zaplogfmt.NewEncoder(fileConfig)
			core := rotatelogCore(logConfig, encoder, alv)
			if core != nil {
				cores = append(cores, core)
			}
		}
	}
	return genLogger(cores, config)
}

func genLogger(cores []zapcore.Core, config *Config) *zap.Logger {
	core := zapcore.NewTee(cores...)
	logger := zap.New(core)
	if len(config.ModuleName) != 0 {
		logger = logger.Named(config.ModuleName)
	}
	if config.DevelopMode {
		logger = logger.WithOptions(zap.Development())
	}
	if config.Caller {
		logger = logger.WithOptions(zap.AddCaller())
	}
	if len(config.StackLevel) != 0 {
		logger = logger.WithOptions(zap.AddStacktrace(logLevel(strings.ToLower(config.StackLevel))))
	}
	// logger = logger.WithOptions(zap.AddCallerSkip(1))
	defer func(logger *zap.Logger) {
		_ = logger.Sync()
	}(logger)
	return logger
}
