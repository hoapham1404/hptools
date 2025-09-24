package logging

import (
	"io"
	"log/slog"
	"os"
	"strings"

	"hptools/internal/config"
)

// NewLogger creates a new structured logger based on configuration
func NewLogger(cfg *config.Config) *slog.Logger {
	var level slog.Level
	switch strings.ToLower(cfg.Log.Level) {
	case "debug":
		level = slog.LevelDebug
	case "info":
		level = slog.LevelInfo
	case "warn", "warning":
		level = slog.LevelWarn
	case "error":
		level = slog.LevelError
	default:
		level = slog.LevelInfo
	}

	var handler slog.Handler
	var writer io.Writer = os.Stdout

	opts := &slog.HandlerOptions{
		Level: level,
	}

	switch strings.ToLower(cfg.Log.Format) {
	case "json":
		handler = slog.NewJSONHandler(writer, opts)
	default:
		handler = slog.NewTextHandler(writer, opts)
	}

	return slog.New(handler)
}

// WithComponent adds component information to the logger
func WithComponent(logger *slog.Logger, component string) *slog.Logger {
	return logger.With("component", component)
}
