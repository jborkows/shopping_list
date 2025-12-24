package logging

import (
	"log/slog"
	"os"
)

func New() *slog.Logger {
	level := new(slog.LevelVar)
	level.Set(slog.LevelInfo)

	switch os.Getenv("LOG_LEVEL") {
	case "debug":
		level.Set(slog.LevelDebug)
	case "warn":
		level.Set(slog.LevelWarn)
	case "error":
		level.Set(slog.LevelError)
	}

	return slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: level}))
}
