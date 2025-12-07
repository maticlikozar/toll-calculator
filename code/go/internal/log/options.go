package log

import (
	"io"

	log "github.com/rs/zerolog"
)

const (
	LogPath        = "path"
	LogMethod      = "method"
	LogSummary     = "summary"
	LogPrincipalID = "principal_id"
	LogComponent   = "component"
	LogCreatedAt   = "created_at"
	LogFunction    = "function"
	LogID          = "id"
)

// Level type.
type Level int8

const (
	// PanicLevel panic level.
	PanicLevel = Level(log.PanicLevel)
	// FatalLevel fatal level.
	FatalLevel = Level(log.FatalLevel)
	// ErrorLevel error level.
	ErrorLevel = Level(log.ErrorLevel)
	// WarnLevel warn level.
	WarnLevel = Level(log.WarnLevel)
	// InfoLevel info level.
	InfoLevel = Level(log.InfoLevel)
	// DebugLevel debug level.
	DebugLevel = Level(log.DebugLevel)
	// TraceLevel trace level.
	TraceLevel = Level(log.TraceLevel)
)

// SetOutput func defines log output.
func SetOutput(writer io.Writer) {
	h.zerolog = h.zerolog.Output(writer)
}

// SetLevel func defines log level.
func SetLevel(level Level) {
	h.zerolog = h.zerolog.Level(log.Level(level))
}

// ParseLevel func parses log level from string.
func ParseLevel(lvl string) Level {
	level, err := log.ParseLevel(lvl)

	if err != nil {
		return TraceLevel
	}

	return Level(level)
}
