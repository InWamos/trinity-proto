package logger

import (
	"log/slog"
	"os"

	"github.com/InWamos/trinity-proto/config"
)

func GetLogger(loggerConfig *config.LoggingConfig) *slog.Logger {
	var logLevel slog.Level
	if err := logLevel.UnmarshalText([]byte(loggerConfig.Level)); err != nil {
		logLevel = slog.LevelInfo
	}
	opts := &slog.HandlerOptions{
		Level:     logLevel,
		AddSource: true,
	}

	handler := slog.NewJSONHandler(os.Stdout, opts)
	return slog.New(handler)
}
