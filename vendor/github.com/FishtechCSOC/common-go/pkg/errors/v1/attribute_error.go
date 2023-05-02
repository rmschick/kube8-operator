package errors

import (
	"fmt"
	"strings"
)

type AttributeError struct {
	Err        error
	Message    string
	Attributes map[string]any
}

func CreateAttributeError(err error, message string, attributes map[string]any) *AttributeError {
	return &AttributeError{
		Err:        err,
		Message:    message,
		Attributes: attributes,
	}
}

func (err *AttributeError) Error() string {
	builder := strings.Builder{}

	builder.WriteString(fmt.Sprintf("%s\n", err.Message))

	for k, v := range err.Attributes {
		builder.WriteString(fmt.Sprintf("%s: %v\n", k, v))
	}

	builder.WriteString(fmt.Sprintf("Error: %s", err.Err.Error()))

	return builder.String()
}
