package go_logger

import (
	"fmt"
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
	if cancel != nil {
		defer cancel()
	}
	if err != nil {
		// zap.Open 失败时不 panic，输出 stderr 告警，返回 nil core 由 genLogger 过滤
		fmt.Fprintf(os.Stderr, "[go-logger] WARNING: lumberjack sink open failed: %v\n", err)
		return nil
	}
	return zapcore.NewCore(
		encoder,
		zapcore.NewMultiWriteSyncer(
			zapcore.AddSync(writeSyncer),
		),
		lv,
	)
}

func rotateLogCore(conf *RotateLogConfig, encoder zapcore.Encoder, lv zap.AtomicLevel) zapcore.Core {
	writeSyncer, cancel, err := zap.Open(NewRotateLogURL(conf))
	if cancel != nil {
		defer cancel()
	}
	if err != nil {
		// zap.Open 失败时不 panic，输出 stderr 告警，返回 nil core 由 genLogger 过滤
		fmt.Fprintf(os.Stderr, "[go-logger] WARNING: rotatelog sink open failed: %v\n", err)
		return nil
	}
	return zapcore.NewCore(
		encoder,
		zapcore.NewMultiWriteSyncer(
			zapcore.AddSync(writeSyncer),
		),
		lv,
	)
}

func consoleLogger(options ...Option) *zap.Logger {
	return buildLogger(options, zapcore.NewConsoleEncoder)
}

func jsonLogger(options ...Option) *zap.Logger {
	return buildLogger(options, zapcore.NewJSONEncoder)
}

func zapLogFmtLogger(options ...Option) *zap.Logger {
	return buildLogger(options, func(cfg zapcore.EncoderConfig) zapcore.Encoder {
		return zaplogfmt.NewEncoder(cfg)
	})
}

