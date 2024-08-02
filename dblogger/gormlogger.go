package dblogger

import (
	"context"
	"errors"
	"fmt"
	"runtime"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	dbLogger "gorm.io/gorm/logger"

	"github.com/jianlu8023/go-logger/pkg/colour"
)

func (l *Logger) LogMode(level dbLogger.LogLevel) dbLogger.Interface {
	newLogger := *l
	switch level {
	case dbLogger.Silent:
		newLogger.LogLevel = OFF
	case dbLogger.Warn:
		newLogger.LogLevel = WARN
	case dbLogger.Info:
		newLogger.LogLevel = INFO
	case dbLogger.Error:
		newLogger.LogLevel = ERROR
	}
	return &newLogger
}

func (l *Logger) Info(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= INFO {
		if l.Colorful {
			for _, s := range strings.Split(fmt.Sprintf(msg, data...), "\n") {
				l.logger(ctx).Sugar().Infof(colour.Blue(s))
			}
		} else {
			l.logger(ctx).Sugar().Infof(msg, data...)
		}
	}
}

func (l *Logger) Warn(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= WARN {
		if l.Colorful {
			for _, s := range strings.Split(fmt.Sprintf(msg, data...), "\n") {
				l.logger(ctx).Sugar().Warnf(colour.Yellow(s))
			}
		} else {
			l.logger(ctx).Sugar().Warnf(msg, data...)
		}
	}
}

func (l *Logger) Error(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= ERROR {
		if l.Colorful {
			for _, s := range strings.Split(fmt.Sprintf(msg, data...), "\n") {
				l.logger(ctx).Sugar().Errorf(colour.Red(s))
			}
		} else {
			l.logger(ctx).Sugar().Errorf(msg, data...)
		}
	}
}

func (l *Logger) Trace(ctx context.Context,
	begin time.Time,
	fc func() (sql string, rowsAffected int64),
	err error) {

	if l.LogLevel <= OFF {
		return
	}

	elapsed := time.Since(begin)
	elapsedStr := fmt.Sprintf("%.3fms", float64(elapsed.Nanoseconds())/1e6)
	switch {
	case err != nil && l.LogLevel >= ERROR && (!errors.Is(err, dbLogger.ErrRecordNotFound) || !l.IgnoreRecordNotFoundError):
		sql, rows := fc()
		if rows == -1 {
			l.Error(ctx, "\n==> 执行语句: %v \n==> 影响行数: %v \n==> 执行耗时: %v \n==> 执行错误: %v", sql, rows, elapsedStr, err)
		} else {
			l.Error(ctx, "\n==> 执行语句: %v \n==> 影响行数: %v \n==> 执行耗时: %v \n==> 执行错误: %v", sql, rows, elapsedStr, err)
		}
	case elapsed > l.SlowThreshold && l.SlowThreshold != 0 && l.LogLevel >= WARN:
		sql, rows := fc()
		slowLog := fmt.Sprintf("SLOW SQL >= %v", l.SlowThreshold)
		if rows == -1 {
			l.Warn(ctx, "\n==> 执行语句: %v \n==> 影响行数: %v \n==> 慢SQL: %v \n==> 执行时间: %v", sql, rows, slowLog, elapsedStr)
		} else {
			l.Warn(ctx, "\n==> 执行语句: %v \n==> 影响行数: %v \n==> 慢SQL: %v \n==> 执行时间: %v", sql, rows, slowLog, elapsedStr)
		}
	case l.LogLevel == INFO:
		sql, rows := fc()
		if rows == -1 {
			l.Info(ctx, "\n==> 执行语句: %v \n==> 影响行数: %v \n==> 执行时间: %v", sql, rows, elapsedStr)
		} else {
			l.Info(ctx, "\n==> 执行语句: %v \n==> 影响行数: %v \n==> 执行时间: %v", sql, rows, elapsedStr)
		}
	}
}

func (l *Logger) logger(ctx context.Context) *zap.Logger {
	logger := l.zapLogger
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
