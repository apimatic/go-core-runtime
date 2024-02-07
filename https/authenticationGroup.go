package https

import "errors"

const SINGLE_AUTH = "single"
const AND_AUTH = "and"
const OR_AUTH = "or"

type AuthGroup struct {
	validatedAuthInterfaces []AuthInterface
	innerAuthGroups         []AuthGroup
	authType                string
	singleAuthKey           string
	authError               error
}

func NewAuth(key string) AuthGroup {
	return AuthGroup{
		singleAuthKey: key,
		authType:      SINGLE_AUTH,
	}
}

func NewOrAuth(authGroup1, authGroup2 AuthGroup, moreAuthGroups ...AuthGroup) AuthGroup {
	return AuthGroup{
		innerAuthGroups: append([]AuthGroup{authGroup1, authGroup2}, moreAuthGroups...),
		authType:        OR_AUTH,
	}
}

func NewAndAuth(authGroup1, authGroup2 AuthGroup, moreAuthGroups ...AuthGroup) AuthGroup {
	return AuthGroup{
		innerAuthGroups: append([]AuthGroup{authGroup1, authGroup2}, moreAuthGroups...),
		authType:        AND_AUTH,
	}
}

func (ag *AuthGroup) appendError(err error) {

	if ag.authError == nil {
		ag.authError = err
	} else {
		if err != nil {
			ag.authError = errors.New(ag.authError.Error() + "\n" + err.Error())
		}
	}
}

func (ag *AuthGroup) validate(authInterfaces map[string]AuthInterface) {
	switch ag.authType {
	case SINGLE_AUTH:
		if val, ok := authInterfaces[ag.singleAuthKey]; ok {
			if val.IsValid() {
				ag.validatedAuthInterfaces = append(ag.validatedAuthInterfaces, val)
			} else {
				ag.authError = internalError{
					Type: "AuthenticationValidation Error",
					Body: val.ErrorMessage(),
					FileInfo: "authenticationGroup.go/validate",
				}
			}
		}
	case OR_AUTH, AND_AUTH:
		for _, authGroup := range ag.innerAuthGroups {
			authGroup.validate(authInterfaces)
			ag.validatedAuthInterfaces = append(ag.validatedAuthInterfaces, authGroup.validatedAuthInterfaces...)

			if ag.authType == OR_AUTH && authGroup.authError == nil {
				return
			}
			ag.appendError(authGroup.authError)
		}
	}
}

func (ag *AuthGroup) apply(cb *defaultCallBuilder) {
	cb.clientError = ag.authError
	for _, authI := range ag.validatedAuthInterfaces {
		cb.intercept(authI.Authenticator())
	}
}
