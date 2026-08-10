package db_logger

import (
	"go.uber.org/zap"
)

const (
	customLoggerKey = "customLogger"
)

type option struct {
	name  string
	value interface{}
}

func (o *option) Name() string       { return o.name }
func (o *option) Value() interface{} { return o.value }

func WithCustomLogger(log *zap.Logger) Option {
	return &option{name: customLoggerKey, value: log}
}
