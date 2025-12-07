package database

import (
	"toll/internal/errlog"
)

type (
	localFlags struct {
		db map[string]*Config
	}
)

var (
	flags *localFlags
)

func Flags(prefix ...string) {
	new(localFlags).Init(prefix...)
}

func (f *localFlags) Validate(name string) error {
	if flags == nil || flags.db == nil {
		return errlog.Errorf("database flags validation error")
	}

	cfg, ok := f.db[name]
	if !ok {
		return errlog.Errorf("no flags found for database: %v", name)
	}

	if err := cfg.Validate(); err != nil {
		return errlog.Errorf("invalid flags for database %v: %w", name, err)
	}

	return nil
}

func (f *localFlags) Init(prefix ...string) *localFlags {
	if flags != nil {
		return flags
	}

	flags = &localFlags{
		db: make(map[string]*Config),
	}

	flags.db["db"] = new(Config).Init("db", prefix[0])

	return flags
}
