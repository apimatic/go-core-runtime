package https

import "fmt"

type internalError struct {
	Body     string
	FileInfo string
}

func (e internalError) Error() string {
	return fmt.Sprintf("Internal Error occured at %v \n %v", e.FileInfo, e.Body)
}
