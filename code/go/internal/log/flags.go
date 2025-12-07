package log

type (
	localFlags struct {
		log *Config
	}
)

var flags *localFlags

func Flags(prefix ...string) {
	new(localFlags).Init(prefix...)
}

func (f *localFlags) Init(prefix ...string) *localFlags {
	if flags != nil {
		return flags
	}

	flags = &localFlags{
		new(Config).Init(prefix...),
	}

	return flags
}
