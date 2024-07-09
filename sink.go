package go_logger

import (
	rotateloggers "github.com/lestrrat-go/file-rotatelogs"
	"go.uber.org/zap"
	"gopkg.in/natefinch/lumberjack.v2"
)

type LumberjackSink struct {
	*lumberjack.Logger
}

func (*LumberjackSink) Sync() error {
	return nil
}

func NewLumberjack(log *lumberjack.Logger) zap.Sink {
	return &LumberjackSink{Logger: log}
}

type RotatelogSink struct {
	*rotateloggers.RotateLogs
}

func (*RotatelogSink) Sync() error {
	return nil
}

func NewRotatelog(log *rotateloggers.RotateLogs) zap.Sink {
	return &RotatelogSink{RotateLogs: log}
}
