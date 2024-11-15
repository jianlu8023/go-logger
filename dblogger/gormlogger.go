package dblogger

import (
	"context"
	"database/sql/driver"
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"time"
	"unicode"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	dbLogger "gorm.io/gorm/logger"

	"github.com/jianlu8023/go-tools/pkg/format/colour"
)

func (l Logger) LogMode(level dbLogger.LogLevel) dbLogger.Interface {
	newLogger := l
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

func (l Logger) Info(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= INFO {
		if l.Colorful {
			for _, s := range strings.Split(fmt.Sprintf(msg, data...), "\n") {
				l.logger(ctx).Sugar().Infof(colour.Blue(s))
			}
			// l.logger(ctx).Sugar().Infof(colour.Blue(fmt.Sprintf(msg, data...)))
		} else {
			l.logger(ctx).Sugar().Infof(msg, data...)
		}
	}
}

func (l Logger) Warn(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= WARN {
		if l.Colorful {
			for _, s := range strings.Split(fmt.Sprintf(msg, data...), "\n") {
				l.logger(ctx).Sugar().Warnf(colour.Yellow(s))
			}
			// l.logger(ctx).Sugar().Warnf(colour.Yellow(fmt.Sprintf(msg, data...)))
		} else {
			l.logger(ctx).Sugar().Warnf(msg, data...)
		}
	}
}

func (l Logger) Error(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= ERROR {
		if l.Colorful {
			for _, s := range strings.Split(fmt.Sprintf(msg, data...), "\n") {
				l.logger(ctx).Sugar().Errorf(colour.Red(s))
			}
			// l.logger(ctx).Sugar().Errorf(colour.Red(fmt.Sprintf(msg, data...)))
		} else {
			l.logger(ctx).Sugar().Errorf(msg, data...)
		}
	}
}

func (l Logger) Trace(ctx context.Context,
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
			// l.Error(ctx, "==> 执行语句: %v", sql)
			// l.Error(ctx, "==> 影响行数: %v", rows)
			// l.Error(ctx, "==> 执行耗时: %v", elapsedStr)
			// l.Error(ctx, "==> 执行错误: %v", err)
			l.Error(ctx, "\n==> 执行语句: %v \n==> 影响行数: %v \n==> 执行耗时: %v \n==> 执行错误: %v", sql, rows, elapsedStr, err)
		} else {
			// l.Error(ctx, "==> 执行语句: %v", sql)
			// l.Error(ctx, "==> 影响行数: %v", rows)
			// l.Error(ctx, "==> 执行耗时: %v", elapsedStr)
			// l.Error(ctx, "==> 执行错误: %v", err)
			l.Error(ctx, "\n==> 执行语句: %v \n==> 影响行数: %v \n==> 执行耗时: %v \n==> 执行错误: %v", sql, rows, elapsedStr, err)
		}
	case elapsed > l.SlowThreshold && l.SlowThreshold != 0 && l.LogLevel >= WARN:
		sql, rows := fc()
		slowLog := fmt.Sprintf("SLOW SQL >= %v", l.SlowThreshold)
		if rows == -1 {
			// l.Warn(ctx, "==> 执行语句: %v", sql)
			// l.Warn(ctx, "==> 影响行数: %v", rows)
			// l.Warn(ctx, "==> 慢SQL: %v", slowLog)
			// l.Warn(ctx, "==> 执行耗时: %v", elapsedStr)

			l.Warn(ctx, "\n==> 执行语句: %v \n==> 影响行数: %v \n==> 慢SQL: %v \n==> 执行时间: %v", sql, rows, slowLog, elapsedStr)
		} else {
			// l.Warn(ctx, "==> 执行语句: %v", sql)
			// l.Warn(ctx, "==> 影响行数: %v", rows)
			// l.Warn(ctx, "==> 慢SQL: %v", slowLog)
			// l.Warn(ctx, "==> 执行耗时: %v", elapsedStr)
			l.Warn(ctx, "\n==> 执行语句: %v \n==> 影响行数: %v \n==> 慢SQL: %v \n==> 执行时间: %v", sql, rows, slowLog, elapsedStr)
		}
	case l.LogLevel == INFO:
		sql, rows := fc()
		if rows == -1 {
			// l.Info(ctx, "==> 执行语句: %v", sql)
			// l.Info(ctx, "==> 影响行数: %v", rows)
			// l.Info(ctx, "==> 执行耗时: %v", elapsedStr)
			l.Info(ctx, "\n==> 执行语句: %v \n==> 影响行数: %v \n==> 执行时间: %v", sql, rows, elapsedStr)
		} else {
			// l.Info(ctx, "==> 执行语句: %v", sql)
			// l.Info(ctx, "==> 影响行数: %v", rows)
			// l.Info(ctx, "==> 执行耗时: %v", elapsedStr)
			l.Info(ctx, "\n==> 执行语句: %v \n==> 影响行数: %v \n==> 执行时间: %v", sql, rows, elapsedStr)
		}
	}
}

