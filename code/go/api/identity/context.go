package identity

import (
	"context"

	"github.com/google/uuid"
)

type authKey string

const (
	// Server principal keys.
	ServerPrincipalKey authKey = "ServerPrincipal"
)

// Get func returns signed api key id from the context.
func Get(ctx context.Context) *uuid.UUID {
	ret, ok := ctx.Value(ServerPrincipalKey).(*uuid.UUID)
	if !ok {
		return nil
	}

	return ret
}

// Set func writes kid to context.
func Set(ctx context.Context, kid *uuid.UUID) context.Context {
	return context.WithValue(ctx, ServerPrincipalKey, kid)
}
