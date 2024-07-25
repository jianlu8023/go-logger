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
	zaplogfmtKey            = "zaplogfmt"
	consoleOutPutKey        = "consoleOutPut"
	fileOutPutKey           = "fileOutPut"
)

func WithFileOutPut() Option                         { return option.NewOption(fileOutPutKey, true) }
func WithConsoleOutPut() Option                      { return option.NewOption(consoleOutPutKey, true) }
func WithRotateLog(config *RotateLogConfig) Option   { return option.NewOption(rotatelogKey, config) }
func WithLumberjack(config *LumberjackConfig) Option { return option.NewOption(lumberjackKey, config) }
func WithJSONFormat() Option                         { return option.NewOption(jsonFormatKey, nil) }
func WithConsoleFormat() Option                      { return option.NewOption(consoleFormatKey, nil) }
func WithZaplogfmtFormat() Option                    { return option.NewOption(zaplogfmtKey, nil) }
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
