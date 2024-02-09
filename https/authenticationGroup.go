package https

import (
	"fmt"
)

const SINGLE_AUTH = "single"
const AND_AUTH = "and"
const OR_AUTH = "or"

type AuthGroup struct {
	validAuthInterfaces []AuthInterface
	innerGroups         []AuthGroup
	authType            string
	singleAuthKey       string
	errMessage          string
}

func NewAuth(key string) AuthGroup {
	return AuthGroup{
		singleAuthKey: key,
		authType:      SINGLE_AUTH,
	}
}

func NewOrAuth(authGroup1, authGroup2 AuthGroup, moreAuthGroups ...AuthGroup) AuthGroup {
	return AuthGroup{
		innerGroups: append([]AuthGroup{authGroup1, authGroup2}, moreAuthGroups...),
		authType:    OR_AUTH,
	}
}

func NewAndAuth(authGroup1, authGroup2 AuthGroup, moreAuthGroups ...AuthGroup) AuthGroup {
	return AuthGroup{
		innerGroups: append([]AuthGroup{authGroup1, authGroup2}, moreAuthGroups...),
		authType:    AND_AUTH,
	}
}

func (ag *AuthGroup) appendIndentedError(errMsg string) {
	if errMsg != "" {
		ag.errMessage += "\n-> " + errMsg
	}
}

func (ag *AuthGroup) validate(authInterfaces map[string]AuthInterface) {
	switch ag.authType {
	case SINGLE_AUTH:
		val, ok := authInterfaces[ag.singleAuthKey]

		if !ok {
			ag.appendIndentedError(fmt.Sprintf("%s is undefined!", ag.singleAuthKey))
			return
		}
		if ok, err := val.Validate(); !ok {
			ag.appendIndentedError(err.Error())
			return
		}
		ag.validAuthInterfaces = append(ag.validAuthInterfaces, val)
	case AND_AUTH, OR_AUTH:
		for _, innerAG := range ag.innerGroups {
			innerAG.validate(authInterfaces)

			ag.validAuthInterfaces = append(ag.validAuthInterfaces, innerAG.validAuthInterfaces...)

			if ag.authType == OR_AUTH && innerAG.errMessage == "" {
				ag.errMessage = ""
				return
			}
			ag.errMessage += innerAG.errMessage
		}
	}
}
