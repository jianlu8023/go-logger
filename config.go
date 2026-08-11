package go_logger

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/jianlu8023/go-logger/v2/internal/define"
)

type Config struct {
	LogLevel    string `json:"logLevel,omitempty"`
	DevelopMode bool   `json:"developMode,omitempty"`
	StackLevel  string `json:"stackLevel,omitempty"`
	ModuleName  string `json:"moduleName,omitempty"`
	Caller      bool   `json:"caller,omitempty"`
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

func (l *LumberjackConfig) String() string {
	marshalString, _ := json.Marshal(l)
	return string(marshalString)
}

func LumberjackDefaultConfig() *LumberjackConfig {
	return &LumberjackConfig{
		FileName:   define.GetLumberjackFileName(),
		MaxSize:    define.GetLumberjackMaxSize(),
		MaxAge:     define.GetLumberjackMaxAge(),
		MaxBackups: define.GetLumberjackMaxBackups(),
		Compress:   define.GetLumberjackCompress(),
		Localtime:  define.GetLumberjackLocaltime(),
	}
}

func NewLumberjackUrl(config *LumberjackConfig) string {
	var (
		fileName   = define.GetLumberjackFileName()
		maxSize    = define.GetLumberjackMaxSize()
		maxAge     = define.GetLumberjackMaxAge()
		maxBackups = define.GetLumberjackMaxBackups()
		compress   = define.GetLumberjackCompress()
		localtime  = define.GetLumberjackLocaltime()
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

	return fmt.Sprintf(define.LumberjackTemplate, fileName, maxSize, maxAge, maxBackups, compress, localtime)
}

type RotateLogConfig struct {
	FileName     string `json:"fileName,omitempty"`
	MaxAge       string `json:"maxAge,omitempty"`
	LocalTime    bool   `json:"localTime,omitempty"`
	RotationTime string `json:"rotationTime,omitempty"`
}

func (r *RotateLogConfig) String() string {
	marshalString, _ := json.Marshal(r)
	return string(marshalString)
}

func RotateLogDefaultConfig() *RotateLogConfig {
	return &RotateLogConfig{
		FileName:     define.GetRotatelogsBaseName(),
		MaxAge:       define.GetRotatelogsRmaxAge().String(),
		LocalTime:    false,
		RotationTime: define.GetRotatelogsRotationTime().String(),
	}
}

func NewRotateLogURL(config *RotateLogConfig) string {
	var (
		baseName     = define.GetRotatelogsBaseName()
		maxAge       = define.GetRotatelogsRmaxAge()
		localtime    = define.GetRotatelogsRlocaltime()
		rotationTime = define.GetRotatelogsRotationTime()
	)

	if nil != config {
		if config.FileName != "" {
			baseName = config.FileName
		}
		if config.MaxAge != "" {
			duration, err := time.ParseDuration(config.MaxAge)
			if err != nil {
				// 解析失败不 panic，输出 stderr 告警，回退到默认值
				_, _ = fmt.Fprintf(os.Stderr, "[go-logger] WARNING: invalid MaxAge %q: %v, using default\n", config.MaxAge, err)
			} else {
				maxAge = duration
			}
		}
		if config.LocalTime == true {
			localtime = time.Local
		} else {
			localtime = time.UTC
		}
		if config.RotationTime != "" {
			duration, err := time.ParseDuration(config.RotationTime)
			if err != nil {
				// 解析失败不 panic，输出 stderr 告警，回退到默认值
				_, _ = fmt.Fprintf(os.Stderr, "[go-logger] WARNING: invalid RotationTime %q: %v, using default\n", config.RotationTime, err)
			} else {
				rotationTime = duration
			}
		}
	}
	// 统一使用 true/false 格式化 localtime 参数，与 lumberjack 保持一致
	isLocal := localtime == time.Local
	return fmt.Sprintf(define.RotateLogsTemplate, baseName, maxAge, isLocal, rotationTime)
}
