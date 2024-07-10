package db

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"runtime"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	gormlogger "gorm.io/gorm/logger"
)

type Logger struct {
	ZapLogger                 *zap.Logger
	LogLevel                  gormlogger.LogLevel
	SlowThreshold             time.Duration
	Colorful                  bool
	IgnoreRecordNotFoundError bool
	ParameterizedQueries      bool
}

func (l *Logger) String() string {
	bytes, _ := json.Marshal(l)

	return string(bytes)
}

func (l *Logger) LogMode(level gormlogger.LogLevel) gormlogger.Interface {
	newLogger := *l
	newLogger.LogLevel = level
	return &newLogger
}

func (l *Logger) Info(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= gormlogger.Info {
		l.logger(ctx).Sugar().Info(msg, data)
	}
}

func (l *Logger) Warn(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= gormlogger.Warn {
		l.logger(ctx).Sugar().Warn(msg, data)
	}
}

func (l *Logger) Error(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= gormlogger.Error {
		l.logger(ctx).Sugar().Error(msg, data)
	}
}

func (l *Logger) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
	if l.LogLevel <= gormlogger.Silent {
		return
	}

	elapsed := time.Since(begin)
	elapsedStr := fmt.Sprintf("%.3fms", float64(elapsed.Nanoseconds())/1e6)
	logger := l.logger(ctx)
	switch {
	case err != nil && l.LogLevel >= gormlogger.Error && (!errors.Is(err, gormlogger.ErrRecordNotFound) || !l.IgnoreRecordNotFoundError):
		fmt.Println("err != nil && l.LogLevel >= gormlogger.Error && (!errors.Is(err, gormlogger.ErrRecordNotFound) || !l.IgnoreRecordNotFoundError)")
		sql, rows := fc()
		if rows == -1 {
			logger.Sugar().Errorf("\n执行时间: %s \n影响行数: %v \n执行错误: %s \n执行语句: %s", elapsedStr, rows, sql, err)
		} else {
			logger.Sugar().Errorf("\n执行时间 %s \n影响行数: %v \n执行错误: %s \n执行语句: %s", elapsedStr, rows, sql, err)
		}
	case elapsed > l.SlowThreshold && l.SlowThreshold != 0 && l.LogLevel >= gormlogger.Warn:
		fmt.Println("elapsed > l.SlowThreshold && l.SlowThreshold != 0 && l.LogLevel >= gormlogger.Warn")
		sql, rows := fc()
		slowLog := fmt.Sprintf("SLOW SQL >= %v", l.SlowThreshold)
		if rows == -1 {
			logger.Sugar().Warnf("\n执行时间 %s \n影响行数: %v \n慢SQL: %v \n执行语句: %s", elapsedStr, rows, slowLog, sql)
		} else {
			logger.Sugar().Warnf("\n执行时间 %s \n影响行数: %v \n慢SQL: %v \n执行语句: %s", elapsedStr, rows, slowLog, sql)
		}
	case l.LogLevel == gormlogger.Info:
		fmt.Println("l.LogLevel == gormlogger.Info")
		sql, rows := fc()
		if rows == -1 {
			logger.Sugar().Infof("\n==> 执行时间 %s \n==> 影响行数: %v \n==> 执行语句: %s", elapsedStr, rows, sql)
		} else {
			logger.Sugar().Infof("\n==> 执行时间 %s \n==> 影响行数: %v \n==> 执行语句: %s", elapsedStr, rows, sql)
		}
	}
}

func NewDevelopDBLogger(config Config) gormlogger.Interface {
	return &Logger{
		ZapLogger:                 config.Logger,
		LogLevel:                  gormlogger.Error,
		SlowThreshold:             config.SlowThreshold,
		Colorful:                  config.Colorful,
		IgnoreRecordNotFoundError: config.IgnoreRecordNotFoundError,
		ParameterizedQueries:      config.ParameterizedQueries,
	}
}

func NewProductionDBLogger(config Config) gormlogger.Interface {
	return &Logger{
		ZapLogger:                 config.Logger,
		LogLevel:                  gormlogger.Info,
		SlowThreshold:             config.SlowThreshold,
		Colorful:                  config.Colorful,
		IgnoreRecordNotFoundError: config.IgnoreRecordNotFoundError,
		ParameterizedQueries:      config.ParameterizedQueries,
	}
}

func (l *Logger) logger(ctx context.Context) *zap.Logger {
	logger := l.ZapLogger
	if ctx != nil {
		if c, ok := ctx.(*gin.Context); ok {
			ctx = c.Request.Context()
		}
		zl := ctx.Value(ctxLoggerKey)
		ctxLogger, ok := zl.(*zap.Logger)
		if ok {
			logger = ctxLogger
		}
	}
	for i := 2; i < 15; i++ {
		_, file, _, ok := runtime.Caller(i)
		switch {
		case !ok:
		case strings.HasSuffix(file, "_test.go"):
		case strings.Contains(file, gormPackage):
		default:
			return logger.WithOptions(zap.AddCallerSkip(i - 1))
		}
	}
	return logger
}
