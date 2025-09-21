package go_logger

import (
	"os"
	"slices"

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

func consoleLogger(options ...Option) *zap.Logger {
	var (
		cores []zapcore.Core

		// consoleConfig 默认是 consoleEncoderConfig
		consoleConfig = consoleEncoderConfig

		// fileConfig 默认是 fileEncoderConfig
		fileConfig = fileEncoderConfig

		encoder zapcore.Encoder
	)

	if len(options) == 0 {
		// fmt.Println("create logger with debug level and console output")

		options = append(options, WithConsoleOutPut())
		options = append(options, WithDefaultLogLevel(debug))
	} else {
		// no log level selected use debug level
		if ok, _ := containsOptions(options, defaultLogLevelKey); !ok {
			// fmt.Println("no log level selected use debug level")
			options = append(options, WithDefaultLogLevel(debug))
		}
		// without console output selected, delete console output option
		if ok, _ := containsOptions(options, withoutConsoleOutPutKey); ok {
			// fmt.Println("without console output selected")
			options = slices.DeleteFunc(options, func(option Option) bool {
				return option.Name() == consoleOutPutKey
			})
			if ok, _ := containsOptions(options, fileOutPutKey); !ok {
				panic("WithOutConsoleOutPut is set, but no output selected")
			}
		} else {
			options = append(options, WithConsoleOutPut())
		}
	}

	// alv := zap.NewAtomicLevel()
	// alv.SetLevel(getLogLevel(options))

	{
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
			panic("WithFileOutPut set, but no output file config")
		}
	} else {
		// 如果没有 WithFileOutPut，则判断是否有
		if ok, _ := containsOptions(options, fileEncoderConfigKey); ok {
			panic("WithFileOutPut not set, but WithFileEncoderConfig set")
		} else if ok, _ = containsOptions(options, lumberjackKey); ok {
			panic("WithFileOutPut not set, but WithLumberjack set")
		} else if ok, _ = containsOptions(options, rotatelogKey); ok {
			panic("WithFileOutPut not set, but WithRotateLog set")
		}
	}

	{
		if ok, _ := containsOptions(options, consoleOutPutKey); ok {
			// default console 输出
			encoder = zapcore.NewConsoleEncoder(consoleConfig)
			consoleLv := zap.NewAtomicLevel()
			consoleLv.SetLevel(getConsoleLogLevel(options))
			cores = append(cores, consoleCore(encoder, consoleLv))
		}
	}

	var (
		moduleName                  = ""
		developMode                 = false
		callerMode                  = false
		stackLogLevel zapcore.Level = -2
		callerSkip                  = 0
	)
	fileLv := zap.NewAtomicLevel()
	fileLv.SetLevel(getFileLogLevel(options))
	for _, option := range options {
		switch option.Name() {
		case lumberjackKey:
			lumberjackConfig := option.Value().(*LumberjackConfig)
			encoder = zapcore.NewConsoleEncoder(fileConfig)
			core := lumberjackCore(lumberjackConfig, encoder, fileLv)
			if core != nil {
				cores = append(cores, core)
			}
		case rotatelogKey:
			logConfig := option.Value().(*RotateLogConfig)
			encoder = zapcore.NewConsoleEncoder(fileConfig)
			core := rotateLogCore(logConfig, encoder, fileLv)
			if core != nil {
				cores = append(cores, core)
			}
		case moduleNameKey:
			moduleName = option.Value().(string)
		case developModeKey:
			developMode = option.Value().(bool)
		case callerKey:
			callerMode = option.Value().(bool)
		case callerSkipKey:
			callerSkip = option.Value().(int)
		case stackLogLevelKey:
			stackLogLevel = option.Value().(zapcore.Level)
		}
	}
	return genLogger(cores, moduleName, developMode, callerMode, callerSkip, stackLogLevel)
}

