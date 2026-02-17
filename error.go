package icfg

import (
	"fmt"
)

const (

	// error types
	NotAStructError      = "passed cfg variable is not a struct"
	MalformedTagError    = "cfg field tag malformed"
	EnvVarNotFound       = "environment variable not found or blank"
	AtoiCastError        = "environment variable could not be cast to int"
	BoolCastError        = "environment variable could not be cast to bool"
	SliceParseError      = "environment variable could not be parsed into slice type"
	UnsupportedTypeError = "cfg field is of unsupported type (supported types are" + supportedConfigTypes + ")"
	NotValidError        = "value is invalid"
	ZeroError            = "value is zero"

	// description variants
	ConfigurationDescription = "Could not configure %d variables, see:\n%s"
	FormattingDescription    = "Could not format %d variables, see:\n%s"
)

type ConfiguratorError struct {
	description string
	variables   []string
	causes      []string
	indices     []int
	num         int
	malformed   bool
}

func NewConfiguratorError(isFormat bool) *ConfiguratorError {
	desc := ConfigurationDescription
	if isFormat {
		desc = FormattingDescription
	}
	return &ConfiguratorError{
		description: desc,
		variables:   make([]string, 0),
		causes:      make([]string, 0),
		indices:     make([]int, 0),
		num:         0,
		malformed:   false,
	}
}

func (ce *ConfiguratorError) Append(variable string, cause string, index int) {
	ce.variables = append(ce.variables, variable)
	ce.causes = append(ce.causes, cause)
	ce.indices = append(ce.indices, index)
	ce.num++
}

func (ce *ConfiguratorError) IsError() bool {
	return ce.num > 0 || ce.malformed
}

func (ce *ConfiguratorError) Malformed(cause string) {
	ce.description = cause
	ce.malformed = true
}

func (ce *ConfiguratorError) Error() string {

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
