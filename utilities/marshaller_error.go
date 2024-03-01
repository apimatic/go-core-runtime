package utilities

import (
	"fmt"
)

type MarshallerError struct {
	structName string
	innerError error
}

func NewMarshallerError(structName string, err error) MarshallerError {
	return MarshallerError{
		structName: structName,
		innerError: err,
	}
}

// Error implements the Error method for the error interface.
// It returns a string representation of the ApiError instance when used in an error context.
func (a MarshallerError) Error() string {
	indent := "\n\t=>"
	switch a.innerError.(type) {
	case MarshallerError:
		indent = "."
	}
	return fmt.Sprintf("%v %v %v", a.structName, indent, a.innerError)
}
