package database

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/jackc/pgx/v5"
	"github.com/jmoiron/sqlx"

	"toll/internal/errlog"
)

type (
	// DB interface defines possible operations on database.
	DB interface {
		Transaction(ctx context.Context, cb func(ctx context.Context, tx TX) error) (err error)

		Get(ctx context.Context, dest interface{}, query string, args ...interface{}) error
		Select(ctx context.Context, dest interface{}, query string, args ...interface{}) error
		Exec(ctx context.Context, query string, args ...interface{}) (res sql.Result, err error)
		NamedExec(ctx context.Context, query string, arg interface{}) (res sql.Result, err error)

		Ping(ctx context.Context) error
		Close(ctx context.Context) error
	}

	// TX interface defines possible operations on database transaction.
	TX interface {
		Exec(query string, args ...interface{}) (res sql.Result, err error)
		NamedExec(query string, arg interface{}) (res sql.Result, err error)
	}

	// db struct holds database connection object with context.
	db struct {
		*sqlx.DB

		metrics *PerformanceMetrics
	}
)

var (
	handle *factory
)

func init() {
	handle = &factory{}
	handle.credentials = make(map[string]credentials)
	handle.instances = make(map[string]*db)
}

// Get func returns database object connected to database.
func Get(name ...string) DB {
	db, err := handle.Get(name...)
	if err != nil {
		panic(errlog.Error(err))
	}

	return db
}

// Transaction func will run callback function in transaction.
func (h *db) Transaction(ctx context.Context, cb func(ctx context.Context, tx TX) error) error {
	tx, err := h.BeginTxx(ctx, &sql.TxOptions{ReadOnly: false})
	if err != nil {
		return errlog.Error(err)
	}

	err = cb(ctx, tx)

	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("tx err: %w, rb err: %w", err, rbErr)
		}

		return errlog.Error(err)
	}

	return tx.Commit()
}

func (h *db) Get(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	t1 := time.Now()

	defer func() {
		h.metrics.queryTime.WithLabelValues("Get").Observe(time.Since(t1).Seconds())
	}()

	err := h.GetContext(ctx, dest, query, args...)

	// Always ignore ErrNoRows.
	if errlog.Is(err, sql.ErrNoRows) {
		return nil
	}

	return errlog.Error(err)
}

func (h *db) Select(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	t1 := time.Now()

	defer func() {
		h.metrics.queryTime.WithLabelValues("Select").Observe(time.Since(t1).Seconds())
	}()

	err := h.SelectContext(ctx, dest, query, args...)

	// Always ignore ErrNoRows.
	if errlog.Is(err, sql.ErrNoRows) {
		return nil
	}

	return errlog.Error(err)
}

func (h *db) Exec(ctx context.Context, query string, args ...interface{}) (res sql.Result, err error) {
	t1 := time.Now()

	defer func() {
		h.metrics.queryTime.WithLabelValues("Exec").Observe(time.Since(t1).Seconds())
	}()

	return h.ExecContext(ctx, query, args...)
}

func (h *db) NamedExec(ctx context.Context, query string, arg interface{}) (res sql.Result, err error) {
	t1 := time.Now()

	defer func() {
		h.metrics.queryTime.WithLabelValues("NamedExec").Observe(time.Since(t1).Seconds())
	}()

	return h.NamedExecContext(ctx, query, arg)
}

func (h *db) Ping(ctx context.Context) error {
	return h.PingContext(ctx)
}

func (h *db) Close(_ context.Context) error {
	return h.DB.Close()
}
