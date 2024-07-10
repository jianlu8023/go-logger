package go_logger

import (
	"os"
	"os/signal"
	"syscall"
	"testing"
	"time"

	"go.uber.org/zap/zapcore"
)

func TestNewLogger(t *testing.T) {
	logger := NewLogger(&Config{
		LogLevel:    "info",
		DevelopMode: true,
	}, WithLumberjack(&LumberjackConfig{
		FileName: "./logs/test.log",
	}), WithFileConfig(zapcore.EncoderConfig{
		MessageKey:     "msg",
		LevelKey:       "level",
		TimeKey:        "time",
		NameKey:        "logger",
		CallerKey:      "caller",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}))
	logger.Info("info log")

}
func TestNewSugaredLogger(t *testing.T) {

	logger := NewSugaredLogger(&Config{
		LogLevel:    "info",
		DevelopMode: false,
	}, WithRotateLog(&RotateLogConfig{
		FileName:  "./logs/rotatelog-test.log",
		LocalTime: true,
	}), WithLumberjack(&LumberjackConfig{
		FileName: "./logs/lumberjack-test.log",
	}))

	ticker := time.NewTicker(time.Second * 10)

	go func() {
		for {
			select {
			case <-ticker.C:
				logger.Infof("info %s", "log")
				logger.Warnf("warn %s", "log")

				logger.Infof("info struct %v", struct {
					Name string
					Age  int
				}{
					Name: "test",
					Age:  18,
				})

				// logger.Errorf("_error %s", errors.New("test _error"))
			}
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	select {
	case <-quit:
		ticker.Stop()
		logger.Info("stop")
	}
}
