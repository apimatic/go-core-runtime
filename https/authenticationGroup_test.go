package https

import (
	"errors"
	"github.com/apimatic/go-core-runtime/internal/assert"
	"net/http"
	"strings"
	"testing"
)

const API_KEY = "api-key"
const API_TOKEN = "api-token"
const API_KEY_MISSING_ERROR = API_KEY + " is empty!"
const API_TOKEN_MISSING_ERROR = API_TOKEN + " is empty!"

type MockHeaderCredentials struct {
	apiKey string
}

func NewMockHeaderCredentials(apiKey string) *MockHeaderCredentials {
	return &MockHeaderCredentials{apiKey: apiKey}
}

func (creds *MockHeaderCredentials) Validate() error {
	if creds.apiKey == "" {
		return errors.New(API_KEY_MISSING_ERROR)
	}

	return nil
}

func (creds *MockHeaderCredentials) Authenticator() HttpInterceptor {
	return func(req *http.Request, next HttpCallExecutor) HttpContext {
		req.Header.Set(API_KEY, creds.apiKey)
		return next(req)
	}
}

type MockQueryCredentials struct {
	apiToken string
}

func NewMockQueryCredentials(apiToken string) *MockQueryCredentials {
	return &MockQueryCredentials{apiToken: apiToken}
}

func (creds *MockQueryCredentials) Validate() error {

	if creds.apiToken == "" {
		return errors.New(API_TOKEN_MISSING_ERROR)
	}

	return nil
}

func (creds *MockQueryCredentials) Authenticator() HttpInterceptor {
	return func(req *http.Request, next HttpCallExecutor) HttpContext {
		query := req.URL.Query()
		query.Add(API_TOKEN, creds.apiToken)
		req.URL.RawQuery = query.Encode()
		return next(req)
	}
}

func AuthenticationError(errMsgs ...string) string {
	var body strings.Builder

	for _, errMsg := range errMsgs {
		body.WriteString("\n-> ")
		body.WriteString(errMsg)
	}

	authError := internalError{
		Type:     AUTHENTICATION_ERROR,
		Body:     body.String(),
		FileInfo: "callBuilder.go/Authenticate",
	}
	return authError.Error()
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

	return GetCallBuilder(ctx, "GET", "/auth", auths)
}

func TestErrorWhenUndefinedAuth(t *testing.T) {
	request := getMockCallBuilderWithAuths()
	request.Authenticate(NewAuth("authThatDoesntExist"))

	_, err := request.Call()

	assert.EqualError(t, err, AuthenticationError("authThatDoesntExist is undefined!"))
}

func TestSuccessfulCallWhenHeaderAuth(t *testing.T) {
	request := getMockCallBuilderWithAuths()
	request.Authenticate(NewAuth("header"))

	httpContext, err := request.Call()

	assert.NoError(t, err)

	header := httpContext.Request.Header

	expected := MockHeaderToken
	actual := header.Get("api-key")

	assert.Equal(t, expected, actual)
}

func TestSuccessfulCallWhenQueryAuth(t *testing.T) {
	request := getMockCallBuilderWithAuths()
	request.Authenticate(NewAuth("query"))

	httpContext, err := request.Call()

	assert.NoError(t, err)

	query := httpContext.Request.URL.Query()

	expected := MockQueryToken
	actual := query.Get("api-token")

	assert.Equal(t, expected, actual)
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

	assert.NoError(t, err)

	headerToken := httpContext.Request.Header.Get(API_KEY)
	assert.Equal(t, MockHeaderToken, headerToken)

	queryToken := httpContext.Request.URL.Query().Get(API_TOKEN)
	assert.Equal(t, MockQueryToken, queryToken)
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

	assert.NoError(t, err)

	headerToken := httpContext.Request.Header.Get(API_KEY)
	queryToken := httpContext.Request.URL.Query().Get(API_TOKEN)

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

	assert.NoError(t, err)

	headerToken := httpContext.Request.Header.Get(API_KEY)
	queryToken := httpContext.Request.URL.Query().Get(API_TOKEN)

	assert.Equal(t, "", headerToken)
	assert.Equal(t, MockQueryToken, queryToken)
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

	assert.NoError(t, err)

	headerToken := httpContext.Request.Header.Get(API_KEY)
	queryToken := httpContext.Request.URL.Query().Get(API_TOKEN)

	assert.Equal(t, "", queryToken)

	assert.Equal(t, MockHeaderToken, headerToken)
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

	assert.NoError(t, err)

	headerToken := httpContext.Request.Header.Get(API_KEY)
	queryToken := httpContext.Request.URL.Query().Get(API_TOKEN)

	assert.Equal(t, "", headerToken)
	assert.Equal(t, MockQueryToken, queryToken)
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

	assert.EqualError(t, err, AuthenticationError(API_KEY_MISSING_ERROR))
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

	assert.EqualError(t, err, AuthenticationError(API_TOKEN_MISSING_ERROR))
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

	assert.EqualError(t, err, AuthenticationError("missingQuery is undefined!"))
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

	assert.EqualError(t, err, AuthenticationError("missingHeader is undefined!"))
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

	assert.EqualError(t, err, AuthenticationError("headerMissing is undefined!", "queryMissing is undefined!"))
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

	assert.EqualError(t, err, AuthenticationError(API_KEY_MISSING_ERROR, API_TOKEN_MISSING_ERROR))
}
