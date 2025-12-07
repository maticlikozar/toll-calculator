package errors

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewAPIError(t *testing.T) {
	t.Parallel()

	err := NewAPIError(500, "internal")
	assert.NotNil(t, err)

	assert.Equal(t, 500, err.StatusCode)
	assert.Equal(t, "internal", err.Title)
	assert.Equal(t, "", err.Detail)
	assert.Equal(t, 0, err.ErrorCode)
}

func TestWithTitle(t *testing.T) {
	t.Parallel()

	err := NewAPIError(500, "internal")
	err = err.WithTitle("not-internal")

	assert.Equal(t, 500, err.StatusCode)
	assert.Equal(t, "not-internal", err.Title)
	assert.Equal(t, "", err.Detail)
	assert.Equal(t, 0, err.ErrorCode)
}

func TestWithDetails(t *testing.T) {
	t.Parallel()

	err := NewAPIError(500, "internal")
	err = err.WithDetails("details")

	assert.Equal(t, 500, err.StatusCode)
	assert.Equal(t, "internal", err.Title)
	assert.Equal(t, "details", err.Detail)
	assert.Equal(t, 0, err.ErrorCode)
}

func TestErrorCode(t *testing.T) {
	t.Parallel()

	err := NewAPIError(500, "internal")
	err = err.WithErrorCode(99)

	assert.Equal(t, 500, err.StatusCode)
	assert.Equal(t, "internal", err.Title)
	assert.Equal(t, "", err.Detail)
	assert.Equal(t, 99, err.ErrorCode)
}
