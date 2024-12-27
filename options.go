package go_logger

import (
	"go.uber.org/zap/zapcore"

	"github.com/jianlu8023/go-logger/v2/internal/option"
)

const (
	stackLogLevelKey        = "stackLogLevel"
	callerKey               = "caller"
	developModeKey          = "developMode"
	moduleNameKey           = "moduleName"
	logLevelKey             = "logLevel"
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
)

// WithLogLevel 设置日志级别
// level: warn  info  debug error fatal panic
func WithLogLevel(level string) Option {
	zlv := logLevel(level)
	return option.NewOption(logLevelKey, zlv)
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

// WithStackLogLevel 设置打印 stack 的日志级别
func WithStackLogLevel(level string) Option {
	zlv := logLevel(level)
	return option.NewOption(stackLogLevelKey, zlv)
}

// WithFileOutPut 输出日志到文件
func WithFileOutPut() Option { return option.NewOption(fileOutPutKey, true) }

// WithOutConsoleOutPut 不输出日志到控制台
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

func containsOptions(options []Option, key string) (bool, Option) {
	var o Option
	exists := false
	for _, opt := range options {
		if opt.Name() == key {
			exists = true
			o = opt
			break
		}
	}
	return exists, o
}

func checkFormat(options []Option) (bool, Option) {
	var jsonExists, consoleExists, zaplogfmtExists bool

	for _, opt := range options {
		if opt.Name() == jsonFormatKey {
			jsonExists = true
		} else if opt.Name() == consoleFormatKey {
			consoleExists = true
		} else if opt.Name() == zaplogfmtKey {
			zaplogfmtExists = true
		}
	}

	if jsonExists && consoleExists && zaplogfmtExists {
		return true, WithConsoleFormat()
	} else if jsonExists && !consoleExists && !zaplogfmtExists {
		return false, WithJSONFormat()
	} else if !jsonExists && consoleExists && !zaplogfmtExists {
		return false, WithConsoleFormat()
	} else if !jsonExists && !consoleExists && zaplogfmtExists {
		return false, WithZaplogfmtFormat()
	} else {
		return false, WithConsoleFormat()
	}

}
