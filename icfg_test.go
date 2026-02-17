package icfg

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"testing"
)

func populateManually() testConfig {
	os.Setenv("BOOL", "true")
	os.Setenv("INT", "8080")
	os.Setenv("STRING", "lol")
	os.Setenv("BOOLSLICE", "true,false,true")
	os.Setenv("INTSLICE", "1,2,3")
	os.Setenv("STRINGSLICE", "a,b,c")

	boolVar, err := strconv.ParseBool(os.Getenv("BOOL"))
	if err != nil {
		panic("could not parse expected bool")
	}
	int64Var, err := strconv.ParseInt(os.Getenv("INT"), 10, 32)
	if err != nil {
		panic("could not parse expected int ")
	}
	intVar := int(int64Var)

	stringVar := os.Getenv("STRING")
	boolSliceVar := make([]bool, 3)
	boolStrings := strings.Split(os.Getenv("BOOLSLICE"), ",")
	for i := range boolSliceVar {
		boolSliceVar[i], err = strconv.ParseBool(boolStrings[i])
		if err != nil {
			panic("could not parse expected bool in bool slice")
		}
	}

	intSliceVar := make([]int, 3)
	intStrings := strings.Split(os.Getenv("INTSLICE"), ",")
	for i := range intSliceVar {
		tempInt, err := strconv.ParseInt(intStrings[i], 10, 32)
		if err != nil {
			panic("could not parse expected int in int slice")
		}
		intSliceVar[i] = int(tempInt)
	}

	stringSliceVar := make([]string, 3)
	stringStrings := strings.Split(os.Getenv("STRINGSLICE"), ",")
	for i := range stringSliceVar {
		stringSliceVar[i] = stringStrings[i]

	}

	return testConfig{
		Bool:        boolVar,
		Int:         intVar,
		String:      stringVar,
		BoolSlice:   boolSliceVar,
		IntSlice:    intSliceVar,
		StringSlice: stringSliceVar,
	}

}

func TestFromEnv(t *testing.T) {

	expected := populateManually()

	cfg, err := FromEnv[testConfig]()
	if err != nil {
		panic(err.Error())
	}

	fmt.Println("Expected:")
	PrintConfig(&expected)
	fmt.Println("Test:")
	PrintConfig(cfg)

}
