package go_logger

import (
	"fmt"
	"time"

	"github.com/jianlu8023/go-logger/v2/internal/colour"
	"go.uber.org/zap/zapcore"
)

const (
	format = "2006-01-02 15:04:05.000"
)

var (
	consoleEncoderConfig = zapcore.EncoderConfig{
		MessageKey:       "msg",
		LevelKey:         "level",
		TimeKey:          "time",
		NameKey:          "logger",
		CallerKey:        "caller",
		StacktraceKey:    "stacktrace",
		FunctionKey:      "func",
		LineEnding:       zapcore.DefaultLineEnding,
		EncodeLevel:      CustomColorCapitalLevelEncoder,
		EncodeTime:       CustomTimeEncoder,
		EncodeDuration:   zapcore.SecondsDurationEncoder,
		EncodeCaller:     zapcore.ShortCallerEncoder,
		EncodeName:       zapcore.FullNameEncoder,
		ConsoleSeparator: "  ",
	}
	fileEncoderConfig = zapcore.EncoderConfig{
		MessageKey:       "msg",
		LevelKey:         "level",
		TimeKey:          "time",
		NameKey:          "logger",
		CallerKey:        "caller",
		StacktraceKey:    "stacktrace",
		FunctionKey:      "func",
		LineEnding:       zapcore.DefaultLineEnding,
		EncodeLevel:      CustomCapitalLevelEncoder,
		EncodeTime:       CustomTimeEncoder,
		EncodeDuration:   zapcore.SecondsDurationEncoder,
		EncodeCaller:     zapcore.ShortCallerEncoder,
		EncodeName:       zapcore.FullNameEncoder,
		ConsoleSeparator: "  ",
	}
)

func CustomTimeEncoder(time time.Time, encoder zapcore.PrimitiveArrayEncoder) {
	encoder.AppendString(time.Format(format))
}

func CustomCapitalLevelEncoder(level zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString("[" + fmt.Sprintf("%-5s", level.CapitalString()) + "]")
}

func CustomColorCapitalLevelEncoder(level zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
	var colorStr string
	switch level {
	case zapcore.DebugLevel:
		colorStr = colour.Magenta(level.CapitalString())
	case zapcore.InfoLevel:
		colorStr = colour.Blue(level.CapitalString() + " ")
	case zapcore.WarnLevel:
		colorStr = colour.Yellow(level.CapitalString() + " ")
	case zapcore.ErrorLevel, zapcore.DPanicLevel, zapcore.PanicLevel, zapcore.FatalLevel:
		colorStr = colour.Red(level.CapitalString())
	default:
	}
	enc.AppendString("[" + colorStr + "]")
}
