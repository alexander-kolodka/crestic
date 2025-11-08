package logger

import (
	"os"
	"time"

	"github.com/rs/zerolog"
)

type Format string

const (
	FormatColor Format = "color" // default: colored console output
	FormatCI    Format = "ci"    // plain text without colors
	FormatJSON  Format = "json"  // JSON output
)

// New creates a new logger with the specified format.
func New(format Format, level zerolog.Level) zerolog.Logger {
	if format == FormatJSON {
		return zerolog.New(os.Stdout).
			Level(level).
			With().
			Timestamp().
			Logger()
	}

	output := zerolog.ConsoleWriter{
		Out:        os.Stdout,
		TimeFormat: time.RFC3339,
		NoColor:    !hasColor(format),
	}

	return zerolog.New(output).
		Level(level).
		With().
		Timestamp().
		Logger()
}

func hasColor(format Format) bool {
	return FormatColor == format
}
