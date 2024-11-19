package bootstrap

import (
	"net/url"
	"strconv"
	"time"

	rotateloggers "github.com/lestrrat-go/file-rotatelogs"
	"go.uber.org/zap"
	"gopkg.in/natefinch/lumberjack.v2"

	"github.com/jianlu8023/go-logger/internal/define"
	lsink "github.com/jianlu8023/go-logger/internal/sink/lumberjack"
	rsink "github.com/jianlu8023/go-logger/internal/sink/rotatelog"
)

func init() {
	if err := zap.RegisterSink(define.Lumberjack, func(url *url.URL) (zap.Sink, error) {
		query := url.Query()
		if query.Has("fileName") {
			define.FileName = query.Get("fileName")
		}
		if query.Has("maxSize") {
			maxsize, err := strconv.Atoi(query.Get("maxSize"))
			if err == nil {
				define.MaxSize = maxsize
			}
		}
		if query.Has("maxBackups") {
			backups, err := strconv.Atoi(query.Get("maxBackups"))
			if err == nil {
				define.MaxBackups = backups
			}
		}
		if query.Has("maxAge") {
			mage, err := strconv.Atoi(query.Get("maxAge"))
			if err == nil {
				define.MaxAge = mage
			}
		}
		if query.Has("compress") {
			if query.Get("compress") == "true" {
				define.Compress = true
			} else {
				define.Compress = false
			}
		}
		if query.Has("localtime") {
			if query.Get("localtime") == "true" {
				define.Localtime = true
			} else {
				define.Localtime = false
			}
		}
		hook := &lumberjack.Logger{
			Filename:   define.FileName,
			MaxSize:    define.MaxSize, // megabytes
			MaxBackups: define.MaxBackups,
			MaxAge:     define.MaxAge, // days
			Compress:   define.Compress,
			LocalTime:  define.Localtime,
		}
		return lsink.NewLumberjack(hook), nil
	}); err != nil {
		panic(err)
	}

	if err := zap.RegisterSink(define.RotateLogs, func(url *url.URL) (zap.Sink, error) {
		query := url.Query()

		if query.Has("fileName") {
			define.RfileName = query.Get("fileName")
			define.BaseName = define.RfileName[:len(define.RfileName)-len(".log")]
			// _%Y-%m-%d %H:%M:%S
			define.RfileName = define.BaseName + ".%Y-%m-%d-%H" + ".log"
			define.BaseName = define.BaseName + ".log"
		}
		if query.Has("rotationTime") {
			rotationDurationStr := query.Get("rotationTime")
			rotationDuration, err := time.ParseDuration(rotationDurationStr)
			if err == nil {
				define.RotationTime = rotationDuration
			}
		}
		if query.Has("maxAge") {
			maxAgeDurationStr := query.Get("maxAge")
			maxAgeDuration, err := time.ParseDuration(maxAgeDurationStr)
			if err == nil {
				define.RmaxAge = maxAgeDuration
			}
		}

		if query.Has("localtime") {

			if query.Get("localtime") == "Local" {
				define.Rlocaltime = time.Local
				define.Rclock = rotateloggers.Local
			} else {
				define.Rlocaltime = time.UTC
				define.Rclock = rotateloggers.UTC
			}

		}

		logs, err := rotateloggers.New(
			define.RfileName,
			rotateloggers.WithLinkName(define.BaseName),
			rotateloggers.WithMaxAge(define.RmaxAge),
			rotateloggers.WithRotationTime(define.RotationTime),
			rotateloggers.WithLocation(define.Rlocaltime),
			rotateloggers.WithClock(define.Rclock),
		)
		if err != nil {
			return nil, err
		}
		return rsink.NewRotateLog(logs), nil
	}); err != nil {
		panic(err)
	}
}
