package go_logger

import (
	"os"
	"os/signal"
	"syscall"
	"testing"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func TestNewLogger(t *testing.T) {
	logger := NewLogger(
		&Config{
			LogLevel:    "info",
			DevelopMode: true,
		},
		WithRotateLog(&RotateLogConfig{
			FileName:  "./logs/rotatelog-db-test.log",
			LocalTime: true,
		}),
		WithRotateLog(RotateLogDefaultConfig()),
		WithLumberjack(&LumberjackConfig{
			FileName: "./logs/lumberjack-db-test.log",
		}),
		WithLumberjack(LumberjackDefaultConfig()),
		WithFileConfig(zapcore.EncoderConfig{
			MessageKey:    "msg",
			LevelKey:      "level",
			TimeKey:       "time",
			NameKey:       "logger",
			CallerKey:     "caller",
			StacktraceKey: "stacktrace",
			LineEnding:    zapcore.DefaultLineEnding,
			EncodeLevel:   zapcore.CapitalLevelEncoder,
			EncodeTime: func(date time.Time, encoder zapcore.PrimitiveArrayEncoder) {
				encoder.AppendString(date.Format("2006-01-02 15:04:05.00000000"))
			},
			EncodeDuration: zapcore.SecondsDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
		}),
		WithConsoleConfig(zapcore.EncoderConfig{
			MessageKey:    "msg",
			LevelKey:      "level",
			TimeKey:       "time",
			NameKey:       "logger",
			CallerKey:     "caller",
			StacktraceKey: "stacktrace",
			LineEnding:    zapcore.DefaultLineEnding,
			EncodeLevel:   zapcore.CapitalLevelEncoder,
			EncodeTime: func(date time.Time, encoder zapcore.PrimitiveArrayEncoder) {
				encoder.AppendString(date.Format("2006-01-02 15:04:05.000"))
			},
			EncodeDuration: zapcore.SecondsDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
		}),
		WithJSONFormat(),
	)

	ticker := time.NewTicker(time.Second * 1)

	go func() {
		for {
			select {

			case <-ticker.C:
				// logger.Info("info log")
				logger.Info("info log", zap.Any("info", time.Now()))
				logger.Warn("warn log", zap.Any("warn", time.Now()))
			}
		}
	}()
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	select {
	case <-quit:
		ticker.Stop()
		logger.Debug("stop logger ticker...")
	}

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
