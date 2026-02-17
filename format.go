package icfg

import (
	"fmt"
	"reflect"
)

func formatConfig(s interface{}) (string, error) {

	formattingError := NewConfiguratorError(true)

	formattedString := ""
	structType := reflect.TypeOf(s).Elem()
	structValue := reflect.ValueOf(s).Elem()

	if structType.Kind() != reflect.Struct {
		formattingError.Malformed(NotAStructError)
		return "", formattingError
	}

	for i := 0; i < structType.NumField(); i++ {
		for _, value := range getConfigTagValues(structType.Field(i).Tag) {
			structField := reflect.Value{}

			if !structValue.IsValid() {
				formattingError.Malformed(NotValidError)
				return formattedString, formattingError
			}

			if !structValue.IsZero() {
				structField = structValue.Field(i)
			}
			formattedString += fmt.Sprintf("%s = %v\n", value, structField)
		}
	}
	return formattedString, nil
}

/*
ConfigString returns a formatted string displaying the names and values of all fields in a config struct referred to by
a pointer.

This function does not work on generic structs, the reflection logic within is only for supported config types.
*/
func ConfigString(config interface{}) string {
	formattedString, err := formatConfig(config)
	if err != nil {
		return formattedString + err.Error()
	}
	return formattedString
}

/*
PrintConfig prints a string returned by ConfigString.

This function does not work on generic structs, the reflection logic within is only for supported config types.
*/
func PrintConfig(config interface{}) {

	formattedString, err := formatConfig(config)
	if err != nil {
		fmt.Println(formattedString + err.Error())
	} else {
		fmt.Println(formattedString)
	}

}
