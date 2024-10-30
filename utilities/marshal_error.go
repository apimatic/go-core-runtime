package utilities

import (
	"errors"
	"fmt"
)

type MarshalError struct {
	structName string
	innerError error
}

func NewMarshalError(structName string, err error) MarshalError {
	return MarshalError{
		structName: structName,
		innerError: err,
	}
}

// Error implements the Error function for the error interface.
// It returns a string representation of the MarshalError instance when used in an error context.
func (a MarshalError) Error() string {
	indent := "\n\t=>"
	var marshalError MarshalError
	if errors.As(a.innerError, &marshalError) {
		indent = "."
	}
	return fmt.Sprintf("%v %v %v", a.structName, indent, a.innerError)
}
