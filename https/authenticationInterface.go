package https

type AuthInterface interface {
	Validate() error
	Authenticator() HttpInterceptor
}
