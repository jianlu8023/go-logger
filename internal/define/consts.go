package define

import (
	"time"

	rotateloggers "github.com/lestrrat-go/file-rotatelogs"
)

const (
	Lumberjack         = "lumberjack"
	LumberjackTemplate = "lumberjack:?fileName=%v&maxSize=%v&maxAge=%v&maxBackups=%v&compress=%v&localtime=%v"
)

// lumberjack 使用
var (
	FileName   = "./logs/lumberjack.log" // 日志文件路径
	MaxSize    = 5                       // MB
	MaxBackups = 7                       // 备份文件的最大数量
	MaxAge     = 30                      // 备份文件最大保留天数
	Compress   = false                   // 是否进行压缩
	Localtime  = false                   // 是否使用本地时间
)

const (
	RotateLogs         = "rotatelogs"
	RotateLogsTemplate = "rotatelogs:?fileName=%v&maxAge=%v&localtime=%v&rotationTime=%v"
)

// rotateLogs 使用
var (
	BaseName     = "./logs/rotatelogs.log"             // 软连接地址
	RfileName    = "./logs/rotatelogs.%Y-%m-%d-%H.log" // 日志文件记录格式
	RotationTime = 3 * time.Hour                       // 日志切割时间间隔
	RmaxAge      = 24 * time.Hour                      // 日志文件最大保留时间
	Rlocaltime   = time.UTC                            // 是否使用本地时间
	Rclock       = rotateloggers.UTC                   // clock时间
)
