package utilities

import (
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

// Error implements the Error method for the error interface.
// It returns a string representation of the ApiError instance when used in an error context.
func (a MarshalError) Error() string {
	indent := "\n\t=>"
	switch a.innerError.(type) {
	case MarshalError:
		indent = "."
	}
	return fmt.Sprintf("%v %v %v", a.structName, indent, a.innerError)
}
