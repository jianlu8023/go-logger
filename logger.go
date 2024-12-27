package go_logger

import (
	"errors"

	"go.uber.org/zap"

	_ "github.com/jianlu8023/go-logger/v2/internal/bootstrap"
)

func NewLogger(options ...Option) *zap.Logger {
	if both, option := checkFormat(options); both {
		panic(errors.New("logger format can not be both console and json"))
	} else {
		switch option.Name() {
		case zaplogfmtKey:
			return zapLogFmtLogger(options...)
		case jsonFormatKey:
			return jsonLogger(options...)
		case consoleFormatKey:
			fallthrough
		default:
			return consoleLogger(options...)
		}
	}
}

func NewSugaredLogger(options ...Option) *zap.SugaredLogger {
	return NewLogger(options...).Sugar()
}
