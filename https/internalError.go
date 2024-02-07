package https

import "fmt"

// internalError represents a custom error type that provides additional information
// about internal errors that occur within the HTTP calling code.
type internalError struct {
	Body     string
	FileInfo string
	Type 	 string
}

// Error returns a formatted error string that includes the file information and the descriptive error message.
func (e internalError) Error() string {
	if e.Type == "" {
		e.Type = "Internal Error"
	}

	return fmt.Sprintf("%v occured at %v \n %v", e.Type, e.FileInfo, e.Body)
}
