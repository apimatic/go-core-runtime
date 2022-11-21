package apiError

import (
	"fmt"
)

// This is the base struct for all exceptions that represent an error response from the server.
type CustomError struct {
	LineNumber 	int			`json:"LineNumber"`
	FileName	string		`json:"FileName"`
	Message		string		`json:"Message"`
	InnerError    	error		`json:"Error"`
}

// Constructor for CustomError.
func NewCustomError(
	lineNumber int,
	fileName string,
	message string,
	innerError error) *CustomError {
	return &CustomError{
		LineNumber: 	lineNumber,
		FileName:       fileName,
		Message:		message,
		InnerError:		innerError,
	}
}

// Implementing the Error method for the error interface.
func (c *CustomError) Error() string {
	return fmt.Sprintf("CustomError occurred in File : %v, at LineNumber : %v, having ErrorMessage : %v, and its InnerError is : %v",
	 c.FileName, c.LineNumber, c.Message, c.InnerError.Error())
}
