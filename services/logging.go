package services

import "os"
import "github.com/withmandala/go-log"

type LogManager struct {
	LogHandler *log.Logger
}

func InitializeLogger() *LogManager {
	manager := &LogManager{}
	manager.LogHandler = log.New(os.Stderr)
	return manager
}

var Logger = InitializeLogger()
