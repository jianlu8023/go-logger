package go_logger

import (
	"errors"

	"go.uber.org/zap"

	_ "github.com/jianlu8023/go-logger/internal/bootstrap"
)

const (
	warn   = "warn"
	info   = "info"
	debug  = "debug"
	_error = "error"
	fatal  = "fatal"
	_panic = "panic"
)

func NewLogger(config *Config, options ...Option) *zap.Logger {
	if both, option := checkFormat(options); both {
		panic(errors.New("logger format can not be both console and json"))
	} else {
		if option.Name() == consoleFormatKey {
			return consoleLogger(config, options...)
		} else {
			return jsonLogger(config, options...)
		}
	}
}

func NewSugaredLogger(config *Config, options ...Option) *zap.SugaredLogger {
	return NewLogger(config, options...).Sugar()
}
