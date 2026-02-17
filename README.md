ICFG (Ian Config) 

This is my lightweight configuration module I use for my Go projects (those in need of configuration, that is). 

It supports a typical JSON-to-struct configurator as well as an env-to-struct functionality which uses reflection to fill a config struct with environment variables in much the same way the JSON decoder does. 

Example Usage : 

```go
// sample configuration struct, the tag "env:" indicates an environment variable and the tag's name indicates that 
//variable's name 
type TestConfig struct {
	Bool        bool     `env:"BOOL"`
	Int         int      `env:"INT"`
	String      string   `env:"STRING"`
	BoolSlice   []bool   `env:"BOOLSLICE"`
	IntSlice    []int    `env:"INTSLICE"`
	StringSlice []string `env:"STRINGSLICE"`
}
```

```go 
// sample code simulating corresponding environment variables and initializing our config struct with them
func main() {
	os.Setenv("BOOL", "true")
	os.Setenv("INT", "8080")
	os.Setenv("STRING", "lol")
	os.Setenv("BOOLSLICE", "true,false,true")
	os.Setenv("INTSLICE", "1,2,3")
	os.Setenv("STRINGSLICE", "a,b,c")
	cfg, err := config.EnvConfigurate[TestConfig]()
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println("Struct values:")
	config.PrintConfig(cfg)
}
```

```terminaloutput
Struct values: 
BOOL = true
INT = 8080
STRING = lol
BOOLSLICE = [true false true]
INTSLICE = [1 2 3]
STRINGSLICE = [a b c]
```