// buildLogger is the common core for console, json, and zaplogfmt loggers.
// encoderFactory creates the encoder based on the given EncoderConfig.
func buildLogger(options []Option, encoderFactory func(zapcore.EncoderConfig) zapcore.Encoder) *zap.Logger {
	var (
		cores []zapcore.Core

		// consoleConfig 默认是 consoleEncoderConfig
		consoleConfig = consoleEncoderConfig

		// fileConfig 默认是 fileEncoderConfig
		fileConfig = fileEncoderConfig

		encoder zapcore.Encoder
	)

	// Ensure default options when none provided
	if len(options) == 0 {
		options = append(options, WithConsoleOutPut())
		options = append(options, WithDefaultLogLevel(debugLevel))
	}

	if !hasOutput(options) {
		options = append(options, WithConsoleOutPut())
	}

	if !hasLogLevel(options) {
		options = append(options, WithDefaultLogLevel(debugLevel))
	}

	optMap := buildOptionMap(options)

	if len(options) != 0 {
		if _, ok := optMap[defaultLogLevelKey]; !ok {
			defaultLevelOpt := WithDefaultLogLevel(debugLevel)
			options = append(options, defaultLevelOpt)
			optMap[defaultLogLevelKey] = defaultLevelOpt
		}
	}

	// 判断 options 中是否有 consoleEncoderConfigKey
	if opt, ok := optMap[consoleEncoderConfigKey]; ok {
		consoleConfig = opt.Value().(zapcore.EncoderConfig)
	}

	// 判断 options 中是否有 fileEncoderConfigKey
	if opt, ok := optMap[fileEncoderConfigKey]; ok {
		fileConfig = opt.Value().(zapcore.EncoderConfig)
	}

	// 判断是否有 console 输出
	// 当 withoutConsoleOutPutKey 存在时，强制不输出 console（优先级高于 consoleOutPutKey）
	if _, without := optMap[withoutConsoleOutPutKey]; !without {
		if _, ok := optMap[consoleOutPutKey]; ok {
			encoder = encoderFactory(consoleConfig)
			consoleLv := zap.NewAtomicLevel()
			consoleLv.SetLevel(getConsoleLogLevel(optMap))
			cores = append(cores, consoleCore(encoder, consoleLv))
		}
	}

	// 判断是否有 file 输出
	if _, ok := optMap[fileOutPutKey]; ok {
		fileLv := zap.NewAtomicLevel()
		fileLv.SetLevel(getFileLogLevel(optMap))

		// fileCoreCount 用于统计实际创建的文件 core 数量。
		// 若用户开启了 WithFileOutPut 但未配置任何 WithLumberjack/WithRotateLog，
		// fileCoreCount 为 0，文件输出会被静默跳过 —— 此时需要 stderr 告警。
		fileCoreCount := 0

		// 遍历原始 options slice（而非 optMap），以支持多个 lumberjack/rotatelog 实例。
		// 这是预期行为：用户可同时传入多个 WithLumberjack/WithRotateLog，
		// 每个都会创建独立的 core，写入不同的文件。
		// 参见 options.go buildOptionMap 的语义说明。
		for _, option := range options {
			switch option.Name() {
			case lumberjackKey:
				lumberjackConfig := option.Value().(*LumberjackConfig)
				encoder = encoderFactory(fileConfig)
				core := lumberjackCore(lumberjackConfig, encoder, fileLv)
				if core != nil {
					cores = append(cores, core)
					fileCoreCount++
				}
			case rotatelogKey:
				logConfig := option.Value().(*RotateLogConfig)
				encoder = encoderFactory(fileConfig)
				core := rotateLogCore(logConfig, encoder, fileLv)
				if core != nil {
					cores = append(cores, core)
					fileCoreCount++
				}
			}
		}

		// 与问题2 设计原则一致：不 panic、不静默吞掉错误、stderr 告警让开发者可感知
		if fileCoreCount == 0 {
			fmt.Fprintln(os.Stderr, "[go-logger] WARNING: WithFileOutPut() is set but no valid WithLumberjack/WithRotateLog configured, file output will be skipped")
		}
	}

	// 提取其他配置参数
	var (
		moduleName                  = ""
		developMode                 = false
		callerMode                  = false
		stackLogLevel zapcore.Level = -2
		callerSkip                  = 0
	)
	if opt, ok := optMap[moduleNameKey]; ok {
		moduleName = opt.Value().(string)
	}
	if opt, ok := optMap[developModeKey]; ok {
		developMode = opt.Value().(bool)
	}
	if opt, ok := optMap[callerKey]; ok {
		callerMode = opt.Value().(bool)
	}
	if opt, ok := optMap[callerSkipKey]; ok {
		callerSkip = opt.Value().(int)
	}
	if opt, ok := optMap[stackLogLevelKey]; ok {
		stackLogLevel = opt.Value().(zapcore.Level)
	}

	return genLogger(cores, moduleName, developMode, callerMode, callerSkip, stackLogLevel)
}
func genLogger(cores []zapcore.Core, moduleName string, developMode, callerMode bool, callerSkip int, stackLogLevel zapcore.Level) *zap.Logger {
	// Filter out nil cores — lumberjackCore/rotateLogCore may return nil on failure
	cores = slices.DeleteFunc(cores, func(c zapcore.Core) bool { return c == nil })

	// 与问题2 设计原则一致：不 panic、不静默吞掉错误、stderr 告警让开发者可感知。
	// 当所有 core 都被过滤掉（如 console 被禁用 + 文件 sink 全部打开失败）时，
	// zapcore.NewTee() 会创建一个 no-op core，所有日志将被静默丢弃。
	if len(cores) == 0 {
		fmt.Fprintln(os.Stderr, "[go-logger] WARNING: no valid log output core configured, all logs will be discarded")
	}

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
	// 不在此处 defer logger.Sync()：Sync 应在 logger 生命周期结束时由用户调用，
	// 创建时 Sync 是无意义的空操作。
	return logger
}
