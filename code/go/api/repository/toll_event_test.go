package repository

import (
	"errors"
	"testing"
	"time"

	"github.com/lib/pq"
	mock "github.com/stretchr/testify/mock"

	database "toll/internal/database/mocks"
	"toll/internal/test"
)

func TestGetAllForLicense_Success(t *testing.T) {
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

	// Run GetAllForLicense() method.
	res, err := repo.GetAllForLicense(t.Context(), license, from)

	test.Match(t, res, err)
}

func TestGetAllForLicense_Error(t *testing.T) {
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

	// Run GetAllForLicense() method.
	_, err := repo.GetAllForLicense(t.Context(), license, from)

	test.Match(t, err)
}

func TestGetAll_Success(t *testing.T) {
	t.Parallel()

	from := time.Now()

	// Define mocked executions.
	events := database.NewMockDB(t)
	events.EXPECT().
		Select(t.Context(), mock.Anything, mock.Anything, from, pq.Array([]string{"plate1", "plate2"})).
		Return(nil)

	// Create mocked event repository.
	repo := TollEvent(events)

	// Run GetAll() method.
	res, err := repo.GetAll(t.Context(), from, []string{"plate1", "plate2"})

	test.Match(t, res, err)
}

func TestGetAll_Error(t *testing.T) {
	t.Parallel()

	from := time.Now()

	errTest := errors.New("test error")

	// Define mocked executions.
	events := database.NewMockDB(t)
	events.EXPECT().
		Select(t.Context(), mock.Anything, mock.Anything, from, pq.Array([]string{"plate1", "plate2"})).
		Return(errTest)

	// Create mocked event repository.
	repo := TollEvent(events)

	// Run GetAll() method.
	_, err := repo.GetAll(t.Context(), from, []string{"plate1", "plate2"})

	test.Match(t, err)
}
