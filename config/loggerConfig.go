package config

import "os"
import "github.com/withmandala/go-log"

type LoggerConfig struct {
	LogHandler *log.Logger
}

func InitializeLogger() *LoggerConfig {
	manager := &LoggerConfig{}
	manager.LogHandler = log.New(os.Stderr)
	return manager
}

var Logger = InitializeLogger()
