package database

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Flags(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	lock.Lock()
	flags = nil
	defer lock.Unlock()

	// Uninitialized.
	assert.Nil(
		flags,
		"uninitialized flags should be nil",
	)
	assert.EqualError(
		flags.Validate("db"),
		"database flags validation error",
		"expected validation error",
	)

	// Init without flags.
	Flags("test")
	assert.NotNil(
		flags,
		"initialized flags should not be nil",
	)
	assert.EqualError(
		flags.Validate("db"),
		"invalid flags for database db: database DSN not set",
		"expected DSN validation error",
	)

	// Hack for env variables.
	flags.db["db"].DSN = "test"

	assert.NotNil(
		flags,
		"initialized flags should not be nil",
	)
	assert.NoError(
		flags.Validate("db"),
		"expected no validation error",
	)

	// Flags already initialized.
	Flags("test")
	assert.NotNil(
		flags,
		"initialized flags should not be nil",
	)
	assert.NoError(
		flags.Validate("db"),
		"expected no validation error",
	)

	// Flags already initialized.
	assert.NotNil(
		flags,
		"initialized flags should not be nil",
	)
	assert.EqualError(
		flags.Validate("test"),
		"no flags found for database: test",
		"expected DSN validation error",
	)
}
