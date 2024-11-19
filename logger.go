package go_logger

import (
	"errors"

	"go.uber.org/zap"

	_ "github.com/jianlu8023/go-logger/internal/bootstrap"
)

func NewLogger(config *Config, options ...Option) *zap.Logger {
	if both, option := checkFormat(options); both {
		panic(errors.New("logger format can not be both console and json"))
	} else {
		switch option.Name() {
		case zaplogfmtKey:
			return zapLogFmtLogger(config, options...)
		case jsonFormatKey:
			return jsonLogger(config, options...)
		case consoleFormatKey:
			fallthrough
		default:
			return consoleLogger(config, options...)
		}
	}
}

func NewSugaredLogger(config *Config, options ...Option) *zap.SugaredLogger {
	return NewLogger(config, options...).Sugar()
}
