package go_logger

import (
	"go.uber.org/zap/zapcore"
)

const (
	warn   = "warn"
	info   = "info"
	debug  = "debug"
	_error = "error"
	fatal  = "fatal"
	_panic = "panic"
)

func logLevel(level string) zapcore.Level {
	switch level {
	case info:
		return zapcore.InfoLevel
	case debug:
		return zapcore.DebugLevel
	case warn:
		return zapcore.WarnLevel
	case _error:
		return zapcore.ErrorLevel
	case _panic:
		return zapcore.PanicLevel
	case fatal:
		return zapcore.FatalLevel
	default:
		return zapcore.InfoLevel
	}
}

// getFileLogLevel 获取输出到文件的日志级别
func getFileLogLevel(options []Option) zapcore.Level {
	for _, opt := range options {
		if opt.Name() == fileLogLevelKey {
			return opt.Value().(zapcore.Level)
		}
	}
	return getLogLevel(options)
}

// getConsoleLogLevel 获取输出到控制台的日志级别
func getConsoleLogLevel(options []Option) zapcore.Level {
	for _, opt := range options {
		if opt.Name() == consoleLogLevelKey {
			return opt.Value().(zapcore.Level)
		}
	}
	return getLogLevel(options)
}

// getLogLevel 从options中获取日志级别
func getLogLevel(options []Option) zapcore.Level {
	for _, opt := range options {
		if opt.Name() == defaultLogLevelKey {
			return opt.Value().(zapcore.Level)
		}
	}
	return zapcore.DebugLevel
}
