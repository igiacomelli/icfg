package icfg

import (
	"encoding/json"
	"os"
	"reflect"
	"strconv"
	"strings"
)

const (

	// tag types
	supportedTagTypes = "env,json"

	// pkg var types
	supportedConfigTypes = "int (any), string, bool"
)

func FromJSON[T any](configPath string) (*T, error) {

	config := new(T)

	file, err := os.Open(configPath)
	if err != nil {
		return nil, err
	}
	decoder := json.NewDecoder(file)
	err = decoder.Decode(config)
	if err != nil {
		return nil, err
	}
	return config, err
}

func FromEnv[T any]() (*T, error) {

	configuratorError := newConfiguratorError(false)

	configType := reflect.TypeFor[T]()
	configPointer := new(T)

	if configType.Kind() != reflect.Struct {
		configuratorError.Malformed(notAStructError)
		return configPointer, configuratorError
	}

	configValue := reflect.ValueOf(configPointer).Elem()
	for i := 0; i < configType.NumField(); i++ {
		configField := configValue.Field(i)
		fieldName := configType.Field(i).Name
		environmentVariableName := configType.Field(i).Tag.Get("env")
		environmentVariableValue := os.Getenv(environmentVariableName)

		if environmentVariableName == "" {
			configuratorError.Append(fieldName, malformedTagError, i)
			continue
		}

		if environmentVariableValue == "" {
			configuratorError.Append(fieldName, envVarNotFound, i)
			continue
		}

		switch configField.Kind() {
		case reflect.String:
			configField.SetString(environmentVariableValue)
		case reflect.Int:
			intValue, err := strconv.Atoi(environmentVariableValue)
			if err != nil {
				configuratorError.Append(fieldName, atoiCastError, i)
			}
			configField.SetInt(int64(intValue))
		case reflect.Bool:
			boolValue, err := strconv.ParseBool(environmentVariableValue)
			if err != nil {
				configuratorError.Append(fieldName, boolCastError, i)
			}
			configField.SetBool(boolValue)
		case reflect.Slice:

			sourceVars := strings.Split(environmentVariableValue, ",")

			sliceType := configField.Type().Elem()

			if sliceType.Kind() == reflect.String {
				configField.Set(reflect.ValueOf(&sourceVars).Elem())
			} else {
				slice := setSlice(sliceType, sourceVars)

				if !slice.IsValid() {
					configuratorError.Append(fieldName, sliceParseError, i)
				} else {
					configField.Set(slice)
				}
			}

		default:
			configuratorError.Append(fieldName, unsupportedTypeError, i)
		}

	}

	if configuratorError.IsError() {
		return nil, configuratorError
	}
	return configPointer, nil
}

func getConfigTagValues(tag reflect.StructTag) []string {
	values := make([]string, 0)
	configTagValues := strings.Split(supportedTagTypes, ",")
	for _, val := range configTagValues {
		tagVal := tag.Get(val)
		if tagVal != "" {
			values = append(values, tagVal)
		}
	}
	return values
}

func setSlice(sliceType reflect.Type, sourceVars []string) reflect.Value {
	sliceLen := len(sourceVars)
	sliceValue := reflect.MakeSlice(reflect.SliceOf(sliceType), sliceLen, sliceLen)
	for i := range sliceLen {
		switch sliceType.Kind() {
		case reflect.Int:
			intValue, err := strconv.Atoi(sourceVars[i])
			if err != nil {
				return reflect.Value{}
			}
			sliceValue.Index(i).Set(reflect.ValueOf(&intValue).Elem())
		case reflect.Bool:
			boolValue, err := strconv.ParseBool(sourceVars[i])
			if err != nil {
				return reflect.Value{}
			}
			sliceValue.Index(i).Set(reflect.ValueOf(&boolValue).Elem())
		default:
		}
	}

	return sliceValue
}
