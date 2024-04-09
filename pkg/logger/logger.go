package logger

import (
	"carinfo/internal/config"
	"log/slog"
	"os"
)

const (
	envDebug = "debug"
	envError = "error"
	envInfo  = "info"
	envWarn  = "warn"
)

func newLogger() *slog.Logger {
	env := config.All.LoggerMode
	var log *slog.Logger

	switch env {
	case envDebug:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
				Level: slog.LevelDebug,
			}),
		)
	case envWarn:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
				Level: slog.LevelWarn,
			}),
		)
	case envError:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
				Level: slog.LevelError,
			}),
		)
	case envInfo:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
				Level: slog.LevelInfo,
			}),
		)
	}
	return log
}

var Log = newLogger()
