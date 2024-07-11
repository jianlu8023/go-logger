package go_logger

import (
	"fmt"
	"testing"
)

func TestNewLumberjackUrl(t *testing.T) {
	fmt.Println(NewLumberjackUrl(&LumberjackConfig{
		FileName: "./logs/test.log",
	}))
}

func TestNewRotateLogURL(t *testing.T) {
	fmt.Println(NewRotateLogURL(&RotateLogConfig{
		LocalTime:    false,
		RotationTime: "4h",
		MaxAge:       "30d",
	}))
}

func TestRotateLogDefaultConfig(t *testing.T) {
	fmt.Println(RotateLogDefaultConfig())
}
