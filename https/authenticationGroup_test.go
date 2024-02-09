package https

import (
	"errors"
	"net/http"
	"testing"
)

type MockHeaderCredentials struct {
	apiKey string
}

func NewMockHeaderCredentials(apiKey string) *MockHeaderCredentials {
	return &MockHeaderCredentials{apiKey: apiKey}
}

func (creds *MockHeaderCredentials) Validate() (bool, error) {
	return creds.apiKey != "", errors.New(creds.ErrorMessage())
}

func (creds *MockHeaderCredentials) Authenticator() HttpInterceptor {
	return func(req *http.Request, next HttpCallExecutor) HttpContext {
		req.Header.Set("api-key", creds.apiKey)
		return next(req)
	}
}

func (creds *MockHeaderCredentials) ErrorMessage() string {

	if creds.apiKey == "" {
		return "api-key is empty!"
	}

	return "Error:  MockHeaderCredentials"
}

type MockQueryCredentials struct {
	apiToken string
}

func NewMockQueryCredentials(apiToken string) *MockQueryCredentials {
	return &MockQueryCredentials{apiToken: apiToken}
}

func (creds *MockQueryCredentials) Validate() (bool, error) {
	return creds.apiToken != "", errors.New(creds.ErrorMessage())
}

func (creds *MockQueryCredentials) Authenticator() HttpInterceptor {
	return func(req *http.Request, next HttpCallExecutor) HttpContext {
		query := req.URL.Query()
		query.Add("api-token", creds.apiToken)
		req.URL.RawQuery = query.Encode()
		return next(req)
	}
}

func (creds *MockQueryCredentials) ErrorMessage() string {

	if creds.apiToken == "" {
		return "api-token is empty!"
	}

	return "Error: MockQueryCredentials"
}

const MockHeaderToken = "1234"
const MockQueryToken = "abcd"

func getMockCallBuilderWithAuths() CallBuilder {
	auths := map[string]AuthInterface{
		"header":         NewMockHeaderCredentials(MockHeaderToken),
		"headerEmptyVal": NewMockHeaderCredentials(""),
		"query":          NewMockQueryCredentials(MockQueryToken),
		"queryEmptyVal":  NewMockQueryCredentials(""),
	}

	// TODO: ctx is a global variable. What should be done with it?
	return GetCallBuilder(ctx, "GET", "/auth", auths)
}

func TestErrorWhenUndefinedAuth(t *testing.T) {
	request := getMockCallBuilderWithAuths()
	request.Authenticate(NewAuth("authThatDoesntExist"))

	_, err := request.Call()

	if err == nil {
		t.Fatalf("Expected an error.")
	}

	expected := "authThatDoesntExist is undefined!"

	if err.Error() != expected {
		t.Errorf("Expected error message: %q. \nGot %q.", expected, err.Error())
	}
}

func TestSuccessfulCallWhenHeaderAuth(t *testing.T) {
	request := getMockCallBuilderWithAuths()
	request.Authenticate(NewAuth("header"))

	httpContext, err := request.Call()

	if err != nil {
		t.Fatalf("Unexpected Error: %s", err.Error())
	}

	header := httpContext.Request.Header

	expected := MockHeaderToken

	if actual := header.Get("api-key"); actual != expected {
		t.Errorf("Expected %q, got %q", expected, actual)
	}
}

func TestSuccessfulCallWhenQueryAuth(t *testing.T) {
	request := getMockCallBuilderWithAuths()
	request.Authenticate(NewAuth("query"))

	httpContext, err := request.Call()

	if err != nil {
		t.Fatalf("Unexpected Error: %s", err.Error())
	}

	query := httpContext.Request.URL.Query()

	expected := MockQueryToken

	if actual := query.Get("api-token"); actual != expected {
		t.Errorf("Expected %q, got %q", expected, actual)
	}
}

func TestSuccessfulCallWhenHeaderAndQueryAuth(t *testing.T) {
	request := getMockCallBuilderWithAuths()
	request.Authenticate(
		NewAndAuth(
			NewAuth("header"),
			NewAuth("query"),
		),
	)

	httpContext, err := request.Call()

	if err != nil {
		t.Fatalf("Unexpected Error: %s", err.Error())
	}

	headerToken := httpContext.Request.Header.Get("api-key")
	if headerToken != MockHeaderToken {
		t.Errorf("Expected header param 'api-token' with value %q, got %q", MockHeaderToken, headerToken)
	}

	queryToken := httpContext.Request.URL.Query().Get("api-token")
	if queryToken != MockQueryToken {
		t.Errorf("Expected query param 'api-key' with value %q, got %q", MockQueryToken, queryToken)
	}
}

func TestSuccessfulCallWhenHeaderOrQueryAuth(t *testing.T) {
	request := getMockCallBuilderWithAuths()
	request.Authenticate(
		NewOrAuth(
			NewAuth("header"),
			NewAuth("query"),
		),
	)

	httpContext, err := request.Call()

	if err != nil {
		t.Fatalf("Unexpected Error: %s", err.Error())
	}

	headerToken := httpContext.Request.Header.Get("api-key")
	queryToken := httpContext.Request.URL.Query().Get("api-token")

	if headerToken != MockHeaderToken && queryToken != MockQueryToken {
		t.Errorf("Expected either header param 'api-key' with value %q"+
			" or query param 'api-token' with value %q. Got neither.",
			MockHeaderToken, MockQueryToken)
	}
}

