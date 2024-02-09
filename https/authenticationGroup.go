package https

import (
	"fmt"
)

const SINGLE_AUTH = "single"
const AND_AUTH = "and"
const OR_AUTH = "or"

type AuthGroup struct {
	validatedAuthInterfaces []AuthInterface
	innerAuthGroups         []AuthGroup
	authType                string
	singleAuthKey           string
	errMsg               	string
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

func (ag *AuthGroup) appendError(errMsg string) {
	if errMsg != "" {
		ag.errMsg = ag.errMsg + "\n-> " + errMsg
	}
}

func (ag *AuthGroup) validate(authInterfaces map[string]AuthInterface) {
	switch ag.authType {
	case SINGLE_AUTH:
		val, ok := authInterfaces[ag.singleAuthKey]

		if !ok {
			ag.errMsg = fmt.Sprintf("%s is undefined!", ag.singleAuthKey)
			return
		}
		if ok, err := val.Validate(); !ok {
			ag.errMsg = err.Error()
			return
		}
		ag.validatedAuthInterfaces = append(ag.validatedAuthInterfaces, val)
	case AND_AUTH, OR_AUTH:
		for _, innerAG := range ag.innerAuthGroups {
			innerAG.validate(authInterfaces)

			ag.validatedAuthInterfaces = append(ag.validatedAuthInterfaces, innerAG.validatedAuthInterfaces...)

			if ag.authType == OR_AUTH && innerAG.errMsg == "" {
				ag.errMsg = ""
				return
			}
			ag.appendError(innerAG.errMsg)
		}
	}
}