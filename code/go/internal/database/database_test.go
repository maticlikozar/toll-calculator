package database

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"

	"toll/internal/errlog"
)

func initSQLmock(name string) (DB, *sql.DB, sqlmock.Sqlmock, error) {
	// Create mocked database.
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		return nil, nil, nil, err
	}

	// Create connection to mocked database with sqlx driver.
	sqlxDB := sqlx.NewDb(mockDB, "sqlmock")
	perf := NewPerformanceMetrics(name)
	// Create database handlers.
	db := &db{
		sqlxDB,
		perf,
	}

	return db, mockDB, mock, err
}

func Test_Transaction_Begin(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	// Get mocked database.
	db, mockDB, mock, err := initSQLmock("begin")
	assert.NoError(
		t,
		err,
		"error creating mocked database, got = %v", err,
	)

	defer mockDB.Close()

	// Begin transaction.
	mock.ExpectBegin()
	mock.ExpectCommit()

	err = db.Transaction(ctx, func(ctx context.Context, tx TX) error {
		return nil
	})

	assert.NoError(
		t,
		err,
		"error starting transaction, got = %v", err,
	)
}

func Test_Transaction_Begin_Err(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	// Get mocked database.
	db, mockDB, mock, err := initSQLmock("begin_err")
	assert.NoError(
		t,
		err,
		"error creating mocked database, got = %v", err,
	)

	defer mockDB.Close()

	// Begin transaction with error.
	errNoTransactions := errlog.New("error starting transaction")
	mock.
		ExpectBegin().
		WillReturnError(errNoTransactions)

	err = db.Transaction(ctx, func(ctx context.Context, tx TX) error {
		return nil
	})
	assert.EqualError(
		t,
		err,
		errNoTransactions.Error(),
		"expected transaction error, got = %v", err,
	)
}

func Test_Transaction_Commit_Err(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	// Get mocked database.
	db, mockDB, mock, err := initSQLmock("commit_err")
	assert.NoError(
		t,
		err,
		"error creating mocked database, got = %v", err,
	)

	defer mockDB.Close()

	// Commit transaction with error.
	errOnCommit := errlog.New("error committing transaction")

	mock.
		ExpectBegin()
	mock.
		ExpectCommit().
		WillReturnError(errOnCommit)

	err = db.Transaction(ctx, func(ctx context.Context, tx TX) error {
		return nil
	})
	assert.EqualError(
		t,
		err,
		errOnCommit.Error(),
		"expected commit error, got = %v", err,
	)
}

func Test_Transaction_Rollback(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	// Get mocked database.
	db, mockDB, mock, err := initSQLmock("rollback")
	assert.NoError(
		t,
		err,
		"error creating mocked database, got = %v", err,
	)

	defer mockDB.Close()

	// Commit transaction with error.
	errTest := errlog.New("test error")

	mock.ExpectBegin()
	mock.ExpectRollback()

	err = db.Transaction(ctx, func(ctx context.Context, tx TX) error {
		return errTest
	})
	assert.EqualError(
		t,
		err,
		errTest.Error(),
		"expected expected test error, got = %v", err,
	)
}

func Test_Transaction_Rollback_Error(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	// Get mocked database.
	db, mockDB, mock, err := initSQLmock("rollback_err")
	assert.NoError(
		t,
		err,
		"error creating mocked database, got = %v", err,
	)

	defer mockDB.Close()

	// Commit transaction with error.
	errTest := errlog.New("test error")
	errOnRollback := errlog.New("rollback error")

	mock.
		ExpectBegin()
	mock.
		ExpectRollback().
		WillReturnError(errOnRollback)

	err = db.Transaction(ctx, func(ctx context.Context, tx TX) error {
		return errTest
	})
	assert.EqualError(
		t,
		err,
		"tx err: test error, rb err: rollback error",
		"expected expected tx and rb error, got = %v", err,
	)
}

func Test_Get(t *testing.T) {
	t.Parallel()

	// Get mocked database.
	db, mockDB, mock, err := initSQLmock("get")
	assert.NoError(
		t,
		err,
	)

	defer mockDB.Close()

	// Execute Get with returned row.
	mock.
		ExpectQuery("SELECT").
		WillReturnRows(sqlmock.NewRows([]string{"col"}).AddRow("val"))

	var dest string

	ctx := context.Background()

	err = db.Get(ctx, &dest, "SELECT")
	assert.NoError(
		t,
		err,
	)

	// Execute Get with returned error.
	selectErr := errors.New("mock select from database error")
	mock.
		ExpectQuery("SELECT").
		WillReturnError(selectErr)

	err = db.Get(ctx, &dest, "SELECT")
	assert.EqualError(
		t,
		err,
		selectErr.Error(),
	)

	// Execute Get with returned no rows error.
	mock.
		ExpectQuery("SELECT").
		WillReturnError(sql.ErrNoRows)

	err = db.Get(ctx, &dest, "SELECT")
	assert.NoError(
		t,
		err,
	)
}

