package go_logger

import (
	"go.uber.org/zap/zapcore"

	"github.com/jianlu8023/go-logger/internal/option"
)

const (
	lumberjackKey           = "lumberjack"
	rotatelogKey            = "rotatelog"
	consoleEncoderConfigKey = "consoleEncoderConfig"
	fileEncoderConfigKey    = "fileEncoderConfig"
	jsonFormatKey           = "jsonFormat"
	consoleFormatKey        = "consoleFormat"
)

func WithLumberjack(config *LumberjackConfig) Option {
	return option.NewOption(lumberjackKey, config)
}

func WithRotateLog(config *RotateLogConfig) Option {
	return option.NewOption(rotatelogKey, config)
}

func WithConsoleConfig(config zapcore.EncoderConfig) Option {
	return option.NewOption(consoleEncoderConfigKey, config)
}

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

func WithJSONFormat() Option {
	return option.NewOption(jsonFormatKey, nil)
}

func WithConsoleFormat() Option {
	return option.NewOption(consoleFormatKey, nil)
}
