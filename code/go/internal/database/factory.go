package database

import (
	"crypto/tls"
	"crypto/x509"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/prometheus/client_golang/prometheus"

	"toll/internal/errlog"
)

var (
	certs     *x509.CertPool
	tlsConfig tls.Config
)

type (
	// credentials is a configuration struct for a database connection.
	credentials struct {
		Connector func() (*sql.DB, error)
		Driver    string
		DSN       string
		TLS       bool
	}

	// factory contains all database credentials and instances.
	factory struct {
		credentials map[string]credentials
		instances   map[string]*db
	}
)

func init() {
	var err error

	certs, err = x509.SystemCertPool()
	if err != nil {
		panic(errlog.Error(err))
	}

	tlsConfig = tls.Config{
		RootCAs: certs,

		MinVersion:         tls.VersionTLS12,
		InsecureSkipVerify: false,
	}
}

func (r *factory) Get(dbName ...string) (*db, error) {
	names := dbName
	if len(names) == 0 {
		names = []string{"db"}
	}

	if len(names) > 1 {
		return nil, fmt.Errorf("no database selected")
	}

	name := names[0]
	if value, ok := r.instances[name]; ok {
		return value, nil
	}

	crd, err := r.getCredentials(name)
	if err != nil {
		return nil, errlog.Error(err)
	}

	// Open sql connection to database.
	var conn *sqlx.DB

	switch true {
	case crd.Connector != nil:
		db, errConn := crd.Connector()
		if errConn != nil {
			return nil, errlog.Error(errConn)
		}

		conn = sqlx.NewDb(db, crd.Driver)
	default:
		conn, err = sqlx.Open(crd.Driver, crd.DSN)
		if err != nil {
			return nil, errlog.Error(err)
		}
	}

	// Define connection idle and timeout time.
	conn.SetConnMaxLifetime(5 * time.Minute)

	// Add prometheus metrics for database statistics.
	perf := NewPerformanceMetrics(name)
	stats := NewStatsMetrics(name, conn)

	prometheus.MustRegister(stats)

	r.instances[name] = &db{
		conn,
		perf,
	}

	return r.instances[name], nil
}

// getCredentials returns credentials for a given db name.
func (r *factory) getCredentials(name string) (*credentials, error) {
	if value, ok := r.credentials[name]; ok {
		return &value, nil
	}

	if err := flags.Validate(name); err != nil {
		panic(errlog.Error(err))
	}

	crd := credentials{
		Driver: flags.db[name].Driver,
		DSN:    flags.db[name].DSN,
		TLS:    flags.db[name].TLS,
	}

	switch crd.Driver {
	case "mysql":
		var opt string
		if crd.TLS {
			opt = "custom_" + name

			err := mysql.RegisterTLSConfig(opt, &tlsConfig)
			if err != nil {
				return nil, errlog.Error(err)
			}
		}

		crd.DSN = r.cleanDsn(crd.Driver, crd.DSN, opt)
	case "pgx":
		connector := func() (*sql.DB, error) {
			cfg, err := pgx.ParseConfig(crd.DSN)
			if err != nil {
				return nil, errlog.Error(err)
			}

			if crd.TLS {
				cfg.TLSConfig = &tlsConfig
			}

			connStr := stdlib.RegisterConnConfig(cfg)

			conn, err := sql.Open("pgx", connStr)
			if err != nil {
				return nil, errlog.Error(err)
			}

			return conn, nil
		}
		crd.Connector = connector
	default:
		return nil, fmt.Errorf("unrecognized database driver")
	}

	r.credentials[name] = crd

	return &crd, nil
}

func (r *factory) cleanDsn(driver string, dsn string, tls string) string {
	addOption := func(s, match, option string) string {
		if !strings.Contains(s, match) {
			s += option
		}

		return s
	}

	switch driver {
	case "mysql":
		dsn = addOption(dsn, "?", "?")
		dsn = addOption(dsn, "collation=", "&collation=utf8mb4_0900_ai_ci")
		dsn = addOption(dsn, "parseTime=", "&parseTime=true")

		if tls != "" {
			dsn = addOption(dsn, "tls=", "&tls="+tls)
		}

		dsn = strings.Replace(dsn, "?&", "?", 1)
	}

	return dsn
}
