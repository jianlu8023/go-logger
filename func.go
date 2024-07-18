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

func lumberjackCore(conf *LumberjackConfig, encoder zapcore.Encoder, lv zapcore.Level) zapcore.Core {
	writeSyncer, cancel, err := zap.Open(NewLumberjackUrl(conf))
	defer func() {
		cancel()
	}()
	if err == nil {
		core := zapcore.NewCore(
			encoder,
			zapcore.NewMultiWriteSyncer(
				zapcore.AddSync(writeSyncer),
			),
			lv,
		)
		return core
	}
	return nil
}

func rotatelogCore(conf *RotateLogConfig, encoder zapcore.Encoder, lv zapcore.Level) zapcore.Core {
	writeSyncer, cancel, err := zap.Open(NewRotateLogURL(conf))
	defer func() {
		cancel()
	}()
	if err == nil {
		core := zapcore.NewCore(
			encoder,
			zapcore.NewMultiWriteSyncer(
				zapcore.AddSync(writeSyncer),
			),
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
	)
	lower := strings.ToLower(config.LogLevel)
	lv := logLevel(lower)

	// 判断options 中实有option的name是consoleEncoderConfigKey
	if ok, opt := containsOptions(options, consoleEncoderConfigKey); ok {
		consoleConfig = opt.Value().(zapcore.EncoderConfig)
	}

	if ok, opt := containsOptions(options, fileEncoderConfigKey); ok {
		fileConfig = opt.Value().(zapcore.EncoderConfig)
	}

	{
		// default console 输出
		// encoder := zaplogfmt.NewEncoder(consoleConfig)
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
			// encoder := zaplogfmt.NewEncoder(fileConfig)
			core := lumberjackCore(lumberjackConfig, encoder, lv)
			if core != nil {
				cores = append(cores, core)
			}
		case rotatelogKey:
			logConfig := option.Value().(*RotateLogConfig)
			encoder := zapcore.NewConsoleEncoder(fileConfig)
			core := rotatelogCore(logConfig, encoder, lv)
			if core != nil {
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

func jsonLogger(config *Config, options ...Option) *zap.Logger {
	var (
		cores []zapcore.Core
		// consoleConfig 默认是consoleEncoderConfig
		consoleConfig = consoleEncoderConfig
		// fileConfig 默认是fileEncoderConfig
		fileConfig = fileEncoderConfig
	)
	lower := strings.ToLower(config.LogLevel)
	lv := logLevel(lower)
	// 判断options 中实有option的name是consoleEncoderConfigKey
	if ok, opt := containsOptions(options, consoleEncoderConfigKey); ok {
		consoleConfig = opt.Value().(zapcore.EncoderConfig)
	}

	if ok, opt := containsOptions(options, fileEncoderConfigKey); ok {
		fileConfig = opt.Value().(zapcore.EncoderConfig)
	}

	{
		encoder := zapcore.NewJSONEncoder(consoleConfig)
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
			encoder := zapcore.NewJSONEncoder(fileConfig)
			core := lumberjackCore(lumberjackConfig, encoder, lv)
			if core != nil {
				cores = append(cores, core)
			}
		case rotatelogKey:
			logConfig := option.Value().(*RotateLogConfig)
			encoder := zapcore.NewJSONEncoder(fileConfig)
			core := rotatelogCore(logConfig, encoder, lv)
			if core != nil {
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

func zaplogfmtLogger(config *Config, options ...Option) *zap.Logger {
	var (
		cores []zapcore.Core
		// consoleConfig 默认是consoleEncoderConfig
		consoleConfig = consoleEncoderConfig
		// fileConfig 默认是fileEncoderConfig
		fileConfig = fileEncoderConfig
	)
	lower := strings.ToLower(config.LogLevel)
	lv := logLevel(lower)
	// 判断options 中实有option的name是consoleEncoderConfigKey
	if ok, opt := containsOptions(options, consoleEncoderConfigKey); ok {
		consoleConfig = opt.Value().(zapcore.EncoderConfig)
	}

	if ok, opt := containsOptions(options, fileEncoderConfigKey); ok {
		fileConfig = opt.Value().(zapcore.EncoderConfig)
	}

	{
		encoder := zaplogfmt.NewEncoder(consoleConfig)
		// encoder := zapcore.NewJSONEncoder(consoleConfig)
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
			// encoder := zapcore.NewJSONEncoder(fileConfig)
			encoder := zaplogfmt.NewEncoder(fileConfig)
			core := lumberjackCore(lumberjackConfig, encoder, lv)
			if core != nil {
				cores = append(cores, core)
			}
		case rotatelogKey:
			logConfig := option.Value().(*RotateLogConfig)
			// encoder := zapcore.NewJSONEncoder(fileConfig)
			encoder := zaplogfmt.NewEncoder(fileConfig)
			core := rotatelogCore(logConfig, encoder, lv)
			if core != nil {
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
