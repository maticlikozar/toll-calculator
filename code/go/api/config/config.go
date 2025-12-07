package config

import (
	"github.com/jnovack/flag"

	"toll/internal/errlog"
)

type (
	Service struct {
		Addr   string
		Domain string
		Debug  bool
	}
)

var svc *Service

func (c *Service) Validate() error {
	if c == nil {
		return nil
	}

	if c.Addr == "" {
		return errlog.New("no service addr is set, can't listen for connections")
	}

	if c.Domain == "" {
		return errlog.New("no domain is set")
	}

	return nil
}

func (*Service) Init(prefix ...string) *Service {
	if svc != nil {
		return svc
	}

	p := func(s string) string {
		return prefix[0] + "_" + s
	}

	svc = new(Service)

	flag.StringVar(
		&svc.Addr,
		p("addr"),
		":8080",
		"Listen address for server",
	)

	flag.BoolVar(
		&svc.Debug,
		p("debug"),
		false,
		"Enable/disable Debug endpoints",
	)

	flag.StringVar(
		&svc.Domain,
		p("domain"),
		"localhost",
		"Environment domain",
	)

	return svc
}
