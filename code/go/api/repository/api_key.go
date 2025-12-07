package repository

import (
	"context"

	"toll/api/types"

	"toll/internal/database"
	"toll/internal/errlog"
)

type (
	// ApiKeyRepository interface with method definitions.
	ApiKeyRepository interface {
		Get(ctx context.Context, keyHash string) (*types.ApiKey, error)
	}

	apiKey struct {
		db database.DB
	}
)

// ApiKey func returns ApiKeyRepository with provided database connection.
func ApiKey(db database.DB) ApiKeyRepository {
	return &apiKey{db: db}
}

// Get func returns an api key for the provided key hash value.
func (r *apiKey) Get(ctx context.Context, keyHash string) (*types.ApiKey, error) {
	query := `
		SELECT
			id,
			key_hash,
			expires_at
		FROM api_key
		WHERE key_hash=$1`

	ret := types.ApiKey{}

	err := r.db.Get(ctx, &ret, query, keyHash)
	if err != nil {
		return nil, errlog.Error(err)
	}

	return &ret, nil
}
