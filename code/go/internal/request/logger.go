package request

import (
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/rs/zerolog"
)

const (
	logLevel = "request"
)

type (
	logger struct {
		zerolog zerolog.Logger
	}
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

// Logger returns a request logger handler.
func Logger(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

		t1 := time.Now()
		defer func() {
			event := h.zerolog.Log()
			event = event.Str("request_id", GetReqIdCtx(r.Context()))
			event = event.Str("host", r.Host)
			event = event.Str("uri", r.RequestURI)
			event = event.Str("proto", r.Proto)
			event = event.Str("method", r.Method)
			event = event.Int("status", ww.Status())
			event = event.Int("bytes", ww.BytesWritten())
			event = event.Float64("duration", time.Since(t1).Seconds())
			event = event.Str("from", r.RemoteAddr)

			event.Send()
		}()

		next.ServeHTTP(ww, r)
	}

	return http.HandlerFunc(fn)
}
