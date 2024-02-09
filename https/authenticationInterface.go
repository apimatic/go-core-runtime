package https

type AuthInterface interface {
	Validate() (bool, error)
	Authenticator() HttpInterceptor
}
