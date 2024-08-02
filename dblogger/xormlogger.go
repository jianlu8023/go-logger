package dblogger

import (
	"context"

	xormlogger "xorm.io/xorm/log"
)

func (l *Logger) BeforeSQL(ctx xormlogger.LogContext) {}

func (l *Logger) AfterSQL(ctx xormlogger.LogContext) {
	// var sessionPart string
	// v := ctx.Ctx.Value(SessionIDKey)
	// if key, ok := v.(string); ok {
	//	 sessionPart = fmt.Sprintf(" [%s]", key)
	// }
	if ctx.ExecuteTime > 0 {
		l.Info(ctx.Ctx, "\n==> 执行语句: %v \n==> 执行参数: %v \n==> 执行时间: %v",
			ctx.SQL, ctx.Args, ctx.ExecuteTime)
	} else {
		l.Info(ctx.Ctx, "\n==> 执行语句: %v \n==> 执行参数: %v \n==> 执行时间: %v",
			ctx.SQL, ctx.Args, ctx.ExecuteTime)
	}
}

func (l *Logger) Debugf(format string, v ...interface{}) {
	if l.LogLevel >= DEBUG {
		l.logger(context.Background()).Sugar().Debugf(format, v...)
	}
}

func (l *Logger) Errorf(format string, v ...interface{}) {
	if l.LogLevel >= ERROR {
		l.logger(context.Background()).Sugar().Errorf(format, v...)
	}
}

func (l *Logger) Infof(format string, v ...interface{}) {
	if l.LogLevel >= INFO {
		l.logger(context.Background()).Sugar().Infof(format, v...)
	}
}

func (l *Logger) Warnf(format string, v ...interface{}) {
	if l.LogLevel >= WARN {
		l.logger(context.Background()).Sugar().Warnf(format, v...)
	}
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
	default:
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
