package df

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
	FileName   = "./logs/lumberjack.log"
	MaxSize    = 5
	MaxBackups = 7
	MaxAge     = 30
	Compress   = true
	Localtime  = true
)

const (
	RotateLogs         = "rotatelogs"
	RotateLogsTemplate = "rotatelogs:?fileName=%v&maxAge=%v&localtime=%v&rotationTime=%v"
)

// rotateLogs 使用
var (
	BaseName     = "./logs/rotatelogs.log"
	RfileName    = "./logs/rotatelogs_%Y-%m-%d %H:%M:%S.log"
	RotationTime = 3 * time.Hour
	RmaxAge      = 24 * time.Hour
	Rlocaltime   = time.Local
	Rclock       = rotateloggers.Local
)
