package icfg

// example/test icfg

type ServerConfig struct {
	Bool        bool     `env:"BOOL"`
	Int         int      `env:"INT"`
	String      string   `env:"STRING"`
	BoolSlice   []bool   `env:"BOOLSLICE"`
	IntSlice    []int    `env:"INTSLICE"`
	StringSlice []string `env:"STRINGSLICE"`
}
