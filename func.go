package go_logger

import (
	"os"
	"strings"

	zaplogfmt "github.com/sykesm/zap-logfmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func consoleCore(encoder zapcore.Encoder, lv zap.AtomicLevel) zapcore.Core {
	return zapcore.NewCore(
		encoder,
		zapcore.NewMultiWriteSyncer(
			zapcore.AddSync(os.Stdout),
		),
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
			zapcore.NewMultiWriteSyncer(
				zapcore.AddSync(writeSyncer),
			),
			lv,
		)
		return core
	}
	return nil
}

func rotateLogCore(conf *RotateLogConfig, encoder zapcore.Encoder, lv zap.AtomicLevel) zapcore.Core {
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

		encoder zapcore.Encoder
	)

	lv := logLevel(strings.ToLower(config.LogLevel))
	alv := zap.NewAtomicLevel()
	alv.SetLevel(lv)

	{
		// 默认带 WithConsoleOutPut
		if ok, _ := containsOptions(options, consoleOutPutKey); !ok {
			options = append(options, WithConsoleOutPut())
		}
		// 判断options 中实有option的name是consoleEncoderConfigKey
		if ok, opt := containsOptions(options, consoleEncoderConfigKey); ok {
			consoleConfig = opt.Value().(zapcore.EncoderConfig)
		}
	}

	// 判断是否有 WithFileOutPut
	if ok, _ := containsOptions(options, fileOutPutKey); ok {
		// 判断options 中实有option的name是fileEncoderConfigKey
		if ok, opt := containsOptions(options, fileEncoderConfigKey); ok {
			fileConfig = opt.Value().(zapcore.EncoderConfig)
		}

		var lumberjack, rotateLog bool
		if ok, _ := containsOptions(options, lumberjackKey); !ok {
			lumberjack = true
		}
		if ok, _ := containsOptions(options, rotatelogKey); !ok {
			rotateLog = true
		}
		if lumberjack && rotateLog {
			panic("WithFileOutPut is set, but no output file config")
		}
	} else {
		// 如果没有 WithFileOutPut，则判断是否有
		if ok, _ := containsOptions(options, fileEncoderConfigKey); ok {
			panic("WithFileOutPut is not set, but WithFileEncoderConfig is set")
		} else if ok, _ = containsOptions(options, lumberjackKey); ok {
			panic("WithFileOutPut is not set, but WithLumberjack is set")
		} else if ok, _ = containsOptions(options, rotatelogKey); ok {
			panic("WithFileOutPut is not set, but WithRotateLog is set")
		}
	}

	{
		if ok, _ := containsOptions(options, consoleOutPutKey); ok {
			// default console 输出
			encoder = zapcore.NewConsoleEncoder(consoleConfig)
			cores = append(cores, consoleCore(encoder, alv))
		}
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
			core := rotateLogCore(logConfig, encoder, alv)
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

	{
		// 默认带 WithConsoleOutPut
		if ok, _ := containsOptions(options, consoleOutPutKey); !ok {
			options = append(options, WithConsoleOutPut())
		}
		// 判断options 中实有option的name是consoleEncoderConfigKey
		if ok, opt := containsOptions(options, consoleEncoderConfigKey); ok {
			consoleConfig = opt.Value().(zapcore.EncoderConfig)
		}
	}

	// 判断是否有 WithFileOutPut
	if ok, _ := containsOptions(options, fileOutPutKey); ok {
		// 判断options 中实有option的name是fileEncoderConfigKey
		if ok, opt := containsOptions(options, fileEncoderConfigKey); ok {
			fileConfig = opt.Value().(zapcore.EncoderConfig)
		}
		var lumberjack, rotateLog bool
		if ok, _ := containsOptions(options, lumberjackKey); !ok {
			lumberjack = true
		}
		if ok, _ := containsOptions(options, rotatelogKey); !ok {
			rotateLog = true
		}
		if lumberjack && rotateLog {
			panic("WithFileOutPut is set, but no output file config")
		}
	} else {
		// 如果没有 WithFileOutPut，则判断是否有
		if ok, _ := containsOptions(options, fileEncoderConfigKey); ok {
			panic("WithFileOutPut is not set, but WithFileEncoderConfig is set")
		} else if ok, _ = containsOptions(options, lumberjackKey); ok {
			panic("WithFileOutPut is not set, but WithLumberjack is set")
		} else if ok, _ = containsOptions(options, rotatelogKey); ok {
			panic("WithFileOutPut is not set, but WithRotateLog is set")
		}
	}

	{
		if ok, _ := containsOptions(options, consoleOutPutKey); ok {
			// default console 输出
			encoder = zapcore.NewJSONEncoder(consoleConfig)
			cores = append(cores, consoleCore(encoder, alv))
		}
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
			core := rotateLogCore(logConfig, encoder, alv)
			if core != nil {
				cores = append(cores, core)
			}
		}
	}
	return genLogger(cores, config)
}

func zapLogFmtLogger(config *Config, options ...Option) *zap.Logger {
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

	{
		// 默认带 WithConsoleOutPut
		if ok, _ := containsOptions(options, consoleOutPutKey); !ok {
			options = append(options, WithConsoleOutPut())
		}

		// 判断options 中实有option的name是consoleEncoderConfigKey
		if ok, opt := containsOptions(options, consoleEncoderConfigKey); ok {
			consoleConfig = opt.Value().(zapcore.EncoderConfig)
		}
	}

	// 判断是否有 WithFileOutPut
	if ok, _ := containsOptions(options, fileOutPutKey); ok {
		// 判断options 中实有option的name是fileEncoderConfigKey
		if ok, opt := containsOptions(options, fileEncoderConfigKey); ok {
			fileConfig = opt.Value().(zapcore.EncoderConfig)
		}
		var lumberjack, rotateLog bool
		if ok, _ := containsOptions(options, lumberjackKey); !ok {
			lumberjack = true
		}
		if ok, _ := containsOptions(options, rotatelogKey); !ok {
			rotateLog = true
		}
		if lumberjack && rotateLog {
			panic("WithFileOutPut is set, but no output file config")
		}
	} else {
		// 如果没有 WithFileOutPut，则判断是否有
		if ok, _ := containsOptions(options, fileEncoderConfigKey); ok {
			panic("WithFileOutPut is not set, but WithFileEncoderConfig is set")
		} else if ok, _ = containsOptions(options, lumberjackKey); ok {
			panic("WithFileOutPut is not set, but WithLumberjack is set")
		} else if ok, _ = containsOptions(options, rotatelogKey); ok {
			panic("WithFileOutPut is not set, but WithRotateLog is set")
		}
	}

	{

		if ok, _ := containsOptions(options, consoleOutPutKey); ok {
			// default console 输出
			encoder = zaplogfmt.NewEncoder(consoleConfig)
			cores = append(cores, consoleCore(encoder, alv))
		}
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
			core := rotateLogCore(logConfig, encoder, alv)
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
