package log

import (
	"fmt"
	"os"

	"github.com/rs/zerolog"
)

type (
	// Logger interface definition.
	Logger interface {
		WithField(key string, val interface{}) Logger
		WithFields(fields Fields) Logger

		Trace(args ...interface{})
		Tracef(format string, args ...interface{})
		Tracee(err error, args ...interface{})

		Debug(args ...interface{})
		Debugf(format string, args ...interface{})
		Debuge(err error, args ...interface{})

		Info(args ...interface{})
		Infof(format string, args ...interface{})
		Infoe(err error, args ...interface{})

		Warn(args ...interface{})
		Warnf(format string, args ...interface{})
		Warne(err error, args ...interface{})

		Error(args ...interface{})
		Errorf(format string, args ...interface{})
		Errore(err error, args ...interface{})

		Fatal(args ...interface{})
		Fatalf(format string, args ...interface{})
		Fatale(err error, args ...interface{})

		Panic(args ...interface{})
		Panicf(format string, args ...interface{})
		Panice(err error, args ...interface{})

		Print(args ...interface{})
		Printf(format string, args ...interface{})
		Printe(err error, args ...interface{})
	}

	// logger struct definition.
	logger struct {
		fields  map[string]interface{}
		zerolog zerolog.Logger
	}

	// Fields definition.
	Fields map[string]interface{}
)

// Global logger instance.
var (
	h = newLogger()
)

func Init() {
	lvl := ParseLevel(flags.log.Level)

	SetLevel(lvl)
}

// newLogger creates instance with default options.
func newLogger() *logger {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnixMs

	lg := zerolog.New(os.Stdout).With().Timestamp().Logger()
	lg.Level(zerolog.TraceLevel)

	// Create logger.
	return &logger{
		fields:  map[string]interface{}{},
		zerolog: lg,
	}
}

// copyLogger returns a deep copy of a logger in a new instance.
func copyLogger(l *logger) *logger {
	nl := &logger{
		fields:  make(map[string]interface{}, len(l.fields)),
		zerolog: l.zerolog,
	}

	for key, val := range l.fields {
		nl.fields[key] = val
	}

	return nl
}

// Noop returns a Logger instance for which all operations are no-op.
func Noop() Logger {
	return &logger{
		fields:  map[string]interface{}{},
		zerolog: zerolog.Nop(),
	}
}

// SetTimeFormat sets the time format used in the log.
// The format must be compatible with the layout supported by Go's standard time package.
func SetTimeFormat(layout string) {
	zerolog.TimeFieldFormat = layout
}

// WithField func returns logger with additional field.
func WithField(key string, val interface{}) Logger {
	t := &logger{
		fields:  make(map[string]interface{}, 1),
		zerolog: h.zerolog,
	}

	t.fields[key] = val

	return t
}

// WithFields func returns logger with additional fields.
func WithFields(fields Fields) Logger {
	t := &logger{
		fields:  make(map[string]interface{}, len(fields)),
		zerolog: h.zerolog,
	}

	for key, val := range fields {
		t.fields[key] = val
	}

	return t
}

// Trace func writes trace log.
func Trace(args ...interface{}) {
	h.zerolog.Trace().Fields(h.fields).Msg(fmt.Sprint(args...))
}

// Tracef func writes formatted trace log.
func Tracef(format string, args ...interface{}) {
	h.zerolog.Trace().Fields(h.fields).Msg(fmt.Sprintf(format, args...))
}

// Tracee func writes trace log with error.
func Tracee(err error, args ...interface{}) {
	h.zerolog.Trace().Err(err).Fields(h.fields).Msg(fmt.Sprint(args...))
}

// Debug func writes debug log.
func Debug(args ...interface{}) {
	h.zerolog.Debug().Fields(h.fields).Msg(fmt.Sprint(args...))
}

// Debugf func writes formatted debug log.
func Debugf(format string, args ...interface{}) {
	h.zerolog.Debug().Fields(h.fields).Msg(fmt.Sprintf(format, args...))
}

// Debuge func writes debug log with error.
func Debuge(err error, args ...interface{}) {
	h.zerolog.Debug().Err(err).Fields(h.fields).Msg(fmt.Sprint(args...))
}

