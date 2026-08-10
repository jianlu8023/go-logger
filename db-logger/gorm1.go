package db_logger

import (
	"context"
	"database/sql/driver"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unicode"
)

func (l Logger) Print(values ...interface{}) {
	messages := LogFormatter(values...)
	stres := make([]string, 0, len(messages))
	for _, message := range messages {
		stres = append(stres, message.(string))
	}
	str := strings.Join(stres, " ")
	l.Info(context.Background(), "%v", str)
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
