package go_logger

import (
	"os"
	"os/signal"
	"syscall"
	"testing"
	"time"
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

	ticker := time.NewTicker(time.Second * 10)

	go func() {
		for {
			select {
			case <-ticker.C:
				logger.Infof("info %s", "log")
				logger.Warnf("warn %s", "log")

				logger.Infof("info struct %v", struct {
					Name string
					Age  int
				}{
					Name: "test",
					Age:  18,
				})

				// logger.Errorf("error %s", errors.New("test error"))
			}
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	select {
	case <-quit:
		ticker.Stop()
		logger.Info("stop")
	}

}
