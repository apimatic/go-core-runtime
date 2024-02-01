package https

type AuthInterface interface {
	IsValid() bool
	Authenticator() HttpInterceptor
	ErrorMessage() string
}
