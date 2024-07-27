package dblogger

import (
	"context"
	"errors"
	"fmt"
	"runtime"
	"strings"
	"time"

	"github.com/bytedance/sonic"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	dbLogger "gorm.io/gorm/logger"
	xormlogger "xorm.io/xorm/log"
)

const (
	SessionIDKey = "__xorm_session_id"
)

type LogLevel int

const (
	OFF = iota + 1
	DEBUG
	INFO
	WARN
	ERROR
)

type Logger struct {
	ZapLogger                 *zap.Logger
	LogLevel                  LogLevel
	SlowThreshold             time.Duration
	Colorful                  bool
	IgnoreRecordNotFoundError bool
	ParameterizedQueries      bool
	showSql                   bool
}

func (l *Logger) String() string {
	bytes, _ := sonic.Marshal(l)
	return string(bytes)
}

func (l *Logger) BeforeSQL(ctx xormlogger.LogContext) {}

func (l *Logger) AfterSQL(ctx xormlogger.LogContext) {
	//var sessionPart string
	//v := ctx.Ctx.Value(SessionIDKey)
	//if key, ok := v.(string); ok {
	//	sessionPart = fmt.Sprintf(" [%s]", key)
	//}
	if ctx.ExecuteTime > 0 {
		l.Info(ctx.Ctx, "\n==> 执行语句: %v \n==> 影响行数: %v \n==> 执行时间: %v\n",
			ctx.SQL, ctx.Args, ctx.ExecuteTime)
	} else {
		l.Info(ctx.Ctx, "\n==> 执行语句: %v \n==> 影响行数: %v \n==> 执行时间: %v\n",
			ctx.SQL, ctx.Args, ctx.ExecuteTime)
	}
}

func (l *Logger) Debugf(format string, v ...interface{}) {
	l.logger(context.Background()).Sugar().Debugf(format, v...)
}

func (l *Logger) Errorf(format string, v ...interface{}) {
	l.logger(context.Background()).Sugar().Errorf(format, v...)
}

func (l *Logger) Infof(format string, v ...interface{}) {
	l.logger(context.Background()).Sugar().Infof(format, v...)
}

func (l *Logger) Warnf(format string, v ...interface{}) {
	l.logger(context.Background()).Sugar().Warnf(format, v...)
}

func (l *Logger) Level() xormlogger.LogLevel {
	switch l.LogLevel {
	case INFO:
		return xormlogger.LOG_INFO
	case WARN:
		return xormlogger.LOG_WARNING
	case ERROR:
		return xormlogger.LOG_ERR
	case OFF:
		return xormlogger.LOG_OFF
	case DEBUG:
		return xormlogger.LOG_DEBUG
	default:
		return xormlogger.LOG_UNKNOWN
	}
}

func (l *Logger) SetLevel(lv xormlogger.LogLevel) {
	switch lv {
	case xormlogger.LOG_DEBUG:
		l.LogLevel = DEBUG
	case xormlogger.LOG_INFO:
		l.LogLevel = INFO
	case xormlogger.LOG_WARNING:
		l.LogLevel = WARN
	case xormlogger.LOG_ERR:
		l.LogLevel = ERROR
	case xormlogger.LOG_OFF:
		l.LogLevel = OFF
	}
}

func (l *Logger) ShowSQL(show ...bool) {
	if len(show) == 0 {
		l.showSql = true
		return
	}
	l.showSql = show[0]
}

func (l *Logger) IsShowSQL() bool {
	return l.showSql
}

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
		l.logger(ctx).Sugar().Infof(msg, data...)
	}
}

func (l *Logger) Warn(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= WARN {
		l.logger(ctx).Sugar().Warnf(msg, data...)
	}
}

func (l *Logger) Error(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= ERROR {
		l.logger(ctx).Sugar().Errorf(msg, data...)
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
			l.Error(ctx, "\n==> 执行语句: %v \n==> 影响行数: %v \n==> 执行耗时: %v \n==> 执行错误: %v\n", sql, rows, elapsedStr, err)
		} else {
			l.Error(ctx, "\n==> 执行语句: %v \n==> 影响行数: %v \n==> 执行耗时: %v \n==> 执行错误: %v\n", sql, rows, elapsedStr, err)
		}
	case elapsed > l.SlowThreshold && l.SlowThreshold != 0 && l.LogLevel >= WARN:
		sql, rows := fc()
		slowLog := fmt.Sprintf("SLOW SQL >= %v", l.SlowThreshold)
		if rows == -1 {
			l.Warn(ctx, "\n==> 执行语句: %v \n==> 影响行数: %v \n==> 慢SQL: %v \n==> 执行时间: %v\n", sql, rows, slowLog, elapsedStr)
		} else {
			l.Warn(ctx, "\n==> 执行语句: %v \n==> 影响行数: %v \n==> 慢SQL: %v \n==> 执行时间: %v\n", sql, rows, slowLog, elapsedStr)
		}
	case l.LogLevel == INFO:
		sql, rows := fc()
		if rows == -1 {
			l.Info(ctx, "\n==> 执行语句: %v \n==> 影响行数: %v \n==> 执行时间: %v\n", sql, rows, elapsedStr)
		} else {
			l.Info(ctx, "\n==> 执行语句: %v \n==> 影响行数: %v \n==> 执行时间: %v\n", sql, rows, elapsedStr)
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

func NewDBLogger(config Config) *Logger {
	return &Logger{
		ZapLogger:                 config.Logger,
		LogLevel:                  config.LogLevel,
		SlowThreshold:             config.SlowThreshold,
		Colorful:                  config.Colorful,
		IgnoreRecordNotFoundError: config.IgnoreRecordNotFoundError,
		ParameterizedQueries:      config.ParameterizedQueries,
		showSql:                   config.ShowSql,
	}
}
