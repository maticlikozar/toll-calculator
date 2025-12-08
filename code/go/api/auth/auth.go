package auth

import (
	"context"

	"github.com/google/uuid"

	apierr "toll/api/handler/errors"
	"toll/api/identity"
	"toll/api/restapi"
	"toll/api/service"
	"toll/api/types"

	"toll/internal/errlog"
	"toll/internal/log"
)

type (
	Authorization interface {
		HandleApiKeyAuth(ctx context.Context, operationName string, t restapi.ApiKeyAuth) (context.Context, error)
	}

	authorizer struct {
		auth service.AuthService
		log  log.Logger
	}
)

var (
	handle *authorizer
)

// Get func returns Authorization Provider handle with configurations.
func Get() Authorization {
	if handle == nil {
		// Create Authorization provider handle with auth service.
		handle = &authorizer{
			auth: service.Authorization,
			log:  log.WithField(types.LogComponent, "api/auth"),
		}
	}

	return handle
}

// HandleApiKeyAuth implements openapi spec security definition for API Key Auth.
func (op *authorizer) HandleApiKeyAuth(ctx context.Context, operationName string, t restapi.ApiKeyAuth) (context.Context, error) {
	ctx, kid, err := op.authApiKey(ctx, t.APIKey)
	if err != nil {
		log.WithFields(errlog.StackLog(err)).Warne(err, "validate api key")

		return ctx, apierr.ErrAPIUnauthorized
	}

	return identity.Set(ctx, kid), nil
}

func (op *authorizer) authApiKey(ctx context.Context, token string) (context.Context, *uuid.UUID, error) {
	kid, err := op.auth.ValidateApiKey(ctx, token)
	if err != nil {
		return ctx, nil, apierr.ErrAPIUnauthorized
	}

	if kid == nil || *kid == uuid.Nil {
		return ctx, nil, apierr.ErrAPIUnauthorized
	}

	return ctx, kid, nil
}
