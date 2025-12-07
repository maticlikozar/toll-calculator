package database

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	lock sync.Mutex
)

func TestConfig_Validate(t *testing.T) {
	t.Parallel()

	// Test Validate with config not set.
	var mockConfig *Config

	err := mockConfig.Validate()
	assert.NoError(
		t,
		err,
	)

	// Test Validate with database driver not set error.
	mockConfig = &Config{}

	err = mockConfig.Validate()
	assert.EqualError(
		t,
		err,
		"database driver not set",
	)

	// Test Validate with database DSN not set error.
	mockConfig = &Config{
		Driver: "mock-driver",
	}

	err = mockConfig.Validate()
	assert.EqualError(
		t,
		err,
		"database DSN not set",
	)

	// Test Validate with config set.
	mockConfig = &Config{
		Driver: "mock-driver",
		DSN:    "mock-dsn",
	}

	err = mockConfig.Validate()
	assert.NoError(
		t,
		err,
	)
}

func TestConfig_Init_SetGlobalConfig(t *testing.T) {
	t.Parallel()

	lock.Lock()
	c := new(Config).Init("mock", "mock")
	defer lock.Unlock()

	// Returned config from init call.
	assert.NotNil(
		t,
		c,
	)

	// Global Config variable should be set.
	assert.NotNil(
		t,
		cfgs["mock"],
	)

	// Both the returned and global variable should be the same.
	assert.Equal(
		t,
		cfgs["mock"],
		c,
	)

	// Calling it again, should return the same.
	c2 := new(Config).Init("mock", "mock")

	// Both the returned and global variable should be the same.
	assert.Equal(
		t,
		cfgs["mock"],
		c2,
	)

	// Calling it again, should return the same.
	c3 := new(Config).Init("other-mock", "mock")

	// Both the returned and global variable should be the same.
	assert.Equal(
		t,
		cfgs["other-mock"],
		c3,
	)

	// Both should still be set.
	assert.NotNil(
		t,
		cfgs["other-mock"],
	)

	assert.NotNil(
		t,
		cfgs["mock"],
	)
}
