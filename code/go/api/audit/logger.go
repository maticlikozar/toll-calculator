package audit

import (
	"context"
	"os"

	"github.com/google/uuid"
	"github.com/ogen-go/ogen/middleware"
	"github.com/rs/zerolog"

	"toll/api/identity"

	"toll/internal/request"
)

const (
	logLevel = "audit"
)

type (
	// Logger represents an audit logger.
	Logger interface {
		Log(ctx context.Context, method string)
		LogP(ctx context.Context, method string, params Parameters)
	}

	logger struct {
		zerolog zerolog.Logger
	}

	Parameters map[string]string
)

var (
	h = newLogger()
)

// newLogger creates instance with default options.
func newLogger() *logger {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnixMs

	lg := zerolog.New(os.Stdout).With().Timestamp().Str("level", logLevel).Logger()

	return &logger{
		zerolog: lg,
	}
}

// Log produces an audit log entry.
func Log(ctx context.Context, method string) {
	h.Log(ctx, method)
}

// LogP produces an audit log entry with additional parameters.
func LogP(ctx context.Context, method string, params Parameters) {
	h.LogP(ctx, method, params)
}

// Log produces an audit log entry.
func (l logger) Log(ctx context.Context, method string) {
	event := l.zerolog.Log()

	// Get request id from context.
	event = event.Str("request_id", request.GetReqIdCtx(ctx))

	// Add user identity to audit logs.
	var kid *uuid.UUID

	if ctx != nil {
		kid = identity.Get(ctx)
	}

	if kid != nil {
		event = event.Str("key_id", kid.String())
	}

	// Add method to log.
	event = event.Str("method", method)

	event.Send()
}

// LogP produces an audit log entry with additional parameters.
func (l logger) LogP(ctx context.Context, method string, params Parameters) {
	event := l.zerolog.Log()

	// Get request id from context.
	event = event.Str("request_id", request.GetReqIdCtx(ctx))

	// Add user identity to audit logs.
	var kid *uuid.UUID

	if ctx != nil {
		kid = identity.Get(ctx)
	}

	if kid != nil {
		event = event.Str("key_id", kid.String())
	}

	// Add method to log.
	event = event.Str("method", method)

	// Add additional parameters.
	for k, v := range params {
		event = event.Str(k, v)
	}

	event.Send()
}

// Middleware returns ogen middleware interface.
func Middleware(req middleware.Request, next middleware.Next) (middleware.Response, error) {
	resp, err := next(req)

	Log(req.Context, req.Raw.Method)

	return resp, err
}
