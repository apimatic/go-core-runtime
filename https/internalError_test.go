package https

import (
	"testing"
)

func TestInternalError(t *testing.T) {
	internalError := internalError{Body: "This is error body", FileInfo: "internalError.go/TestInternalError"}
	expected := "Internal Error occured at internalError.go/TestInternalError \n This is error body"
	if internalError.Error() != expected {
		t.Errorf("Failed:\nExpected: %v\nGot: %v", expected, internalError.Error())
	}
}
