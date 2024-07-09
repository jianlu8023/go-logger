package go_logger

import (
	"fmt"
	"net/url"
	"strconv"
	"time"

	rotateloggers "github.com/lestrrat-go/file-rotatelogs"
	"go.uber.org/zap"
	"gopkg.in/natefinch/lumberjack.v2"
)

func init() {

	fmt.Println("init")
	if err := zap.RegisterSink(Lumberjack, func(url *url.URL) (zap.Sink, error) {
		query := url.Query()
		if query.Has("fileName") {
			fileName = query.Get("fileName")
		}
		if query.Has("maxSize") {
			maxsize, err := strconv.Atoi(query.Get("maxSize"))
			if err == nil {
				maxSize = maxsize
			}
		}
		if query.Has("maxBackups") {
			backups, err := strconv.Atoi(query.Get("maxBackups"))
			if err == nil {
				maxBackups = backups
			}
		}
		if query.Has("maxAge") {
			mage, err := strconv.Atoi(query.Get("maxAge"))
			if err == nil {
				maxAge = mage
			}
		}
		if query.Has("compress") {
			if query.Get("compress") == "true" {
				compress = true
			} else {
				compress = false
			}
		}
		if query.Has("localtime") {
			if query.Get("localtime") == "true" {
				localtime = true
			} else {
				localtime = false
			}
		}
		hook := &lumberjack.Logger{
			Filename:   fileName,
			MaxSize:    maxSize, // megabytes
			MaxBackups: maxBackups,
			MaxAge:     maxAge, // days
			Compress:   compress,
			LocalTime:  localtime,
		}
		return NewLumberjack(hook), nil
	}); err != nil {
		panic(err)
	}

	if err := zap.RegisterSink(RotateLogs, func(url *url.URL) (zap.Sink, error) {
		query := url.Query()

		if query.Has("fileName") {
			rfileName = query.Get("fileName")
			baseName = rfileName[:len(fileName)-len(".log")]
			rfileName = baseName + "_%Y-%m-%d %H:%M:%S" + ".log"
			baseName = baseName + ".log"
		}
		if query.Has("rotationTime") {
			rotationDurationStr := query.Get("rotationTime")
			rotationDuration, err := time.ParseDuration(rotationDurationStr)
			if err == nil {
				rotationTime = rotationDuration
			}
		}
		if query.Has("maxAge") {
			maxAgeDurationStr := query.Get("maxAge")
			maxAgeDuration, err := time.ParseDuration(maxAgeDurationStr)
			if err == nil {
				rmaxAge = maxAgeDuration
			}
		}

		if query.Has("localtime") {

			if query.Get("localtime") == "true" {
				rlocaltime = time.Local
				rclock = rotateloggers.Local
			} else {
				rlocaltime = time.UTC
				rclock = rotateloggers.UTC
			}

		}

		logs, err := rotateloggers.New(
			rfileName,
			rotateloggers.WithLinkName(baseName),
			rotateloggers.WithMaxAge(rmaxAge),
			rotateloggers.WithRotationTime(rotationTime),
			rotateloggers.WithLocation(rlocaltime),
			rotateloggers.WithClock(rclock),
		)
		if err != nil {
			return nil, err
		}
		return NewRotatelog(logs), nil
	}); err != nil {
		panic(err)
	}
}
