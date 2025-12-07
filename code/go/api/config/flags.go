package config

import "toll/internal/errlog"

type (
	AppFlags struct {
		Svc *Service
	}
)

var flags *AppFlags

func Get() *AppFlags {
	return flags
}

func (c *AppFlags) Validate() error {
	if c == nil {
		return errlog.New("API service flags are not initialized, need to call Flags()")
	}

	if err := c.Svc.Validate(); err != nil {
		return err
	}

	return nil
}

func Flags(prefix ...string) {
	if flags != nil {
		return
	}

	if len(prefix) == 0 {
		panic("API.Flags() needs prefix on first call")
	}

	flags = &AppFlags{
		new(Service).Init(prefix...),
	}
}