func Test_Select(t *testing.T) {
	t.Parallel()

	// Get mocked database.
	db, mockDB, mock, err := initSQLmock("select")
	assert.NoError(
		t,
		err,
	)

	defer mockDB.Close()

	// Execute Select with returned rows.
	mock.
		ExpectQuery("SELECT").
		WillReturnRows(sqlmock.NewRows([]string{"col"}).AddRow("val"))

	var dest []string

	ctx := context.Background()

	err = db.Select(ctx, &dest, "SELECT")
	assert.NoError(
		t,
		err,
	)

	// Execute Select with returned error.
	selectErr := errors.New("mock select from database error")
	mock.
		ExpectQuery("SELECT").
		WillReturnError(selectErr)

	err = db.Select(ctx, &dest, "SELECT")
	assert.EqualError(
		t,
		err,
		selectErr.Error(),
	)

	// Execute Select with returned no rows error.
	mock.
		ExpectQuery("SELECT").
		WillReturnError(sql.ErrNoRows)

	err = db.Select(ctx, &dest, "SELECT")
	assert.NoError(
		t,
		err,
	)
}

func Test_Exec(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	// Get mocked database.
	db, mockDB, mock, err := initSQLmock("exec")
	assert.NoError(
		t,
		err,
		"error creating mocked database",
	)

	defer mockDB.Close()

	// Execute select.
	mock.
		ExpectExec("SELECT").
		WillReturnResult(sqlmock.NewResult(0, 0))

	_, err = db.Exec(ctx, "SELECT")
	assert.NoError(t, err, "error starting transaction")

	// Execute select in transaction.
	mock.
		ExpectExec("SELECT").
		WillReturnResult(sqlmock.NewResult(0, 0))

	_, err = db.Exec(ctx, "SELECT")
	assert.NoError(
		t,
		err,
		"error starting transaction, got = %v", err,
	)

	// Execute select with err.
	errSelect := errlog.New("select error")
	mock.
		ExpectExec("SELECT").
		WillReturnError(errSelect)

	_, err = db.Exec(ctx, "SELECT")

	assert.EqualError(
		t,
		err,
		errSelect.Error(),
		"expected err := %v, got = %v", errSelect, err,
	)
}

func Test_NamedExec(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	// Get mocked database.
	db, mockDB, mock, err := initSQLmock("namedexec")
	assert.NoError(
		t,
		err,
		"error creating mocked database",
	)

	defer mockDB.Close()

	testArg := map[string]interface{}{
		"testArg": "testValue",
	}
	// Execute select.
	mock.
		ExpectExec("SELECT").
		WillReturnResult(sqlmock.NewResult(0, 0))

	_, err = db.NamedExec(ctx, "SELECT", testArg)
	assert.NoError(t, err, "error executing query")

	// Execute select in transaction.
	mock.
		ExpectExec("SELECT").
		WillReturnResult(sqlmock.NewResult(0, 0))

	_, err = db.NamedExec(ctx, "SELECT", testArg)
	assert.NoError(
		t,
		err,
		"expected no error, got = %v", err,
	)

	// Execute select with err.
	errSelect := errlog.New("select error")
	mock.
		ExpectExec("SELECT").
		WillReturnError(errSelect)

	_, err = db.NamedExec(ctx, "SELECT", testArg)

	assert.EqualError(
		t,
		err,
		errSelect.Error(),
		"expected err := %v, got = %v", errSelect, err,
	)
}

func Test_Ping(t *testing.T) {
	t.Parallel()

	// Get mocked database.
	db, mockDB, _, err := initSQLmock("ping")
	assert.NoError(
		t,
		err,
	)

	defer mockDB.Close()

	ctx := context.Background()

	// Execute Ping.
	err = db.Ping(ctx)
	assert.NoError(
		t,
		err,
	)
}

func Test_Close(t *testing.T) {
	t.Parallel()

	// Get mocked database.
	db, mockDB, mock, err := initSQLmock("close")
	assert.NoError(
		t,
		err,
	)

	defer mockDB.Close()

	ctx := context.Background()

	mock.ExpectClose()

	// Execute Close.
	err = db.Close(ctx)
	assert.NoError(
		t,
		err,
	)
}

func Test_Get_DB(t *testing.T) {
	t.Parallel()

	lock.Lock()
	handle = &factory{
		credentials: map[string]credentials{
			"test": {
				Driver: "mysql",
				DSN:    "mock-dns",
			},
		},
		instances: map[string]*db{
			"test": {},
		},
	}

	defer lock.Unlock()

	Get("test")

	func() {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("Expected test to panic.")
			}
		}()

		Get("mysql", "pgx")
	}()
}
