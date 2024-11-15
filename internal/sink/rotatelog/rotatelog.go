package rotatelog

import (
	rotateloggers "github.com/lestrrat-go/file-rotatelogs"
	"go.uber.org/zap"
)

type Sink struct {
	*rotateloggers.RotateLogs
}

func (*Sink) Sync() error {
	return nil
}

func NewRotatelog(log *rotateloggers.RotateLogs) zap.Sink {
	return &Sink{RotateLogs: log}
}
