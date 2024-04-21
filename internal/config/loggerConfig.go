package config

import (
	"os"

	"github.com/withmandala/go-log"
)

type LoggerConfig struct {
	LogHandler *log.Logger
}

func GetLogger() *log.Logger {

	manager := &LoggerConfig{}
	manager.LogHandler = log.
		New(os.Stderr).
		WithDebug().
		WithColor().
		WithTimestamp().
		NoQuiet()

	defer func() {
		manager.LogHandler = nil
		manager = nil
	}()

	return manager.LogHandler
}
