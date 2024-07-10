package go_logger

import (
	"github.com/jianlu8023/go-logger/internal/option"
)

const (
	lumberjackKey = "lumberjack"
)

func WithLumberjack(config *LumberjackConfig) Option {
	return option.NewOption(lumberjackKey, config)
}
