package dblogger

import (
	"context"
	"time"

	gormlogger "gorm.io/gorm/logger"
	xormlogger "xorm.io/xorm/log"
)

type GormLoggerInterface interface {
	LogMode(gormlogger.LogLevel) gormlogger.Interface
	Info(context.Context, string, ...interface{})
	Warn(context.Context, string, ...interface{})
	Error(context.Context, string, ...interface{})
	Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error)
}

type XormLoggerInterface interface {
	BeforeSQL(context xormlogger.LogContext) // only invoked when IsShowSQL is true
	AfterSQL(context xormlogger.LogContext)
	Debugf(format string, v ...interface{})
	Errorf(format string, v ...interface{})
	Infof(format string, v ...interface{})
	Warnf(format string, v ...interface{})

	Level() xormlogger.LogLevel
	SetLevel(l xormlogger.LogLevel)
	ShowSQL(show ...bool)
	IsShowSQL() bool
}
