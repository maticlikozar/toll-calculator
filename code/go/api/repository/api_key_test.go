package repository

import (
	"context"
	"errors"
	"testing"

	mock "github.com/stretchr/testify/mock"

	"toll/api/types"
	database "toll/internal/database/mocks"
	"toll/internal/test"
)

func TestApiKey_Get(t *testing.T) {
	t.Parallel()

	// Define mocked executions.
	keyHash := "mock-api-key-hash"
	apiKey := types.ApiKey{}
	ctx := context.Background()

	db := database.NewMockDB(t)
	db.
		EXPECT().
		Get(ctx, &apiKey, mock.Anything, keyHash).
		Run(func(_ context.Context, dest interface{}, query string, args ...interface{}) {
			e := dest.(*types.ApiKey)
			e.KeyHash = keyHash
		}).
		Return(nil)

	// Create repository with mocked dependencies.
	repo := ApiKey(db)

	// Run Get() function and make assertions.
	ret, err := repo.Get(ctx, keyHash)

	test.Match(t, keyHash, ret.KeyHash, err)
}

func TestApiKey_Get_Error(t *testing.T) {
	t.Parallel()

	// Define mocked executions.
	keyHash := "mock-api-key-hash"
	apiKey := types.ApiKey{}

	errTest := errors.New("mock get data from database error")

	db := database.NewMockDB(t)
	db.
		EXPECT().
		Get(t.Context(), &apiKey, mock.Anything, keyHash).
		Return(errTest)

	// Create repository with mocked dependencies.
	repo := ApiKey(db)

	// Run Get() function and make assertions.
	_, err := repo.Get(t.Context(), keyHash)
	
	test.Match(t, err)
}
