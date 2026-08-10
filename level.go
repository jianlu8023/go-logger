package go_logger

import (
	"go.uber.org/zap/zapcore"
	"strings"
)

const (
	warnLevel  = "warn"
	infoLevel  = "info"
	debugLevel = "debug"
	errorLevel = "error"
	fatalLevel = "fatal"
	panicLevel = "panic"
)

func logLevel(level string) zapcore.Level {
	level = strings.ToLower(level)
	switch level {
	case infoLevel:
		return zapcore.InfoLevel
	case debugLevel:
		return zapcore.DebugLevel
	case warnLevel:
		return zapcore.WarnLevel
	case errorLevel:
		return zapcore.ErrorLevel
	case panicLevel:
		return zapcore.PanicLevel
	case fatalLevel:
		return zapcore.FatalLevel
	default:
		return zapcore.InfoLevel
	}
}

// getFileLogLevel 获取输出到文件的日志级别
func getFileLogLevel(optMap map[string]Option) zapcore.Level {
	if opt, ok := optMap[fileLogLevelKey]; ok {
		return opt.Value().(zapcore.Level)
	}
	return getLogLevel(optMap)
}

// getConsoleLogLevel 获取输出到控制台的日志级别
func getConsoleLogLevel(optMap map[string]Option) zapcore.Level {
	if opt, ok := optMap[consoleLogLevelKey]; ok {
		return opt.Value().(zapcore.Level)
	}
	return getLogLevel(optMap)
}

// getLogLevel 从options中获取日志级别
func getLogLevel(optMap map[string]Option) zapcore.Level {
	if opt, ok := optMap[defaultLogLevelKey]; ok {
		return opt.Value().(zapcore.Level)
	}
	return zapcore.DebugLevel
}
