package logger

import (
	"log"
	"log/slog"
	"orgstructure/internal/config"
	"os"
)

func New(logConfig config.LogConfig) *slog.Logger {
	var logFile *os.File
	if logConfig.LogFile == "stdout" {
		logFile = os.Stdout
	} else {
		var err error
		logFile, err = os.Create(logConfig.LogFile)
		if err != nil {
			log.Fatalf("can't create log file: %s", err)
		}
	}

	var logLevel slog.Level
	switch logConfig.LogLevel {
	case "debug":
		logLevel = slog.LevelDebug
	case "info":
		logLevel = slog.LevelInfo
	case "warn":
		logLevel = slog.LevelWarn
	case "error":
		logLevel = slog.LevelError
	default:
		log.Fatal("incorrect log level")
	}

	return slog.New(slog.NewTextHandler(logFile, &slog.HandlerOptions{Level: logLevel}))
}
