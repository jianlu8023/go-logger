package go_logger

import (
	"errors"
	"testing"
)

func TestNewLogger(t *testing.T) {
	logger := NewLogger(&Config{
		Mode:        []string{"stdout", "file", "date"},
		LogLevel:    "info",
		DevelopMode: true,
	})
	logger.Info("info log")

}

func TestNewSugaredLogger(t *testing.T) {

	logger := NewSugaredLogger(&Config{
		Mode: []string{
			"stdout", "file", "date",
		},
		LogLevel:    "info",
		DevelopMode: false,
	})

	logger.Infof("info %s", "log")
	logger.Warnf("warn %s", "log")

	logger.Infof("info struct %v", struct {
		Name string
		Age  int
	}{
		Name: "test",
		Age:  18,
	})

	logger.Errorf("error %s", errors.New("test error"))
}
