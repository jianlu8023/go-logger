package bootstrap

import (
	"net/url"
	"strconv"
	"time"

	rotateloggers "github.com/lestrrat-go/file-rotatelogs"
	"go.uber.org/zap"
	"gopkg.in/natefinch/lumberjack.v2"

	lsink "github.com/jianlu8023/go-logger/internal/sink/lumberjack"
	rsink "github.com/jianlu8023/go-logger/internal/sink/rotatelog"
	"github.com/jianlu8023/go-logger/pkg/df"
)

func init() {
	if err := zap.RegisterSink(df.Lumberjack, func(url *url.URL) (zap.Sink, error) {
		query := url.Query()
		if query.Has("fileName") {
			df.FileName = query.Get("fileName")
		}
		if query.Has("maxSize") {
			maxsize, err := strconv.Atoi(query.Get("maxSize"))
			if err == nil {
				df.MaxSize = maxsize
			}
		}
		if query.Has("maxBackups") {
			backups, err := strconv.Atoi(query.Get("maxBackups"))
			if err == nil {
				df.MaxBackups = backups
			}
		}
		if query.Has("maxAge") {
			mage, err := strconv.Atoi(query.Get("maxAge"))
			if err == nil {
				df.MaxAge = mage
			}
		}
		if query.Has("compress") {
			if query.Get("compress") == "true" {
				df.Compress = true
			} else {
				df.Compress = false
			}
		}
		if query.Has("localtime") {
			if query.Get("localtime") == "true" {
				df.Localtime = true
			} else {
				df.Localtime = false
			}
		}
		hook := &lumberjack.Logger{
			Filename:   df.FileName,
			MaxSize:    df.MaxSize, // megabytes
			MaxBackups: df.MaxBackups,
			MaxAge:     df.MaxAge, // days
			Compress:   df.Compress,
			LocalTime:  df.Localtime,
		}
		return lsink.NewLumberjack(hook), nil
	}); err != nil {
		panic(err)
	}

	if err := zap.RegisterSink(df.RotateLogs, func(url *url.URL) (zap.Sink, error) {
		query := url.Query()

		if query.Has("fileName") {
			df.RfileName = query.Get("fileName")
			df.BaseName = df.RfileName[:len(df.RfileName)-len(".log")]
			// _%Y-%m-%d %H:%M:%S
			df.RfileName = df.BaseName + ".%Y-%m-%d-%H" + ".log"
			df.BaseName = df.BaseName + ".log"
		}
		if query.Has("rotationTime") {
			rotationDurationStr := query.Get("rotationTime")
			rotationDuration, err := time.ParseDuration(rotationDurationStr)
			if err == nil {
				df.RotationTime = rotationDuration
			}
		}
		if query.Has("maxAge") {
			maxAgeDurationStr := query.Get("maxAge")
			maxAgeDuration, err := time.ParseDuration(maxAgeDurationStr)
			if err == nil {
				df.RmaxAge = maxAgeDuration
			}
		}

		if query.Has("localtime") {

			if query.Get("localtime") == "Local" {
				df.Rlocaltime = time.Local
				df.Rclock = rotateloggers.Local
			} else {
				df.Rlocaltime = time.UTC
				df.Rclock = rotateloggers.UTC
			}

		}

		logs, err := rotateloggers.New(
			df.RfileName,
			rotateloggers.WithLinkName(df.BaseName),
			rotateloggers.WithMaxAge(df.RmaxAge),
			rotateloggers.WithRotationTime(df.RotationTime),
			rotateloggers.WithLocation(df.Rlocaltime),
			rotateloggers.WithClock(df.Rclock),
		)
		if err != nil {
			return nil, err
		}
		return rsink.NewRotatelog(logs), nil
	}); err != nil {
		panic(err)
	}
}
