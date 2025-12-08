package request

import (
	"context"
	"net/http"

	"github.com/google/uuid"
)

type ctxKey string

const (
	requestIdKey ctxKey = "request_id"
	requestId    string = "X-Request-ID"
)

// RequestId middleware generates or propagates X-Request-ID.
func RequestId(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rid := r.Header.Get(requestId)

		// If no request ID present, generate one.
		if rid == "" {
			rid = uuid.New().String()
			r.Header.Set(requestId, rid)
		}

		ctx := context.WithValue(r.Context(), requestIdKey, rid)
		h.ServeHTTP(w, r.WithContext(ctx))
	})
}

func GetReqId(r *http.Request) string {
	if r == nil || r.Header == nil {
		return ""
	}

	return r.Header.Get(requestId)
}

func GetReqIdCtx(ctx context.Context) string {
	val := ctx.Value(requestIdKey)

	rid, ok := val.(string)
	if !ok {
		return ""
	}

	return rid
}
