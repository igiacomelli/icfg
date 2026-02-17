package icfg

import (
	"fmt"
)

const (

	// error types
	notAStructError      = "passed cfg variable is not a struct"
	malformedTagError    = "cfg field tag malformed"
	envVarNotFound       = "environment variable not found or blank"
	atoiCastError        = "environment variable could not be cast to int"
	boolCastError        = "environment variable could not be cast to bool"
	sliceParseError      = "environment variable could not be parsed into slice type"
	unsupportedTypeError = "cfg field is of unsupported type (supported types are" + supportedConfigTypes + ")"
	notValidError        = "value is invalid"
	zeroError            = "value is zero"

	// description variants
	configurationDescription = "Could not configure %d variables, see:\n%s"
	formattingDescription    = "Could not format %d variables, see:\n%s"
)

type configuratorError struct {
	description string
	variables   []string
	causes      []string
	indices     []int
	num         int
	malformed   bool
}

func newConfiguratorError(isFormat bool) *configuratorError {
	desc := configurationDescription
	if isFormat {
		desc = formattingDescription
	}
	return &configuratorError{
		description: desc,
		variables:   make([]string, 0),
		causes:      make([]string, 0),
		indices:     make([]int, 0),
		num:         0,
		malformed:   false,
	}
}

func (ce *configuratorError) Append(variable string, cause string, index int) {
	ce.variables = append(ce.variables, variable)
	ce.causes = append(ce.causes, cause)
	ce.indices = append(ce.indices, index)
	ce.num++
}

func (ce *configuratorError) IsError() bool {
	return ce.num > 0 || ce.malformed
}

func (ce *configuratorError) Malformed(cause string) {
	ce.description = cause
	ce.malformed = true
}

func (ce *configuratorError) Error() string {

	if ce.malformed {
		return fmt.Sprintf("Config is malformed, see: %s", ce.description)
	}

	rowTemplate := "\tname=%s\t\tindex=%d\t\tcause=%s\n"
	rows := ""
	for i := 0; i < ce.num; i++ {
		rows += fmt.Sprintf(rowTemplate, ce.variables[i], ce.indices[i], ce.causes[i])
	}
	ce.description = fmt.Sprintf(ce.description, ce.num, rows)
	return ce.description
}
