package logger

import (
	"log/slog"
	"os"
)

// Logger wraps slog.Logger to allow future replacement.
type Logger struct {
	*slog.Logger
}

// New creates a structured JSON logger.
func New() *Logger {
	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	})

	return &Logger{
		Logger: slog.New(handler),
	}
}