func jsonLogger(options ...Option) *zap.Logger {
	var (
		cores []zapcore.Core
		// consoleConfig 默认是consoleEncoderConfig
		consoleConfig = consoleEncoderConfig
		// fileConfig 默认是fileEncoderConfig
		fileConfig = fileEncoderConfig
		encoder    zapcore.Encoder
	)
	if len(options) == 0 {
		// fmt.Println("no options selected")
		// fmt.Println("create logger with debug level and console output")

		options = append(options, WithConsoleOutPut())
		options = append(options, WithDefaultLogLevel(debug))
	} else {
		// no log level selected use debug level
		if ok, _ := containsOptions(options, defaultLogLevelKey); !ok {
			// fmt.Println("no log level selected use debug level")
			options = append(options, WithDefaultLogLevel(debug))
		}
		// without console output selected, delete console output option
		if ok, _ := containsOptions(options, withoutConsoleOutPutKey); ok {
			// fmt.Println("without console output selected")
			options = slices.DeleteFunc(options, func(option Option) bool {
				return option.Name() == consoleOutPutKey
			})
			if ok, _ := containsOptions(options, fileOutPutKey); !ok {
				panic("WithOutConsoleOutPut is set, but no output selected")
			}
		} else {
			options = append(options, WithConsoleOutPut())
		}
	}

	// alv := zap.NewAtomicLevel()
	// alv.SetLevel(getLogLevel(options))

	{
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
			consoleLv := zap.NewAtomicLevel()
			consoleLv.SetLevel(getConsoleLogLevel(options))
			cores = append(cores, consoleCore(encoder, consoleLv))
		}
	}

	var (
		moduleName                  = ""
		developMode                 = false
		callerMode                  = false
		stackLogLevel zapcore.Level = -2
		callerSkip                  = 0
	)
	fileLv := zap.NewAtomicLevel()
	fileLv.SetLevel(getFileLogLevel(options))
	for _, option := range options {
		switch option.Name() {
		case lumberjackKey:
			lumberjackConfig := option.Value().(*LumberjackConfig)
			encoder = zapcore.NewJSONEncoder(fileConfig)
			core := lumberjackCore(lumberjackConfig, encoder, fileLv)
			if core != nil {
				cores = append(cores, core)
			}
		case rotatelogKey:
			logConfig := option.Value().(*RotateLogConfig)
			encoder = zapcore.NewJSONEncoder(fileConfig)
			core := rotateLogCore(logConfig, encoder, fileLv)
			if core != nil {
				cores = append(cores, core)
			}
		case moduleNameKey:
			moduleName = option.Value().(string)
		case developModeKey:
			developMode = option.Value().(bool)
		case callerKey:
			callerMode = option.Value().(bool)
		case callerSkipKey:
			callerSkip = option.Value().(int)
		case stackLogLevelKey:
			stackLogLevel = option.Value().(zapcore.Level)
		}
	}
	return genLogger(cores, moduleName, developMode, callerMode, callerSkip, stackLogLevel)
}

func zapLogFmtLogger(options ...Option) *zap.Logger {
	var (
		cores []zapcore.Core
		// consoleConfig 默认是consoleEncoderConfig
		consoleConfig = consoleEncoderConfig
		// fileConfig 默认是fileEncoderConfig
		fileConfig = fileEncoderConfig
		encoder    zapcore.Encoder
	)
	if len(options) == 0 {
		// fmt.Println("no options selected")
		// fmt.Println("create logger with debug level and console output")

		options = append(options, WithConsoleOutPut())
		options = append(options, WithDefaultLogLevel(debug))
	} else {
		// no log level selected use debug level
		if ok, _ := containsOptions(options, defaultLogLevelKey); !ok {
			// fmt.Println("no log level selected use debug level")
			options = append(options, WithDefaultLogLevel(debug))
		}
		// without console output selected, delete console output option
		if ok, _ := containsOptions(options, withoutConsoleOutPutKey); ok {
			// fmt.Println("without console output selected")
			options = slices.DeleteFunc(options, func(option Option) bool {
				return option.Name() == consoleOutPutKey
			})
			if ok, _ := containsOptions(options, fileOutPutKey); !ok {
				panic("WithOutConsoleOutPut is set, but no output selected")
			}
		} else {
			options = append(options, WithConsoleOutPut())
		}
	}

	// alv := zap.NewAtomicLevel()
	// alv.SetLevel(getLogLevel(options))

	{
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
			consoleLv := zap.NewAtomicLevel()
			consoleLv.SetLevel(getConsoleLogLevel(options))
			cores = append(cores, consoleCore(encoder, consoleLv))
		}
	}
	var (
		moduleName                  = ""
		developMode                 = false
		callerMode                  = false
		stackLogLevel zapcore.Level = -2
		callerSkip                  = 0
	)
	fileLv := zap.NewAtomicLevel()
	fileLv.SetLevel(getFileLogLevel(options))
	for _, option := range options {
		switch option.Name() {
		case lumberjackKey:
			lumberjackConfig := option.Value().(*LumberjackConfig)
			encoder = zaplogfmt.NewEncoder(fileConfig)
			core := lumberjackCore(lumberjackConfig, encoder, fileLv)
			if core != nil {
				cores = append(cores, core)
			}
		case rotatelogKey:
			logConfig := option.Value().(*RotateLogConfig)
			encoder = zaplogfmt.NewEncoder(fileConfig)
			core := rotateLogCore(logConfig, encoder, fileLv)
			if core != nil {
				cores = append(cores, core)
			}
		case moduleNameKey:
			moduleName = option.Value().(string)
		case developModeKey:
			developMode = option.Value().(bool)
		case callerKey:
			callerMode = option.Value().(bool)
		case callerSkipKey:
			callerSkip = option.Value().(int)
		case stackLogLevelKey:
			stackLogLevel = option.Value().(zapcore.Level)
		}
	}
	return genLogger(cores, moduleName, developMode, callerMode, callerSkip, stackLogLevel)
}

func genLogger(cores []zapcore.Core, moduleName string, developMode, callerMode bool, callerSkip int, stackLogLevel zapcore.Level) *zap.Logger {
	core := zapcore.NewTee(cores...)
	logger := zap.New(core)
	if len(moduleName) != 0 {
		logger = logger.Named(moduleName)
	}
	if developMode {
		logger = logger.WithOptions(zap.Development())
	}
	if callerMode {
		logger = logger.WithOptions(zap.AddCaller())
		logger = logger.WithOptions(zap.AddCallerSkip(callerSkip))
	}
	if stackLogLevel > -2 {
		logger = logger.WithOptions(zap.AddStacktrace(stackLogLevel))
	}
	defer func(logger *zap.Logger) {
		_ = logger.Sync()
	}(logger)
	return logger
}
