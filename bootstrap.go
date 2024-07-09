package go_logger

import (
	"net/url"
	"strconv"

	"go.uber.org/zap"
	"gopkg.in/natefinch/lumberjack.v2"
)

func init() {
	zap.RegisterSink(Lumberjack, func(url *url.URL) (zap.Sink, error) {
		query := url.Query()
		var (
			fileName   = "./logs/sdk.log"
			maxSize    = 5
			maxBackups = 7
			maxAge     = 30
			compress   = true
			localtime  = true
		)
		if query.Has("filename") {
			fileName = query.Get("filename")
		}
		if query.Has("maxsize") {
			atoi, err := strconv.Atoi(query.Get("maxsize"))
			if err == nil {
				maxSize = atoi
			}
		}
		if query.Has("maxbackups") {
			atoi, err := strconv.Atoi(query.Get("maxbackups"))
			if err == nil {
				maxBackups = atoi
			}
		}
		if query.Has("maxage") {
			atoi, err := strconv.Atoi(query.Get("maxage"))
			if err == nil {
				maxAge = atoi
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
	})
}
