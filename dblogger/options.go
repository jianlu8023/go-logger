package dblogger

import (
	"go.uber.org/zap"

	"github.com/jianlu8023/go-logger/internal/option"
)

const (
	customLoggerKey = "customLogger"
)

func WithCustomLogger(log *zap.Logger) Option {
	return option.NewOption(customLoggerKey, log)
}