// Info func writes info log.
func Info(args ...interface{}) {
	h.zerolog.Info().Fields(h.fields).Msg(fmt.Sprint(args...))
}

// Infof func writes formatted info log.
func Infof(format string, args ...interface{}) {
	h.zerolog.Info().Fields(h.fields).Msg(fmt.Sprintf(format, args...))
}

// Infoe func writes info log with error.
func Infoe(err error, args ...interface{}) {
	h.zerolog.Info().Err(err).Fields(h.fields).Msg(fmt.Sprint(args...))
}

// Warn func writes warn log.
func Warn(args ...interface{}) {
	h.zerolog.Warn().Fields(h.fields).Msg(fmt.Sprint(args...))
}

// Warnf func writes formatted warn log.
func Warnf(format string, args ...interface{}) {
	h.zerolog.Warn().Fields(h.fields).Msg(fmt.Sprintf(format, args...))
}

// Warne func writes warn log with error.
func Warne(err error, args ...interface{}) {
	h.zerolog.Warn().Err(err).Fields(h.fields).Msg(fmt.Sprint(args...))
}

// Error func writes error log.
func Error(args ...interface{}) {
	h.zerolog.Error().Fields(h.fields).Msg(fmt.Sprint(args...))
}

// Errorf func writes formatted error log.
func Errorf(format string, args ...interface{}) {
	h.zerolog.Error().Fields(h.fields).Msg(fmt.Sprintf(format, args...))
}

// Errore func writes error log with error.
func Errore(err error, args ...interface{}) {
	h.zerolog.Error().Err(err).Fields(h.fields).Msg(fmt.Sprint(args...))
}

// Fatal func writes fatal log.
func Fatal(args ...interface{}) {
	h.zerolog.Fatal().Fields(h.fields).Msg(fmt.Sprint(args...))
}

// Fatalf func writes formatted fatal log.
func Fatalf(format string, args ...interface{}) {
	h.zerolog.Fatal().Fields(h.fields).Msg(fmt.Sprintf(format, args...))
}

// Fatale func writes fatal log with error.
func Fatale(err error, args ...interface{}) {
	h.zerolog.Fatal().Err(err).Fields(h.fields).Msg(fmt.Sprint(args...))
}

// Panic func writes panic log.
func Panic(args ...interface{}) {
	h.zerolog.Panic().Fields(h.fields).Msg(fmt.Sprint(args...))
}

// Panicf func writes formatted panic log.
func Panicf(format string, args ...interface{}) {
	h.zerolog.Panic().Fields(h.fields).Msg(fmt.Sprintf(format, args...))
}

// Panice func writes panic log with error.
func Panice(err error, args ...interface{}) {
	h.zerolog.Panic().Err(err).Fields(h.fields).Msg(fmt.Sprint(args...))
}

// Print func prints log.
func Print(args ...interface{}) {
	h.zerolog.Print(args...)
}

// Printf func prints formatted log.
func Printf(format string, args ...interface{}) {
	h.zerolog.Printf(format, args...)
}

// Printe func writes debug log with error.
func Printe(err error, args ...interface{}) {
	h.zerolog.Debug().Err(err).Msg(fmt.Sprint(args...))
}

// WithField func returns logger with additional field.
func (l *logger) WithField(key string, val interface{}) Logger {
	nl := copyLogger(l)

	nl.fields[key] = val

	return nl
}

// WithFields func returns logger with additional fields.
func (l *logger) WithFields(fields Fields) Logger {
	nl := copyLogger(l)

	for key, val := range fields {
		nl.fields[key] = val
	}

	return nl
}

// Trace func writes trace log.
func (l *logger) Trace(args ...interface{}) {
	l.zerolog.Trace().Fields(l.fields).Msg(fmt.Sprint(args...))
}

// Tracef func writes formatted trace log.
func (l *logger) Tracef(format string, args ...interface{}) {
	l.zerolog.Trace().Fields(l.fields).Msg(fmt.Sprintf(format, args...))
}

// Tracee func writes trace log with error.
func (l *logger) Tracee(err error, args ...interface{}) {
	l.zerolog.Trace().Err(err).Fields(l.fields).Msg(fmt.Sprint(args...))
}

