package repository

import (
	"errors"
	"testing"
	"time"

	mock "github.com/stretchr/testify/mock"

	database "toll/internal/database/mocks"
	"toll/internal/test"
)

func TestGetConnectivity_Success(t *testing.T) {
	t.Parallel()

	license := "vk-123"

	from := time.Now()

	// Define mocked executions.
	events := database.NewMockDB(t)
	events.EXPECT().
		Select(t.Context(), mock.Anything, mock.Anything, license, from).
		Return(nil)

	// Create mocked event repository.
	repo := TollEvent(events)

	// Run GetAll() method.
	res, err := repo.GetAll(t.Context(), license, from)

	test.Match(t, res, err)
}

func TestGetConnectivity_Error(t *testing.T) {
	t.Parallel()

	license := "vk-123"

	from := time.Now()

	errTest := errors.New("test error")

	// Define mocked executions.
	events := database.NewMockDB(t)
	events.EXPECT().
		Select(t.Context(), mock.Anything, mock.Anything, license, from).
		Return(errTest)

	// Create mocked event repository.
	repo := TollEvent(events)

	// Run GetAll() method.
	_, err := repo.GetAll(t.Context(), license, from)

	test.Match(t, err)
}
