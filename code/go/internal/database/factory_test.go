package database

import (
	"database/sql"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Factory(t *testing.T) {
	t.Parallel()

	lock.Lock()

	factory := &factory{
		credentials: map[string]credentials{
			"test": {
				Driver: "mysql",
				DSN:    "mock-dsn",
			},
			"db": {
				Driver: "mysql",
				DSN:    "toll:toll@tcp(db:3306)/toll",
			},
			"conn-err": {
				Driver:    "mysql",
				DSN:       "mock-dsn",
				Connector: func() (*sql.DB, error) { return nil, errors.New("mock connection error") },
			},
		},
		instances: map[string]*db{
			"test": {},
		},
	}

	flags = &localFlags{
		db: map[string]*Config{
			"mysql": {
				Driver: "mysql",
				DSN:    "toll:toll@tcp(db:3306)/toll",
				TLS:    true,
			},
			"mysql-invalid-dsn": {
				Driver: "mysql",
				DSN:    "invalid-dsn",
			},
			"pgx": {
				Driver: "pgx",
				DSN:    "postgres://toll:toll@tolldb:5432/toll",
				TLS:    true,
			},
			"pgx-invalid-dsn": {
				Driver: "pgx",
				DSN:    "invalid-dsn",
			},
			"unrecognized-driver": {
				Driver: "unrecognized",
				DSN:    "mock-dsn",
			},
		},
	}

	defer lock.Unlock()

	Flags()

	// Test Get database from factory instances.
	_, err := factory.Get("test")
	assert.NoError(
		t,
		err,
	)

	// Test Get default database from factory credentials.
	_, err = factory.Get()
	assert.NoError(
		t,
		err,
	)

	// Test Get mysql database from flags credentials.
	_, err = factory.Get("mysql")
	assert.NoError(
		t,
		err,
	)

	// Test Get pgx database from flags credentials.
	_, err = factory.Get("pgx")
	assert.NoError(
		t,
		err,
	)

	// Test no database selected error.
	_, err = factory.Get("mysql", "pgx")
	assert.EqualError(
		t,
		err,
		"no database selected",
	)

	// Test unrecognized driver error.
	_, err = factory.Get("unrecognized-driver")
	assert.EqualError(
		t,
		err,
		"unrecognized database driver",
	)

	// Test connection error.
	_, err = factory.Get("conn-err")
	assert.EqualError(
		t,
		err,
		"mock connection error",
	)

	// Test invalid DSN error - mysql driver.
	_, err = factory.Get("mysql-invalid-dsn")
	assert.EqualError(
		t,
		err,
		"invalid DSN: missing the slash separating the database name",
	)

	// Test invalid DSN error - pgx driver.
	_, err = factory.Get("pgx-invalid-dsn")
	assert.EqualError(
		t,
		err,
		"cannot parse `invalid-dsn`: failed to parse as keyword/value (invalid keyword/value)",
	)
}
