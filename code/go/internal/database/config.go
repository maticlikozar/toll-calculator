package database

import (
	"github.com/jnovack/flag"

	"toll/internal/errlog"
)

type (
	Config struct {
		Driver string
		DSN    string
		TLS    bool
	}
)

var cfgs map[string]*Config

func (c *Config) Validate() error {
	if c == nil {
		return nil
	}

	if c.Driver == "" {
		return errlog.New("database driver not set")
	}

	if c.DSN == "" {
		return errlog.New("database DSN not set")
	}

	return nil
}

func (*Config) Init(name string, prefix string) *Config {
	if cfgs == nil {
		cfgs = make(map[string]*Config)
	}

	// Specific database config.
	cfg, ok := cfgs[name]
	if ok && cfg != nil {
		return cfg
	}

	if !ok {
		cfg = new(Config)
		cfgs[name] = cfg
	}

	p := func(s string) string {
		return prefix + "_" + name + "_" + s
	}

	flag.StringVar(
		&cfg.Driver,
		p("driver"),
		"mysql",
		"Database driver name (required)",
	)

	flag.StringVar(
		&cfg.DSN,
		p("dsn"),
		"",
		"Database data source name (required)",
	)

	flag.BoolVar(
		&cfg.TLS,
		p("tls"),
		false,
		"Enable TLS traffic encryption",
	)

	return cfg
}
