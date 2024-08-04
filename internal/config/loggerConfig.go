package config

import (
	"os"

	"github.com/withmandala/go-log"
)

type LoggerConfig struct {
	LogHandler *log.Logger
}

func GetLogger() *log.Logger {

	// year, month, day := time.Now().UTC().Date()
	// today := fmt.Sprintf("%d-%s-%d", day, month.String(), year)
	// output := fmt.Sprintf("logs/%s.log", today)

	// f, err := os.Create(output)
	// if err != nil {
	// 	fmt.Println(err)
	// 	return nil
	// }

	manager := &LoggerConfig{}
	manager.LogHandler = log.
		New(os.Stdout).
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
