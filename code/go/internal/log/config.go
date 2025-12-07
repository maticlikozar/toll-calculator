package log

import (
	"github.com/jnovack/flag"
)

type (
	Config struct {
		Level string
	}
)

var cfg *Config

func (c *Config) Validate() error {
	if c == nil {
		return nil
	}

	return nil
}

func (*Config) Init(prefix ...string) *Config {
	if cfg != nil {
		return cfg
	}

	p := func(s string) string {
		return prefix[0] + "_" + s
	}

	cfg := new(Config)

	flag.StringVar(
		&cfg.Level,
		p("log_level"),
		"",
		"Log level (trace,debug,info,warn,error,fatal)",
	)

	return cfg
}
