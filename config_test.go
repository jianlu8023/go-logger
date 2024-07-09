package go_logger

import (
	"fmt"
	"testing"
)

func TestNewLumberjackUrl(t *testing.T) {
	fmt.Println(NewLumberjackUrl(&LumberjackConfig{
		FileName:   "./logs/test.log",
		MaxAge:     5,
		MaxBackups: 7,
		MaxSize:    5,
		Localtime:  true,
		Compress:   true,
	}))
}

func TestNewRotateLogURL(t *testing.T) {
	fmt.Println(NewRotateLogURL(&RotateLogConfig{
		rfileName,
		"30d", true, "3h",
	}))
}
