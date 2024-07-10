package go_logger

import (
	"fmt"

	"github.com/jianlu8023/go-logger/pkg/df"
)

type Config struct {
	LogLevel    string `json:"logLevel,omitempty"`
	DevelopMode bool   `json:"developMode,omitempty"`
}

type LumberjackConfig struct {
	// 文件名
	FileName string `json:"fileName"`
	// MaxSize is the maximum size in megabytes of the log file before it gets
	// rotated. It defaults to 100 megabytes.
	MaxSize int `json:"maxSize"`
	// MaxBackups is the maximum number of old log files to retain.  The default
	// is to retain all old log files (though MaxAge may still cause them to get
	// deleted.)
	MaxBackups int `json:"maxBackups"`
	// MaxAge is the maximum number of days to retain old log files based on the
	// timestamp encoded in their filename.  Note that a day is defined as 24
	// hours and may not exactly correspond to calendar days due to daylight
	// savings, leap seconds, etc. The default is not to remove old log files
	// based on age.
	MaxAge int `json:"maxAge"`
	// Compress determines if the rotated log files should be compressed
	// using gzip. The default is not to perform compression.
	Compress bool `json:"compress"`
	// LocalTime determines if the time used for formatting the timestamps in
	// backup files is the computer's local time.  The default is to use UTC
	// time.
	Localtime bool `json:"localtime"`
}

func NewLumberjackUrl(config *LumberjackConfig) string {
	dst := make([]byte, len(df.LumberjackTemplate))
	copy(dst, df.LumberjackTemplate)
	if nil == config {
		return fmt.Sprintf(string(dst), df.FileName, df.MaxSize, df.MaxAge, df.MaxBackups,
			df.Compress, df.Localtime)
	} else {
		return fmt.Sprintf(string(dst), config.FileName, config.MaxSize, config.MaxAge, config.MaxBackups, config.Compress, config.Localtime)
	}
}

type RotateLogConfig struct {
	FileName     string `json:"fileName,omitempty"`
	MaxAge       string `json:"maxAge,omitempty"`
	LocalTime    bool   `json:"localTime,omitempty"`
	RotationTime string `json:"rotationTime,omitempty"`
}

func NewRotateLogURL(config *RotateLogConfig) string {
	dst := make([]byte, len(df.RotateLogsTemplate))
	copy(dst, df.RotateLogsTemplate)
	var (
		baseName     interface{}
		maxAge       interface{}
		localtime    interface{}
		rotationTime interface{}
	)

	if nil == config {
		baseName = df.BaseName
		maxAge = df.RmaxAge
		localtime = df.Rlocaltime
		rotationTime = df.RotationTime
	} else {
		if config.FileName != "" {
			baseName = config.FileName
		} else {
			baseName = df.BaseName
		}
		if config.MaxAge != "" {
			maxAge = config.MaxAge
		} else {
			maxAge = df.MaxAge
		}
		if config.LocalTime == true {
			localtime = true
		} else {
			localtime = false
		}
		if config.RotationTime != "" {
			rotationTime = config.RotationTime
		} else {
			rotationTime = df.RotationTime
		}
	}

	return fmt.Sprintf(string(dst), baseName, maxAge, localtime, rotationTime)
}