func TestSuccessfulCallWhenEmptyHeaderOrQueryAuth(t *testing.T) {
	request := getMockCallBuilderWithAuths()
	request.Authenticate(
		NewOrAuth(
			NewAuth("headerEmptyVal"),
			NewAuth("query"),
		),
	)

	httpContext, err := request.Call()

	if err != nil {
		t.Fatalf("Unexpected Error: %s", err.Error())
	}

	headerToken := httpContext.Request.Header.Get("api-key")
	queryToken := httpContext.Request.URL.Query().Get("api-token")

	if headerToken != "" {
		t.Errorf("Expected no header param. Got %q.", headerToken)
	}

	if queryToken != MockQueryToken {
		t.Errorf("Expected query param 'api-token' with value %q. Got %q.", MockQueryToken, queryToken)
	}
}

func TestSuccessfulCallWhenHeaderOrMissingQueryAuth(t *testing.T) {
	request := getMockCallBuilderWithAuths()
	request.Authenticate(
		NewOrAuth(
			NewAuth("header"),
			NewAuth("queryMissing"),
		),
	)

	httpContext, err := request.Call()

	if err != nil {
		t.Fatalf("Unexpected Error: %s", err.Error())
	}

	headerToken := httpContext.Request.Header.Get("api-key")
	queryToken := httpContext.Request.URL.Query().Get("api-token")

	if queryToken != "" {
		t.Errorf("Expected no query param. Got %q.", queryToken)
	}

	if headerToken != MockHeaderToken {
		t.Errorf("Expected header param 'api-key' with value %q. Got %q.", MockHeaderToken, headerToken)
	}
}

func TestSuccessfulCallWhenMissingHeaderOrQueryAuth(t *testing.T) {
	request := getMockCallBuilderWithAuths()
	request.Authenticate(
		NewOrAuth(
			NewAuth("headerMissing"),
			NewAuth("query"),
		),
	)

	httpContext, err := request.Call()

	if err != nil {
		t.Fatalf("Unexpected Error: %s", err.Error())
	}

	headerToken := httpContext.Request.Header.Get("api-key")
	queryToken := httpContext.Request.URL.Query().Get("api-token")

	if headerToken != "" {
		t.Errorf("Expected no header param. Got %q.", headerToken)
	}

	if queryToken != MockQueryToken {
		t.Errorf("Expected query param 'api-token' with value %q. Got %q.", MockQueryToken, queryToken)
	}
}

func TestErrorWhenHeaderWithEmptyValueAndQueryAuth(t *testing.T) {
	request := getMockCallBuilderWithAuths()
	request.Authenticate(
		NewAndAuth(
			NewAuth("headerEmptyVal"),
			NewAuth("query"),
		),
	)

	_, err := request.Call()

	if err == nil {
		t.Fatalf("Expected an error.")
	}

	expected := "api-key is empty!"

	if err.Error() != expected {
		t.Errorf("Expected error message: %q. Got: \n%s.", expected, err.Error())
	}
}

func TestErrorWhenHeaderAndQueryWithEmptyValueAuth(t *testing.T) {
	request := getMockCallBuilderWithAuths()
	request.Authenticate(
		NewAndAuth(
			NewAuth("header"),
			NewAuth("queryEmptyVal"),
		),
	)

	_, err := request.Call()

	if err == nil {
		t.Fatalf("Expected an error.")
	}

	expected := "api-token is empty!"

	if err.Error() != expected {
		t.Errorf("Expected error message: %q. Got: \n%s.", expected, err.Error())
	}
}

func TestErrorWhenHeaderAndMissingQueryAuth(t *testing.T) {
	request := getMockCallBuilderWithAuths()
	request.Authenticate(
		NewAndAuth(
			NewAuth("header"),
			NewAuth("missingQuery"),
		),
	)

	_, err := request.Call()

	if err == nil {
		t.Fatalf("Expected an error.")
	}

	expected := "missingQuery is undefined!"

	if err.Error() != expected {
		t.Errorf("Expected error message: %q. Got: \n%s.", expected, err.Error())
	}
}

func TestErrorWhenMissingHeaderAndQueryAuth(t *testing.T) {
	request := getMockCallBuilderWithAuths()
	request.Authenticate(
		NewAndAuth(
			NewAuth("missingHeader"),
			NewAuth("query"),
		),
	)

	_, err := request.Call()

	if err == nil {
		t.Fatalf("Expected an error.")
	}

	expected := "missingHeader is undefined!"

	if err.Error() != expected {
		t.Errorf("Expected error message: %q. Got: \n%s.", expected, err.Error())
	}
}

func TestErrorWhenHeaderOrQueryAuthBothAreMissing(t *testing.T) {
	request := getMockCallBuilderWithAuths()
	request.Authenticate(
		NewOrAuth(
			NewAuth("headerMissing"),
			NewAuth("queryMissing"),
		),
	)

	_, err := request.Call()

	if err == nil {
		t.Fatalf("Expected an error.")
	}

	expected := "at least one valid auth credential must be provided"

	if err.Error() != expected {
		t.Errorf("Expected error message: %q. Got: \n%s.", expected, err.Error())
	}
}

func TestErrorWhenHeaderOrQueryAuthBothAreEmpty(t *testing.T) {
	request := getMockCallBuilderWithAuths()
	request.Authenticate(
		NewOrAuth(
			NewAuth("headerEmptyVal"),
			NewAuth("queryEmptyVal"),
		),
	)

	_, err := request.Call()

	if err == nil {
		t.Fatalf("Expected an error.")
	}

	expected := "at least one valid auth credential must be provided"

	if err.Error() != expected {
		t.Errorf("Expected error message: %q. Got: \n%s.", expected, err.Error())
	}
}
