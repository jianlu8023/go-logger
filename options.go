package go_logger

import (
	"go.uber.org/zap/zapcore"

	"github.com/jianlu8023/go-logger/v2/internal/option"
)

const (
	stackLogLevelKey        = "stackLogLevel"
	callerKey               = "caller"
	callerSkipKey           = "callerSkip"
	developModeKey          = "developMode"
	moduleNameKey           = "moduleName"
	defaultLogLevelKey      = "logLevel"
	consoleLogLevelKey      = "consoleLogLevel"
	fileLogLevelKey         = "fileLogLevel"
	lumberjackKey           = "lumberjack"
	rotatelogKey            = "rotatelog"
	consoleEncoderConfigKey = "consoleEncoderConfig"
	fileEncoderConfigKey    = "fileEncoderConfig"
	jsonFormatKey           = "jsonFormat"
	consoleFormatKey        = "consoleFormat"
	zaplogfmtKey            = "zaplogfmt"
	consoleOutPutKey        = "consoleOutPut"
	withoutConsoleOutPutKey = "withoutConsoleOutPut"
	fileOutPutKey           = "fileOutPut"
	withoutFileOutPutKey    = "withoutFileOutPut"
)

// WithDefaultLogLevel 设置日志级别
// level: warn  info  debug error fatal panic
func WithDefaultLogLevel(level string) Option {
	zlv := logLevel(level)
	return option.NewOption(defaultLogLevelKey, zlv)
}

// WithConsoleLogLevel 设置控制台日志级别
// level: warn  info  debug error fatal panic
func WithConsoleLogLevel(level string) Option {
	zlv := logLevel(level)
	return option.NewOption(consoleLogLevelKey, zlv)
}

// WithFileLogLevel 设置文件日志级别
// level: warn  info  debug error fatal panic
func WithFileLogLevel(level string) Option {
	zlv := logLevel(level)
	return option.NewOption(fileLogLevelKey, zlv)
}

// WithStackLogLevel 设置打印 stack 的日志级别
func WithStackLogLevel(level string) Option {
	zlv := logLevel(level)
	return option.NewOption(stackLogLevelKey, zlv)
}

// WithModuleName 设置moduleName
func WithModuleName(name string) Option {
	return option.NewOption(moduleNameKey, name)
}

// WithDevelopMode 开启开发模式
func WithDevelopMode() Option {
	return option.NewOption(developModeKey, true)
}

// WithCaller 开启 caller
func WithCaller() Option {
	return option.NewOption(callerKey, true)
}

// WithCallerSkip 设置 caller skip
func WithCallerSkip(skip int) Option {
	return option.NewOption(callerSkipKey, skip)
}

// WithFileOutPut 输出日志到文件
func WithFileOutPut() Option { return option.NewOption(fileOutPutKey, true) }

// WithoutFileOutPut 不输出日志到文件（优先级高于 WithFileOutPut）
func WithoutFileOutPut() Option { return option.NewOption(withoutFileOutPutKey, true) }

// WithOutConsoleOutPut 不输出日志到控制台（优先级高于 WithConsoleOutPut）
func WithOutConsoleOutPut() Option { return option.NewOption(withoutConsoleOutPutKey, true) }

// WithConsoleOutPut 输出日志到控制台
func WithConsoleOutPut() Option { return option.NewOption(consoleOutPutKey, true) }

// WithRotateLog 使用 rotatelog 进行日志切割
func WithRotateLog(config *RotateLogConfig) Option { return option.NewOption(rotatelogKey, config) }

// WithLumberjack 使用 lumberjack 进行日志切割
func WithLumberjack(config *LumberjackConfig) Option { return option.NewOption(lumberjackKey, config) }

// WithJSONFormat 使用 json 格式化日志
func WithJSONFormat() Option { return option.NewOption(jsonFormatKey, nil) }

// WithConsoleFormat 使用 console 格式化日志
func WithConsoleFormat() Option { return option.NewOption(consoleFormatKey, nil) }

// WithZaplogfmtFormat 使用 zaplogfmt 格式化日志
func WithZaplogfmtFormat() Option { return option.NewOption(zaplogfmtKey, nil) }

// WithConsoleConfig 设置控制台日志格式
func WithConsoleConfig(config zapcore.EncoderConfig) Option {
	return option.NewOption(consoleEncoderConfigKey, config)
}

// WithFileConfig 设置文件日志格式
func WithFileConfig(config zapcore.EncoderConfig) Option {
	return option.NewOption(fileEncoderConfigKey, config)
}

// buildOptionMap 将 options 转为 map，相同 Name 的 option 后者覆盖前者。
//
// 语义说明：
//   - 此 map 仅用于"单值配置"的快速查找（如 defaultLogLevelKey、moduleNameKey、
//     consoleOutPutKey、fileOutPutKey、formatKey 等），重复时只生效最后一个。
//   - 对于"可多实例配置"（lumberjackKey、rotatelogKey），不应从此 map 读取，
//     而应直接遍历原始 options slice，以支持多文件输出场景。
//     参见 func.go buildLogger 中对 lumberjack/rotatelog 的处理。
func buildOptionMap(options []Option) map[string]Option {
	optMap := make(map[string]Option, len(options))
	for _, opt := range options {
		optMap[opt.Name()] = opt
	}
	return optMap
}

func checkFormat(options []Option) (bool, Option) {
	var formats []Option
	for _, opt := range options {
		switch opt.Name() {
		case jsonFormatKey, consoleFormatKey, zaplogfmtKey:
			formats = append(formats, opt)
		}
	}

	switch len(formats) {
	case 0:
		return false, WithConsoleFormat()
	case 1:
		return false, formats[0]
	default:
		return true, WithConsoleFormat()
	}
}

func hasOutput(options []Option) bool {
	optMap := buildOptionMap(options)

	// console 输出：默认有效，仅当 withoutConsoleOutPutKey 存在时无效
	consoleEffective := true
	if _, without := optMap[withoutConsoleOutPutKey]; without {
		consoleEffective = false
	}

	// file 输出：仅当 fileOutPutKey 存在 且 withoutFileOutPutKey 不存在时有效
	fileEffective := false
	if _, without := optMap[withoutFileOutPutKey]; !without {
		if _, with := optMap[fileOutPutKey]; with {
			fileEffective = true
		}
	}

	return consoleEffective || fileEffective
}

func hasLogLevel(options []Option) bool {
	optMap := buildOptionMap(options)
	if _, ok1 := optMap[defaultLogLevelKey]; !ok1 {
		return false
	}
	return true
}
