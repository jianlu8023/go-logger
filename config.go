package go_logger

import (
	"fmt"
	"time"

	rotateloggers "github.com/lestrrat-go/file-rotatelogs"
)

type Config struct {
	// Mode file or date or stdout
	// file use lumberjack
	// date use rotatelogs
	Mode        []string `json:"mode,omitempty"`
	LogLevel    string   `json:"logLevel,omitempty"`
	DevelopMode bool     `json:"developMode,omitempty"`
}

const (
	Lumberjack         = "lumberjack"
	lumberjackTemplate = "lumberjack:?fileName=%v&maxSize=%v&maxAge=%v&maxBackups=%v&compress=%v&localtime=%v"
)

// lumberjack 使用
var (
	fileName   = "./logs/lumberjack.log"
	maxSize    = 5
	maxBackups = 7
	maxAge     = 30
	compress   = true
	localtime  = true
)

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
	dst := make([]byte, len(lumberjackTemplate))
	copy(dst, lumberjackTemplate)
	if nil == config {
		return fmt.Sprintf(string(dst), fileName, maxSize, maxAge, maxBackups, compress, localtime)
	} else {
		return fmt.Sprintf(string(dst), config.FileName, config.MaxSize, config.MaxAge, config.MaxBackups, config.Compress, config.Localtime)
	}
}

const (
	RotateLogs         = "rotatelogs"
	rotateLogsTemplate = "rotatelogs:?fileName=%v&maxAge=%v&localtime=%v&rotationTime=%v"
)

// rotateLogs 使用
var (
	baseName     = "./logs/rotatelogs.log"
	rfileName    = "./logs/rotatelogs_%Y-%m-%d %H:%M:%S.log"
	rotationTime = 3 * time.Hour
	rmaxAge      = 24 * time.Hour
	rlocaltime   = time.Local
	rclock       = rotateloggers.Local
)

type RotateLogConfig struct {
	FileName     string `json:"fileName,omitempty"`
	MaxAge       string `json:"maxAge,omitempty"`
	LocalTime    bool   `json:"localTime,omitempty"`
	RotationTime string `json:"rotationTime,omitempty"`
}

func NewRotateLogURL(config *RotateLogConfig) string {
	dst := make([]byte, len(rotateLogsTemplate))
	copy(dst, rotateLogsTemplate)

	if nil == config {
		return fmt.Sprintf(string(dst), baseName, rmaxAge, rlocaltime, rotationTime)
	} else {

		return fmt.Sprintf(string(dst), config.FileName, config.MaxAge, config.LocalTime, config.RotationTime)
	}
}
