package go_logger

import (
	"fmt"
	"os"

	"go.uber.org/zap"

	_ "github.com/jianlu8023/go-logger/v2/internal/bootstrap"
)

// NewLogger 根据传入的 options 创建并返回一个可用的 *zap.Logger。
//
// 设计说明：
//   - logger 是程序基础设施，初始化失败不应导致程序崩溃，因此本函数永不 panic。
//   - 当检测到配置冲突（如同时指定多种日志格式）时，降级为 console 格式，
//     并通过 stderr 输出告警，便于开发者在开发阶段发现问题。
//   - 仍然只返回 *zap.Logger（不返回 error），调用方可安全使用返回值。
func NewLogger(options ...Option) *zap.Logger {
	if both, option := checkFormat(options); both {
		// 格式冲突：不 panic，降级为 console 格式，并通过 stderr 告警
		fmt.Fprintln(os.Stderr, "[go-logger] WARNING: multiple log formats specified, falling back to console format")
		return consoleLogger(options...)
	} else {
		switch option.Name() {
		case zaplogfmtKey:
			return zapLogFmtLogger(options...)
		case jsonFormatKey:
			return jsonLogger(options...)
		case consoleFormatKey:
			fallthrough
		default:
			return consoleLogger(options...)
		}
	}
}

func NewSugaredLogger(options ...Option) *zap.SugaredLogger {
	return NewLogger(options...).Sugar()
}