// Debug func writes debug log.
func (l *logger) Debug(args ...interface{}) {
	l.zerolog.Debug().Fields(l.fields).Msg(fmt.Sprint(args...))
}

// Debugf func writes formatted debug log.
func (l *logger) Debugf(format string, args ...interface{}) {
	l.zerolog.Debug().Fields(l.fields).Msg(fmt.Sprintf(format, args...))
}

// Debuge func writes debug log with error.
func (l *logger) Debuge(err error, args ...interface{}) {
	l.zerolog.Debug().Err(err).Fields(l.fields).Msg(fmt.Sprint(args...))
}

// Info func writes info log.
func (l *logger) Info(args ...interface{}) {
	l.zerolog.Info().Fields(l.fields).Msg(fmt.Sprint(args...))
}

// Infof func writes formatted info log.
func (l *logger) Infof(format string, args ...interface{}) {
	l.zerolog.Info().Fields(l.fields).Msg(fmt.Sprintf(format, args...))
}

// Infoe func writes info log with error.
func (l *logger) Infoe(err error, args ...interface{}) {
	l.zerolog.Info().Err(err).Fields(l.fields).Msg(fmt.Sprint(args...))
}

// Warn func writes warn log.
func (l *logger) Warn(args ...interface{}) {
	l.zerolog.Warn().Fields(l.fields).Msg(fmt.Sprint(args...))
}

// Warnf func writes formatted warn log.
func (l *logger) Warnf(format string, args ...interface{}) {
	l.zerolog.Warn().Fields(l.fields).Msg(fmt.Sprintf(format, args...))
}

// Warne func writes warn log with error.
func (l *logger) Warne(err error, args ...interface{}) {
	l.zerolog.Warn().Err(err).Fields(l.fields).Msg(fmt.Sprint(args...))
}

// Error func writes error log.
func (l *logger) Error(args ...interface{}) {
	l.zerolog.Error().Fields(l.fields).Msg(fmt.Sprint(args...))
}

// Errorf func writes formatted error log.
func (l *logger) Errorf(format string, args ...interface{}) {
	l.zerolog.Error().Fields(l.fields).Msg(fmt.Sprintf(format, args...))
}

// Errore func writes error log with error.
func (l *logger) Errore(err error, args ...interface{}) {
	l.zerolog.Error().Err(err).Fields(l.fields).Msg(fmt.Sprint(args...))
}

// Fatal func writes fatal log.
func (l *logger) Fatal(args ...interface{}) {
	l.zerolog.Fatal().Fields(l.fields).Msg(fmt.Sprint(args...))
}

// Fatalf func writes formatted fatal log.
func (l *logger) Fatalf(format string, args ...interface{}) {
	l.zerolog.Fatal().Fields(l.fields).Msg(fmt.Sprintf(format, args...))
}

// Fatale func writes fatal log with error.
func (l *logger) Fatale(err error, args ...interface{}) {
	l.zerolog.Fatal().Err(err).Fields(l.fields).Msg(fmt.Sprint(args...))
}

// Panic func writes panic log.
func (l *logger) Panic(args ...interface{}) {
	l.zerolog.Panic().Fields(l.fields).Msg(fmt.Sprint(args...))
}

// Panicf func writes formatted panic log.
func (l *logger) Panicf(format string, args ...interface{}) {
	l.zerolog.Panic().Fields(l.fields).Msg(fmt.Sprintf(format, args...))
}

// Panice func writes panic log with error.
func (l *logger) Panice(err error, args ...interface{}) {
	l.zerolog.Panic().Err(err).Fields(l.fields).Msg(fmt.Sprint(args...))
}

// Print func prints log.
func (l *logger) Print(args ...interface{}) {
	l.zerolog.Debug().Fields(l.fields).Msg(fmt.Sprint(args...))
}

// Printf func prints formatted log.
func (l *logger) Printf(format string, args ...interface{}) {
	l.zerolog.Debug().Fields(l.fields).Msg(fmt.Sprintf(format, args...))
}

// Printe func writes debug log with error.
func (l *logger) Printe(err error, args ...interface{}) {
	l.zerolog.Debug().Err(err).Fields(l.fields).Msg(fmt.Sprint(args...))
}
