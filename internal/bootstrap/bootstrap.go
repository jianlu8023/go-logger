package bootstrap

import (
	"fmt"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	rotateloggers "github.com/lestrrat-go/file-rotatelogs"
	"go.uber.org/zap"
	"gopkg.in/natefinch/lumberjack.v2"

	"github.com/jianlu8023/go-logger/v2/internal/define"
	lsink "github.com/jianlu8023/go-logger/v2/internal/sink/lumberjack"
	rsink "github.com/jianlu8023/go-logger/v2/internal/sink/rotatelog"
)

func init() {
	if err := zap.RegisterSink(define.Lumberjack, func(url *url.URL) (zap.Sink, error) {
		// Start from a clean snapshot of defaults — no race with concurrent open calls.
		cfg := define.SnapshotLumberjackDefaults()

		query := url.Query()
		if query.Has("fileName") {
			cfg.FileName = query.Get("fileName")
		}
		if query.Has("maxSize") {
			maxsize, err := strconv.Atoi(query.Get("maxSize"))
			if err != nil {
				// 解析失败不 panic，输出 stderr 告警，回退到默认值
				fmt.Fprintf(os.Stderr, "[go-logger] WARNING: invalid maxSize %q: %v, using default\n", query.Get("maxSize"), err)
			} else {
				cfg.MaxSize = maxsize
			}
		}
		if query.Has("maxBackups") {
			backups, err := strconv.Atoi(query.Get("maxBackups"))
			if err != nil {
				fmt.Fprintf(os.Stderr, "[go-logger] WARNING: invalid maxBackups %q: %v, using default\n", query.Get("maxBackups"), err)
			} else {
				cfg.MaxBackups = backups
			}
		}
		if query.Has("maxAge") {
			mage, err := strconv.Atoi(query.Get("maxAge"))
			if err != nil {
				fmt.Fprintf(os.Stderr, "[go-logger] WARNING: invalid maxAge %q: %v, using default\n", query.Get("maxAge"), err)
			} else {
				cfg.MaxAge = mage
			}
		}
		if query.Has("compress") {
			if query.Get("compress") == "true" {
				cfg.Compress = true
			} else {
				cfg.Compress = false
			}
		}
		if query.Has("localtime") {
			if query.Get("localtime") == "true" {
				cfg.Localtime = true
			} else {
				cfg.Localtime = false
			}
		}

		hook := &lumberjack.Logger{
			Filename:   cfg.FileName,
			MaxSize:    cfg.MaxSize, // megabytes
			MaxBackups: cfg.MaxBackups,
			MaxAge:     cfg.MaxAge, // days
			Compress:   cfg.Compress,
			LocalTime:  cfg.Localtime,
		}
		return lsink.NewLumberjack(hook), nil
	}); err != nil {
		panic(err)
	}

	if err := zap.RegisterSink(define.RotateLogs, func(url *url.URL) (zap.Sink, error) {
		// Start from a clean snapshot of defaults — no race with concurrent open calls.
		cfg := define.SnapshotRotatelogsDefaults()

		query := url.Query()

		if query.Has("fileName") {
			cfg.RfileName = query.Get("fileName")
			// 安全剥离 .log 后缀，避免文件名不以 .log 结尾时切片越界 panic
			baseName := strings.TrimSuffix(cfg.RfileName, ".log")
			// _%Y-%m-%d %H:%M:%S
			cfg.RfileName = baseName + ".%Y-%m-%d-%H" + ".log"
			cfg.BaseName = baseName + ".log"
		}
		if query.Has("rotationTime") {
			rotationDurationStr := query.Get("rotationTime")
			rotationDuration, err := time.ParseDuration(rotationDurationStr)
			if err != nil {
				fmt.Fprintf(os.Stderr, "[go-logger] WARNING: invalid rotationTime %q: %v, using default\n", rotationDurationStr, err)
			} else {
				cfg.RotationTime = rotationDuration
			}
		}
		if query.Has("maxAge") {
			maxAgeDurationStr := query.Get("maxAge")
			maxAgeDuration, err := time.ParseDuration(maxAgeDurationStr)
			if err != nil {
				fmt.Fprintf(os.Stderr, "[go-logger] WARNING: invalid maxAge %q: %v, using default\n", maxAgeDurationStr, err)
			} else {
				cfg.RmaxAge = maxAgeDuration
			}
		}

		if query.Has("localtime") {
			if query.Get("localtime") == "true" {
				cfg.Rlocaltime = time.Local
				cfg.Rclock = rotateloggers.Local
			} else {
				cfg.Rlocaltime = time.UTC
				cfg.Rclock = rotateloggers.UTC
			}
		}

		logs, err := rotateloggers.New(
			cfg.RfileName,
			rotateloggers.WithLinkName(cfg.BaseName),
			rotateloggers.WithMaxAge(cfg.RmaxAge),
			rotateloggers.WithRotationTime(cfg.RotationTime),
			rotateloggers.WithLocation(cfg.Rlocaltime),
			rotateloggers.WithClock(cfg.Rclock),
		)
		if err != nil {
			return nil, err
		}
		return rsink.NewRotateLog(logs), nil
	}); err != nil {
		panic(err)
	}
}
