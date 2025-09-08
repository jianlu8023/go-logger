package go_logger

import (
	"errors"
	"os"
	"os/signal"
	"syscall"
	"testing"
	"time"
	
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func TestNewLoggerWithoutConsoleOutPut(t *testing.T) {
	logger := NewLogger(
		WithOutConsoleOutPut(),
	)
	logger.Debug("test debug")
}

func TestNewLoggerWithNoFile(t *testing.T) {
	logger := NewLogger()
	ticker(logger)
}

func TestNewLoggerWithLumberJack(t *testing.T) {
	logger := NewLogger(
		WithDefaultLogLevel("info"),
		WithModuleName("[app]"),
		WithStackLogLevel("error"),
		WithCaller(),
		WithFileOutPut(),
		WithLumberjack(&LumberjackConfig{
			FileName:   "./logs/lumberjack-only-logger.log",
			Localtime:  true,
			Compress:   true,
			MaxSize:    5,
			MaxAge:     30,
			MaxBackups: 7,
		}),
		WithLumberjack(LumberjackDefaultConfig()),
		WithConsoleFormat(),
	)
	ticker(logger)
	
}

func TestNewLogger(t *testing.T) {
	logger := NewLogger(
		
		WithDefaultLogLevel("debug"),
		WithDevelopMode(),
		WithModuleName("[sdk]"),
		WithCaller(),
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
			FunctionKey:   "func",
			LineEnding:    zapcore.DefaultLineEnding,
			EncodeLevel:   CustomColorCapitalLevelEncoder,
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
			FunctionKey:   "func",
			LineEnding:    zapcore.DefaultLineEnding,
			EncodeLevel:   CustomColorCapitalLevelEncoder,
			EncodeTime: func(date time.Time, encoder zapcore.PrimitiveArrayEncoder) {
				encoder.AppendString(date.Format("2006-01-02 15:04:05.000"))
			},
			EncodeDuration: zapcore.SecondsDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
		}),
		WithConsoleFormat(),
		WithFileOutPut(),
		WithConsoleOutPut(),
	)
	ticker(logger)
	
}

func TestNewSugaredLogger(t *testing.T) {
	
	logger := NewSugaredLogger(
		WithFileOutPut(),
		WithDefaultLogLevel("DEBUG"),
		WithDevelopMode(),
		WithModuleName("[app]"),
		WithStackLogLevel("error"),
		WithCaller(),
		WithRotateLog(&RotateLogConfig{
			FileName:  "./logs/rotatelog-sugared.log",
			LocalTime: true,
		}),
		WithLumberjack(&LumberjackConfig{
			FileName:  "./logs/lumberjack-sugared.log",
			Localtime: true,
		}),
	)
	logger.Errorf("error info %s", errors.New("this is a test error"))
	tickerSugared(logger)
	
}

func tickerSugared(logger *zap.SugaredLogger) {
	ticker := time.NewTicker(time.Second * 1)
	
	go func() {
		for {
			select {
			case <-ticker.C:
				logger.Infof("info %s", "log")
				logger.Infof("info struct %v", struct {
					Name string
					Age  int
				}{
					Name: "test",
					Age:  18,
				})
				logger.Warnf("warn %s", "log")
				logger.Debugf("debug %s", "log")
				logger.Errorf("error %s", "log")
				logger.Errorf("_error %s", errors.New("test _error"))
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
				log.Debug("debug log", zap.Any("debug", time.Now()))
				log.Error("error log", zap.Any("error", time.Now()))
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
