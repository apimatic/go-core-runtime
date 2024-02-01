package https

import (
	"errors"
	"strings"
)

const SINGLE_AUTH = "single"
const AND_AUTH = "and"
const OR_AUTH = "or"

type AuthGroup struct {
	validatedAuthInterfaces []AuthInterface
	innerAuthGroups         []AuthGroup
	authType                string
	singleAuthKey           string
}

func NewAuth(key string) AuthGroup {
	return AuthGroup{
		singleAuthKey: key,
		authType:      SINGLE_AUTH,
	}
}

func NewOrAuth(authGroups ...AuthGroup) AuthGroup {
	return AuthGroup{
		innerAuthGroups: authGroups,
		authType:        OR_AUTH,
	}
}

func NewAndAuth(authGroups ...AuthGroup) AuthGroup {
	return AuthGroup{
		innerAuthGroups: authGroups,
		authType:        AND_AUTH,
	}
}

func (ag AuthGroup) validate(authInterfaces map[string]AuthInterface, errorList []string) {

	switch ag.authType {
	case SINGLE_AUTH:
		if val, ok := authInterfaces[ag.singleAuthKey]; ok {
			if val.IsValid() {
				ag.validatedAuthInterfaces = append(ag.validatedAuthInterfaces, val)
			} else {
				errorList = append(errorList, val.ErrorMessage())
			}
		}
	case OR_AUTH:
		for _, authGroup := range ag.innerAuthGroups {

			authGroup.validate(authInterfaces, errorList)

			if len(errorList) == 0 {
				return
			}
		}
	case AND_AUTH:
		for _, authGroup := range ag.innerAuthGroups {
			authGroup.validate(authInterfaces, errorList)
		}
	}

	if len(errorList) > 0 {
		errors.New(strings.Join(errorList, "\n ==> "))
	}
}

func (ag AuthGroup) apply(cb CallBuilder) {
	for _, authI := range ag.validatedAuthInterfaces {
		cb.intercept(authI.Authenticator())
	}
}