func (l Logger) logger(ctx context.Context) *zap.Logger {
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

func (l Logger) Print(values ...interface{}) {
	messages := LogFormatter(values...)
	l.Info(context.Background(), "%v", messages)
}

var LogFormatter = func(values ...interface{}) (messages []interface{}) {
	if len(values) > 1 {
		var (
			sql             string
			formattedValues []string
			level           = values[0]
		)

		// messages = []interface{}{source, currentTime}

		// if len(values) == 2 {
		// remove the line break
		// currentTime = currentTime[1:]
		// remove the brackets
		// source = fmt.Sprintf("\033[35m%v\033[0m", values[1])

		// messages = []interface{}{currentTime, source}
		// }

		if level == "sql" {
			// duration
			messages = append(messages, fmt.Sprintf("执行耗时: %.2fms\t", float64(values[2].(time.Duration).Nanoseconds()/1e4)/100.0))

			// sql
			for _, value := range values[4].([]interface{}) {
				indirectValue := reflect.Indirect(reflect.ValueOf(value))
				if indirectValue.IsValid() {
					value = indirectValue.Interface()
					if t, ok := value.(time.Time); ok {
						if t.IsZero() {
							formattedValues = append(formattedValues, fmt.Sprintf("'%v'", "0000-00-00 00:00:00"))
						} else {
							formattedValues = append(formattedValues, fmt.Sprintf("'%v'", t.Format("2006-01-02 15:04:05")))
						}
					} else if b, ok := value.([]byte); ok {
						if str := string(b); isPrintable(str) {
							formattedValues = append(formattedValues, fmt.Sprintf("'%v'", str))
						} else {
							formattedValues = append(formattedValues, "'<binary>'")
						}
					} else if r, ok := value.(driver.Valuer); ok {
						if value, err := r.Value(); err == nil && value != nil {
							formattedValues = append(formattedValues, fmt.Sprintf("'%v'", value))
						} else {
							formattedValues = append(formattedValues, "NULL")
						}
					} else {
						switch value.(type) {
						case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, float32, float64, bool:
							formattedValues = append(formattedValues, fmt.Sprintf("%v", value))
						default:
							formattedValues = append(formattedValues, fmt.Sprintf("'%v'", value))
						}
					}
				} else {
					formattedValues = append(formattedValues, "NULL")
				}
			}

			// differentiate between $n placeholders or else treat like ?
			if numericPlaceHolderRegexp.MatchString(values[3].(string)) {
				sql = values[3].(string)
				for index, value := range formattedValues {
					placeholder := fmt.Sprintf(`\$%d([^\d]|$)`, index+1)
					sql = regexp.MustCompile(placeholder).ReplaceAllString(sql, value+"$1")
				}
			} else {
				formattedValuesLength := len(formattedValues)
				for index, value := range sqlRegexp.Split(values[3].(string), -1) {
					sql += value
					if index < formattedValuesLength {
						sql += formattedValues[index]
					}
				}
			}

			messages = append(messages, fmt.Sprintf("执行SQL: %v\t", sql))

			messages = append(messages, fmt.Sprintf("影响行数: %v", strconv.FormatInt(values[5].(int64), 10)))
		} else {
			// messages = append(messages, "\033[31;1m")
			messages = append(messages, values[2:]...)
			// messages = append(messages, "\033[0m")
		}
	}

	return
}

func isPrintable(s string) bool {
	for _, r := range s {
		if !unicode.IsPrint(r) {
			return false
		}
	}
	return true
}
