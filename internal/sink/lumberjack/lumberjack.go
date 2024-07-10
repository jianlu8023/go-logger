package lumberjack

import (
	"go.uber.org/zap"
	"gopkg.in/natefinch/lumberjack.v2"
)

type Sink struct {
	*lumberjack.Logger
}

func (*Sink) Sync() error {
	return nil
}

func NewLumberjack(log *lumberjack.Logger) zap.Sink {
	return &Sink{Logger: log}
}
