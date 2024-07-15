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

func TestNewLoggerWithNoFile(t *testing.T) {
	logger := NewLogger(&Config{LogLevel: "info", DevelopMode: true})
	ticker(logger)
}

func TestNewLoggerWithLumberJack(t *testing.T) {
	logger := NewLogger(
		&Config{
			LogLevel:    "info",
			DevelopMode: true,
		},
		WithLumberjack(&LumberjackConfig{
			FileName:   "./logs/lumberjack-only-logger.log",
			Localtime:  true,
			Compress:   true,
			MaxSize:    5,
			MaxAge:     30,
			MaxBackups: 7,
		}),
		WithLumberjack(LumberjackDefaultConfig()),
	)
	ticker(logger)

}

func TestNewLogger(t *testing.T) {
	logger := NewLogger(
		&Config{
			LogLevel:    "info",
			DevelopMode: true,
		},
		WithRotateLog(&RotateLogConfig{
			FileName:  "./logs/rotatelog-logger.log",
			LocalTime: true,
		}),
		WithRotateLog(RotateLogDefaultConfig()),
		WithLumberjack(&LumberjackConfig{
			FileName: "./logs/lumberjack-logger.log",
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
		WithConsoleFormat(),
	)
	ticker(logger)

}
func TestNewSugaredLogger(t *testing.T) {

	logger := NewSugaredLogger(&Config{
		LogLevel:    "info",
		DevelopMode: false,
	}, WithRotateLog(&RotateLogConfig{
		FileName:  "./logs/rotatelog-sugared.log",
		LocalTime: true,
	}), WithLumberjack(&LumberjackConfig{
		FileName:  "./logs/lumberjack-sugared.log",
		Localtime: true,
	}))

	tickerSugared(logger)

}

func tickerSugared(logger *zap.SugaredLogger) {
	ticker := time.NewTicker(time.Second * 1)

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

func ticker(log *zap.Logger) {
	ticker := time.NewTicker(time.Second * 1)
	go func() {
		for {
			select {
			case <-ticker.C:
				// logger.Info("info log")
				log.Info("info log", zap.Any("info", time.Now()))
				log.Warn("warn log", zap.Any("warn", time.Now()))
			}
		}
	}()
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	select {
	case <-quit:
		ticker.Stop()
		log.Debug("stop logger ticker...")
	}
}
