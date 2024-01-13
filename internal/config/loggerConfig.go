package config

import (
	"os"

	"github.com/withmandala/go-log"
)

type LoggerConfig struct {
	LogHandler *log.Logger
}

var logger *LoggerConfig

func InitializeLogger() *LoggerConfig {

	if logger != nil {
		return logger
	}

	manager := &LoggerConfig{}
	manager.LogHandler = log.New(os.Stderr).WithDebug().WithColor().WithTimestamp().NoQuiet()
	logger = manager

	return manager
}
