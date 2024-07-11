package go_logger

import (
	"fmt"
	"time"

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

func LumberjackDefaultConfig() *LumberjackConfig {
	return &LumberjackConfig{
		FileName:   df.FileName,
		MaxSize:    df.MaxSize,
		MaxAge:     df.MaxAge,
		MaxBackups: df.MaxBackups,
		Compress:   df.Compress,
		Localtime:  df.Localtime,
	}
}

func NewLumberjackUrl(config *LumberjackConfig) string {
	dst := make([]byte, len(df.LumberjackTemplate))
	copy(dst, df.LumberjackTemplate)
	var (
		fileName   = df.FileName
		maxSize    = df.MaxSize
		maxAge     = df.MaxAge
		maxBackups = df.MaxBackups
		compress   = df.Compress
		localtime  = df.Localtime
	)
	if nil != config {
		if config.FileName != "" {
			fileName = config.FileName
		}
		if config.MaxSize != 0 {
			maxSize = config.MaxSize
		}
		if config.MaxAge != 0 {
			maxAge = config.MaxAge
		}
		if config.MaxBackups != 0 {
			maxBackups = config.MaxBackups
		}
		compress = config.Compress
		localtime = config.Localtime
	}

	return fmt.Sprintf(string(dst), fileName, maxSize, maxAge, maxBackups,
		compress, localtime)
}

type RotateLogConfig struct {
	FileName     string `json:"fileName,omitempty"`
	MaxAge       string `json:"maxAge,omitempty"`
	LocalTime    bool   `json:"localTime,omitempty"`
	RotationTime string `json:"rotationTime,omitempty"`
}

func RotateLogDefaultConfig() *RotateLogConfig {
	return &RotateLogConfig{
		FileName:     df.BaseName,
		MaxAge:       df.RmaxAge.String(),
		LocalTime:    false,
		RotationTime: df.RotationTime.String(),
	}
}

func NewRotateLogURL(config *RotateLogConfig) string {
	dst := make([]byte, len(df.RotateLogsTemplate))
	copy(dst, df.RotateLogsTemplate)
	var (
		baseName     = df.BaseName
		maxAge       = df.RmaxAge
		localtime    = df.Rlocaltime
		rotationTime = df.RotationTime
	)

	if nil != config {
		if config.FileName != "" {
			baseName = config.FileName
		}
		if config.MaxAge != "" {
			duration, err := time.ParseDuration(config.MaxAge)
			if err == nil {
				maxAge = duration
			}
			// maxAge = config.MaxAge
		}
		if config.LocalTime == true {
			localtime = time.Local
		} else {
			localtime = time.UTC
		}
		if config.RotationTime != "" {
			duration, err := time.ParseDuration(config.RotationTime)
			if err == nil {
				rotationTime = duration
			}
			// rotationTime = config.RotationTime
		}
	}
	return fmt.Sprintf(string(dst), baseName, maxAge, localtime, rotationTime)
}
