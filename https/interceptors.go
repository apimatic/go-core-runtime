package https

import (
	"net/http"
)

type HttpCallExecutor func(request *http.Request) HttpContext

type HttpInterceptor func(request *http.Request, next HttpCallExecutor) HttpContext

func PassThroughInterceptor(
	req *http.Request,
	next HttpCallExecutor,
) HttpContext {
	return next(req)
}

func CallHttpInterceptors(
	interceptors []HttpInterceptor,
	client HttpCallExecutor,
) HttpCallExecutor {
	var next HttpCallExecutor = client
	for index := len(interceptors) - 1; index >= 0; index-- {
		current := interceptors[index]
		last := next
		next = func(req *http.Request) HttpContext {
			return current(req, last)
		}
	}
	return next
}
