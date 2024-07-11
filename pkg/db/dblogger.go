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
		l.logger(ctx).Sugar().Infof(msg, data...)
	}
}

func (l *Logger) Warn(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= gormlogger.Warn {
		l.logger(ctx).Sugar().Warnf(msg, data...)
	}
}

func (l *Logger) Error(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= gormlogger.Error {
		l.logger(ctx).Sugar().Errorf(msg, data...)
	}
}

func (l *Logger) Trace(ctx context.Context,
	begin time.Time,
	fc func() (sql string, rowsAffected int64),
	err error) {

	if l.LogLevel <= gormlogger.Silent {
		return
	}

	elapsed := time.Since(begin)
	elapsedStr := fmt.Sprintf("%.3fms", float64(elapsed.Nanoseconds())/1e6)
	switch {
	case err != nil && l.LogLevel >= gormlogger.Error && (!errors.Is(err, gormlogger.ErrRecordNotFound) || !l.IgnoreRecordNotFoundError):
		sql, rows := fc()
		if rows == -1 {
			l.Error(ctx, "\n==> 执行语句: %v \n==> 影响行数: %v \n==> 执行耗时: %v \n==> 执行错误: %v\n", sql, rows, elapsedStr, err)
		} else {
			l.Error(ctx, "\n==> 执行语句: %v \n==> 影响行数: %v \n==> 执行耗时: %v \n==> 执行错误: %v\n", sql, rows, elapsedStr, err)
		}
	case elapsed > l.SlowThreshold && l.SlowThreshold != 0 && l.LogLevel >= gormlogger.Warn:
		sql, rows := fc()
		slowLog := fmt.Sprintf("SLOW SQL >= %v", l.SlowThreshold)
		if rows == -1 {
			l.Warn(ctx, "\n==> 执行语句 %v \n==> 影响行数: %v \n==> 慢SQL: %v \n==> 执行时间: %v\n", sql, rows, slowLog, elapsedStr)
		} else {
			l.Warn(ctx, "\n==> 执行语句 %v \n==> 影响行数: %v \n==> 慢SQL: %v \n==> 执行时间: %v\n", sql, rows, slowLog, elapsedStr)
		}
	case l.LogLevel == gormlogger.Info:
		sql, rows := fc()
		if rows == -1 {
			l.Info(ctx, "\n==> 执行语句 %v \n==> 影响行数: %v \n==> 执行时间: %v\n", sql, rows, elapsedStr)
		} else {
			l.Info(ctx, "\n==> 执行语句 %v \n==> 影响行数: %v \n==> 执行时间: %v\n", sql, rows, elapsedStr)
		}
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
		case strings.Contains(file, zapgormPackage):
		default:
			return logger.WithOptions(zap.AddCallerSkip(i - 1))
		}
	}
	return logger
}

func NewDBLogger(config Config) gormlogger.Interface {
	level := dbMap[config.LogLevel]
	return &Logger{
		ZapLogger:                 config.Logger,
		LogLevel:                  level,
		SlowThreshold:             config.SlowThreshold,
		Colorful:                  config.Colorful,
		IgnoreRecordNotFoundError: config.IgnoreRecordNotFoundError,
		ParameterizedQueries:      config.ParameterizedQueries,
	}
}
