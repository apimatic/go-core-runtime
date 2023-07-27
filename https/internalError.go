package https

import "fmt"

// internalError represents a custom error type that provides additional information
// about internal errors that occur within the HTTP calling code.
type internalError struct {
	Body     string
	FileInfo string
}

// Error returns a formatted error string that includes the file information and the descriptive error message.
func (e internalError) Error() string {
	return fmt.Sprintf("Internal Error occured at %v \n %v", e.FileInfo, e.Body)
}
