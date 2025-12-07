package service

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"time"

	"github.com/google/uuid"
	"github.com/pkg/errors"

	"toll/api/repository"

	"toll/internal/database"
	"toll/internal/errlog"
)

var (
	// ErrIncorrectApiKey is returned when api key is not found in the database.
	ErrIncorrectApiKey = errors.New("incorrect api key")

	// ErrApiKeyExpired is returned when api key is expired.
	ErrApiKeyExpired = errors.New("api key is expired")

	// ErrUserIdNotSet is returned when user id is missing in api key object.
	ErrApiKeyUserIdNotSet = errors.New("missing user id in api key")
)

type (
	// AuthService interface with method definitions.
	AuthService interface {
		ValidateApiKey(ctx context.Context, key string) (*uuid.UUID, error)
	}

	auth struct {
		keys repository.ApiKeyRepository
	}
)

// Auth func returns new AuthService with new background context.
func Auth() AuthService {
	db := database.Get()

	return &auth{
		keys: repository.ApiKey(db),
	}
}

func (svc *auth) ValidateApiKey(ctx context.Context, key string) (*uuid.UUID, error) {
	h := sha256.New()
	h.Write([]byte(key))

	keyHash := base64.RawURLEncoding.EncodeToString(h.Sum(nil))

	apiKey, err := svc.keys.Get(ctx, keyHash)
	if err != nil {
		return nil, errlog.Error(err)
	}

	if apiKey == nil || apiKey.Id == uuid.Nil {
		return nil, errlog.Error(ErrIncorrectApiKey)
	}

	if apiKey.ExpiresAt.Before(time.Now()) {
		return nil, errlog.Error(ErrApiKeyExpired)
	}

	return &apiKey.Id, nil
}
