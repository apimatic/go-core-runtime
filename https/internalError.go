package https

import "fmt"

const INTERNAL_ERROR = "Internal Error"
const AUTHENTICATION_ERROR = "Authentication Error"

// internalError represents a custom error type that provides additional information
// about internal errors that occur within the HTTP calling code.
type internalError struct {
	Body     string
	FileInfo string
	Type     string
}

// Error returns a formatted error string that includes the file information and the descriptive error message.
func (e internalError) Error() string {
	if e.Type == AUTHENTICATION_ERROR {
		return fmt.Sprintf("%v occured at %v due to following errors:%v", e.Type, e.FileInfo, e.Body)
	}

	return fmt.Sprintf("%v occurred at %v \n %v", INTERNAL_ERROR, e.FileInfo, e.Body)
}
