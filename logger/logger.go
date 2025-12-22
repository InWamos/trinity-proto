package logger

import (
	"log/slog"
	"os"

	"github.com/InWamos/trinity-proto/config"
)

func GetLogger(loggerConfig *config.LoggingConfig) *slog.Logger {
	var logLevel slog.Level
	var handler slog.Handler
	if err := logLevel.UnmarshalText([]byte(loggerConfig.Level)); err != nil {
		logLevel = slog.LevelInfo
	}
	// Plain text for debug env, JSON for prod env
	if logLevel == slog.LevelDebug {
		handler = slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level:     logLevel,
			AddSource: true,
		})
	} else {
		handler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: logLevel,
		})
	}

	return slog.New(handler)
}
